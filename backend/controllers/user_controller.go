package controllers

import (
	"net/http"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var ratings []models.Rating
	config.DB.Where("user_id = ?", c.Param("id")).Find(&ratings)

	var comments []models.Comment
	config.DB.Where("user_id = ?", c.Param("id")).Find(&comments)

	c.JSON(http.StatusOK, gin.H{
		"user":     user,
		"ratings":  ratings,
		"comments": comments,
	})
}
