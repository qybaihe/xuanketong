package config

import (
	"log"
	"os"
	"path/filepath"
	"xuan-ke-tong/models"

	// 暂先使用不依赖cgo的SQLite驱动，解决Windows适配问题
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Create data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		panic("Failed to create data directory: " + err.Error())
	}

	dbPath := filepath.Join(dataDir, "test.db")
	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// 自动迁移
	// err = database.AutoMigrate(&models.User{}, &models.Course{}, &models.Rating{}, &models.Comment{}, &models.EvaluationRequest{})
	if err := database.AutoMigrate(&models.User{}); err != nil {
		log.Printf("GORM迁移User失败: %v", err)
	}
	if err := database.AutoMigrate(&models.Course{}); err != nil {
		log.Printf("GORM迁移Course失败: %v", err)
	}
	if err := database.AutoMigrate(&models.Rating{}); err != nil {
		log.Printf("GORM迁移Rating失败: %v", err)
	}
	if err := database.AutoMigrate(&models.Comment{}); err != nil {
		log.Printf("GORM迁移Comment失败: %v", err)
	}
	if err := database.AutoMigrate(&models.EvaluationRequest{}); err != nil {
		log.Printf("GORM迁移EvaluationRequest失败: %v", err)
	}

	DB = database
}
