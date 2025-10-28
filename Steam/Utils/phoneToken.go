// phoneToken.go - Steam手机令牌(Steam Guard移动认证器)相关功能
// 实现Steam Guard移动认证器的核心算法，包括验证码生成和确认参数计算
package Utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"

	"github.com/JovanniChen/SteamDB/Steam/Param"
)

// chars Steam验证码使用的字符集
// 排除了容易混淆的字符(如0/O, 1/I/L等)，提高用户输入准确性
const chars = "23456789BCDFGHJKMNPQRTVWXY"

// MaFile Steam移动验证器文件结构
type MaFile struct {
	SharedSecret   string `json:"shared_secret"`
	IdentitySecret string `json:"identity_secret"`
	DeviceID       string `json:"device_id"`
	Session        struct {
		SteamID int64 `json:"steamid"`
	} `json:"Session"`
}

// PhoneToken 手机令牌处理器
type PhoneToken struct {
	MaFile *MaFile
}

// LoadMaFile 加载maFile文件
func LoadMaFile(fileContent string) (*PhoneToken, error) {
	var maFile MaFile
	if err := json.Unmarshal([]byte(fileContent), &maFile); err != nil {
		return nil, fmt.Errorf("解析maFile失败: %v", err)
	}

	return &PhoneToken{MaFile: &maFile}, nil
}

// GenerateConfirmationHashForTime 计算手机操作需要的密文
func (pt *PhoneToken) GenerateConfirmationHashForTime(timestamp int64, tag string) (string, error) {
	identitySecret, err := base64.StdEncoding.DecodeString(pt.MaFile.IdentitySecret)
	if err != nil {
		return "", fmt.Errorf("解析identity_secret失败: %v", err)
	}

	// 计算buffer大小
	bufferSize := 8
	if tag != "" {
		if len(tag) > 32 {
			bufferSize = 8 + 32
		} else {
			bufferSize = 8 + len(tag)
		}
	}

	// 创建buffer
	buffer := make([]byte, bufferSize)

	// 写入时间戳（大端序）
	for i := 7; i >= 0; i-- {
		buffer[i] = byte(timestamp & 0xFF)
		timestamp >>= 8
	}

	// 写入tag
	if tag != "" {
		copy(buffer[8:], []byte(tag))
	}

	// 计算HMAC
	mac := hmac.New(sha1.New, identitySecret)
	mac.Write(buffer)
	hashedData := mac.Sum(nil)

	// 编码并URL编码
	encodedData := base64.StdEncoding.EncodeToString(hashedData)

	return url.QueryEscape(encodedData), nil
}

// bufferizeSecret 将密钥字符串转换为字节数组
// 支持十六进制和Base64两种编码格式的密钥
// 参数：secret - 密钥字符串(hex或base64格式)
// 返回值：解码后的字节数组
func bufferizeSecret(secret string) []byte {
	if secret == "" {
		return []byte(secret)
	}

	// 如果是40字符长度，尝试作为十六进制解码
	if len(secret) == 40 {
		if _, err := hex.DecodeString(secret); err == nil {
			decoded, _ := hex.DecodeString(secret)
			return decoded
		}
	}

	// 否则作为Base64解码
	decoded, _ := base64.StdEncoding.DecodeString(secret)
	return decoded
}

// GenerateAuthCode 生成Steam Guard验证码
// 基于TOTP(Time-based One-Time Password)算法生成5位验证码
// 参数：
//
//	secret - Steam Guard共享密钥
//	time - Unix时间戳
//
// 返回值：5位验证码字符串
func GenerateAuthCode(secret string, time int64) string {
	// 1. 准备密钥字节数组
	secretBytes := bufferizeSecret(secret)

	// 2. 构建时间计数器(每30秒为一个周期)
	b := make([]byte, 8)

	// 高4字节置零(相当于Node.js的writeUInt32BE(0, 0))
	for i := 0; i < 4; i++ {
		b[i] = 0
	}

	// 低4字节存储时间计数器(相当于Node.js的writeUInt32BE(Math.floor(time / 30), 4))
	timeDiv30 := uint32(math.Floor(float64(time) / 30))
	for i := 3; i >= 0; i-- {
		b[4+i] = byte(timeDiv30 & 0xff)
		timeDiv30 >>= 8
	}

	// 3. 使用HMAC-SHA1计算哈希值
	hmacHash := hmac.New(sha1.New, secretBytes)
	hmacHash.Write(b)
	hmacData := hmacHash.Sum(nil)

	// 4. 动态截取(Dynamic Truncation)
	start := hmacData[19] & 0x0f         // 取最后一个字节的低4位作为偏移量
	hmacData = hmacData[start : start+4] // 截取4字节

	// 5. 转换为32位无符号整数
	var fullcode uint32
	for i := 0; i < 4; i++ {
		fullcode = (fullcode << 8) | uint32(hmacData[i])
	}
	fullcode &= 0x7fffffff // 清除符号位

	// 6. 转换为5位验证码
	code := ""
	for i := 0; i < 5; i++ {
		code += string(chars[fullcode%uint32(len(chars))])
		fullcode /= uint32(len(chars))
	}
	return code
}

