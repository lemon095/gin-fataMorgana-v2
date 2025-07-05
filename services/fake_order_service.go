package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"
)

// FakeOrderService 假订单生成服务
type FakeOrderService struct {
	orderRepo       *database.OrderRepository
	groupBuyRepo    *database.GroupBuyRepository
	memberLevelRepo *database.MemberLevelRepository
	config          *FakeOrderConfig
	periodCache     map[string]string // 缓存时间段对应的期号
}

// FakeOrderConfig 假订单配置
type FakeOrderConfig struct {
	MinOrders       int     `yaml:"min_orders"`
	MaxOrders       int     `yaml:"max_orders"`
	PurchaseRatio   float64 `yaml:"purchase_ratio"`
	TaskMinCount    int     `yaml:"task_min_count"`
	TaskMaxCount    int     `yaml:"task_max_count"`
	TimeWindow      TimeWindowConfig `yaml:"time_window"`
}

// TimeWindowConfig 时间窗口配置
type TimeWindowConfig struct {
	BeforeMinutes int `yaml:"before_minutes"`
	AfterMinutes  int `yaml:"after_minutes"`
	TotalWindow   int `yaml:"total_window"`
}

// GenerationStats 生成统计
type GenerationStats struct {
	TotalGenerated    int64         `json:"total_generated"`
	PurchaseOrders    int64         `json:"purchase_orders"`
	GroupBuyOrders    int64         `json:"group_buy_orders"`
	LastGeneration    time.Time     `json:"last_generation"`
	AverageTime       time.Duration `json:"average_time"`
	TotalAmount       float64       `json:"total_amount"`
	TotalProfit       float64       `json:"total_profit"`
}

// NewFakeOrderService 创建新的假订单生成服务
func NewFakeOrderService(config *FakeOrderConfig) *FakeOrderService {
	return &FakeOrderService{
		orderRepo:       database.NewOrderRepository(),
		groupBuyRepo:    database.NewGroupBuyRepository(),
		memberLevelRepo: database.NewMemberLevelRepository(database.DB),
		config:          config,
		periodCache:     make(map[string]string),
	}
}

// GenerateFakeOrders 生成假订单
func (s *FakeOrderService) GenerateFakeOrders(count int) (*GenerationStats, error) {
	startTime := time.Now()
	ctx := context.Background()

	// 生成随机订单数量
	if count <= 0 {
		count = rand.Intn(s.config.MaxOrders-s.config.MinOrders+1) + s.config.MinOrders
	}

	log.Printf("开始生成 %d 条假订单", count)

	// 预加载期数数据到缓存
	if err := s.preloadPeriodData(); err != nil {
		log.Printf("预加载期数数据失败: %v", err)
	}

	var purchaseOrders []*models.Order
	var groupBuyOrders []*models.GroupBuy
	var totalAmount, totalProfit float64

	// 生成订单
	for i := 0; i < count; i++ {
		if rand.Float64() < s.config.PurchaseRatio {
			// 生成购买单
			order := s.generatePurchaseOrder()
			purchaseOrders = append(purchaseOrders, order)
			totalAmount += order.Amount
			totalProfit += order.ProfitAmount
		} else {
			// 生成拼单
			groupBuy := s.generateGroupBuyOrder()
			groupBuyOrders = append(groupBuyOrders, groupBuy)
			totalAmount += groupBuy.PerPersonAmount
		}
	}

	// 逐个插入购买单
	if len(purchaseOrders) > 0 {
		for _, order := range purchaseOrders {
			if err := s.orderRepo.CreateOrder(ctx, order); err != nil {
				log.Printf("插入购买单失败: %v", err)
				continue
			}
		}
		log.Printf("成功插入 %d 条购买单", len(purchaseOrders))
	}

	// 逐个插入拼单
	if len(groupBuyOrders) > 0 {
		for _, groupBuy := range groupBuyOrders {
			if err := s.groupBuyRepo.Create(ctx, groupBuy); err != nil {
				log.Printf("插入拼单失败: %v", err)
				continue
			}
		}
		log.Printf("成功插入 %d 条拼单", len(groupBuyOrders))
	}

	duration := time.Since(startTime)

	stats := &GenerationStats{
		TotalGenerated: int64(count),
		PurchaseOrders: int64(len(purchaseOrders)),
		GroupBuyOrders: int64(len(groupBuyOrders)),
		LastGeneration: time.Now(),
		AverageTime:    duration,
		TotalAmount:    totalAmount,
		TotalProfit:    totalProfit,
	}

	log.Printf("假订单生成完成: 总数=%d, 购买单=%d, 拼单=%d, 总金额=%.2f, 总利润=%.2f, 耗时=%v",
		stats.TotalGenerated, stats.PurchaseOrders, stats.GroupBuyOrders,
		stats.TotalAmount, stats.TotalProfit, stats.AverageTime)

	return stats, nil
}

