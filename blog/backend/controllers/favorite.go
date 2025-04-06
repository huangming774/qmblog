package controllers

import (
	"blog/config"
	"blog/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 添加收藏
func AddFavorite(c *gin.Context) {
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

	// 获取文章ID
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID不能为空"})
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查是否已收藏
	var existingFavorite models.Favorite
	result := config.DB.Where("user_id = ? AND post_id = ?", userModel.ID, post.ID).First(&existingFavorite)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已经收藏过该文章"})
		return
	}

	// 创建收藏记录
	favorite := models.Favorite{
		UserID: userModel.ID,
		PostID: post.ID,
	}

	if err := config.DB.Create(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加收藏失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "收藏成功",
		"data":    favorite,
	})
}

// 取消收藏
func RemoveFavorite(c *gin.Context) {
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

	// 获取收藏ID
	favoriteID := c.Param("id")
	if favoriteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "收藏ID不能为空"})
		return
	}

	// 检查收藏记录是否存在且属于当前用户
	var favorite models.Favorite
	if err := config.DB.Where("id = ? AND user_id = ?", favoriteID, userModel.ID).First(&favorite).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "收藏记录不存在或不属于当前用户"})
		return
	}

	// 删除收藏记录
	if err := config.DB.Delete(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消收藏失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "取消收藏成功",
	})
}

// 检查是否已收藏
func CheckFavorite(c *gin.Context) {
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

	// 获取文章ID
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID不能为空"})
		return
	}

	// 检查是否已收藏
	var favorite models.Favorite
	result := config.DB.Where("user_id = ? AND post_id = ?", userModel.ID, postID).First(&favorite)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"favorited": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"favorited": true,
		"favorite":  favorite,
	})
}
