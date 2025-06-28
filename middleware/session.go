package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// SessionMiddleware 会话管理中间件
func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		requestID := generateRequestID()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// 设置请求开始时间
		c.Set("start_time", time.Now())

		// 检查登录态
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// 有token，尝试验证
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
				tokenString := tokenParts[1]
				claims, err := utils.ValidateToken(tokenString)
				if err == nil {
					// token有效，设置用户信息
					c.Set("user_id", claims.UserID)
					c.Set("username", claims.Username)
					c.Set("claims", claims)
					c.Set("is_authenticated", true)
				} else {
					// token无效，但不阻止请求继续
					c.Set("is_authenticated", false)
				}
			} else {
				c.Set("is_authenticated", false)
			}
		} else {
			c.Set("is_authenticated", false)
		}

		c.Next()
	}
}

// CheckLoginStatus 检查登录状态的中间件
func CheckLoginStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAuthenticated := IsAuthenticated(c)

		// 将登录状态添加到响应头中
		if isAuthenticated {
			c.Header("X-Auth-Status", "authenticated")
			c.Header("X-User-ID", strconv.FormatUint(uint64(GetCurrentUser(c)), 10))
		} else {
			c.Header("X-Auth-Status", "unauthenticated")
		}

		c.Next()
	}
}

// RequireLogin 要求登录的中间件
func RequireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsAuthenticated(c) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":      401,
				"message":   "请先登录",
				"error":     "LOGIN_REQUIRED",
				"timestamp": time.Now().Unix(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// OptionalLogin 可选登录的中间件
func OptionalLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不强制要求登录，但会设置用户信息（如果有的话）
		c.Next()
	}
}

// GetLoginStatus 获取登录状态
func GetLoginStatus(c *gin.Context) gin.H {
	isAuthenticated := IsAuthenticated(c)

	status := gin.H{
		"is_authenticated": isAuthenticated,
		"timestamp":        time.Now().Unix(),
	}

	if isAuthenticated {
		status["user_id"] = GetCurrentUser(c)
		status["username"] = GetCurrentUsername(c)
	}

	return status
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + utils.RandomString(8)
}
