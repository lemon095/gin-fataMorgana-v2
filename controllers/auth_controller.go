package controllers

import (
	"gin-fataMorgana/middleware"
	"gin-fataMorgana/models"
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	userService *services.UserService
}

// NewAuthController 创建认证控制器实例
func NewAuthController() *AuthController {
	return &AuthController{
		userService: services.NewUserService(),
	}
}

// Register 用户注册
func (ac *AuthController) Register(c *gin.Context) {
	var req models.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	user, err := ac.userService.Register(&req)
	if err != nil {
		switch err.Error() {
		case "邮箱已被注册":
			utils.EmailAlreadyExists(c)
		case "该邮箱已被删除，无法重新注册":
			utils.ErrorWithMessage(c, utils.CodeUserAlreadyExists, err.Error())
		case "两次输入的密码不一致":
			utils.ErrorWithMessage(c, utils.CodeValidationFailed, err.Error())
		default:
			utils.ErrorWithMessage(c, utils.CodeOperationFailed, err.Error())
		}
		return
	}

	utils.SuccessWithMessage(c, "用户注册成功", gin.H{
		"user": user,
	})
}

// Login 用户登录
func (ac *AuthController) Login(c *gin.Context) {
	var req models.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取客户端IP地址
	clientIP := c.ClientIP()

	// 获取User-Agent
	userAgent := c.GetHeader("User-Agent")

	tokens, err := ac.userService.Login(&req, clientIP, userAgent)
	if err != nil {
		switch err.Error() {
		case "邮箱或密码错误":
			utils.LoginFailed(c)
		case "账户已被删除，无法登录":
			utils.ErrorWithMessage(c, utils.CodeUserNotFound, err.Error())
		case "账户已被禁用，无法登录":
			utils.AccountLocked(c)
		default:
			utils.ErrorWithMessage(c, utils.CodeOperationFailed, err.Error())
		}
		return
	}

	utils.SuccessWithMessage(c, "登录成功", gin.H{
		"tokens": tokens,
	})
}

// RefreshToken 刷新访问令牌
func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	tokens, err := ac.userService.RefreshToken(req.RefreshToken)
	if err != nil {
		switch err.Error() {
		case "无效的刷新令牌":
			utils.TokenInvalid(c)
		case "用户不存在":
			utils.UserNotFound(c)
		case "账户已被删除，无法刷新令牌":
			utils.ErrorWithMessage(c, utils.CodeUserNotFound, err.Error())
		case "账户已被禁用，无法刷新令牌":
			utils.AccountLocked(c)
		default:
			utils.ErrorWithMessage(c, utils.CodeOperationFailed, err.Error())
		}
		return
	}

	utils.SuccessWithMessage(c, "令牌刷新成功", gin.H{
		"tokens": tokens,
	})
}

// GetProfile 获取当前用户信息
func (ac *AuthController) GetProfile(c *gin.Context) {
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	user, err := ac.userService.GetUserByID(userID)
	if err != nil {
		switch err.Error() {
		case "用户不存在":
			utils.UserNotFound(c)
		case "用户已被删除":
			utils.ErrorWithMessage(c, utils.CodeUserNotFound, err.Error())
		default:
			utils.ErrorWithMessage(c, utils.CodeOperationFailed, err.Error())
		}
		return
	}

	utils.Success(c, gin.H{
		"user": user,
	})
}

// Logout 用户登出（客户端需要删除本地存储的token）
func (ac *AuthController) Logout(c *gin.Context) {
	// 这里只是返回成功消息，客户端需要删除本地token
	utils.SuccessWithMessage(c, "登出成功", nil)
}

// BindBankCard 绑定银行卡
func (ac *AuthController) BindBankCard(c *gin.Context) {
	var req services.BindBankCardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	user, err := ac.userService.BindBankCard(&req)
	if err != nil {
		switch err.Error() {
		case "用户不存在":
			utils.UserNotFound(c)
		case "账户已被删除，无法绑定银行卡":
			utils.ErrorWithMessage(c, utils.CodeUserNotFound, err.Error())
		case "账户已被禁用，无法绑定银行卡":
			utils.AccountLocked(c)
		case "银行名称不能为空", "持卡人姓名不能为空", "银行卡号不能为空", "卡类型不能为空":
			utils.ErrorWithMessage(c, utils.CodeValidationFailed, err.Error())
		case "银行卡号长度不正确", "卡类型不正确，支持的类型：借记卡、信用卡、储蓄卡":
			utils.ErrorWithMessage(c, utils.CodeValidationFailed, err.Error())
		default:
			utils.ErrorWithMessage(c, utils.CodeOperationFailed, err.Error())
		}
		return
	}

	utils.SuccessWithMessage(c, "银行卡绑定成功", gin.H{
		"user": user,
	})
}

// GetBankCardInfo 获取银行卡信息
func (ac *AuthController) GetBankCardInfo(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		utils.InvalidParamsWithMessage(c, "用户ID不能为空")
		return
	}

	bankCardInfo, err := ac.userService.GetBankCardInfo(uid)
	if err != nil {
		switch err.Error() {
		case "用户不存在":
			utils.UserNotFound(c)
		case "账户已被删除":
			utils.ErrorWithMessage(c, utils.CodeUserNotFound, err.Error())
		case "账户已被禁用":
			utils.AccountLocked(c)
		case "用户未绑定银行卡":
			utils.ErrorWithMessage(c, utils.CodeNotFound, err.Error())
		default:
			utils.ErrorWithMessage(c, utils.CodeOperationFailed, err.Error())
		}
		return
	}

	utils.Success(c, gin.H{
		"bank_card_info": bankCardInfo,
	})
}
