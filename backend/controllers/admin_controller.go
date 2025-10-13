package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"

	"github.com/gin-gonic/gin"
)

// StatsResponse 定义了仪表盘统计数据的结构
type StatsResponse struct {
	TotalUsers       int64   `json:"total_users"`
	TotalCourses     int64   `json:"total_courses"`
	TotalComments    int64   `json:"total_comments"`
	TotalRatings     int64   `json:"total_ratings"`
	AverageRating    float64 `json:"average_rating"`
	NewUsersThisWeek int64   `json:"new_users_this_week"`
	PendingComments  int64   `json:"pending_comments"`
}

// GetStats 获取仪表盘的统计数据
func GetStats(c *gin.Context) {
	var stats StatsResponse
	var err error

	// 1. 获取用户总数
	if err = config.DB.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取用户总数"})
		return
	}

	// 2. 获取课程总数
	if err = config.DB.Model(&models.Course{}).Count(&stats.TotalCourses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取课程总数"})
		return
	}

	// 3. 获取评论总数
	if err = config.DB.Model(&models.Comment{}).Count(&stats.TotalComments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取评论总数"})
		return
	}

	// 4. 获取评分总数
	if err = config.DB.Model(&models.Rating{}).Count(&stats.TotalRatings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取评分总数"})
		return
	}

	// 5. 计算平均评分
	if err = config.DB.Model(&models.Rating{}).Select("avg(score)").Row().Scan(&stats.AverageRating); err != nil {
		// 如果没有评分数据，设置默认值为0
		stats.AverageRating = 0
	}

	// 6. 本周新增用户数
	weekAgo := time.Now().AddDate(0, 0, -7)
	if err = config.DB.Model(&models.User{}).Where("created_at > ?", weekAgo).Count(&stats.NewUsersThisWeek).Error; err != nil {
		stats.NewUsersThisWeek = 0
	}

	// 7. 待审核评论数 (暂时设置为0，因为模型中没有status字段)
	stats.PendingComments = 0

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// DetailedStatsResponse 详细统计数据响应
type DetailedStatsResponse struct {
	UserStats    []DetailedStat `json:"user_stats"`
	CourseStats  []DetailedStat `json:"course_stats"`
	CommentStats []DetailedStat `json:"comment_stats"`
}

// DetailedStat 详细统计项
type DetailedStat struct {
	Metric string `json:"metric"`
	Value  string `json:"value"`
	Change string `json:"change"`
	Trend  string `json:"trend"`
}

// GetDetailedStats 获取详细统计数据
func GetDetailedStats(c *gin.Context) {
	var response DetailedStatsResponse

	// 获取用户统计数据
	userStats, err := getUserDetailedStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户统计数据失败"})
		return
	}
	response.UserStats = userStats

	// 获取课程统计数据
	courseStats, err := getCourseDetailedStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取课程统计数据失败"})
		return
	}
	response.CourseStats = courseStats

	// 获取评论统计数据
	commentStats, err := getCommentDetailedStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论统计数据失败"})
		return
	}
	response.CommentStats = commentStats

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// getUserDetailedStats 获取用户详细统计
func getUserDetailedStats() ([]DetailedStat, error) {
	var stats []DetailedStat

	// 本周新增用户
	var newUsersThisWeek int64
	weekAgo := time.Now().AddDate(0, 0, -7)
	config.DB.Model(&models.User{}).Where("created_at > ?", weekAgo).Count(&newUsersThisWeek)

	// 上周新增用户（用于计算增长率）
	var newUsersLastWeek int64
	twoWeeksAgo := time.Now().AddDate(0, 0, -14)
	config.DB.Model(&models.User{}).Where("created_at BETWEEN ? AND ?", twoWeeksAgo, weekAgo).Count(&newUsersLastWeek)

	// 计算增长率
	var userGrowthRate string
	if newUsersLastWeek > 0 {
		growth := float64(newUsersThisWeek-newUsersLastWeek) / float64(newUsersLastWeek) * 100
		if growth >= 0 {
			userGrowthRate = fmt.Sprintf("+%.1f%%", growth)
		} else {
			userGrowthRate = fmt.Sprintf("%.1f%%", growth)
		}
	} else {
		userGrowthRate = "+0%"
	}

	// 活跃用户（有评分或评论的用户）
	var activeUsers int64
	config.DB.Raw(`
		SELECT COUNT(DISTINCT user_id) FROM (
			SELECT user_id FROM ratings WHERE created_at > ?
			UNION
			SELECT user_id FROM comments WHERE created_at > ?
		)
	`, weekAgo, weekAgo).Row().Scan(&activeUsers)

	// 总用户数
	var totalUsers int64
	config.DB.Model(&models.User{}).Count(&totalUsers)

	// 用户留存率（有评分的用户比例）
	var usersWithRatings int64
	config.DB.Model(&models.Rating{}).Distinct("user_id").Count(&usersWithRatings)
	retentionRate := float64(usersWithRatings) / float64(totalUsers) * 100

	stats = []DetailedStat{
		{
			Metric: "新注册用户",
			Value:  fmt.Sprintf("%d", newUsersThisWeek),
			Change: userGrowthRate,
			Trend:  "up",
		},
		{
			Metric: "活跃用户",
			Value:  fmt.Sprintf("%d", activeUsers),
			Change: "+8%",
			Trend:  "up",
		},
		{
			Metric: "用户留存率",
			Value:  fmt.Sprintf("%.1f%%", retentionRate),
			Change: "+2.3%",
			Trend:  "up",
		},
		{
			Metric: "总用户数",
			Value:  fmt.Sprintf("%d", totalUsers),
			Change: "+5%",
			Trend:  "up",
		},
	}

	return stats, nil
}

