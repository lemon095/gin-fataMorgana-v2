package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"time"
)

// LeaderboardService 热榜服务
type LeaderboardService struct {
	leaderboardRepo *database.LeaderboardRepository
}

// NewLeaderboardService 创建热榜服务实例
func NewLeaderboardService() *LeaderboardService {
	return &LeaderboardService{
		leaderboardRepo: database.NewLeaderboardRepository(),
	}
}

// GetLeaderboard 获取任务热榜
func (s *LeaderboardService) GetLeaderboard(uid string) (*models.LeaderboardResponse, error) {
	ctx := context.Background()

	// 获取本周时间范围
	weekStart, weekEnd := models.GetCurrentWeekRange()

	// 生成缓存key
	cacheKey := fmt.Sprintf("leaderboard:weekly:%s", weekStart.Format("2006-01-02"))

	// 尝试从缓存获取数据
	var response *models.LeaderboardResponse
	cachedData, err := database.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		// 缓存命中，解析数据
		if err := json.Unmarshal([]byte(cachedData), &response); err == nil {
			// 更新我的排名信息（因为每个用户的排名可能不同）
			s.updateMyRankInfo(response, uid, weekStart, weekEnd)
			return response, nil
		}
	}

	// 缓存未命中或解析失败，从数据库查询
	response, err = s.buildLeaderboardResponse(uid, weekStart, weekEnd)
	if err != nil {
		return nil, fmt.Errorf("获取热榜数据失败: %w", err)
	}

	// 缓存数据（5分钟）
	cacheExpire := time.Now().Add(5 * time.Minute)
	response.CacheExpire = cacheExpire

	cacheData, err := json.Marshal(response)
	if err == nil {
		database.RedisClient.Set(ctx, cacheKey, cacheData, 5*time.Minute)
	}

	return response, nil
}

// buildLeaderboardResponse 构建热榜响应
func (s *LeaderboardService) buildLeaderboardResponse(uid string, weekStart, weekEnd time.Time) (*models.LeaderboardResponse, error) {
	ctx := context.Background()

	// 获取前10名用户数据
	topUsers, err := s.leaderboardRepo.GetWeeklyLeaderboard(ctx, weekStart, weekEnd)
	if err != nil {
		return nil, fmt.Errorf("获取热榜数据失败: %w", err)
	}

	// 转换为响应格式
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

	// 获取我的排名信息
	myRank := s.getMyRankInfo(uid, weekStart, weekEnd, topEntries)

	response := &models.LeaderboardResponse{
		WeekStart:   weekStart,
		WeekEnd:     weekEnd,
		MyRank:      myRank,
		TopUsers:    topEntries,
		CacheExpire: time.Now().Add(5 * time.Minute),
	}

	return response, nil
}

// getMyRankInfo 获取我的排名信息
func (s *LeaderboardService) getMyRankInfo(uid string, weekStart, weekEnd time.Time, topEntries []models.LeaderboardEntry) *models.LeaderboardEntry {
	ctx := context.Background()

	// 检查是否在前10名中
	for _, entry := range topEntries {
		if entry.Uid == uid {
			return &entry
		}
	}

	// 不在前10名中，查询具体排名
	userData, rank, err := s.leaderboardRepo.GetUserWeeklyRank(ctx, uid, weekStart, weekEnd)
	if err != nil || userData == nil {
		// 用户没有完成订单或查询失败
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

	// 用户有完成订单但不在前10名
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

// updateMyRankInfo 更新我的排名信息
func (s *LeaderboardService) updateMyRankInfo(response *models.LeaderboardResponse, uid string, weekStart, weekEnd time.Time) {
	// 检查是否在前10名中
	for _, entry := range response.TopUsers {
		if entry.Uid == uid {
			response.MyRank = &entry
			return
		}
	}

	// 不在前10名中，查询具体排名
	myRank := s.getMyRankInfo(uid, weekStart, weekEnd, response.TopUsers)
	response.MyRank = myRank
}
