package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"
	"log"
	"time"
)

// LeaderboardCacheService 热榜缓存服务
type LeaderboardCacheService struct {
	leaderboardRepo *database.LeaderboardRepository
}

// CachedLeaderboardData 缓存的排行榜数据
type CachedLeaderboardData struct {
	WeekStart   time.Time          `json:"week_start"`
	WeekEnd     time.Time          `json:"week_end"`
	TopUsers    []models.LeaderboardEntry `json:"top_users"`
	NextUpdate  time.Time          `json:"next_update"`
	CacheTime   time.Time          `json:"cache_time"`
}

// NewLeaderboardCacheService 创建热榜缓存服务实例
func NewLeaderboardCacheService() *LeaderboardCacheService {
	return &LeaderboardCacheService{
		leaderboardRepo: database.NewLeaderboardRepository(),
	}
}

// UpdateLeaderboardCache 更新热榜缓存
func (s *LeaderboardCacheService) UpdateLeaderboardCache() error {
	ctx := context.Background()
	weekStart, weekEnd := models.GetCurrentWeekRange()
	
	log.Printf("🔄 [缓存服务] 开始更新热榜缓存")
	log.Printf("🔄 [缓存服务] 时间范围: %s 到 %s", weekStart.Format("2006-01-02 15:04:05"), weekEnd.Format("2006-01-02 15:04:05"))

	// 查询前10名用户数据
	topUsers, err := s.leaderboardRepo.GetWeeklyLeaderboard(ctx, weekStart, weekEnd)
	if err != nil {
		log.Printf("❌ [缓存服务] 查询排行榜数据失败: %v", err)
		return utils.NewAppError(utils.CodeDatabaseError, "获取热榜数据失败")
	}

	// 转换为LeaderboardEntry格式
	var topEntries []models.LeaderboardEntry
	for i, user := range topUsers {
		entry := models.LeaderboardEntry{
			ID:          uint(i + 1),
			Uid:         user.Uid,
			Username:    models.MaskUsername(user.Username),
			CompletedAt: user.CompletedAt,
			OrderCount:  user.OrderCount,
			TotalAmount: user.TotalAmount,
			TotalProfit: user.TotalProfit,
			Rank:        i + 1,
			IsRank:      true,
		}
		topEntries = append(topEntries, entry)
	}

	// 计算下次更新时间（当前时间+6分钟）
	nextUpdate := time.Now().Add(6 * time.Minute)
	
	// 构建缓存数据
	cacheData := &CachedLeaderboardData{
		WeekStart:  weekStart,
		WeekEnd:    weekEnd,
		TopUsers:   topEntries,
		NextUpdate: nextUpdate,
		CacheTime:  time.Now(),
	}

	// 序列化缓存数据
	jsonData, err := json.Marshal(cacheData)
	if err != nil {
		log.Printf("❌ [缓存服务] 序列化缓存数据失败: %v", err)
		return utils.NewAppError(utils.CodeDatabaseError, "序列化缓存数据失败")
	}

	// 生成缓存键
	cacheKey := fmt.Sprintf("leaderboard:weekly:%s", weekStart.Format("2006-01-02"))

	// 存储到Redis，设置6分钟过期时间
	err = database.RedisClient.Set(ctx, cacheKey, jsonData, 6*time.Minute).Err()
	if err != nil {
		log.Printf("❌ [缓存服务] 存储缓存数据失败: %v", err)
		return utils.NewAppError(utils.CodeDatabaseError, "存储缓存数据失败")
	}

	log.Printf("✅ [缓存服务] 热榜缓存更新成功，缓存了 %d 条数据，下次更新时间: %s", 
		len(topEntries), nextUpdate.Format("2006-01-02 15:04:05"))

	return nil
}

