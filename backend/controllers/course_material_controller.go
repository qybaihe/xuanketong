package controllers

import (
	"net/http"
	"strconv"
	"time"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"

	"github.com/gin-gonic/gin"
)

// CreateCourseMaterialRequest 创建课程资料请求
type CreateCourseMaterialRequest struct {
	MaterialName   string `json:"materialName" binding:"required"`
	Type           string `json:"type" binding:"required,oneof=cloud note"`
	Service        string `json:"service" binding:"required"`
	Link           string `json:"link" binding:"required"`
	AccessPassword string `json:"accessPassword"`
	Description    string `json:"description"`
}

// CourseMaterialResponse 课程资料响应
type CourseMaterialResponse struct {
	ID             uint      `json:"id"`
	CourseID       uint      `json:"courseId"`
	CourseName     string    `json:"courseName"`
	MaterialName   string    `json:"materialName"`
	Type           string    `json:"type"`
	Service        string    `json:"service"`
	Link           string    `json:"link"`
	AccessPassword string    `json:"accessPassword"`
	Description    string    `json:"description"`
	UploaderID     uint      `json:"uploaderId"`
	UploaderName   string    `json:"uploaderName"`
	UploadTime     time.Time `json:"uploadTime"`
	DownloadCount  int       `json:"downloadCount"`
}

// GetCourseMaterials 获取课程资料列表
func GetCourseMaterials(c *gin.Context) {
	courseIDStr := c.Param("courseId")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的课程ID"})
		return
	}

	var materials []models.CourseMaterial
	if err := config.DB.Where("course_id = ? AND is_active = ?", courseID, true).
		Preload("Uploader").
		Order("upload_time DESC").
		Find(&materials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取资料列表失败"})
		return
	}

	var response []CourseMaterialResponse
	for _, material := range materials {
		response = append(response, CourseMaterialResponse{
			ID:             material.ID,
			CourseID:       material.CourseID,
			MaterialName:   material.MaterialName,
			Type:           material.Type,
			Service:        material.Service,
			Link:           material.Link,
			AccessPassword: material.AccessPassword,
			Description:    material.Description,
			UploaderID:     material.UploaderID,
			UploaderName:   material.Uploader.Nickname,
			UploadTime:     material.UploadTime,
			DownloadCount:  material.DownloadCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total":     len(response),
			"materials": response,
		},
	})
}

// CreateCourseMaterial 创建课程资料
func CreateCourseMaterial(c *gin.Context) {
	// 获取当前用户
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	courseIDStr := c.Param("courseId")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的课程ID"})
		return
	}

	var req CreateCourseMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证课程是否存在
	var course models.Course
	if err := config.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	// 创建资料记录
	material := models.CourseMaterial{
		CourseID:       uint(courseID),
		MaterialName:   req.MaterialName,
		Type:           req.Type,
		Service:        req.Service,
		Link:           req.Link,
		AccessPassword: req.AccessPassword,
		Description:    req.Description,
		UploaderID:     userID.(uint),
		UploadTime:     time.Now(),
		DownloadCount:  0,
		IsActive:       true,
	}

	if err := config.DB.Create(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建资料失败"})
		return
	}

	// 获取上传者信息
	var uploader models.User
	config.DB.First(&uploader, material.UploaderID)

	response := CourseMaterialResponse{
		ID:             material.ID,
		CourseID:       material.CourseID,
		MaterialName:   material.MaterialName,
		Type:           material.Type,
		Service:        material.Service,
		Link:           material.Link,
		AccessPassword: material.AccessPassword,
		Description:    material.Description,
		UploaderID:     material.UploaderID,
		UploaderName:   uploader.Nickname,
		UploadTime:     material.UploadTime,
		DownloadCount:  material.DownloadCount,
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
	})
}

