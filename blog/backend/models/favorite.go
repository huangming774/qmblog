package models

import (
	"time"

	"gorm.io/gorm"
)

// 收藏模型
type Favorite struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"userId" gorm:"not null;index"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	PostID    uint           `json:"postId" gorm:"not null;index"`
	Post      Post           `json:"post" gorm:"foreignKey:PostID"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
