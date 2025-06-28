package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// LogLevel 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// LogLevelString 日志级别字符串映射
var LogLevelString = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// Logger 日志记录器
type Logger struct {
	level LogLevel
}

// NewLogger 创建新的日志记录器
func NewLogger(level LogLevel) *Logger {
	return &Logger{level: level}
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// shouldLog 判断是否应该记录日志
func (l *Logger) shouldLog(level LogLevel) bool {
	return level >= l.level
}

// formatMessage 格式化日志消息
func (l *Logger) formatMessage(level LogLevel, requestID, message string, args ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	levelStr := LogLevelString[level]
	
	// 格式化消息
	formattedMessage := message
	if len(args) > 0 {
		formattedMessage = fmt.Sprintf(message, args...)
	}
	
	// 如果有requestID，添加到日志中
	if requestID != "" {
		return fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, levelStr, requestID, formattedMessage)
	}
	
	return fmt.Sprintf("[%s] [%s] %s", timestamp, levelStr, formattedMessage)
}

// Debug 调试日志
func (l *Logger) Debug(requestID, message string, args ...interface{}) {
	if l.shouldLog(DEBUG) {
		log.Printf(l.formatMessage(DEBUG, requestID, message, args...))
	}
}

// Info 信息日志
func (l *Logger) Info(requestID, message string, args ...interface{}) {
	if l.shouldLog(INFO) {
		log.Printf(l.formatMessage(INFO, requestID, message, args...))
	}
}

// Warn 警告日志
func (l *Logger) Warn(requestID, message string, args ...interface{}) {
	if l.shouldLog(WARN) {
		log.Printf(l.formatMessage(WARN, requestID, message, args...))
	}
}

// Error 错误日志
func (l *Logger) Error(requestID, message string, args ...interface{}) {
	if l.shouldLog(ERROR) {
		log.Printf(l.formatMessage(ERROR, requestID, message, args...))
	}
}

// Fatal 致命错误日志
func (l *Logger) Fatal(requestID, message string, args ...interface{}) {
	if l.shouldLog(FATAL) {
		log.Printf(l.formatMessage(FATAL, requestID, message, args...))
		os.Exit(1)
	}
}

// 全局日志记录器实例
var globalLogger = NewLogger(INFO)

// SetGlobalLogLevel 设置全局日志级别
func SetGlobalLogLevel(level LogLevel) {
	globalLogger.SetLevel(level)
}

// GetRequestID 从gin上下文获取requestID
func GetRequestID(c *gin.Context) string {
	if c == nil {
		return ""
	}
	
	requestID, exists := c.Get("request_id")
	if !exists {
		return ""
	}
	
	if id, ok := requestID.(string); ok {
		return id
	}
	
	return ""
}

// 全局日志函数（简化调用）

// LogDebug 调试日志
func LogDebug(c *gin.Context, message string, args ...interface{}) {
	requestID := GetRequestID(c)
	globalLogger.Debug(requestID, message, args...)
}

// LogInfo 信息日志
func LogInfo(c *gin.Context, message string, args ...interface{}) {
	requestID := GetRequestID(c)
	globalLogger.Info(requestID, message, args...)
}

// LogWarn 警告日志
func LogWarn(c *gin.Context, message string, args ...interface{}) {
	requestID := GetRequestID(c)
	globalLogger.Warn(requestID, message, args...)
}

// LogError 错误日志
func LogError(c *gin.Context, message string, args ...interface{}) {
	requestID := GetRequestID(c)
	globalLogger.Error(requestID, message, args...)
}

// LogFatal 致命错误日志
func LogFatal(c *gin.Context, message string, args ...interface{}) {
	requestID := GetRequestID(c)
	globalLogger.Fatal(requestID, message, args...)
}

// 业务日志函数

// LogUserLogin 用户登录日志
func LogUserLogin(c *gin.Context, userID, username, email, ip string, success bool) {
	requestID := GetRequestID(c)
	status := "失败"
	if success {
		status = "成功"
	}
	
	message := fmt.Sprintf("用户登录 %s - 用户ID: %s, 用户名: %s, 邮箱: %s, IP: %s", 
		status, userID, username, MaskEmail(email), ip)
	
	if success {
		globalLogger.Info(requestID, message)
	} else {
		globalLogger.Warn(requestID, message)
	}
}

// LogUserRegister 用户注册日志
func LogUserRegister(c *gin.Context, userID, username, email, ip string) {
	requestID := GetRequestID(c)
	message := fmt.Sprintf("用户注册 - 用户ID: %s, 用户名: %s, 邮箱: %s, IP: %s", 
		userID, username, MaskEmail(email), ip)
	globalLogger.Info(requestID, message)
}

// LogWalletOperation 钱包操作日志
func LogWalletOperation(c *gin.Context, userID, operation, amount, description string) {
	requestID := GetRequestID(c)
	message := fmt.Sprintf("钱包操作 - 用户ID: %s, 操作: %s, 金额: %s, 描述: %s", 
		userID, operation, amount, description)
	globalLogger.Info(requestID, message)
}

// LogDatabaseError 数据库错误日志
func LogDatabaseError(c *gin.Context, operation string, err error) {
	requestID := GetRequestID(c)
	message := fmt.Sprintf("数据库错误 - 操作: %s, 错误: %v", operation, err)
	globalLogger.Error(requestID, message)
}

// LogSecurityEvent 安全事件日志
func LogSecurityEvent(c *gin.Context, event, details string) {
	requestID := GetRequestID(c)
	message := fmt.Sprintf("安全事件 - 事件: %s, 详情: %s", event, details)
	globalLogger.Warn(requestID, message)
} 