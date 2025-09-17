package main

import (
	"fmt"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"
	"xuan-ke-tong/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	config.ConnectDatabase()

	// 首先检查表是否存在，如果不存在则创建
	if err := ensureTablesExist(); err != nil {
		panic("数据库初始化失败: " + err.Error())
	}
	fmt.Println("数据库初始化成功")

	// GORM自动迁移（可选，确保模型同步）
	if err := config.DB.AutoMigrate(&models.User{}, &models.Course{}, &models.Rating{}, &models.Comment{}); err != nil {
		fmt.Printf("GORM自动迁移失败: %v\n", err)
	} else {
		fmt.Println("GORM数据库迁移成功")
	}

	// 添加种子数据
	seedData()

	routes.AuthRoutes(r)
	routes.CourseRoutes(r)
	routes.RatingRoutes(r)
	routes.CommentRoutes(r)
	routes.AdminRoutes(r)
	routes.UserRoutes(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 添加管理员路由日志
	r.GET("/admin-routes", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Admin routes are registered",
		})
	})

	// 打印所有注册的路由
	routes := r.Routes()
	println("Registered routes:")
	for _, route := range routes {
		println(route.Method, route.Path)
	}

	// 添加一个简单的测试路由
	r.POST("/api/v1/admin/courses/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Test route works"})
	})

	// 添加种子数据路由
	r.GET("/api/v1/seed", func(c *gin.Context) {
		seedData()
		c.JSON(200, gin.H{"message": "Seed data added successfully"})
	})

	r.Run(":8080") // listen and serve on 0.0.0:8080
}

// ensureTablesExist 确保数据库表存在
func ensureTablesExist() error {
	// 尝试手动创建表结构
	err := config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS courses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			grade TEXT,
			semester TEXT,
			subject TEXT,
			teacher TEXT,
			credits INTEGER,
			image_url TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	if err != nil {
		return fmt.Errorf("failed to create courses table: %v", err)
	}

	err = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			nickname TEXT,
			avatar TEXT,
			role TEXT DEFAULT 'user',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	err = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS ratings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			course_id INTEGER,
			score REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	if err != nil {
		return fmt.Errorf("failed to create ratings table: %v", err)
	}

	err = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			course_id INTEGER,
			content TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	if err != nil {
		return fmt.Errorf("failed to create comments table: %v", err)
	}

	return nil
}

func seedData() {
	// 创建管理员账号
	adminPassword := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Failed to hash admin password: %v\n", err)
	} else {
		adminUser := models.User{
			Username: "admin",
			Password: string(hashedPassword),
			Email:    "admin@example.com",
			Nickname: "Administrator",
			Role:     "admin",
		}
		// 使用 GORM 的 FirstOrCreate 来避免重复创建
		if err := config.DB.Where(models.User{Username: "admin"}).FirstOrCreate(&adminUser).Error; err != nil {
			fmt.Printf("Failed to create admin user: %v\n", err)
		} else {
			fmt.Println("Admin user created or already exists")
		}
	}

	// 清除现有数据
	config.DB.Exec("DELETE FROM ratings")
	config.DB.Exec("DELETE FROM comments")
	config.DB.Exec("DELETE FROM courses")

	// 创建测试课程
	courses := []models.Course{
		{
			Name:        "高等数学A",
			Description: "本课程主要介绍微积分、线性代数等数学基础知识，培养学生的数学思维能力和解决实际问题的能力。",
			Teacher:     "张教授",
			Credits:     4,
			Grade:       "大一",
			Semester:    "第一学期",
			Subject:     "数学",
			ImageURL:    "https://picsum.photos/seed/math1/400/200.jpg",
		},
		{
			Name:        "大学物理",
			Description: "涵盖力学、热学、电磁学、光学等基础物理知识，为后续专业课程打下坚实基础。",
			Teacher:     "李教授",
			Credits:     3,
			Grade:       "大一",
			Semester:    "第二学期",
			Subject:     "物理",
			ImageURL:    "https://picsum.photos/seed/physics1/400/200.jpg",
		},
		{
			Name:        "程序设计基础",
			Description: "学习C语言程序设计，掌握基本的编程思想和算法设计方法。",
			Teacher:     "王教授",
			Credits:     3,
			Grade:       "大一",
			Semester:    "第一学期",
			Subject:     "计算机",
			ImageURL:    "https://picsum.photos/seed/programming1/400/200.jpg",
		},
		{
			Name:        "数据结构",
			Description: "学习各种数据结构的原理和实现，包括线性表、树、图等，以及相关的算法设计。",
			Teacher:     "陈教授",
			Credits:     4,
			Grade:       "大二",
			Semester:    "第一学期",
			Subject:     "计算机",
			ImageURL:    "https://picsum.photos/seed/datastructure1/400/200.jpg",
		},
		{
			Name:        "英语听说",
			Description: "提高英语听力和口语表达能力，培养跨文化交际能力。",
			Teacher:     "Smith教授",
			Credits:     2,
			Grade:       "大一",
			Semester:    "第二学期",
			Subject:     "英语",
			ImageURL:    "https://picsum.photos/seed/english1/400/200.jpg",
		},
		{
			Name:        "线性代数",
			Description: "学习矩阵理论、线性方程组、向量空间等线性代数基础知识。",
			Teacher:     "赵教授",
			Credits:     3,
			Grade:       "大一",
			Semester:    "第二学期",
			Subject:     "数学",
			ImageURL:    "https://picsum.photos/seed/linearalgebra1/400/200.jpg",
		},
	}

	// 插入课程
	for i := range courses {
		if err := config.DB.Create(&courses[i]).Error; err != nil {
			fmt.Printf("创建课程失败: %v\n", err)
			continue
		}
		fmt.Printf("创建课程: %s\n", courses[i].Name)
	}

	// 获取创建的课程ID并为每个课程添加评分
	var createdCourses []models.Course
	if err := config.DB.Find(&createdCourses).Error; err != nil {
		fmt.Printf("获取课程失败: %v\n", err)
	} else {
		// 为每个课程创建评分
		ratingData := map[int][]float64{
			1: {4.5, 4.0, 4.8, 3.5, 4.2},  // 高等数学A的评分
			2: {3.8, 4.1, 3.9, 4.3},       // 大学物理的评分
			3: {4.9, 4.7, 4.8, 4.6, 4.5}, // 程序设计基础的评分
			4: {4.2, 4.4, 4.1},             // 数据结构的评分
			5: {3.5, 3.8, 3.2, 3.6},       // 英语听说的评分
			6: {4.0, 4.3, 3.9},             // 线性代数的评分
		}

		for _, course := range createdCourses {
			courseIndex := int(course.ID - 42) // 减去偏移量，因为课程ID从43开始
			if scores, exists := ratingData[courseIndex]; exists {
				for i, score := range scores {
					rating := models.Rating{
						UserID:   uint(i + 1), // 用户ID从1开始
						CourseID: course.ID,
						Score:    score,
					}
					if err := config.DB.Create(&rating).Error; err != nil {
						fmt.Printf("为课程%s创建评分失败: %v\n", course.Name, err)
					} else {
						fmt.Printf("为课程%s创建评分: %.1f\n", course.Name, score)
					}
				}
			}
		}
	}

	fmt.Println("测试数据添加完成！")
}
