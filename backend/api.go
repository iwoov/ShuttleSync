package main

import (
	"crypto/md5"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/valyala/fastjson"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// fetchLoginInfo 获取token的接口 post
func fetchLoginInfo(client *resty.Client, oauthToken string) (*fastjson.Value, error) {
	return retryRequest(3, time.Millisecond*500, func() (*fastjson.Value, error) {
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		data := map[string]string{}
		tokenBaseUrl := "/api/login"
		sign := getSign(tokenBaseUrl, data, timeStamp)
		refer := "http://www.tyys.zju.edu.cn/venue/login?oauth_token=" + oauthToken
		headers := map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"Host":         "www.tyys.zju.edu.cn",
			"Dnt":          "1",
			"Origin":       "http://www.tyys.zju.edu.cn",
			"Referer":      refer,
			"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Edg/128.0.0.0",
			"App-key":      appKey,
			"Oauth-token":  oauthToken,
			"Sign":         sign,
			"Timestamp":    timeStamp,
		}
		reps, err := client.R().
			SetHeaders(headers).
			SetFormData(data).
			Post(tokenUrl)

		if err != nil {
			return nil, fmt.Errorf("请求错误: %v", err)
		}

		if reps.StatusCode() != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d", reps.StatusCode())
		}
		var p fastjson.Parser
		return p.ParseBytes(reps.Body())
	})
}

// fetchBuddies 获取同伴码的接口 post
func fetchBuddies(client *resty.Client) (*fastjson.Value, error) {
	return retryRequest(3, time.Millisecond*500, func() (*fastjson.Value, error) {
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		data := map[string]string{}
		parsedURL, err := url.Parse(buddiesUrl)
		if err != nil {
			return nil, err
		}
		sign := getSign(strings.Replace(parsedURL.Path, "/venue-server", "", 1), data, timeStamp)
		headers := map[string]string{
			"App-key":    appKey,
			"User-Agent": userAgent,
			"Sign":       sign,
			"Timestamp":  timeStamp,
		}
		reps, err := client.R().
			SetHeaders(headers).
			SetFormData(data).
			Post(buddiesUrl)

		if err != nil {
			return nil, fmt.Errorf("请求错误: %v", err)
		}
		if reps.StatusCode() != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d", reps.StatusCode())
		}
		var p fastjson.Parser
		return p.ParseBytes(reps.Body())
	})
}

// fetchPersonalInfo 获取个人信息
func fetchPersonalInfo(client *resty.Client, personalUrl string) (*fastjson.Value, error) {
	return retryRequest(3, time.Millisecond*500, func() (*fastjson.Value, error) {
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		data := map[string]string{"nocache": timeStamp}
		sign := getSign(personalUrl, data, timeStamp)
		personalInfoUrl := baseUrl + personalUrl
		headers := map[string]string{
			"App-key": appKey, "User-Agent": userAgent,
			"Sign": sign, "Timestamp": timeStamp,
		}
		reps, err := client.R().
			SetHeaders(headers).
			SetQueryParam("nocache", timeStamp).
			Get(personalInfoUrl)
		if err != nil {
			return nil, fmt.Errorf("请求错误: %v", err)
		}
		if reps.StatusCode() != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d", reps.StatusCode())
		}
		var p fastjson.Parser
		return p.ParseBytes(reps.Body())
	})
}

// fetchBuddiesList 获取用户的所有已添加的同伴信息的接口 get
func fetchBuddiesList(client *resty.Client) (*fastjson.Value, error) {
	return retryRequest(3, time.Millisecond*500, func() (*fastjson.Value, error) {
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		data := map[string]string{
			"page": "0", "size": "20", "nocache": timeStamp,
		}
		parsedURL, err := url.Parse(buddiesListUrl)
		if err != nil {
			return nil, err
		}
		sign := getSign(strings.Replace(parsedURL.Path, "/venue-server", "", 1), data, timeStamp)
		headers := map[string]string{
			"App-Key": appKey, "User-Agent": userAgent,
			"Sign": sign, "Timestamp": timeStamp,
		}
		reps, err := client.R().
			SetHeaders(headers).
			SetQueryParams(data).
			Get(buddiesListUrl)
		if err != nil {
			return nil, fmt.Errorf("请求错误: %v", err)
		}
		if reps.StatusCode() != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d", reps.StatusCode())
		}
		var p fastjson.Parser
		return p.ParseBytes(reps.Body())
	})
}