// GetCachedLeaderboardData 获取缓存的热榜数据
func (s *LeaderboardCacheService) GetCachedLeaderboardData() (*CachedLeaderboardData, error) {
	ctx := context.Background()
	weekStart, _ := models.GetCurrentWeekRange()
	
	// 生成缓存键
	cacheKey := fmt.Sprintf("leaderboard:weekly:%s", weekStart.Format("2006-01-02"))

	// 从Redis获取缓存数据
	cachedData, err := database.RedisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		log.Printf("⚠️ [缓存服务] 缓存未命中: %v", err)
		return nil, err
	}

	// 反序列化缓存数据
	var cacheData CachedLeaderboardData
	if err := json.Unmarshal([]byte(cachedData), &cacheData); err != nil {
		log.Printf("❌ [缓存服务] 反序列化缓存数据失败: %v", err)
		return nil, utils.NewAppError(utils.CodeDatabaseError, "反序列化缓存数据失败")
	}

	log.Printf("✅ [缓存服务] 成功获取缓存数据，包含 %d 条记录", len(cacheData.TopUsers))
	return &cacheData, nil
}

// GetUserRankFromCache 从缓存数据中获取用户排名
func (s *LeaderboardCacheService) GetUserRankFromCache(uid string, cachedData *CachedLeaderboardData) *models.LeaderboardEntry {
	// 先检查用户是否在缓存的前10名中
	for _, entry := range cachedData.TopUsers {
		if entry.Uid == uid {
			log.Printf("✅ [缓存服务] 用户 %s 在前10名中，排名第%d", uid, entry.Rank)
			return &entry
		}
	}

	// 如果用户不在前10名中，需要实时查询用户排名
	log.Printf("🔍 [缓存服务] 用户 %s 不在前10名中，需要实时查询排名", uid)
	
	ctx := context.Background()
	userData, rank, err := s.leaderboardRepo.GetUserWeeklyRank(ctx, uid, cachedData.WeekStart, cachedData.WeekEnd)
	if err != nil {
		log.Printf("❌ [缓存服务] 查询用户排名失败: %v", err)
		return s.getDefaultUserRankInfo(uid)
	}

	if userData == nil {
		log.Printf("⚠️ [缓存服务] 用户 %s 没有完成任何订单", uid)
		return s.getDefaultUserRankInfo(uid)
	}

	log.Printf("✅ [缓存服务] 用户 %s 排名第%d，完成订单数=%d，总金额=%.2f", 
		uid, rank, userData.OrderCount, userData.TotalAmount)

	return &models.LeaderboardEntry{
		ID:          uint(rank),
		Uid:         userData.Uid,
		Username:    models.MaskUsername(userData.Username),
		CompletedAt: userData.CompletedAt,
		OrderCount:  userData.OrderCount,
		TotalAmount: userData.TotalAmount,
		TotalProfit: userData.TotalProfit,
		Rank:        rank,
		IsRank:      false,
	}
}

// getDefaultUserRankInfo 获取默认用户排名信息
func (s *LeaderboardCacheService) getDefaultUserRankInfo(uid string) *models.LeaderboardEntry {
	log.Printf("🔍 [缓存服务] 获取用户 %s 的默认排名信息", uid)
	userRepo := database.NewUserRepository()
	user, err := userRepo.FindByUid(context.Background(), uid)
	if err != nil {
		log.Printf("❌ [缓存服务] 查询用户信息失败: %v", err)
		return &models.LeaderboardEntry{
			ID:          999,
			Uid:         uid,
			Username:    "",
			CompletedAt: time.Time{},
			OrderCount:  0,
			TotalAmount: 0,
			TotalProfit: 0,
			Rank:        999,
			IsRank:      false,
		}
	}

	log.Printf("✅ [缓存服务] 用户 %s 默认排名信息：用户名=%s", uid, user.Username)
	return &models.LeaderboardEntry{
		ID:          999,
		Uid:         user.Uid,
		Username:    models.MaskUsername(user.Username),
		CompletedAt: time.Time{},
		OrderCount:  0,
		TotalAmount: 0,
		TotalProfit: 0,
		Rank:        999,
		IsRank:      false,
	}
} 