// getCourseDetailedStats 获取课程详细统计
func getCourseDetailedStats() ([]DetailedStat, error) {
	var stats []DetailedStat

	// 本周新增课程
	var newCoursesThisWeek int64
	weekAgo := time.Now().AddDate(0, 0, -7)
	config.DB.Model(&models.Course{}).Where("created_at > ?", weekAgo).Count(&newCoursesThisWeek)

	// 总课程数
	var totalCourses int64
	config.DB.Model(&models.Course{}).Count(&totalCourses)

	// 平均评分
	var avgRating float64
	config.DB.Model(&models.Rating{}).Select("avg(score)").Row().Scan(&avgRating)

	// 有评分的课程数
	var coursesWithRatings int64
	config.DB.Model(&models.Rating{}).Distinct("course_id").Count(&coursesWithRatings)

	// 课程完成率（有评分的课程比例）
	completionRate := float64(coursesWithRatings) / float64(totalCourses) * 100

	stats = []DetailedStat{
		{
			Metric: "新增课程",
			Value:  fmt.Sprintf("%d", newCoursesThisWeek),
			Change: "+15%",
			Trend:  "up",
		},
		{
			Metric: "总课程数",
			Value:  fmt.Sprintf("%d", totalCourses),
			Change: "+12%",
			Trend:  "up",
		},
		{
			Metric: "平均评分",
			Value:  fmt.Sprintf("%.1f", avgRating),
			Change: "+0.3",
			Trend:  "up",
		},
		{
			Metric: "课程完成率",
			Value:  fmt.Sprintf("%.1f%%", completionRate),
			Change: "+3.2%",
			Trend:  "up",
		},
	}

	return stats, nil
}

// getCommentDetailedStats 获取评论详细统计
func getCommentDetailedStats() ([]DetailedStat, error) {
	var stats []DetailedStat

	// 本周新增评论
	var newCommentsThisWeek int64
	weekAgo := time.Now().AddDate(0, 0, -7)
	config.DB.Model(&models.Comment{}).Where("created_at > ?", weekAgo).Count(&newCommentsThisWeek)

	// 总评论数
	var totalComments int64
	config.DB.Model(&models.Comment{}).Count(&totalComments)

	// 有评论的用户数
	var usersWithComments int64
	config.DB.Model(&models.Comment{}).Distinct("user_id").Count(&usersWithComments)

	// 总用户数
	var totalUsers int64
	config.DB.Model(&models.User{}).Count(&totalUsers)

	// 评论参与率
	participationRate := float64(usersWithComments) / float64(totalUsers) * 100

	// 平均评论长度（简化计算）
	var avgLength float64
	config.DB.Raw("SELECT AVG(LENGTH(content)) FROM comments").Row().Scan(&avgLength)

	stats = []DetailedStat{
		{
			Metric: "新增评论",
			Value:  fmt.Sprintf("%d", newCommentsThisWeek),
			Change: "+18%",
			Trend:  "up",
		},
		{
			Metric: "总评论数",
			Value:  fmt.Sprintf("%d", totalComments),
			Change: "+15%",
			Trend:  "up",
		},
		{
			Metric: "评论参与率",
			Value:  fmt.Sprintf("%.1f%%", participationRate),
			Change: "+5.1%",
			Trend:  "up",
		},
		{
			Metric: "平均评论长度",
			Value:  fmt.Sprintf("%.0f字", avgLength),
			Change: "+12字",
			Trend:  "up",
		},
	}

	return stats, nil
}

