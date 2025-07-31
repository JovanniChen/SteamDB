// Model包 - Steam登录相关的数据模型定义
// 包含登录过程中使用的各种数据结构和响应格式
package Model

import (
	"encoding/hex"
	"math/big"
	"strconv"
)

// SteamPublicKey Steam RSA公钥结构体
// 用于存储从Steam服务器获取的RSA公钥信息，用于密码加密
type SteamPublicKey struct {
	PublicKeyExp string `json:"publickey_exp,omitempty"`   // 公钥指数(十六进制字符串)
	PublicKeyMod string `json:"publickey_mod,omitempty"`   // 公钥模数(十六进制字符串)
	SteamID      uint64 `json:"steamid,string,omitempty"`  // Steam用户ID
	Success      bool   `json:"success"`                   // 请求是否成功
	Timestamp    uint64 `json:"timestamp,string,omitempty"` // 时间戳
	TokenGID     string `json:"token_gid,omitempty"`       // 令牌GID
}

// Modulus 获取RSA公钥的模数
// 将十六进制字符串转换为大整数，用于RSA加密计算
// 返回值：*big.Int类型的模数和可能的错误
func (spk SteamPublicKey) Modulus() (*big.Int, error) {
	// 将十六进制字符串解码为字节数组
	by, er := hex.DecodeString(spk.PublicKeyMod)
	if er != nil {
		return nil, er
	}
	
	// 创建大整数并设置字节值
	bi := big.NewInt(0)
	return bi.SetBytes(by), nil
}

// Exponent 获取RSA公钥的指数
// 将十六进制字符串转换为整数，用于RSA加密计算
// 返回值：int64类型的指数和可能的错误
func (spk SteamPublicKey) Exponent() (int64, error) {
	return strconv.ParseInt(spk.PublicKeyExp, 16, 0)
}

// message 通用响应消息结构体
// 大部分Steam API响应的基础结构
type message struct {
	Success bool   `json:"success"`           // 操作是否成功
	Message string `json:"message,omitempty"` // 错误或状态消息
}

// RefreshResponse 令牌刷新响应结构体
// 用于处理JWT令牌刷新操作的响应
type RefreshResponse struct {
	message
	Cookie map[string]string `json:"cookie,omitempty"` // 返回的Cookie信息
}

// LoginResponse 登录响应结构体
// 用于处理登录请求的响应数据
type LoginResponse struct {
	message
	Data struct {
		ClientID  string // 客户端ID
		RequestID string // 请求ID
		SteamID   string // Steam用户ID
	}
}

// FinalizeResponse 登录完成响应结构体
// 用于处理登录流程最终步骤的响应
type FinalizeResponse struct {
	message
	SteamID      string `json:"steamID"` // Steam用户ID
	TransferInfo []struct {              // 传输信息数组
		Url    string `json:"url"` // 传输URL
		Params struct {           // 传输参数
			Nonce string `json:"nonce"` // 一次性随机数
			Auth  string `json:"auth"`  // 认证信息
		} `json:"params"`
	} `json:"transfer_info"`
}

// CheckLoginResponse 检查登录状态响应结构体
// 用于验证登录状态和获取会话信息
type CheckLoginResponse struct {
	message
	Url  string `json:"url"` // 重定向URL
	Data struct {           // 登录数据
		SteamLoginSecure string `json:"steamLoginSecure"` // Steam安全登录令牌
		SessionId        string `json:"sessionid"`        // 会话ID
	}
}
