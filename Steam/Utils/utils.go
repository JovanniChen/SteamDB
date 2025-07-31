// Utils包 - 通用工具函数集合
// 提供安全的随机字符串生成功能，用于密码学和安全相关操作
package Utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
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
//   length - 生成字符串的长度
//   charset - 字符集（例如："0123456789"、"abcdefghijklmnopqrstuvwxyz"）
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
