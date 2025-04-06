package main

import (
	"blog/config"
	"blog/controllers"
	"blog/middlewares"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	config.InitDB()
	defer config.CloseDB()

	// 初始化Redis连接
	config.InitRedis()

	// 执行数据库迁移
	config.RunMigrations()

	// 初始化Gin框架
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 初始化路由
	setupRoutes(r)

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 文章相关路由
		v1.GET("/posts", controllers.GetPosts)
		v1.GET("/posts/:id", controllers.GetPost)
		v1.POST("/posts", middlewares.AuthMiddleware(), controllers.CreatePost)
		v1.PUT("/posts/:id", middlewares.AuthMiddleware(), controllers.UpdatePost)
		v1.DELETE("/posts/:id", middlewares.AuthMiddleware(), controllers.DeletePost)

		// 用户认证相关路由
		v1.POST("/auth/login", controllers.Login)
		v1.POST("/auth/register", controllers.Register)

		// 评论相关路由
		v1.GET("/posts/:id/comments", controllers.GetComments)
		v1.POST("/posts/:id/comments", middlewares.AuthMiddleware(), controllers.CreateComment)
		v1.PUT("/comments/:id", middlewares.AuthMiddleware(), controllers.UpdateComment)
		v1.DELETE("/comments/:id", middlewares.AuthMiddleware(), controllers.DeleteComment)

		// 用户资料相关路由
		v1.GET("/user/profile", middlewares.AuthMiddleware(), controllers.GetUserProfile)
		v1.POST("/user/profile", middlewares.AuthMiddleware(), controllers.UpdateUserProfile)
		v1.PUT("/user/password", middlewares.AuthMiddleware(), controllers.UpdateUserPassword)
		v1.PUT("/user/theme", middlewares.AuthMiddleware(), controllers.UpdateThemeSettings)

		// 用户文章与评论
		v1.GET("/user/posts", middlewares.AuthMiddleware(), controllers.GetUserPosts)
		v1.GET("/user/comments", middlewares.AuthMiddleware(), controllers.GetUserComments)

		// 用户收藏
		v1.GET("/user/favorites", middlewares.AuthMiddleware(), controllers.GetUserFavorites)
		v1.POST("/posts/:id/favorite", middlewares.AuthMiddleware(), controllers.AddFavorite)
		v1.GET("/posts/:id/favorite", middlewares.AuthMiddleware(), controllers.CheckFavorite)
		v1.DELETE("/favorites/:id", middlewares.AuthMiddleware(), controllers.RemoveFavorite)

		// 用户通知
		v1.GET("/user/notifications", middlewares.AuthMiddleware(), controllers.GetNotifications)
		v1.PUT("/user/notifications/:id/read", middlewares.AuthMiddleware(), controllers.MarkNotificationAsRead)
		v1.PUT("/user/notifications/read-all", middlewares.AuthMiddleware(), controllers.MarkAllNotificationsAsRead)
		v1.DELETE("/user/notifications/:id", middlewares.AuthMiddleware(), controllers.DeleteNotification)

		// 分类和标签
		v1.GET("/categories", controllers.GetCategories)
		v1.GET("/categories/:id", controllers.GetCategory)
		v1.POST("/categories", middlewares.AuthMiddleware(), middlewares.AdminMiddleware(), controllers.AddCategory)
		v1.PUT("/categories/:id", middlewares.AuthMiddleware(), middlewares.AdminMiddleware(), controllers.UpdateCategory)
		v1.DELETE("/categories/:id", middlewares.AuthMiddleware(), middlewares.AdminMiddleware(), controllers.DeleteCategory)

		v1.GET("/tags", controllers.GetTags)
		v1.GET("/tags/popular", controllers.GetPopularTags)
		v1.GET("/tags/:id", controllers.GetTag)
		v1.POST("/tags", middlewares.AuthMiddleware(), middlewares.AdminMiddleware(), controllers.AddTag)
		v1.PUT("/tags/:id", middlewares.AuthMiddleware(), middlewares.AdminMiddleware(), controllers.UpdateTag)
		v1.DELETE("/tags/:id", middlewares.AuthMiddleware(), middlewares.AdminMiddleware(), controllers.DeleteTag)
	}
}

// 这些函数将在各自的控制器文件中实现
func GetPosts(c *gin.Context)         {}
func GetPost(c *gin.Context)          {}
func CreatePost(c *gin.Context)       {}
func UpdatePost(c *gin.Context)       {}
func DeletePost(c *gin.Context)       {}
func Login(c *gin.Context)            {}
func Register(c *gin.Context)         {}
func GetComments(c *gin.Context)      {}
func CreateComment(c *gin.Context)    {}
func DeleteComment(c *gin.Context)    {}
func AuthMiddleware() gin.HandlerFunc { return func(c *gin.Context) { c.Next() } }
