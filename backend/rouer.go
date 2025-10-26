package main

import (
	"shuttlesync/auth"
	"shuttlesync/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func router() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Disable automatic redirects that may cause 301 for SPA paths
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	// 使用默认的CORS配置
	router.Use(cors.Default())

	// 创建路由组/api
	api := router.Group("/api")

	{
		// ========== 认证相关路由（无需鉴权） ==========
		authGroup := api.Group("/auth")
		{
			// 用户登录
			authGroup.POST("/login", handlers.Login)
			// 刷新 access token
			authGroup.POST("/refresh", handlers.RefreshToken)
			// 用户登出
			authGroup.POST("/logout", handlers.Logout)
			// 获取当前用户信息（需要鉴权）
			authGroup.GET("/me", auth.AuthMiddleware(), handlers.GetCurrentUser)
		}

		// ========== 用户路由组 ==========
		user := api.Group("/user")
		{
			// 网站用户-用户注册（无需鉴权）
			user.POST("/register", userRegister)

			// 以下接口需要鉴权
			user.Use(auth.AuthMiddleware())

			// 网站用户-所有用户（需要管理员权限）
			user.GET("/all", auth.AdminMiddleware(), getAllUser)
			// 网站用户-更新密码
			user.PATCH("/password", updateUserPassword)
			// 网站用户-更新验证码api
			user.PATCH("/captcha_api", updateCaptchaAPI)
		}

		// ========== 预约账号路由组（需要鉴权） ==========
		account := api.Group("/account")
		account.Use(auth.AuthMiddleware())
		{
			// 预约账号-获取账号信息
			account.GET("/list", getReservationAccountList)
			// 预约账号-添加账号
			account.POST("/add", addReservationAccount)
			// 预约账号-更新账号
			account.PATCH("/update", updateReservationAccount)
			// 预约账号-删除账号
			account.DELETE("/delete", deleteReservationAccount)
		}

		// ========== 预约任务路由组（需要鉴权） ==========
		task := api.Group("/task")
		task.Use(auth.AuthMiddleware())
		{
			// 预约任务-获取任务
			task.GET("/list", getTaskList)
			// 预约任务-提交预约
			task.POST("/add", addTask)
			// 预约任务-预约取消
			task.GET("/cancel", cancelTask)
		}

		// ========== 浙大体艺路由组（需要鉴权） ==========
		tyys := api.Group("/tyys")
		tyys.Use(auth.AuthMiddleware())
		{
			// 浙大体艺-预约码获取
			tyys.POST("/buddy_num", getBuddyNum)
			// 浙大体艺-网站登录
			tyys.POST("/login", tyysLogin)
			// 浙大体艺-预约码获取
			tyys.GET("/qr_code", getReservationCode)
		}

	}

	// Serve embedded frontend (built assets copied to backend/web)
	registerFrontendRoutes(router)

	router.Run(":4050")
}
