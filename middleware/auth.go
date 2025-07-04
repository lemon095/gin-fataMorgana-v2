package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"gin-fataMorgana/database"
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "缺少认证令牌",
				"error":   "MISSING_TOKEN",
			})
			c.Abort()
			return
		}

		// 检查Authorization头格式
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌格式错误",
				"error":   "INVALID_TOKEN_FORMAT",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// 使用新的验证方法，包含黑名单检查
		ctx := context.Background()
		tokenService := services.NewTokenService()
		claims, err := tokenService.ValidateTokenWithBlacklist(ctx, tokenString)
		if err != nil {
			// 根据错误类型返回不同的错误信息
			errorMessage := "认证失败"
			errorCode := "AUTH_FAILED"
			
			if strings.Contains(err.Error(), "已在其他设备登录") {
				errorMessage = err.Error()
				errorCode = "TOKEN_REVOKED"
			} else if strings.Contains(err.Error(), "已过期") {
				errorMessage = "令牌已过期，请重新登录"
				errorCode = "TOKEN_EXPIRED"
			} else if strings.Contains(err.Error(), "无效的令牌") {
				errorMessage = "无效的认证令牌"
				errorCode = "INVALID_TOKEN"
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": errorMessage,
				"error":   errorCode,
				"timestamp": time.Now().Unix(),
			})
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", claims.UserID)
		c.Set("uid", claims.Uid)
		c.Set("username", claims.Username)
		c.Set("claims", claims)
		c.Set("is_authenticated", true)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选的认证中间件（不强制要求登录）
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// 检查Bearer前缀
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := tokenParts[1]

		// 验证令牌
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("claims", claims)

		c.Next()
	}
}

// GetCurrentUser 获取当前用户ID
func GetCurrentUser(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint)
}

// GetCurrentUID 获取当前用户UID
func GetCurrentUID(c *gin.Context) string {
	uid, exists := c.Get("uid")
	if !exists {
		return ""
	}
	return uid.(string)
}

// GetCurrentUsername 获取当前用户名
func GetCurrentUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}

// IsAuthenticated 检查用户是否已认证
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// RequireAuth 要求认证的中间件（强制登录）
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsAuthenticated(c) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "需要登录才能访问此接口",
				"error":   "LOGIN_REQUIRED",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RegisterOpenMiddleware 检查注册开关
func RegisterOpenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过OPTIONS请求的注册开关检查
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		ctx := context.Background()
		exists, err := database.RedisClient.Exists(ctx, "dmin_system_isOpen").Result()
		if err != nil {
			utils.ErrorWithMessage(c, utils.CodeOperationFailed, "系统繁忙，请稍后再试")
			c.Abort()
			return
		}
		if exists > 0 {
			utils.ErrorWithMessage(c, utils.CodeRegisterClosed, "当前系统不允许注册")
			c.Abort()
			return
		}
		c.Next()
	}
}
