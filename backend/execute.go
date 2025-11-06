package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/robfig/cron/v3"
	"github.com/valyala/fastjson"
)

// getBuddyId 获取同伴id 不同于userid
func getBuddyId(client *resty.Client, taskInfo *TaskInfoDb) error {
	// 获取同伴id 不同于userid
	buddiesListResp, _ := fetchBuddiesList(client)
	buddiesList := buddiesListResp.Get("data").GetArray("content")
	for _, buddy := range buddiesList {
		if strconv.FormatInt(buddy.GetInt64("userId"), 10) == taskInfo.BuddyUserID {
			taskInfo.BuddyId = strconv.FormatInt(buddy.GetInt64("id"), 10)
			return nil
		}
	}
	return fmt.Errorf("获取同伴id失败")
}

// getVenueInfo 获取场地信息
func getVenueInfo(client *resty.Client, taskInfo *TaskInfoDb) error {
	venueSiteTimeInfoRes, err := fetchVenueInfo(client, taskInfo.ReservationDate, taskInfo.VenueSiteID)
	if err != nil {
		return fmt.Errorf("获取场地信息失败：%v", err)
	}
	taskInfo.OrderToken = string(venueSiteTimeInfoRes.Get("data").GetStringBytes("token"))
	// 通过场地名称获取场地id
	siteId, getSiteIdErr := getSiteId(venueSiteTimeInfoRes, *taskInfo)
	if getSiteIdErr != nil {
		return fmt.Errorf("获取场地id失败：%v", getSiteIdErr)
	}
	taskInfo.SiteId = siteId
	// 通过预约时间获取时间id
	timeId, getTimeIdErr := getTimeId(venueSiteTimeInfoRes, *taskInfo)
	if getTimeIdErr != nil {
		return fmt.Errorf("获取时间id失败：%v", getTimeIdErr)
	}
	taskInfo.TimeId = timeId
	return nil
}

// 通过场地名称获取场地id
func getSiteId(venueSiteTimeInfoRes *fastjson.Value, taskInfo TaskInfoDb) (string, error) {
	siteInfoList := venueSiteTimeInfoRes.Get("data").Get("reservationDateSpaceInfo").GetArray(taskInfo.ReservationDate)
	for _, siteInfo := range siteInfoList {
		siteName := string(siteInfo.GetStringBytes("spaceName"))
		if siteName == taskInfo.SiteName {
			siteId := strconv.FormatInt(siteInfo.GetInt64("id"), 10) // 将 int64 转换为 string
			return siteId, nil
		}
	}
	return "", fmt.Errorf("获取场地id失败")
}

// 通过预约时间获取时间id
func getTimeId(venueSiteTimeInfoRes *fastjson.Value, taskInfo TaskInfoDb) (string, error) {
	spaceTimeInfoList := venueSiteTimeInfoRes.Get("data").GetArray("spaceTimeInfo")
	for _, spaceTimeInfo := range spaceTimeInfoList {
		timeBegin := string(spaceTimeInfo.GetStringBytes("beginTime")) // 获取字符串并去掉引号
		timeId := spaceTimeInfo.GetInt64("id")
		if timeBegin == taskInfo.ReservationTime {
			return fmt.Sprintf("%d", timeId), nil // 将 int64 转换为 string
		}
	}
	return "", fmt.Errorf("获取时间id失败")
}

func confirmVenue(client *resty.Client, taskInfo *TaskInfoDb) error {
	reservationOrderJson := fmt.Sprintf("[{\"spaceId\":%s,\"timeId\":%s,\"venueSpaceGroupId\":null}]", taskInfo.SiteId, taskInfo.TimeId)
	dataFormVenue := map[string]string{
		"venueSiteId":          taskInfo.VenueSiteID,
		"reservationDate":      taskInfo.ReservationDate,
		"weekStartDate":        taskInfo.ReservationDate,
		"reservationOrderJson": reservationOrderJson, // 5435 18-19 5436 19-20 5437 20-21 5438 21-22
		"token":                taskInfo.OrderToken,
	}
	_, fetchVenueErr := fetchVenue(client, dataFormVenue)
	if fetchVenueErr != nil {
		return fmt.Errorf("获取场地信息失败：%v", fetchVenueErr)
	}
	return nil
}

