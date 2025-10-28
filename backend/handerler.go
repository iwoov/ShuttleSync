package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 用户注册 /api/user/register
func userRegister(c *gin.Context) {
	var userInfo UserDb
	err := c.ShouldBindJSON(&userInfo)
	createUserErr := createUser(userInfo.Username, userInfo.Password)
	if createUserErr != nil {
		c.JSON(200, gin.H{
			"message": "注册失败",
			"data":    createUserErr.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    "注册成功!",
	})
	if err != nil {
		c.JSON(200, gin.H{
			"message": "注册失败",
			"data":    err.Error(),
		})
	}
}

// 用户登录 /api/user/login
func userLogin(c *gin.Context) {
	var userInfo UserDb
	c.ShouldBindJSON(&userInfo)
	taskInfos, captchaAPI, isAdmin, checkUserErr := checkUser(userInfo.Username, userInfo.Password)
	if checkUserErr != nil {
		c.JSON(200, gin.H{
			"message": checkUserErr.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data": gin.H{
			"username":   userInfo.Username,
			"captchaAPI": captchaAPI,
			"isAdmin":    isAdmin,
			"taskInfos":  taskInfos,
		},
	})
}

// 更新验证码API /api/user/captcha_api
func updateCaptchaAPI(c *gin.Context) {
	// 从上下文获取当前登录用户
	username, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{
			"message": "未认证",
		})
		return
	}

	var requestData struct {
		CaptchaAPI string `json:"captcha_api"`
	}
	c.ShouldBindJSON(&requestData)

	fmt.Println(username, requestData.CaptchaAPI)
	updateCaptchaAPIErr := updateCaptchaAPIFromDB(username.(string), requestData.CaptchaAPI)
	if updateCaptchaAPIErr != nil {
		c.JSON(200, gin.H{
			"message": "更新验证码API失败",
			"data":    updateCaptchaAPIErr.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    "更新验证码API成功!",
	})
}

// 更新密码  /api/user/change_password
func updateUserPassword(c *gin.Context) {
	// 从上下文获取当前登录用户
	username, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{
			"message": "未认证",
		})
		return
	}

	var changePasswordInfo struct {
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}
	c.ShouldBind(&changePasswordInfo)

	changePasswordErr := changePassword(username.(string), changePasswordInfo.Password, changePasswordInfo.NewPassword)
	if changePasswordErr != nil {
		c.JSON(200, gin.H{
			"message": "修改密码失败",
			"data":    changePasswordErr.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    "修改密码成功!",
	})
}

