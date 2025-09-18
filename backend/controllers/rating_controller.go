package controllers

import (
	"fmt"
	"net/http"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"

	"github.com/gin-gonic/gin"
)

// CreateRating handles the creation of a new rating
func CreateRating(c *gin.Context) {
	var input struct {
		Score       float64 `json:"score" binding:"required"`
		Difficulty  float64 `json:"difficulty" binding:"required"`
		Usefulness  float64 `json:"usefulness" binding:"required"`
		Teaching    float64 `json:"teaching" binding:"required"`
		Content     string  `json:"content" binding:"required"`
		IsAnonymous bool    `json:"isAnonymous"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 从URL参数中获取courseId
	courseIDStr := c.Param("id")
	if courseIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	// 转换courseID为uint类型
	var courseID uint
	if _, err := fmt.Sscanf(courseIDStr, "%d", &courseID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Course ID"})
		return
	}

	rating := models.Rating{
		UserID:     userID.(uint),
		CourseID:   courseID,
		Score:      input.Score,
		Difficulty: input.Difficulty,
		Usefulness: input.Usefulness,
		Teaching:   input.Teaching,
		Content:    input.Content,
	}

	if err := config.DB.Create(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating created successfully", "data": rating})
}

func GetRatingsByCourse(c *gin.Context) {
	var ratings []models.Rating
	if err := config.DB.Preload("User").Where("course_id = ?", c.Param("id")).Find(&ratings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ratings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ratings})
}
