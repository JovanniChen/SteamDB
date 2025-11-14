// Errors包 - 错误处理和堆栈跟踪功能
// 提供统一的错误处理机制，包含详细的堆栈信息用于调试
package Errors

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
)

// 错误消息模板定义
var (
	errorGettingKey   = "异常错误, 原因: %s"            // 通用错误消息模板
	errorResponseCode = "请求返回异常 responseCode: %s" // HTTP响应错误模板
	unavailable       = "异常错误"                    // 服务不可用错误消息
)

// printStack 打印错误堆栈信息
// 获取当前调用堆栈并格式化错误信息，便于调试和问题定位
// 参数：
//
//	errKey - 错误消息模板
//	errStr - 具体错误内容
//
// 返回值：包含堆栈信息的格式化错误对象
func printStack(errKey, errStr string) error {
	var buf [4096]byte // 堆栈缓冲区，最大4KB

	// 获取调用堆栈信息
	n := runtime.Stack(buf[:], false) // false表示只获取当前goroutine的堆栈

	// 格式化错误消息
	message := fmt.Sprintf(errKey, errStr)
	message += "\n" + string(buf[:n]) // 添加堆栈信息

	return fmt.Errorf(errKey, message)
}

// Error 创建通用错误
// 用于创建包含堆栈信息的通用错误对象
// 参数：errStr - 错误描述字符串
// 返回值：包含堆栈信息的错误对象
func Error(errStr string) error {
	return printStack(errorGettingKey, errStr)
}

// ResponseError 创建HTTP响应错误
// 用于处理HTTP请求响应异常的情况
// 参数：errStr - HTTP状态码
// 返回值：包含堆栈信息的HTTP错误对象
func ResponseError(errStr int) error {
	return printStack(errorResponseCode, strconv.Itoa(errStr))
}

// Unavailable 创建服务不可用错误
// 用于表示服务暂时不可用的情况
// 返回值：服务不可用错误对象
func Unavailable() error {
	return printStack(unavailable, "")
}

// ErrRateLimited 表示API请求速率受限（HTTP 429）
var ErrNewRequest = errors.New("程序内部错误[NewRequest]")
var ErrRetryRequest = errors.New("程序内部错误[RetryRequest]")
var ErrGzipReader = errors.New("程序内部错误[GzipReader]")
var ErrIOReadAll = errors.New("程序内部错误[IOReadAll]")
var ErrRateLimited = errors.New("steam api rate limited (429)")
var ErrServerError = errors.New("steam server error(502)")
var ErrAuthorizationFailed = errors.New("steam authorization failed(401/403)")
var ErrAccountBan = errors.New("您的帐户当前无法使用社区市场。")

func IsAccountBan(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrAccountBan)
}

// IsRateLimitError 检查错误是否为速率限制错误
func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrRateLimited)
}

func IsServerError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrServerError)
}
