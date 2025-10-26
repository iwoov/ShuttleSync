package migrations

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// RefreshToken 定义 refresh token 的数据库模型（与 models/token.go 中的定义一致）
type RefreshToken struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	Username  string `gorm:"not null;index"`
	Token     string `gorm:"not null;unique"`
	ExpiresAt int64  `gorm:"not null"` // 使用 Unix 时间戳
	CreatedAt int64  `gorm:"autoCreateTime;not null"`
	IsRevoked bool   `gorm:"not null;default:false"`
}

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

	// 自动迁移 RefreshToken 表
	err = db.AutoMigrate(&RefreshToken{})
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
