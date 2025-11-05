package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/valyala/fastjson"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 全局 cron 调度器管理
var (
	bargainCronJobs = make(map[string]*cron.Cron) // taskID -> cron instance
	bargainMutex    sync.RWMutex
)

// createBargainTask 创建捡漏任务
func createBargainTask(req BargainTaskRequest, username string) (*BargainTaskDb, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}

	// 验证两个账号是否存在且属于当前用户
	var account1, account2 UserInfoDb
	if err := db.Where("id = ? AND user = ? AND is_delete = false", req.AccountID1, username).First(&account1).Error; err != nil {
		return nil, fmt.Errorf("账号1不存在或无权限")
	}
	if err := db.Where("id = ? AND user = ? AND is_delete = false", req.AccountID2, username).First(&account2).Error; err != nil {
		return nil, fmt.Errorf("账号2不存在或无权限")
	}

	// 验证预约日期必须在已开放时间内（当前时间之后）
	reservationDate, err := time.Parse("2006-01-02", req.ReservationDate)
	if err != nil {
		return nil, fmt.Errorf("预约日期格式错误")
	}
	if reservationDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, fmt.Errorf("预约日期必须在当前日期之后")
	}

	// 解析截止时间（可选）
	var deadline time.Time
	if req.Deadline != "" {
		deadline, err = time.Parse("2006-01-02 15:04:05", req.Deadline)
		if err != nil {
			return nil, fmt.Errorf("截止时间格式错误，应为 YYYY-MM-DD HH:mm:ss")
		}
		// 验证截止时间必须在当前时间之后
		if deadline.Before(time.Now()) {
			return nil, fmt.Errorf("截止时间必须在当前时间之后")
		}
		// 验证截止时间不能晚于预约日期的23:59:59
		maxDeadline := reservationDate.Add(24*time.Hour - time.Second)
		if deadline.After(maxDeadline) {
			return nil, fmt.Errorf("截止时间不能晚于预约日期当天")
		}
	}

	// 创建任务
	taskID := uuid.New().String()
	task := &BargainTaskDb{
		User:            username,
		TaskID:          taskID,
		AccountID1:      req.AccountID1,
		AccountID2:      req.AccountID2,
		VenueSiteID:     req.VenueSiteID,
		ReservationDate: req.ReservationDate,
		SiteName:        req.SiteName,
		ReservationTime: req.ReservationTime,
		ScanInterval:    req.ScanInterval,
		Deadline:        deadline,
		Status:          "active",
		SuccessCount:    0,
		ScanCount:       0,
	}

	if err := db.Create(task).Error; err != nil {
		return nil, fmt.Errorf("创建任务失败: %v", err)
	}

	// 启动定时扫描任务
	if err := startBargainScheduler(taskID); err != nil {
		log.Printf("启动定时任务失败: %v", err)
		// 不返回错误，任务已创建，可以稍后重试
	}

	return task, nil
}

// startBargainScheduler 启动捡漏任务调度器
func startBargainScheduler(taskID string) error {
	db, err := openDB()
	if err != nil {
		return err
	}

	var task BargainTaskDb
	if err := db.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		return fmt.Errorf("任务不存在")
	}

	if task.Status != "active" {
		return fmt.Errorf("任务状态不是active，无法启动")
	}

	// 创建 cron 实例
	c := cron.New()

	// 根据扫描间隔设置定时任务
	// 每N分钟执行一次
	cronSpec := fmt.Sprintf("*/%d * * * *", task.ScanInterval)

	_, err = c.AddFunc(cronSpec, func() {
		if err := scanAndReserve(taskID); err != nil {
			log.Printf("捡漏扫描失败 [TaskID: %s]: %v", taskID, err)
		}
	})

	if err != nil {
		return fmt.Errorf("添加定时任务失败: %v", err)
	}

	// 启动调度器
	c.Start()

	// 保存到全局map
	bargainMutex.Lock()
	bargainCronJobs[taskID] = c
	bargainMutex.Unlock()

	log.Printf("捡漏任务调度器已启动 [TaskID: %s, 间隔: %d分钟]", taskID, task.ScanInterval)
	return nil
}

// stopBargainScheduler 停止捡漏任务调度器
func stopBargainScheduler(taskID string) {
	bargainMutex.Lock()
	defer bargainMutex.Unlock()

	if c, exists := bargainCronJobs[taskID]; exists {
		c.Stop()
		delete(bargainCronJobs, taskID)
		log.Printf("捡漏任务调度器已停止 [TaskID: %s]", taskID)
	}
}

