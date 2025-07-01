package services

import (
	"gin-fataMorgana/models"
)

// LeaderboardService 热榜服务
type LeaderboardService struct{}

// NewLeaderboardService 创建热榜服务实例
func NewLeaderboardService() *LeaderboardService {
	return &LeaderboardService{}
}

// GetLeaderboard 获取任务热榜
func (s *LeaderboardService) GetLeaderboard(userID int64) (*models.LeaderboardResponse, error) {
	// 生成固定假数据
	rankingList := s.generateMockData()
	
	// 查找当前用户的数据
	myData := s.findMyData(userID, rankingList)
	
	return &models.LeaderboardResponse{
		RankingList: rankingList,
		MyData:      myData,
	}, nil
}

// generateMockData 生成固定假数据
func (s *LeaderboardService) generateMockData() []models.LeaderboardEntry {
	// 创建固定假数据
	mockData := []models.LeaderboardEntry{
		{
			Rank:       1,
			UserID:     1001,
			Username:   "张三",
			Amount:     15800.50,
			OrderCount: 156,
			Profit:     3200.80,
			CreatedAt:  "2024-01-15 14:30:00",
		},
		{
			Rank:       2,
			UserID:     1002,
			Username:   "李四",
			Amount:     14200.30,
			OrderCount: 142,
			Profit:     2850.60,
			CreatedAt:  "2024-01-15 13:45:00",
		},
		{
			Rank:       3,
			UserID:     1003,
			Username:   "王五",
			Amount:     12800.75,
			OrderCount: 128,
			Profit:     2560.15,
			CreatedAt:  "2024-01-15 12:20:00",
		},
		{
			Rank:       4,
			UserID:     1004,
			Username:   "赵六",
			Amount:     11500.20,
			OrderCount: 115,
			Profit:     2300.40,
			CreatedAt:  "2024-01-15 11:15:00",
		},
		{
			Rank:       5,
			UserID:     1005,
			Username:   "钱七",
			Amount:     10200.80,
			OrderCount: 102,
			Profit:     2040.16,
			CreatedAt:  "2024-01-15 10:30:00",
		},
		{
			Rank:       6,
			UserID:     1006,
			Username:   "孙八",
			Amount:     8900.45,
			OrderCount: 89,
			Profit:     1780.09,
			CreatedAt:  "2024-01-15 09:45:00",
		},
		{
			Rank:       7,
			UserID:     1007,
			Username:   "周九",
			Amount:     7600.30,
			OrderCount: 76,
			Profit:     1520.06,
			CreatedAt:  "2024-01-15 08:20:00",
		},
		{
			Rank:       8,
			UserID:     1008,
			Username:   "吴十",
			Amount:     6500.60,
			OrderCount: 65,
			Profit:     1300.12,
			CreatedAt:  "2024-01-15 07:30:00",
		},
		{
			Rank:       9,
			UserID:     1009,
			Username:   "郑十一",
			Amount:     5400.25,
			OrderCount: 54,
			Profit:     1080.05,
			CreatedAt:  "2024-01-15 06:15:00",
		},
		{
			Rank:       10,
			UserID:     1010,
			Username:   "王十二",
			Amount:     4300.90,
			OrderCount: 43,
			Profit:     860.18,
			CreatedAt:  "2024-01-15 05:00:00",
		},
	}
	
	return mockData
}

// findMyData 查找当前用户的数据
func (s *LeaderboardService) findMyData(userID int64, rankingList []models.LeaderboardEntry) *models.MyLeaderboardData {
	// 在排行榜中查找当前用户
	for _, entry := range rankingList {
		if entry.UserID == userID {
			return &models.MyLeaderboardData{
				Rank:     entry.Rank,
				UserID:   userID,
				IsRanked: true,
				Entry:    &entry,
			}
		}
	}
	
	// 如果没上榜，返回未上榜状态
	return &models.MyLeaderboardData{
		Rank:     0,
		UserID:   userID,
		IsRanked: false,
		Entry:    nil,
	}
} 