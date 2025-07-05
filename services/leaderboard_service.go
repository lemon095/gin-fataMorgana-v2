package services

import (
	"context"
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
	weekStart, weekEnd := models.GetCurrentWeekRange()

	response, err := s.buildLeaderboardResponse(uid, weekStart, weekEnd)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "获取热榜数据失败")
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
		WeekStart: weekStart,
		WeekEnd:   weekEnd,
		MyRank:    myRank,
		TopUsers:  topEntries,
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