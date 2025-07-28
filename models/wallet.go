package models

import (
	"time"

	"github.com/gin-gonic/gin"
)

// 钱包状态常量
const (
	WalletStatusNormal     = 1 // 正常
	WalletStatusFrozen     = 0 // 冻结
	WalletStatusNoWithdraw = 2 // 无法提现
)

// Wallet 钱包模型
type Wallet struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Uid          string    `gorm:"uniqueIndex;not null;size:8;comment:用户唯一ID" json:"uid"`                // 用户ID
	Balance      float64   `gorm:"type:decimal(15,2);default:0.00;not null;comment:钱包余额" json:"balance"` // 总余额
	Status       int       `gorm:"default:1;comment:钱包状态 1:正常 0:冻结 2:无法提现" json:"status"`                // 状态：1-正常，0-冻结，2-无法提现
	Currency     string    `gorm:"default:'PHP';size:3;comment:货币类型" json:"currency"`                    // 货币类型
	LastActiveAt time.Time `gorm:"autoUpdateTime;comment:最后活跃时间" json:"last_active_at"`                  // 最后活跃时间
	CreatedAt    time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (Wallet) TableName() string {
	return "wallets"
}

// TableComment 表注释
func (Wallet) TableComment() string {
	return "钱包表 - 存储用户钱包信息，包括余额等"
}

// ToResponse 转换为响应格式
func (w *Wallet) ToResponse() gin.H {
	return gin.H{
		"id":             w.ID,
		"uid":            w.Uid,
		"balance":        w.Balance,
		"status":         w.Status,
		"currency":       w.Currency,
		"last_active_at": w.LastActiveAt,
		"created_at":     w.CreatedAt,
		"updated_at":     w.UpdatedAt,
	}
}

// GetAvailableBalance 获取可用余额
func (w *Wallet) GetAvailableBalance() float64 {
	return w.Balance // 没有冻结余额，可用余额等于总余额
}

// Recharge 充值（不统计收入）
func (w *Wallet) Recharge(amount float64) {
	w.Balance += amount
	// 充值不算收入，只是资金转移
}

// Withdraw 提现（不统计支出）
func (w *Wallet) Withdraw(amount float64) error {
	if w.Balance < amount {
		return ErrInsufficientBalance
	}
	w.Balance -= amount
	// 提现不算支出，只是资金转移
	return nil
}

// IsActive 检查钱包是否激活
func (w *Wallet) IsActive() bool {
	return w.Status == WalletStatusNormal
}

// IsFrozen 检查钱包是否冻结
func (w *Wallet) IsFrozen() bool {
	return w.Status == WalletStatusFrozen
}

// IsNoWithdraw 检查钱包是否无法提现
func (w *Wallet) IsNoWithdraw() bool {
	return w.Status == WalletStatusNoWithdraw
}

// CanWithdraw 检查钱包是否可以提现
func (w *Wallet) CanWithdraw() bool {
	return w.Status != WalletStatusNoWithdraw && w.Status != WalletStatusFrozen
}

// CanOperate 检查钱包是否可以操作（充值、消费等）
func (w *Wallet) CanOperate() bool {
	return w.Status != WalletStatusFrozen
}

// GetStatusName 获取状态名称
func (w *Wallet) GetStatusName() string {
	statusNames := map[int]string{
		WalletStatusNormal:     "正常",
		WalletStatusFrozen:     "冻结",
		WalletStatusNoWithdraw: "无法提现",
	}
	return statusNames[w.Status]
}

// UpdateLastActive 更新最后活跃时间
func (w *Wallet) UpdateLastActive() {
	w.LastActiveAt = time.Now()
}

// 错误定义
var (
	ErrInsufficientBalance = &WalletError{Message: "余额不足"}
)

// WalletError 钱包错误
type WalletError struct {
	Message string
}

func (e *WalletError) Error() string {
	return e.Message
}