func getCaptchaToken(client *resty.Client, taskInfo *TaskInfoDb) error {
	// 获取验证码
	captchaRes, fetchCaptchaErr := fetchCaptcha(client)
	if fetchCaptchaErr != nil {
		return fmt.Errorf("获取验证码失败：%v", fetchCaptchaErr)
	}
	// 获取验证码的 secretKey 和 token
	secretKey := string(captchaRes.GetStringBytes("data", "repData", "secretKey"))
	CaptchaToken := string(captchaRes.GetStringBytes("data", "repData", "token"))
	// 获取验证码图片的base64
	originalImageBase64 := string(captchaRes.GetStringBytes("data", "repData", "originalImageBase64"))
	// 获取验证码的坐标，第三方打码平台返回数据
	positionJsonStr, getCaptchaErr := Captcha(originalImageBase64, getCaptchaWord(captchaRes), taskInfo.CaptchaAPI)
	if getCaptchaErr != nil {
		return fmt.Errorf("验证码打码失败：%v", getCaptchaErr)
	}
	// 将坐标转换为加密数据
	pointJson, getPointJsonErr := getPointJson(positionJsonStr, secretKey)
	if getPointJsonErr != nil {
		return fmt.Errorf("获取验证码坐标失败：%v", getPointJsonErr)
	}
	// 提交验证码

	checkCaptchaResponse, checkCaptchaErr := checkCaptchaResult(client, pointJson, CaptchaToken)
	if checkCaptchaErr != nil {
		return fmt.Errorf("验证码校验失败：%v", checkCaptchaErr)
	}
	//判断验证码是否正确
	if string(checkCaptchaResponse.GetStringBytes("message")) == "OK" {
		captchaVerification, getCaptchaVerificationErr := GetCaptchaVerification(secretKey, CaptchaToken, positionJsonStr)
		if getCaptchaVerificationErr != nil {
			return fmt.Errorf("获取验证码认证失败：%v", getCaptchaVerificationErr)
		}
		taskInfo.CaptchaVerification = captchaVerification
	} else {
		fmt.Println("验证码校验失败！重试...")
	}
	return nil
}

func checkCaptchaResult(client *resty.Client, pointJson string, CaptchaToken string) (*fastjson.Value, error) {
	dataFormCheck := map[string]string{
		"captchaType": "clickWord",
		"pointJson":   pointJson,
		"token":       CaptchaToken,
	}
	checkCaptchaResponse, _ := checkCaptcha(client, dataFormCheck)
	return checkCaptchaResponse, nil
}

// GetCaptchaWord 获取Word
func getCaptchaWord(GetCaptchaRes *fastjson.Value) string {
	wordList := GetCaptchaRes.GetArray("data", "repData", "wordList")
	// 将 []*Value 切片转换为逗号分隔的字符串
	var builder strings.Builder
	for i, v := range wordList {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(strings.Trim(v.String(), `"`))
	}
	wordListString := builder.String()
	return wordListString
}

func submitOrder(client *resty.Client, taskInfo *TaskInfoDb) (*fastjson.Value, error) {
	reservationOrderJson := fmt.Sprintf("[{\"spaceId\":%s,\"timeId\":%s,\"venueSpaceGroupId\":null}]", taskInfo.SiteId, taskInfo.TimeId)
	//提交订单信息的data
	dataFormSubmit := map[string]string{
		"venueSiteId":          taskInfo.VenueSiteID,
		"reservationDate":      taskInfo.ReservationDate,
		"reservationOrderJson": reservationOrderJson,
		"phone":                taskInfo.UserPhone,
		"buddyIds":             taskInfo.BuddyId,
		"weekStartDate":        taskInfo.ReservationDate,
		"captchaVerification":  taskInfo.CaptchaVerification,
		"buddyNo":              taskInfo.BuddyNum,
		"isOfflineTicket":      "1",
		"token":                taskInfo.OrderToken,
	}
	submitResponse, fetchSubmitErr := fetchSubmit(client, dataFormSubmit)
	if fetchSubmitErr != nil {
		return nil, fmt.Errorf("提交订单信息失败：%v", fetchSubmitErr)
	}
	return submitResponse, nil
}

func payOrder(client *resty.Client, taskInfo TaskInfoDb) error {
	//提交订单信息的data
	dataFormPay := map[string]string{
		"venueTradeNo": taskInfo.TradeNo,
		"isApp":        "0",
	}
	_, err := fetchPay(client, dataFormPay)
	if err != nil {
		return err
	}
	return nil
}

