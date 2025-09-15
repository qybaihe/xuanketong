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
