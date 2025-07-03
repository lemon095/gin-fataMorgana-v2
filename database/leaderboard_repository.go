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

// GetWeeklyLeaderboard 获取本周热榜数据
func (r *LeaderboardRepository) GetWeeklyLeaderboard(ctx context.Context, weekStart, weekEnd time.Time) ([]WeeklyLeaderboardData, error) {
	var results []WeeklyLeaderboardData

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
		WHERE o.status = 'success' 
		AND o.updated_at >= ? 
		AND o.updated_at <= ?
		GROUP BY o.uid, u.username
		ORDER BY order_count DESC, total_amount DESC, completed_at ASC
		LIMIT 10
	`

	err := r.db.WithContext(ctx).Raw(query, weekStart, weekEnd).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetUserWeeklyRank 获取用户本周排名
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
		WHERE o.status = 'success' 
		AND o.uid = ?
		AND o.updated_at >= ? 
		AND o.updated_at <= ?
		GROUP BY o.uid, u.username
	`

	err := r.db.WithContext(ctx).Raw(userQuery, uid, weekStart, weekEnd).Scan(&userData).Error
	if err != nil {
		return nil, 0, err
	}

	// 如果用户没有完成订单，返回nil
	if userData.Uid == "" {
		return nil, 0, nil
	}

	// 获取用户排名
	rankQuery := `
		SELECT COUNT(*) + 1 as rank
		FROM (
			SELECT 
				o.uid,
				COUNT(*) as order_count,
				SUM(o.amount) as total_amount,
				MAX(o.updated_at) as completed_at
			FROM orders o
			WHERE o.status = 'success' 
			AND o.updated_at >= ? 
			AND o.updated_at <= ?
			GROUP BY o.uid
			HAVING (
				order_count > (SELECT COUNT(*) FROM orders WHERE status = 'success' AND uid = ? AND updated_at >= ? AND updated_at <= ?)
				OR (
					order_count = (SELECT COUNT(*) FROM orders WHERE status = 'success' AND uid = ? AND updated_at >= ? AND updated_at <= ?)
					AND total_amount > (SELECT SUM(amount) FROM orders WHERE status = 'success' AND uid = ? AND updated_at >= ? AND updated_at <= ?)
				)
				OR (
					order_count = (SELECT COUNT(*) FROM orders WHERE status = 'success' AND uid = ? AND updated_at >= ? AND updated_at <= ?)
					AND total_amount = (SELECT SUM(amount) FROM orders WHERE status = 'success' AND uid = ? AND updated_at >= ? AND updated_at <= ?)
					AND completed_at < (SELECT MAX(updated_at) FROM orders WHERE status = 'success' AND uid = ? AND updated_at >= ? AND updated_at <= ?)
				)
			)
		) as ranking
	`

	var rank int
	err = r.db.WithContext(ctx).Raw(rankQuery,
		weekStart, weekEnd, uid, weekStart, weekEnd,
		uid, weekStart, weekEnd, uid, weekStart, weekEnd,
		uid, weekStart, weekEnd, uid, weekStart, weekEnd,
		uid, weekStart, weekEnd).Scan(&rank).Error

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
