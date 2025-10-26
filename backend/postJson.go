package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// Point 结构体用于表示坐标点
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// 将给定的数据使用 AES 算法加密，并返回 Base64 编码后的字符串
func encryptWithAES(data string, secretKey string) (string, error) {
	keyBytes := []byte(secretKey)
	dataBytes := []byte(data)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	paddedData := pkcs7Pad(dataBytes, block.BlockSize())
	ciphertext := make([]byte, len(paddedData))

	ecbEncrypt(ciphertext, paddedData, block)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

func ecbEncrypt(dst, src []byte, block cipher.Block) {
	for len(src) > 0 {
		block.Encrypt(dst, src[:block.BlockSize()])
		src = src[block.BlockSize():]
		dst = dst[block.BlockSize():]
	}
}

// GetCaptchaVerification 获取验证数据
func GetCaptchaVerification(secretKey, backToken string, positionJsonStr string) (string, error) {
	data := fmt.Sprintf("%s---%s", backToken, positionJsonStr)

	if secretKey != "" {
		encryptedData, err := encryptWithAES(data, secretKey)
		if err != nil {
			return "", err
		}
		return encryptedData, nil
	}

	return data, nil
}

// getPointJson 数据
func getPointJson(positionJsonStr, secretKey string) (string, error) {
	// 如果密钥存在，则对 JSON 数据进行加密
	if secretKey != "" {
		encryptedData, err := encryptWithAES(positionJsonStr, secretKey)
		if err != nil {
			return "", err
		}
		return encryptedData, nil
	}

	// 否则，直接返回 JSON 数据
	return "", nil
}
