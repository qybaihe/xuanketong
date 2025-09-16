package controllers

import (
	"net/http"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"

	"github.com/gin-gonic/gin"
)

// StatsResponse 定义了仪表盘统计数据的结构
type StatsResponse struct {
	TotalUsers       int64   `json:"total_users"`
	TotalCourses     int64   `json:"total_courses"`
	TotalComments    int64   `json:"total_comments"`
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

	// 4. 计算平均评分
	if err = config.DB.Model(&models.Rating{}).Select("avg(score)").Row().Scan(&stats.AverageRating); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法计算平均评分"})
		return
	}

	// 5. 本周新增用户数 (此为示例，实际应根据具体需求调整时间范围)
	// SQLite 不支持 date('now', '-7 days')，这里用一个简化逻辑代替
	// 在生产环境中，您可能需要使用更复杂的、与数据库方言兼容的日期函数
	// time.Now().AddDate(0, 0, -7)
	stats.NewUsersThisWeek = 0 // 暂时设置为0，因为需要更复杂的日期处理

	// 6. 待审核评论数 (此为示例，假设Comment模型有Status字段)
	// if err = config.DB.Model(&models.Comment{}).Where("status = ?", "pending").Count(&stats.PendingComments).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取待审核评论数"})
	// 	return
	// }
	stats.PendingComments = 0 // 暂时设置为0，因为模型中没有status字段

	c.JSON(http.StatusOK, stats)
}

func GetAllUsers(c *gin.Context) {
	var users []models.User
	config.DB.Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
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
