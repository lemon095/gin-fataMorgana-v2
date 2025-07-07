package utils

import "fmt"

/*
Redis Key 管理工具

本文件统一管理项目中所有Redis Key的生成规则，确保Key的命名规范和一致性。

Key命名规范：
1. 使用冒号(:)分隔不同层级
2. 格式：{业务前缀}:{子类型}:{具体标识}
3. 示例：wallet:balance:user123

Key分类：
- wallet: 钱包相关（余额、锁、空值缓存）
- order: 订单相关（缓存、锁）
- user: 用户相关（配置、设置）
- email: 邮箱存在检查
- username: 用户名存在检查
- invite_code: 邀请码存在检查
- leaderboard: 排行榜相关
- group_buy: 拼单相关
- announcement: 公告相关
- config: 配置相关
- session: 会话相关
- rate_limit: 限流相关
- global: 全局系统相关

使用示例：
```go
// 生成钱包锁Key
lockKey := utils.RedisKeys.GenerateWalletLockKey("user123")
// 结果: wallet:lock:user123

// 生成订单缓存Key
cacheKey := utils.RedisKeys.GenerateOrderCacheKey("order_12345")
// 结果: order:cache:order_12345

// 生成限流Key
rateLimitKey := utils.RedisKeys.GenerateRateLimitKey("192.168.1.1", "1m")
// 结果: rate_limit:192.168.1.1:1m
```
*/

// Redis Key 前缀常量
const (
	// 钱包相关前缀
	// 示例: wallet:balance:user123, wallet:lock:user123, wallet:empty:user123
	WALLET_PREFIX = "wallet"
	
	// 订单相关前缀
	// 示例: order:cache:order_12345, order:lock:order_12345
	ORDER_PREFIX = "order"
	
	// 用户相关前缀
	// 示例: user:profile:user123, user:settings:user123
	USER_PREFIX = "user"
	
	// 邮箱相关前缀
	// 示例: email:test@example.com:exists
	EMAIL_PREFIX = "email"
	
	// 用户名相关前缀
	// 示例: username:john_doe:exists
	USERNAME_PREFIX = "username"
	
	// 邀请码相关前缀
	// 示例: invite_code:ABC123:exists
	INVITE_CODE_PREFIX = "invite_code"
	
	// 排行榜相关前缀
	// 示例: leaderboard:daily, leaderboard:weekly, leaderboard:lock:daily
	LEADERBOARD_PREFIX = "leaderboard"
	
	// 拼单相关前缀
	// 示例: group_buy:cache:group_123, group_buy:lock:group_123
	GROUP_BUY_PREFIX = "group_buy"
	
	// 公告相关前缀
	// 示例: announcement:cache:banner, announcement:cache:notice
	ANNOUNCEMENT_PREFIX = "announcement"
	
	// 配置相关前缀
	// 示例: config:cache:amount_config, config:cache:member_level
	CONFIG_PREFIX = "config"
	
	// 会话相关前缀
	// 示例: session:abc123def456
	SESSION_PREFIX = "session"
	
	// 限流相关前缀
	// 示例: rate_limit:192.168.1.1:1m, rate_limit:user123:1h
	RATE_LIMIT_PREFIX = "rate_limit"
)

// RedisKeyManager Redis Key管理器
type RedisKeyManager struct{}

// NewRedisKeyManager 创建Redis Key管理器
func NewRedisKeyManager() *RedisKeyManager {
	return &RedisKeyManager{}
}

// 钱包相关Key生成方法

// GenerateWalletLockKey 生成钱包分布式锁Key
// 示例: wallet:lock:user123
// 用途: 用于跨进程的钱包余额操作锁，确保同一用户的余额操作串行化
func (r *RedisKeyManager) GenerateWalletLockKey(uid string) string {
	return fmt.Sprintf("%s:lock:%s", WALLET_PREFIX, uid)
}

// GenerateWalletBalanceKey 生成钱包余额缓存Key
// 示例: wallet:balance:user123
// 用途: 缓存用户钱包余额数据，不过期，通过事件驱动更新
func (r *RedisKeyManager) GenerateWalletBalanceKey(uid string) string {
	return fmt.Sprintf("%s:balance:%s", WALLET_PREFIX, uid)
}