// fetchVenueInfo 获取紫金港风雨操场 羽毛球场地的可预约时间接口 get 获取重要参数token
func fetchVenueInfo(client *resty.Client, reservationDate string, venueSiteId string) (*fastjson.Value, error) {
	return retryRequest(3, time.Millisecond*500, func() (*fastjson.Value, error) {
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		data := map[string]string{
			"venueSiteId": venueSiteId,
			"searchDate":  reservationDate,
			"nocache":     timeStamp,
		}
		parsedURL, err := url.Parse(venueInfoUrl)
		if err != nil {
			return nil, err
		}
		sign := getSign(strings.Replace(parsedURL.Path, "/venue-server", "", 1), data, timeStamp)
		headers := map[string]string{
			"App-Key": appKey, "User-Agent": userAgent,
			"Sign": sign, "Timestamp": timeStamp,
		}
		reps, err := client.R().
			SetHeaders(headers).
			SetQueryParams(data).
			Get(venueInfoUrl)
		if err != nil {
			return nil, fmt.Errorf("请求错误: %v", err)
		}
		if reps.StatusCode() != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d", reps.StatusCode())
		}
		var p fastjson.Parser
		return p.ParseBytes(reps.Body())
	})
}

// fetchVenue 选择预约场地的接口 post
func fetchVenue(client *resty.Client, data map[string]string) (*fastjson.Value, error) {
	return retryRequest(3, time.Millisecond*500, func() (*fastjson.Value, error) {
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		parsedURL, err := url.Parse(venueOrderUrl)
		if err != nil {
			return nil, err
		}
		sign := getSign(strings.Replace(parsedURL.Path, "/venue-server", "", 1), data, timeStamp)
		headers := map[string]string{
			"App-Key": appKey, "User-Agent": userAgent,
			"Sign": sign, "Timestamp": timeStamp,
			"Content-Type": "application/x-www-form-urlencoded",
		}
		reps, err := client.R().
			SetHeaders(headers).
			SetFormData(data).
			Post(venueOrderUrl)
		if err != nil {
			return nil, fmt.Errorf("请求错误: %v", err)
		}
		if reps.StatusCode() != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d", reps.StatusCode())
		}
		var p fastjson.Parser
		return p.ParseBytes(reps.Body())
	})
}

// fetchSubmit 提交订单信息的接口 post
func fetchSubmit(client *resty.Client, data map[string]string) (*fastjson.Value, error) {
	return retryRequest(3, time.Millisecond*500, func() (*fastjson.Value, error) {
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		submitBaseUrl := "/api/reservation/order/submit"
		sign := getSign(submitBaseUrl, data, timeStamp)
		headers := map[string]string{
			"App-Key": appKey, "User-Agent": userAgent,
			"Sign": sign, "Timestamp": timeStamp,
			"Content-Type": "application/x-www-form-urlencoded",
		}
		reps, err := client.R().
			SetHeaders(headers).
			SetFormData(data).
			Post(submitUrl)
		if err != nil {
			return nil, fmt.Errorf("请求错误: %v", err)
		}
		if reps.StatusCode() != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d", reps.StatusCode())
		}
		var p fastjson.Parser
		return p.ParseBytes(reps.Body())
	})
}

