package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

const CustomUrl = "http://api.jfbym.com/api/YmServer/customApi"

func commonVerify(image string, ApiKey string) string {
	//# 数英汉字类型
	//# 通用数英1-4位 10110
	//# 通用数英5-8位 10111
	//# 通用数英9~11位 10112
	//# 通用数英12位及以上 10113
	//# 通用数英1~6位plus 10103
	//# 定制-数英5位~qcs 9001
	//# 定制-纯数字4位 193
	//# 中文类型
	//# 通用中文字符1~2位 10114
	//# 通用中文字符 3~5位 10115
	//# 通用中文字符6~8位 10116
	//# 通用中文字符9位及以上 10117
	//# 定制-XX西游苦行中文字符 10107
	//# 计算类型
	//# 通用数字计算题 50100
	//# 通用中文计算题 50101
	//# 定制-计算题 cni 452

	config := map[string]any{}
	config["image"] = image
	config["type"] = "10110"
	config["token"] = ApiKey
	configData, _ := json.Marshal(config)
	body := bytes.NewBuffer([]byte(configData))
	resp, err := http.Post(CustomUrl, "application/json;charset=utf-8", body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data), err)
	return string(data)
}

func slideVerify(slideImage string, backgroundImage string, ApiKey string) string {
	//# 滑块类型
	//# 通用双图滑块  20111
	// # slide_image 需要识别图片的小图片的base64字符串
	// # background_image 需要识别图片的背景图片的base64字符串(背景图需还原)
	config := map[string]any{}
	config["slide_image"] = slideImage
	config["background_image"] = backgroundImage
	config["type"] = "20111"
	config["token"] = ApiKey
	configData, _ := json.Marshal(config)
	body := bytes.NewBuffer([]byte(configData))
	resp, err := http.Post(CustomUrl, "application/json;charset=utf-8", body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data), err)
	return string(data)
}

func sinSlideVerify(image string, ApiKey string) string {
	//# 滑块类型
	//# 通用单图滑块(截图)  20110
	config := map[string]any{}
	config["image"] = image
	config["type"] = "20110"
	config["token"] = ApiKey
	configData, _ := json.Marshal(config)
	body := bytes.NewBuffer([]byte(configData))
	resp, err := http.Post(CustomUrl, "application/json;charset=utf-8", body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data), err)
	return string(data)
}

func trafficSlideVerify(seed string, data string, href string, ApiKey string) string {
	//# 滑块类型
	//# 定制-滑块协议slide_traffic  900010
	config := map[string]any{}
	config["seed"] = seed
	config["data"] = data
	config["href"] = href
	config["type"] = "900010"
	config["token"] = ApiKey
	configData, _ := json.Marshal(config)
	body := bytes.NewBuffer([]byte(configData))
	resp, err := http.Post(CustomUrl, "application/json;charset=utf-8", body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	_data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(_data), err)
	return string(data)
}

func clickVerify(image string, extra any, ApiKey string) []byte {
	//# 通用任意点选1~4个坐标 30009
	//# 通用文字点选1(extra,点选文字逗号隔开,原图) 30100
	//# 定制-文字点选2(extra="click",原图) 30103
	//# 定制-单图文字点选 30102
	//# 定制-图标点选1(原图) 30104
	//# 定制-图标点选2(原图,extra="icon") 30105
	//# 定制-语序点选1(原图,extra="phrase") 30106
	//# 定制-语序点选2(原图) 30107
	//# 定制-空间推理点选1(原图,extra="请点击xxx") 30109
	//# 定制-空间推理点选1(原图,extra="请_点击_小尺寸绿色物体。") 30110
	//# 定制-tx空间点选(extra="请点击侧对着你的字母") 50009
	//# 定制-tt_空间点选 30101
	//# 定制-推理拼图1(原图,extra="交换2个图块") 30108
	//# 定制-xy4九宫格点选(原图,label_image,image) 30008
	// # 定制-文字点选3(extra="je4_click") 30112
	//# 定制-图标点选3(extra="je4_icon") 30113
	//# 定制-语序点选3(extra="je4_phrase") 30114
	//# 如有其他未知类型,请联系我们
	config := map[string]any{}
	config["image"] = image
	if extra != nil {
		config["extra"] = extra
	}
	// config["type"] = os.Getenv("CAPTCHA_TYPE")
	config["type"] = "300010"
	config["token"] = ApiKey
	configData, _ := json.Marshal(config)
	body := bytes.NewBuffer([]byte(configData))
	resp, _ := http.Post(CustomUrl, "application/json;charset=utf-8", body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	data, _ := io.ReadAll(resp.Body)
	return data
}

func rotate(image string, ApiKey string) string {
	// # 定制-X度单图旋转  90007
	config := map[string]any{}

	//# 定制-Tt双图旋转,2张图,内圈图,外圈图  90004
	//config["out_ring_image"] = image
	//config["inner_circle_image"] = image
	//config["type"] = "90004"

	config["image"] = image
	config["type"] = "90007"
	config["token"] = ApiKey
	configData, _ := json.Marshal(config)
	body := bytes.NewBuffer([]byte(configData))
	resp, err := http.Post(CustomUrl, "application/json;charset=utf-8", body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data), err)
	return string(data)
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}

type Data struct {
	Code      int    `json:"code"`
	CaptchaId string `json:"captchaId"`
	RecordId  string `json:"recordId"`
	Data      string `json:"data"`
}

