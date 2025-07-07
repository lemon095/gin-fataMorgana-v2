package controllers

import (
	"context"
	"gin-fataMorgana/database"
	"gin-fataMorgana/middleware"
	"gin-fataMorgana/models"
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// WalletController 钱包控制器
type WalletController struct {
	walletService *services.WalletService
}

// NewWalletController 创建钱包控制器实例
func NewWalletController() *WalletController {
	return &WalletController{
		walletService: services.NewWalletService(),
	}
}

// GetUserTransactions 获取用户资金记录
func (wc *WalletController) GetUserTransactions(c *gin.Context) {
	var req models.GetTransactionsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, utils.CodeInvalidParams, "参数验证失败")
		return
	}

	// 获取当前用户ID
	userID, err := middleware.GetUserID(c)
	if err != nil {
		middleware.ErrorResponse(c, utils.CodeAuth, "用户未认证")
		return
	}

	// 根据user_id查询uid，确保获取正确的uid
	userRepo := database.NewUserRepository()
	var user models.User
	err = userRepo.FindByID(context.Background(), uint(userID), &user)
	if err != nil {
		middleware.ErrorResponse(c, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	// 构建服务请求
	serviceReq := &services.GetUserTransactionsRequest{
		Uid:      user.Uid,
		Page:     req.Page,
		PageSize: req.PageSize,
		Type:     req.Type,
	}

	// 调用服务
	response, err := wc.walletService.GetUserTransactions(serviceReq)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			middleware.ErrorResponse(c, appErr.Code, appErr.Message)
		} else {
			middleware.ErrorResponse(c, utils.CodeDatabaseError, "获取交易记录失败")
		}
		return
	}

	middleware.SuccessResponse(c, response)
}

// GetWallet 获取钱包信息
func (wc *WalletController) GetWallet(c *gin.Context) {
	var req models.GetWalletRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, utils.CodeInvalidParams, "参数验证失败")
		return
	}

	// 获取当前用户ID
	userID, err := middleware.GetUserID(c)
	if err != nil {
		middleware.ErrorResponse(c, utils.CodeAuth, "用户未认证")
		return
	}

	// 根据user_id查询uid，确保获取正确的uid
	userRepo := database.NewUserRepository()
	var user models.User
	err = userRepo.FindByID(context.Background(), uint(userID), &user)
	if err != nil {
		middleware.ErrorResponse(c, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	wallet, err := wc.walletService.GetWallet(user.Uid)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			middleware.ErrorResponse(c, appErr.Code, appErr.Message)
		} else {
			middleware.ErrorResponse(c, utils.CodeDatabaseError, "获取钱包信息失败")
		}
		return
	}

	middleware.SuccessResponse(c, wallet.ToResponse())
}

// CreateWallet 创建钱包
func (wc *WalletController) CreateWallet(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		utils.InvalidParamsWithMessage(c, "用户ID不能为空")
		return
	}

	wallet, err := wc.walletService.CreateWallet(uid)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, wallet.ToResponse())
}

// Recharge 充值申请
func (wc *WalletController) Recharge(c *gin.Context) {
	var req struct {
		Uid         string  `json:"uid" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}

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

	// 根据user_id查询uid，确保只能操作自己的钱包
	userRepo := database.NewUserRepository()
	var user models.User
	err := userRepo.FindByID(context.Background(), userID, &user)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	// 校验uid是否与当前登录用户匹配
	if req.Uid != user.Uid {
		utils.ErrorWithMessage(c, utils.CodeForbidden, "只能操作自己的钱包")
		return
	}

	transactionNo, err := wc.walletService.Recharge(req.Uid, req.Amount, req.Description)
	if err != nil {
		// 检查是否是AppError类型
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorWithMessage(c, appErr.Code, appErr.Message)
		} else {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		}
		return
	}

	utils.SuccessWithMessage(c, "充值申请已提交", gin.H{"transaction_no": transactionNo})
}

// AddProfit 添加利润
func (wc *WalletController) AddProfit(c *gin.Context) {
	var req struct {
		Uid         string  `json:"uid" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
		OrderNo     string  `json:"order_no"` // 关联订单号（可选）
	}

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

	// 根据user_id查询uid，确保只能操作自己的钱包
	userRepo := database.NewUserRepository()
	var user models.User
	err := userRepo.FindByID(context.Background(), userID, &user)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	// 校验uid是否与当前登录用户匹配
	if req.Uid != user.Uid {
		utils.ErrorWithMessage(c, utils.CodeForbidden, "只能操作自己的钱包")
		return
	}

	transactionNo, err := wc.walletService.CreateProfitTransaction(context.Background(), req.Uid, req.Amount, req.Description, req.OrderNo)
	if err != nil {
		// 检查是否是AppError类型
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorWithMessage(c, appErr.Code, appErr.Message)
		} else {
			utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		}
		return
	}

	utils.SuccessWithMessage(c, "利润添加成功", gin.H{"transaction_no": transactionNo})
}

// RequestWithdraw 申请提现
func (wc *WalletController) RequestWithdraw(c *gin.Context) {
	var req services.WithdrawRequest

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

	// 根据user_id查询uid，确保只能操作自己的钱包
	userRepo := database.NewUserRepository()
	var user models.User
	err := userRepo.FindByID(context.Background(), userID, &user)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	response, err := wc.walletService.RequestWithdraw(&req, user.Uid)
	if err != nil {
		// 检查是否是AppError类型
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorWithMessage(c, appErr.Code, appErr.Message)
		} else {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		}
		return
	}

	utils.Success(c, response)
}

// GetWithdrawSummary 获取提现汇总信息
func (wc *WalletController) GetWithdrawSummary(c *gin.Context) {
	var req models.GetWithdrawSummaryRequest

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

	// 根据user_id查询uid，确保获取正确的uid
	userRepo := database.NewUserRepository()
	var user models.User
	err := userRepo.FindByID(context.Background(), userID, &user)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	summary, err := wc.walletService.GetWithdrawSummary(user.Uid)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, summary)
}

// GetTransactionDetail 获取交易详情
func (wc *WalletController) GetTransactionDetail(c *gin.Context) {
	var req struct {
		TransactionNo string `json:"transaction_no" binding:"required"`
	}

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

	// 根据user_id查询uid，确保只能查看自己的交易
	userRepo := database.NewUserRepository()
	var user models.User
	err := userRepo.FindByID(context.Background(), userID, &user)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	// 先获取交易记录，检查是否属于当前用户
	walletRepo := database.NewWalletRepository()
	transaction, err := walletRepo.GetTransactionByNo(context.Background(), req.TransactionNo)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeTransactionDetailGetFailed, "获取交易详情失败")
		return
	}

	// 验证交易是否属于当前用户
	if transaction.Uid != user.Uid {
		utils.ErrorWithMessage(c, utils.CodeForbidden, "只能查看自己的交易详情")
		return
	}

	// 构建服务请求
	serviceReq := &services.GetTransactionDetailRequest{
		TransactionNo: req.TransactionNo,
	}

	// 调用服务
	response, err := wc.walletService.GetTransactionDetail(serviceReq)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, response)
}
