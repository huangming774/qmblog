package main

import (
	"blog/config"
	"blog/models"
	"fmt"
	"math/rand"
	"time"
)

// 初始化数据库测试数据
func main() {
	// 初始化数据库连接
	config.InitDB()
	defer config.CloseDB()

	// 创建测试用户
	createTestUsers()

	// 创建测试文章和标签
	createTestPosts()

	fmt.Println("数据库初始化完成！")
}

// 创建测试用户
func createTestUsers() {
	var count int64
	config.DB.Model(&models.User{}).Count(&count)

	if count < 3 {
		users := []models.User{
			{
				Username: "admin",
				Email:    "admin@example.com",
				Password: "admin123",
				Role:     "admin",
				Avatar:   "https://ui-avatars.com/api/?name=Admin&background=random",
			},
			{
				Username: "demo",
				Email:    "demo@example.com",
				Password: "demo123",
				Role:     "user",
				Avatar:   "https://ui-avatars.com/api/?name=Demo&background=random",
			},
			{
				Username: "test",
				Email:    "test@example.com",
				Password: "test123",
				Role:     "user",
				Avatar:   "https://ui-avatars.com/api/?name=Test&background=random",
			},
		}

		for _, user := range users {
			config.DB.Where("email = ?", user.Email).FirstOrCreate(&user)
		}

		fmt.Println("创建测试用户成功！")
	}
}

// 创建测试文章和标签
func createTestPosts() {
	var count int64
	config.DB.Model(&models.Post{}).Count(&count)

	if count < 5 {
		// 获取用户
		var users []models.User
		config.DB.Find(&users)

		if len(users) == 0 {
			fmt.Println("没有可用用户，无法创建文章")
			return
		}

		// 创建标签
		tags := []models.Tag{
			{Name: "Go"},
			{Name: "Vue"},
			{Name: "React"},
			{Name: "PostgreSQL"},
			{Name: "Redis"},
		}

		for i := range tags {
			config.DB.Where("name = ?", tags[i].Name).FirstOrCreate(&tags[i])
		}

		// 创建文章
		rand.Seed(time.Now().UnixNano())
		posts := []models.Post{
			{
				Title:     "Go语言入门教程",
				Content:   "Go（又称Golang）是Google开发的一种静态强类型、编译型、并发型，并具有垃圾回收功能的编程语言。\n\nGo语言的特点是：\n1. 简洁、快速、安全\n2. 并行、有趣、开源\n3. 内存管理、数组安全、编译迅速\n\n本教程将带你从零开始学习Go语言...",
				Summary:   "这是一篇关于Go语言入门的教程，从零开始学习Go编程。",
				Cover:     "https://via.placeholder.com/800x400?text=Go+Programming",
				Status:    "published",
				UserID:    users[rand.Intn(len(users))].ID,
				ViewCount: uint(rand.Intn(1000)),
			},
			{
				Title:     "Vue3最佳实践指南",
				Content:   "Vue.js是一个用于构建用户界面的渐进式框架。与其它大型框架不同的是，Vue被设计为可以自底向上逐层应用。\n\nVue 3带来了许多新特性：\n1. Composition API\n2. Teleport组件\n3. Fragments\n4. 更好的TypeScript支持\n\n本文将介绍Vue 3的最佳实践...",
				Summary:   "探讨Vue 3的新特性与最佳实践，助你构建更高效的前端应用。",
				Cover:     "https://via.placeholder.com/800x400?text=Vue3+Best+Practices",
				Status:    "published",
				UserID:    users[rand.Intn(len(users))].ID,
				ViewCount: uint(rand.Intn(1000)),
			},
			{
				Title:     "使用React构建现代化网站",
				Content:   "React是一个用于构建用户界面的JavaScript库。React使创建交互式UI变得轻而易举。\n\n为你的应用程序的每个状态设计简单的视图，当你的数据改变时，React 将有效地更新和正确的渲染组件。\n\n声明式视图使你的代码更可预测，更容易调试...",
				Summary:   "学习如何使用React生态系统构建现代化、高性能的web应用。",
				Cover:     "https://via.placeholder.com/800x400?text=Modern+React",
				Status:    "published",
				UserID:    users[rand.Intn(len(users))].ID,
				ViewCount: uint(rand.Intn(1000)),
			},
			{
				Title:     "PostgreSQL高级技巧",
				Content:   "PostgreSQL是一个功能强大的开源对象关系数据库系统，它使用和扩展了SQL语言，并结合了许多安全存储和扩展最复杂数据工作负载的功能。\n\n本文将分享一些PostgreSQL的高级技巧...",
				Summary:   "掌握PostgreSQL高级特性和优化技巧，提升数据库性能。",
				Cover:     "https://via.placeholder.com/800x400?text=PostgreSQL+Advanced",
				Status:    "draft",
				UserID:    users[rand.Intn(len(users))].ID,
				ViewCount: uint(rand.Intn(100)),
			},
			{
				Title:     "Redis缓存策略详解",
				Content:   "Redis是一个开源（BSD许可）的，内存中的数据结构存储系统，它可以用作数据库、缓存和消息中间件。\n\n本文将详细介绍Redis的缓存策略...",
				Summary:   "深入理解Redis缓存机制，设计最佳缓存策略，优化应用性能。",
				Cover:     "https://via.placeholder.com/800x400?text=Redis+Caching",
				Status:    "published",
				UserID:    users[rand.Intn(len(users))].ID,
				ViewCount: uint(rand.Intn(1000)),
			},
		}

		// 保存文章并关联标签
		for i := range posts {
			result := config.DB.Create(&posts[i])
			if result.Error != nil {
				fmt.Printf("创建文章失败: %v\n", result.Error)
				continue
			}

			// 随机关联1-3个标签
			numTags := rand.Intn(3) + 1
			selectedTags := make(map[int]bool)

			for j := 0; j < numTags; j++ {
				tagIdx := rand.Intn(len(tags))
				if !selectedTags[tagIdx] {
					selectedTags[tagIdx] = true
					config.DB.Model(&posts[i]).Association("Tags").Append(&tags[tagIdx])
				}
			}

			// 添加一些评论
			numComments := rand.Intn(5) + 1
			for k := 0; k < numComments; k++ {
				comment := models.Comment{
					Content: []string{
						"感谢分享这篇文章，很有帮助！",
						"内容非常详尽，学到了很多。",
						"有没有更多相关的资料推荐？",
						"这个观点我有不同看法...",
						"写得太好了，期待更多此类内容。",
					}[rand.Intn(5)],
					UserID: users[rand.Intn(len(users))].ID,
					PostID: posts[i].ID,
				}
				config.DB.Create(&comment)
			}
		}

		fmt.Println("创建测试文章和标签成功！")
	}
}
