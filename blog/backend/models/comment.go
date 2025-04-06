package models

import (
	"time"

	"gorm.io/gorm"
)

// 评论模型
type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	UserID    uint           `json:"userId" gorm:"not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	PostID    uint           `json:"postId" gorm:"not null"`
	Post      Post           `json:"-" gorm:"foreignKey:PostID"`
	ParentID  *uint          `json:"parentId" gorm:"default:null"` // 父评论ID，用于回复功能
	Parent    *Comment       `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies   []Comment      `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
