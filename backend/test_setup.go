//go:build ignore

package main

import (
	"fmt"
	"xuan-ke-tong/config"
	"xuan-ke-tong/models"
)

func setupTestData() error {
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

	// 添加测试课程
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
	}

	// 插入课程
	for _, course := range courses {
		if err := config.DB.Create(&course).Error; err != nil {
			fmt.Printf("创建课程失败: %v\n", err)
			continue
		}
		fmt.Printf("创建课程: %s\n", course.Name)
	}

	fmt.Println("测试数据设置完成！")
	return nil
}
