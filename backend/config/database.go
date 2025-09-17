package config

import (
	"os"
	"path/filepath"

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
	DB = database
}
