package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 创建一个sqlite数据库，名称database.db 创建一张表 task_info
func InitDb() {
	dbPath := "./database.db"
	// 检查数据库文件是否存在
	dbExists := fileExists(dbPath)

	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if !dbExists {
		fmt.Println("Database file not found, creating database...")
		// 数据库文件不存在时，创建表
		db.AutoMigrate(&TaskInfoDb{}, &UserInfoDb{}, &UserDb{}, &BargainTaskDb{}, &BargainLogDb{})
		fmt.Println("Tables created successfully")
	} else {
		fmt.Println("Database file exists, checking tables...")

		// 数据库文件存在，检查并创建表
		err := db.AutoMigrate(&TaskInfoDb{}, &UserInfoDb{}, &BargainTaskDb{}, &BargainLogDb{})
		if err != nil {
			fmt.Println("Error during table check:", err)
		} else {
			fmt.Println("Tables checked and created if necessary")
		}
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

// 根据 task_id 更新指定字段的值
func updateTask(taskID string, fieldName string, newValue any) error {
	dbPath := "./database.db"
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 更新指定 task_id 的指定字段的值
	err = db.Model(&TaskInfoDb{}).Where("task_id = ?", taskID).Update(fieldName, newValue).Error
	if err != nil {
		return err
	}

	return nil
}

// addTask 增加一条任务数据
func addTaskFromDB(taskInfo TaskInfo) error {
	dbPath := "./database.db"
	// 连接 SQLite 数据
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	instantReservation, err := getInstantReservation(taskInfo)
	if err != nil {
		return fmt.Errorf("error parsing reservation date: %v", err)
	}

	// 创建任务信息
	task := TaskInfoDb{
		User:               taskInfo.User,
		Username:           taskInfo.Username,
		Password:           taskInfo.Password,
		UserPhone:          taskInfo.UserPhone,
		CaptchaAPI:         taskInfo.CaptchaApi,
		BuddyUserID:        taskInfo.BuddyUserId,
		BuddyNum:           taskInfo.BuddyNum,
		VenueSiteID:        taskInfo.VenueSiteId,
		ReservationDate:    taskInfo.ReservationDate,
		ReservationTime:    taskInfo.ReservationTime,
		SiteName:           taskInfo.SiteName,
		TaskID:             taskInfo.TaskId,
		CreateTime:         time.Now(),
		IsFinished:         false,
		InstantReservation: instantReservation,
		Autocancel:         false,
		ReservationStatus:  true,
	}

	// 插入任务信息到数据库
	result := db.Create(&task)
	if result.Error != nil {
		fmt.Println("Error adding task:", result.Error)
	} else {
		log.Println("Task added successfully")
	}
	return nil
}

func getInstantReservation(taskInfo TaskInfo) (bool, error) {
	// 获取当前时间戳（秒）
	currentTimestamp := time.Now().Unix()

	// 加载中国时区
	cst, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return false, fmt.Errorf("error loading timezone: %v", err)
	}

	// 解析 ReservationDate
	reservationDate, err := time.ParseInLocation("2006-01-02", taskInfo.ReservationDate, cst)
	if err != nil {
		return false, fmt.Errorf("error parsing reservation date: %v", err)
	}

	// 设置预约日期的早上 9:00（中国时间）
	reservationDateTime := time.Date(
		reservationDate.Year(),
		reservationDate.Month(),
		reservationDate.Day(),
		9, 0, 0, 0,
		cst,
	)
	reservationTimestamp := reservationDateTime.Unix()

	// 计算时间差（秒）
	timeDiff := reservationTimestamp - currentTimestamp

	// 判断时间差是否大于 48 小时
	if timeDiff < 48*3600 {
		return true, nil
	} else {
		return false, nil
	}
}

func createUser(username, password string) error {
	dbPath := "./database.db"
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	isAdmin := false
	if username == "Moonstone" {
		isAdmin = true
	}
	// 向表user_db中插入一条数据
	user := UserDb{
		Username: username,
		Password: password,
		IsAdmin:  isAdmin,
		IsDelete: false,
	}
	db.Create(&user)
	return nil
}

func checkUser(username, password string) ([]TaskInfoDb, string, bool, error) {
	dbPath := "./database.db"
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 查询表user_db中是否存在该用户名,如果不存在返回信息：用户名不存在，如果存在
	var user UserDb
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, "", false, fmt.Errorf("error querying database: %v", result.Error)
	}
	if user.ID == 0 {
		return nil, "", false, fmt.Errorf("用户名不存在")
	} else {
		captchaAPI := user.CaptchaAPI
		if user.Password == password {
			// 查询用户是否为管理员
			isAdmin := user.IsAdmin
			// 通过用户名查询task_info_db表中的数据
			taskInfos, err := getTaskListFromDb(username)
			if err != nil {
				return nil, captchaAPI, false, fmt.Errorf("error querying database: %v", err)
			}
			return taskInfos, captchaAPI, isAdmin, nil
		} else {
			return nil, "", false, fmt.Errorf("密码错误")
		}
	}
}

func getUsers() ([]UserDb, error) {
	dbPath := "./database.db"
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	// 查询表user_db中的所有用户
	var users []UserDb
	result := db.Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("error querying database: %v", result.Error)
	}

	return users, nil
}

// getTaskListFromDb 获取指定用户的所有任务数据
func getTaskListFromDb(username string) ([]TaskInfoDb, error) {
	dbPath := "./database.db"
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	// 查询表task_info_db
	var taskInfos []TaskInfoDb
	// 返回所有的数据条目
	result := db.Where("user = ?", username).Find(&taskInfos)
	if result.Error != nil {
		return nil, fmt.Errorf("error querying database: %v", result.Error)
	}

	return taskInfos, nil
}

