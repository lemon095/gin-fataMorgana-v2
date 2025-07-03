package models

import (
	"time"
)

// PurchaseConfig 价格配置
type PurchaseConfig struct {
	LikeAmount     float64   `json:"like_amount"`
	ShareAmount    float64   `json:"share_amount"`
	ForwardAmount  float64   `json:"forward_amount"`
	FavoriteAmount float64   `json:"favorite_amount"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// PeriodListResponse 期数列表响应
type PeriodListResponse struct {
	ID             uint    `json:"id"`
	PeriodNumber   string  `json:"period_number"`
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	Status         string  `json:"status"`
	IsExpired      bool    `json:"is_expired"`
	RemainingTime  int64   `json:"remaining_time"`
	LikeAmount     float64 `json:"like_amount"`
	ShareAmount    float64 `json:"share_amount"`
	ForwardAmount  float64 `json:"forward_amount"`
	FavoriteAmount float64 `json:"favorite_amount"`
}
