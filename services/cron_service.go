package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/robfig/cron/v3"
)

// CronService 定时任务服务
type CronService struct {
	cron                    *cron.Cron
	fakeOrderService        *FakeOrderService
	dataCleanupService      *DataCleanupService
	leaderboardCacheService *LeaderboardCacheService
	config                  *CronConfig
	orderEntryID            cron.EntryID
	cleanupEntryID          cron.EntryID
	leaderboardEntryID      cron.EntryID
}

// CronConfig 定时任务配置
type CronConfig struct {
	Enabled             bool    `yaml:"enabled"`
	OrderCronExpr       string  `yaml:"order_cron_expr"`       // 订单生成定时表达式
	CleanupCronExpr     string  `yaml:"cleanup_cron_expr"`     // 数据清理定时表达式
	LeaderboardCronExpr string  `yaml:"leaderboard_cron_expr"` // 热榜缓存更新定时表达式
	MinOrders           int     `yaml:"min_orders"`
	MaxOrders           int     `yaml:"max_orders"`
	PurchaseRatio       float64 `yaml:"purchase_ratio"`
	TaskMinCount        int     `yaml:"task_min_count"`
	TaskMaxCount        int     `yaml:"task_max_count"`
	RetentionDays       int     `yaml:"retention_days"`
}

// NewCronService 创建新的定时任务服务
func NewCronService(config *CronConfig) *CronService {
	// 创建假订单配置
	fakeOrderConfig := &FakeOrderConfig{
		MinOrders:     config.MinOrders,
		MaxOrders:     config.MaxOrders,
		PurchaseRatio: config.PurchaseRatio,
		TaskMinCount:  config.TaskMinCount,
		TaskMaxCount:  config.TaskMaxCount,
		TimeWindow: TimeWindowConfig{
			BeforeMinutes: 5,
			AfterMinutes:  5,
			TotalWindow:   10,
		},
	}

	// 创建数据清理配置
	cleanupConfig := &DataCleanupConfig{
		RetentionDays: config.RetentionDays,
	}

	return &CronService{
		cron:                    cron.New(cron.WithSeconds()),
		fakeOrderService:        NewFakeOrderService(fakeOrderConfig),
		dataCleanupService:      NewDataCleanupService(cleanupConfig),
		leaderboardCacheService: NewLeaderboardCacheService(),
		config:                  config,
	}
}

// Start 启动定时任务服务
func (s *CronService) Start() error {
	if !s.config.Enabled {
		return nil
	}

	// 启动订单生成定时任务
	if err := s.StartFakeOrderCron(); err != nil {
		return err
	}

	// 启动数据清理定时任务
	if err := s.StartCleanupCron(); err != nil {
		return err
	}

	// 启动热榜缓存更新定时任务
	if err := s.StartLeaderboardCacheCron(); err != nil {
		return err
	}

	// 启动cron调度器
	s.cron.Start()

	return nil
}

// Stop 停止定时任务服务
func (s *CronService) Stop() {
	if s.cron != nil {
		s.cron.Stop()
	}
}

// StartFakeOrderCron 启动假订单生成定时任务
func (s *CronService) StartFakeOrderCron() error {
	if s.config.OrderCronExpr == "" {
		s.config.OrderCronExpr = "0 */5 * * * *" // 默认每5分钟（包含秒）
	}

	entryID, err := s.cron.AddFunc(s.config.OrderCronExpr, s.generateFakeOrders)
	if err != nil {
		return err
	}

	s.orderEntryID = entryID
	return nil
}

// StopFakeOrderCron 停止假订单生成定时任务
func (s *CronService) StopFakeOrderCron() {
	if s.orderEntryID != 0 {
		s.cron.Remove(s.orderEntryID)
		s.orderEntryID = 0
	}
}

// StartCleanupCron 启动数据清理定时任务
func (s *CronService) StartCleanupCron() error {
	if s.config.CleanupCronExpr == "" {
		s.config.CleanupCronExpr = "0 0 2 * * *" // 默认每天凌晨2点（包含秒）
	}

	entryID, err := s.cron.AddFunc(s.config.CleanupCronExpr, s.cleanupOldData)
	if err != nil {
		return err
	}

	s.cleanupEntryID = entryID
	return nil
}

