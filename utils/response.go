package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 完整的错误码定义
const (
	CodeSuccess         = 0    // 成功
	CodeError           = 1    // 一般错误
	CodeAuth            = 401  // 认证错误
	CodeForbidden       = 403  // 禁止访问
	CodeNotFound        = 404  // 资源不存在
	CodeValidation      = 422  // 数据验证错误
	CodeServer          = 500  // 服务器错误
	CodeDatabaseError   = 1001 // 数据库错误
	CodeRedisError      = 1002 // Redis错误
	CodeInvalidParams   = 1003 // 参数错误
	CodeOperationFailed = 1004 // 操作失败
	CodeUserNotFound    = 1005 // 用户不存在
	CodeUserAlreadyExists = 1006 // 用户已存在
	CodeValidationFailed = 1007 // 验证失败
	CodeAccountLocked   = 1008 // 账户锁定
	CodeRegisterClosed  = 1009 // 注册关闭
)

// ResponseMessage 完整的响应消息映射
var ResponseMessage = map[int]string{
	CodeSuccess:         "操作成功",
	CodeError:           "操作失败",
	CodeAuth:            "认证失败",
	CodeForbidden:       "禁止访问",
	CodeNotFound:        "资源不存在",
	CodeValidation:      "数据验证失败",
	CodeServer:          "服务器内部错误",
	CodeDatabaseError:   "数据库操作失败",
	CodeRedisError:      "Redis操作失败",
	CodeInvalidParams:   "参数错误",
	CodeOperationFailed: "操作失败",
	CodeUserNotFound:    "用户不存在",
	CodeUserAlreadyExists: "用户已存在",
	CodeValidationFailed: "验证失败",
	CodeAccountLocked:   "账户已被锁定",
	CodeRegisterClosed:  "当前系统不允许注册",
}

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      CodeSuccess,
		Message:   ResponseMessage[CodeSuccess],
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      CodeSuccess,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// Error 错误响应
func Error(c *gin.Context, code int) {
	message := ResponseMessage[code]
	if message == "" {
		message = "未知错误"
	}

	c.JSON(getHTTPStatus(code), Response{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}

// ErrorWithMessage 带自定义消息的错误响应
func ErrorWithMessage(c *gin.Context, code int, message string) {
	c.JSON(getHTTPStatus(code), Response{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *gin.Context, code int, data interface{}) {
	message := ResponseMessage[code]
	if message == "" {
		message = "未知错误"
	}

	c.JSON(getHTTPStatus(code), Response{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// getHTTPStatus 根据错误码获取HTTP状态码
func getHTTPStatus(code int) int {
	switch code {
	case CodeSuccess:
		return http.StatusOK
	case CodeAuth:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeValidation:
		return http.StatusUnprocessableEntity
	case CodeServer:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

// 常用响应函数

// InvalidParams 参数错误
func InvalidParams(c *gin.Context) {
	Error(c, CodeValidation)
}

// InvalidParamsWithMessage 带消息的参数错误
func InvalidParamsWithMessage(c *gin.Context, message string) {
	ErrorWithMessage(c, CodeValidation, message)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context) {
	Error(c, CodeAuth)
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context) {
	Error(c, CodeForbidden)
}

// NotFound 资源不存在
func NotFound(c *gin.Context) {
	Error(c, CodeNotFound)
}

// InternalError 内部服务器错误
func InternalError(c *gin.Context) {
	Error(c, CodeServer)
}

// DatabaseError 数据库错误
func DatabaseError(c *gin.Context) {
	Error(c, CodeDatabaseError)
}

// RedisError Redis错误
func RedisError(c *gin.Context) {
	Error(c, CodeRedisError)
}

// LoginFailed 登录失败
func LoginFailed(c *gin.Context) {
	ErrorWithMessage(c, CodeAuth, "邮箱或密码错误")
}

// UserNotFound 用户不存在
func UserNotFound(c *gin.Context) {
	ErrorWithMessage(c, CodeNotFound, "用户不存在")
}

// UserAlreadyExists 用户已存在
func UserAlreadyExists(c *gin.Context) {
	ErrorWithMessage(c, CodeValidation, "用户已存在")
}

// EmailAlreadyExists 邮箱已存在
func EmailAlreadyExists(c *gin.Context) {
	ErrorWithMessage(c, CodeValidation, "邮箱已被注册")
}

// TokenExpired Token过期
func TokenExpired(c *gin.Context) {
	ErrorWithMessage(c, CodeAuth, "Token已过期")
}

// TokenInvalid Token无效
func TokenInvalid(c *gin.Context) {
	ErrorWithMessage(c, CodeAuth, "Token无效")
}

// InviteCodeInvalid 邀请码无效
func InviteCodeInvalid(c *gin.Context) {
	ErrorWithMessage(c, CodeValidation, "邀请码无效")
}

// AccountLocked 账户锁定
func AccountLocked(c *gin.Context) {
	ErrorWithMessage(c, CodeAccountLocked, "账户已被锁定")
}