// 获取所有网站用户 /api/user/all
func getAllUser(c *gin.Context) {
	usersList, getUsersErr := getUsers()
	if getUsersErr != nil {
		c.JSON(200, gin.H{
			"message": getUsersErr.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    usersList,
	})
}

// /api/account/
// 添加预约账户 /api/account/add
func addReservationAccount(c *gin.Context) {
	// 从上下文获取当前登录用户
	currentUser, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{
			"message": "未认证",
		})
		return
	}

	var reservationAccount struct {
		Lable    string `json:"lable"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	c.ShouldBind(&reservationAccount)

	addRevserErr := addRevserUser(reservationAccount.Lable, reservationAccount.Username, reservationAccount.Password, currentUser.(string))
	if addRevserErr != nil {
		c.JSON(200, gin.H{
			"message": "添加失败",
			"data":    addRevserErr.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    "添加预约用户成功!",
	})
}

// 获取预约账户列表 /api/account/list
func getReservationAccountList(c *gin.Context) {
	// 从上下文获取当前登录用户
	currentUser, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{
			"message": "未认证",
		})
		return
	}

	revserUsers, getRevserUserErr := getRevserUser(currentUser.(string))
	if getRevserUserErr != nil {
		c.JSON(200, gin.H{
			"message": "获取预约用户列表失败",
			"data":    getRevserUserErr.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    revserUsers,
	})
}

// 更新预约账户信息 /api/account/update
func updateReservationAccount(c *gin.Context) {
	var reservationAccount ReservationAccount
	c.ShouldBind(&reservationAccount)
	updateRevserUserErr := updateRevserUser(reservationAccount.Username, "lable", reservationAccount.Lable)
	if updateRevserUserErr != nil {
		c.JSON(200, gin.H{
			"message": "更新预约用户Lable失败",
			"data":    updateRevserUserErr.Error(),
		})
		return
	}
	updateRevserUserErr = updateRevserUser(reservationAccount.Username, "password", reservationAccount.Password)
	if updateRevserUserErr != nil {
		c.JSON(200, gin.H{
			"message": "更新预约用户Password失败",
			"data":    updateRevserUserErr.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    "更新预约用户信息成功!",
	})
}

// 删除预约账户信息 /api/account/delete
func deleteReservationAccount(c *gin.Context) {
	var reservationAccount ReservationAccount
	c.ShouldBind(&reservationAccount)
	deleteRevserErr := updateRevserUser(reservationAccount.Username, "IsDelete", true)
	if deleteRevserErr != nil {
		c.JSON(200, gin.H{
			"message": "删除预约用户失败",
			"data":    deleteRevserErr.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    "删除预约用户成功!",
	})
}

// /api/task/
// 获取预约任务列表 /api/task/list
func getTaskList(c *gin.Context) {
	// 从上下文获取当前登录用户
	currentUser, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{
			"message": "未认证",
		})
		return
	}

	taskInfoList, getTaskErr := getTaskListFromDb(currentUser.(string))
	if getTaskErr != nil {
		c.JSON(200, gin.H{
			"message": getTaskErr.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data": gin.H{
			"username":  currentUser.(string),
			"taskInfos": taskInfoList,
		},
	})
}

// 获取所有用户的预约任务（管理员） /api/task/all
func getAllTaskList(c *gin.Context) {
	// 仅作为安全兜底：路由已加 AdminMiddleware，这里再检查一次可选的上下文标记
	if isAdminVal, exists := c.Get("is_admin"); !exists || !isAdminVal.(bool) {
		c.JSON(403, gin.H{
			"message": "需要管理员权限",
		})
		return
	}

	taskInfoList, getTaskErr := getAllTaskListFromDb()
	if getTaskErr != nil {
		c.JSON(200, gin.H{
			"message": getTaskErr.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data": gin.H{
			"taskInfos": taskInfoList,
		},
	})
}

// 提交预约 /api/task/add
func addTask(c *gin.Context) {
	// 从上下文获取当前登录用户
	currentUser, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{
			"message": "未认证",
		})
		return
	}

	var taskInfo TaskInfo
	c.ShouldBind(&taskInfo)
	// 确保任务归属于当前登录用户
	taskInfo.User = currentUser.(string)

	taskResultErr := addReserveTask(taskInfo)
	if taskResultErr != nil {
		c.JSON(200, gin.H{
			"message": "预约添加失败",
			"data":    taskResultErr.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    "预约提交成功!",
	})
}

// 取消任务 /api/task/cancel
func cancelTask(c *gin.Context) {
	taskId := c.Query("task_id")
	cancelTaskErr := cancelTaskFromDB(taskId)
	if cancelTaskErr != nil {
		c.JSON(200, gin.H{
			"message": "取消任务失败",
			"data":    cancelTaskErr.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    "取消任务成功!",
	})
}

// /api/tyys
// 登录浙大体艺网站 /api/tyys/login
func tyysLogin(c *gin.Context) {
	var tyysAccount TyysAccount
	c.ShouldBind(&tyysAccount)
	_, userInfo, err := NewClient(tyysAccount.Username, tyysAccount.Password)
	// 失败兜底：err 或 userInfo 为空 或 userId 为 0
	if err != nil || userInfo == nil || userInfo.UserId == 0 {
		var dataMsg string
		if err != nil {
			dataMsg = err.Error()
		} else {
			dataMsg = "用户名或密码错误"
		}
		c.JSON(200, gin.H{
			"message": "登录失败",
			"data":    dataMsg,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    userInfo,
	})
}

// 获取同伴码 /api/tyys/buddy_num
func getBuddyNum(c *gin.Context) {
	var tyysAccount TyysAccount
	c.ShouldBind(&tyysAccount)
	_, BuddyInfo, err := NewClient(tyysAccount.Username, tyysAccount.Password)
	// 失败兜底：err 或 BuddyInfo 为空 或 buddyNum 为空
	if err != nil || BuddyInfo == nil || BuddyInfo.BuddyNum == "" {
		var dataMsg string
		if err != nil {
			dataMsg = err.Error()
		} else {
			dataMsg = "登录失败或未获取到同伴码"
		}
		c.JSON(200, gin.H{
			"message": "获取同伴码失败",
			"data":    dataMsg,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    map[string]string{"buddy_num": BuddyInfo.BuddyNum, "buddy_id": strconv.FormatInt(BuddyInfo.UserId, 10)},
	})
}

// 获取预约码 /api/tyys/qr_code
func getReservationCode(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	orderId := c.Query("order_id")
	orderCodeBase64, orderCodeErr := tyysGetOrderCode(username, password, orderId)
	if orderCodeErr != nil {
		c.JSON(200, gin.H{
			"message": "查询预约码失败",
			"data":    orderCodeErr.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    orderCodeBase64,
	})
}
