package models

// LeaderboardEntry 热榜条目
type LeaderboardEntry struct {
	Rank        int     `json:"rank"`         // 排名
	UserID      int64   `json:"user_id"`      // 用户ID
	Username    string  `json:"username"`     // 用户名
	Amount      float64 `json:"amount"`       // 金额
	OrderCount  int     `json:"order_count"`  // 完成单数
	Profit      float64 `json:"profit"`       // 利润金额
	CreatedAt   string  `json:"created_at"`   // 时间
}

// LeaderboardResponse 热榜响应
type LeaderboardResponse struct {
	RankingList []LeaderboardEntry `json:"ranking_list"` // 排行榜列表
	MyData      *MyLeaderboardData `json:"my_data"`      // 我的数据
}

// MyLeaderboardData 我的热榜数据
type MyLeaderboardData struct {
	Rank      int   `json:"rank"`       // 排名
	UserID    int64 `json:"user_id"`    // 用户ID
	IsRanked  bool  `json:"is_ranked"`  // 是否上榜
	Entry     *LeaderboardEntry `json:"entry,omitempty"` // 上榜时的详细信息
}

// LeaderboardRequest 热榜请求
type LeaderboardRequest struct {
	UserID int64 `json:"user_id" binding:"required"` // 当前用户ID
} 