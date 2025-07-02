package services

import (
	"context"
	"errors"
	"fmt"
	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"
	"time"
)

// WalletService 钱包服务
type WalletService struct {
	walletRepo *database.WalletRepository
	userRepo   *database.UserRepository
}

// NewWalletService 创建钱包服务实例
func NewWalletService() *WalletService {
	return &WalletService{
		walletRepo: database.NewWalletRepository(),
		userRepo:   database.NewUserRepository(),
	}
}

// GetUserTransactionsRequest 获取用户资金记录请求
type GetUserTransactionsRequest struct {
	Uid      string `json:"uid" binding:"required"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// GetUserTransactionsResponse 获取用户资金记录响应
type GetUserTransactionsResponse struct {
	Transactions []models.WalletTransactionResponse `json:"transactions"`
	Pagination   PaginationInfo                     `json:"pagination"`
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

// GetUserTransactions 获取用户资金记录
func (s *WalletService) GetUserTransactions(req *GetUserTransactionsRequest) (*GetUserTransactionsResponse, error) {
	ctx := context.Background()

	// 设置默认分页参数
	page := req.Page
	pageSize := req.PageSize

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 验证分页参数
	if pageSize > 100 {
		return nil, errors.New("每页数量不能超过100")
	}

	// 获取交易记录
	transactions, total, err := s.walletRepo.GetUserTransactions(ctx, req.Uid, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取用户资金记录失败: %w", err)
	}

	// 转换为响应格式
	var transactionResponses []models.WalletTransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, transaction.ToResponse())
	}

	// 计算分页信息
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	pagination := PaginationInfo{
		CurrentPage: page,
		PageSize:    pageSize,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}

	return &GetUserTransactionsResponse{
		Transactions: transactionResponses,
		Pagination:   pagination,
	}, nil
}

// CreateWallet 创建钱包
func (s *WalletService) CreateWallet(uid string) (*models.Wallet, error) {
	ctx := context.Background()

	// 检查钱包是否已存在
	existingWallet, err := s.walletRepo.FindWalletByUid(ctx, uid)
	if err == nil && existingWallet != nil {
		return existingWallet, nil // 钱包已存在，直接返回
	}

	// 创建新钱包
	wallet := &models.Wallet{
		Uid:      uid,
		Balance:  0.00,
		Status:   1,
		Currency: "CNY",
	}

	if err := s.walletRepo.CreateWallet(ctx, wallet); err != nil {
		return nil, fmt.Errorf("创建钱包失败: %w", err)
	}

	return wallet, nil
}

// GetWallet 获取钱包信息
func (s *WalletService) GetWallet(uid string) (*models.Wallet, error) {
	ctx := context.Background()

	wallet, err := s.walletRepo.FindWalletByUid(ctx, uid)
	if err != nil {
		// 钱包不存在，自动创建
		wallet = &models.Wallet{
			Uid:      uid,
			Balance:  0.00,
			Status:   1,
			Currency: "CNY",
		}
		
		if err := s.walletRepo.CreateWallet(ctx, wallet); err != nil {
			return nil, fmt.Errorf("创建钱包失败: %w", err)
		}
	}

	return wallet, nil
}

// CreateTransaction 创建交易记录
func (s *WalletService) CreateTransaction(transaction *models.WalletTransaction) error {
	ctx := context.Background()

	// 生成交易流水号
	if transaction.TransactionNo == "" {
		transaction.TransactionNo = s.generateTransactionNo()
	}

	// 设置默认状态
	if transaction.Status == "" {
		transaction.Status = models.TransactionStatusSuccess
	}

	if err := s.walletRepo.CreateTransaction(ctx, transaction); err != nil {
		return fmt.Errorf("创建交易记录失败: %w", err)
	}

	return nil
}

// generateTransactionNo 生成交易流水号
func (s *WalletService) generateTransactionNo() string {
	// 格式：TX + 年月日 + 时分秒 + 4位随机数
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := utils.RandomString(4)
	return fmt.Sprintf("TX%s%s", timestamp, random)
}

// Recharge 充值
func (s *WalletService) Recharge(uid string, amount float64, description string, operatorUid string) (string, error) {
	ctx := context.Background()

	// 获取钱包，如果不存在则自动创建
	wallet, err := s.walletRepo.FindWalletByUid(ctx, uid)
	if err != nil {
		// 钱包不存在，自动创建
		wallet = &models.Wallet{
			Uid:      uid,
			Balance:  0.00,
			Status:   1,
			Currency: "CNY",
		}
		
		if err := s.walletRepo.CreateWallet(ctx, wallet); err != nil {
			return "", fmt.Errorf("创建钱包失败: %w", err)
		}
	}

	// 检查钱包状态
	if !wallet.IsActive() {
		return "", fmt.Errorf("钱包已被冻结，无法充值")
	}

	// 检查充值金额是否合理
	if amount <= 0 {
		return "", fmt.Errorf("充值金额必须大于0")
	}

	// 检查是否超过单笔充值限额（可选，这里设置100万）
	if amount > 1000000 {
		return "", fmt.Errorf("单笔充值金额不能超过100万元")
	}

	// 记录交易前余额
	balanceBefore := wallet.Balance

	// 生成交易流水号
	transactionNo := s.generateTransactionNo()

	// 创建充值交易记录（不增加余额）
	transaction := &models.WalletTransaction{
		TransactionNo: transactionNo,
		Uid:           uid,
		Type:          models.TransactionTypeRecharge,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceBefore, // 余额不变
		Status:        models.TransactionStatusPending, // 设置为待处理状态
		Description:   description,
		OperatorUid:   operatorUid,
	}

	// 创建交易记录
	if err := s.CreateTransaction(transaction); err != nil {
		return "", fmt.Errorf("创建充值申请失败: %w", err)
	}

	return transactionNo, nil
}

// WithdrawRequest 提现申请请求
type WithdrawRequest struct {
	Uid         string  `json:"uid" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
	BankCardInfo string `json:"bank_card_info"` // 提现银行卡信息
	Password    string  `json:"password" binding:"required"` // 登录密码
}