// GenerateConfirmationQueryParams 生成确认查询参数
func (pt *PhoneToken) GenerateConfirmationQueryParams(timestamp int64, tag string) (Param.Params, error) {
	if pt.MaFile.DeviceID == "" {
		return nil, fmt.Errorf("设备ID不存在")
	}

	confirmationHash, err := pt.GenerateConfirmationHashForTime(timestamp, tag)
	if err != nil {
		return nil, err
	}

	params := Param.Params{}
	params.SetString("p", pt.MaFile.DeviceID)
	params.SetInt64("a", pt.MaFile.Session.SteamID)
	params.SetString("k", confirmationHash)
	params.SetInt64("t", timestamp)
	params.SetString("m", "react")
	params.SetString("tag", tag)

	return params, nil
}

// GenerateConfirmationQueryParams 生成市场确认操作的查询参数
// 用于Steam市场交易确认、礼品确认等操作
// 参数：
//
//	deviceID - 设备ID(Android设备标识符)
//	identitySecret - 身份验证密钥
//	steamid - Steam用户ID
//	time - Unix时间戳
//	tag - 操作标签("conf"用于确认, "cancel"用于取消等)
//
// 返回值：查询参数映射和可能的错误
func GenerateConfirmationQueryParams(deviceID, identitySecret, steamid string, time int64, tag string) (Param.Params, error) {
	if deviceID == "" {
		return nil, errors.New("Device ID is not present")
	}

	params := Param.Params{}
	params.SetString("p", deviceID)
	params.SetString("a", steamid)
	params.SetString("k", GenerateConfirmationHashForTime(identitySecret, time, tag))
	params.SetInt64("t", time)
	params.SetString("m", "react")
	params.SetString("tag", tag)

	// return map[string]string{
	// 	"p":   deviceID,                                                   // 设备ID
	// 	"a":   steamid,                                                    // Steam用户ID
	// 	"k":   GenerateConfirmationHashForTime(identitySecret, time, tag), // 验证哈希
	// 	"t":   strconv.FormatInt(time, 10),                                // 时间戳
	// 	"m":   "react",                                                    // 固定值
	// 	"tag": tag,                                                        // 操作标签
	// }, nil

	return params, nil
}

// GenerateConfirmationHashForTime 生成市场确认操作的验证哈希
// 计算手机令牌操作所需的HMAC签名，用于证明操作的合法性
// 参数：
//
//	identitySecret - Base64编码的身份验证密钥
//	time - Unix时间戳
//	tag - 操作标签字符串
//
// 返回值：Base64编码的HMAC-SHA1签名
func GenerateConfirmationHashForTime(identitySecret string, time int64, tag string) string {
	// 解码身份验证密钥
	decode, _ := base64.StdEncoding.DecodeString(identitySecret)

	// 计算数据数组长度(8字节时间戳 + 标签长度，最大32字节)
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

	// 构建待签名数据数组
	array := make([]byte, n2)

	// 前8字节存储时间戳(大端序)
	for i := 7; i >= 0; i-- {
		array[i] = byte(time & 0xff)
		time >>= 8
	}

	// 后续字节存储标签
	if tag != "" {
		copy(array[8:], tag)
	}

	// 使用HMAC-SHA1计算签名
	hmacHash := hmac.New(sha1.New, decode)
	hmacHash.Write(array)
	hashedData := hmacHash.Sum(nil)

	// 返回Base64编码的签名
	encodedData := url.QueryEscape(base64.StdEncoding.EncodeToString(hashedData))
	return fmt.Sprintf("%s", encodedData)
}
