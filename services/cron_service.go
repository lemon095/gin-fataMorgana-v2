package services

import (
	"fmt"
	"log"
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
		log.Println("❌ 定时任务服务已禁用")
		return nil
	}

	log.Println("🚀 启动定时任务服务...")

	// 启动订单生成定时任务
	log.Println("⏰ 启动订单生成定时任务...")
	if err := s.StartFakeOrderCron(); err != nil {

		return err
	}

	// 启动数据清理定时任务
	log.Println("🧹 启动数据清理定时任务...")
	if err := s.StartCleanupCron(); err != nil {
		log.Printf("❌ 启动数据清理定时任务失败: %v", err)
		return err
	}

	// 启动热榜缓存更新定时任务
	log.Println("🏆 启动热榜缓存更新定时任务...")
	if err := s.StartLeaderboardCacheCron(); err != nil {
		log.Printf("❌ 启动热榜缓存更新定时任务失败: %v", err)
		return err
	}

	// 启动cron调度器
	log.Println("⚙️  启动cron调度器...")
	s.cron.Start()

	log.Println("✅ 定时任务服务启动成功")
	// 获取下次执行时间
	if s.orderEntryID != 0 {
		entries := s.cron.Entries()
		for _, entry := range entries {
			if entry.ID == s.orderEntryID {
				log.Printf("📅 下次订单生成时间: %s", entry.Next.Format("2006-01-02 15:04:05"))
				break
			}
		}
	}
	return nil
}

// Stop 停止定时任务服务
func (s *CronService) Stop() {
	if s.cron != nil {
		log.Println("停止定时任务服务...")
		s.cron.Stop()
		log.Println("定时任务服务已停止")
	}
}

// StartFakeOrderCron 启动假订单生成定时任务
func (s *CronService) StartFakeOrderCron() error {
	if s.config.OrderCronExpr == "" {
		s.config.OrderCronExpr = "0 */5 * * * *" // 默认每5分钟（包含秒）
	}

	log.Printf("⏰ 验证cron表达式: %s", s.config.OrderCronExpr)

	entryID, err := s.cron.AddFunc(s.config.OrderCronExpr, s.generateFakeOrders)
	if err != nil {
		log.Printf("❌ cron表达式验证失败: %v", err)
		return err
	}

	s.orderEntryID = entryID
	log.Printf("✅ 假订单生成定时任务已启动，表达式: %s", s.config.OrderCronExpr)
	return nil
}

// StopFakeOrderCron 停止假订单生成定时任务
func (s *CronService) StopFakeOrderCron() {
	if s.orderEntryID != 0 {
		s.cron.Remove(s.orderEntryID)
		s.orderEntryID = 0
		log.Println("假订单生成定时任务已停止")
	}
}

// StartCleanupCron 启动数据清理定时任务
func (s *CronService) StartCleanupCron() error {
	if s.config.CleanupCronExpr == "" {
		s.config.CleanupCronExpr = "0 0 2 * * *" // 默认每天凌晨2点（包含秒）
	}

	log.Printf("🧹 验证清理cron表达式: %s", s.config.CleanupCronExpr)

	entryID, err := s.cron.AddFunc(s.config.CleanupCronExpr, s.cleanupOldData)
	if err != nil {
		log.Printf("❌ 清理cron表达式验证失败: %v", err)
		return err
	}

	s.cleanupEntryID = entryID
	log.Printf("✅ 数据清理定时任务已启动，表达式: %s", s.config.CleanupCronExpr)
	return nil
}

// StopCleanupCron 停止数据清理定时任务
func (s *CronService) StopCleanupCron() {
	if s.cleanupEntryID != 0 {
		s.cron.Remove(s.cleanupEntryID)
		s.cleanupEntryID = 0
		log.Println("数据清理定时任务已停止")
	}
}

// StartLeaderboardCacheCron 启动热榜缓存更新定时任务
func (s *CronService) StartLeaderboardCacheCron() error {
	if s.config.LeaderboardCronExpr == "" {
		s.config.LeaderboardCronExpr = "0 */5 * * * *" // 默认每5分钟（包含秒）
	}

	log.Printf("🏆 验证热榜缓存cron表达式: %s", s.config.LeaderboardCronExpr)

	entryID, err := s.cron.AddFunc(s.config.LeaderboardCronExpr, s.updateLeaderboardCache)
	if err != nil {
		log.Printf("❌ 热榜缓存cron表达式验证失败: %v", err)
		return err
	}

	s.leaderboardEntryID = entryID
	log.Printf("✅ 热榜缓存更新定时任务已启动，表达式: %s", s.config.LeaderboardCronExpr)
	return nil
}

