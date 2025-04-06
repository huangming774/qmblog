package controllers

import (
	"blog/config"
	"blog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取所有分类
func GetCategories(c *gin.Context) {
	var categories []models.Category

	// 查询所有分类
	if err := config.DB.Order("name ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

// 获取单个分类及其文章
func GetCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类ID不能为空"})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
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

	// 查询分类下的文章
	var posts []models.Post
	var total int64

	query := config.DB.Model(&models.Post{}).
		Joins("JOIN post_categories ON posts.id = post_categories.post_id").
		Where("post_categories.category_id = ? AND posts.status = ?", category.ID, "published")

	// 获取总数
	query.Count(&total)

	// 获取文章列表
	query.Preload("User").Preload("Tags").
		Order("posts.created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"category": category,
		"posts": gin.H{
			"data":  posts,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// 添加分类（管理员）
func AddCategory(c *gin.Context) {
	// 验证请求
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证失败: " + err.Error()})
		return
	}

	// 检查分类名是否已存在
	var existingCategory models.Category
	if result := config.DB.Where("name = ?", req.Name).First(&existingCategory); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类名已存在"})
		return
	}

	// 创建新分类
	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建分类失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "分类创建成功",
		"data":    category,
	})
}

// 更新分类（管理员）
func UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类ID不能为空"})
		return
	}

	// 验证请求
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证失败: " + err.Error()})
		return
	}

	// 检查分类是否存在
	var category models.Category
	if err := config.DB.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	// 检查是否有其他分类使用相同名称
	var existingCategory models.Category
	if result := config.DB.Where("name = ? AND id != ?", req.Name, category.ID).First(&existingCategory); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类名已被其他分类使用"})
		return
	}

	// 更新分类
	category.Name = req.Name
	category.Description = req.Description

	if err := config.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新分类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "分类更新成功",
		"data":    category,
	})
}

// 删除分类（管理员）
func DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类ID不能为空"})
		return
	}

	// 检查分类是否存在
	var category models.Category
	if err := config.DB.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	// 检查是否有文章关联该分类
	var count int64
	config.DB.Table("post_categories").Where("category_id = ?", category.ID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类下有关联的文章，无法删除"})
		return
	}

	// 删除分类
	if err := config.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除分类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "分类删除成功",
	})
}
