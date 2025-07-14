package services

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"
)

// FakeOrderService 假订单生成服务
type FakeOrderService struct {
	orderRepo    *database.OrderRepository
	groupBuyRepo *database.GroupBuyRepository
	config       *FakeOrderConfig
	periodCache  map[string]string // 缓存时间段对应的期号
}

// FakeOrderConfig 假订单配置
type FakeOrderConfig struct {
	MinOrders     int              `yaml:"min_orders"`
	MaxOrders     int              `yaml:"max_orders"`
	PurchaseRatio float64          `yaml:"purchase_ratio"`
	TaskMinCount  int              `yaml:"task_min_count"`
	TaskMaxCount  int              `yaml:"task_max_count"`
	TimeWindow    TimeWindowConfig `yaml:"time_window"`
}

// TimeWindowConfig 时间窗口配置
type TimeWindowConfig struct {
	BeforeMinutes int `yaml:"before_minutes"`
	AfterMinutes  int `yaml:"after_minutes"`
	TotalWindow   int `yaml:"total_window"`
}

// GenerationStats 生成统计
type GenerationStats struct {
	TotalGenerated int64         `json:"total_generated"`
	PurchaseOrders int64         `json:"purchase_orders"`
	GroupBuyOrders int64         `json:"group_buy_orders"`
	LastGeneration time.Time     `json:"last_generation"`
	AverageTime    time.Duration `json:"average_time"`
	TotalAmount    float64       `json:"total_amount"`
	TotalProfit    float64       `json:"total_profit"`
}

// NewFakeOrderService 创建新的假订单生成服务
func NewFakeOrderService(config *FakeOrderConfig) *FakeOrderService {
	return &FakeOrderService{
		orderRepo:    database.NewOrderRepository(),
		groupBuyRepo: database.NewGroupBuyRepository(),
		config:       config,
		periodCache:  make(map[string]string),
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

	// 预加载期数数据到缓存
	if err := s.preloadPeriodData(); err != nil {
		return nil, err
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
				continue
			}
		}
	}

	// 逐个插入拼单
	if len(groupBuyOrders) > 0 {
		for _, groupBuy := range groupBuyOrders {
			if err := s.groupBuyRepo.Create(ctx, groupBuy); err != nil {
				continue
			}
		}
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

	return stats, nil
}

// generatePurchaseOrder 生成购买单
func (s *FakeOrderService) generatePurchaseOrder() *models.Order {
	// 生成随机创建时间（10分钟窗口）
	createdAt := s.generateRandomTime()

	// 随机选择1-4个类型，每个类型数量为1
	likeCount := 0
	shareCount := 0
	followCount := 0
	favoriteCount := 0

	// 随机选择类型数量（1-4个）
	typeCount := rand.Intn(4) + 1

	// 创建类型数组并随机打乱
	types := []string{"like", "share", "follow", "favorite"}
	rand.Shuffle(len(types), func(i, j int) {
		types[i], types[j] = types[j], types[i]
	})

	// 选择前typeCount个类型，数量设为1
	for i := 0; i < typeCount; i++ {
		switch types[i] {
		case "like":
			likeCount = 1
		case "share":
			shareCount = 1
		case "follow":
			followCount = 1
		case "favorite":
			favoriteCount = 1
		}
	}

	// 生成总金额（10万到1000万之间）
	totalAmount := float64(rand.Intn(9900000) + 100000) // 100000-10000000

	// 假购买订单不计算利润金额
	profitAmount := 0.0

	// 随机选择状态
	status := s.getRandomPurchaseStatus()

	// 根据状态设置过期时间
	expireTime := s.getStatusBasedExpireTime(status, createdAt)

	order := &models.Order{
		OrderNo:        utils.GenerateSystemOrderNo(),
		Uid:            utils.GenerateSystemUID(),
		PeriodNumber:   s.getPeriodNumberByTime(createdAt),
		Amount:         totalAmount,
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

	// 随机选择1-4个类型，每个类型数量为1
	likeCount := 0
	shareCount := 0
	followCount := 0
	favoriteCount := 0

	// 随机选择类型数量（1-4个）
	typeCount := rand.Intn(4) + 1

	// 创建类型数组并随机打乱
	types := []string{"like", "share", "follow", "favorite"}
	rand.Shuffle(len(types), func(i, j int) {
		types[i], types[j] = types[j], types[i]
	})

	// 选择前typeCount个类型，数量设为1
	for i := 0; i < typeCount; i++ {
		switch types[i] {
		case "like":
			likeCount = 1
		case "share":
			shareCount = 1
		case "follow":
			followCount = 1
		case "favorite":
			favoriteCount = 1
		}
	}

	// 随机生成单价（1万到10万之间）
	unitPrice := float64(rand.Intn(90000) + 10000) // 10000-100000

	// 计算总任务数量
	totalTaskCount := likeCount + shareCount + followCount + favoriteCount

	// 计算总金额：单价 × 总任务数量
	totalAmount := unitPrice * float64(totalTaskCount)

	// 随机生成参与人数和目标人数
	currentParticipants := rand.Intn(3) + 1 // 1-3人
	targetParticipants := rand.Intn(5) + 3  // 3-7人

	// 计算人均金额：总金额 ÷ 目标人数
	perPersonAmount := totalAmount / float64(targetParticipants)

	// 随机生成利润比例（110%-160%）
	profitMargin := float64(rand.Intn(51)+110) / 100.0 // 1.10 - 1.60

	// 随机选择状态
	status := s.getRandomGroupBuyStatus()

	// 根据状态设置截止时间
	deadline := s.getGroupBuyDeadline(status, createdAt)

	groupBuy := &models.GroupBuy{
		GroupBuyNo:          utils.GenerateSystemGroupBuyNo(),
		Uid:                 utils.GenerateSystemUID(),
		CreatorUid:          utils.GenerateSystemUID(), // 创建者UID
		CurrentParticipants: currentParticipants,
		TargetParticipants:  targetParticipants,
		GroupBuyType:        models.GroupBuyTypeNormal,
		TotalAmount:         totalAmount,
		PaidAmount:          perPersonAmount * float64(currentParticipants),
		PerPersonAmount:     perPersonAmount,
		ProfitMargin:        profitMargin, // 添加利润比例
		Status:              status,
		CreatedAt:           createdAt,
		UpdatedAt:           createdAt,
		Deadline:            deadline,
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

// getPurchaseConfig 获取价格配置（已废弃，不再使用缓存的价格配置）
func (s *FakeOrderService) getPurchaseConfig() *models.PurchaseConfig {
	// 新的逻辑不再使用缓存的价格配置
	return &models.PurchaseConfig{
		LikeAmount:     0.1,
		ShareAmount:    0.2,
		ForwardAmount:  0.3,
		FavoriteAmount: 0.4,
	}
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

// calculateGroupBuyAmount 计算拼单人均金额（已废弃，新的逻辑直接生成随机金额）
func (s *FakeOrderService) calculateGroupBuyAmount() float64 {
	// 新的逻辑直接生成1万到10万之间的随机金额
	return float64(rand.Intn(90000) + 10000) // 10000-100000
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
