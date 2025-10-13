package routes

import (
	"xuan-ke-tong/controllers"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(router *gin.Engine) {
	router.POST("/api/v1/comments", controllers.CreateComment)
	router.GET("/api/v1/courses/:id/comments", controllers.GetCommentsByCourse)
	router.GET("/api/v1/comments/latest", controllers.GetLatestComments) // 获取最新评论
}
