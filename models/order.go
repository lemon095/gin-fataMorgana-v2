package models

import (
	"errors"
	"time"
)

// OrderStatus 订单状态枚举
const (
	OrderStatusPending   = "pending"   // 待处理
	OrderStatusSuccess   = "success"   // 成功
	OrderStatusFailed    = "failed"    // 失败
	OrderStatusCancelled = "cancelled" // 已取消
	OrderStatusExpired   = "expired"   // 已过期
)

// TaskStatus 任务完成状态枚举
const (
	TaskStatusPending = "pending" // 待完成
	TaskStatusSuccess = "success" // 已完成
)

// Order 订单表
type Order struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNo        string    `json:"order_no" gorm:"uniqueIndex;not null;size:32;comment:订单编号"`
	Uid            string    `json:"uid" gorm:"not null;size:8;index;comment:用户唯一ID"`
	PeriodNumber   string    `json:"period_number" gorm:"not null;size:32;comment:期号"`
	Amount         float64   `json:"amount" gorm:"type:decimal(15,2);not null;comment:订单金额"`
	ProfitAmount   float64   `json:"profit_amount" gorm:"type:decimal(15,2);not null;comment:利润金额"`
	Status         string    `json:"status" gorm:"not null;size:20;default:'pending';index;comment:订单状态"`
	ExpireTime     time.Time `json:"expire_time" gorm:"not null;index;comment:订单剩余时间"`
	LikeCount      int       `json:"like_count" gorm:"not null;default:0;comment:点赞数"`
	ShareCount     int       `json:"share_count" gorm:"not null;default:0;comment:转发数"`
	FollowCount    int       `json:"follow_count" gorm:"not null;default:0;comment:关注数"`
	FavoriteCount  int       `json:"favorite_count" gorm:"not null;default:0;comment:收藏数"`
	LikeStatus     string    `json:"like_status" gorm:"not null;size:20;default:'pending';comment:点赞完成状态"`
	ShareStatus    string    `json:"share_status" gorm:"not null;size:20;default:'pending';comment:转发完成状态"`
	FollowStatus   string    `json:"follow_status" gorm:"not null;size:20;default:'pending';comment:关注完成状态"`
	FavoriteStatus string    `json:"favorite_status" gorm:"not null;size:20;default:'pending';comment:收藏完成状态"`
	AuditorUid     string    `json:"auditor_uid" gorm:"size:8;index;comment:审核员ID"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime;index;comment:创建时间"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// TableComment 表注释
func (Order) TableComment() string {
	return "订单表 - 记录用户创建的订单信息，包括任务要求和完成状态"
}

// OrderResponse 订单响应
type OrderResponse struct {
	ID                 uint      `json:"id"`
	OrderNo            string    `json:"order_no"`
	Uid                string    `json:"uid"`
	Number             string    `json:"period_number"`
	Amount             float64   `json:"amount"`
	ProfitAmount       float64   `json:"profit_amount"`
	Status             string    `json:"status"`
	StatusName         string    `json:"status_name"`
	ExpireTime         time.Time `json:"expire_time"`
	LikeCount          int       `json:"like_count"`
	ShareCount         int       `json:"share_count"`
	FollowCount        int       `json:"follow_count"`
	FavoriteCount      int       `json:"favorite_count"`
	LikeStatus         string    `json:"like_status"`
	LikeStatusName     string    `json:"like_status_name"`
	ShareStatus        string    `json:"share_status"`
	ShareStatusName    string    `json:"share_status_name"`
	FollowStatus       string    `json:"follow_status"`
	FollowStatusName   string    `json:"follow_status_name"`
	FavoriteStatus     string    `json:"favorite_status"`
	FavoriteStatusName string    `json:"favorite_status_name"`
	AuditorUid         string    `json:"auditor_uid"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	IsExpired          bool      `json:"is_expired"`
	RemainingTime      int64     `json:"remaining_time"` // 剩余时间（秒）
}

// ToResponse 转换为响应格式
func (o *Order) ToResponse() OrderResponse {
	return OrderResponse{
		ID:                 o.ID,
		OrderNo:            o.OrderNo,
		Uid:                o.Uid,
		Number:             o.PeriodNumber,
		Amount:             o.Amount,
		ProfitAmount:       o.ProfitAmount,
		Status:             o.Status,
		StatusName:         o.GetStatusName(),
		ExpireTime:         o.ExpireTime,
		LikeCount:          o.LikeCount,
		ShareCount:         o.ShareCount,
		FollowCount:        o.FollowCount,
		FavoriteCount:      o.FavoriteCount,
		LikeStatus:         o.LikeStatus,
		LikeStatusName:     o.GetTaskStatusName(o.LikeStatus),
		ShareStatus:        o.ShareStatus,
		ShareStatusName:    o.GetTaskStatusName(o.ShareStatus),
		FollowStatus:       o.FollowStatus,
		FollowStatusName:   o.GetTaskStatusName(o.FollowStatus),
		FavoriteStatus:     o.FavoriteStatus,
		FavoriteStatusName: o.GetTaskStatusName(o.FavoriteStatus),
		AuditorUid:         o.AuditorUid,
		CreatedAt:          o.CreatedAt,
		UpdatedAt:          o.UpdatedAt,
		IsExpired:          o.IsExpired(),
		RemainingTime:      o.GetRemainingTime(),
	}
}

