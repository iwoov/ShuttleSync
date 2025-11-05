package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// createBargainTaskHandler 创建捡漏任务
// POST /api/bargain/create
func createBargainTaskHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "未授权",
		})
		return
	}

	var req BargainTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证两个账号不能相同
	if req.AccountID1 == req.AccountID2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "两个预约账号不能相同",
		})
		return
	}

	task, err := createBargainTask(req, username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "创建任务失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "捡漏任务创建成功",
		"data":    task,
	})
}

// getBargainTasksHandler 获取用户的捡漏任务列表
// GET /api/bargain/list
func getBargainTasksHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "未授权",
		})
		return
	}

	tasks, err := getBargainTasksByUser(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取任务列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    tasks,
	})
}

// getBargainTaskDetailHandler 获取捡漏任务详情
// GET /api/bargain/:id
func getBargainTaskDetailHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "未授权",
		})
		return
	}

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "任务ID不能为空",
		})
		return
	}

	task, err := getBargainTaskDetail(taskID, username.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "任务不存在",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    task,
	})
}

// cancelBargainTaskHandler 取消捡漏任务
// DELETE /api/bargain/:id
func cancelBargainTaskHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "未授权",
		})
		return
	}

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "任务ID不能为空",
		})
		return
	}

	if err := cancelBargainTask(taskID, username.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "取消任务失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "任务已取消",
	})
}

// getBargainLogsHandler 获取捡漏任务的扫描日志
// GET /api/bargain/:id/logs
func getBargainLogsHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "未授权",
		})
		return
	}

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "任务ID不能为空",
		})
		return
	}

	logs, err := getBargainLogs(taskID, username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取日志失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    logs,
	})
}

// getAllBargainTasksHandler 获取所有捡漏任务（管理员）
// GET /api/bargain/all
func getAllBargainTasksHandler(c *gin.Context) {
	isAdmin, exists := c.Get("is_admin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "需要管理员权限",
		})
		return
	}

	tasks, err := getAllBargainTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取任务列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    tasks,
	})
}

// updateBargainTaskHandler 更新捡漏任务
// PUT /api/bargain/:id
func updateBargainTaskHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "未授权",
		})
		return
	}

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "任务ID不能为空",
		})
		return
	}

	var req BargainTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证两个账号不能相同
	if req.AccountID1 == req.AccountID2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "两个预约账号不能相同",
		})
		return
	}

	if err := updateBargainTask(taskID, req, username.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "更新任务失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "任务更新成功",
	})
}