// StopCleanupCron 停止数据清理定时任务
func (s *CronService) StopCleanupCron() {
	if s.cleanupEntryID != 0 {
		s.cron.Remove(s.cleanupEntryID)
		s.cleanupEntryID = 0
	}
}

// StartLeaderboardCacheCron 启动热榜缓存更新定时任务
func (s *CronService) StartLeaderboardCacheCron() error {
	if s.config.LeaderboardCronExpr == "" {
		s.config.LeaderboardCronExpr = "0 */5 * * * *" // 默认每5分钟（包含秒）
	}

	entryID, err := s.cron.AddFunc(s.config.LeaderboardCronExpr, s.updateLeaderboardCache)
	if err != nil {
		return err
	}

	s.leaderboardEntryID = entryID
	return nil
}

// StopLeaderboardCacheCron 停止热榜缓存更新定时任务
func (s *CronService) StopLeaderboardCacheCron() {
	if s.leaderboardEntryID != 0 {
		s.cron.Remove(s.leaderboardEntryID)
		s.leaderboardEntryID = 0
	}
}

// generateFakeOrders 生成假订单（定时任务回调函数）
func (s *CronService) generateFakeOrders() {
	defer func() {
		if r := recover(); r != nil {
			// 这里不记录日志，因为不是关键错误
		}
	}()

	startTime := time.Now()

	// 生成随机订单数量
	count := 0
	if s.config.MaxOrders > s.config.MinOrders {
		count = rand.Intn(s.config.MaxOrders-s.config.MinOrders+1) + s.config.MinOrders
	} else {
		count = s.config.MinOrders
	}

	// 生成假订单
	stats, err := s.fakeOrderService.GenerateFakeOrders(count)
	if err != nil {
		// 记录错误但不影响主流程
		return
	}

	duration := time.Since(startTime)
	// 这里不记录日志，因为不是关键信息
	_ = stats
	_ = duration
}

// cleanupOldData 清理旧数据（定时任务回调函数）
func (s *CronService) cleanupOldData() {
	defer func() {
		if r := recover(); r != nil {
			// 这里不记录日志，因为不是关键错误
		}
	}()

	startTime := time.Now()

	// 清理旧数据
	stats, err := s.dataCleanupService.CleanupOldSystemOrders()
	if err != nil {
		// 记录错误但不影响主流程
		return
	}

	duration := time.Since(startTime)
	// 这里不记录日志，因为不是关键信息
	_ = stats
	_ = duration
}

// updateLeaderboardCache 更新热榜缓存（定时任务回调函数）
func (s *CronService) updateLeaderboardCache() {
	defer func() {
		if r := recover(); r != nil {
			// 这里不记录日志，因为不是关键错误
		}
	}()

	startTime := time.Now()

	// 更新热榜缓存
	err := s.leaderboardCacheService.UpdateLeaderboardCache()
	if err != nil {
		// 记录错误但不影响主流程
		return
	}

	duration := time.Since(startTime)
	// 这里不记录日志，因为不是关键信息
	_ = duration
}

// GetCronStatus 获取定时任务状态
func (s *CronService) GetCronStatus() map[string]interface{} {
	entries := s.cron.Entries()
	status := make(map[string]interface{})

	for i, entry := range entries {
		status[fmt.Sprintf("task_%d", i)] = map[string]interface{}{
			"next_run": entry.Next,
			"prev_run": entry.Prev,
			"schedule": fmt.Sprintf("%v", entry.Schedule),
		}
	}

	return status
}

// ManualGenerateOrders 手动生成订单
func (s *CronService) ManualGenerateOrders(count int) (*GenerationStats, error) {
	return s.fakeOrderService.GenerateFakeOrders(count)
}

// ManualCleanup 手动清理数据
func (s *CronService) ManualCleanup() (*CleanupStats, error) {
	return s.dataCleanupService.CleanupOldSystemOrders()
}

// ManualUpdateLeaderboardCache 手动更新热榜缓存
func (s *CronService) ManualUpdateLeaderboardCache() error {
	return s.leaderboardCacheService.UpdateLeaderboardCache()
}
