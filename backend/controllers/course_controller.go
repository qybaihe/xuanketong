package controllers

import (
	"net/http"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"

	"github.com/gin-gonic/gin"
)

func CreateCourse(c *gin.Context) {
    var input struct {
        Name           string  `json:"Name"`
        Description    string  `json:"Description"`
        Grade          string  `json:"Grade"`
        Semester       string  `json:"Semester"`
        Subject        string  `json:"Subject"`
        Teacher        string  `json:"Teacher"`
        Credits        float64 `json:"Credits"`
        ImageURL       string  `json:"ImageURL"`
        CourseUnitName string  `json:"CourseUnitName"`
        CourseCategory string  `json:"CourseCategory"`
    }

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    course := models.Course{
        Name:           input.Name,
        Description:    input.Description,
        Grade:          input.Grade,
        Semester:       input.Semester,
        Subject:        input.Subject,
        Teacher:        input.Teacher,
        Credits:        input.Credits,
        ImageURL:       input.ImageURL,
        CourseUnitName: input.CourseUnitName,
        CourseCategory: input.CourseCategory,
    }

	if err := config.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": course})
}

func GetCourses(c *gin.Context) {
	var courses []models.Course
	query := config.DB

	if grade := c.Query("grade"); grade != "" {
		query = query.Where("grade = ?", grade)
	}
	if semester := c.Query("semester"); semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if subject := c.Query("subject"); subject != "" {
		query = query.Where("subject = ?", subject)
	}
	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if teacher := c.Query("teacher"); teacher != "" {
		query = query.Where("teacher LIKE ?", "%"+teacher+"%")
	}

	if err := query.Unscoped().Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get courses"})
		return
	}

	// 获取每个课程的平均评分和评分分布
	type CourseWithRating struct {
		models.Course
		AverageRating      float64     `json:"averageRating"`
		TotalRatings       int         `json:"totalRatings"`
		AverageDifficulty  float64     `json:"averageDifficulty"`  // 平均难度
		RatingDistribution map[int]int `json:"ratingDistribution"` // 1-5星评分分布
		MaterialsCount     int         `json:"materialsCount"`     // 资料数量
		HotComments        []string    `json:"hotComments"`        // 热门评论
		HotMaterials       []string    `json:"hotMaterials"`       // 热门资料
	}

	var coursesWithRatings []CourseWithRating

	for _, course := range courses {
		var avgRating float64
		var avgDifficulty float64
		var totalRatings int64

		// 计算平均评分和平均难度
		config.DB.Model(&models.Rating{}).
			Where("course_id = ?", course.ID).
			Select("AVG(score) as avg_rating, AVG(difficulty) as avg_difficulty, COUNT(*) as total").
			Row().
			Scan(&avgRating, &avgDifficulty, &totalRatings)

		if avgRating == 0 {
			avgRating = 0
		}
		if avgDifficulty == 0 {
			avgDifficulty = 0
		}

		// 获取评分分布（1-5星）
		ratingDistribution := make(map[int]int)
		for i := 1; i <= 5; i++ {
			var count int64
			config.DB.Model(&models.Rating{}).
				Where("course_id = ? AND score >= ? AND score < ?", course.ID, float64(i)-0.5, float64(i)+0.5).
				Count(&count)
			ratingDistribution[i] = int(count)
		}

		// 获取资料数量
		var materialsCount int64
		config.DB.Model(&models.CourseMaterial{}).
			Where("course_id = ? AND is_active = ?", course.ID, true).
			Count(&materialsCount)

		// 获取热门评论（按评分排序，取前3条有内容的评论）
		var hotComments []string
		var ratings []models.Rating
		config.DB.Where("course_id = ? AND content != '' AND content IS NOT NULL", course.ID).
			Order("score DESC, created_at DESC").
			Limit(3).
			Find(&ratings)

		for _, rating := range ratings {
			if len(rating.Content) > 50 {
				hotComments = append(hotComments, rating.Content[:50]+"...")
			} else {
				hotComments = append(hotComments, rating.Content)
			}
		}

		// 获取热门资料（按下载量排序，取前3个）
		var hotMaterials []string
		var materials []models.CourseMaterial
		config.DB.Where("course_id = ? AND is_active = ?", course.ID, true).
			Order("download_count DESC").
			Limit(3).
			Find(&materials)

		for _, material := range materials {
			hotMaterials = append(hotMaterials, material.MaterialName)
		}

		coursesWithRatings = append(coursesWithRatings, CourseWithRating{
			Course:             course,
			AverageRating:      avgRating,
			TotalRatings:       int(totalRatings),
			AverageDifficulty:  avgDifficulty,
			RatingDistribution: ratingDistribution,
			MaterialsCount:     int(materialsCount),
			HotComments:        hotComments,
			HotMaterials:       hotMaterials,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": coursesWithRatings})
}

func GetCourse(c *gin.Context) {
	var course models.Course
	if err := config.DB.Where("id = ?", c.Param("id")).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": course})
}

func UpdateCourse(c *gin.Context) {
	var course models.Course
	if err := config.DB.First(&course, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	var input models.Course
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 只更新允许的字段
    updateData := map[string]interface{}{
		"Name":           input.Name,
		"Description":    input.Description,
		"Grade":          input.Grade,
		"Semester":       input.Semester,
		"Subject":        input.Subject,
		"Teacher":        input.Teacher,
		"Credits":        input.Credits,
		"ImageURL":       input.ImageURL,
		"CourseUnitName": input.CourseUnitName,
        "CourseCategory": input.CourseCategory,
	}

	if err := config.DB.Model(&course).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": course})
}

func DeleteCourse(c *gin.Context) {
	var course models.Course
	if err := config.DB.First(&course, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	if err := config.DB.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

// BatchCreateCourses 批量创建课程（管理员）
func BatchCreateCourses(c *gin.Context) {
    var inputs []struct {
        Name           string  `json:"Name"`
        Description    string  `json:"Description"`
        Grade          string  `json:"Grade"`
        Semester       string  `json:"Semester"`
        Subject        string  `json:"Subject"`
        Teacher        string  `json:"Teacher"`
        Credits        float64 `json:"Credits"`
        ImageURL       string  `json:"ImageURL"`
        CourseUnitName string  `json:"CourseUnitName"`
        CourseCategory string  `json:"CourseCategory"`
    }

	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(inputs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty payload"})
		return
	}

	courses := make([]models.Course, 0, len(inputs))
	for _, in := range inputs {
        courses = append(courses, models.Course{
			Name:           in.Name,
			Description:    in.Description,
			Grade:          in.Grade,
			Semester:       in.Semester,
			Subject:        in.Subject,
			Teacher:        in.Teacher,
			Credits:        in.Credits,
			ImageURL:       in.ImageURL,
			CourseUnitName: in.CourseUnitName,
            CourseCategory: in.CourseCategory,
		})
	}

	tx := config.DB.Begin()
	if err := tx.Create(&courses).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to batch create courses"})
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit batch create"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": courses, "count": len(courses)})
}
