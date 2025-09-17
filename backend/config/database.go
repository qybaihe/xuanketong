package config

import (
	"gorm.io/gorm"

	"xuan-ke-tong/config/database"
)

var DB *gorm.DB

var dbConnectors = map[string]func() *gorm.DB{
	"sqlite": database.ConnectDatabaseSqlite,
}

func ConnectDatabase() {
	dbType := database.EnvOrDefault("XKT_DB_TYPE", "sqlite")

	if connector, exists := dbConnectors[dbType]; exists {
		DB = connector()
	} else {
		panic("Unsupported database type: " + dbType)
	}
}
