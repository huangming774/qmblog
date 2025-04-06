package models

import (
	"time"

	"gorm.io/gorm"
)

// 文章模型
type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Summary   string         `json:"summary" gorm:"size:500"`
	Cover     string         `json:"cover" gorm:"size:255"`
	Status    string         `json:"status" gorm:"size:20;default:'draft'"` // draft, published
	UserID    uint           `json:"userId" gorm:"not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Tags      []Tag          `json:"tags" gorm:"many2many:post_tags;"`
	Comments  []Comment      `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	ViewCount uint           `json:"viewCount" gorm:"default:0"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// 标签模型
type Tag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"uniqueIndex;size:50;not null"`
	Posts     []Post         `json:"posts,omitempty" gorm:"many2many:post_tags;"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
