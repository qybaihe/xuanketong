package main

import (
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"
	"xuan-ke-tong/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{}, &models.Course{}, &models.Rating{}, &models.Comment{})

	routes.AuthRoutes(r)
	routes.CourseRoutes(r)
	routes.RatingRoutes(r)
	routes.CommentRoutes(r)
	routes.AdminRoutes(r)
	routes.UserRoutes(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 添加管理员路由日志
	r.GET("/admin-routes", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Admin routes are registered",
		})
	})

	// 打印所有注册的路由
	routes := r.Routes()
	println("Registered routes:")
	for _, route := range routes {
		println(route.Method, route.Path)
	}

	// 添加一个简单的测试路由
	r.POST("/api/v1/admin/courses/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Test route works"})
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
