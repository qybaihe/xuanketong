package controllers

import (
	"net/http"
	"strconv"
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

// GetAllUsers 获取所有用户列表（带分页和搜索）
func GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	role := c.Query("role")

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
		AverageRating     float64 `json:"average_rating"`
		UsersWithComment  int64   `json:"users_with_comment"`
	}

	// 总用户数
	config.DB.Model(&models.User{}).Count(&stats.TotalUsers)

	// 管理员数
	config.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&stats.TotalAdmins)

	// 普通用户数
	config.DB.Model(&models.User{}).Where("role = ?", "user").Count(&stats.TotalRegularUsers)

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

