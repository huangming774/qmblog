package controllers

import (
	"blog/config"
	"blog/models"
	"blog/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 创建文章请求
type CreatePostRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Summary string   `json:"summary"`
	Cover   string   `json:"cover"`
	Status  string   `json:"status" binding:"required,oneof=draft published"`
	Tags    []string `json:"tags"`
}

// 更新文章请求
type UpdatePostRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Summary string   `json:"summary"`
	Cover   string   `json:"cover"`
	Status  string   `json:"status" binding:"omitempty,oneof=draft published"`
	Tags    []string `json:"tags"`
}

// 获取所有文章
func GetPosts(c *gin.Context) {
	// 查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.DefaultQuery("status", "published")
	tag := c.Query("tag")

	offset := (page - 1) * pageSize
	var posts []models.Post
	var total int64
	query := config.DB.Model(&models.Post{})

	// 只获取已发布文章或者当前用户的草稿
	if status == "all" {
		// 如果用户是管理员，允许查看所有文章
		user, exists := c.Get("user")
		if exists {
			userModel, ok := user.(models.User)
			if ok && userModel.Role == "admin" {
				// 管理员可以查看所有状态
			} else {
				// 非管理员只能查看自己的草稿和所有已发布文章
				query = query.Where("status = 'published' OR (status = 'draft' AND user_id = ?)", userModel.ID)
			}
		} else {
			// 未登录用户只能查看已发布文章
			query = query.Where("status = 'published'")
		}
	} else {
		query = query.Where("status = ?", status)
	}

	// 按标签过滤
	if tag != "" {
		query = query.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
			Joins("JOIN tags ON post_tags.tag_id = tags.id").
			Where("tags.name = ?", tag)
	}

	// 使用WaitGroup等待两个并发操作完成：1.获取总数 2.获取分页数据
	var wg sync.WaitGroup
	wg.Add(2)

	var queryErr error

	// 并发获取总数
	go func() {
		defer wg.Done()
		queryClone := query
		if err := queryClone.Count(&total).Error; err != nil {
			queryErr = err
		}
	}()

	// 并发获取分页数据
	go func() {
		defer wg.Done()
		queryClone := query
		if err := queryClone.Preload("User").Preload("Tags").
			Order("created_at DESC").
			Offset(offset).
			Limit(pageSize).
			Find(&posts).Error; err != nil {
			queryErr = err
		}
	}()

	// 等待所有goroutine完成
	wg.Wait()

	// 处理可能出现的错误
	if queryErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败: " + queryErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  posts,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// 获取单篇文章
func GetPost(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()

	// 定义Redis缓存键
	postCacheKey := fmt.Sprintf("post:%s", id)
	viewCacheKey := fmt.Sprintf("post_view:%s", id)

	// 尝试从Redis缓存获取文章
	var post models.Post
	postData, err := config.Redis.HGetAll(ctx, postCacheKey).Result()

	// 如果缓存存在且不为空
	if err == nil && len(postData) > 0 {
		// 从缓存提取基本字段
		post.ID = uint(utils.StringToUint(postData["id"]))
		post.Title = postData["title"]
		post.Content = postData["content"]
		post.Summary = postData["summary"]
		post.Cover = postData["cover"]
		post.Status = postData["status"]
		post.UserID = utils.StringToUint(postData["user_id"])
		post.ViewCount = utils.StringToUint(postData["view_count"])
		post.CreatedAt, _ = time.Parse(time.RFC3339, postData["created_at"])
		post.UpdatedAt, _ = time.Parse(time.RFC3339, postData["updated_at"])

		// 使用goroutine异步增加阅读计数，不阻塞主流程
		go func() {
			config.Redis.Incr(ctx, viewCacheKey)

			// 检查是否需要同步到数据库
			count, err := config.Redis.Get(ctx, viewCacheKey).Int()
			if err == nil && count%10 == 0 {
				// 每10次同步到数据库
				config.DB.Model(&models.Post{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + ?", 10))
				// 重置计数器
				config.Redis.Set(ctx, viewCacheKey, 0, time.Hour*24)
				// 更新缓存中的阅读量
				config.Redis.HIncrBy(ctx, postCacheKey, "view_count", 10)
			}
		}()
	} else {
		// 缓存不存在，从数据库获取
		result := config.DB.Preload("User").Preload("Tags").Preload("Comments").Preload("Comments.User").First(&post, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}

		// 使用goroutine异步设置缓存，不阻塞主流程
		go func(p models.Post) {
			// 将文章存入Redis缓存
			cacheData := map[string]interface{}{
				"id":         fmt.Sprintf("%d", p.ID),
				"title":      p.Title,
				"content":    p.Content,
				"summary":    p.Summary,
				"cover":      p.Cover,
				"status":     p.Status,
				"user_id":    fmt.Sprintf("%d", p.UserID),
				"view_count": fmt.Sprintf("%d", p.ViewCount),
				"created_at": p.CreatedAt.Format(time.RFC3339),
				"updated_at": p.UpdatedAt.Format(time.RFC3339),
			}

			// 设置缓存，过期时间24小时
			config.Redis.HMSet(ctx, postCacheKey, cacheData)
			config.Redis.Expire(ctx, postCacheKey, time.Hour*24)

			// 设置阅读计数器
			config.Redis.Set(ctx, viewCacheKey, 0, time.Hour*24)
			config.Redis.Incr(ctx, viewCacheKey)
		}(post)
	}

	c.JSON(http.StatusOK, post)
}

// 创建文章
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证失败: " + err.Error()})
		return
	}

	// 获取当前用户
	user, _ := c.Get("user")
	userModel := user.(models.User)

	// 创建文章
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		Summary: req.Summary,
		Cover:   req.Cover,
		Status:  req.Status,
		UserID:  userModel.ID,
	}

	// 开始事务
	tx := config.DB.Begin()

	// 保存文章
	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败: " + err.Error()})
		return
	}

	// 处理标签
	if len(req.Tags) > 0 {
		for _, tagName := range req.Tags {
			var tag models.Tag
			// 查找或创建标签
			if tx.Where("name = ?", tagName).First(&tag).Error != nil {
				tag = models.Tag{Name: tagName}
				if err := tx.Create(&tag).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败: " + err.Error()})
					return
				}
			}
			// 添加标签到文章
			if err := tx.Model(&post).Association("Tags").Append(&tag); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "添加标签失败: " + err.Error()})
				return
			}
		}
	}

	// 提交事务
	tx.Commit()

	c.JSON(http.StatusCreated, post)
}

