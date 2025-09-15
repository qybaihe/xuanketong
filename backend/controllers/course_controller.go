package controllers

import (
	"net/http"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"

	"github.com/gin-gonic/gin"
)

func CreateCourse(c *gin.Context) {
	var input models.Course
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course := models.Course{
		Name:        input.Name,
		Description: input.Description,
		Grade:       input.Grade,
		Semester:    input.Semester,
		Subject:     input.Subject,
		Teacher:     input.Teacher,
		Credits:     input.Credits,
		ImageURL:    input.ImageURL,
	}

	if err := config.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course created successfully"})
}

func GetCourses(c *gin.Context) {
	var courses []models.Course
	query := config.DB

	if grade := c.Query("grade"); grade != "" {
		query = query.Where("grade = ?", grade)
	}
	if semester := c.Query("semester"); semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if subject := c.Query("subject"); subject != "" {
		query = query.Where("subject = ?", subject)
	}

	if err := query.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": courses})
}

func GetCourse(c *gin.Context) {
	var course models.Course
	if err := config.DB.Where("id = ?", c.Param("id")).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": course})
}

func UpdateCourse(c *gin.Context) {
	var course models.Course
	if err := config.DB.Where("id = ?", c.Param("id")).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	var input models.Course
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&course).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": course})
}

func DeleteCourse(c *gin.Context) {
	var course models.Course
	if err := config.DB.Where("id = ?", c.Param("id")).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	config.DB.Delete(&course)
	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}
