package controllers

import (
	"blog/config"
	"blog/models"
	"blog/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 更新用户资料请求
type UpdateProfileRequest struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Website  string `json:"website"`
	Github   string `json:"github"`
	Twitter  string `json:"twitter"`
}

// 更新用户密码请求
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=NewPassword"`
}

// 更新主题设置请求
type UpdateThemeRequest struct {
	DarkMode   *bool  `json:"darkMode"`
	ThemeColor string `json:"themeColor"`
	FontSize   string `json:"fontSize"`
}

// 获取当前用户资料
func GetUserProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 获取完整的用户信息
	var fullUser models.User
	if err := config.DB.First(&fullUser, userModel.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 获取用户文章数量
	var postCount int64
	config.DB.Model(&models.Post{}).Where("user_id = ?", fullUser.ID).Count(&postCount)

	// 获取用户评论数量
	var commentCount int64
	config.DB.Model(&models.Comment{}).Where("user_id = ?", fullUser.ID).Count(&commentCount)

	c.JSON(http.StatusOK, gin.H{
		"id":           fullUser.ID,
		"username":     fullUser.Username,
		"email":        fullUser.Email,
		"avatar":       fullUser.Avatar,
		"bio":          fullUser.Bio,
		"website":      fullUser.Website,
		"github":       fullUser.Github,
		"twitter":      fullUser.Twitter,
		"role":         fullUser.Role,
		"createdAt":    fullUser.CreatedAt,
		"postCount":    postCount,
		"commentCount": commentCount,
	})
}

// 更新用户资料
func UpdateUserProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 获取表单数据
	username := c.PostForm("username")
	bio := c.PostForm("bio")
	website := c.PostForm("website")
	github := c.PostForm("github")
	twitter := c.PostForm("twitter")

	// 获取当前用户记录
	var dbUser models.User
	if err := config.DB.First(&dbUser, userModel.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 检查用户名是否已被其他用户使用
	if username != "" && username != dbUser.Username {
		var existingUser models.User
		if result := config.DB.Where("username = ? AND id != ?", username, dbUser.ID).First(&existingUser); result.Error == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "该用户名已被使用"})
			return
		}
		dbUser.Username = username
	}

	// 更新其他字段
	if bio != "" {
		dbUser.Bio = bio
	}
	if website != "" {
		dbUser.Website = website
	}
	if github != "" {
		dbUser.Github = github
	}
	if twitter != "" {
		dbUser.Twitter = twitter
	}

	// 头像上传处理
	file, err := c.FormFile("avatar")
	if err == nil {
		// 创建上传目录
		uploadDir := "uploads/avatars"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
			return
		}

		// 生成唯一文件名
		fileExt := filepath.Ext(file.Filename)
		fileName := uuid.NewString() + fileExt
		filePath := filepath.Join(uploadDir, fileName)

		// 保存文件
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存头像失败"})
			return
		}

		// 更新头像URL
		dbUser.Avatar = "/" + filePath
	}

	// 保存更新
	if err := config.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户资料失败"})
		return
	}

	// 更新令牌中的用户名
	if username != "" && username != userModel.Username {
		token, err := utils.GenerateToken(dbUser.ID, dbUser.Username, dbUser.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成新令牌失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "资料更新成功",
			"token":   token,
			"user": gin.H{
				"id":       dbUser.ID,
				"username": dbUser.Username,
				"email":    dbUser.Email,
				"avatar":   dbUser.Avatar,
				"role":     dbUser.Role,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "资料更新成功",
		"user": gin.H{
			"id":       dbUser.ID,
			"username": dbUser.Username,
			"email":    dbUser.Email,
			"avatar":   dbUser.Avatar,
			"role":     dbUser.Role,
		},
	})
}

// 更新用户密码
func UpdateUserPassword(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证失败: " + err.Error()})
		return
	}

	// 获取当前用户记录
	var dbUser models.User
	if err := config.DB.First(&dbUser, userModel.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 验证当前密码
	if !dbUser.CheckPassword(req.CurrentPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "当前密码不正确"})
		return
	}

	// 更新密码
	if err := dbUser.UpdatePassword(req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
		return
	}

	// 保存更新
	if err := config.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
		return
	}

	// 重新生成令牌
	token, err := utils.GenerateToken(dbUser.ID, dbUser.Username, dbUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成新令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "密码更新成功",
		"token":   token,
	})
}

