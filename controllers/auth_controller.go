package controllers

import (
	"context"
	"encoding/json"
	"gin-fataMorgana/database"
	"gin-fataMorgana/middleware"
	"gin-fataMorgana/models"
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	userService             *services.UserService
	operationFailureService *services.OperationFailureService
}

// NewAuthController 创建认证控制器实例
func NewAuthController() *AuthController {
	return &AuthController{
		userService:             services.NewUserService(),
		operationFailureService: services.NewOperationFailureService(),
	}
}

// Register godoc
// @Summary 用户注册
// @Description 用户注册接口，需要提供邮箱、密码和邀请码
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.UserRegisterRequest true "注册请求参数"
// @Success 200 {object} utils.Response{data=models.UserResponse} "注册成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 422 {object} utils.Response "验证失败"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	// 读取原始请求体
	body, _ := ioutil.ReadAll(c.Request.Body)

	c.Request.Body = ioutil.NopCloser(strings.NewReader(string(body)))

	var req models.UserRegisterRequest

	// 解析JSON请求
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError(c, "JSON解析失败: %v", err)
		utils.HandleValidationError(c, err)
		return
	}

	// 如果结构体有BankCardInfo字段且为空，赋默认值
	type bankCardInfoSetter interface {
		SetBankCardInfoDefault()
	}
	if setter, ok := any(&req).(bankCardInfoSetter); ok {
		setter.SetBankCardInfoDefault()
	}

	user, err := ac.userService.Register(&req)
	if err != nil {
		utils.LogError(c, "注册失败: 账号=%s, 邀请码=%s, 错误原因=%s",
			req.Account, req.InviteCode, err.Error())

		// 记录操作失败
		ac.operationFailureService.RecordRegisterFailure(c.Request.Context(), req, gin.H{
			"error": err.Error(),
			"code":  utils.CodeOperationFailed,
		})

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

	// 解析JSON请求
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError(c, "登录JSON解析失败: %v", err)
		utils.HandleValidationError(c, err)
		return
	}

	// 获取客户端IP地址
	clientIP := c.ClientIP()
	// 获取User-Agent
	userAgent := c.GetHeader("User-Agent")

	tokens, err := ac.userService.Login(&req, clientIP, userAgent)
	if err != nil {
		utils.LogError(c, "登录失败: 账号=%s, 密码=%s, 错误原因=%s",
			req.Account, req.Password, err.Error())

		// 记录操作失败
		ac.operationFailureService.RecordLoginFailure(c.Request.Context(), nil, req, gin.H{
			"error": err.Error(),
			"code":  utils.CodeOperationFailed,
		})

		switch err.Error() {
		case "邮箱或密码错误":
			utils.LoginFailed(c)
		case "账户已被删除，无法登录":
			utils.ErrorWithMessage(c, utils.CodeUserNotFound, err.Error())
		case "账户已被禁用，无法登录":
			utils.AccountLocked(c)
		case "账户待审核，无法登录":
			utils.ErrorWithMessage(c, utils.CodeUserPendingApproval, err.Error())
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
		utils.HandleValidationError(c, err)
		return
	}

	tokens, err := ac.userService.RefreshToken(&req)
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
	var req models.GetProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

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

// Logout 用户登出（撤销当前token）
func (ac *AuthController) Logout(c *gin.Context) {
	// 获取当前用户信息
	uid := middleware.GetCurrentUID(c)
	if uid == "" {
		utils.Unauthorized(c)
		return
	}

	// 获取当前token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.SuccessWithMessage(c, "登出成功", nil)
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		utils.SuccessWithMessage(c, "登出成功", nil)
		return
	}

	tokenString := tokenParts[1]

	// 撤销用户会话
	ctx := context.Background()
	tokenService := services.NewTokenService()

	// 将当前token加入黑名单
	err := tokenService.AddTokenToBlacklist(ctx, tokenString)
	if err != nil {
		// 记录错误但不影响登出流程
		utils.LogError(c, "将token加入黑名单失败: %v", err)
	}

	// 撤销用户会话
	err = tokenService.RevokeUserSession(ctx, uid)
	if err != nil {
		// 记录错误但不影响登出流程
		utils.LogError(c, "撤销用户会话失败: %v", err)
	}

	utils.SuccessWithMessage(c, "登出成功", nil)
}

// BindBankCard 绑定银行卡
func (ac *AuthController) BindBankCard(c *gin.Context) {
	var req services.BindBankCardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	// 根据user_id查询uid，确保只能操作自己的账户
	userRepo := database.NewUserRepository()
	var user models.User
	err := userRepo.FindByID(context.Background(), userID, &user)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	userResponse, err := ac.userService.BindBankCard(&req, user.Uid)
	if err != nil {
		// 记录操作失败
		ac.operationFailureService.RecordBankCardBindFailure(c.Request.Context(), user.Uid, req, gin.H{
			"error": err.Error(),
			"code":  utils.CodeOperationFailed,
		})

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
		"user": userResponse,
	})
}

// GetBankCardInfo 获取银行卡信息
func (ac *AuthController) GetBankCardInfo(c *gin.Context) {
	var req models.GetBankCardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

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

	// 解析银行卡信息
	var bankCardInfo models.BankCardInfo
	if user.BankCardInfo != "" {
		if err := json.Unmarshal([]byte(user.BankCardInfo), &bankCardInfo); err != nil {
			utils.ErrorWithMessage(c, utils.CodeOperationFailed, "银行卡信息解析失败")
			return
		}

		// 对银行卡号进行脱敏处理
		if bankCardInfo.CardNumber != "" {
			bankCardInfo.CardNumber = utils.MaskBankCard(bankCardInfo.CardNumber)
		}

		// 对持卡人姓名进行脱敏处理
		if bankCardInfo.CardHolder != "" {
			bankCardInfo.CardHolder = utils.MaskName(bankCardInfo.CardHolder)
		}
	}

	utils.Success(c, gin.H{
		"bank_card_info": bankCardInfo,
	})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 用户修改登录密码
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param request body models.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/auth/change-password [post]
func (ac *AuthController) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	// 获取当前用户UID
	uid := middleware.GetCurrentUID(c)
	if uid == "" {
		utils.Unauthorized(c)
		return
	}

	// 调用服务层修改密码
	err := ac.userService.ChangePassword(&req, uid)
	if err != nil {
		switch err.Error() {
		case "用户不存在":
			utils.UserNotFound(c)
		case "用户已被删除，无法修改密码":
			utils.ErrorWithMessage(c, utils.CodeUserNotFound, err.Error())
		case "账户已被禁用，无法修改密码":
			utils.AccountLocked(c)
		case "当前密码错误":
			utils.ErrorWithMessage(c, utils.CodeValidationFailed, err.Error())
		case "新密码不能与当前密码相同":
			utils.ErrorWithMessage(c, utils.CodeValidationFailed, err.Error())
		case "新密码不能为空", "新密码长度不能少于6位", "新密码长度不能超过50位":
			utils.ErrorWithMessage(c, utils.CodeValidationFailed, err.Error())
		default:
			utils.ErrorWithMessage(c, utils.CodeOperationFailed, err.Error())
		}
		return
	}

	// 返回成功响应
	utils.SuccessWithMessage(c, "密码修改成功", nil)
}
