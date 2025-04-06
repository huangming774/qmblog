package controllers

import (
	"blog/config"
	"blog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取所有标签
func GetTags(c *gin.Context) {
	var tags []models.Tag

	// 查询所有标签
	if err := config.DB.Order("name ASC").Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tags,
	})
}

// 获取热门标签（按文章数量排序）
func GetPopularTags(c *gin.Context) {
	// 获取查询数量，默认10个
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// 查询标签和关联的文章数量
	type TagCount struct {
		models.Tag
		PostCount int `json:"postCount"`
	}

	var tagCounts []TagCount

	// 使用原始SQL查询，按文章数量排序
	err := config.DB.Raw(`
		SELECT t.*, COUNT(pt.post_id) as post_count 
		FROM tags t 
		LEFT JOIN post_tags pt ON t.id = pt.tag_id 
		GROUP BY t.id 
		ORDER BY post_count DESC, t.name ASC 
		LIMIT ?
	`, limit).Scan(&tagCounts).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取热门标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tagCounts,
	})
}

// 获取单个标签及其文章
func GetTag(c *gin.Context) {
	tagID := c.Param("id")
	if tagID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标签ID不能为空"})
		return
	}

	var tag models.Tag
	if err := config.DB.First(&tag, tagID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 查询标签下的文章
	var posts []models.Post
	var total int64

	query := config.DB.Model(&models.Post{}).
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Where("post_tags.tag_id = ? AND posts.status = ?", tag.ID, "published")

	// 获取总数
	query.Count(&total)

	// 获取文章列表
	query.Preload("User").Preload("Tags").
		Order("posts.created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"tag": tag,
		"posts": gin.H{
			"data":  posts,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// 添加标签（管理员）
func AddTag(c *gin.Context) {
	// 验证请求
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证失败: " + err.Error()})
		return
	}

	// 检查标签名是否已存在
	var existingTag models.Tag
	if result := config.DB.Where("name = ?", req.Name).First(&existingTag); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该标签名已存在"})
		return
	}

	// 创建新标签
	tag := models.Tag{
		Name: req.Name,
	}

	if err := config.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "标签创建成功",
		"data":    tag,
	})
}

// 更新标签（管理员）
func UpdateTag(c *gin.Context) {
	tagID := c.Param("id")
	if tagID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标签ID不能为空"})
		return
	}

	// 验证请求
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证失败: " + err.Error()})
		return
	}

	// 检查标签是否存在
	var tag models.Tag
	if err := config.DB.First(&tag, tagID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	// 检查是否有其他标签使用相同名称
	var existingTag models.Tag
	if result := config.DB.Where("name = ? AND id != ?", req.Name, tag.ID).First(&existingTag); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该标签名已被其他标签使用"})
		return
	}

	// 更新标签
	tag.Name = req.Name

	if err := config.DB.Save(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "标签更新成功",
		"data":    tag,
	})
}

// 删除标签（管理员）
func DeleteTag(c *gin.Context) {
	tagID := c.Param("id")
	if tagID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标签ID不能为空"})
		return
	}

	// 检查标签是否存在
	var tag models.Tag
	if err := config.DB.First(&tag, tagID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	// 先删除文章与标签的关联
	if err := config.DB.Exec("DELETE FROM post_tags WHERE tag_id = ?", tag.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除标签关联失败"})
		return
	}

	// 删除标签
	if err := config.DB.Delete(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "标签删除成功",
	})
}