// StopLeaderboardCacheCron 停止热榜缓存更新定时任务
func (s *CronService) StopLeaderboardCacheCron() {
	if s.leaderboardEntryID != 0 {
		s.cron.Remove(s.leaderboardEntryID)
		s.leaderboardEntryID = 0
		log.Println("热榜缓存更新定时任务已停止")
	}
}

// generateFakeOrders 生成假订单（定时任务回调函数）
func (s *CronService) generateFakeOrders() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("生成假订单时发生panic: %v", r)
		}
	}()

	log.Println("=== 开始执行假订单生成定时任务 ===")
	log.Printf("当前时间: %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Printf("定时任务配置: 最小订单数=%d, 最大订单数=%d, 购买单比例=%.2f",
		s.config.MinOrders, s.config.MaxOrders, s.config.PurchaseRatio)

	startTime := time.Now()

	// 生成随机订单数量
	count := 0
	if s.config.MaxOrders > s.config.MinOrders {
		count = rand.Intn(s.config.MaxOrders-s.config.MinOrders+1) + s.config.MinOrders
	} else {
		count = s.config.MinOrders
	}

	log.Printf("本次将生成 %d 条假订单", count)

	// 生成假订单
	log.Println("开始调用假订单生成服务...")
	stats, err := s.fakeOrderService.GenerateFakeOrders(count)
	if err != nil {
		log.Printf("❌ 生成假订单失败: %v", err)
		return
	}

	duration := time.Since(startTime)
	log.Printf("✅ 假订单生成定时任务完成: 总数=%d, 购买单=%d, 拼单=%d, 总金额=%.2f, 总利润=%.2f, 耗时=%v",
		stats.TotalGenerated, stats.PurchaseOrders, stats.GroupBuyOrders,
		stats.TotalAmount, stats.TotalProfit, duration)
	log.Println("=== 假订单生成定时任务结束 ===")
}

// cleanupOldData 清理旧数据（定时任务回调函数）
func (s *CronService) cleanupOldData() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("清理旧数据时发生panic: %v", r)
		}
	}()

	log.Println("开始执行数据清理定时任务...")
	startTime := time.Now()

	// 清理旧数据
	stats, err := s.dataCleanupService.CleanupOldSystemOrders()
	if err != nil {
		log.Printf("清理旧数据失败: %v", err)
		return
	}

	duration := time.Since(startTime)
	log.Printf("数据清理定时任务完成: 删除订单=%d, 删除拼单=%d, 耗时=%v",
		stats.DeletedOrders, stats.DeletedGroupBuys, duration)
}

// updateLeaderboardCache 更新热榜缓存（定时任务回调函数）
func (s *CronService) updateLeaderboardCache() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("更新热榜缓存时发生panic: %v", r)
		}
	}()

	log.Println("=== 开始执行热榜缓存更新定时任务 ===")
	log.Printf("当前时间: %s", time.Now().Format("2006-01-02 15:04:05"))

	startTime := time.Now()

	// 更新热榜缓存
	err := s.leaderboardCacheService.UpdateLeaderboardCache()
	if err != nil {
		log.Printf("❌ 更新热榜缓存失败: %v", err)
		return
	}

	duration := time.Since(startTime)
	log.Printf("✅ 热榜缓存更新定时任务完成，耗时=%v", duration)
	log.Println("=== 热榜缓存更新定时任务结束 ===")
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
	log.Printf("手动生成 %d 条假订单", count)
	return s.fakeOrderService.GenerateFakeOrders(count)
}

// ManualCleanup 手动清理数据
func (s *CronService) ManualCleanup() (*CleanupStats, error) {
	log.Println("手动清理旧数据")
	return s.dataCleanupService.CleanupOldSystemOrders()
}

// ManualUpdateLeaderboardCache 手动更新热榜缓存
func (s *CronService) ManualUpdateLeaderboardCache() error {
	log.Println("手动更新热榜缓存")
	return s.leaderboardCacheService.UpdateLeaderboardCache()
}