// generatePurchaseOrder 生成购买单
func (s *FakeOrderService) generatePurchaseOrder() *models.Order {
	// 生成随机创建时间（10分钟窗口）
	createdAt := s.generateRandomTime()
	
	// 生成任务数量
	likeCount := rand.Intn(s.config.TaskMaxCount-s.config.TaskMinCount+1) + s.config.TaskMinCount
	shareCount := rand.Intn(s.config.TaskMaxCount-s.config.TaskMinCount+1) + s.config.TaskMinCount
	followCount := rand.Intn(s.config.TaskMaxCount-s.config.TaskMinCount+1) + s.config.TaskMinCount
	favoriteCount := rand.Intn(s.config.TaskMaxCount-s.config.TaskMinCount+1) + s.config.TaskMinCount

	// 获取价格配置
	purchaseConfig := s.getPurchaseConfig()

	// 计算订单金额
	amount := float64(likeCount)*purchaseConfig.LikeAmount +
		float64(shareCount)*purchaseConfig.ShareAmount +
		float64(followCount)*purchaseConfig.ForwardAmount +
		float64(favoriteCount)*purchaseConfig.FavoriteAmount

	// 随机选择用户等级计算利润
	profitAmount := s.calculateProfitAmount(amount)

	// 随机选择状态
	status := s.getRandomPurchaseStatus()
	
	// 根据状态设置过期时间
	expireTime := s.getStatusBasedExpireTime(status, createdAt)

	order := &models.Order{
		OrderNo:        utils.GenerateSystemOrderNo(),
		Uid:            utils.GenerateSystemUID(),
		PeriodNumber:   s.getPeriodNumberByTime(createdAt),
		Amount:         amount,
		ProfitAmount:   profitAmount,
		Status:         status,
		ExpireTime:     expireTime,
		LikeCount:      likeCount,
		ShareCount:     shareCount,
		FollowCount:    followCount,
		FavoriteCount:  favoriteCount,
		LikeStatus:     s.getTaskStatus(likeCount, status),
		ShareStatus:    s.getTaskStatus(shareCount, status),
		FollowStatus:   s.getTaskStatus(followCount, status),
		FavoriteStatus: s.getTaskStatus(favoriteCount, status),
		IsSystemOrder:  true,
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt, // 确保更新时间也是过去时间
	}

	return order
}

// generateGroupBuyOrder 生成拼单
func (s *FakeOrderService) generateGroupBuyOrder() *models.GroupBuy {
	// 生成随机创建时间
	createdAt := s.generateRandomTime()
	
	// 基于价格配置计算人均金额
	perPersonAmount := s.calculateGroupBuyAmount()

	// 随机选择状态
	status := s.getRandomGroupBuyStatus()
	
	// 根据状态设置截止时间
	deadline := s.getGroupBuyDeadline(status, createdAt)

	// 随机生成参与人数和目标人数
	currentParticipants := rand.Intn(3) + 1 // 1-3人
	targetParticipants := rand.Intn(5) + 3  // 3-7人
	
	// 计算总金额
	totalAmount := perPersonAmount * float64(targetParticipants)

	groupBuy := &models.GroupBuy{
		GroupBuyNo:        utils.GenerateSystemGroupBuyNo(),
		Uid:               utils.GenerateSystemUID(),
		CreatorUid:        utils.GenerateSystemUID(), // 创建者UID
		CurrentParticipants: currentParticipants,
		TargetParticipants:  targetParticipants,
		GroupBuyType:      models.GroupBuyTypeNormal,
		TotalAmount:       totalAmount,
		PaidAmount:        perPersonAmount * float64(currentParticipants),
		PerPersonAmount:   perPersonAmount,
		Status:            status,
		CreatedAt:         createdAt,
		UpdatedAt:         createdAt,
		Deadline:          deadline,
	}

	return groupBuy
}

// generateRandomTime 生成随机时间（过去10分钟到未来10分钟）
func (s *FakeOrderService) generateRandomTime() time.Time {
	now := time.Now()
	
	// 时间窗口：当前时间前后各10分钟（过去10分钟到未来10分钟）
	startTime := now.Add(-10 * time.Minute)
	endTime := now.Add(10 * time.Minute)
	
	// 计算时间差
	timeDiff := endTime.Sub(startTime)
	
	// 生成随机时间偏移
	randomOffset := time.Duration(rand.Int63n(int64(timeDiff)))
	
	return startTime.Add(randomOffset)
}

