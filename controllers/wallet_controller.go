package controllers

import (
	"fmt"
	"gin-fataMorgana/middleware"
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"
	"strconv"

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
	// 获取参数
	uid := c.Query("uid")
	if uid == "" {
		utils.InvalidParamsWithMessage(c, "用户ID不能为空")
		return
	}

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 构建请求
	req := &services.GetUserTransactionsRequest{
		Uid:      uid,
		Page:     page,
		PageSize: pageSize,
	}

	// 调用服务
	response, err := wc.walletService.GetUserTransactions(req)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, response)
}

// GetWallet 获取钱包信息
func (wc *WalletController) GetWallet(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		utils.InvalidParamsWithMessage(c, "用户ID不能为空")
		return
	}

	wallet, err := wc.walletService.GetWallet(uid)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, wallet.ToResponse())
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

// Recharge 充值
func (wc *WalletController) Recharge(c *gin.Context) {
	var req struct {
		Uid         string  `json:"uid" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID作为操作员
	operatorUid := middleware.GetCurrentUser(c)
	operatorUidStr := "system"
	if operatorUid != 0 {
		operatorUidStr = strconv.FormatUint(uint64(operatorUid), 10)
	}

	err := wc.walletService.Recharge(req.Uid, req.Amount, req.Description, operatorUidStr)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "充值成功", nil)
}

// RequestWithdraw 申请提现
func (wc *WalletController) RequestWithdraw(c *gin.Context) {
	var req services.WithdrawRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID作为操作员
	operatorUid := middleware.GetCurrentUser(c)
	operatorUidStr := "system"
	if operatorUid != 0 {
		operatorUidStr = strconv.FormatUint(uint64(operatorUid), 10)
	}

	response, err := wc.walletService.RequestWithdraw(&req, operatorUidStr)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, response)
}

// ConfirmWithdraw 确认提现
func (wc *WalletController) ConfirmWithdraw(c *gin.Context) {
	var req struct {
		TransactionNo string `json:"transaction_no" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID作为操作员
	operatorUid := middleware.GetCurrentUser(c)
	operatorUidStr := "system"
	if operatorUid != 0 {
		operatorUidStr = strconv.FormatUint(uint64(operatorUid), 10)
	}

	err := wc.walletService.ConfirmWithdraw(req.TransactionNo, operatorUidStr)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "提现确认成功", nil)
}

// CancelWithdraw 取消提现
func (wc *WalletController) CancelWithdraw(c *gin.Context) {
	var req struct {
		TransactionNo string `json:"transaction_no" binding:"required"`
		Reason        string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID作为操作员
	operatorUid := middleware.GetCurrentUser(c)
	operatorUidStr := "system"
	if operatorUid != 0 {
		operatorUidStr = strconv.FormatUint(uint64(operatorUid), 10)
	}

	err := wc.walletService.CancelWithdraw(req.TransactionNo, operatorUidStr, req.Reason)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "提现取消成功", nil)
}

// GetWithdrawSummary 获取提现汇总信息
func (wc *WalletController) GetWithdrawSummary(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		utils.InvalidParamsWithMessage(c, "用户ID不能为空")
		return
	}

	summary, err := wc.walletService.GetWithdrawSummary(uid)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, summary)
}

// RechargeApply 充值申请
func (wc *WalletController) RechargeApply(c *gin.Context) {
	var req struct {
		Uid         string  `json:"uid" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}
	operatorUid := middleware.GetCurrentUser(c)
	operatorUidStr := "system"
	if operatorUid != 0 {
		operatorUidStr = fmt.Sprintf("%d", operatorUid)
	}
	transactionNo, err := wc.walletService.RechargeApply(req.Uid, req.Amount, req.Description, operatorUidStr)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeOperationFailed, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "充值申请已提交", gin.H{"transaction_no": transactionNo})
}

// RechargeConfirm 充值确认
func (wc *WalletController) RechargeConfirm(c *gin.Context) {
	var req struct {
		TransactionNo string `json:"transaction_no" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}
	operatorUid := middleware.GetCurrentUser(c)
	operatorUidStr := "system"
	if operatorUid != 0 {
		operatorUidStr = fmt.Sprintf("%d", operatorUid)
	}
	if err := wc.walletService.RechargeConfirm(req.TransactionNo, operatorUidStr); err != nil {
		utils.ErrorWithMessage(c, utils.CodeOperationFailed, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "充值已到账", nil)
}