// getAllTaskListFromDb 获取所有用户的所有任务数据（管理员）
func getAllTaskListFromDb() ([]TaskInfoDb, error) {
	dbPath := "./database.db"
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	// 查询所有任务，按创建时间倒序
	var taskInfos []TaskInfoDb
	result := db.Order("create_time desc").Find(&taskInfos)
	if result.Error != nil {
		return nil, fmt.Errorf("error querying database: %v", result.Error)
	}

	return taskInfos, nil
}

// addRevserUser 增加一条预约用户数据
func addRevserUser(lable, username, password, user string) error {
	dbPath := "./database.db"
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}
	// 检查是否存在该用户名
	var userInfo UserInfoDb
	result := db.Where("username = ?", username).First(&userInfo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 用户不存在，可以创建新用户
			newUser := UserInfoDb{
				Lable:    lable,
				Username: username,
				Password: password,
				IsDelete: false,
				User:     user,
			}
			if err := db.Create(&newUser).Error; err != nil {
				return fmt.Errorf("error creating new user: %v", err)
			}
			return nil
		}
		// 其他查询错误
		return fmt.Errorf("error querying database: %v", result.Error)
	}

	// 用户存在，检查是否被标记为删除
	if userInfo.IsDelete {
		userInfo.IsDelete = false
		userInfo.Lable = lable
		userInfo.Password = password
		userInfo.User = user
		if err := db.Save(&userInfo).Error; err != nil {
			return fmt.Errorf("error reactivating user: %v", err)
		}
		return nil
	}

	// 用户存在且未被删除
	return fmt.Errorf("用户名已存在")
}

// updateRevserUser 更新一条预约用户数据
func updateRevserUser(username string, fieldName string, newValue any) error {
	dbPath := "./database.db"
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 更新指定 task_id 的指定字段的值
	err = db.Model(&UserInfoDb{}).Where("username = ?", username).Update(fieldName, newValue).Error
	if err != nil {
		return err
	}

	return nil
}

// getRevserUser 获取某位用户的所有预约用户数据
func getRevserUser(user string) ([]UserInfoDb, error) {
	dbPath := "./database.db"
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	// 查询表user_info_db
	var userInfos []UserInfoDb
	// 返回所的数据条目，条件是user匹配且is_delete为false
	result := db.Where("user = ? AND is_delete = ?", user, false).Find(&userInfos)
	if result.Error != nil {
		return nil, fmt.Errorf("error querying database: %v", result.Error)
	}

	return userInfos, nil
}

// updateCaptchaAPI 更新指定用户的验证码API
func updateCaptchaAPIFromDB(user, captchaAPI string) error {
	dbPath := "./database.db"
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	// 更新指定用户的验证码API
	err = db.Model(&UserDb{}).Where("username = ?", user).Update("captcha_api", captchaAPI).Error
	if err != nil {
		return fmt.Errorf("error updating captcha API: %v", err)
	}
	return nil
}

// changePassword 修改指定用户的密码
func changePassword(username, password, newPassword string) error {
	dbPath := "./database.db"
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}
	// 检查旧密码是否正确
	var user UserDb
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return fmt.Errorf("error querying database: %v", result.Error)
	}
	if user.Password != password {
		return fmt.Errorf("旧密码错误")
	}

	// 更新指定用户的密码
	err = db.Model(&UserDb{}).Where("username = ?", username).Update("password", newPassword).Error
	if err != nil {
		return fmt.Errorf("error updating password: %v", err)
	}
	return nil
}

// getTaskInfoByTaskId 根据taskId查询task_info_db表中的数据
func getTaskInfoByTaskId(taskId string) (TaskInfoDb, error) {
	dbPath := "./database.db"
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return TaskInfoDb{}, fmt.Errorf("failed to connect database: %v", err)
	}

	// 查询表task_info_db
	var taskInfo TaskInfoDb
	result := db.Where("task_id = ?", taskId).First(&taskInfo)
	if result.Error != nil {
		return TaskInfoDb{}, fmt.Errorf("error querying database: %v", result.Error)
	}

	return taskInfo, nil
}

// cancelTask 取消指定任务
func cancelTaskFromDB(taskId string) error {
	dbPath := "./database.db"
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	// 查询表task_info_db
	var taskInfo TaskInfoDb
	result := db.Where("task_id = ?", taskId).First(&taskInfo)
	if result.Error != nil {
		return fmt.Errorf("error querying database: %v", result.Error)
	}

	// 登录
	client, _, err := NewClient(taskInfo.Username, taskInfo.Password)
	if err != nil {
		return fmt.Errorf("error creating client: %v", err)
	}

	// 取消订单
	dataFormCancel := map[string]string{
		"venueTradeNo": taskInfo.TradeNo,
	}
	cancelResponse, cancelErr := cancel(client, dataFormCancel)
	if cancelErr != nil {
		return fmt.Errorf("error canceling task: %v", cancelErr)
	}
	if string(cancelResponse.GetStringBytes("message")) == "success" {
		// 更新数据库中的数据
		updateTaskInfoErr := updateTask(taskId, "reservation_status", false) //false 表示取消预约
		if updateTaskInfoErr != nil {
			return fmt.Errorf("更新订单状态失败：%v", updateTaskInfoErr)
		}
	}
	return nil
}