// 更新文章
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证失败: " + err.Error()})
		return
	}

	// 查找文章
	var post models.Post
	if result := config.DB.First(&post, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限
	user, _ := c.Get("user")
	userModel := user.(models.User)
	if post.UserID != userModel.ID && userModel.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改此文章"})
		return
	}

	// 开始事务
	tx := config.DB.Begin()

	// 更新文章
	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.Cover != "" {
		updates["cover"] = req.Cover
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := tx.Model(&post).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败: " + err.Error()})
		return
	}

	// 处理标签
	if len(req.Tags) > 0 {
		// 清除旧标签
		if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "清除标签失败: " + err.Error()})
			return
		}

		// 添加新标签
		for _, tagName := range req.Tags {
			var tag models.Tag
			// 查找或创建标签
			if tx.Where("name = ?", tagName).First(&tag).Error != nil {
				tag = models.Tag{Name: tagName}
				if err := tx.Create(&tag).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败: " + err.Error()})
					return
				}
			}
			// 添加标签到文章
			if err := tx.Model(&post).Association("Tags").Append(&tag); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "添加标签失败: " + err.Error()})
				return
			}
		}
	}

	// 提交事务
	tx.Commit()

	// 重新加载文章
	tx.Preload("User").Preload("Tags").First(&post, id)

	// 使用goroutine异步执行缓存删除，不阻塞主流程
	ctx := context.Background()
	postCacheKey := fmt.Sprintf("post:%s", id)
	go func() {
		// 更新完成后，删除缓存，强制下次请求重新加载
		config.Redis.Del(ctx, postCacheKey)
	}()

	c.JSON(http.StatusOK, post)
}

// 删除文章
func DeletePost(c *gin.Context) {
	id := c.Param("id")

	// 查找文章
	var post models.Post
	if result := config.DB.First(&post, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限
	user, _ := c.Get("user")
	userModel := user.(models.User)
	if post.UserID != userModel.ID && userModel.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此文章"})
		return
	}

	// 删除文章（软删除）
	if result := config.DB.Delete(&post); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败: " + result.Error.Error()})
		return
	}

	// 使用goroutine异步执行缓存删除，不阻塞主流程
	ctx := context.Background()
	postCacheKey := fmt.Sprintf("post:%s", id)
	viewCacheKey := fmt.Sprintf("post_view:%s", id)
	go func() {
		// 删除缓存
		config.Redis.Del(ctx, postCacheKey, viewCacheKey)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "文章已删除"})
}
