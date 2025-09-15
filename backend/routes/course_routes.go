package routes

import (
	"xuan-ke-tong/controllers"

	"github.com/gin-gonic/gin"
)

func CourseRoutes(router *gin.Engine) {
	router.POST("/api/v1/courses", controllers.CreateCourse)
	router.GET("/api/v1/courses", controllers.GetCourses)
	router.GET("/api/v1/courses/:id", controllers.GetCourse)
	router.PUT("/api/v1/courses/:id", controllers.UpdateCourse)
	router.DELETE("/api/v1/courses/:id", controllers.DeleteCourse)
}
