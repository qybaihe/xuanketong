package controllers

import (
	"net/http"
	"strconv"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"

	"github.com/gin-gonic/gin"
)

// CreateCourseNoteRequest 创建课程笔记请求
type CreateCourseNoteRequest struct {
	Content  string `json:"content" binding:"required"`
	IsPublic bool   `json:"isPublic"`
}

// UpdateCourseNoteRequest 更新课程笔记请求
type UpdateCourseNoteRequest struct {
	Content  string `json:"content" binding:"required"`
	IsPublic bool   `json:"isPublic"`
}

// CourseNoteResponse 课程笔记响应
type CourseNoteResponse struct {
	ID         uint   `json:"id"`
	CourseID   uint   `json:"courseId"`
	Content    string `json:"content"`
	AuthorID   uint   `json:"authorId"`
	AuthorName string `json:"authorName"`
	IsPublic   bool   `json:"isPublic"`
	IsOwner    bool   `json:"isOwner"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

// GetCourseNotes 获取课程笔记列表
func GetCourseNotes(c *gin.Context) {
	courseIDStr := c.Param("courseId")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的课程ID"})
		return
	}

	// 获取当前用户ID（可选）
	var currentUserID uint
	if userID, exists := c.Get("userId"); exists {
		currentUserID = userID.(uint)
	}

	var notes []models.CourseNote
	if err := config.DB.Where("course_id = ?", courseID).
		Preload("Author").
		Order("created_at DESC").
		Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取笔记列表失败"})
		return
	}

	var response []CourseNoteResponse
	for _, note := range notes {
		// 如果笔记不是公开的，只有作者本人可以看到
		if !note.IsPublic && currentUserID != note.AuthorID {
			continue
		}

		response = append(response, CourseNoteResponse{
			ID:         note.ID,
			CourseID:   note.CourseID,
			Content:    note.Content,
			AuthorID:   note.AuthorID,
			AuthorName: note.Author.Nickname,
			IsPublic:   note.IsPublic,
			IsOwner:    currentUserID == note.AuthorID,
			CreatedAt:  note.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:  note.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"notes": response,
		},
	})
}

// CreateCourseNote 创建课程笔记
func CreateCourseNote(c *gin.Context) {
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

	var req CreateCourseNoteRequest
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

	// 创建笔记记录
	note := models.CourseNote{
		CourseID: uint(courseID),
		Content:  req.Content,
		AuthorID: userID.(uint),
		IsPublic: req.IsPublic,
	}

	if err := config.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建笔记失败"})
		return
	}

	// 获取作者信息
	var author models.User
	config.DB.First(&author, note.AuthorID)

	response := CourseNoteResponse{
		ID:         note.ID,
		CourseID:   note.CourseID,
		Content:    note.Content,
		AuthorID:   note.AuthorID,
		AuthorName: author.Nickname,
		IsPublic:   note.IsPublic,
		IsOwner:    true,
		CreatedAt:  note.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  note.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
	})
}

// UpdateCourseNote 更新课程笔记
func UpdateCourseNote(c *gin.Context) {
	// 获取当前用户
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	noteIDStr := c.Param("noteId")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}

	var note models.CourseNote
	if err := config.DB.First(&note, noteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	// 检查权限：只有作者可以修改
	if note.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限修改此笔记"})
		return
	}

	var req UpdateCourseNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新笔记内容
	note.Content = req.Content
	note.IsPublic = req.IsPublic

	if err := config.DB.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新笔记失败"})
		return
	}

	// 获取作者信息
	var author models.User
	config.DB.First(&author, note.AuthorID)

	response := CourseNoteResponse{
		ID:         note.ID,
		CourseID:   note.CourseID,
		Content:    note.Content,
		AuthorID:   note.AuthorID,
		AuthorName: author.Nickname,
		IsPublic:   note.IsPublic,
		IsOwner:    true,
		CreatedAt:  note.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  note.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// DeleteCourseNote 删除课程笔记
func DeleteCourseNote(c *gin.Context) {
	// 获取当前用户
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	noteIDStr := c.Param("noteId")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}

	var note models.CourseNote
	if err := config.DB.First(&note, noteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	// 检查权限：只有作者可以删除
	if note.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此笔记"})
		return
	}

	if err := config.DB.Delete(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除笔记失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "笔记删除成功",
	})
}

// GetUserNotes 获取用户的笔记列表
func GetUserNotes(c *gin.Context) {
	// 获取当前用户
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var notes []models.CourseNote
	if err := config.DB.Where("author_id = ?", userID).
		Preload("Course").
		Order("updated_at DESC").
		Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取笔记列表失败"})
		return
	}

	var response []CourseNoteResponse
	for _, note := range notes {
		response = append(response, CourseNoteResponse{
			ID:         note.ID,
			CourseID:   note.CourseID,
			Content:    note.Content,
			AuthorID:   note.AuthorID,
			AuthorName: note.Author.Nickname,
			IsPublic:   note.IsPublic,
			IsOwner:    true,
			CreatedAt:  note.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:  note.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"notes": response,
		},
	})
}
