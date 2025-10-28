package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/go-resty/resty/v2"
)

func NewClient(username, password string) (*resty.Client, *UserInfo, error) {
	// 实例化 UserInfo 结构体
	userInfo := UserInfo{
		Username: username,
		Password: password,
	}
	client := resty.New()

	client.SetRedirectPolicy(resty.NoRedirectPolicy()) // 禁用自动重定向

	// 预认证 获取密码加密参数
	execution, pubKey := preAuth(client)
	// 加密密码
	password_encrypt := encryptPassword(pubKey, userInfo.Password)
	// 获取tiketUrl
	tiketUrl := getTiketUrl(client, userInfo.Username, password_encrypt, execution)
	if tiketUrl == "" {
		return nil, nil, fmt.Errorf("login failed: unexpected login response (check username/password)")
	}
	log.Println(tiketUrl)
	// 获取jsessionIdUrl
	jsessionIdUrl := getJsessionIdUrl(client, tiketUrl)
	if jsessionIdUrl == "" {
		return nil, nil, fmt.Errorf("login failed: session redirect missing")
	}
	log.Println(jsessionIdUrl)
	// 获取oauthTokenUrl
	oauthTokenUrl := getOauthToken(client, jsessionIdUrl)
	if oauthTokenUrl == "" {
		return nil, nil, fmt.Errorf("login failed: oauth token redirect missing")
	}
	log.Println(oauthTokenUrl)
	// 获取oauthToken
	oauthToken, err := extractOauthToken(oauthTokenUrl)
	if err != nil {
		return nil, nil, err
	}
	// 获取登录信息
	loginInfoResponse, err := fetchLoginInfo(client, oauthToken)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	client.SetHeader("cgauthorization", string(loginInfoResponse.GetStringBytes("data", "token", "access_token")))

	// 更新 UserInfo 结构体
	userInfo.UserId = loginInfoResponse.GetInt64("data", "userId")
	userInfo.Name = string(loginInfoResponse.GetStringBytes("data", "name"))
	userInfo.Role = loginInfoResponse.GetInt64("data", "role")
	// 获取同伴码
	buddiesResponse, err := fetchBuddies(client)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	userInfo.BuddyNum = string(buddiesResponse.GetStringBytes("data"))
	// 获取个人信息
	var personalInfoUrl string
	switch userInfo.Role {
	case 3:
		personalInfoUrl = studentUrl + strconv.FormatInt(userInfo.UserId, 10)
	case 4:
		personalInfoUrl = teacherUrl + strconv.FormatInt(userInfo.UserId, 10)
	}
	personalInfoResponse, err := fetchPersonalInfo(client, personalInfoUrl)
	if err != nil {
		return nil, nil, err
	}
	userInfo.Phone = string(personalInfoResponse.GetStringBytes("data", "phone"))

	return client, &userInfo, nil
}

// extractOauthToken 提取URL中的oauth_token值
func extractOauthToken(urlStr string) (string, error) {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	queryParams := parsedUrl.Query()
	oauthToken := queryParams.Get("oauth_token")
	if oauthToken == "" {
		return "", fmt.Errorf("oauth_token not found in URL")
	}
	return oauthToken, nil
}

// 预认证 获取execution和pubKey
func preAuth(client *resty.Client) (string, *PubKey) {
	// 获取execution
	execution := getExecution(client)

	// 获取RSA加密公钥
	pubKey, _ := getPubKey(client)

	return *execution, pubKey
}

func getExecution(client *resty.Client) *string {
	headers := map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   userAgent,
	}
	resp, err := client.R().SetHeaders(headers).Get(loginUrl)
	if err != nil {
		log.Println(err)
		return nil
	}
	reg, _ := regexp.Compile("name=\"execution\" value=\"(.*?)\"")
	execution := reg.FindStringSubmatch(string(resp.Body()))[1]
	return &execution

}

func getPubKey(client *resty.Client) (*PubKey, error) {
	resp, err := client.R().SetHeader("Content-Type", "application/json").Get(pubKeyUrl)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	getPubKey := PubKey{}
	err = json.Unmarshal(resp.Body(), &getPubKey)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &getPubKey, nil
}

// 对密码进行RSA加密的函数, golang不支持NoPadding模式，自己实现
func encryptPassword(pubKey *PubKey, password string) string {
	Modulus, _ := parse2bigInt(pubKey.Modulus)
	Exponent, _ := parse2bigInt(pubKey.Exponent)
	content := big.NewInt(0).SetBytes([]byte(password))
	var z big.Int
	z.Exp(content, Exponent, Modulus)
	return fmt.Sprintf("%x", &z)
}

// 将字符串转化为big.Int的函数
func parse2bigInt(s string) (i *big.Int, err error) {
	i = new(big.Int)
	i.SetString(s, 16)
	return
}

func getTiketUrl(client *resty.Client, username string, password string, execution string) string {
	formData := map[string]string{
		"username":  username,
		"password":  password,
		"authcode":  "",
		"execution": execution,
		"_eventId":  "submit",
	}
	headers := map[string]string{
		"User-Agent":   userAgent,
		"Content-Type": "application/x-www-form-urlencoded",
	}

	resp, err := client.R().
		SetHeaders(headers).
		SetFormData(formData).
		Post(loginUrl)

	// 首先检查是否是重定向
	if resp != nil && resp.StatusCode() == http.StatusFound {
		location := resp.Header().Get("Location")
		return location
	}

	// 然后处理其他可能的错误
	if err != nil {
		log.Println("Error:", err)
		return ""
	}

	// 如果不是重定向，也不是错误，那么可能是意外的响应
	log.Printf("Unexpected status code: %d\n", resp.StatusCode())
	return ""
}

func getJsessionIdUrl(client *resty.Client, tiketUrl string) string {
	headers := map[string]string{
		"User-Agent": userAgent,
	}

	resp, err := client.R().
		SetHeaders(headers).
		Get(tiketUrl)

	/// 首先检查是否是重定向
	if resp != nil && resp.StatusCode() == http.StatusFound {
		location := resp.Header().Get("Location")
		return location
	}

	// 然后处理其他可能的错误
	if err != nil {
		log.Println("Error:", err)
		return ""
	}

	// 如果不是重定向，也不是错误，那么可能是意外的响应
	log.Printf("Unexpected status code: %d\n", resp.StatusCode())
	return ""
}

func getOauthToken(client *resty.Client, jsessionIdUrl string) string {
	headers := map[string]string{
		"User-Agent": userAgent,
	}

	resp, err := client.R().
		SetHeaders(headers).
		Get(jsessionIdUrl)
	// 首先检查是否是重定向
	if resp != nil && resp.StatusCode() == http.StatusFound {
		location := resp.Header().Get("Location")
		return location
	}

	// 然后处理其他可能的错误
	if err != nil {
		log.Println("Error:", err)
		return ""
	}

	// 如果不是重定向，也不是错误，那么可能是意外的响应
	log.Printf("Unexpected status code: %d\n", resp.StatusCode())
	return ""
}
