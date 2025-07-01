package models

import (
	"time"
)

// OrderStatus 订单状态枚举
const (
	OrderStatusPending   = "pending"   // 待处理
	OrderStatusSuccess   = "success"   // 成功
	OrderStatusFailed    = "failed"    // 失败
	OrderStatusCancelled = "cancelled" // 已取消
)

// Order 订单模型
type Order struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNo     string    `json:"order_no" gorm:"uniqueIndex;not null;size:32;comment:订单编号"`
	Uid         string    `json:"uid" gorm:"not null;size:8;index;comment:用户唯一ID"`
	BuyAmount   float64   `json:"buy_amount" gorm:"type:decimal(15,2);not null;comment:买入金额"`
	ProfitAmount float64  `json:"profit_amount" gorm:"type:decimal(15,2);default:0.00;comment:利润金额"`
	Status      string    `json:"status" gorm:"not null;size:20;default:'pending';index;comment:订单状态"`
	Description string    `json:"description" gorm:"size:200;comment:订单描述"`
	Remark      string    `json:"remark" gorm:"size:500;comment:备注信息"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;index;comment:创建时间"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// TableComment 表注释
func (Order) TableComment() string {
	return "订单表 - 记录用户订单信息，包括买入金额、利润金额等"
}

// OrderResponse 订单响应
type OrderResponse struct {
	ID           uint      `json:"id"`
	OrderNo      string    `json:"order_no"`
	Uid          string    `json:"uid"`
	BuyAmount    float64   `json:"buy_amount"`
	ProfitAmount float64   `json:"profit_amount"`
	Status       string    `json:"status"`
	StatusName   string    `json:"status_name"`
	Description  string    `json:"description"`
	Remark       string    `json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
}

// ToResponse 转换为响应格式
func (o *Order) ToResponse() OrderResponse {
	return OrderResponse{
		ID:           o.ID,
		OrderNo:      o.OrderNo,
		Uid:          o.Uid,
		BuyAmount:    o.BuyAmount,
		ProfitAmount: o.ProfitAmount,
		Status:       o.Status,
		StatusName:   o.GetStatusName(),
		Description:  o.Description,
		Remark:       o.Remark,
		CreatedAt:    o.CreatedAt,
	}
}

// GetStatusName 获取订单状态名称
func (o *Order) GetStatusName() string {
	statusNames := map[string]string{
		OrderStatusPending:   "待处理",
		OrderStatusSuccess:   "成功",
		OrderStatusFailed:    "失败",
		OrderStatusCancelled: "已取消",
	}
	return statusNames[o.Status]
}

// IsSuccess 检查订单是否成功
func (o *Order) IsSuccess() bool {
	return o.Status == OrderStatusSuccess
}

// IsPending 检查订单是否待处理
func (o *Order) IsPending() bool {
	return o.Status == OrderStatusPending
}

// IsFailed 检查订单是否失败
func (o *Order) IsFailed() bool {
	return o.Status == OrderStatusFailed
}

// IsCancelled 检查订单是否已取消
func (o *Order) IsCancelled() bool {
	return o.Status == OrderStatusCancelled
}

// GetTotalAmount 获取订单总金额（买入金额 + 利润金额）
func (o *Order) GetTotalAmount() float64 {
	return o.BuyAmount + o.ProfitAmount
}

// OrderListRequest 订单列表请求
type OrderListRequest struct {
	Page     int `json:"page" binding:"min=1"`      // 页码，从1开始
	PageSize int `json:"page_size" binding:"min=1,max=100"` // 每页大小，最大100
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Orders     []OrderResponse `json:"orders"`
	Pagination PaginationInfo  `json:"pagination"`
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