package controllers

import (
	"net/http"
	"time"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"

	"github.com/gin-gonic/gin"
)

// EnhancedHomeStatsResponse 定义了增强版首页统计数据的结构
type EnhancedHomeStatsResponse struct {
	OverviewData       OverviewStats        `json:"overview_data"`
	TopRatedCourses    []TopRatedCourse     `json:"top_rated_courses"`
	RecentCourses      []RecentCourse       `json:"recent_courses"`
	PopularCourses     []PopularCourse      `json:"popular_courses"`
	UserActivityStats  UserActivityResponse `json:"user_activity_stats"`
	CourseDistribution CourseDistribution   `json:"course_distribution"`
	MonthlyStats       []MonthlyStat        `json:"monthly_stats"`
}

// OverviewStats 总体统计
type OverviewStats struct {
	TotalCourses        int64   `json:"total_courses"`
	TotalUsers          int64   `json:"total_users"`
	TotalRatings        int64   `json:"total_ratings"`
	TotalComments       int64   `json:"total_comments"`
	AverageRating       float64 `json:"average_rating"`
	ActiveUsersThisWeek int64   `json:"active_users_this_week"`
}

// TopRatedCourse 最高评分课程
type TopRatedCourse struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Teacher       string  `json:"teacher"`
	AverageRating float64 `json:"average_rating"`
	TotalRatings  int64   `json:"total_ratings"`
	ImageURL      string  `json:"image_url"`
	Subject       string  `json:"subject"`
	Grade         string  `json:"grade"`
}

