package Errors

import (
	"fmt"
	"runtime"
	"strconv"
)

var (
	errorGettingKey   = "异常错误, 原因: %s"
	errorResponseCode = "请求返回异常 responseCode: %s"
	unavailable       = "异常错误"
)

func printStack(errKey, errStr string) error {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	message := fmt.Sprintf(errKey, errStr)
	message += "\n" + string(buf[:n])
	return fmt.Errorf(errKey, message)
}

func Error(errStr string) error {
	return printStack(errorGettingKey, errStr)
}

func ResponseError(errStr int) error {
	return printStack(errorResponseCode, strconv.Itoa(errStr))
}

func Unavailable() error {
	return printStack(unavailable, "")
}
