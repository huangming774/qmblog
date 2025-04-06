package config

import (
	"blog/models"
	"log"
)

// 执行数据库迁移
func RunMigrations() {
	// 自动迁移模型到数据库
	err := DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Tag{},
		&models.Comment{},
		&models.Category{},
		&models.Favorite{},
		&models.Notification{},
	)

	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库迁移完成")

	// 创建管理员账户（如果不存在）
	createAdminUser()
}

// 创建管理员账户
func createAdminUser() {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

	if count == 0 {
		admin := models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "admin123", // 在生产环境中使用更强的密码
			Role:     "admin",
			Avatar:   "https://ui-avatars.com/api/?name=Admin&background=random",
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Printf("创建管理员账户失败: %v", err)
			return
		}

		log.Println("创建管理员账户成功")
	}
}