// WithdrawResponse 提现申请响应
type WithdrawResponse struct {
	TransactionNo string  `json:"transaction_no"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	Message       string  `json:"message"`
}

// RequestWithdraw 申请提现（锁定金额）
func (s *WalletService) RequestWithdraw(req *WithdrawRequest, operatorUid string) (*WithdrawResponse, error) {
	ctx := context.Background()

	// 验证用户密码
	user, err := s.userRepo.FindByUid(ctx, req.Uid)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 检查用户状态
	if !user.IsActive() {
		return nil, fmt.Errorf("用户账户已被禁用，无法提现")
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, fmt.Errorf("登录密码错误")
	}

	// 获取钱包，如果不存在则自动创建
	wallet, err := s.walletRepo.FindWalletByUid(ctx, req.Uid)
	if err != nil {
		// 钱包不存在，自动创建
		wallet = &models.Wallet{
			Uid:      req.Uid,
			Balance:  0.00,
			Status:   1,
			Currency: "CNY",
		}
		
		if err := s.walletRepo.CreateWallet(ctx, wallet); err != nil {
			return nil, fmt.Errorf("创建钱包失败: %w", err)
		}
	}

	// 检查钱包状态
	if !wallet.IsActive() {
		return nil, fmt.Errorf("钱包已被冻结，无法提现")
	}

	// 检查提现金额是否合理
	if req.Amount <= 0 {
		return nil, fmt.Errorf("提现金额必须大于0")
	}

	// 检查总余额是否足够
	if wallet.Balance < req.Amount {
		return nil, fmt.Errorf("总余额不足，当前总余额: %.2f，申请提现: %.2f", wallet.Balance, req.Amount)
	}

	// 检查是否超过单笔提现限额（可选，这里设置100万）
	if req.Amount > 1000000 {
		return nil, fmt.Errorf("单笔提现金额不能超过100万元")
	}

	// 检查是否超过每日提现限额（可选，这里设置500万）
	// 这里可以添加每日提现限额的检查逻辑
	// dailyWithdrawLimit := 5000000.0
	// if err := s.checkDailyWithdrawLimit(ctx, req.Uid, req.Amount, dailyWithdrawLimit); err != nil {
	//     return nil, err
	// }

	// 记录交易前余额
	balanceBefore := wallet.Balance

	// 直接扣减余额
	if err := wallet.Withdraw(req.Amount); err != nil {
		return nil, fmt.Errorf("扣减余额失败: %w", err)
	}

	// 更新钱包
	if err := s.walletRepo.UpdateWallet(ctx, wallet); err != nil {
		return nil, fmt.Errorf("更新钱包失败: %w", err)
	}

	// 生成交易流水号
	transactionNo := s.generateTransactionNo()

	// 创建提现交易记录
	transaction := &models.WalletTransaction{
		TransactionNo: transactionNo,
		Uid:           req.Uid,
		Type:          models.TransactionTypeWithdraw,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  wallet.Balance,
		Status:        models.TransactionStatusPending, // 设置为待处理状态
		Description:   req.Description,
		Remark:        req.BankCardInfo, // 将银行卡信息存储到备注字段
		OperatorUid:   operatorUid,
	}

	// 创建交易记录
	if err := s.CreateTransaction(transaction); err != nil {
		// 如果创建交易记录失败，需要回滚扣减的余额
		wallet.Recharge(req.Amount)
		s.walletRepo.UpdateWallet(ctx, wallet)
		return nil, fmt.Errorf("创建提现申请失败: %w", err)
	}

	return &WithdrawResponse{
		TransactionNo: transactionNo,
		Amount:        req.Amount,
		Status:        models.TransactionStatusPending,
		Message:       "提现申请已提交，等待处理",
	}, nil
}

// checkDailyWithdrawLimit 检查每日提现限额
func (s *WalletService) checkDailyWithdrawLimit(ctx context.Context, uid string, amount float64, dailyLimit float64) error {
	// 获取今日已申请的提现总额
	today := time.Now().Format("2006-01-02")
	
	// 查询今日的提现申请记录
	transactions, _, err := s.walletRepo.GetTransactionsByDateRange(ctx, uid, today, today, 1, 1000)
	if err != nil {
		return fmt.Errorf("查询今日提现记录失败: %w", err)
	}

	// 计算今日已申请的提现总额
	var todayTotal float64
	for _, tx := range transactions {
		if tx.Type == models.TransactionTypeWithdraw && 
		   (tx.Status == models.TransactionStatusPending || tx.Status == models.TransactionStatusSuccess) {
			todayTotal += tx.Amount
		}
	}

	// 检查是否超过每日限额
	if todayTotal+amount > dailyLimit {
		return fmt.Errorf("超过每日提现限额，今日已申请: %.2f，本次申请: %.2f，每日限额: %.2f", 
			todayTotal, amount, dailyLimit)
	}

	return nil
}

// GetWithdrawSummary 获取提现汇总信息
func (s *WalletService) GetWithdrawSummary(uid string) (map[string]interface{}, error) {
	ctx := context.Background()

	// 获取钱包信息
	wallet, err := s.walletRepo.FindWalletByUid(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("获取钱包失败: %w", err)
	}

	// 获取今日提现申请
	today := time.Now().Format("2006-01-02")
	todayTransactions, _, err := s.walletRepo.GetTransactionsByDateRange(ctx, uid, today, today, 1, 1000)
	if err != nil {
		return nil, fmt.Errorf("查询今日提现记录失败: %w", err)
	}

	// 计算今日提现统计
	var todayPendingTotal float64
	var todaySuccessTotal float64
	var todayCancelledTotal float64
	var pendingCount int
	var successCount int
	var cancelledCount int

	for _, tx := range todayTransactions {
		if tx.Type == models.TransactionTypeWithdraw {
			switch tx.Status {
			case models.TransactionStatusPending:
				todayPendingTotal += tx.Amount
				pendingCount++
			case models.TransactionStatusSuccess:
				todaySuccessTotal += tx.Amount
				successCount++
			case models.TransactionStatusCancelled:
				todayCancelledTotal += tx.Amount
				cancelledCount++
			}
		}
	}

	// 获取所有待处理的提现申请
	pendingTransactions, _, err := s.walletRepo.GetTransactionsByType(ctx, uid, models.TransactionTypeWithdraw, 1, 1000)
	if err != nil {
		return nil, fmt.Errorf("查询待处理提现记录失败: %w", err)
	}

	// 计算待处理提现统计
	var totalPendingAmount float64
	for _, tx := range pendingTransactions {
		if tx.Status == models.TransactionStatusPending {
			totalPendingAmount += tx.Amount
		}
	}

	return map[string]interface{}{
		"wallet_info": map[string]interface{}{
			"total_balance":     wallet.Balance,
			"available_balance": wallet.GetAvailableBalance(),
		},
		"today_withdraw": map[string]interface{}{
			"pending_total":   todayPendingTotal,
			"success_total":   todaySuccessTotal,
			"cancelled_total": todayCancelledTotal,
			"pending_count":   pendingCount,
			"success_count":   successCount,
			"cancelled_count": cancelledCount,
		},
		"total_pending": map[string]interface{}{
			"amount": totalPendingAmount,
			"count":  len(pendingTransactions),
		},
		"limits": map[string]interface{}{
			"single_limit":  1000000.0, // 单笔限额
			"daily_limit":   5000000.0, // 每日限额
			"remaining_today": 5000000.0 - todayPendingTotal - todaySuccessTotal,
		},
	}, nil
}

// GetTransactionDetailRequest 获取交易详情请求
type GetTransactionDetailRequest struct {
	TransactionNo string `json:"transaction_no" binding:"required"`
}

// GetTransactionDetail 获取交易详情
func (s *WalletService) GetTransactionDetail(req *GetTransactionDetailRequest) (*models.WalletTransactionResponse, error) {
	ctx := context.Background()

	// 获取交易记录
	transaction, err := s.walletRepo.GetTransactionByNo(ctx, req.TransactionNo)
	if err != nil {
		return nil, fmt.Errorf("获取交易详情失败: %w", err)
	}

	// 转换为响应格式
	response := transaction.ToResponse()

	return &response, nil
}
