package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseCode 响应码定义
const (
	// 成功响应码
	CodeSuccess = 0

	// 客户端错误码 (1000-1999)
	CodeInvalidParams       = 1000 // 参数错误
	CodeValidationFailed    = 1001 // 验证失败
	CodeUnauthorized        = 1002 // 未授权
	CodeForbidden           = 1003 // 禁止访问
	CodeNotFound            = 1004 // 资源不存在
	CodeMethodNotAllowed    = 1005 // 方法不允许
	CodeRequestTimeout      = 1006 // 请求超时
	CodeTooManyRequests     = 1007 // 请求过于频繁
	CodeConflict            = 1008 // 资源冲突
	CodeUnprocessableEntity = 1009 // 无法处理的实体

	// 认证相关错误码 (2000-2099)
	CodeTokenExpired       = 2000 // Token过期
	CodeTokenInvalid       = 2001 // Token无效
	CodeTokenMissing       = 2002 // Token缺失
	CodeLoginFailed        = 2003 // 登录失败
	CodeUserNotFound       = 2004 // 用户不存在
	CodePasswordIncorrect  = 2005 // 密码错误
	CodeUserAlreadyExists  = 2006 // 用户已存在
	CodeEmailAlreadyExists = 2007 // 邮箱已存在
	CodeInviteCodeInvalid  = 2008 // 邀请码无效
	CodeAccountLocked      = 2009 // 账户被锁定
	CodeSessionExpired     = 2010 // 会话过期

	// 业务逻辑错误码 (3000-3999)
	CodeOperationFailed       = 3000 // 操作失败
	CodeResourceBusy          = 3001 // 资源繁忙
	CodeInsufficientFunds     = 3002 // 余额不足
	CodeLimitExceeded         = 3003 // 超出限制
	CodeInvalidOperation      = 3004 // 无效操作
	CodeDataInconsistent      = 3005 // 数据不一致
	CodeBusinessRuleViolation = 3006 // 违反业务规则

	// 服务器错误码 (5000-5999)
	CodeInternalError      = 5000 // 内部服务器错误
	CodeDatabaseError      = 5001 // 数据库错误
	CodeRedisError         = 5002 // Redis错误
	CodeExternalAPIError   = 5003 // 外部API错误
	CodeServiceUnavailable = 5004 // 服务不可用
	CodeGatewayTimeout     = 5005 // 网关超时
	CodeConfigError        = 5006 // 配置错误
)

// ResponseMessage 响应消息映射
var ResponseMessage = map[int]string{
	// 成功响应
	CodeSuccess: "操作成功",

	// 客户端错误
	CodeInvalidParams:       "参数错误",
	CodeValidationFailed:    "数据验证失败",
	CodeUnauthorized:        "未授权访问",
	CodeForbidden:           "禁止访问",
	CodeNotFound:            "资源不存在",
	CodeMethodNotAllowed:    "请求方法不允许",
	CodeRequestTimeout:      "请求超时",
	CodeTooManyRequests:     "请求过于频繁，请稍后再试",
	CodeConflict:            "资源冲突",
	CodeUnprocessableEntity: "无法处理的请求",

	// 认证相关错误
	CodeTokenExpired:       "Token已过期",
	CodeTokenInvalid:       "Token无效",
	CodeTokenMissing:       "Token缺失",
	CodeLoginFailed:        "登录失败",
	CodeUserNotFound:       "用户不存在",
	CodePasswordIncorrect:  "密码错误",
	CodeUserAlreadyExists:  "用户已存在",
	CodeEmailAlreadyExists: "邮箱已被注册",
	CodeInviteCodeInvalid:  "邀请码无效",
	CodeAccountLocked:      "账户已被锁定",
	CodeSessionExpired:     "会话已过期",

	// 业务逻辑错误
	CodeOperationFailed:       "操作失败",
	CodeResourceBusy:          "资源繁忙，请稍后再试",
	CodeInsufficientFunds:     "余额不足",
	CodeLimitExceeded:         "超出限制",
	CodeInvalidOperation:      "无效操作",
	CodeDataInconsistent:      "数据不一致",
	CodeBusinessRuleViolation: "违反业务规则",

	// 服务器错误
	CodeInternalError:      "内部服务器错误",
	CodeDatabaseError:      "数据库错误",
	CodeRedisError:         "Redis错误",
	CodeExternalAPIError:   "外部服务错误",
	CodeServiceUnavailable: "服务暂时不可用",
	CodeGatewayTimeout:     "网关超时",
	CodeConfigError:        "配置错误",
}

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 响应码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	response := Response{
		Code:    CodeSuccess,
		Message: ResponseMessage[CodeSuccess],
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	response := Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

// Error 错误响应
func Error(c *gin.Context, code int) {
	message, exists := ResponseMessage[code]
	if !exists {
		message = "未知错误"
	}

	response := Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}

	// 根据错误码确定HTTP状态码
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, response)
}

// ErrorWithMessage 带自定义消息的错误响应
func ErrorWithMessage(c *gin.Context, code int, message string) {
	response := Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}

	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, response)
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *gin.Context, code int, data interface{}) {
	message, exists := ResponseMessage[code]
	if !exists {
		message = "未知错误"
	}

	response := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}

	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, response)
}

// getHTTPStatus 根据业务错误码获取HTTP状态码
func getHTTPStatus(code int) int {
	switch {
	case code == CodeSuccess:
		return http.StatusOK
	case code >= 1000 && code < 2000:
		return http.StatusBadRequest
	case code >= 2000 && code < 2100:
		return http.StatusUnauthorized
	case code >= 3000 && code < 4000:
		return http.StatusUnprocessableEntity
	case code >= 5000 && code < 6000:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// 常用响应函数

// InvalidParams 参数错误
func InvalidParams(c *gin.Context) {
	Error(c, CodeInvalidParams)
}

// InvalidParamsWithMessage 带消息的参数错误
func InvalidParamsWithMessage(c *gin.Context, message string) {
	ErrorWithMessage(c, CodeInvalidParams, message)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context) {
	Error(c, CodeUnauthorized)
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
	Error(c, CodeInternalError)
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
	Error(c, CodeLoginFailed)
}

// UserNotFound 用户不存在
func UserNotFound(c *gin.Context) {
	Error(c, CodeUserNotFound)
}

// UserAlreadyExists 用户已存在
func UserAlreadyExists(c *gin.Context) {
	Error(c, CodeUserAlreadyExists)
}

// EmailAlreadyExists 邮箱已存在
func EmailAlreadyExists(c *gin.Context) {
	Error(c, CodeEmailAlreadyExists)
}

// TokenExpired Token过期
func TokenExpired(c *gin.Context) {
	Error(c, CodeTokenExpired)
}

// TokenInvalid Token无效
func TokenInvalid(c *gin.Context) {
	Error(c, CodeTokenInvalid)
}

// InviteCodeInvalid 邀请码无效
func InviteCodeInvalid(c *gin.Context) {
	Error(c, CodeInviteCodeInvalid)
}

// AccountLocked 账户被锁定
func AccountLocked(c *gin.Context) {
	Error(c, CodeAccountLocked)
}

// SessionExpired 会话过期
func SessionExpired(c *gin.Context) {
	Error(c, CodeSessionExpired)
}