// scanAndReserve 扫描场地并尝试预约
func scanAndReserve(taskID string) error {
	db, err := openDB()
	if err != nil {
		return err
	}

	// 获取任务信息
	var task BargainTaskDb
	if err := db.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		return fmt.Errorf("任务不存在")
	}

	// 检查任务状态
	if task.Status != "active" {
		stopBargainScheduler(taskID)
		return fmt.Errorf("任务已停止或完成")
	}

	// 检查是否超过截止时间
	currentTime := time.Now()
	if !task.Deadline.IsZero() && currentTime.After(task.Deadline) {
		// 超过用户设置的截止时间
		failureReason := fmt.Sprintf("超过截止时间 %s，未能预约到场地", task.Deadline.Format("2006-01-02 15:04:05"))
		db.Model(&task).Updates(map[string]interface{}{
			"status":         "failed",
			"failure_reason": failureReason,
		})
		logBargainScan(taskID, 0, false, "任务失败", failureReason)
		stopBargainScheduler(taskID)
		log.Printf("捡漏任务失败 [TaskID: %s]: %s", taskID, failureReason)
		return fmt.Errorf(failureReason)
	}

	// 检查是否超过预约日期
	reservationDate, err := time.Parse("2006-01-02", task.ReservationDate)
	if err == nil {
		// 预约日期的结束时间是当天23:59:59
		reservationEndTime := reservationDate.Add(24*time.Hour - time.Second)
		if currentTime.After(reservationEndTime) {
			// 已经超过预约日期
			failureReason := fmt.Sprintf("已超过预约日期 %s，预约窗口已关闭", task.ReservationDate)
			db.Model(&task).Updates(map[string]interface{}{
				"status":         "failed",
				"failure_reason": failureReason,
			})
			logBargainScan(taskID, 0, false, "任务失败", failureReason)
			stopBargainScheduler(taskID)
			log.Printf("捡漏任务失败 [TaskID: %s]: %s", taskID, failureReason)
			return fmt.Errorf(failureReason)
		}
	}

	// 更新扫描次数和时间
	db.Model(&task).Updates(map[string]interface{}{
		"scan_count":     task.ScanCount + 1,
		"last_scan_time": time.Now(),
	})

	// 获取两个账号信息
	var account1, account2 UserInfoDb
	if err := db.Where("id = ?", task.AccountID1).First(&account1).Error; err != nil {
		logBargainScan(taskID, 0, false, "获取账号1失败", err.Error())
		return err
	}
	if err := db.Where("id = ?", task.AccountID2).First(&account2).Error; err != nil {
		logBargainScan(taskID, 0, false, "获取账号2失败", err.Error())
		return err
	}

	// 获取用户的验证码API
	var user UserDb
	if err := db.Where("username = ?", task.User).First(&user).Error; err != nil {
		logBargainScan(taskID, 0, false, "获取用户信息失败", err.Error())
		return err
	}

	log.Printf("开始扫描场地 [TaskID: %s, 日期: %s]", taskID, task.ReservationDate)

	// 使用第一个账号登录并获取场地信息
	client, err := NewClient(account1.Username, account1.Password)
	if err != nil {
		logBargainScan(taskID, 0, false, "账号1登录失败", err.Error())
		return err
	}

	// 获取场地信息
	venueInfo, err := fetchVenueInfo(client, task.ReservationDate, task.VenueSiteID)
	if err != nil {
		logBargainScan(taskID, 0, false, "获取场地信息失败", err.Error())
		return err
	}

	// 解析可用的场地和时间段
	availableSlots := findAvailableSlots(venueInfo, task.SiteName, task.ReservationTime)

	if len(availableSlots) == 0 {
		logBargainScan(taskID, 0, false, "没有可用场地", "")
		log.Printf("未找到可用场地 [TaskID: %s]", taskID)
		return nil
	}

	log.Printf("发现 %d 个可用场地 [TaskID: %s]", len(availableSlots), taskID)
	logBargainScan(taskID, len(availableSlots), false, fmt.Sprintf("发现%d个可用场地", len(availableSlots)), "")

	// 随机选择一个可用场地
	// 使用时间戳作为随机种子，确保每次扫描都有不同的随机结果
	randomIndex := rand.Intn(len(availableSlots))
	slot := availableSlots[randomIndex]
	log.Printf("随机选择场地 [索引: %d/%d, 场地: %s, 时间: %s]",
		randomIndex+1, len(availableSlots), slot.SiteName, slot.TimeSlot)

	// 获取账号2的BuddyNum和BuddyUserID
	// 需要登录账号2来获取其BuddyNum
	client2, buddyInfo, err := NewClient(account2.Username, account2.Password)
	if err != nil {
		logBargainScan(taskID, len(availableSlots), false, "账号2登录失败", err.Error())
		log.Printf("账号2登录失败 [TaskID: %s]: %v", taskID, err)
		return err
	}
	_ = client2 // 避免未使用变量警告

	if buddyInfo == nil || buddyInfo.BuddyNum == "" {
		errMsg := "无法获取账号2的同伴码"
		logBargainScan(taskID, len(availableSlots), false, errMsg, "")
		log.Printf("%s [TaskID: %s]", errMsg, taskID)
		return fmt.Errorf(errMsg)
	}

	log.Printf("获取同伴信息成功 [TaskID: %s, BuddyUserID: %d, BuddyNum: %s]",
		taskID, buddyInfo.UserId, buddyInfo.BuddyNum)

	// 使用账号1作为主预约账号，账号2提供同伴信息
	// 两个账号一起预约同一个场地
	success, err := executeReservation(account1, account2, buddyInfo, user.CaptchaAPI, task, slot)
	if success {
		// 预约成功，更新任务状态
		db.Model(&task).Updates(map[string]interface{}{
			"success_count": task.SuccessCount + 1,
			"status":        "completed",
		})
		logBargainScan(taskID, len(availableSlots), true,
			fmt.Sprintf("预约成功 [场地: %s, 时间: %s]", slot.SiteName, slot.TimeSlot), "")

		// 停止定时任务
		stopBargainScheduler(taskID)
		log.Printf("捡漏任务完成 [TaskID: %s, 场地: %s, 时间: %s]",
			taskID, slot.SiteName, slot.TimeSlot)
	} else {
		logBargainScan(taskID, len(availableSlots), false,
			fmt.Sprintf("预约失败: %v", err), "")
		log.Printf("预约失败 [TaskID: %s]: %v", taskID, err)
	}

	return nil
}

