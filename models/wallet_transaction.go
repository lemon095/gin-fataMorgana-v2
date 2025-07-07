package models

import (
	"fmt"
	"time"
)

// TransactionType 交易类型枚举
const (
	TransactionTypeRecharge = "recharge"  // 充值
	TransactionTypeWithdraw = "withdraw"  // 提现
	TransactionTypeOrderBuy = "purchase"  // 购买
	TransactionTypeGroupBuy = "group_buy" // 拼单
	TransactionTypeProfit   = "profit"    // 利润
)

// TransactionStatus 交易状态枚举
const (
	TransactionStatusPending   = "pending"   // 待处理
	TransactionStatusSuccess   = "success"   // 成功
	TransactionStatusFailed    = "failed"    // 失败
	TransactionStatusCancelled = "cancelled" // 已取消
)

// WalletTransaction 钱包交易流水表
type WalletTransaction struct {
	ID             uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	TransactionNo  string  `json:"transaction_no" gorm:"uniqueIndex;not null;size:32;comment:交易流水号"`
	Uid            string  `json:"uid" gorm:"not null;size:8;index;comment:用户唯一ID"`
	Type           string  `json:"type" gorm:"not null;size:20;index;comment:交易类型"`
	Amount         float64 `json:"amount" gorm:"type:decimal(15,2);not null;comment:交易金额"`
	BalanceBefore  float64 `json:"balance_before" gorm:"type:decimal(15,2);not null;comment:交易前余额"`
	BalanceAfter   float64 `json:"balance_after" gorm:"type:decimal(15,2);not null;comment:交易后余额"`
	Status         string  `json:"status" gorm:"not null;size:20;default:'success';index;comment:交易状态"`
	Description    string  `json:"description" gorm:"size:200;comment:交易描述"`
	Remark         string  `json:"remark" gorm:"size:500;comment:备注信息"`
	RelatedOrderNo string  `json:"related_order_no" gorm:"size:32;index;comment:关联订单号"`

	OperatorUid string    `json:"operator_uid" gorm:"size:8;index;comment:操作员ID"`
	IPAddress   string    `json:"ip_address" gorm:"size:45;comment:操作IP地址"`
	UserAgent   string    `json:"user_agent" gorm:"size:500;comment:用户代理"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;index;comment:创建时间"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}

// TableComment 表注释
func (WalletTransaction) TableComment() string {
	return "钱包交易流水表 - 记录所有钱包交易明细，包括充值、提现、购买、拼单等操作"
}

// WalletTransactionResponse 交易流水响应
type WalletTransactionResponse struct {
	ID             uint    `json:"id"`
	TransactionNo  string  `json:"transaction_no"`
	Uid            string  `json:"uid"`
	Type           string  `json:"type"`
	TypeName       string  `json:"type_name"`
	Amount         float64 `json:"amount"`
	BalanceBefore  float64 `json:"balance_before"`
	BalanceAfter   float64 `json:"balance_after"`
	Status         string  `json:"status"`
	StatusName     string  `json:"status_name"`
	Description    string  `json:"description"`
	Remark         string  `json:"remark"`
	RelatedOrderNo string  `json:"related_order_no"`

	OperatorUid string    `json:"operator_uid"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToResponse 转换为响应格式
func (t *WalletTransaction) ToResponse() WalletTransactionResponse {
	return WalletTransactionResponse{
		ID:             t.ID,
		TransactionNo:  t.TransactionNo,
		Uid:            t.Uid,
		Type:           t.Type,
		TypeName:       t.GetTypeName(),
		Amount:         t.Amount,
		BalanceBefore:  t.BalanceBefore,
		BalanceAfter:   t.BalanceAfter,
		Status:         t.Status,
		StatusName:     t.GetStatusName(),
		Description:    t.Description,
		Remark:         t.Remark,
		RelatedOrderNo: t.RelatedOrderNo,

		OperatorUid: t.OperatorUid,
		IPAddress:   t.IPAddress,
		UserAgent:   t.UserAgent,
		CreatedAt:   t.CreatedAt,
	}
}

// GetTypeName 获取交易类型名称
func (t *WalletTransaction) GetTypeName() string {
	typeNames := map[string]string{
		TransactionTypeRecharge: "充值",
		TransactionTypeWithdraw: "提现",
		TransactionTypeOrderBuy: "购买订单",
		TransactionTypeGroupBuy: "拼单",
		TransactionTypeProfit:   "利润",
	}
	return typeNames[t.Type]
}

// GetStatusName 获取交易状态名称
func (t *WalletTransaction) GetStatusName() string {
	statusNames := map[string]string{
		TransactionStatusPending:   "待处理",
		TransactionStatusSuccess:   "成功",
		TransactionStatusFailed:    "失败",
		TransactionStatusCancelled: "已取消",
	}
	return statusNames[t.Status]
}

// IsSuccess 检查交易是否成功
func (t *WalletTransaction) IsSuccess() bool {
	return t.Status == TransactionStatusSuccess
}

// IsPending 检查交易是否待处理
func (t *WalletTransaction) IsPending() bool {
	return t.Status == TransactionStatusPending
}

// IsFailed 检查交易是否失败
func (t *WalletTransaction) IsFailed() bool {
	return t.Status == TransactionStatusFailed
}

// IsCancelled 检查交易是否已取消
func (t *WalletTransaction) IsCancelled() bool {
	return t.Status == TransactionStatusCancelled
}

// GetAmountDisplay 获取金额显示（带正负号）
func (t *WalletTransaction) GetAmountDisplay() string {
	switch t.Type {
	case TransactionTypeRecharge, TransactionTypeProfit:
		return "+" + formatAmount(t.Amount)
	case TransactionTypeWithdraw, TransactionTypeOrderBuy, TransactionTypeGroupBuy:
		return "-" + formatAmount(t.Amount)
	default:
		return formatAmount(t.Amount)
	}
}

// formatAmount 格式化金额
func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

// 交易类型说明：
//
// 1. recharge (充值) - 用户从银行卡充值到钱包
// 2. withdraw (提现) - 用户从钱包提现到银行卡
// 3. purchase (购买) - 用户购买订单
// 4. group_buy (拼单) - 用户参与拼单
// 5. profit (利润) - 用户获得利润收入
//
// 交易状态说明：
//
// 1. pending (待处理) - 交易已创建，等待处理
// 2. success (成功) - 交易处理成功
// 3. failed (失败) - 交易处理失败
// 4. cancelled (已取消) - 交易已取消

// GetTransactionsResponse 获取交易记录响应
type GetTransactionsResponse struct {
	Transactions []WalletTransaction `json:"transactions"`
	Total        int64               `json:"total"`
	Page         int                 `json:"page"`
	PageSize     int                 `json:"page_size"`
}

// WithdrawResponse 提现响应
type WithdrawResponse struct {
	TransactionNo string  `json:"transaction_no"`
	Amount        float64 `json:"amount"`
	Balance       float64 `json:"balance"`
	Status        string  `json:"status"`
}

// TransactionDetail 交易详情
type TransactionDetail struct {
	TransactionNo string    `json:"transaction_no"`
	Uid           string    `json:"uid"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	Balance       float64   `json:"balance"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// WithdrawSummary 提现汇总
type WithdrawSummary struct {
	TotalWithdrawAmount float64 `json:"total_withdraw_amount"`
	TotalWithdrawCount  int64   `json:"total_withdraw_count"`
	PendingAmount       float64 `json:"pending_amount"`
	PendingCount        int64   `json:"pending_count"`
	SuccessAmount       float64 `json:"success_amount"`
	SuccessCount        int64   `json:"success_count"`
	FailedAmount        float64 `json:"failed_amount"`
	FailedCount         int64   `json:"failed_count"`
}
