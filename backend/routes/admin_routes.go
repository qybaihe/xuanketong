package routes

import (
	"xuan-ke-tong/controllers"
	"xuan-ke-tong/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.RequireAdminMiddleware())
	{
		admin.GET("/stats", controllers.GetStats)
		admin.GET("/stats/enhanced", controllers.GetEnhancedHomeStats) // 增强版首页统计
		admin.GET("/user-stats", controllers.GetUserStats)

		// 用户管理路由
		admin.GET("/users", controllers.GetAllUsers)
		admin.GET("/users/:id", controllers.GetUserByID)
		admin.PUT("/users/:id", controllers.UpdateUser)
		admin.DELETE("/users/:id", controllers.DeleteUser)

		admin.GET("/ratings", controllers.GetAllRatings)
		admin.GET("/comments", controllers.GetAllComments)

		// 课程管理路由
		admin.POST("/courses", controllers.CreateCourse)
		admin.GET("/courses", controllers.GetCourses)
		admin.GET("/courses/:id", controllers.GetCourse)
		admin.PUT("/courses/:id", controllers.UpdateCourse)
		admin.DELETE("/courses/:id", controllers.DeleteCourse)

		// 测试路由
		admin.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Admin routes are working"})
		})
	}

	// 添加根级别的测试路由
	router.GET("/admin-test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "AdminRoutes function was called"})
	})
}