// 更新主题设置
func UpdateThemeSettings(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	var req UpdateThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证失败: " + err.Error()})
		return
	}

	// 获取当前用户记录
	var dbUser models.User
	if err := config.DB.First(&dbUser, userModel.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 构建主题设置JSON
	themeSettings := make(map[string]interface{})

	// 从现有设置解析
	if dbUser.ThemeSettings != "" {
		if err := utils.ParseJSON(dbUser.ThemeSettings, &themeSettings); err != nil {
			// 如果解析失败，创建新的设置
			themeSettings = make(map[string]interface{})
		}
	}

	// 更新设置
	if req.DarkMode != nil {
		themeSettings["darkMode"] = *req.DarkMode
	}
	if req.ThemeColor != "" {
		themeSettings["themeColor"] = req.ThemeColor
	}
	if req.FontSize != "" {
		themeSettings["fontSize"] = req.FontSize
	}

	// 保存设置
	themeSettingsJSON, err := utils.ToJSON(themeSettings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "序列化主题设置失败"})
		return
	}

	dbUser.ThemeSettings = themeSettingsJSON

	// 保存更新
	if err := config.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新主题设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "主题设置更新成功",
		"settings": themeSettings,
	})
}

// 获取用户文章列表
func GetUserPosts(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
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

	// 获取筛选参数
	status := c.Query("status")
	keyword := c.Query("keyword")

	// 构建查询
	query := config.DB.Model(&models.Post{}).Where("user_id = ?", userModel.ID)

	// 应用筛选条件
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		search := fmt.Sprintf("%%%s%%", keyword)
		query = query.Where("title LIKE ? OR content LIKE ?", search, search)
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询文章列表
	var posts []models.Post
	query.Preload("Tags").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"data":  posts,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// 获取用户评论列表
func GetUserComments(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
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

	// 获取筛选参数
	postID := c.Query("postId")
	keyword := c.Query("keyword")

	// 构建查询
	query := config.DB.Model(&models.Comment{}).Where("user_id = ?", userModel.ID)

	// 应用筛选条件
	if postID != "" {
		query = query.Where("post_id = ?", postID)
	}
	if keyword != "" {
		search := fmt.Sprintf("%%%s%%", keyword)
		query = query.Where("content LIKE ?", search)
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询评论列表
	var comments []models.Comment
	query.Preload("Post").Preload("Replies").Preload("Parent").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments)

	// 扩展评论数据，添加文章标题
	type CommentWithPostTitle struct {
		models.Comment
		PostTitle string `json:"postTitle"`
		ReplyTo   string `json:"replyTo,omitempty"`
	}

	result := make([]CommentWithPostTitle, 0, len(comments))
	for _, comment := range comments {
		item := CommentWithPostTitle{
			Comment:   comment,
			PostTitle: comment.Post.Title,
		}
		if comment.Parent != nil {
			// 获取父评论作者
			var parentUser models.User
			config.DB.First(&parentUser, comment.Parent.UserID)
			item.ReplyTo = parentUser.Username + ": " + comment.Parent.Content
			if len(item.ReplyTo) > 50 {
				item.ReplyTo = item.ReplyTo[:50] + "..."
			}
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  result,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// 获取用户收藏列表
func GetUserFavorites(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "12"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 12
	}
	offset := (page - 1) * pageSize

	// 获取筛选参数
	categoryID := c.Query("categoryId")
	tagID := c.Query("tagId")
	keyword := c.Query("keyword")

	// 基础查询
	query := config.DB.Model(&models.Favorite{}).Where("user_id = ?", userModel.ID)

	// 应用筛选 - 这需要连接查询
	if categoryID != "" || tagID != "" || keyword != "" {
		query = query.Joins("JOIN posts ON favorites.post_id = posts.id")

		if categoryID != "" {
			query = query.Joins("JOIN post_categories ON posts.id = post_categories.post_id").
				Where("post_categories.category_id = ?", categoryID)
		}

		if tagID != "" {
			query = query.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
				Where("post_tags.tag_id = ?", tagID)
		}

		if keyword != "" {
			search := fmt.Sprintf("%%%s%%", keyword)
			query = query.Where("posts.title LIKE ? OR posts.content LIKE ?", search, search)
		}
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询收藏列表
	var favorites []models.Favorite
	query.Preload("Post").Preload("Post.User").Preload("Post.Tags").
		Order("favorites.created_at DESC").
		Offset(offset).Limit(pageSize).Find(&favorites)

	c.JSON(http.StatusOK, gin.H{
		"data":  favorites,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}
