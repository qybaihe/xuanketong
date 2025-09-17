//go:build !cgo

package database

import (
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func ConnectDatabaseSqlite() *gorm.DB {
	// Create data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		panic("Failed to create data directory: " + err.Error())
	}

	dbPath := filepath.Join(dataDir, EnvOrDefault("XKT_SQLITE_DB_NAME", "test.db"))
	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	return database
}
