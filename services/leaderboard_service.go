package services

import (
	"context"
	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"
	"log"
	"time"
)

type LeaderboardService struct {
	leaderboardRepo *database.LeaderboardRepository
	leaderboardCacheService *LeaderboardCacheService
}

func NewLeaderboardService() *LeaderboardService {
	return &LeaderboardService{
		leaderboardRepo: database.NewLeaderboardRepository(),
		leaderboardCacheService: NewLeaderboardCacheService(),
	}
}

func (s *LeaderboardService) GetLeaderboard(uid string) (*models.LeaderboardResponse, error) {
	log.Printf("🔍 [排行榜] 开始查询排行榜数据")
	log.Printf("🔍 [排行榜] 用户UID: %s", uid)
	log.Printf("🔍 [排行榜] 当前时间: %s", time.Now().Format("2006-01-02 15:04:05"))

	// 尝试从缓存获取数据
	cachedData, err := s.leaderboardCacheService.GetCachedLeaderboardData()
	if err != nil {
		log.Printf("⚠️ [排行榜] 缓存未命中，从数据库查询: %v", err)
		// 缓存未命中，从数据库查询
		weekStart, weekEnd := models.GetCurrentWeekRange()
		response, err := s.buildLeaderboardResponse(uid, weekStart, weekEnd)
		if err != nil {
			log.Printf("❌ [排行榜] 构建排行榜响应失败: %v", err)
		return nil, utils.NewAppError(utils.CodeDatabaseError, "获取热榜数据失败")
	}
		log.Printf("✅ [排行榜] 从数据库查询完成，返回 %d 条数据", len(response.TopUsers))
		return response, nil
	}

	// 从缓存获取数据成功
	log.Printf("✅ [排行榜] 从缓存获取数据成功，包含 %d 条记录", len(cachedData.TopUsers))
	
	// 获取用户排名信息
	myRank := s.leaderboardCacheService.GetUserRankFromCache(uid, cachedData)
	
	response := &models.LeaderboardResponse{
		WeekStart:  cachedData.WeekStart,
		WeekEnd:    cachedData.WeekEnd,
		MyRank:     myRank,
		TopUsers:   cachedData.TopUsers,
		NextUpdate: cachedData.NextUpdate,
	}

	log.Printf("✅ [排行榜] 从缓存查询完成，返回 %d 条数据，下次更新时间: %s", 
		len(response.TopUsers), cachedData.NextUpdate.Format("2006-01-02 15:04:05"))
	return response, nil
}

func (s *LeaderboardService) buildLeaderboardResponse(uid string, weekStart, weekEnd time.Time) (*models.LeaderboardResponse, error) {
	ctx := context.Background()
	
	log.Printf("🔍 [排行榜] 开始查询前10名用户数据")
	topUsers, err := s.leaderboardRepo.GetWeeklyLeaderboard(ctx, weekStart, weekEnd)
	if err != nil {
		log.Printf("❌ [排行榜] 查询前10名用户失败: %v", err)
		return nil, utils.NewAppError(utils.CodeDatabaseError, "获取热榜数据失败")
	}
	
	log.Printf("✅ [排行榜] 查询到 %d 个用户的数据", len(topUsers))
	
	// 输出每个用户的详细信息
	for i, user := range topUsers {
		log.Printf("📊 [排行榜] 第%d名: UID=%s, 用户名=%s, 订单数=%d, 总金额=%.2f, 总利润=%.2f, 完成时间=%s", 
			i+1, user.Uid, user.Username, user.OrderCount, user.TotalAmount, user.TotalProfit, 
			user.CompletedAt.Format("2006-01-02 15:04:05"))
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
	
	log.Printf("🔍 [排行榜] 开始查询用户 %s 的排名信息", uid)
	myRank := s.getMyRankInfo(uid, weekStart, weekEnd, topEntries)
	
	response := &models.LeaderboardResponse{
		WeekStart: weekStart,
		WeekEnd:   weekEnd,
		MyRank:    myRank,
		TopUsers:  topEntries,
	}
	
	log.Printf("✅ [排行榜] 排行榜响应构建完成")
	return response, nil
}

func (s *LeaderboardService) getMyRankInfo(uid string, weekStart, weekEnd time.Time, topEntries []models.LeaderboardEntry) *models.LeaderboardEntry {
	ctx := context.Background()
	
	// 先检查是否在前10名中
	for _, entry := range topEntries {
		if entry.Uid == uid {
			log.Printf("✅ [排行榜] 用户 %s 在前10名中，排名第%d", uid, entry.Rank)
			return &entry
		}
	}
	
	log.Printf("🔍 [排行榜] 用户 %s 不在前10名中，查询具体排名", uid)
	userData, rank, err := s.leaderboardRepo.GetUserWeeklyRank(ctx, uid, weekStart, weekEnd)
	if err != nil {
		log.Printf("❌ [排行榜] 查询用户排名失败: %v", err)
		return s.getDefaultUserRankInfo(uid)
	}
	
	if userData == nil {
		log.Printf("⚠️ [排行榜] 用户 %s 没有完成任何订单", uid)
		return s.getDefaultUserRankInfo(uid)
	}
	
	log.Printf("✅ [排行榜] 用户 %s 排名第%d，完成订单数=%d，总金额=%.2f", 
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

func (s *LeaderboardService) getDefaultUserRankInfo(uid string) *models.LeaderboardEntry {
	log.Printf("🔍 [排行榜] 获取用户 %s 的默认排名信息", uid)
	userRepo := database.NewUserRepository()
	user, err := userRepo.FindByUid(context.Background(), uid)
	if err != nil {
		log.Printf("❌ [排行榜] 查询用户信息失败: %v", err)
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
	
	log.Printf("✅ [排行榜] 用户 %s 默认排名信息：用户名=%s", uid, user.Username)
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