package models

import (
	"time"
)

// LotteryPeriodStatus 期数状态枚举
const (
	LotteryPeriodStatusPending = "pending" // 待开始
	LotteryPeriodStatusActive  = "active"  // 进行中
	LotteryPeriodStatusClosed  = "closed"  // 已结束
)

// LotteryPeriod 游戏期数表
// 对应数据库表 lottery_periods
// 仅保留SQL定义的字段
type LotteryPeriod struct {
	ID               uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PeriodNumber     string    `json:"period_number" gorm:"uniqueIndex:uk_period_number;not null;size:20;comment:期数编号"`
	TotalOrderAmount float64   `json:"total_order_amount" gorm:"type:decimal(15,2);not null;default:0.00;comment:本期购买订单金额"`
	Status           string    `json:"status" gorm:"not null;size:20;default:'pending';index:idx_status;comment:期数状态: pending-待开始, active-进行中, closed-已结束"`
	LotteryResult    *string   `json:"lottery_result" gorm:"size:50;comment:开奖结果"`
	OrderStartTime   time.Time `json:"order_start_time" gorm:"not null;index:idx_order_start_time;comment:订单开始时间"`
	OrderEndTime     time.Time `json:"order_end_time" gorm:"not null;index:idx_order_end_time;comment:订单结束时间"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime;index:idx_created_at;comment:创建时间"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (LotteryPeriod) TableName() string {
	return "lottery_periods"
}

// TableComment 表注释
func (LotteryPeriod) TableComment() string {
	return "游戏期数表 - 记录每期的编号、订单金额、状态、开奖结果和时间信息"
}

// IsActive 检查期数是否活跃（在开始时间和结束时间范围内）
func (lp *LotteryPeriod) IsActive() bool {
	now := time.Now()
	return now.After(lp.OrderStartTime) && now.Before(lp.OrderEndTime)
}

// IsExpired 检查期数是否已过期
func (lp *LotteryPeriod) IsExpired() bool {
	return time.Now().After(lp.OrderEndTime)
}

// IsPending 检查期数是否待开始
func (lp *LotteryPeriod) IsPending() bool {
	return time.Now().Before(lp.OrderStartTime)
}

// IsValidTimeRange 检查期数时间范围是否有效
func (lp *LotteryPeriod) IsValidTimeRange() bool {
	return lp.OrderStartTime.Before(lp.OrderEndTime)
}

// GetTimeRangeError 获取时间范围错误信息
func (lp *LotteryPeriod) GetTimeRangeError() string {
	if !lp.IsValidTimeRange() {
		return "期数开始时间不能晚于结束时间"
	}
	return ""
}

// GetStatus 获取期数状态
func (lp *LotteryPeriod) GetStatus() string {
	if lp.IsPending() {
		return LotteryPeriodStatusPending
	} else if lp.IsActive() {
		return LotteryPeriodStatusActive
	} else {
		return LotteryPeriodStatusClosed
	}
}

// GetStatusName 获取状态名称
func (lp *LotteryPeriod) GetStatusName() string {
	statusNames := map[string]string{
		LotteryPeriodStatusPending: "待开始",
		LotteryPeriodStatusActive:  "进行中",
		LotteryPeriodStatusClosed:  "已结束",
	}
	return statusNames[lp.Status]
}

// ToOrderResponse 转换为订单响应格式
func (lp *LotteryPeriod) ToOrderResponse() OrderResponse {
	return OrderResponse{
		ID:           lp.ID,
		OrderNo:      lp.PeriodNumber,
		Uid:          "", // 期数没有特定用户
		Amount:       lp.TotalOrderAmount,
		ProfitAmount: 0, // 期数没有利润金额
		Status:       lp.GetStatus(),
		StatusName:   lp.GetStatusName(),
		ExpireTime:   lp.OrderEndTime,
		CreatedAt:    lp.CreatedAt,
		UpdatedAt:    lp.UpdatedAt,
		IsExpired:    lp.IsExpired(),
		RemainingTime: func() int64 {
			if lp.IsExpired() {
				return 0
			}
			return int64(time.Until(lp.OrderEndTime).Seconds())
		}(),
	}
}
