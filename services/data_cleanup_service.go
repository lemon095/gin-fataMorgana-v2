package services

import (
	"context"
	"time"

	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
)

// DataCleanupService 数据清理服务
type DataCleanupService struct {
	orderRepo    *database.OrderRepository
	groupBuyRepo *database.GroupBuyRepository
	config       *DataCleanupConfig
}

// DataCleanupConfig 数据清理配置
type DataCleanupConfig struct {
	RetentionDays int `yaml:"retention_days"` // 保留天数
}

// CleanupStats 清理统计
type CleanupStats struct {
	DeletedOrders    int64         `json:"deleted_orders"`
	DeletedGroupBuys int64         `json:"deleted_group_buys"`
	LastCleanup      time.Time     `json:"last_cleanup"`
	CleanupTime      time.Duration `json:"cleanup_time"`
}

// NewDataCleanupService 创建新的数据清理服务
func NewDataCleanupService(config *DataCleanupConfig) *DataCleanupService {
	return &DataCleanupService{
		orderRepo:    database.NewOrderRepository(),
		groupBuyRepo: database.NewGroupBuyRepository(),
		config:       config,
	}
}

// CleanupOldSystemOrders 清理旧的系统订单
func (s *DataCleanupService) CleanupOldSystemOrders() (*CleanupStats, error) {
	startTime := time.Now()
	ctx := context.Background()

	// 计算清理时间点
	cutoffTime := time.Now().AddDate(0, 0, -s.config.RetentionDays)

	// 清理系统订单
	deletedOrders, err := s.cleanupSystemOrders(ctx, cutoffTime)
	if err != nil {
		return nil, err
	}

	// 清理系统拼单
	deletedGroupBuys, err := s.cleanupSystemGroupBuys(ctx, cutoffTime)
	if err != nil {
		return nil, err
	}

	duration := time.Since(startTime)

	stats := &CleanupStats{
		DeletedOrders:    deletedOrders,
		DeletedGroupBuys: deletedGroupBuys,
		LastCleanup:      time.Now(),
		CleanupTime:      duration,
	}

	return stats, nil
}

// cleanupSystemOrders 清理系统订单
func (s *DataCleanupService) cleanupSystemOrders(ctx context.Context, cutoffTime time.Time) (int64, error) {
	var totalDeleted int64
	batchSize := 1000

	for {
		// 分批删除，避免锁表
		result := database.DB.WithContext(ctx).
			Where("is_system_order = ? AND created_at < ?", true, cutoffTime).
			Limit(batchSize).
			Delete(&models.Order{})

		if result.Error != nil {
			return totalDeleted, result.Error
		}

		deleted := result.RowsAffected
		totalDeleted += deleted

		// 如果没有更多数据需要删除，退出循环
		if deleted < int64(batchSize) {
			break
		}

		// 短暂休息，避免对数据库造成压力
		time.Sleep(100 * time.Millisecond)
	}

	return totalDeleted, nil
}

// cleanupSystemGroupBuys 清理系统拼单
func (s *DataCleanupService) cleanupSystemGroupBuys(ctx context.Context, cutoffTime time.Time) (int64, error) {
	var totalDeleted int64
	batchSize := 1000

	for {
		// 分批删除，避免锁表
		result := database.DB.WithContext(ctx).
			Where("created_at < ?", cutoffTime).
			Limit(batchSize).
			Delete(&models.GroupBuy{})

		if result.Error != nil {
			return totalDeleted, result.Error
		}

		deleted := result.RowsAffected
		totalDeleted += deleted

		// 如果没有更多数据需要删除，退出循环
		if deleted < int64(batchSize) {
			break
		}

		// 短暂休息，避免对数据库造成压力
		time.Sleep(100 * time.Millisecond)
	}

	return totalDeleted, nil
}

// GetCleanupStats 获取清理统计
func (s *DataCleanupService) GetCleanupStats() (*CleanupStats, error) {
	// 这里可以实现获取历史清理统计信息的逻辑
	return &CleanupStats{}, nil
}
