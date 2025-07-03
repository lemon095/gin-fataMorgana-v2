package models

import (
	"time"
)

// GroupBuyStatus 拼单状态枚举
const (
	GroupBuyStatusNotStarted = "not_started" // 未开启
	GroupBuyStatusPending    = "pending"     // 进行中
	GroupBuyStatusSuccess    = "success"     // 已完成
)

// GroupBuyType 拼单类型枚举
const (
	GroupBuyTypeNormal = "normal" // 普通拼单
	GroupBuyTypeFlash  = "flash"  // 限时拼单
	GroupBuyTypeVIP    = "vip"    // VIP拼单
)

// GroupBuy 拼单表
type GroupBuy struct {
	ID                  uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupBuyNo          string    `json:"group_buy_no" gorm:"uniqueIndex;not null;size:32;comment:拼单编号"`
	OrderNo             *string   `json:"order_no" gorm:"size:32;index;comment:关联订单编号"`
	CreatorUid          string    `json:"creator_uid" gorm:"not null;size:8;index;comment:创建用户ID"`
	Uid                 string    `json:"uid" gorm:"size:8;index;comment:参与用户ID"`
	CurrentParticipants int       `json:"current_participants" gorm:"not null;default:1;comment:当前参与人数"`
	TargetParticipants  int       `json:"target_participants" gorm:"not null;default:2;comment:目标参与人数"`
	GroupBuyType        string    `json:"group_buy_type" gorm:"not null;size:20;default:'normal';index;comment:拼单类型"`
	TotalAmount         float64   `json:"total_amount" gorm:"type:decimal(15,2);not null;comment:拼单总金额"`
	PaidAmount          float64   `json:"paid_amount" gorm:"type:decimal(15,2);not null;default:0;comment:已付款金额"`
	PerPersonAmount     float64   `json:"per_person_amount" gorm:"type:decimal(15,2);not null;comment:每人需要付款金额"`
	Deadline            time.Time `json:"deadline" gorm:"not null;index;comment:拼单截止时间"`
	Status              string    `json:"status" gorm:"not null;size:20;default:'not_started';index;comment:拼单状态"`
	Description         string    `json:"description" gorm:"type:text;comment:拼单描述"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime;index;comment:创建时间"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (GroupBuy) TableName() string {
	return "group_buys"
}

// TableComment 表注释
func (GroupBuy) TableComment() string {
	return "拼单表 - 记录拼单信息，包括参与人数、付款金额、截止时间等"
}

// GetGroupBuyDetailResponse 获取拼单详情响应
type GetGroupBuyDetailResponse struct {
	HasData             bool      `json:"has_data"`             // 是否有数据
	GroupBuyNo          string    `json:"group_buy_no"`         // 拼单编号
	GroupBuyType        string    `json:"group_buy_type"`       // 拼单类型
	TotalAmount         float64   `json:"total_amount"`         // 拼单总金额
	CurrentParticipants int       `json:"current_participants"` // 当前参与人数
	TargetParticipants  int       `json:"target_participants"`  // 最大参与人数
	PaidAmount          float64   `json:"paid_amount"`          // 已付款金额
	PerPersonAmount     float64   `json:"per_person_amount"`    // 每人需要付款金额
	RemainingAmount     float64   `json:"remaining_amount"`     // 还需要付款的金额
	Deadline            time.Time `json:"deadline"`             // 截止时间
}

// GetGroupBuyListResponse 获取拼单列表响应
type GetGroupBuyListResponse struct {
	HasData    bool                        `json:"has_data"`    // 是否有数据
	GroupBuys  []GetGroupBuyDetailResponse `json:"group_buys"`  // 拼单列表
	TotalCount int                         `json:"total_count"` // 总数量
}

// ToDetailResponse 转换为详情响应格式
func (g *GroupBuy) ToDetailResponse() GetGroupBuyDetailResponse {
	// 计算还需要付款的金额
	remainingAmount := g.TotalAmount - g.PaidAmount
	if remainingAmount < 0 {
		remainingAmount = 0
	}

	return GetGroupBuyDetailResponse{
		HasData:             true,
		GroupBuyNo:          g.GroupBuyNo,
		GroupBuyType:        g.GroupBuyType,
		TotalAmount:         g.TotalAmount,
		CurrentParticipants: g.CurrentParticipants,
		TargetParticipants:  g.TargetParticipants,
		PaidAmount:          g.PaidAmount,
		PerPersonAmount:     g.PerPersonAmount,
		RemainingAmount:     remainingAmount,
		Deadline:            g.Deadline,
	}
}

// JoinGroupBuyRequest 确认参与拼单请求
type JoinGroupBuyRequest struct {
	GroupBuyNo string `json:"group_buy_no" binding:"required"` // 拼单编号
}

// JoinGroupBuyResponse 确认参与拼单响应
type JoinGroupBuyResponse struct {
	OrderID uint `json:"order_id"` // 订单ID
}
