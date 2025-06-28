package models

import (
	"time"
)

// Wallet 钱包模型
type Wallet struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Uid           string    `json:"uid" gorm:"uniqueIndex;not null;size:8;comment:用户唯一ID"`
	Balance       float64   `json:"balance" gorm:"type:decimal(15,2);default:0.00;not null;comment:钱包余额"`
	FrozenBalance float64   `json:"frozen_balance" gorm:"type:decimal(15,2);default:0.00;not null;comment:冻结余额"`
	TotalIncome   float64   `json:"total_income" gorm:"type:decimal(15,2);default:0.00;not null;comment:总收入"`
	TotalExpense  float64   `json:"total_expense" gorm:"type:decimal(15,2);default:0.00;not null;comment:总支出"`
	Status        int       `json:"status" gorm:"default:1;comment:钱包状态 1:正常 0:冻结"`
	Currency      string    `json:"currency" gorm:"default:'CNY';size:3;comment:货币类型"`
	LastActiveAt  time.Time `json:"last_active_at" gorm:"autoUpdateTime;comment:最后活跃时间"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (Wallet) TableName() string {
	return "wallets"
}

// WalletResponse 钱包响应
type WalletResponse struct {
	ID               uint      `json:"id"`
	Uid              string    `json:"uid"`
	Balance          float64   `json:"balance"`
	FrozenBalance    float64   `json:"frozen_balance"`
	AvailableBalance float64   `json:"available_balance"`
	TotalIncome      float64   `json:"total_income"`
	TotalExpense     float64   `json:"total_expense"`
	Status           int       `json:"status"`
	Currency         string    `json:"currency"`
	LastActiveAt     time.Time `json:"last_active_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (w *Wallet) ToResponse() WalletResponse {
	return WalletResponse{
		ID:               w.ID,
		Uid:              w.Uid,
		Balance:          w.Balance,
		FrozenBalance:    w.FrozenBalance,
		AvailableBalance: w.GetAvailableBalance(),
		TotalIncome:      w.TotalIncome,
		TotalExpense:     w.TotalExpense,
		Status:           w.Status,
		Currency:         w.Currency,
		LastActiveAt:     w.LastActiveAt,
		CreatedAt:        w.CreatedAt,
		UpdatedAt:        w.UpdatedAt,
	}
}

// GetAvailableBalance 获取可用余额
func (w *Wallet) GetAvailableBalance() float64 {
	return w.Balance - w.FrozenBalance
}

// AddBalance 增加余额（通用方法）
func (w *Wallet) AddBalance(amount float64) {
	w.Balance += amount
}

// Recharge 充值（不统计收入）
func (w *Wallet) Recharge(amount float64) {
	w.Balance += amount
	// 充值不算收入，只是资金转移
}

// AddIncome 增加收入（统计收入）
func (w *Wallet) AddIncome(amount float64) {
	w.Balance += amount
	w.TotalIncome += amount
}

// SubtractBalance 减少余额（通用方法）
func (w *Wallet) SubtractBalance(amount float64) error {
	availableBalance := w.GetAvailableBalance()
	if availableBalance < amount {
		return ErrInsufficientBalance
	}
	w.Balance -= amount
	return nil
}

// Withdraw 提现（不统计支出）
func (w *Wallet) Withdraw(amount float64) error {
	availableBalance := w.GetAvailableBalance()
	if availableBalance < amount {
		return ErrInsufficientBalance
	}
	w.Balance -= amount
	// 提现不算支出，只是资金转移
	return nil
}

// AddExpense 增加支出（统计支出）
func (w *Wallet) AddExpense(amount float64) error {
	availableBalance := w.GetAvailableBalance()
	if availableBalance < amount {
		return ErrInsufficientBalance
	}
	w.Balance -= amount
	w.TotalExpense += amount
	return nil
}

// FreezeBalance 冻结余额
func (w *Wallet) FreezeBalance(amount float64) error {
	availableBalance := w.GetAvailableBalance()
	if availableBalance < amount {
		return ErrInsufficientBalance
	}
	w.FrozenBalance += amount
	return nil
}

// UnfreezeBalance 解冻余额
func (w *Wallet) UnfreezeBalance(amount float64) error {
	if w.FrozenBalance < amount {
		return ErrInsufficientFrozenBalance
	}
	w.FrozenBalance -= amount
	return nil
}

// HasSufficientBalance 检查可用余额是否充足
func (w *Wallet) HasSufficientBalance(amount float64) bool {
	return w.GetAvailableBalance() >= amount
}

// GetBalance 获取总余额
func (w *Wallet) GetBalance() float64 {
	return w.Balance
}

// GetFrozenBalance 获取冻结余额
func (w *Wallet) GetFrozenBalance() float64 {
	return w.FrozenBalance
}

// IsActive 检查钱包是否激活
func (w *Wallet) IsActive() bool {
	return w.Status == 1
}

// Freeze 冻结钱包
func (w *Wallet) Freeze() {
	w.Status = 0
}

// Unfreeze 解冻钱包
func (w *Wallet) Unfreeze() {
	w.Status = 1
}

// 错误定义
var (
	ErrInsufficientBalance       = &WalletError{Message: "余额不足"}
	ErrInsufficientFrozenBalance = &WalletError{Message: "冻结余额不足"}
	ErrWalletFrozen              = &WalletError{Message: "钱包已被冻结"}
)

// WalletError 钱包错误
type WalletError struct {
	Message string
}

func (e *WalletError) Error() string {
	return e.Message
}

// 业务场景说明：
//
// 1. 充值 (Recharge) - 用户从银行卡充值到钱包
//    - 增加余额
//    - 不统计收入（因为只是资金转移）
//    - 示例：用户充值100元，余额+100，总收入不变
//
// 2. 收入 (AddIncome) - 用户获得真正的收入
//    - 增加余额
//    - 统计收入（如：返现、奖励、收益等）
//    - 示例：用户获得返现10元，余额+10，总收入+10
//
// 3. 提现 (Withdraw) - 用户从钱包提现到银行卡
//    - 减少余额
//    - 不统计支出（因为只是资金转移）
//    - 示例：用户提现50元，余额-50，总支出不变
//
// 4. 支出 (AddExpense) - 用户进行消费
//    - 减少余额
//    - 统计支出（如：购物、服务费等）
//    - 示例：用户购物消费30元，余额-30，总支出+30
//
// 5. 冻结 (FreezeBalance) - 冻结部分余额
//    - 冻结余额增加
//    - 可用余额减少
//    - 示例：冻结20元用于订单，可用余额-20，冻结余额+20
//
// 使用示例：
// wallet := &Wallet{Uid: "12345678", Balance: 100.00}
//
// // 充值
// wallet.Recharge(50.00)           // 余额: 150.00, 总收入: 0.00
//
// // 获得返现
// wallet.AddIncome(10.00)          // 余额: 160.00, 总收入: 10.00
//
// // 购物消费
// wallet.AddExpense(30.00)         // 余额: 130.00, 总支出: 30.00
//
// // 提现
// wallet.Withdraw(20.00)           // 余额: 110.00, 总支出: 30.00 (不变)
//
// // 冻结余额
// wallet.FreezeBalance(15.00)      // 可用余额: 95.00, 冻结余额: 15.00
