package controllers

import (
	"blog/config"
	"blog/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建评论请求
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required"`
	ParentID *uint  `json:"parentId"`
}

// 更新评论请求
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// 获取文章评论
func GetComments(c *gin.Context) {
	postID := c.Param("id")

	var comments []models.Comment
	result := config.DB.Where("post_id = ? AND parent_id IS NULL", postID).
		Preload("User").
		Preload("Replies").
		Preload("Replies.User").
		Order("created_at DESC").
		Find(&comments)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// 创建评论
func CreateComment(c *gin.Context) {
	postID := c.Param("id")
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证失败: " + err.Error()})
		return
	}

	// 检查文章是否存在
	var post models.Post
	if result := config.DB.First(&post, postID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 获取当前用户
	user, _ := c.Get("user")
	userModel := user.(models.User)

	// 如果是回复评论，检查父评论是否存在
	var parentComment models.Comment
	var parentAuthorID uint

	if req.ParentID != nil {
		if result := config.DB.First(&parentComment, req.ParentID); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "父评论不存在"})
			return
		}
		// 确保父评论属于当前文章
		if parentComment.PostID != post.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的父评论ID"})
			return
		}
		parentAuthorID = parentComment.UserID
	}

	// 创建评论
	comment := models.Comment{
		Content:  req.Content,
		UserID:   userModel.ID,
		PostID:   post.ID,
		ParentID: req.ParentID,
	}

	if result := config.DB.Create(&comment).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return
	}

	// 重新加载评论以获取关联数据
	config.DB.Preload("User").First(&comment, comment.ID)

	// 发送通知
	if req.ParentID == nil {
		// 如果是对文章的评论，通知文章作者
		CreateCommentNotification(post.UserID, userModel.ID, comment.ID, post.ID, post.Title)
	} else {
		// 如果是对评论的回复，通知原评论作者
		CreateReplyNotification(parentAuthorID, userModel.ID, comment.ID, post.ID, post.Title)
	}

	c.JSON(http.StatusCreated, comment)
}

// 更新评论
func UpdateComment(c *gin.Context) {
	commentID := c.Param("id")

	// 验证请求
	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证失败: " + err.Error()})
		return
	}

	// 查找评论
	var comment models.Comment
	if result := config.DB.First(&comment, commentID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 检查权限
	user, _ := c.Get("user")
	userModel := user.(models.User)
	if comment.UserID != userModel.ID && userModel.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权编辑此评论"})
		return
	}

	// 更新评论
	comment.Content = req.Content
	if err := config.DB.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "评论已更新",
		"data":    comment,
	})
}

// 删除评论
func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")

	// 查找评论
	var comment models.Comment
	if result := config.DB.First(&comment, commentID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 检查权限
	user, _ := c.Get("user")
	userModel := user.(models.User)
	if comment.UserID != userModel.ID && userModel.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此评论"})
		return
	}

	// 删除评论及其所有回复
	tx := config.DB.Begin()

	// 删除回复
	if err := tx.Where("parent_id = ?", comment.ID).Delete(&models.Comment{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除回复失败: " + err.Error()})
		return
	}

	// 删除相关通知
	if err := tx.Where("comment_id = ?", comment.ID).Delete(&models.Notification{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除相关通知失败"})
		return
	}

	// 删除主评论
	if err := tx.Delete(&comment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除评论失败: " + err.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "评论已删除",
	})
}
