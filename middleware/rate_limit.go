package middleware

import (
	"sync"
	"time"

	"gin-fataMorgana/utils"

	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// RateLimiter 限流器结构
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int           // 时间窗口内允许的请求数
	window   time.Duration // 时间窗口
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// isAllowed 检查是否允许请求
func (rl *RateLimiter) isAllowed(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// 获取该key的请求记录
	requests, exists := rl.requests[key]
	if !exists {
		requests = []time.Time{}
	}

	// 清理过期的请求记录
	var validRequests []time.Time
	for _, reqTime := range requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// 检查是否超过限制
	if len(validRequests) >= rl.limit {
		return false
	}

	// 添加当前请求
	validRequests = append(validRequests, now)
	rl.requests[key] = validRequests

	return true
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		// 跳过OPTIONS请求的限流检查
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 获取客户端IP
		clientIP := c.ClientIP()
		if clientIP == "" {
			clientIP = "unknown"
		}

		// 检查是否允许请求
		if !limiter.isAllowed(clientIP) {
			utils.ErrorWithMessage(c, utils.CodeRateLimitExceeded, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// 预定义的限流配置
const (
	// 登录限流：每分钟最多5次
	LoginRateLimit  = 5
	LoginRateWindow = 1 * time.Minute

	// 注册限流：每小时最多3次
	RegisterRateLimit  = 3
	RegisterRateWindow = 1 * time.Hour

	// 提现限流：每小时最多2次
	WithdrawRateLimit  = 2
	WithdrawRateWindow = 1 * time.Hour

	// 通用限流：每分钟最多60次
	GeneralRateLimit  = 60
	GeneralRateWindow = 1 * time.Minute
)

// LoginRateLimitMiddleware 登录限流中间件（账号级别，优先用account参数）
func LoginRateLimitMiddleware() gin.HandlerFunc {
	limiter := NewRateLimiter(10, 1*time.Minute) // 1分钟10次

	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		var key string
		if c.Request.URL.Path == "/api/v2/auth/login" && c.Request.Method == "POST" {
			raw, _ := ioutil.ReadAll(c.Request.Body)
			var body struct {
				Account string `json:"account"`
			}
			_ = json.Unmarshal(raw, &body)
			if body.Account != "" {
				key = body.Account
			}
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(raw))
		}
		if key == "" {
			key = c.ClientIP()
		}

		if !limiter.isAllowed(key) {
			utils.ErrorWithMessage(c, utils.CodeRateLimitExceeded, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// 通用账号级别限流中间件（优先用account参数）
func AccountRateLimitMiddleware() gin.HandlerFunc {
	limiter := NewRateLimiter(60, 1*time.Minute) // 1分钟60次

	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		var key string
		// 尝试从请求体获取account字段
		raw, _ := ioutil.ReadAll(c.Request.Body)
		var body struct {
			Account string `json:"account"`
		}
		_ = json.Unmarshal(raw, &body)
		if body.Account != "" {
			key = body.Account
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(raw))
		if key == "" {
			key = c.ClientIP()
		}

		if !limiter.isAllowed(key) {
			utils.ErrorWithMessage(c, utils.CodeRateLimitExceeded, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// 统一注册、提现、通用限流为账号级别
func RegisterRateLimitMiddleware() gin.HandlerFunc {
	return AccountRateLimitMiddleware()
}
func WithdrawRateLimitMiddleware() gin.HandlerFunc {
	return AccountRateLimitMiddleware()
}
func GeneralRateLimitMiddleware() gin.HandlerFunc {
	return AccountRateLimitMiddleware()
}
