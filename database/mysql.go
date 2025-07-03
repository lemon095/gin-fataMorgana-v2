package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gin-fataMorgana/config"
	"gin-fataMorgana/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitMySQL 初始化MySQL数据库连接
func InitMySQL() error {
	cfg := config.GlobalConfig.Database

	// 获取数据库连接字符串
	dsn := cfg.GetDSN()

	// 配置GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层的sql.DB对象
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池参数（从配置文件读取）
	maxIdleConns := cfg.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = 10 // 默认值
	}

	maxOpenConns := cfg.MaxOpenConns
	if maxOpenConns == 0 {
		maxOpenConns = 100 // 默认值
	}

	connMaxLifetime := time.Duration(cfg.ConnMaxLifetime) * time.Second
	if connMaxLifetime == 0 {
		connMaxLifetime = time.Hour // 默认值
	}

	connMaxIdleTime := time.Duration(cfg.ConnMaxIdleTime) * time.Second
	if connMaxIdleTime == 0 {
		connMaxIdleTime = 30 * time.Minute // 默认值
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	DB = db
	log.Printf("MySQL数据库连接成功 - 连接池配置: 最大空闲=%d, 最大连接=%d, 连接最大生存时间=%v, 连接空闲超时=%v",
		maxIdleConns, maxOpenConns, connMaxLifetime, connMaxIdleTime)
	return nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 自动迁移表结构
	err := DB.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.WalletTransaction{},
		&models.AdminUser{},
		&models.UserLoginLog{},
		&models.Order{},
		&models.AmountConfig{},
		&models.Announcement{},
		&models.AnnouncementBanner{},
		&models.GroupBuy{},
		&models.MemberLevel{},
		&models.LotteryPeriod{},
	)
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 添加表注释
	if err := addTableComments(); err != nil {
		log.Printf("添加表注释失败: %v", err)
	}

	// 为拼单表添加注释
	sqlDB, err := DB.DB()
	if err == nil {
		_, _ = sqlDB.Exec("ALTER TABLE `group_buys` COMMENT = '拼单表 - 记录拼单信息，包括参与人数、付款金额、截止时间等'")
	}

	log.Println("数据库表迁移完成")
	return nil
}

// addTableComments 添加表注释
func addTableComments() error {
	// 获取数据库连接
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 表注释映射
	tableComments := map[string]string{
		"users":                "用户表 - 存储用户基本信息、认证信息、银行卡信息、经验值、信用分等",
		"wallets":              "钱包表 - 存储用户钱包信息，包括余额、冻结余额、总收入、总支出等",
		"wallet_transactions":  "钱包交易流水表 - 记录所有钱包交易明细，包括充值、提现、购买、拼单等操作",
		"user_login_logs":      "用户登录日志表 - 记录用户登录历史，包括登录时间、IP地址、设备信息、登录状态等",
		"admin_users":          "邀请码管理表 - 存储邀请码信息，用于用户注册时的邀请码校验，默认角色为业务员(4)",
		"amount_config":        "金额配置表 - 存储充值、提现等操作的金额配置，支持排序和激活状态管理",
		"announcements":        "公告表 - 存储系统公告信息，支持富文本内容，包括标题、纯文本内容、富文本内容、标签、状态等",
		"announcement_banners": "公告图片表 - 存储公告相关的图片信息，支持排序和跳转链接",
		"member_level":         "用户等级配置表 - 存储用户等级配置信息，包括等级、经验值范围、返现比例等",
		"lottery_periods":      "游戏期数表 - 记录每期的编号、订单金额、状态和时间信息",
	}

	// 为每个表添加注释
	for tableName, comment := range tableComments {
		query := fmt.Sprintf("ALTER TABLE `%s` COMMENT = '%s'", tableName, comment)
		if _, err := sqlDB.Exec(query); err != nil {
			log.Printf("为表 %s 添加注释失败: %v", tableName, err)
		} else {
			log.Printf("为表 %s 添加注释成功", tableName)
		}
	}

	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("获取数据库实例失败: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDBStats 获取数据库连接池统计信息
func GetDBStats() map[string]interface{} {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return nil
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

// HealthCheck 数据库健康检查
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	return sqlDB.Ping()
}

// Transaction 事务包装器
func Transaction(fn func(tx *gorm.DB) error) error {
	return DB.Transaction(fn)
}

// TransactionWithContext 带上下文的事务包装器
func TransactionWithContext(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return DB.WithContext(ctx).Transaction(fn)
}

// QueryCache 查询缓存中间件
type QueryCache struct {
	DB *gorm.DB
}

// NewQueryCache 创建查询缓存实例
func NewQueryCache() *QueryCache {
	return &QueryCache{DB: DB}
}

// WithCache 带缓存的查询
func (qc *QueryCache) WithCache(ctx context.Context, key string, expiration time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	// 先尝试从Redis获取缓存
	if cached, err := RedisClient.Get(ctx, key).Result(); err == nil {
		return cached, nil
	}

	// 缓存未命中，执行查询
	result, err := fn()
	if err != nil {
		return nil, err
	}

	// 缓存结果
	RedisClient.Set(ctx, key, result, expiration)
	return result, nil
}

// InvalidateCache 清除缓存
func (qc *QueryCache) InvalidateCache(ctx context.Context, pattern string) error {
	keys, err := RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return RedisClient.Del(ctx, keys...).Err()
	}
	return nil
}

// GetQueryStats 获取查询统计信息
func GetQueryStats() map[string]interface{} {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return nil
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections":   stats.MaxOpenConnections,
		"open_connections":       stats.OpenConnections,
		"in_use":                 stats.InUse,
		"idle":                   stats.Idle,
		"wait_count":             stats.WaitCount,
		"wait_duration":          stats.WaitDuration.String(),
		"max_idle_closed":        stats.MaxIdleClosed,
		"max_lifetime_closed":    stats.MaxLifetimeClosed,
		"cache_hit_ratio":        calculateCacheHitRatio(stats),
		"connection_utilization": float64(stats.InUse) / float64(stats.MaxOpenConnections) * 100,
	}
}

// calculateCacheHitRatio 计算缓存命中率
func calculateCacheHitRatio(stats sql.DBStats) float64 {
	totalClosed := stats.MaxIdleClosed + stats.MaxLifetimeClosed
	if totalClosed == 0 {
		return 100.0
	}
	return float64(stats.MaxIdleClosed) / float64(totalClosed) * 100
}

// OptimizeQueries 查询优化建议
func OptimizeQueries() map[string]interface{} {
	return map[string]interface{}{
		"recommendations": []string{
			"使用索引优化查询",
			"避免SELECT *，只查询需要的字段",
			"使用批量操作减少数据库往返",
			"合理使用连接池参数",
			"启用查询缓存",
			"使用预处理语句",
		},
		"cache_settings": map[string]interface{}{
			"prep_stmt_cache_size":      256,
			"prep_stmt_cache_sql_limit": 2048,
			"redis_cache_ttl":           "5-10分钟",
		},
		"connection_pool_settings": map[string]interface{}{
			"max_open_conns":     "根据并发量调整",
			"max_idle_conns":     "max_open_conns的10-20%",
			"conn_max_lifetime":  "1小时",
			"conn_max_idle_time": "30分钟",
		},
	}
}
