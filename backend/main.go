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
	r.Run() // listen and serve on 0.0.0.0:8080
}
