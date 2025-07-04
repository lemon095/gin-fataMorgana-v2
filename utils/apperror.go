package utils

// AppError 应用错误类型
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError 创建新的应用错误
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewAppErrorWithCode 根据错误码创建应用错误
func NewAppErrorWithCode(code int) *AppError {
	message := ResponseMessage[code]
	if message == "" {
		message = "未知错误"
	}
	return &AppError{
		Code:    code,
		Message: message,
	}
} 