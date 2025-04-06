package controllers

import (
	"blog/config"
	"blog/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取用户通知列表
func GetNotifications(c *gin.Context) {
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
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 获取筛选参数
	isRead := c.Query("isRead")
	notificationType := c.Query("type")

	// 构建查询
	query := config.DB.Model(&models.Notification{}).Where("user_id = ?", userModel.ID)

	// 应用筛选条件
	if isRead != "" {
		if isRead == "true" {
			query = query.Where("is_read = ?", true)
		} else if isRead == "false" {
			query = query.Where("is_read = ?", false)
		}
	}

	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 获取未读通知数
	var unreadCount int64
	config.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userModel.ID, false).Count(&unreadCount)

	// 查询通知列表
	var notifications []models.Notification
	query.Preload("Sender").Preload("Post").Preload("Comment").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&notifications)

	c.JSON(http.StatusOK, gin.H{
		"data":        notifications,
		"total":       total,
		"unreadCount": unreadCount,
		"page":        page,
		"size":        pageSize,
	})
}

// 将通知标记为已读
func MarkNotificationAsRead(c *gin.Context) {
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

	// 获取通知ID
	notificationID := c.Param("id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "通知ID不能为空"})
		return
	}

	// 检查通知是否存在且属于当前用户
	var notification models.Notification
	if err := config.DB.Where("id = ? AND user_id = ?", notificationID, userModel.ID).First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在或不属于当前用户"})
		return
	}

	// 更新为已读
	notification.IsRead = true
	if err := config.DB.Save(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新通知状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "通知已标记为已读",
	})
}

// 将所有通知标记为已读
func MarkAllNotificationsAsRead(c *gin.Context) {
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

	// 更新所有未读通知为已读
	if err := config.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userModel.ID, false).
		Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新通知状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "所有通知已标记为已读",
	})
}

// 删除通知
func DeleteNotification(c *gin.Context) {
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

	// 获取通知ID
	notificationID := c.Param("id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "通知ID不能为空"})
		return
	}

	// 检查通知是否存在且属于当前用户
	var notification models.Notification
	if err := config.DB.Where("id = ? AND user_id = ?", notificationID, userModel.ID).First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在或不属于当前用户"})
		return
	}

	// 删除通知
	if err := config.DB.Delete(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除通知失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "通知已删除",
	})
}

// 创建新通知（系统内部使用）
func CreateNotification(userID uint, notificationType string, content string, senderID *uint, postID *uint, commentID *uint, redirectURL string) error {
	notification := models.Notification{
		Type:        notificationType,
		Content:     content,
		UserID:      userID,
		SenderID:    senderID,
		PostID:      postID,
		CommentID:   commentID,
		RedirectURL: redirectURL,
		IsRead:      false,
	}

	return config.DB.Create(&notification).Error
}

// 创建评论通知
func CreateCommentNotification(postAuthorID uint, commenterID uint, commentID uint, postID uint, postTitle string) error {
	// 如果评论者就是文章作者，不发送通知
	if postAuthorID == commenterID {
		return nil
	}

	// 查询评论者信息
	var commenter models.User
	if err := config.DB.First(&commenter, commenterID).Error; err != nil {
		return err
	}

	content := fmt.Sprintf("%s 评论了你的文章《%s》", commenter.Username, postTitle)
	redirectURL := fmt.Sprintf("/posts/%d#comment-%d", postID, commentID)

	return CreateNotification(
		postAuthorID,
		models.NotificationTypeComment,
		content,
		&commenterID,
		&postID,
		&commentID,
		redirectURL,
	)
}

// 创建回复通知
func CreateReplyNotification(parentCommentAuthorID uint, replierID uint, commentID uint, postID uint, postTitle string) error {
	// 如果回复者就是原评论作者，不发送通知
	if parentCommentAuthorID == replierID {
		return nil
	}

	// 查询回复者信息
	var replier models.User
	if err := config.DB.First(&replier, replierID).Error; err != nil {
		return err
	}

	content := fmt.Sprintf("%s 回复了你在《%s》中的评论", replier.Username, postTitle)
	redirectURL := fmt.Sprintf("/posts/%d#comment-%d", postID, commentID)

	return CreateNotification(
		parentCommentAuthorID,
		models.NotificationTypeReply,
		content,
		&replierID,
		&postID,
		&commentID,
		redirectURL,
	)
}
