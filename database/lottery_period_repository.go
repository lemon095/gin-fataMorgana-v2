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
func (r *LotteryPeriodRepository) GetCurrentPeriod(ctx context.Context) (*models.LotteryPeriod, error) {
	var period models.LotteryPeriod

	// 按 status = 'active' 查询当前活跃期数，获取最新的一条
	err := r.db.WithContext(ctx).
		Where("status = ?", "active").
		Order("created_at DESC").
		First(&period).Error

	if err != nil {
		// 如果没有活跃期数，尝试获取最近的期数
		err = r.db.WithContext(ctx).
			Order("created_at DESC").
			First(&period).Error
		
		if err != nil {
			return nil, err
		}
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

// UpdatePeriodStatus 更新期数状态
func (r *LotteryPeriodRepository) UpdatePeriodStatus(ctx context.Context) error {
	now := time.Now().UTC()
	
	// 更新已过期的期数为 closed
	err := r.db.WithContext(ctx).
		Model(&models.LotteryPeriod{}).
		Where("order_end_time <= ? AND status != ?", now, "closed").
		Update("status", "closed").Error
	if err != nil {
		return err
	}
	
	// 更新正在进行的期数为 active
	err = r.db.WithContext(ctx).
		Model(&models.LotteryPeriod{}).
		Where("order_start_time <= ? AND order_end_time > ? AND status != ?", now, now, "active").
		Update("status", "active").Error
	if err != nil {
		return err
	}
	
	// 更新待开始的期数为 pending
	err = r.db.WithContext(ctx).
		Model(&models.LotteryPeriod{}).
		Where("order_start_time > ? AND status != ?", now, "pending").
		Update("status", "pending").Error
	
	return err
} 