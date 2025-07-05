package database

import (
	"context"
	"gin-fataMorgana/models"
	"time"
)

// LeaderboardRepository 热榜仓库
type LeaderboardRepository struct {
	*BaseRepository
}

// NewLeaderboardRepository 创建热榜仓库实例
func NewLeaderboardRepository() *LeaderboardRepository {
	return &LeaderboardRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// WeeklyLeaderboardData 周热榜数据
type WeeklyLeaderboardData struct {
	Uid         string    `json:"uid"`
	Username    string    `json:"username"`
	CompletedAt time.Time `json:"completed_at"`
	OrderCount  int       `json:"order_count"`
	TotalAmount float64   `json:"total_amount"`
	TotalProfit float64   `json:"total_profit"`
}

// GetWeeklyLeaderboard 获取本周热榜数据（不使用窗口函数）
func (r *LeaderboardRepository) GetWeeklyLeaderboard(ctx context.Context, weekStart, weekEnd time.Time) ([]WeeklyLeaderboardData, error) {
	var results []WeeklyLeaderboardData

	// 使用简单的GROUP BY和ORDER BY，不使用窗口函数
	query := `
		SELECT 
			o.uid,
			u.username,
			MAX(o.updated_at) as completed_at,
			COUNT(*) as order_count,
			SUM(o.amount) as total_amount,
			SUM(o.profit_amount) as total_profit
		FROM orders o
		JOIN users u ON o.uid = u.uid
		WHERE o.status = ? 
		AND o.updated_at >= ? 
		AND o.updated_at <= ?
		GROUP BY o.uid, u.username
		ORDER BY 
			order_count DESC,
			total_amount DESC,
			completed_at ASC
		LIMIT 10
	`

	err := r.db.WithContext(ctx).Raw(query, "success", weekStart, weekEnd).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetUserWeeklyRank 获取用户本周排名（不使用窗口函数）
func (r *LeaderboardRepository) GetUserWeeklyRank(ctx context.Context, uid string, weekStart, weekEnd time.Time) (*WeeklyLeaderboardData, int, error) {
	var userData WeeklyLeaderboardData

	// 获取用户数据
	userQuery := `
		SELECT 
			o.uid,
			u.username,
			MAX(o.updated_at) as completed_at,
			COUNT(*) as order_count,
			SUM(o.amount) as total_amount,
			SUM(o.profit_amount) as total_profit
		FROM orders o
		JOIN users u ON o.uid = u.uid
		WHERE o.status = ? 
		AND o.uid = ?
		AND o.updated_at >= ? 
		AND o.updated_at <= ?
		GROUP BY o.uid, u.username
	`

	err := r.db.WithContext(ctx).Raw(userQuery, "success", uid, weekStart, weekEnd).Scan(&userData).Error
	if err != nil {
		return nil, 0, err
	}

	// 如果用户没有完成订单，返回nil
	if userData.Uid == "" {
		return nil, 0, nil
	}

	// 计算用户排名：统计比当前用户成绩更好的用户数量
	rankQuery := `
		SELECT COUNT(*) + 1 as rank
		FROM (
			SELECT 
				o.uid,
				COUNT(*) as order_count,
				SUM(o.amount) as total_amount,
				MAX(o.updated_at) as completed_at
			FROM orders o
			WHERE o.status = ? 
			AND o.updated_at >= ? 
			AND o.updated_at <= ?
			GROUP BY o.uid
			HAVING 
				order_count > ? 
				OR (order_count = ? AND total_amount > ?)
				OR (order_count = ? AND total_amount = ? AND completed_at < ?)
		) as better_users
	`

	var rank int
	err = r.db.WithContext(ctx).Raw(rankQuery,
		"success", weekStart, weekEnd,
		userData.OrderCount,
		userData.OrderCount, userData.TotalAmount,
		userData.OrderCount, userData.TotalAmount, userData.CompletedAt).Scan(&rank).Error

	if err != nil {
		return nil, 0, err
	}

	return &userData, rank, nil
}

// GetUserByUid 根据UID获取用户信息
func (r *LeaderboardRepository) GetUserByUid(ctx context.Context, uid string) (*models.User, error) {
	var user models.User
	err := r.FindByCondition(ctx, map[string]interface{}{"uid": uid}, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
