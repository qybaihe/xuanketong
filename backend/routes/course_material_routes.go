package routes

import (
	"xuan-ke-tong/controllers"
	"xuan-ke-tong/middleware"

	"github.com/gin-gonic/gin"
)

func CourseMaterialRoutes(r *gin.Engine) {
	// 课程资料相关路由 - 使用不同的路径避免冲突
	materialGroup := r.Group("/api/v1/course-materials")
	{
		// 获取课程资料列表（公开访问）
		materialGroup.GET("/course/:courseId", controllers.GetCourseMaterials)

		// 获取热门资料（公开访问）
		materialGroup.GET("/popular", controllers.GetPopularMaterials)

		// 需要认证的路由
		authGroup := materialGroup.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			// 创建课程资料
			authGroup.POST("/course/:courseId", controllers.CreateCourseMaterial)
		}
	}

	// 单个资料操作路由
	materialDetailGroup := r.Group("/api/v1/materials/:materialId")
	{
		// 获取资料详情（公开访问）
		materialDetailGroup.GET("", controllers.GetCourseMaterial)

		// 重定向到资料链接并统计浏览次数（公开访问）
		materialDetailGroup.GET("/redirect", controllers.RedirectToMaterial)

		// 增加下载次数（公开访问）
		materialDetailGroup.POST("/download", controllers.IncrementDownloadCount)

		// 需要认证的路由
		authDetailGroup := materialDetailGroup.Group("")
		authDetailGroup.Use(middleware.AuthMiddleware())
		{
			// 更新课程资料
			authDetailGroup.PUT("", controllers.UpdateCourseMaterial)

			// 删除课程资料
			authDetailGroup.DELETE("", controllers.DeleteCourseMaterial)
		}
	}
}
