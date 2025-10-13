package controllers

import (
	"fmt"
	"net/http"
	"strconv"
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
	if err := config.DB.Preload("User").Where("course_id = ?", c.Param("id")).Order("created_at DESC").Find(&ratings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ratings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ratings})
}

// GetRating 获取单个评价
func GetRating(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 获取评价ID
	ratingIDStr := c.Param("id")
	ratingID, err := strconv.ParseUint(ratingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating ID"})
		return
	}

	// 查找评价
	var rating models.Rating
	if err := config.DB.Preload("User").Where("id = ?", ratingID).First(&rating).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found"})
		return
	}

	// 检查权限：只有评价的作者可以查看
	if rating.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own ratings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rating})
}

// UpdateRating 更新评价
func UpdateRating(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 获取评价ID
	ratingIDStr := c.Param("id")
	ratingID, err := strconv.ParseUint(ratingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating ID"})
		return
	}

	// 查找评价
	var rating models.Rating
	if err := config.DB.Where("id = ?", ratingID).First(&rating).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found"})
		return
	}

	// 检查权限：只有评价的作者可以更新
	if rating.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own ratings"})
		return
	}

	// 解析请求数据
	var input struct {
		Score       float64 `json:"score" binding:"required"`
		Difficulty  float64 `json:"difficulty" binding:"required"`
		Usefulness  float64 `json:"usefulness" binding:"required"`
		Teaching    float64 `json:"teaching" binding:"required"`
		Content     string  `json:"content"`
		IsAnonymous bool    `json:"isAnonymous"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新评价数据
	rating.Score = input.Score
	rating.Difficulty = input.Difficulty
	rating.Usefulness = input.Usefulness
	rating.Teaching = input.Teaching
	rating.Content = input.Content

	if err := config.DB.Save(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating updated successfully", "data": rating})
}

// DeleteRating 删除评价
func DeleteRating(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 获取评价ID
	ratingIDStr := c.Param("id")
	ratingID, err := strconv.ParseUint(ratingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating ID"})
		return
	}

	// 查找评价
	var rating models.Rating
	if err := config.DB.Where("id = ?", ratingID).First(&rating).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found"})
		return
	}

	// 检查权限：只有评价的作者可以删除
	if rating.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own ratings"})
		return
	}

	// 删除评价
	if err := config.DB.Delete(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating deleted successfully"})
}
