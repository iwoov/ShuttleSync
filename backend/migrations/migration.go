package migrations

import (
	"fmt"
	"os"

	"shuttlesync/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// RunMigrations 执行数据库迁移
func RunMigrations() {
	dbPath := "./database.db"

	// 检查数据库文件是否存在
	dbExists := fileExists(dbPath)

	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if !dbExists {
		fmt.Println("Database file not found, creating database with all tables...")
	} else {
		fmt.Println("Database file exists, running migrations...")
	}

	// 自动迁移 RefreshToken 表（复用 models 包中的定义，避免字段不一致）
	err = db.AutoMigrate(&models.RefreshToken{})
	if err != nil {
		fmt.Println("Error migrating RefreshToken table:", err)
	} else {
		fmt.Println("RefreshToken table migrated successfully")
	}
}

// fileExists 检查文件是否存在
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
