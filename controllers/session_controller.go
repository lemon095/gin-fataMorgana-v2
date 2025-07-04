package controllers

import (
	"time"

	"gin-fataMorgana/middleware"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// SessionController 会话控制器
type SessionController struct{}

// NewSessionController 创建会话控制器实例
func NewSessionController() *SessionController {
	return &SessionController{}
}

// CheckLoginStatus 检查登录状态
func (sc *SessionController) CheckLoginStatus(c *gin.Context) {
	var req models.GetSessionStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	status := middleware.GetLoginStatus(c)

	utils.SuccessWithMessage(c, "获取登录状态成功", status)
}

// GetCurrentUserInfo 获取当前用户信息
func (sc *SessionController) GetCurrentUserInfo(c *gin.Context) {
	var req models.GetSessionUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	if !middleware.IsAuthenticated(c) {
		utils.Unauthorized(c)
		return
	}

	userID := middleware.GetCurrentUser(c)
	username := middleware.GetCurrentUsername(c)

	utils.SuccessWithMessage(c, "获取用户信息成功", gin.H{
		"user_id":    userID,
		"username":   username,
		"login_time": time.Now().Unix(),
	})
}

// Logout 用户登出
func (sc *SessionController) Logout(c *gin.Context) {
	// 在实际项目中，这里可以将token加入黑名单
	// 或者更新用户的最后登出时间等

	utils.SuccessWithMessage(c, "登出成功", gin.H{
		"logout_time": time.Now().Unix(),
	})
}

// RefreshSession 刷新会话
func (sc *SessionController) RefreshSession(c *gin.Context) {
	if !middleware.IsAuthenticated(c) {
		utils.Unauthorized(c)
		return
	}

	// 在实际项目中，这里可以更新用户的最后活动时间
	// 或者延长会话有效期等

	utils.SuccessWithMessage(c, "会话刷新成功", gin.H{
		"refresh_time": time.Now().Unix(),
		"user_id":      middleware.GetCurrentUser(c),
		"username":     middleware.GetCurrentUsername(c),
	})
}
