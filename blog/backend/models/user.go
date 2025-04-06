package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 用户模型
type User struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Username      string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email         string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Password      string         `json:"-" gorm:"size:100;not null"` // 不在JSON中暴露密码
	Avatar        string         `json:"avatar" gorm:"size:255"`
	Bio           string         `json:"bio" gorm:"size:500"`
	Website       string         `json:"website" gorm:"size:255"`
	Github        string         `json:"github" gorm:"size:100"`
	Twitter       string         `json:"twitter" gorm:"size:100"`
	ThemeSettings string         `json:"themeSettings" gorm:"type:text"`
	Role          string         `json:"role" gorm:"size:20;default:'user'"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// 创建用户前的钩子 - 用于密码加密
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// 更新密码
func (u *User) UpdatePassword(newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