// RecentCourse 最新课程
type RecentCourse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Teacher     string    `json:"teacher"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Subject     string    `json:"subject"`
	Grade       string    `json:"grade"`
	CreatedAt   time.Time `json:"created_at"`
}

// PopularCourse 热门课程
type PopularCourse struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	Teacher           string  `json:"teacher"`
	AverageRating     float64 `json:"average_rating"`
	TotalRatings      int64   `json:"total_ratings"`
	ImageURL          string  `json:"image_url"`
	Subject           string  `json:"subject"`
	Grade             string  `json:"grade"`
	StudentEngagement int64   `json:"student_engagement"` // 学生参与度指标
}

// CourseDistribution 课程分布统计
type CourseDistribution struct {
	BySubject  map[string]int64 `json:"by_subject"`
	ByGrade    map[string]int64 `json:"by_grade"`
	BySemester map[string]int64 `json:"by_semester"`
}

// MonthlyStat 月度统计
type MonthlyStat struct {
	Month      string `json:"month"`       // 月份，格式：2024-01
	Courses    int64  `json:"courses"`     // 新增课程数
	Users      int64  `json:"users"`       // 新增用户数
	Ratings    int64  `json:"ratings"`     // 新增评分数
	Comments   int64  `json:"comments"`    // 新增评论数
	TotalScore int64  `json:"total_score"` // 总评分（用于计算平均）
}

// GetEnhancedHomeStats 获取增强版首页统计数据
func GetEnhancedHomeStats(c *gin.Context) {
	var response EnhancedHomeStatsResponse

	// 1. 获取总体统计数据
	if err := getOverviewStats(&response.OverviewData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取总体统计数据失败"})
		return
	}

	// 2. 获取评分最高的课程（前6个）
	if err := getTopRatedCourses(&response.TopRatedCourses, 6); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评分最高课程失败"})
		return
	}

	// 3. 获取最新课程（前6个）
	if err := getRecentCourses(&response.RecentCourses, 6); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取最新课程失败"})
		return
	}

	// 4. 获取最受欢迎课程（前6个）
	if err := getPopularCourses(&response.PopularCourses, 6); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取热门课程失败"})
		return
	}

	// 5. 获取用户活动统计
	if err := getUserActivityStats(&response.UserActivityStats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户活动统计失败"})
		return
	}

	// 6. 获取课程分布统计
	if err := getCourseDistribution(&response.CourseDistribution); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取课程分布统计失败"})
		return
	}

	// 7. 获取月度统计（最近6个月）
	if err := getMonthlyStats(&response.MonthlyStats, 6); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取月度统计失败"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// getOverviewStats 获取总体统计数据
func getOverviewStats(stats *OverviewStats) error {
	// 总课程数
	if err := config.DB.Model(&models.Course{}).Count(&stats.TotalCourses).Error; err != nil {
		return err
	}

	// 总用户数
	if err := config.DB.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return err
	}

	// 总评分数
	if err := config.DB.Model(&models.Rating{}).Count(&stats.TotalRatings).Error; err != nil {
		return err
	}

	// 总评论数
	if err := config.DB.Model(&models.Comment{}).Count(&stats.TotalComments).Error; err != nil {
		return err
	}

	// 平均评分
	config.DB.Model(&models.Rating{}).Select("avg(score)").Row().Scan(&stats.AverageRating)

	// 本周活跃用户数（基于有评分或评论的用户）
	var activeUsers int64
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	// 查询最近7天内有评分或评论的用户数
	config.DB.Raw(`
		SELECT COUNT(DISTINCT user_id) as active_users
		FROM (
			SELECT user_id FROM ratings WHERE created_at > ?
			UNION
			SELECT user_id FROM comments WHERE created_at > ?
		)
	`, sevenDaysAgo, sevenDaysAgo).Row().Scan(&activeUsers)

	stats.ActiveUsersThisWeek = activeUsers

	return nil
}

// getTopRatedCourses 获取评分最高的课程
func getTopRatedCourses(courses *[]TopRatedCourse, limit int) error {
	err := config.DB.Table("courses").
		Select("courses.id, courses.name, courses.teacher, courses.image_url, courses.subject, courses.grade, AVG(ratings.score) as average_rating, COUNT(ratings.id) as total_ratings").
		Joins("LEFT JOIN ratings ON courses.id = ratings.course_id").
		Group("courses.id").
		Having("COUNT(ratings.id) > 0").
		Order("average_rating DESC").
		Limit(limit).
		Scan(courses).Error

	return err
}

// getRecentCourses 获取最新课程
func getRecentCourses(courses *[]RecentCourse, limit int) error {
	err := config.DB.Table("courses").
		Select("courses.id, courses.name, courses.teacher, courses.description, courses.image_url, courses.subject, courses.grade, courses.created_at").
		Order("courses.created_at DESC").
		Limit(limit).
		Scan(courses).Error

	return err
}

// getPopularCourses 获取热门课程（基于评分数量和平均评分的综合指标）
func getPopularCourses(courses *[]PopularCourse, limit int) error {
	err := config.DB.Table("courses").
		Select(`
			courses.id, 
			courses.name, 
			courses.teacher, 
			courses.image_url, 
			courses.subject, 
			courses.grade, 
			AVG(ratings.score) as average_rating, 
			COUNT(ratings.id) as total_ratings,
			(COUNT(ratings.id) * AVG(ratings.score) + 10) as student_engagement
		`).
		Joins("LEFT JOIN ratings ON courses.id = ratings.course_id").
		Group("courses.id").
		Having("COUNT(ratings.id) > 0").
		Order("student_engagement DESC, average_rating DESC").
		Limit(limit).
		Scan(courses).Error

	return err
}

// getUserActivityStats 获取用户活动统计
func getUserActivityStats(stats *UserActivityResponse) error {
	// 调用现有的用户活动统计函数
	return GetUserActivityData(stats)
}

// getCourseDistribution 获取课程分布统计
func getCourseDistribution(distribution *CourseDistribution) error {
	// 按科目分布
	bySubject := make(map[string]int64)
	rows, err := config.DB.Model(&models.Course{}).
		Select("subject, COUNT(*) as count").
		Group("subject").
		Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var subject string
		var count int64
		if err := rows.Scan(&subject, &count); err == nil {
			bySubject[subject] = count
		}
	}

	// 按年级分布
	byGrade := make(map[string]int64)
	rows, err = config.DB.Model(&models.Course{}).
		Select("grade, COUNT(*) as count").
		Group("grade").
		Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var grade string
		var count int64
		if err := rows.Scan(&grade, &count); err == nil {
			byGrade[grade] = count
		}
	}

	// 按学期分布
	bySemester := make(map[string]int64)
	rows, err = config.DB.Model(&models.Course{}).
		Select("semester, COUNT(*) as count").
		Group("semester").
		Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var semester string
		var count int64
		if err := rows.Scan(&semester, &count); err == nil {
			bySemester[semester] = count
		}
	}

	distribution.BySubject = bySubject
	distribution.ByGrade = byGrade
	distribution.BySemester = bySemester

	return nil
}

// getMonthlyStats 获取月度统计数据
func getMonthlyStats(stats *[]MonthlyStat, months int) error {
	// 获取最近6个月的月度统计
	for i := months - 1; i >= 0; i-- {
		month := time.Now().AddDate(0, -i, 0)
		monthStr := month.Format("2006-01")

		var stat MonthlyStat
		stat.Month = monthStr

		// 计算当月新增课程数
		firstDay := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
		lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Second)

		config.DB.Model(&models.Course{}).
			Where("created_at BETWEEN ? AND ?", firstDay, lastDay).
			Count(&stat.Courses)

		// 计算当月新增用户数
		config.DB.Model(&models.User{}).
			Where("created_at BETWEEN ? AND ?", firstDay, lastDay).
			Count(&stat.Users)

		// 计算当月新增评分数
		config.DB.Model(&models.Rating{}).
			Where("created_at BETWEEN ? AND ?", firstDay, lastDay).
			Count(&stat.Ratings)

		// 计算当月新增评论数
		config.DB.Model(&models.Comment{}).
			Where("created_at BETWEEN ? AND ?", firstDay, lastDay).
			Count(&stat.Comments)

		// 计算当月的总评分（用于后续计算平均值）
		var totalScore float64
		config.DB.Model(&models.Rating{}).
			Where("created_at BETWEEN ? AND ?", firstDay, lastDay).
			Select("COALESCE(SUM(score), 0)").Row().Scan(&totalScore)
		stat.TotalScore = int64(totalScore)

		*stats = append(*stats, stat)
	}

	return nil
}

// UserActivityResponse 用户活动统计响应
type UserActivityResponse struct {
	ActiveUsersToday     int64 `json:"active_users_today"`
	ActiveUsersThisWeek  int64 `json:"active_users_this_week"`
	ActiveUsersThisMonth int64 `json:"active_users_this_month"`
	NewUsersToday        int64 `json:"new_users_today"`
	NewUsersThisWeek     int64 `json:"new_users_this_week"`
	NewUsersThisMonth    int64 `json:"new_users_this_month"`
}

// UserActivityData 获取用户活动数据 - 这是供内部调用的实际实现
func UserActivityData(response *UserActivityResponse) error {
	now := time.Now()

	// 今日活跃用户（有评分或评论的用户）
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var todayActive int64
	config.DB.Raw(`
		SELECT COUNT(DISTINCT user_id) FROM (
			SELECT user_id FROM ratings WHERE created_at > ?
			UNION  
			SELECT user_id FROM comments WHERE created_at > ?
		)
	`, todayStart, todayStart).Row().Scan(&todayActive)
	response.ActiveUsersToday = todayActive

	// 本周活跃用户
	weekStart := now.AddDate(0, 0, -7)
	var weekActive int64
	config.DB.Raw(`
		SELECT COUNT(DISTINCT user_id) FROM (
			SELECT user_id FROM ratings WHERE created_at > ?
			UNION  
			SELECT user_id FROM comments WHERE created_at > ?
		)
	`, weekStart, weekStart).Row().Scan(&weekActive)
	response.ActiveUsersThisWeek = weekActive

	// 本月活跃用户
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var monthActive int64
	config.DB.Raw(`
		SELECT COUNT(DISTINCT user_id) FROM (
			SELECT user_id FROM ratings WHERE created_at > ?
			UNION  
			SELECT user_id FROM comments WHERE created_at > ?
		)
	`, monthStart, monthStart).Row().Scan(&monthActive)
	response.ActiveUsersThisMonth = monthActive

	// 今日新用户
	var todayNew int64
	config.DB.Model(&models.User{}).Where("created_at > ?", todayStart).Count(&todayNew)
	response.NewUsersToday = todayNew

	// 本周新用户
	var weekNew int64
	config.DB.Model(&models.User{}).Where("created_at > ?", weekStart).Count(&weekNew)
	response.NewUsersThisWeek = weekNew

	// 本月新用户
	var monthNew int64
	config.DB.Model(&models.User{}).Where("created_at > ?", monthStart).Count(&monthNew)
	response.NewUsersThisMonth = monthNew

	return nil
}

// GetUserActivityData 供其他控制器调用的用户活动数据函数
func GetUserActivityData(response *UserActivityResponse) error {
	return UserActivityData(response)
}
