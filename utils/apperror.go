package utils

import "fmt"

// AppError 统一业务错误类型
// Code: 业务错误码，Message: 用户可读消息，Err: 底层错误
// 可用于service层返回

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[code=%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[code=%d] %s", e.Code, e.Message)
}

// NewAppError 创建业务错误
func NewAppError(code int, msg string, err error) *AppError {
	return &AppError{Code: code, Message: msg, Err: err}
}

// WrapAppError 包装底层错误
func WrapAppError(code int, msg string, err error) *AppError {
	return &AppError{Code: code, Message: msg, Err: err}
}

// IsAppError 判断是否为AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError 获取AppError及其属性
func GetAppError(err error) (int, string) {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code, appErr.Message
	}
	return 1, "操作失败"
}