// GenerateWalletVersionKey 生成钱包版本号Key
// 示例: wallet:version:user123
// 用途: 缓存钱包数据版本号，用于数据一致性控制
func (r *RedisKeyManager) GenerateWalletVersionKey(uid string) string {
	return fmt.Sprintf("%s:version:%s", WALLET_PREFIX, uid)
}

// GenerateWalletEmptyKey 生成钱包空值缓存Key
// 示例: wallet:empty:user123
// 用途: 缓存空值，防止缓存穿透，过期时间10分钟
func (r *RedisKeyManager) GenerateWalletEmptyKey(uid string) string {
	return fmt.Sprintf("%s:empty:%s", WALLET_PREFIX, uid)
}

// 订单相关Key生成方法

// GenerateOrderCacheKey 生成订单缓存Key
// 示例: order:cache:order_12345
// 用途: 缓存订单数据，提高查询性能
func (r *RedisKeyManager) GenerateOrderCacheKey(orderID string) string {
	return fmt.Sprintf("%s:cache:%s", ORDER_PREFIX, orderID)
}

// GenerateOrderLockKey 生成订单锁Key
// 示例: order:lock:order_12345
// 用途: 订单操作的分布式锁，防止并发操作冲突
func (r *RedisKeyManager) GenerateOrderLockKey(orderID string) string {
	return fmt.Sprintf("%s:lock:%s", ORDER_PREFIX, orderID)
}

// 用户相关Key生成方法

// GenerateUserLoginTimeKey 生成用户登录时间Key
// 示例: user:login_time:user123
// 用途: 记录用户最后登录时间，用于钱包缓存过期策略
func (r *RedisKeyManager) GenerateUserLoginTimeKey(uid string) string {
	return fmt.Sprintf("%s:login_time:%s", USER_PREFIX, uid)
}

// GenerateEmailExistsKey 生成邮箱存在检查缓存Key
// 示例: email:test@example.com:exists
// 用途: 缓存邮箱是否存在的检查结果，避免重复查询数据库
func (r *RedisKeyManager) GenerateEmailExistsKey(email string) string {
	return fmt.Sprintf("%s:%s:exists", EMAIL_PREFIX, email)
}

// GenerateUsernameExistsKey 生成用户名存在检查缓存Key
// 示例: username:john_doe:exists
// 用途: 缓存用户名是否存在的检查结果，避免重复查询数据库
func (r *RedisKeyManager) GenerateUsernameExistsKey(username string) string {
	return fmt.Sprintf("%s:%s:exists", USERNAME_PREFIX, username)
}

// GenerateInviteCodeExistsKey 生成邀请码存在检查缓存Key
// 示例: invite_code:ABC123:exists
// 用途: 缓存邀请码是否存在的检查结果，避免重复查询数据库
func (r *RedisKeyManager) GenerateInviteCodeExistsKey(inviteCode string) string {
	return fmt.Sprintf("%s:%s:exists", INVITE_CODE_PREFIX, inviteCode)
}

// 排行榜相关Key生成方法

// GenerateLeaderboardKey 生成排行榜Key
// 示例: leaderboard:daily, leaderboard:weekly, leaderboard:monthly
// 用途: 缓存排行榜数据，使用Redis的有序集合存储
func (r *RedisKeyManager) GenerateLeaderboardKey(leaderboardType string) string {
	return fmt.Sprintf("%s:%s", LEADERBOARD_PREFIX, leaderboardType)
}

// GenerateLeaderboardLockKey 生成排行榜锁Key
// 示例: leaderboard:lock:daily, leaderboard:lock:weekly
// 用途: 排行榜更新的分布式锁，防止并发更新冲突
func (r *RedisKeyManager) GenerateLeaderboardLockKey(leaderboardType string) string {
	return fmt.Sprintf("%s:lock:%s", LEADERBOARD_PREFIX, leaderboardType)
}

// 拼单相关Key生成方法

// GenerateGroupBuyCacheKey 生成拼单缓存Key
// 示例: group_buy:cache:group_123, group_buy:cache:group_456
// 用途: 缓存拼单数据，提高查询性能
func (r *RedisKeyManager) GenerateGroupBuyCacheKey(groupBuyID string) string {
	return fmt.Sprintf("%s:cache:%s", GROUP_BUY_PREFIX, groupBuyID)
}

// GenerateGroupBuyLockKey 生成拼单锁Key
// 示例: group_buy:lock:group_123, group_buy:lock:group_456
// 用途: 拼单操作的分布式锁，防止并发操作冲突
func (r *RedisKeyManager) GenerateGroupBuyLockKey(groupBuyID string) string {
	return fmt.Sprintf("%s:lock:%s", GROUP_BUY_PREFIX, groupBuyID)
}