// UpdateCourseMaterial 更新课程资料
func UpdateCourseMaterial(c *gin.Context) {
	// 获取当前用户
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	materialIDStr := c.Param("materialId")
	materialID, err := strconv.ParseUint(materialIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的资料ID"})
		return
	}

	var material models.CourseMaterial
	if err := config.DB.First(&material, materialID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "资料不存在"})
		return
	}

	// 检查权限：只有上传者可以修改
	if material.UploaderID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限修改此资料"})
		return
	}

	var req CreateCourseMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新资料信息
	material.MaterialName = req.MaterialName
	material.Type = req.Type
	material.Service = req.Service
	material.Link = req.Link
	material.AccessPassword = req.AccessPassword
	material.Description = req.Description

	if err := config.DB.Save(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新资料失败"})
		return
	}

	// 获取上传者信息
	var uploader models.User
	config.DB.First(&uploader, material.UploaderID)

	response := CourseMaterialResponse{
		ID:             material.ID,
		CourseID:       material.CourseID,
		MaterialName:   material.MaterialName,
		Type:           material.Type,
		Service:        material.Service,
		Link:           material.Link,
		AccessPassword: material.AccessPassword,
		Description:    material.Description,
		UploaderID:     material.UploaderID,
		UploaderName:   uploader.Nickname,
		UploadTime:     material.UploadTime,
		DownloadCount:  material.DownloadCount,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// DeleteCourseMaterial 删除课程资料
func DeleteCourseMaterial(c *gin.Context) {
	// 获取当前用户
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	materialIDStr := c.Param("materialId")
	materialID, err := strconv.ParseUint(materialIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的资料ID"})
		return
	}

	var material models.CourseMaterial
	if err := config.DB.First(&material, materialID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "资料不存在"})
		return
	}

	// 检查权限：只有上传者可以删除
	if material.UploaderID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此资料"})
		return
	}

	// 软删除：设置is_active为false
	material.IsActive = false
	if err := config.DB.Save(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除资料失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "资料删除成功",
	})
}

// GetCourseMaterial 获取单个课程资料详情
func GetCourseMaterial(c *gin.Context) {
	materialIDStr := c.Param("materialId")
	materialID, err := strconv.ParseUint(materialIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的资料ID"})
		return
	}

	var material models.CourseMaterial
	if err := config.DB.Where("id = ? AND is_active = ?", materialID, true).
		Preload("Uploader").
		First(&material).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "资料不存在"})
		return
	}

	response := CourseMaterialResponse{
		ID:             material.ID,
		CourseID:       material.CourseID,
		MaterialName:   material.MaterialName,
		Type:           material.Type,
		Service:        material.Service,
		Link:           material.Link,
		AccessPassword: material.AccessPassword,
		Description:    material.Description,
		UploaderID:     material.UploaderID,
		UploaderName:   material.Uploader.Nickname,
		UploadTime:     material.UploadTime,
		DownloadCount:  material.DownloadCount,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// RedirectToMaterial 重定向到资料链接并增加浏览次数
func RedirectToMaterial(c *gin.Context) {
	materialIDStr := c.Param("materialId")
	materialID, err := strconv.ParseUint(materialIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的资料ID"})
		return
	}

	var material models.CourseMaterial
	if err := config.DB.Where("id = ? AND is_active = ?", materialID, true).First(&material).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "资料不存在"})
		return
	}

	// 增加浏览次数
	material.DownloadCount++
	if err := config.DB.Save(&material).Error; err != nil {
		// 即使更新浏览次数失败，也要重定向到链接
		c.Redirect(http.StatusTemporaryRedirect, material.Link)
		return
	}

	// 重定向到实际的资料链接
	c.Redirect(http.StatusTemporaryRedirect, material.Link)
}

// GetPopularMaterials 获取热门资料（按浏览量排序）
func GetPopularMaterials(c *gin.Context) {
	// 获取查询参数
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	var materials []models.CourseMaterial
	if err := config.DB.Where("is_active = ?", true).
		Preload("Course").
		Preload("Uploader").
		Order("download_count DESC").
		Limit(limit).
		Find(&materials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取热门资料失败"})
		return
	}

	// 构建响应数据
	var response []CourseMaterialResponse
	for _, material := range materials {
		response = append(response, CourseMaterialResponse{
			ID:             material.ID,
			CourseID:       material.CourseID,
			CourseName:     material.Course.Name,
			MaterialName:   material.MaterialName,
			Type:           material.Type,
			Service:        material.Service,
			Link:           material.Link,
			AccessPassword: material.AccessPassword,
			Description:    material.Description,
			UploaderID:     material.UploaderID,
			UploaderName:   material.Uploader.Nickname,
			UploadTime:     material.UploadTime,
			DownloadCount:  material.DownloadCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// IncrementDownloadCount 增加下载次数（保留用于其他用途）
func IncrementDownloadCount(c *gin.Context) {
	materialIDStr := c.Param("materialId")
	materialID, err := strconv.ParseUint(materialIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的资料ID"})
		return
	}

	var material models.CourseMaterial
	if err := config.DB.First(&material, materialID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "资料不存在"})
		return
	}

	// 增加下载次数
	material.DownloadCount++
	if err := config.DB.Save(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新下载次数失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"downloadCount": material.DownloadCount,
		},
	})
}
