package routes

import (
	"xuan-ke-tong/controllers"
	"xuan-ke-tong/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/api/v1/auth/register", controllers.Register)
	router.POST("/api/v1/auth/login", controllers.Login)
	router.GET("/api/v1/auth/me", middleware.AuthMiddleware(), controllers.GetCurrentUser)
}
