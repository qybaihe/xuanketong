package controllers

import (
	"math/rand"
	"net/http"
	"time"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"

	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	var input struct {
		CourseID uint   `json:"courseId" binding:"required"`
		Content  string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	comment := models.Comment{
		UserID:   userID.(uint),
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
	if err := config.DB.Preload("User").Where("course_id = ?", c.Param("id")).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

// GetLatestComments 获取最新评论（用于首页展示）
// 从有评论的课程中随机抽取三条评论
func GetLatestComments(c *gin.Context) {
	// 首先获取所有有评论的课程ID
	var courseIDs []uint
	if err := config.DB.Model(&models.Comment{}).Distinct("course_id").Pluck("course_id", &courseIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve courses with comments"})
		return
	}

	// 如果没有课程有评论，返回空数组
	if len(courseIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
		return
	}

	// 从有评论的课程中随机选择最多3个课程
	selectedCourseIDs := courseIDs
	if len(courseIDs) > 3 {
		// 使用随机算法选择3个课程
		rand.Seed(time.Now().UnixNano())
		selectedCourseIDs = make([]uint, 3)
		copy(selectedCourseIDs, courseIDs)

		// Fisher-Yates 洗牌算法的简化版本
		for i := len(selectedCourseIDs) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			selectedCourseIDs[i], selectedCourseIDs[j] = selectedCourseIDs[j], selectedCourseIDs[i]
		}
		selectedCourseIDs = selectedCourseIDs[:3]
	}

	// 从选中的课程中获取评论
	var comments []models.Comment
	if err := config.DB.Preload("User").Where("course_id IN ?", selectedCourseIDs).Order("created_at DESC").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	// 如果评论数量超过3条，只取前3条
	if len(comments) > 3 {
		comments = comments[:3]
	}

	// 构建返回数据，包含课程信息
	type CommentWithCourse struct {
		models.Comment
		CourseName    string `json:"courseName"`
		CourseTeacher string `json:"courseTeacher"`
	}

	var commentsWithCourse []CommentWithCourse
	for _, comment := range comments {
		// 单独查询课程信息
		var course models.Course
		if err := config.DB.First(&course, comment.CourseID).Error; err == nil {
			commentsWithCourse = append(commentsWithCourse, CommentWithCourse{
				Comment:       comment,
				CourseName:    course.Name,
				CourseTeacher: course.Teacher,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": commentsWithCourse})
}
