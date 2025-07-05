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

type LeaderboardService struct {
	leaderboardRepo *database.LeaderboardRepository
}

func NewLeaderboardService() *LeaderboardService {
	return &LeaderboardService{
		leaderboardRepo: database.NewLeaderboardRepository(),
	}
}

func (s *LeaderboardService) GetLeaderboard(uid string) (*models.LeaderboardResponse, error) {
	ctx := context.Background()
	weekStart, weekEnd := models.GetCurrentWeekRange()
	cacheKey := fmt.Sprintf("leaderboard:weekly:%s", weekStart.Format("2006-01-02"))

	var response *models.LeaderboardResponse
	cachedData, err := database.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		if err := json.Unmarshal([]byte(cachedData), &response); err == nil {
			s.updateMyRankInfo(response, uid, weekStart, weekEnd)
			return response, nil
		}
	}

	response, err = s.buildLeaderboardResponse(uid, weekStart, weekEnd)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "获取热榜数据失败")
	}

	cacheExpire := time.Now().UTC().Add(5 * time.Minute)
	response.CacheExpire = cacheExpire

	cacheData, err := json.Marshal(response)
	if err == nil {
		database.RedisClient.Set(ctx, cacheKey, cacheData, 5*time.Minute)
	}

	return response, nil
}

func (s *LeaderboardService) buildLeaderboardResponse(uid string, weekStart, weekEnd time.Time) (*models.LeaderboardResponse, error) {
	ctx := context.Background()
	topUsers, err := s.leaderboardRepo.GetWeeklyLeaderboard(ctx, weekStart, weekEnd)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "获取热榜数据失败")
	}
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
	myRank := s.getMyRankInfo(uid, weekStart, weekEnd, topEntries)
	response := &models.LeaderboardResponse{
		WeekStart:   weekStart,
		WeekEnd:     weekEnd,
		MyRank:      myRank,
		TopUsers:    topEntries,
		CacheExpire: time.Now().UTC().Add(5 * time.Minute),
	}
	return response, nil
}

func (s *LeaderboardService) getMyRankInfo(uid string, weekStart, weekEnd time.Time, topEntries []models.LeaderboardEntry) *models.LeaderboardEntry {
	ctx := context.Background()
	for _, entry := range topEntries {
		if entry.Uid == uid {
			return &entry
		}
	}
	userData, rank, err := s.leaderboardRepo.GetUserWeeklyRank(ctx, uid, weekStart, weekEnd)
	if err != nil || userData == nil {
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

func (s *LeaderboardService) getDefaultUserRankInfo(uid string) *models.LeaderboardEntry {
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

func (s *LeaderboardService) updateMyRankInfo(response *models.LeaderboardResponse, uid string, weekStart, weekEnd time.Time) {
	for _, entry := range response.TopUsers {
		if entry.Uid == uid {
			response.MyRank = &entry
			return
		}
	}
	myRank := s.getMyRankInfo(uid, weekStart, weekEnd, response.TopUsers)
	response.MyRank = myRank
}

// ClearCache 清除排行榜缓存
func (s *LeaderboardService) ClearCache() error {
	ctx := context.Background()
	weekStart, _ := models.GetCurrentWeekRange()
	cacheKey := fmt.Sprintf("leaderboard:weekly:%s", weekStart.Format("2006-01-02"))
	
	// 删除缓存
	err := database.RedisClient.Del(ctx, cacheKey).Err()
	if err != nil {
		return utils.NewAppError(utils.CodeDatabaseError, "清除缓存失败")
	}
	
	return nil
}