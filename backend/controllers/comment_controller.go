package controllers

import (
	"net/http"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"

	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	var input models.Comment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := models.Comment{
		UserID:   input.UserID,
		CourseID: input.CourseID,
		Content:  input.Content,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully"})
}

func GetCommentsByCourse(c *gin.Context) {
	var comments []models.Comment
	config.DB.Where("course_id = ?", c.Param("id")).Find(&comments)

	// 获取用户信息
	type CommentWithUser struct {
		models.Comment
		Username string `json:"username"`
		Nickname string `json:"nickname"`
	}

	var result []CommentWithUser
	for _, comment := range comments {
		var user models.User
		if err := config.DB.Where("id = ?", comment.UserID).First(&user).Error; err == nil {
			result = append(result, CommentWithUser{
				Comment:  comment,
				Username: user.Username,
				Nickname: user.Nickname,
			})
		} else {
			// 如果用户不存在，使用默认值
			result = append(result, CommentWithUser{
				Comment:  comment,
				Username: "用户" + string(comment.UserID),
				Nickname: "用户" + string(comment.UserID),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
