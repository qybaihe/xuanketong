package routes

import (
	"xuan-ke-tong/controllers"
	"xuan-ke-tong/middleware"

	"github.com/gin-gonic/gin"
)

func CourseNoteRoutes(r *gin.Engine) {
	// 课程笔记相关路由 - 使用不同的路径避免冲突
	noteGroup := r.Group("/api/v1/course-notes")
	{
		// 获取课程笔记列表（公开访问，但只显示公开的笔记）
		noteGroup.GET("/course/:courseId", controllers.GetCourseNotes)

		// 需要认证的路由
		authGroup := noteGroup.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			// 创建课程笔记
			authGroup.POST("/course/:courseId", controllers.CreateCourseNote)
		}
	}

	// 单个笔记操作路由
	noteDetailGroup := r.Group("/api/v1/notes/:noteId")
	{
		// 获取笔记详情（使用可选认证中间件）
		noteDetailGroup.GET("", middleware.OptionalAuthMiddleware(), controllers.GetCourseNote)

		// 需要认证的路由
		authNoteDetailGroup := noteDetailGroup.Group("")
		authNoteDetailGroup.Use(middleware.AuthMiddleware())
		{
			// 更新课程笔记
			authNoteDetailGroup.PUT("", controllers.UpdateCourseNote)

			// 删除课程笔记
			authNoteDetailGroup.DELETE("", controllers.DeleteCourseNote)
		}
	}

	// 用户笔记相关路由
	userNoteGroup := r.Group("/api/v1/users/notes")
	userNoteGroup.Use(middleware.AuthMiddleware())
	{
		// 获取当前用户的笔记列表
		userNoteGroup.GET("", controllers.GetUserNotes)
	}
}
