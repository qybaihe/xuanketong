package routes

import (
	"xuan-ke-tong/controllers"
	"xuan-ke-tong/middleware"

	"github.com/gin-gonic/gin"
)

func RatingRoutes(router *gin.Engine) {
	router.POST("/api/v1/ratings", middleware.AuthMiddleware(), controllers.CreateRating)
	router.POST("/api/v1/courses/:id/ratings", middleware.AuthMiddleware(), controllers.CreateRating)
	router.GET("/api/v1/courses/:id/ratings", controllers.GetRatingsByCourse)
	router.GET("/api/v1/ratings/:id", middleware.AuthMiddleware(), controllers.GetRating)
	router.PUT("/api/v1/ratings/:id", middleware.AuthMiddleware(), controllers.UpdateRating)
	router.DELETE("/api/v1/ratings/:id", middleware.AuthMiddleware(), controllers.DeleteRating)
}
