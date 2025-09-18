package controllers

import (
	"net/http"
	"strconv"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"

	"github.com/gin-gonic/gin"
)

// CreateEvaluationRequest 处理POST /api/evaluation-requests请求，创建新的求评价请求
func CreateEvaluationRequest(c *gin.Context) {
	// 获取当前用户信息
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	var input struct {
		CourseID uint `json:"courseId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "courseId无效或缺失"})
		return
	}

	// 检查课程是否存在
	var course models.Course
	if err := config.DB.First(&course, input.CourseID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "课程不存在"})
		return
	}

	// 检查是否已经存在活跃的求评价请求
	var existingRequest models.EvaluationRequest
	if err := config.DB.Where("user_id = ? AND course_id = ? AND status = ?", userID, input.CourseID, "pending").First(&existingRequest).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "该课程已经存在一个活跃的求评价请求"})
		return
	}

	// 创建新的求评价请求
	evaluationRequest := models.EvaluationRequest{
		UserID:   userID.(uint),
		CourseID: input.CourseID,
		Status:   "pending",
	}

	if err := config.DB.Create(&evaluationRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建求评价请求失败"})
		return
	}

	// 加载关联的用户和课程信息
	config.DB.Preload("User").Preload("Course").First(&evaluationRequest, evaluationRequest.ID)

	// 返回响应数据
	response := gin.H{
		"id":        evaluationRequest.ID,
		"userId":    evaluationRequest.UserID,
		"courseId":  evaluationRequest.CourseID,
		"status":    evaluationRequest.Status,
		"createdAt": evaluationRequest.CreatedAt,
		"user": gin.H{
			"id":       evaluationRequest.User.ID,
			"nickname": evaluationRequest.User.Nickname,
		},
		"course": gin.H{
			"id":   evaluationRequest.Course.ID,
			"name": evaluationRequest.Course.Name,
		},
	}

	c.JSON(http.StatusCreated, response)
}

// GetEvaluationRequests 处理GET /api/evaluation-requests请求，获取所有状态为pending的求评价列表
func GetEvaluationRequests(c *gin.Context) {
	// 获取分页参数
	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if parsedPageSize, err := strconv.Atoi(ps); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询求评价请求列表
	var evaluationRequests []models.EvaluationRequest
	var total int64

	// 先获取总数
	config.DB.Model(&models.EvaluationRequest{}).Where("status = ?", "pending").Count(&total)

	// 获取分页数据，包含关联的用户和课程信息
	if err := config.DB.Preload("User").Preload("Course").
		Where("status = ?", "pending").
		Limit(pageSize).Offset(offset).
		Find(&evaluationRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取求评价列表失败"})
		return
	}

	// 构建响应数据
	var items []gin.H
	for _, request := range evaluationRequests {
		item := gin.H{
			"id":        request.ID,
			"userId":    request.UserID,
			"courseId":  request.CourseID,
			"status":    request.Status,
			"createdAt": request.CreatedAt,
			"user": gin.H{
				"id":       request.User.ID,
				"nickname": request.User.Nickname,
			},
			"course": gin.H{
				"id":       request.Course.ID,
				"name":     request.Course.Name,
				"teacher":  request.Course.Teacher,
				"imageURL": request.Course.ImageURL,
			},
		}
		items = append(items, item)
	}

	response := gin.H{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"items":    items,
	}

	c.JSON(http.StatusOK, response)
}
