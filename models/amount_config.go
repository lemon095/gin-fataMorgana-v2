package models

import (
	"time"
)

// AmountConfig 金额配置
type AmountConfig struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Type        string    `json:"type" gorm:"not null;size:20;index;comment:配置类型: recharge-充值, withdraw-提现"`
	Amount      float64   `json:"amount" gorm:"not null;type:decimal(10,2);comment:金额"`
	Description string    `json:"description" gorm:"size:100;comment:描述"`
	IsActive    bool      `json:"is_active" gorm:"not null;default:1;comment:是否激活"`
	SortOrder   int       `json:"sort_order" gorm:"not null;default:0;comment:排序"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (AmountConfig) TableName() string {
	return "amount_config"
}

// AmountConfigRequest 金额配置请求
type AmountConfigRequest struct {
	Type string `json:"type" binding:"required,oneof=recharge withdraw"` // 配置类型
}

// AmountConfigResponse 金额配置响应
type AmountConfigResponse struct {
	ID          int64   `json:"id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	IsActive    bool    `json:"is_active"`
	SortOrder   int     `json:"sort_order"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (ac *AmountConfig) ToResponse() *AmountConfigResponse {
	return &AmountConfigResponse{
		ID:          ac.ID,
		Type:        ac.Type,
		Amount:      ac.Amount,
		Description: ac.Description,
		IsActive:    ac.IsActive,
		SortOrder:   ac.SortOrder,
		CreatedAt:   ac.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   ac.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
} 