// getPurchaseConfig 获取价格配置
func (s *FakeOrderService) getPurchaseConfig() *models.PurchaseConfig {
	// 从Redis获取价格配置，如果获取失败则使用默认配置
	ctx := context.Background()
	configJSON, err := database.RedisClient.Get(ctx, "purchase_config").Result()
	if err != nil {
		// 返回默认配置
		return &models.PurchaseConfig{
			LikeAmount:     0.1,
			ShareAmount:    0.2,
			ForwardAmount:  0.3,
			FavoriteAmount: 0.4,
		}
	}

	var config models.PurchaseConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return &models.PurchaseConfig{
			LikeAmount:     0.1,
			ShareAmount:    0.2,
			ForwardAmount:  0.3,
			FavoriteAmount: 0.4,
		}
	}

	return &config
}

// calculateProfitAmount 计算利润金额
func (s *FakeOrderService) calculateProfitAmount(amount float64) float64 {
	ctx := context.Background()
	
	// 随机选择用户等级（1-10级）
	randomLevel := rand.Intn(10) + 1
	
	// 根据等级获取返现比例
	level, err := s.memberLevelRepo.GetByLevel(ctx, randomLevel)
	if err != nil {
		// 如果获取失败，使用默认比例5%
		return amount * 0.05
	}
	
	// 计算利润金额：订单金额 × (返现比例 / 100)
	profitAmount := amount * (level.CashbackRatio / 100.0)
	return profitAmount
}

// getRandomPurchaseStatus 获取随机购买单状态
func (s *FakeOrderService) getRandomPurchaseStatus() string {
	randNum := rand.Float64()
	
	if randNum < 0.6 {
		return models.OrderStatusPending // 60% 进行中
	} else if randNum < 0.9 {
		return models.OrderStatusSuccess // 30% 已完成
	} else {
		return models.OrderStatusCancelled // 10% 已关闭
	}
}

// getRandomGroupBuyStatus 获取随机拼单状态
func (s *FakeOrderService) getRandomGroupBuyStatus() string {
	randNum := rand.Float64()
	
	if randNum < 0.2 {
		return models.GroupBuyStatusNotStarted // 20% 待开始
	} else if randNum < 0.7 {
		return models.GroupBuyStatusPending // 50% 进行中
	} else {
		return models.GroupBuyStatusSuccess // 30% 已完成
	}
}

// getStatusBasedExpireTime 根据状态设置过期时间
func (s *FakeOrderService) getStatusBasedExpireTime(status string, createdAt time.Time) time.Time {
	switch status {
	case models.OrderStatusPending:
		// 进行中：创建时间 + 5-15分钟
		return createdAt.Add(time.Duration(rand.Intn(10)+5) * time.Minute)
	case models.OrderStatusSuccess:
		// 已完成：创建时间 + 1-3分钟（快速完成）
		return createdAt.Add(time.Duration(rand.Intn(2)+1) * time.Minute)
	case models.OrderStatusCancelled:
		// 已关闭：创建时间 + 1-2分钟（快速关闭）
		return createdAt.Add(time.Duration(rand.Intn(1)+1) * time.Minute)
	default:
		return createdAt.Add(5 * time.Minute)
	}
}

// getGroupBuyDeadline 获取拼单截止时间
func (s *FakeOrderService) getGroupBuyDeadline(status string, createdAt time.Time) time.Time {
	switch status {
	case models.GroupBuyStatusNotStarted:
		// 待开始：创建时间 + 10-30分钟
		return createdAt.Add(time.Duration(rand.Intn(20)+10) * time.Minute)
	case models.GroupBuyStatusPending:
		// 进行中：创建时间 + 5-15分钟
		return createdAt.Add(time.Duration(rand.Intn(10)+5) * time.Minute)
	case models.GroupBuyStatusSuccess:
		// 已完成：创建时间 + 2-5分钟（快速完成）
		return createdAt.Add(time.Duration(rand.Intn(3)+2) * time.Minute)
	default:
		return createdAt.Add(10 * time.Minute)
	}
}

// getTaskStatus 获取任务状态
func (s *FakeOrderService) getTaskStatus(count int, orderStatus string) string {
	if count == 0 {
		return models.TaskStatusSuccess // 任务数为0时直接完成
	}
	
	// 如果订单状态是已完成，任务状态也应该是已完成
	if orderStatus == models.OrderStatusSuccess {
		return models.TaskStatusSuccess
	}
	
	// 如果订单状态是已关闭，任务状态也应该是已关闭
	if orderStatus == models.OrderStatusCancelled {
		return models.TaskStatusCancelled
	}
	
	// 如果订单状态是进行中，根据概率设置任务状态
	randNum := rand.Float64()
	if randNum < 0.3 {
		return models.TaskStatusSuccess // 30% 已完成
	} else {
		return models.TaskStatusPending // 70% 待完成
	}
}

