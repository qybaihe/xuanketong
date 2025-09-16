package main

import (
	"fmt"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"
	"xuan-ke-tong/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{}, &models.Course{}, &models.Rating{}, &models.Comment{})

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

func seedData() {
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

	// 为课程添加评分
	ratings := []models.Rating{
		{UserID: 1, CourseID: 1, Score: 4.5},
		{UserID: 2, CourseID: 1, Score: 4.0},
		{UserID: 3, CourseID: 1, Score: 4.8},
		{UserID: 4, CourseID: 1, Score: 3.5},
		{UserID: 5, CourseID: 1, Score: 4.2},

		{UserID: 1, CourseID: 2, Score: 3.8},
		{UserID: 2, CourseID: 2, Score: 4.1},
		{UserID: 3, CourseID: 2, Score: 3.9},
		{UserID: 4, CourseID: 2, Score: 4.3},

		{UserID: 1, CourseID: 3, Score: 4.9},
		{UserID: 2, CourseID: 3, Score: 4.7},
		{UserID: 3, CourseID: 3, Score: 4.8},
		{UserID: 4, CourseID: 3, Score: 4.6},
		{UserID: 5, CourseID: 3, Score: 4.5},

		{UserID: 1, CourseID: 4, Score: 4.2},
		{UserID: 2, CourseID: 4, Score: 4.4},
		{UserID: 3, CourseID: 4, Score: 4.1},

		{UserID: 1, CourseID: 5, Score: 3.5},
		{UserID: 2, CourseID: 5, Score: 3.8},
		{UserID: 3, CourseID: 5, Score: 3.2},
		{UserID: 4, CourseID: 5, Score: 3.6},

		{UserID: 1, CourseID: 6, Score: 4.0},
		{UserID: 2, CourseID: 6, Score: 4.3},
		{UserID: 3, CourseID: 6, Score: 3.9},
	}

	// 插入评分
	for i := range ratings {
		if err := config.DB.Create(&ratings[i]).Error; err != nil {
			fmt.Printf("创建评分失败: %v\n", err)
			continue
		}
		fmt.Printf("为课程ID %d 创建评分: %.1f\n", ratings[i].CourseID, ratings[i].Score)
	}

	fmt.Println("测试数据添加完成！")
}