// 公告相关Key生成方法

// GenerateAnnouncementCacheKey 生成公告缓存Key
// 示例: announcement:cache:banner, announcement:cache:notice, announcement:cache:popup
// 用途: 缓存公告数据，减少数据库查询
func (r *RedisKeyManager) GenerateAnnouncementCacheKey(announcementType string) string {
	return fmt.Sprintf("%s:cache:%s", ANNOUNCEMENT_PREFIX, announcementType)
}

// 配置相关Key生成方法

// GenerateConfigCacheKey 生成配置缓存Key
// 示例: config:cache:amount_config, config:cache:member_level, config:cache:system_config
// 用途: 缓存系统配置数据，避免频繁查询数据库
func (r *RedisKeyManager) GenerateConfigCacheKey(configType string) string {
	return fmt.Sprintf("%s:cache:%s", CONFIG_PREFIX, configType)
}

// 会话相关Key生成方法

// GenerateSessionKey 生成会话Key
// 示例: session:abc123def456, session:xyz789uvw012
// 用途: 存储用户会话信息，支持会话过期和自动清理
func (r *RedisKeyManager) GenerateSessionKey(sessionID string) string {
	return fmt.Sprintf("%s:%s", SESSION_PREFIX, sessionID)
}

// 限流相关Key生成方法

// GenerateRateLimitKey 生成限流Key
// 示例: rate_limit:192.168.1.1:1m, rate_limit:user123:1h, rate_limit:api_login:1d
// 用途: 实现API限流，支持按IP、用户ID、API接口等维度限流
func (r *RedisKeyManager) GenerateRateLimitKey(identifier string, window string) string {
	return fmt.Sprintf("%s:%s:%s", RATE_LIMIT_PREFIX, identifier, window)
}

// 全局Key生成方法（不分类）

// GenerateGlobalLockKey 生成全局锁Key
// 示例: global:lock:system_maintenance, global:lock:data_migration
// 用途: 全局系统锁，用于系统维护、数据迁移等全局操作
func (r *RedisKeyManager) GenerateGlobalLockKey(lockName string) string {
	return fmt.Sprintf("global:lock:%s", lockName)
}

// GenerateGlobalCacheKey 生成全局缓存Key
// 示例: global:cache:system_status, global:cache:app_version
// 用途: 全局系统缓存，存储系统状态、应用版本等全局信息
func (r *RedisKeyManager) GenerateGlobalCacheKey(cacheName string) string {
	return fmt.Sprintf("global:cache:%s", cacheName)
}

// GenerateGlobalCounterKey 生成全局计数器Key
// 示例: global:counter:total_users, global:counter:daily_orders
// 用途: 全局计数器，统计用户总数、日订单数等全局数据
func (r *RedisKeyManager) GenerateGlobalCounterKey(counterName string) string {
	return fmt.Sprintf("global:counter:%s", counterName)
}

// 工具方法

// GetKeyPattern 获取Key模式（用于模糊匹配）
// 示例: wallet:*, order:*, user:*
// 用途: 用于Redis的SCAN命令，批量操作或清理特定前缀的Key
func (r *RedisKeyManager) GetKeyPattern(prefix string) string {
	return fmt.Sprintf("%s:*", prefix)
}

// GetWalletKeyPattern 获取钱包相关Key模式
// 示例: wallet:*
// 用途: 匹配所有钱包相关的Key，用于批量清理钱包缓存
func (r *RedisKeyManager) GetWalletKeyPattern() string {
	return r.GetKeyPattern(WALLET_PREFIX)
}

// GetOrderKeyPattern 获取订单相关Key模式
// 示例: order:*
// 用途: 匹配所有订单相关的Key，用于批量清理订单缓存
func (r *RedisKeyManager) GetOrderKeyPattern() string {
	return r.GetKeyPattern(ORDER_PREFIX)
}

// GetUserKeyPattern 获取用户相关Key模式
// 示例: user:*
// 用途: 匹配所有用户相关的Key，用于批量清理用户缓存
func (r *RedisKeyManager) GetUserKeyPattern() string {
	return r.GetKeyPattern(USER_PREFIX)
}

// 全局实例
var RedisKeys = NewRedisKeyManager() 