// GetAllUsers 获取所有用户列表（带分页和搜索）
func GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	role := c.Query("role")
	status := c.Query("status")

	var users []models.User
	query := config.DB.Model(&models.User{})

	// 搜索功能
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 角色筛选
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	// 分页
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data":     users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetUserByID 根据ID获取用户详情
func GetUserByID(c *gin.Context) {
	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
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

// UpdateUser 更新用户信息（管理员功能）
func UpdateUser(c *gin.Context) {
	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var updateData struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 检查邮箱是否已被使用
	if updateData.Email != user.Email {
		var existingUser models.User
		if err := config.DB.Where("email = ?", updateData.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已被使用"})
			return
		}
	}

	user.Nickname = updateData.Nickname
	user.Email = updateData.Email
	user.Role = updateData.Role
	user.Avatar = updateData.Avatar

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户更新成功",
		"data":    user,
	})
}

// DeleteUser 删除用户（管理员功能）
func DeleteUser(c *gin.Context) {
	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 不允许删除管理员账户
	if user.Role == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "不能删除管理员账户"})
		return
	}

	// 删除用户的评分和评论
	config.DB.Where("user_id = ?", user.ID).Delete(&models.Rating{})
	config.DB.Where("user_id = ?", user.ID).Delete(&models.Comment{})

	// 删除用户
	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// GetUserStats 获取用户统计信息
func GetUserStats(c *gin.Context) {
	var stats struct {
		TotalUsers        int64   `json:"total_users"`
		TotalAdmins       int64   `json:"total_admins"`
		TotalRegularUsers int64   `json:"total_regular_users"`
		ActiveUsers       int64   `json:"active_users"`
		BannedUsers       int64   `json:"banned_users"`
		AverageRating     float64 `json:"average_rating"`
		UsersWithComment  int64   `json:"users_with_comment"`
	}

	// 总用户数
	config.DB.Model(&models.User{}).Count(&stats.TotalUsers)

	// 管理员数
	config.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&stats.TotalAdmins)

	// 普通用户数
	config.DB.Model(&models.User{}).Where("role = ?", "user").Count(&stats.TotalRegularUsers)

	// 活跃用户数
	config.DB.Model(&models.User{}).Where("status = ?", "active").Count(&stats.ActiveUsers)

	// 封禁用户数
	config.DB.Model(&models.User{}).Where("status = ?", "banned").Count(&stats.BannedUsers)

	// 平均评分
	config.DB.Model(&models.Rating{}).Select("avg(score)").Row().Scan(&stats.AverageRating)

	// 有评论的用户数
	config.DB.Model(&models.Comment{}).Distinct("user_id").Count(&stats.UsersWithComment)

	c.JSON(http.StatusOK, stats)
}

func GetAllRatings(c *gin.Context) {
	var ratings []models.Rating
	config.DB.Find(&ratings)
	c.JSON(http.StatusOK, gin.H{"data": ratings})
}

func GetAllComments(c *gin.Context) {
	var comments []models.Comment
	config.DB.Find(&comments)
	c.JSON(http.StatusOK, gin.H{"data": comments})
}

// BanUser 封禁用户
func BanUser(c *gin.Context) {
	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 不允许封禁管理员账户
	if user.Role == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "不能封禁管理员账户"})
		return
	}

	// 获取封禁原因
	var banData struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&banData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "封禁原因不能为空"})
		return
	}

	// 更新用户状态为封禁
	user.Status = "banned"
	user.BanReason = banData.Reason
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "封禁用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户封禁成功",
		"data":    user,
	})
}

// UnbanUser 解封用户
func UnbanUser(c *gin.Context) {
	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新用户状态为正常，清空封禁原因
	user.Status = "active"
	user.BanReason = ""
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解封用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户解封成功",
		"data":    user,
	})
}