func generateTaskId(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func addReserveTask(taskInfo TaskInfo) error {
	// 随机创建一个taskId
	taskId := generateTaskId(6) // 16 characters long
	taskInfo.TaskId = taskId
	// 添加任务到数据库
	addTaskErr := addTaskFromDB(taskInfo)
	if addTaskErr != nil {
		return fmt.Errorf("任务添加失败%v", addTaskErr)
	}
	// 任务添加成功
	log.Printf("任务添加成功，任务id为：%s", taskId)
	// 启动该预约任务
	creatTtyysReserveTask(taskId)
	return nil
}

func creatTtyysReserveTask(taskId string) error {
	// 根据taskId 查询数据库，获取对应taskId 的数据
	taskInfo, err := getTaskInfoByTaskId(taskId)
	if err != nil {
		return fmt.Errorf("查询数据库失败：%v", err)
	}
	// 如果是立即预约，则直接预约
	if taskInfo.InstantReservation {
		// 启动一个协程执行预约任务
		log.Println("立即预约为true，执行立即预约任务")
		go func(info TaskInfoDb) {
			if err := tyysReserveTask(&info, true); err != nil {
				log.Printf("预约任务执行失败：%v", err)
			}
		}(taskInfo)
	} else {
		// 启动一个协程执行延时预约任务
		log.Println("立即预约为false，执行延时预约任务")
		tyysReserveDelayTask(taskInfo)
	}
	return nil
}

func tyysReserveDelayTask(taskInfo TaskInfoDb) error {
	// 解析预约时间
	reservationTime, err := time.ParseInLocation("2006-01-02 15:04:05", taskInfo.ReservationDate+" 08:59:00", time.Local)
	if err != nil {
		return fmt.Errorf("解析时间出错: %v", err)
	}

	// 将时间往前推48小时
	executionTime := reservationTime.Add(-48 * time.Hour)

	// 计算调度时间，精确到秒
	cronTime := fmt.Sprintf("%d %d %d %d %d *", executionTime.Second(), executionTime.Minute(), executionTime.Hour(), executionTime.Day(), executionTime.Month())

	// 创建支持秒级调度的 cron 实例
	c := cron.New(cron.WithSeconds())

	// 添加任务到调度
	_, err = c.AddFunc(cronTime, func() {
		log.Printf("执行预约任务，原定时间: %v", executionTime)

		const maxRetries = 3
		var lastErr error

		for i := 0; i < maxRetries; i++ {
			taskAttempt := taskInfo
			err := tyysReserveTask(&taskAttempt, false)
			if err == nil {
				log.Printf("预约任务执行成功")
				return
			}

			lastErr = err
			log.Printf("预约任务执行失败 (尝试 %d/%d): %v", i+1, maxRetries, err)

			if i < maxRetries-1 {
				time.Sleep(2 * time.Second) // 重试之间等待2秒
			}
		}

		log.Printf("预约任务在 %d 次尝试后仍然失败: %v", maxRetries, lastErr)
	})
	if err != nil {
		return fmt.Errorf("调度任务出错: %v", err)
	}

	c.Start()

	log.Printf("任务已安排，将在 %v 执行", executionTime)
	return nil
}

func tyysGetOrderCode(username, password, tradeId string) (string, error) {
	client, _, _ := NewClient(username, password)
	orderCodeResponse, orderCodeErr := fetchOderCode(client, tradeId)
	if orderCodeErr != nil {
		return "", orderCodeErr
	}
	orderCode := string(orderCodeResponse.GetStringBytes("data", "orderNoCode"))
	return orderCode, nil
}

func tyysReserveTask(taskInfo *TaskInfoDb, isInstantReservation bool) error {
	// 在函数开始时使用 defer 确保在函数返回时更新状态
	defer func() {
		if err := updateTask(taskInfo.TaskID, "is_finished", true); err != nil {
			log.Printf("更新任务状态失败：%v", err)
		}
	}()

	// 每次执行前重置关键状态
	taskInfo.TradeNo = ""
	taskInfo.OrderId = ""
	taskInfo.ReservationStatus = false

	// 登录预约账号
	client, _, loginErr := NewClient(taskInfo.Username, taskInfo.Password)
	if loginErr != nil {
		return fmt.Errorf("登录失败：%v", loginErr)
	}
	log.Printf("登录成功，用户名：%s", taskInfo.Username)

	// 获取同伴id
	getBuddyErr := getBuddyId(client, taskInfo)
	if getBuddyErr != nil {
		return fmt.Errorf("获取同伴id失败：%v", getBuddyErr)
	}
	log.Printf("获取同伴id成功，同伴id：%s", taskInfo.BuddyUserID)
	if !isInstantReservation {
		// 判断当前时间是否大于9点
		currentTime := time.Now()
		nineAmOneSecond := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 9, 0, 1, 0, currentTime.Location())
		if currentTime.Before(nineAmOneSecond) {
			// 当前时间小于9:00:01，等待到9:00:01
			time.Sleep(nineAmOneSecond.Sub(currentTime))
		}
		log.Printf("开始提交预约信息，当前时间：%s", time.Now().Format("2006-01-02 15:04:05"))
	}
	//  获取预约场地信息
	getVenueInfoErr := getVenueInfo(client, taskInfo)
	if getVenueInfoErr != nil {
		return fmt.Errorf("获取场地信息失败：%v", getVenueInfoErr)
	}
	log.Printf("获取场地信息成功，场地信息：%s", taskInfo.SiteName)
	const captchaAttempts = 3
	var lastErr error
	var lastMessage string

	for attempt := 1; attempt <= captchaAttempts; attempt++ {
		if attempt > 1 {
			log.Printf("准备重新提交预约订单 [尝试: %d/%d]", attempt, captchaAttempts)
			time.Sleep(time.Second * time.Duration(attempt))
		}

		if err := confirmVenue(client, taskInfo); err != nil {
			lastErr = err
			log.Printf("提交预约信息失败 [尝试: %d/%d]：%v", attempt, captchaAttempts, err)
			continue
		}
		log.Printf("提交预约信息成功，场地信息：%s", taskInfo.SiteName)

		if err := getCaptchaToken(client, taskInfo); err != nil {
			lastErr = err
			log.Printf("获取验证码失败 [尝试: %d/%d]：%v", attempt, captchaAttempts, err)
			continue
		}

		time.Sleep(1500 * time.Millisecond)

		submitResponse, err := submitOrder(client, taskInfo)
		if err != nil {
			lastErr = err
			log.Printf("提交订单失败 [尝试: %d/%d]：%v", attempt, captchaAttempts, err)
			continue
		}

		if submitResponse.GetInt64("code") != 200 {
			message := submitResponse.GetStringBytes("message")
			lastMessage = string(message)
			log.Printf("提交订单失败 [尝试: %d/%d]：%s", attempt, captchaAttempts, message)
			if strings.Contains(lastMessage, "验证码") && attempt < captchaAttempts {
				continue
			}
			return fmt.Errorf("提交订单失败：%s", lastMessage)
		}

		taskInfo.TradeNo = string(submitResponse.GetStringBytes("data", "orderInfo", "tradeNo"))
		log.Printf("提交订单成功，trade_num：%s", taskInfo.TradeNo)

		oderDetailResponse, oderDetailErr := fetchOderdetail(client, taskInfo.TradeNo)
		if oderDetailErr != nil {
			return fmt.Errorf("获取订单详情失败：%v", oderDetailErr)
		}

		taskInfo.OrderId = strconv.FormatInt(oderDetailResponse.Get("data").GetInt64("orderId"), 10)
		log.Printf("获取订单详情成功，order_id：%s", taskInfo.OrderId)

		if payOrderErr := payOrder(client, *taskInfo); payOrderErr != nil {
			return fmt.Errorf("支付订单失败：%v", payOrderErr)
		}
		log.Printf("支付订单成功，订单号：%s", taskInfo.TradeNo)

		taskInfo.ReservationStatus = true

		if updateErr := updateTask(taskInfo.TaskID, "reservation_status", true); updateErr != nil {
			return fmt.Errorf("更新订单状态失败：%v", updateErr)
		}
		if updateErr := updateTask(taskInfo.TaskID, "trade_no", taskInfo.TradeNo); updateErr != nil {
			return fmt.Errorf("更新任务信息失败：%v", updateErr)
		}
		if updateErr := updateTask(taskInfo.TaskID, "order_id", taskInfo.OrderId); updateErr != nil {
			return fmt.Errorf("更新任务信息失败：%v", updateErr)
		}

		log.Printf("恭喜%s，您的场地预约成功啦！%s %s %s", taskInfo.Username, taskInfo.ReservationDate, taskInfo.SiteName, taskInfo.ReservationTime)
		return nil
	}

	if lastMessage != "" {
		return fmt.Errorf("提交订单失败：%s", lastMessage)
	}
	if lastErr != nil {
		return fmt.Errorf("提交订单失败：%v", lastErr)
	}
	return fmt.Errorf("提交订单失败：未知原因")
}