// fetchPay 支付订单的接口 post
func fetchPay(client *resty.Client, data map[string]string) (*fastjson.Value, error) {
	return retryRequest(3, time.Millisecond*500, func() (*fastjson.Value, error) {
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		payBaseUrl := "/api/venue/finances/order/pay"
		sign := getSign(payBaseUrl, data, timeStamp)
		headers := map[string]string{
			"App-Key": appKey, "User-Agent": userAgent,
			"Sign": sign, "Timestamp": timeStamp,
			"Content-Type": "application/x-www-form-urlencoded",
		}
		reps, err := client.R().
			SetHeaders(headers).
			SetFormData(data).
			Post(PayUrl)
		if err != nil {
			return nil, fmt.Errorf("请求错误: %v", err)
		}
		if reps.StatusCode() != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d", reps.StatusCode())
		}
		var p fastjson.Parser
		return p.ParseBytes(reps.Body())
	})
}

// fetchOrderList 获取订单的接口 get
func fetchOrderList(client *resty.Client) (*fastjson.Value, error) {
	return retryRequest(3, time.Millisecond*500, func() (*fastjson.Value, error) {
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		data := map[string]string{
			"page": "0", "size": "20", "nocache": timeStamp,
		}
		parsedURL, err := url.Parse(oderUrl)
		if err != nil {
			return nil, err
		}
		sign := getSign(strings.Replace(parsedURL.Path, "/venue-server", "", 1), data, timeStamp)
		headers := map[string]string{
			"App-Key": appKey, "User-Agent": userAgent,
			"sign": sign, "timestamp": timeStamp,
		}
		reps, err := client.R().
			SetHeaders(headers).
			SetQueryParams(data).
			Get(oderUrl)
		if err != nil {
			return nil, fmt.Errorf("请求错误: %v", err)
		}
		if reps.StatusCode() != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d", reps.StatusCode())
		}
		var p fastjson.Parser
		return p.ParseBytes(reps.Body())
	})
}

// cancel 支付订单的接口 post
func cancel(client *resty.Client, data map[string]string) (*fastjson.Value, error) {
	return retryRequest(3, time.Millisecond*500, func() (*fastjson.Value, error) {
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		cancelBaseUrl := "/api/venue/finances/order/cancel"
		sign := getSign(cancelBaseUrl, data, timeStamp)
		headers := map[string]string{
			"App-Key": appKey, "User-Agent": userAgent,
			"Sign": sign, "Timestamp": timeStamp,
			"Content-Type": "application/x-www-form-urlencoded",
		}
		reps, err := client.R().
			SetHeaders(headers).
			SetFormData(data).
			Post(cancelUrl)
		if err != nil {
			return nil, fmt.Errorf("请求错误: %v", err)
		}
		if reps.StatusCode() != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d", reps.StatusCode())
		}
		var p fastjson.Parser
		return p.ParseBytes(reps.Body())
	})
}

// fetchOderCode 获取预约码的接口
func fetchOderCode(client *resty.Client, tradeId string) (*fastjson.Value, error) {
	return retryRequest(3, time.Millisecond*500, func() (*fastjson.Value, error) {
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		data := map[string]string{"nocache": timeStamp}
		parsedURL, err := url.Parse(oderCodeUrl + tradeId)
		if err != nil {
			return nil, err
		}
		sign := getSign(strings.Replace(parsedURL.Path, "/venue-server", "", 1), data, timeStamp)
		headers := map[string]string{
			"App-Key": appKey, "User-Agent": userAgent,
			"Sign": sign, "Timestamp": timeStamp,
		}
		reps, err := client.R().
			SetHeaders(headers).
			SetQueryParams(data).
			Get(oderCodeUrl + tradeId)
		if err != nil {
			return nil, fmt.Errorf("请求错误: %v", err)
		}
		if reps.StatusCode() != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d", reps.StatusCode())
		}
		var p fastjson.Parser
		return p.ParseBytes(reps.Body())
	})
}

