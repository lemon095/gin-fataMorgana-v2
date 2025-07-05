package database

import (
	"context"
	"gin-fataMorgana/models"
	"log"
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

	log.Printf("🔍 [数据库] 开始查询排行榜数据")
	log.Printf("🔍 [数据库] 时间范围: %s 到 %s", weekStart.Format("2006-01-02 15:04:05"), weekEnd.Format("2006-01-02 15:04:05"))

	// 使用LEFT JOIN，为系统订单生成虚拟用户名
	query := `
		SELECT 
			o.uid,
			CASE 
				WHEN o.is_system_order = 1 THEN CONCAT(
					CHAR(65 + (ASCII(SUBSTRING(o.uid, 1, 1)) % 26)),
					CHAR(65 + (ASCII(SUBSTRING(o.uid, 2, 1)) % 26)),
					'**',
					CHAR(65 + (ASCII(SUBSTRING(o.uid, -2, 1)) % 26)),
					CHAR(65 + (ASCII(SUBSTRING(o.uid, -1, 1)) % 26))
				)
				ELSE COALESCE(u.username, '未知用户')
			END as username,
			MAX(o.updated_at) as completed_at,
			COUNT(*) as order_count,
			SUM(o.amount) as total_amount,
			SUM(o.profit_amount) as total_profit
		FROM orders o
		LEFT JOIN users u ON o.uid = u.uid
		WHERE o.status = ? 
		AND o.updated_at >= ? 
		AND o.updated_at <= ?
		GROUP BY o.uid, 
			CASE 
				WHEN o.is_system_order = 1 THEN CONCAT(
					CHAR(65 + (ASCII(SUBSTRING(o.uid, 1, 1)) % 26)),
					CHAR(65 + (ASCII(SUBSTRING(o.uid, 2, 1)) % 26)),
					'**',
					CHAR(65 + (ASCII(SUBSTRING(o.uid, -2, 1)) % 26)),
					CHAR(65 + (ASCII(SUBSTRING(o.uid, -1, 1)) % 26))
				)
				ELSE COALESCE(u.username, '未知用户')
			END
		ORDER BY 
			order_count DESC,
			total_amount DESC,
			completed_at ASC
		LIMIT 10;
	`

	log.Printf("🔍 [数据库] 执行SQL查询: %s", query)
	log.Printf("🔍 [数据库] 查询参数: status=success, weekStart=%s, weekEnd=%s", 
		weekStart.Format("2006-01-02 15:04:05"), weekEnd.Format("2006-01-02 15:04:05"))

	err := r.db.WithContext(ctx).Raw(query, "success", weekStart, weekEnd).Scan(&results).Error
	if err != nil {
		log.Printf("❌ [数据库] 查询排行榜数据失败: %v", err)
		return nil, err
	}

	log.Printf("✅ [数据库] 查询成功，返回 %d 条记录", len(results))
	
	// 输出查询结果的详细信息
	for i, result := range results {
		log.Printf("📊 [数据库] 结果%d: UID=%s, 用户名=%s, 订单数=%d, 总金额=%.2f, 总利润=%.2f", 
			i+1, result.Uid, result.Username, result.OrderCount, result.TotalAmount, result.TotalProfit)
	}

	return results, nil
}

