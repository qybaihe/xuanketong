package routes

import (
	"xuan-ke-tong/controllers"
	"xuan-ke-tong/middleware"

	"github.com/gin-gonic/gin"
)

func EvaluationRequestRoutes(router *gin.Engine) {
	// 获取求评价列表 - 公开访问
	router.GET("/api/evaluation-requests", controllers.GetEvaluationRequests)

	// 创建求评价请求 - 需要认证
	router.POST("/api/evaluation-requests", middleware.AuthMiddleware(), controllers.CreateEvaluationRequest)
}