// GetStatusName 获取订单状态名称
func (o *Order) GetStatusName() string {
	statusNames := map[string]string{
		OrderStatusPending:   "待处理",
		OrderStatusSuccess:   "成功",
		OrderStatusFailed:    "失败",
		OrderStatusCancelled: "已取消",
		OrderStatusExpired:   "已过期",
	}
	return statusNames[o.Status]
}

// GetTaskStatusName 获取任务状态名称
func (o *Order) GetTaskStatusName(status string) string {
	statusNames := map[string]string{
		TaskStatusPending: "待完成",
		TaskStatusSuccess: "已完成",
	}
	return statusNames[status]
}

// IsExpired 检查订单是否已过期
func (o *Order) IsExpired() bool {
	return time.Now().UTC().After(o.ExpireTime)
}

// GetRemainingTime 获取剩余时间（秒）
func (o *Order) GetRemainingTime() int64 {
	if o.IsExpired() {
		return 0
	}
	return int64(time.Until(o.ExpireTime).Seconds())
}

// IsAllTasksCompleted 检查所有任务是否已完成
func (o *Order) IsAllTasksCompleted() bool {
	return o.LikeStatus == TaskStatusSuccess &&
		o.ShareStatus == TaskStatusSuccess &&
		o.FollowStatus == TaskStatusSuccess &&
		o.FavoriteStatus == TaskStatusSuccess
}

// IsAllTasksZero 检查所有任务数是否都为0
func (o *Order) IsAllTasksZero() bool {
	return o.LikeCount == 0 && o.ShareCount == 0 && o.FollowCount == 0 && o.FavoriteCount == 0
}

// HasAnyTask 检查是否有任何任务
func (o *Order) HasAnyTask() bool {
	return o.LikeCount > 0 || o.ShareCount > 0 || o.FollowCount > 0 || o.FavoriteCount > 0
}

// InitializeTaskStatuses 初始化任务状态
func (o *Order) InitializeTaskStatuses() {
	// 如果点赞数为0，设置为已完成
	if o.LikeCount == 0 {
		o.LikeStatus = TaskStatusSuccess
	} else {
		o.LikeStatus = TaskStatusPending
	}

	// 如果转发数为0，设置为已完成
	if o.ShareCount == 0 {
		o.ShareStatus = TaskStatusSuccess
	} else {
		o.ShareStatus = TaskStatusPending
	}

	// 如果关注数为0，设置为已完成
	if o.FollowCount == 0 {
		o.FollowStatus = TaskStatusSuccess
	} else {
		o.FollowStatus = TaskStatusPending
	}

	// 如果收藏数为0，设置为已完成
	if o.FavoriteCount == 0 {
		o.FavoriteStatus = TaskStatusSuccess
	} else {
		o.FavoriteStatus = TaskStatusPending
	}
}

// ValidateOrderData 验证订单数据
func (o *Order) ValidateOrderData() error {
	// 检查总金额不能为0
	if o.Amount <= 0 {
		return errors.New("订单金额必须大于0")
	}

	// 检查利润金额不能为负数
	if o.ProfitAmount < 0 {
		return errors.New("利润金额不能为负数")
	}

	// 检查不能所有任务都为0
	if o.IsAllTasksZero() {
		return errors.New("至少需要有一个任务数量大于0")
	}

	// 检查任务数量不能为负数
	if o.LikeCount < 0 || o.ShareCount < 0 || o.FollowCount < 0 || o.FavoriteCount < 0 {
		return errors.New("任务数量不能为负数")
	}

	return nil
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	ProfitAmount  float64 `json:"profit_amount" binding:"required,gte=0"`
	LikeCount     int     `json:"like_count" binding:"gte=0"`
	ShareCount    int     `json:"share_count" binding:"gte=0"`
	FollowCount   int     `json:"follow_count" binding:"gte=0"`
	FavoriteCount int     `json:"favorite_count" binding:"gte=0"`
}

// OrderStatusType 订单状态类型枚举
const (
	OrderStatusTypeInProgress = 1 // 进行中
	OrderStatusTypeCompleted  = 2 // 已完成
	OrderStatusTypeAll        = 3 // 全部
)

// GetOrderListRequest 获取订单列表请求
type GetOrderListRequest struct {
	Page     int `json:"page" binding:"min=1"`
	PageSize int `json:"page_size" binding:"min=1"`
	Status   int `json:"status" binding:"min=1,max=3"` // 1:进行中 2:已完成 3:全部
}

// GetStatusByType 根据状态类型获取对应的状态值
func GetStatusByType(statusType int) string {
	switch statusType {
	case OrderStatusTypeInProgress:
		return OrderStatusPending
	case OrderStatusTypeCompleted:
		return OrderStatusSuccess
	case OrderStatusTypeAll:
		return "" // 空字符串表示查询期数数据
	default:
		return ""
	}
}

// GetStatusTypeName 获取状态类型名称
func GetStatusTypeName(statusType int) string {
	switch statusType {
	case OrderStatusTypeInProgress:
		return "进行中"
	case OrderStatusTypeCompleted:
		return "已完成"
	case OrderStatusTypeAll:
		return "期数数据"
	default:
		return "未知"
	}
}

// GetOrderDetailRequest 获取订单详情请求
type GetOrderDetailRequest struct {
	OrderNo string `json:"order_no" binding:"required"`
}

// OrderListRequest 订单列表请求
type OrderListRequest struct {
	Page     int `json:"page" binding:"min=1"`              // 页码，从1开始
	PageSize int `json:"page_size" binding:"min=1"`         // 每页大小，最小1
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