// VenueSlot 可用场地时间段
type VenueSlot struct {
	SiteID   string
	SiteName string
	TimeID   string
	TimeSlot string
}

// findAvailableSlots 查找可用的场地时间段
func findAvailableSlots(venueInfo *fastjson.Value, siteName, reservationTime string) []VenueSlot {
	var slots []VenueSlot

	data := venueInfo.Get("data")
	if data == nil {
		return slots
	}

	// 遍历所有场地
	spaces := data.GetArray("space")
	for _, space := range spaces {
		spaceName := string(space.GetStringBytes("spaceName"))
		spaceID := string(space.GetStringBytes("id"))

		// 如果指定了场地名称，只查找匹配的场地
		if siteName != "" && spaceName != siteName {
			continue
		}

		// 遍历该场地的所有时间段
		times := space.GetArray("timeInfo")
		for _, timeInfo := range times {
			status := timeInfo.GetInt("status")
			// status == 1 表示可预约
			if status == 1 {
				beginTime := string(timeInfo.GetStringBytes("beginTime"))
				endTime := string(timeInfo.GetStringBytes("endTime"))
				timeSlot := beginTime + "-" + endTime
				timeID := string(timeInfo.GetStringBytes("id"))

				// 如果指定了时间段，只查找匹配的时间段
				if reservationTime != "" && beginTime != reservationTime {
					continue
				}

				slots = append(slots, VenueSlot{
					SiteID:   spaceID,
					SiteName: spaceName,
					TimeID:   timeID,
					TimeSlot: timeSlot,
				})
			}
		}
	}

	return slots
}

// executeReservation 执行预约
// account1: 主预约账号
// account2: 同伴账号
// buddyInfo: 账号2的同伴信息（包含BuddyNum和UserId）
func executeReservation(account1 UserInfoDb, account2 UserInfoDb, buddyInfo *UserInfo, captchaAPI string, task BargainTaskDb, slot VenueSlot) (bool, error) {
	// 创建临时TaskInfo用于复用现有预约逻辑
	// 账号1作为主预约账号，账号2提供同伴码
	taskInfo := TaskInfoDb{
		User:            task.User,
		Username:        account1.Username,  // 主预约账号
		Password:        account1.Password,
		UserPhone:       account1.Phone,
		CaptchaAPI:      captchaAPI,
		BuddyUserID:     fmt.Sprintf("%d", buddyInfo.UserId), // 账号2的UserID
		BuddyNum:        buddyInfo.BuddyNum,                  // 账号2的BuddyNum
		VenueSiteID:     task.VenueSiteID,
		ReservationDate: task.ReservationDate,
		ReservationTime: slot.TimeSlot,
		SiteName:        slot.SiteName,
		TaskID:          task.TaskID + "_" + uuid.New().String()[:8],
		IsFinished:      false,
		InstantReservation: true,
		SiteId:          slot.SiteID,
		TimeId:          slot.TimeID,
	}

	log.Printf("开始预约 [主账号: %s, 同伴: %s, 场地: %s, 时间: %s]",
		account1.Username, account2.Username, slot.SiteName, slot.TimeSlot)

	// 使用现有的预约逻辑
	if err := tyysReserveTask(taskInfo, true); err != nil {
		return false, err
	}

	return true, nil
}

