package controllers

import (
	"net/http"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"

	"github.com/gin-gonic/gin"
)

func CreateCourse(c *gin.Context) {
	var input struct {
		Name        string `json:"Name"`
		Description string `json:"Description"`
		Grade       string `json:"Grade"`
		Semester    string `json:"Semester"`
		Subject     string `json:"Subject"`
		Teacher     string `json:"Teacher"`
		Credits     int    `json:"Credits"`
		ImageURL    string `json:"ImageURL"`
	}

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

	c.JSON(http.StatusOK, gin.H{"data": course})
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

	if err := query.Unscoped().Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get courses"})
		return
	}

	// 获取每个课程的平均评分
	type CourseWithRating struct {
		models.Course
		AverageRating float64 `json:"averageRating"`
		TotalRatings  int     `json:"totalRatings"`
	}

	var coursesWithRatings []CourseWithRating

	for _, course := range courses {
		var avgRating float64
		var totalRatings int64

		// 计算平均评分
		config.DB.Model(&models.Rating{}).
			Where("course_id = ?", course.ID).
			Select("AVG(score) as avg_rating, COUNT(*) as total").
			Row().
			Scan(&avgRating, &totalRatings)

		if avgRating == 0 {
			avgRating = 0
		}

		coursesWithRatings = append(coursesWithRatings, CourseWithRating{
			Course:        course,
			AverageRating: avgRating,
			TotalRatings:  int(totalRatings),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": coursesWithRatings})
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
	if err := config.DB.First(&course, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	var input models.Course
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 只更新允许的字段
	updateData := map[string]interface{}{
		"Name":        input.Name,
		"Description": input.Description,
		"Grade":       input.Grade,
		"Semester":    input.Semester,
		"Subject":     input.Subject,
		"Teacher":     input.Teacher,
		"Credits":     input.Credits,
		"ImageURL":    input.ImageURL,
	}

	if err := config.DB.Model(&course).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": course})
}

func DeleteCourse(c *gin.Context) {
	var course models.Course
	if err := config.DB.First(&course, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	if err := config.DB.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}
