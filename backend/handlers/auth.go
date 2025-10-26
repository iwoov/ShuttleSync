package handlers

import (
	"errors"
	"fmt"
	"time"

	"shuttlesync/auth"
	"shuttlesync/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// UserDb 用户数据库模型（与 main 包中的定义一致）
type UserDb struct {
	ID         uint      `gorm:"primaryKey"`
	Username   string    `gorm:"not null;unique"`
	Password   string    `gorm:"not null"`
	CaptchaAPI string    `gorm:"not null"`
	IsAdmin    bool      `gorm:"not null;default:false"`
	IsDelete   bool      `gorm:"not null;default:false"`
	CreatedAt  time.Time `gorm:"autoCreateTime;not null"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshRequest 刷新 token 请求结构
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest 登出请求结构
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
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

// Login 用户登录接口
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 连接数据库
	db, err := getDB()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "数据库连接失败",
			"error":   err.Error(),
		})
		return
	}

	// 查询用户
	var user UserDb
	result := db.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{
				"message": "用户名不存在",
			})
			return
		}
		c.JSON(500, gin.H{
			"message": "查询用户失败",
			"error":   result.Error.Error(),
		})
		return
	}

	// 验证密码
	if user.Password != req.Password {
		c.JSON(200, gin.H{
			"message": "密码错误",
		})
		return
	}

	// 生成 Access Token
	accessToken, err := auth.GenerateAccessToken(user.Username, user.IsAdmin)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "生成 Access Token 失败",
			"error":   err.Error(),
		})
		return
	}

	// 生成 Refresh Token
	refreshToken, err := auth.GenerateRefreshToken(user.Username)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "生成 Refresh Token 失败",
			"error":   err.Error(),
		})
		return
	}

	// 保存 Refresh Token 到数据库
	expiresAt := time.Now().Add(auth.RefreshTokenExpiry)
	err = models.SaveRefreshToken(user.ID, user.Username, refreshToken, expiresAt)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "保存 Refresh Token 失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回 tokens 和用户信息
	c.JSON(200, gin.H{
		"message": "success",
		"data": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"username":      user.Username,
			"is_admin":      user.IsAdmin,
			"captcha_api":   user.CaptchaAPI,
			"expires_in":    int(auth.AccessTokenExpiry.Seconds()),
		},
	})
}

// RefreshToken 刷新 access token 接口
func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证 Refresh Token 是否存在且有效
	refreshTokenRecord, err := models.GetRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Refresh Token 无效或已过期",
			"error":   err.Error(),
		})
		return
	}

	// 解析 Refresh Token 获取用户信息
	claims, err := auth.ParseToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Refresh Token 解析失败",
			"error":   err.Error(),
		})
		return
	}

	// 获取用户信息（需要 is_admin 字段）
	db, err := getDB()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "数据库连接失败",
			"error":   err.Error(),
		})
		return
	}

	var user UserDb
	result := db.Where("username = ?", claims.Username).First(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "查询用户失败",
			"error":   result.Error.Error(),
		})
		return
	}

	// 生成新的 Access Token
	newAccessToken, err := auth.GenerateAccessToken(user.Username, user.IsAdmin)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "生成新 Access Token 失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data": gin.H{
			"access_token": newAccessToken,
			"username":     refreshTokenRecord.Username,
			"expires_in":   int(auth.AccessTokenExpiry.Seconds()),
		},
	})
}

// Logout 用户登出接口
func Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 撤销 Refresh Token
	err := models.RevokeRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "登出失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data":    "登出成功",
	})
}

// GetCurrentUser 获取当前登录用户信息（需要鉴权）
func GetCurrentUser(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{
			"message": "未认证",
		})
		return
	}

	isAdmin, _ := c.Get("is_admin")

	// 从数据库获取完整用户信息
	db, err := getDB()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "数据库连接失败",
		})
		return
	}

	var user UserDb
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "查询用户失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data": gin.H{
			"username":    user.Username,
			"is_admin":    isAdmin,
			"captcha_api": user.CaptchaAPI,
			"created_at":  user.CreatedAt,
		},
	})
}

// CleanupExpiredTokens 清理过期的 token（可以通过定时任务或管理接口调用）
func CleanupExpiredTokens(c *gin.Context) {
	err := models.CleanExpiredTokens()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "清理过期 token 失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data":    "清理完成",
	})
}

// 内部函数：根据用户名获取用户ID（用于其他handler）
func GetUserIDByUsername(username string) (uint, error) {
	db, err := getDB()
	if err != nil {
		return 0, err
	}

	var user UserDb
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return 0, fmt.Errorf("用户不存在")
	}

	return user.ID, nil
}
