package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "未提供认证令牌",
				"code":    "NO_TOKEN",
			})
			c.Abort()
			return
		}

		// 检查 Bearer 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "认证令牌格式错误",
				"code":    "INVALID_TOKEN_FORMAT",
			})
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "认证令牌无效或已过期",
				"code":    "INVALID_TOKEN",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.IsAdmin)

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件（需要在 AuthMiddleware 之后使用）
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取 is_admin
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "需要管理员权限",
				"code":    "FORBIDDEN",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware 可选的鉴权中间件（Token 有效则注入用户信息，无效也继续）
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				claims, err := ParseToken(parts[1])
				if err == nil {
					c.Set("username", claims.Username)
					c.Set("is_admin", claims.IsAdmin)
				}
			}
		}
		c.Next()
	}
}
