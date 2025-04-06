package models

import (
	"time"

	"gorm.io/gorm"
)

// 通知类型
const (
	NotificationTypeComment = "comment"
	NotificationTypeReply   = "reply"
	NotificationTypeLike    = "like"
	NotificationTypeSystem  = "system"
)

// 通知模型
type Notification struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Type        string         `json:"type" gorm:"size:20;not null"` // comment, reply, like, system
	Content     string         `json:"content" gorm:"type:text;not null"`
	IsRead      bool           `json:"isRead" gorm:"default:false"`
	UserID      uint           `json:"userId" gorm:"not null;index"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	SenderID    *uint          `json:"senderId"` // 发送者用户ID，可以为空（系统通知）
	Sender      *User          `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
	PostID      *uint          `json:"postId"` // 关联的文章ID，可以为空
	Post        *Post          `json:"post,omitempty" gorm:"foreignKey:PostID"`
	CommentID   *uint          `json:"commentId"` // 关联的评论ID，可以为空
	Comment     *Comment       `json:"comment,omitempty" gorm:"foreignKey:CommentID"`
	RedirectURL string         `json:"redirectUrl" gorm:"size:255"` // 点击通知跳转链接
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
