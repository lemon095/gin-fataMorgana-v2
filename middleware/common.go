package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gin-fataMorgana/utils"
)

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
	Total    int64 `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ResponseWrapper 响应包装器
type ResponseWrapper struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 记录请求日志
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if raw != "" {
			path = path + "?" + raw
		}

		utils.LogInfo(c, "HTTP Request - method: %s, path: %s, status: %d, latency: %v, client_ip: %s, user_agent: %s",
			method, path, status, latency, clientIP, c.Request.UserAgent())
	}
}

// ErrorHandlerMiddleware 错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			utils.LogError(c, "Request error - error: %s, path: %s", err.Error(), c.Request.URL.Path)
			
			// 根据错误类型返回相应的响应
			if appErr, ok := err.Err.(*utils.AppError); ok {
				c.JSON(http.StatusOK, ResponseWrapper{
					Code:    appErr.Code,
					Message: appErr.Message,
					TraceID: getTraceID(c),
				})
				return
			}

			// 默认错误响应
			c.JSON(http.StatusInternalServerError, ResponseWrapper{
				Code:    500,
				Message: "Internal server error",
				TraceID: getTraceID(c),
			})
		}
	}
}

// PaginationMiddleware 分页中间件
func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pagination PaginationInfo
		
		// 从查询参数获取分页信息
		if err := c.ShouldBindQuery(&pagination); err != nil {
			pagination = PaginationInfo{Page: 1, PageSize: 10}
		}

		// 设置默认值
		if pagination.Page <= 0 {
			pagination.Page = 1
		}
		if pagination.PageSize <= 0 {
			pagination.PageSize = 10
		}
		if pagination.PageSize > 100 {
			pagination.PageSize = 100
		}

		c.Set("pagination", pagination)
		c.Next()
	}
}

// GetPagination 获取分页信息
func GetPagination(c *gin.Context) PaginationInfo {
	if pagination, exists := c.Get("pagination"); exists {
		return pagination.(PaginationInfo)
	}
	return PaginationInfo{Page: 1, PageSize: 10}
}

// ApplyPagination 应用分页到数据库查询
func ApplyPagination(db *gorm.DB, pagination PaginationInfo) *gorm.DB {
	offset := (pagination.Page - 1) * pagination.PageSize
	return db.Offset(offset).Limit(pagination.PageSize)
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResponseWrapper{
		Code:    0,
		Message: "success",
		Data:    data,
		TraceID: getTraceID(c),
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, ResponseWrapper{
		Code:    code,
		Message: message,
		TraceID: getTraceID(c),
	})
}

// getTraceID 获取追踪ID
func getTraceID(c *gin.Context) string {
	if traceID, exists := c.Get("trace_id"); exists {
		return traceID.(string)
	}
	return ""
}

// SetTraceID 设置追踪ID
func SetTraceID(c *gin.Context, traceID string) {
	c.Set("trace_id", traceID)
}

// GetUserID 获取用户ID
func GetUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, utils.NewAppError(401, "User not authenticated")
	}
	
	// 处理不同类型的user_id
	switch v := userID.(type) {
	case uint:
		return v, nil
	case int64:
		return uint(v), nil
	case int:
		return uint(v), nil
	case float64:
		return uint(v), nil
	default:
		return 0, utils.NewAppError(400, "Invalid user ID type")
	}
}

// GetUserIDFromParam 从路径参数获取用户ID
func GetUserIDFromParam(c *gin.Context) (int64, error) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		return 0, utils.NewAppError(400, "User ID is required")
	}
	
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, utils.NewAppError(400, "Invalid user ID format")
	}
	
	return userID, nil
} 