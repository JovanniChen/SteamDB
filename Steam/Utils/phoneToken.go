package Utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
)

const chars = "23456789BCDFGHJKMNPQRTVWXY"

// bufferizeSecret 将秘密字符串转换为字节切片
func bufferizeSecret(secret string) []byte {
	if secret == "" {
		return []byte(secret)
	}
	if len(secret) == 40 {
		if _, err := hex.DecodeString(secret); err == nil {
			decoded, _ := hex.DecodeString(secret)
			return decoded
		}
	}
	decoded, _ := base64.StdEncoding.DecodeString(secret)
	return decoded
}

// GenerateAuthCode 生成认证码
func GenerateAuthCode(secret string, time int64) string {
	secretBytes := bufferizeSecret(secret)
	b := make([]byte, 8)
	// 这里相当于 Node.js  里的 b.writeUInt32BE(0,  0);
	for i := 0; i < 4; i++ {
		b[i] = 0
	}
	// 这里相当于 Node.js  里的 b.writeUInt32BE(Math.floor(time  / 30), 4);
	timeDiv30 := uint32(math.Floor(float64(time) / 30))
	for i := 3; i >= 0; i-- {
		b[4+i] = byte(timeDiv30 & 0xff)
		timeDiv30 >>= 8
	}

	hmacHash := hmac.New(sha1.New, secretBytes)
	hmacHash.Write(b)
	hmacData := hmacHash.Sum(nil)

	start := hmacData[19] & 0x0f
	hmacData = hmacData[start : start+4]

	var fullcode uint32
	for i := 0; i < 4; i++ {
		fullcode = (fullcode << 8) | uint32(hmacData[i])
	}
	fullcode &= 0x7fffffff

	code := ""
	for i := 0; i < 5; i++ {
		code += string(chars[fullcode%uint32(len(chars))])
		fullcode /= uint32(len(chars))
	}
	return code
}

// GenerateConfirmationQueryParams 生成确认查询参数
func GenerateConfirmationQueryParams(deviceID, identitySecret, steamid string, time int64, tag string) (map[string]string, error) {
	if deviceID == "" {
		return nil, errors.New("Device ID is not present")
	}
	return map[string]string{
		"p":   deviceID,
		"a":   steamid,
		"k":   GenerateConfirmationHashForTime(identitySecret, time, tag),
		"t":   strconv.FormatInt(time, 10),
		"m":   "react",
		"tag": tag,
	}, nil
}

// GenerateConfirmationHashForTime 计算手机操作需要的密文
func GenerateConfirmationHashForTime(identitySecret string, time int64, tag string) string {
	decode, _ := base64.StdEncoding.DecodeString(identitySecret)
	var n2 int
	if tag != "" {
		if len(tag) > 32 {
			n2 = 8 + 32
		} else {
			n2 = 8 + len(tag)
		}
	} else {
		n2 = 8
	}
	array := make([]byte, n2)
	for i := 7; i >= 0; i-- {
		array[i] = byte(time & 0xff)
		time >>= 8
	}
	if tag != "" {
		copy(array[8:], tag)
	}

	hmacHash := hmac.New(sha1.New, decode)
	hmacHash.Write(array)
	hashedData := hmacHash.Sum(nil)
	encodedData := base64.StdEncoding.EncodeToString(hashedData)
	return fmt.Sprintf("%s", encodedData)
}
