package Logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// LogLevel 日志级别枚举
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// Logger 日志器结构体
type Logger struct {
	level  LogLevel
	logger *log.Logger
}

// 全局日志器实例
var globalLogger *Logger

// init 初始化全局日志器
func init() {
	globalLogger = &Logger{
		level:  DEBUG,                     // 默认日志级别为 DEBUG，便于开发调试
		logger: log.New(os.Stdout, "", 0), // 去掉默认的标志，我们自己格式化
	}
}

// SetLevel 设置日志级别
func SetLevel(level LogLevel) {
	globalLogger.level = level
}

// SetOutput 设置日志输出位置
func SetOutput(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	globalLogger.logger.SetOutput(file)
	return nil
}

// getCallerInfo 获取调用者信息
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2) // 跳过当前函数和日志方法
	if !ok {
		return "unknown:0"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// formatLog 格式化日志消息
func formatLog(level string, callerInfo string, message string) string {
	now := time.Now().Format("2006/01/02 15:04:05")
	return fmt.Sprintf("[SteamDB] [%s] [%s] %s: %s", now, level, callerInfo, message)
}

// Debug 输出调试日志
func Debug(v ...any) {
	if globalLogger.level <= DEBUG {
		callerInfo := getCallerInfo()
		message := fmt.Sprint(v...)
		logMessage := formatLog("DEBUG", callerInfo, message)
		globalLogger.logger.Println(logMessage)
		globalLogger.logger.Println()
	}
}

// Debugf 格式化输出调试日志
func Debugf(format string, v ...any) {
	if globalLogger.level <= DEBUG {
		callerInfo := getCallerInfo()
		message := fmt.Sprintf(format, v...)
		logMessage := formatLog("DEBUG", callerInfo, message)
		globalLogger.logger.Println(logMessage)
	}
}

// Info 输出信息日志
func Info(v ...any) {
	if globalLogger.level <= INFO {
		callerInfo := getCallerInfo()
		message := fmt.Sprint(v...)
		logMessage := formatLog("INFO", callerInfo, message)
		globalLogger.logger.Println(logMessage)
	}
}

// Infof 格式化输出信息日志
func Infof(format string, v ...any) {
	if globalLogger.level <= INFO {
		callerInfo := getCallerInfo()
		message := fmt.Sprintf(format, v...)
		logMessage := formatLog("INFO", callerInfo, message)
		globalLogger.logger.Println(logMessage)
	}
}

// Warn 输出警告日志
func Warn(v ...any) {
	if globalLogger.level <= WARN {
		callerInfo := getCallerInfo()
		message := fmt.Sprint(v...)
		logMessage := formatLog("WARN", callerInfo, message)
		globalLogger.logger.Println(logMessage)
	}
}

// Warnf 格式化输出警告日志
func Warnf(format string, v ...any) {
	if globalLogger.level <= WARN {
		callerInfo := getCallerInfo()
		message := fmt.Sprintf(format, v...)
		logMessage := formatLog("WARN", callerInfo, message)
		globalLogger.logger.Println(logMessage)
	}
}

// Error 输出错误日志
func Error(v ...any) {
	if globalLogger.level <= ERROR {
		callerInfo := getCallerInfo()
		message := fmt.Sprint(v...)
		logMessage := formatLog("ERROR", callerInfo, message)
		globalLogger.logger.Println(logMessage)
	}
}

// Errorf 格式化输出错误日志
func Errorf(format string, v ...any) {
	if globalLogger.level <= ERROR {
		callerInfo := getCallerInfo()
		message := fmt.Sprintf(format, v...)
		logMessage := formatLog("ERROR", callerInfo, message)
		globalLogger.logger.Println(logMessage)
	}
}

// Print 输出普通日志（兼容 fmt.Print）
func Print(v ...any) {
	callerInfo := getCallerInfo()
	message := fmt.Sprint(v...)
	logMessage := formatLog("", callerInfo, message)
	globalLogger.logger.Println(logMessage)
}

// Printf 格式化输出普通日志（兼容 fmt.Printf）
func Printf(format string, v ...any) {
	callerInfo := getCallerInfo()
	message := fmt.Sprintf(format, v...)
	logMessage := formatLog("", callerInfo, message)
	globalLogger.logger.Println(logMessage)
}

// Println 输出普通日志并换行（兼容 fmt.Println）
func Println(v ...any) {
	callerInfo := getCallerInfo()
	message := fmt.Sprint(v...)
	logMessage := formatLog("", callerInfo, message)
	globalLogger.logger.Println(logMessage)
}
