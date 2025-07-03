package database

import (
	"context"
	"gin-fataMorgana/models"
	"time"
)

// LotteryPeriodRepository 期数仓库
type LotteryPeriodRepository struct {
	*BaseRepository
}

// NewLotteryPeriodRepository 创建期数仓库实例
func NewLotteryPeriodRepository() *LotteryPeriodRepository {
	return &LotteryPeriodRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// GetCurrentPeriod 获取当前活跃期数
// 查询条件：当前时间在订单开始时间和订单结束时间范围内
func (r *LotteryPeriodRepository) GetCurrentPeriod(ctx context.Context) (*models.LotteryPeriod, error) {
	var period models.LotteryPeriod

	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("order_start_time <= ? AND order_end_time > ?", now, now).
		Order("created_at DESC").
		First(&period).Error

	if err != nil {
		return nil, err
	}

	return &period, nil
}

// GetPeriodByNumber 根据期数编号获取期数
func (r *LotteryPeriodRepository) GetPeriodByNumber(ctx context.Context, periodNumber string) (*models.LotteryPeriod, error) {
	var period models.LotteryPeriod

	err := r.db.WithContext(ctx).
		Where("period_number = ?", periodNumber).
		First(&period).Error

	if err != nil {
		return nil, err
	}

	return &period, nil
} 