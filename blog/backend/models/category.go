package models

import (
	"time"

	"gorm.io/gorm"
)

// 分类模型
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;size:50;not null"`
	Description string         `json:"description" gorm:"size:200"`
	Posts       []Post         `json:"posts,omitempty" gorm:"many2many:post_categories;"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