// calculateGroupBuyAmount 计算拼单人均金额
func (s *FakeOrderService) calculateGroupBuyAmount() float64 {
	// 获取价格配置
	purchaseConfig := s.getPurchaseConfig()
	
	// 基于价格配置计算合理的拼单金额
	// 拼单金额应该是一个合理的任务组合价格
	maxTaskCountPerType := rand.Intn(5) + 3 // 每种任务类型的最大数量：3-7个
	
	// 随机选择任务类型组合
	taskTypes := []string{"like", "share", "follow", "favorite"}
	selectedTasks := make([]string, 0)
	
	// 随机选择2-4种任务类型
	numTaskTypes := rand.Intn(3) + 2 // 2-4种任务类型
	for i := 0; i < numTaskTypes; i++ {
		taskType := taskTypes[rand.Intn(len(taskTypes))]
		if !contains(selectedTasks, taskType) {
			selectedTasks = append(selectedTasks, taskType)
		}
	}
	
	// 计算人均金额
	var totalAmount float64
	for _, taskType := range selectedTasks {
		taskCount := rand.Intn(maxTaskCountPerType) + 1
		switch taskType {
		case "like":
			totalAmount += float64(taskCount) * purchaseConfig.LikeAmount
		case "share":
			totalAmount += float64(taskCount) * purchaseConfig.ShareAmount
		case "follow":
			totalAmount += float64(taskCount) * purchaseConfig.ForwardAmount
		case "favorite":
			totalAmount += float64(taskCount) * purchaseConfig.FavoriteAmount
		}
	}
	
	// 确保金额在合理范围内（5.00-50.00）
	if totalAmount < 5.0 {
		totalAmount = 5.0
	} else if totalAmount > 50.0 {
		totalAmount = 50.0
	}
	
	return totalAmount
}

// contains 检查切片是否包含元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// preloadPeriodData 预加载期数数据到缓存
func (s *FakeOrderService) preloadPeriodData() error {
	ctx := context.Background()
	periodRepo := database.NewLotteryPeriodRepository()
	
	// 清空缓存
	s.periodCache = make(map[string]string)
	
	// 获取当前时间前后30分钟的时间范围
	now := time.Now()
	startTime := now.Add(-30 * time.Minute)
	endTime := now.Add(30 * time.Minute)
	
	// 查询这个时间范围内的所有期数
	periods, err := periodRepo.GetPeriodsByTimeRange(ctx, startTime, endTime)
	if err != nil {
		return err
	}
	
	// 将期数数据缓存到内存中
	for _, period := range periods {
		// 使用期数的时间范围作为key
		key := fmt.Sprintf("%s_%s", period.OrderStartTime.Format("2006-01-02 15:04:05"), 
			period.OrderEndTime.Format("2006-01-02 15:04:05"))
		s.periodCache[key] = period.PeriodNumber
	}
	
	log.Printf("预加载了 %d 个期数到缓存", len(periods))
	return nil
}

// getPeriodNumberByTime 根据时间获取对应的期号（使用缓存）
func (s *FakeOrderService) getPeriodNumberByTime(targetTime time.Time) string {
	// 首先尝试从缓存中查找
	for key, periodNumber := range s.periodCache {
		// 解析key中的时间范围
		parts := strings.Split(key, "_")
		if len(parts) == 2 {
			startTime, _ := time.Parse("2006-01-02 15:04:05", parts[0])
			endTime, _ := time.Parse("2006-01-02 15:04:05", parts[1])
			
			// 检查目标时间是否在这个范围内
			if targetTime.After(startTime) && targetTime.Before(endTime) {
				return periodNumber
			}
		}
	}
	
	// 如果缓存中没有找到，回退到数据库查询
	ctx := context.Background()
	periodRepo := database.NewLotteryPeriodRepository()
	
	period, err := periodRepo.GetPeriodByTime(ctx, targetTime)
	if err != nil {
		// 如果获取失败，使用目标时间生成期号
		return targetTime.Format("20240101")
	}
	
	return period.PeriodNumber
}

// getCurrentPeriodNumber 获取当前期号
func (s *FakeOrderService) getCurrentPeriodNumber() string {
	return s.getPeriodNumberByTime(time.Now())
}

// GetGenerationStats 获取生成统计
func (s *FakeOrderService) GetGenerationStats() (*GenerationStats, error) {
	// 这里可以实现获取历史统计信息的逻辑
	return &GenerationStats{}, nil
} 