// fetchCaptcha 获取金港风雨操场 羽毛球场地的可预约时间接口 get 获取重要参数token
func fetchCaptcha(client *resty.Client) (*fastjson.Value, error) {
	return retryRequest(3, time.Millisecond*500, func() (*fastjson.Value, error) {
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		data := map[string]string{
			"captchaType": "clickWord",
			"clientUid":   "point-591fcbdd-da39-40d9-b188-91af1b5d7171",
			"ts":          timeStamp,
			"nocache":     timeStamp,
		}
		parsedURL, err := url.Parse(captchaUrl)
		if err != nil {
			return nil, err
		}
		sign := getSign(strings.Replace(parsedURL.Path, "/venue-server", "", 1), data, timeStamp)
		headers := map[string]string{
			"App-Key": appKey, "User-Agent": userAgent,
			"Sign": sign, "Timestamp": timeStamp,
		}
		reps, err := client.R().
			SetHeaders(headers).
			SetQueryParams(data).
			Get(captchaUrl)
		if err != nil {
			return nil, fmt.Errorf("请求错误: %v", err)
		}
		if reps.StatusCode() != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d", reps.StatusCode())
		}
		var p fastjson.Parser
		return p.ParseBytes(reps.Body())
	})
}

// checkCaptcha 检查captcha接口 post
func checkCaptcha(client *resty.Client, data map[string]string) (*fastjson.Value, error) {
	return retryRequest(3, time.Millisecond*500, func() (*fastjson.Value, error) {
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		checkCaptchaBaseUrl := "/api/captcha/check"
		sign := getSign(checkCaptchaBaseUrl, data, timeStamp)
		headers := map[string]string{
			"App-Key": appKey, "User-Agent": userAgent,
			"Sign": sign, "Timestamp": timeStamp,
			"Content-Type": "application/x-www-form-urlencoded",
		}
		reps, err := client.R().
			SetHeaders(headers).
			SetFormData(data).
			Post(checkCaptchaUrl)
		if err != nil {
			return nil, fmt.Errorf("请求错误: %v", err)
		}
		if reps.StatusCode() != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d", reps.StatusCode())
		}
		var p fastjson.Parser
		return p.ParseBytes(reps.Body())
	})
}

// fetchOderdetail 获取订单详情的接口 get
func fetchOderdetail(client *resty.Client, tradeNo string) (*fastjson.Value, error) {
	return retryRequest(3, time.Millisecond*500, func() (*fastjson.Value, error) {
		timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		data := map[string]string{
			"venueTradeNo": tradeNo,
			"nocache":      timeStamp,
		}
		parsedURL, err := url.Parse(oderDetailUrl)
		if err != nil {
			return nil, err
		}
		sign := getSign(strings.Replace(parsedURL.Path, "/venue-server", "", 1), data, timeStamp)
		headers := map[string]string{
			"App-Key": appKey, "User-Agent": userAgent,
			"Sign": sign, "Timestamp": timeStamp,
		}
		reps, err := client.R().
			SetHeaders(headers).
			SetQueryParams(data).
			Get(oderDetailUrl)
		if err != nil {
			return nil, fmt.Errorf("请求错误: %v", err)
		}
		if reps.StatusCode() != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d", reps.StatusCode())
		}
		var p fastjson.Parser
		return p.ParseBytes(reps.Body())
	})
}

// 获取sign
func getSign(pathUrl string, data map[string]string, timeStamp string) string {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	dataInfo := ""
	for _, key := range keys {
		value := data[key]
		if value != "" && value != "null" && value != "undefined" {
			dataInfo += key + value
		}
	}
	concatenatedString := para + pathUrl + dataInfo + timeStamp + " " + para
	sign := fmt.Sprintf("%x", md5.Sum([]byte(concatenatedString)))
	return sign
}

// 通用的重试函数
func retryRequest(maxRetries int, retryDelay time.Duration, f func() (*fastjson.Value, error)) (*fastjson.Value, error) {
	for i := 0; i < maxRetries; i++ {
		v, err := f()
		if err == nil {
			return v, nil
		}
		if i == maxRetries-1 {
			return nil, err
		}
		time.Sleep(retryDelay)
	}
	return nil, fmt.Errorf("超过最大重试次数")
}
