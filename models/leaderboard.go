package models

import (
	"time"
)

// LeaderboardEntry 热榜条目
type LeaderboardEntry struct {
	ID          uint      `json:"id"`
	Uid         string    `json:"uid"`
	Username    string    `json:"username"`     // 脱敏后的用户名
	CompletedAt time.Time `json:"completed_at"` // 最新完成时间
	OrderCount  int       `json:"order_count"`  // 完成订单数量
	TotalAmount float64   `json:"total_amount"` // 总金额
	TotalProfit float64   `json:"total_profit"` // 总利润
	Rank        int       `json:"rank"`         // 排名
	IsRank      bool      `json:"is_rank"`      // 是否在榜单上
}

// LeaderboardResponse 热榜响应
type LeaderboardResponse struct {
	WeekStart   time.Time          `json:"week_start"`   // 本周开始时间
	WeekEnd     time.Time          `json:"week_end"`     // 本周结束时间
	MyRank      *LeaderboardEntry  `json:"my_rank"`      // 我的排名信息
	TopUsers    []LeaderboardEntry `json:"top_users"`    // 前10名用户
	CacheExpire time.Time          `json:"cache_expire"` // 缓存过期时间
}

// GetWeekStart 获取本周开始时间（周一）
func GetWeekStart(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	} else {
		weekday = weekday - 1
	}

	// 计算到本周一的天数
	daysToMonday := int(weekday)

	// 获取本周一的时间
	weekStart := t.AddDate(0, 0, -daysToMonday)
	return time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
}

// GetWeekEnd 获取本周结束时间（周日）
func GetWeekEnd(weekStart time.Time) time.Time {
	return weekStart.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
}

// GetCurrentWeekRange 获取当前周的时间范围
func GetCurrentWeekRange() (time.Time, time.Time) {
	now := time.Now()
	weekStart := GetWeekStart(now)
	weekEnd := GetWeekEnd(weekStart)
	return weekStart, weekEnd
}

// MaskUsername 对用户名进行脱敏处理（优化版本）
func MaskUsername(username string) string {
	if len(username) <= 1 {
		return username // 如果用户名太短（只有1个字符），直接返回
	}
	
	// 对于2个字符的用户名，在中间加*
	if len(username) == 2 {
		return username[:1] + "*" + username[1:]
	}
	
	// 对于3-4个字符的用户名，只显示首尾
	if len(username) <= 4 {
		first := username[:1]
		last := username[len(username)-1:]
		return first + "*" + last
	}
	
	// 对于5个字符及以上的用户名，显示首尾各2个字符
	first := username[:2]
	last := username[len(username)-2:]
	return first + "*" + last
}

// LeaderboardQuery 热榜查询参数
type LeaderboardQuery struct {
	Uid string `json:"uid" binding:"required"` // 当前用户ID
}