func googleVerify(googleKey string, pageUrl string, ApiKey string) string {
	config := map[string]any{}
	//第一步，创建验证码任务
	//:param
	//:return taskId : string 创建成功的任务ID
	url1 := "http://api.jfbym.com/api/YmServer/funnelApi"
	config["token"] = ApiKey
	config["type"] = "40010" // v2
	//config["type"] = "40011" // v3
	config["googlekey"] = googleKey
	config["pageurl"] = pageUrl
	config["enterprise"] = 0 //是否为企业版
	config["invisible"] = 0  //是否为可见类型
	config["data-s"] = ""    //## V2+企业如果能找到，找不到传空字符串
	//config["action"] = ""    //## #V3必传
	//config["min_score"] = "0.8" //#V3才支持的可选参数
	configData, _ := json.Marshal(config)
	body := bytes.NewBuffer([]byte(configData))
	resp, err := http.Post(url1, "application/json;charset=utf-8", body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data), err)
	//data := `{"msg": "识别成功", "code": 10000, "data": {"code": 0, "captchaId": "51436618130", "recordId": "74892"}}`
	fmt.Println(data)
	var res Result
	_ = json.Unmarshal([]byte(data), &res)
	captchaId := res.Data.CaptchaId
	recordId := res.Data.RecordId

	config1 := map[string]any{}
	timeOut := 120
	times := 0
	for {
		if times > timeOut {
			fmt.Println("超时")
			break
		}
		url2 := "http://api.jfbym.com/api/YmServer/funnelApiResult"
		config1["token"] = ApiKey
		config1["captchaId"] = captchaId
		config1["recordId"] = recordId
		configData, _ := json.Marshal(config1)
		body := bytes.NewBuffer([]byte(configData))
		resp, err := http.Post(url2, "application/json;charset=utf-8", body)
		if err != nil {
			fmt.Printf("polling captcha result failed: %v\n", err)
			time.Sleep(5 * time.Second)
			times += 5
			continue
		}
		data, readErr := io.ReadAll(resp.Body)
		closeErr := resp.Body.Close()
		if closeErr != nil {
			fmt.Printf("closing captcha result response failed: %v\n", closeErr)
		}
		if readErr != nil {
			fmt.Printf("reading captcha result failed: %v\n", readErr)
			time.Sleep(5 * time.Second)
			times += 5
			continue
		}
		//data := `{"msg": "请求成功", "code": 10001, "data": {"data": "03AGdBq2611GTOgA2v9HUpMMEUE70p6dwOtYyHJQK4xhdKF0Y8ouSGsFZt647SpJvZ22qinYrm6"}}`
		fmt.Println(data)
		var res Result
		_ = json.Unmarshal([]byte(data), &res)
		if res.Code != 10001 {
			time.Sleep(5 * time.Second)
			times += 5
			continue
		}
		return res.Data.Data

	}
	return string(data)
}

func hcaptchaVerify(siteKey string, siteUrl string, ApiKey string) string {
	//# Hcaptcha
	//# 请保证购买相应服务后请求对应 verify_type
	//# verify_type="50001"
	config := map[string]any{}
	config["site_key"] = siteKey
	config["site_url"] = siteUrl
	config["type"] = "50001"
	config["token"] = ApiKey
	configData, _ := json.Marshal(config)
	body := bytes.NewBuffer([]byte(configData))
	resp, err := http.Post(CustomUrl, "application/json;charset=utf-8", body)
	if err != nil {
		return fmt.Sprintf("hcaptcha verify request failed: %v", err)
	}
	defer resp.Body.Close()
	data, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return fmt.Sprintf("reading hcaptcha verify response failed: %v", readErr)
	}
	fmt.Println(string(data), err)
	return string(data)
}

func funCaptchaVerify(siteKey string, siteUrl string, ApiKey string) string {
	//# Hcaptcha
	//# 请保证购买相应服务后请求对应 verify_type
	//# verify_type="40007"
	config := map[string]any{}
	config["publickey"] = siteKey
	config["pageurl"] = siteUrl
	config["type"] = "40007"
	config["token"] = ApiKey
	configData, _ := json.Marshal(config)
	body := bytes.NewBuffer([]byte(configData))
	resp, err := http.Post(CustomUrl, "application/json;charset=utf-8", body)
	if err != nil {
		return fmt.Sprintf("funcaptcha verify request failed: %v", err)
	}
	defer resp.Body.Close()
	data, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return fmt.Sprintf("reading funcaptcha verify response failed: %v", readErr)
	}
	fmt.Println(string(data), err)
	return string(data)
}

// 从 clickVerify 函数返回的数据中提取 Point 切片
func extractPointsFromData(data []byte) ([]Point, error) {
	// 将字节数组转换为字符串
	jsonString := string(data)

	// 定义正则表达式，用于提取所有点坐标字符串
	re := regexp.MustCompile(`\d+,\d+`)

	// 从字符串中提取所有点坐标字符串
	match := re.FindAllString(jsonString, -1)
	if match == nil {
		return nil, fmt.Errorf("Failed to extract point data from JSON string: %s", jsonString)
	}

	// 将点坐标字符串解析为 Point 切片
	points := make([]Point, 0)
	for _, coord := range match {
		parts := strings.Split(coord, ",")
		if len(parts) != 2 {
			continue
		}
		x, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		y, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		points = append(points, Point{X: x, Y: y})
	}

	return points, nil
}

func Captcha(b64Str string, extra string, ApiKey string) (string, error) {
	data := clickVerify(b64Str, extra, ApiKey)

	// 从 clickVerify 函数返回的数据中提取点坐标数据
	points, err := extractPointsFromData(data)
	if err != nil {
		return "", err
	}

	// 将 Point 切片转换为 JSON 字符串
	jsonData, err := json.Marshal(points)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
