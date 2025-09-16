package routes

import (
	"xuan-ke-tong/controllers"

	"github.com/gin-gonic/gin"
)

func RatingRoutes(router *gin.Engine) {
	router.POST("/api/v1/ratings", controllers.CreateRating)
	router.GET("/api/v1/courses/:id/ratings", controllers.GetRatingsByCourse)
}