// logBargainScan 记录捡漏扫描日志
func logBargainScan(taskID string, availableSlots int, success bool, message, details string) {
	db, err := openDB()
	if err != nil {
		log.Printf("记录日志失败: %v", err)
		return
	}

	logEntry := BargainLogDb{
		TaskID:         taskID,
		AvailableSlots: availableSlots,
		Success:        success,
		Message:        message,
		Details:        details,
	}

	if err := db.Create(&logEntry).Error; err != nil {
		log.Printf("保存日志失败: %v", err)
	}
}

// getBargainTasksByUser 获取用户的所有捡漏任务
func getBargainTasksByUser(username string) ([]BargainTaskDb, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}

	var tasks []BargainTaskDb
	if err := db.Where("user = ?", username).Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

// getBargainTaskByID 根据ID获取捡漏任务
func getBargainTaskByID(taskID string, username string) (*BargainTaskDb, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}

	var task BargainTaskDb
	if err := db.Where("task_id = ? AND user = ?", taskID, username).First(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

// cancelBargainTask 取消捡漏任务
func cancelBargainTask(taskID string, username string) error {
	db, err := openDB()
	if err != nil {
		return err
	}

	var task BargainTaskDb
	if err := db.Where("task_id = ? AND user = ?", taskID, username).First(&task).Error; err != nil {
		return fmt.Errorf("任务不存在")
	}

	if task.Status == "cancelled" || task.Status == "completed" {
		return fmt.Errorf("任务已结束")
	}

	// 停止定时任务
	stopBargainScheduler(taskID)

	// 更新任务状态
	if err := db.Model(&task).Update("status", "cancelled").Error; err != nil {
		return err
	}

	logBargainScan(taskID, 0, false, "任务已取消", "")
	return nil
}

// getBargainLogs 获取捡漏任务的日志
func getBargainLogs(taskID string, username string) ([]BargainLogDb, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}

	// 验证任务是否属于该用户
	var task BargainTaskDb
	if err := db.Where("task_id = ? AND user = ?", taskID, username).First(&task).Error; err != nil {
		return nil, fmt.Errorf("任务不存在或无权限")
	}

	var logs []BargainLogDb
	if err := db.Where("task_id = ?", taskID).Order("scan_time DESC").Limit(100).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

// restartActiveBargainTasks 重启所有活跃的捡漏任务（用于服务器启动时）
func restartActiveBargainTasks() {
	db, err := openDB()
	if err != nil {
		log.Printf("重启捡漏任务失败: %v", err)
		return
	}

	var tasks []BargainTaskDb
	if err := db.Where("status = ?", "active").Find(&tasks).Error; err != nil {
		log.Printf("查询活跃任务失败: %v", err)
		return
	}

	for _, task := range tasks {
		if err := startBargainScheduler(task.TaskID); err != nil {
			log.Printf("重启任务失败 [TaskID: %s]: %v", task.TaskID, err)
		} else {
			log.Printf("任务已重启 [TaskID: %s]", task.TaskID)
		}
	}

	log.Printf("共重启 %d 个活跃的捡漏任务", len(tasks))
}

// openDB 打开数据库连接
func openDB() (*gorm.DB, error) {
	dbPath := "./database.db"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}
	return db, nil
}

// BargainTaskDetail 捡漏任务详细信息（包含账号信息）
type BargainTaskDetail struct {
	BargainTaskDb
	Account1Label string `json:"account_1_label"`
	Account2Label string `json:"account_2_label"`
}

// getBargainTaskDetail 获取捡漏任务详细信息
func getBargainTaskDetail(taskID string, username string) (*BargainTaskDetail, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}

	var task BargainTaskDb
	if err := db.Where("task_id = ? AND user = ?", taskID, username).First(&task).Error; err != nil {
		return nil, err
	}

	// 获取账号标签
	var account1, account2 UserInfoDb
	db.Where("id = ?", task.AccountID1).First(&account1)
	db.Where("id = ?", task.AccountID2).First(&account2)

	detail := &BargainTaskDetail{
		BargainTaskDb: task,
		Account1Label: account1.Lable,
		Account2Label: account2.Lable,
	}

	return detail, nil
}

// getAllBargainTasks 获取所有捡漏任务（管理员）
func getAllBargainTasks() ([]BargainTaskDb, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}

	var tasks []BargainTaskDb
	if err := db.Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}
