package Utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

// 生成十六进制随机字符串（长度=2*n）
func SafeHexString(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

// 生成自定义字符集的随机字符串
func SafeCustomString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[idx.Int64()]
	}
	return string(b)
}
