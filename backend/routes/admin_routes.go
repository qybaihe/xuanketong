package routes

import (
	"xuan-ke-tong/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {
	admin := router.Group("/api/v1/admin")
	{
		admin.GET("/stats", controllers.GetStats)
		admin.GET("/users", controllers.GetAllUsers)
		admin.GET("/ratings", controllers.GetAllRatings)
		admin.GET("/comments", controllers.GetAllComments)
	}
}
