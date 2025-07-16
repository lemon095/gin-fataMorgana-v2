package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"
	"time"
)

// LeaderboardCacheService 热榜缓存服务
type LeaderboardCacheService struct {
	leaderboardRepo *database.LeaderboardRepository
}

// CachedLeaderboardData 缓存的排行榜数据
type CachedLeaderboardData struct {
	WeekStart  time.Time                 `json:"week_start"`
	WeekEnd    time.Time                 `json:"week_end"`
	TopUsers   []models.LeaderboardEntry `json:"top_users"`
	NextUpdate time.Time                 `json:"next_update"`
	CacheTime  time.Time                 `json:"cache_time"`
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

	// 查询前10名用户数据
	topUsers, err := s.leaderboardRepo.GetWeeklyLeaderboard(ctx, weekStart, weekEnd)
	if err != nil {

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

		return utils.NewAppError(utils.CodeDatabaseError, "序列化缓存数据失败")
	}

	// 生成缓存键
	cacheKey := fmt.Sprintf("leaderboard:weekly:%s", weekStart.Format("2006-01-02"))

	// 存储到Redis，设置6分钟过期时间
	err = database.RedisClient.Set(ctx, cacheKey, jsonData, 6*time.Minute).Err()
	if err != nil {

		return utils.NewAppError(utils.CodeDatabaseError, "存储缓存数据失败")
	}

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
		return nil, err
	}

	// 反序列化缓存数据
	var cacheData CachedLeaderboardData
	if err := json.Unmarshal([]byte(cachedData), &cacheData); err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "反序列化缓存数据失败")
	}
	return &cacheData, nil
}

// GetUserRankFromCache 从缓存数据中获取用户排名
func (s *LeaderboardCacheService) GetUserRankFromCache(uid string, cachedData *CachedLeaderboardData) *models.LeaderboardEntry {
	// 先检查用户是否在缓存的前10名中
	for _, entry := range cachedData.TopUsers {
		if entry.Uid == uid {
			return &entry
		}
	}

	// 如果用户不在前10名中，需要实时查询用户排名

	ctx := context.Background()
	userData, rank, err := s.leaderboardRepo.GetUserWeeklyRank(ctx, uid, cachedData.WeekStart, cachedData.WeekEnd)
	if err != nil {
		return s.getDefaultUserRankInfo(uid)
	}

	if userData == nil {
		return s.getDefaultUserRankInfo(uid)
	}

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
	userRepo := database.NewUserRepository()
	user, err := userRepo.FindByUid(context.Background(), uid)
	if err != nil {
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