// GetUserWeeklyRank 获取用户本周排名（不使用窗口函数）
func (r *LeaderboardRepository) GetUserWeeklyRank(ctx context.Context, uid string, weekStart, weekEnd time.Time) (*WeeklyLeaderboardData, int, error) {
	var userData WeeklyLeaderboardData

	log.Printf("🔍 [数据库] 开始查询用户 %s 的排名", uid)
	log.Printf("🔍 [数据库] 时间范围: %s 到 %s", weekStart.Format("2006-01-02 15:04:05"), weekEnd.Format("2006-01-02 15:04:05"))

	// 获取用户数据
	userQuery := `
		SELECT 
			o.uid,
			CASE 
				WHEN o.is_system_order = 1 THEN CONCAT(
					CHAR(65 + (ASCII(SUBSTRING(o.uid, 1, 1)) % 26)),
					CHAR(65 + (ASCII(SUBSTRING(o.uid, 2, 1)) % 26)),
					'**',
					CHAR(65 + (ASCII(SUBSTRING(o.uid, -2, 1)) % 26)),
					CHAR(65 + (ASCII(SUBSTRING(o.uid, -1, 1)) % 26))
				)
				ELSE COALESCE(u.username, '未知用户')
			END as username,
			MAX(o.updated_at) as completed_at,
			COUNT(*) as order_count,
			SUM(o.amount) as total_amount,
			SUM(o.profit_amount) as total_profit
		FROM orders o
		LEFT JOIN users u ON o.uid = u.uid
		WHERE o.status = ? 
		AND o.uid = ?
		AND o.updated_at >= ? 
		AND o.updated_at <= ?
		GROUP BY o.uid, 
			CASE 
				WHEN o.is_system_order = 1 THEN CONCAT(
					CHAR(65 + (ASCII(SUBSTRING(o.uid, 1, 1)) % 26)),
					CHAR(65 + (ASCII(SUBSTRING(o.uid, 2, 1)) % 26)),
					'**',
					CHAR(65 + (ASCII(SUBSTRING(o.uid, -2, 1)) % 26)),
					CHAR(65 + (ASCII(SUBSTRING(o.uid, -1, 1)) % 26))
				)
				ELSE COALESCE(u.username, '未知用户')
			END
	`

	log.Printf("🔍 [数据库] 执行用户数据查询: %s", userQuery)
	log.Printf("🔍 [数据库] 用户查询参数: status=success, uid=%s, weekStart=%s, weekEnd=%s", 
		uid, weekStart.Format("2006-01-02 15:04:05"), weekEnd.Format("2006-01-02 15:04:05"))

	err := r.db.WithContext(ctx).Raw(userQuery, "success", uid, weekStart, weekEnd).Scan(&userData).Error
	if err != nil {
		log.Printf("❌ [数据库] 查询用户数据失败: %v", err)
		return nil, 0, err
	}

	// 如果用户没有完成订单，返回nil
	if userData.Uid == "" {
		log.Printf("⚠️ [数据库] 用户 %s 没有找到任何订单数据", uid)
		return nil, 0, nil
	}

	log.Printf("✅ [数据库] 用户 %s 数据: 订单数=%d, 总金额=%.2f, 总利润=%.2f", 
		uid, userData.OrderCount, userData.TotalAmount, userData.TotalProfit)

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

	log.Printf("🔍 [数据库] 执行排名计算查询: %s", rankQuery)
	log.Printf("🔍 [数据库] 排名查询参数: status=success, weekStart=%s, weekEnd=%s, orderCount=%d, totalAmount=%.2f, completedAt=%s", 
		weekStart.Format("2006-01-02 15:04:05"), weekEnd.Format("2006-01-02 15:04:05"), 
		userData.OrderCount, userData.TotalAmount, userData.CompletedAt.Format("2006-01-02 15:04:05"))

	var rank int
	err = r.db.WithContext(ctx).Raw(rankQuery,
		"success", weekStart, weekEnd,
		userData.OrderCount,
		userData.OrderCount, userData.TotalAmount,
		userData.OrderCount, userData.TotalAmount, userData.CompletedAt).Scan(&rank).Error

	if err != nil {
		log.Printf("❌ [数据库] 计算用户排名失败: %v", err)
		return nil, 0, err
	}

	log.Printf("✅ [数据库] 用户 %s 排名计算完成: 第%d名", uid, rank)

	return &userData, rank, nil
}

// GetUserByUid 根据UID获取用户信息
func (r *LeaderboardRepository) GetUserByUid(ctx context.Context, uid string) (*models.User, error) {
	var user models.User
	err := r.FindByCondition(ctx, map[string]interface{}{"uid": uid}, &user)
	if err != nil {
		log.Printf("❌ [数据库] 查询用户信息失败: %v", err)
		return nil, err
	}
	log.Printf("✅ [数据库] 查询用户信息成功: UID=%s, 用户名=%s", uid, user.Username)
	return &user, nil
}
