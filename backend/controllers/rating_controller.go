package controllers

import (
	"net/http"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"

	"github.com/gin-gonic/gin"
)

func CreateRating(c *gin.Context) {
	var input models.Rating
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rating := models.Rating{
		UserID:   input.UserID,
		CourseID: input.CourseID,
		Score:    input.Score,
	}

	if err := config.DB.Create(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating created successfully"})
}

func GetRatingsByCourse(c *gin.Context) {
	var ratings []models.Rating
	config.DB.Where("course_id = ?", c.Param("id")).Find(&ratings)
	c.JSON(http.StatusOK, gin.H{"data": ratings})
}
