package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

// 初始化PostgreSQL数据库连接
func InitDB() {
	dsn := "host=localhost user=postgres password=postgres dbname=blog port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	log.Println("成功连接PostgreSQL数据库")
}

// 初始化Redis连接
func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 无密码
		DB:       0,  // 默认DB
	})

	ctx := context.Background()
	_, err := Redis.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("连接Redis失败: %v", err)
	}

	log.Println("成功连接Redis数据库")
}

// 关闭数据库连接
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Printf("获取SQL DB实例失败: %v", err)
		return
	}
	sqlDB.Close()

	if Redis != nil {
		Redis.Close()
	}
}
