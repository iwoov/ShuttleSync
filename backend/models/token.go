package models

import (
	"errors"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// RefreshToken 定义 refresh token 的数据库模型
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`  // 关联用户ID
	Username  string    `gorm:"not null;index" json:"username"` // 用户名（冗余字段，方便查询）
	Token     string    `gorm:"not null;unique" json:"token"`   // Refresh Token 值
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`     // 过期时间
	CreatedAt time.Time `gorm:"autoCreateTime;not null" json:"created_at"`
	IsRevoked bool      `gorm:"not null;default:false" json:"is_revoked"` // 是否已撤销
}

// 获取数据库连接
func getDB() (*gorm.DB, error) {
	dbPath := "./database.db"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database")
	}
	return db, nil
}

// SaveRefreshToken 保存 refresh token 到数据库
func SaveRefreshToken(userID uint, username, token string, expiresAt time.Time) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	refreshToken := RefreshToken{
		UserID:    userID,
		Username:  username,
		Token:     token,
		ExpiresAt: expiresAt,
		IsRevoked: false,
	}

	return db.Create(&refreshToken).Error
}

// GetRefreshToken 根据 token 值获取记录
func GetRefreshToken(token string) (*RefreshToken, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}

	var refreshToken RefreshToken
	result := db.Where("token = ? AND is_revoked = ?", token, false).First(&refreshToken)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token not found or revoked")
		}
		return nil, result.Error
	}

	// 检查是否过期
	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	return &refreshToken, nil
}

// RevokeRefreshToken 撤销 refresh token（登出时使用）
func RevokeRefreshToken(token string) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	return db.Model(&RefreshToken{}).Where("token = ?", token).Update("is_revoked", true).Error
}

// RevokeAllUserTokens 撤销用户的所有 refresh token（例如修改密码后）
func RevokeAllUserTokens(username string) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	return db.Model(&RefreshToken{}).Where("username = ?", username).Update("is_revoked", true).Error
}

// CleanExpiredTokens 清理过期的 token（可以定期执行）
func CleanExpiredTokens() error {
	db, err := getDB()
	if err != nil {
		return err
	}

	return db.Where("expires_at < ?", time.Now()).Delete(&RefreshToken{}).Error
}
