// Utils包 - 通用工具函数集合
// 提供安全的随机字符串生成功能，用于密码学和安全相关操作
package Utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SafeHexString 生成安全的十六进制随机字符串
// 使用加密安全的随机数生成器，适用于密钥、令牌等安全场景
// 参数：n - 字节长度（最终字符串长度为2*n）
// 返回值：十六进制编码的随机字符串
// 示例：SafeHexString(16) 生成32字符的十六进制字符串
func SafeHexString(n int) string {
	bytes := make([]byte, n)

	// 使用加密安全的随机数填充字节数组
	if _, err := rand.Read(bytes); err != nil {
		panic(err) // 随机数生成失败时终止程序，确保安全性
	}

	// 将字节数组编码为十六进制字符串
	return hex.EncodeToString(bytes)
}

// SafeCustomString 生成自定义字符集的安全随机字符串
// 从指定字符集中随机选择字符组成字符串，适用于密码、验证码等场景
// 参数：
//
//	length - 生成字符串的长度
//	charset - 字符集（例如："0123456789"、"abcdefghijklmnopqrstuvwxyz"）
//
// 返回值：由指定字符集组成的随机字符串
// 示例：SafeCustomString(8, "0123456789") 生成8位数字验证码
func SafeCustomString(length int, charset string) string {
	b := make([]byte, length)

	// 为每个位置从字符集中随机选择字符
	for i := range b {
		// 生成[0, len(charset))范围内的安全随机数
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[idx.Int64()]
	}

	return string(b)
}

func WalletConvert(wallet string) int {
	// 正则表达式匹配数字（整数或小数）
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	// 查找所有匹配的数字
	numbers := re.FindAllString(wallet, -1)
	// 将所有数字拼接起来
	moneyStr := strings.Join(numbers, "")

	// 转换为浮点数
	money, err := strconv.ParseFloat(moneyStr, 64)
	if err != nil {
		// 如果转换失败，返回0
		return 0
	}

	// 乘以100后转换为整数（以分为单位）
	return int(money * 100)
}

func FriendCodeToSteamID64(friendCode uint32) uint64 {
	const base = uint64(76561197960265728)
	return base + uint64(friendCode)
}

func SteamID64ToFriendCode(id64 uint64) uint32 {
	const base = uint64(76561197960265728)
	return uint32(id64 - base)
}

func StructToURLValues(input any) url.Values {
	values := url.Values{}
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("form")

		if tag == "" || tag == "-" {
			continue
		}

		addField(values, tag, field)
	}

	return values
}

// addField 处理各种类型
func addField(values url.Values, key string, v reflect.Value) {
	if !v.IsValid() {
		return
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		values.Set(key, v.String())

	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		values.Set(key, strconv.FormatInt(v.Int(), 10))

	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		values.Set(key, strconv.FormatUint(v.Uint(), 10))

	case reflect.Float64, reflect.Float32:
		values.Set(key, strconv.FormatFloat(v.Float(), 'f', -1, 64))

	case reflect.Bool:
		if v.Bool() {
			values.Set(key, "1")
		} else {
			values.Set(key, "0")
		}

	case reflect.Struct:
		if t, ok := v.Interface().(time.Time); ok {
			values.Set(key, t.Format(time.RFC3339))
		} else {
			// 嵌套 struct：递归
			t := v.Type()
			for i := 0; i < v.NumField(); i++ {
				subField := v.Field(i)
				subType := t.Field(i)
				subTag := subType.Tag.Get("form")

				if subTag != "" && subTag != "-" {
					addField(values, subTag, subField)
				}
			}
		}
	}
}
