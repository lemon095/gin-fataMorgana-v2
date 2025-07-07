package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"gin-fataMorgana/config"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"

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
		return utils.NewAppError(utils.CodeDBConnectFailed, "连接数据库失败")
	}

	// 获取底层的sql.DB对象
	sqlDB, err := db.DB()
	if err != nil {
		return utils.NewAppError(utils.CodeDBInstanceFailed, "获取数据库实例失败")
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
		return utils.NewAppError(utils.CodeDBNotInitialized, "数据库未初始化")
	}

	log.Println("🚀 开始数据库迁移...")

	// 第一步：自动迁移表结构
	log.Println("📋 第一步：创建/更新表结构...")
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
		return utils.NewAppError(utils.CodeDBMigrationFailed, "数据库迁移失败")
	}
	log.Println("✅ 表结构迁移完成")

	// 第二步：添加表注释
	log.Println("📝 第二步：添加表注释...")
	if err := addTableComments(); err != nil {
		log.Printf("⚠️  添加表注释失败: %v", err)
	} else {
		log.Println("✅ 表注释添加完成")
	}

	// 第三步：检测和创建复合索引和优化索引
	log.Println("🔍 第三步：检测和创建优化索引...")
	
	// ===== 索引自动创建功能 =====
	// 方法1：简单注释方式（当前使用）
	// 如需禁用索引自动创建，请注释下面的代码块
	/*
	if err := createOptimizedIndexes(); err != nil {
		log.Printf("⚠️  创建优化索引失败: %v", err)
	} else {
		log.Println("✅ 索引检测和创建完成")
	}
	*/
	// ===== 索引自动创建功能结束 =====
	
	// 如需启用索引自动创建，请取消注释上面的代码块，并注释下面这行
	log.Println("⏭️  索引自动创建已禁用，跳过索引创建")
	
	// 方法2：条件编译方式（可选）
	// 如需使用条件编译，请：
	// 1. 注释掉上面的简单注释代码
	// 2. 取消注释下面的条件编译代码
	// 3. 编译时使用：go build -tags=autoindex 启用索引创建
	// 4. 编译时使用：go build 禁用索引创建
	/*
	// +build autoindex
	if err := createOptimizedIndexes(); err != nil {
		log.Printf("⚠️  创建优化索引失败: %v", err)
	} else {
		log.Println("✅ 索引检测和创建完成")
	}
	// +build !autoindex
	log.Println("⏭️  索引自动创建已禁用，跳过索引创建")
	*/

	// 第四步：为拼单表添加注释
	log.Println("📝 第四步：添加特殊表注释...")
	sqlDB, err := DB.DB()
	if err == nil {
		_, _ = sqlDB.Exec("ALTER TABLE `group_buys` COMMENT = '拼单表 - 记录拼单信息，包括参与人数、付款金额、截止时间等'")
		log.Println("✅ 特殊表注释添加完成")
	}

	log.Println("🎉 数据库迁移全部完成！")
	return nil
}

// createOptimizedIndexes 创建优化的复合索引
func createOptimizedIndexes() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 索引定义结构
	type IndexDef struct {
		TableName string
		IndexName string
		Columns   string
		SQL       string
	}

	// 索引创建SQL语句
	indexDefs := []IndexDef{
		// users表优化索引
		{TableName: "users", IndexName: "idx_users_uid_status_deleted_at", Columns: "uid, status, deleted_at", SQL: "CREATE INDEX idx_users_uid_status_deleted_at ON users(uid, status, deleted_at)"},
		{TableName: "users", IndexName: "idx_users_username_deleted_at", Columns: "username, deleted_at", SQL: "CREATE INDEX idx_users_username_deleted_at ON users(username, deleted_at)"},
		{TableName: "users", IndexName: "idx_users_email_deleted_at", Columns: "email, deleted_at", SQL: "CREATE INDEX idx_users_email_deleted_at ON users(email, deleted_at)"},
		{TableName: "users", IndexName: "idx_users_phone_deleted_at", Columns: "phone, deleted_at", SQL: "CREATE INDEX idx_users_phone_deleted_at ON users(phone, deleted_at)"},
		{TableName: "users", IndexName: "idx_users_invited_by_deleted_at", Columns: "invited_by, deleted_at", SQL: "CREATE INDEX idx_users_invited_by_deleted_at ON users(invited_by, deleted_at)"},
		{TableName: "users", IndexName: "idx_users_status_deleted_at", Columns: "status, deleted_at", SQL: "CREATE INDEX idx_users_status_deleted_at ON users(status, deleted_at)"},
		{TableName: "users", IndexName: "idx_users_created_at", Columns: "created_at", SQL: "CREATE INDEX idx_users_created_at ON users(created_at)"},
		
		// orders表优化索引
		{TableName: "orders", IndexName: "idx_orders_uid_status_created_at", Columns: "uid, status, created_at", SQL: "CREATE INDEX idx_orders_uid_status_created_at ON orders(uid, status, created_at)"},
		{TableName: "orders", IndexName: "idx_orders_status_updated_at", Columns: "status, updated_at", SQL: "CREATE INDEX idx_orders_status_updated_at ON orders(status, updated_at)"},
		{TableName: "orders", IndexName: "idx_orders_period_number_status", Columns: "period_number, status", SQL: "CREATE INDEX idx_orders_period_number_status ON orders(period_number, status)"},
		{TableName: "orders", IndexName: "idx_orders_expire_time_status", Columns: "expire_time, status", SQL: "CREATE INDEX idx_orders_expire_time_status ON orders(expire_time, status)"},
		{TableName: "orders", IndexName: "idx_orders_auditor_uid_status", Columns: "auditor_uid, status", SQL: "CREATE INDEX idx_orders_auditor_uid_status ON orders(auditor_uid, status)"},
		{TableName: "orders", IndexName: "idx_orders_is_system_order_status", Columns: "is_system_order, status", SQL: "CREATE INDEX idx_orders_is_system_order_status ON orders(is_system_order, status)"},
		
		// wallet_transactions表优化索引
		{TableName: "wallet_transactions", IndexName: "idx_wallet_transactions_uid_type_status", Columns: "uid, type, status", SQL: "CREATE INDEX idx_wallet_transactions_uid_type_status ON wallet_transactions(uid, type, status)"},
		{TableName: "wallet_transactions", IndexName: "idx_wallet_transactions_uid_created_at", Columns: "uid, created_at", SQL: "CREATE INDEX idx_wallet_transactions_uid_created_at ON wallet_transactions(uid, created_at)"},
		{TableName: "wallet_transactions", IndexName: "idx_wallet_transactions_type_status_created_at", Columns: "type, status, created_at", SQL: "CREATE INDEX idx_wallet_transactions_type_status_created_at ON wallet_transactions(type, status, created_at)"},
		{TableName: "wallet_transactions", IndexName: "idx_wallet_transactions_transaction_no", Columns: "transaction_no", SQL: "CREATE INDEX idx_wallet_transactions_transaction_no ON wallet_transactions(transaction_no)"},
		{TableName: "wallet_transactions", IndexName: "idx_wallet_transactions_amount", Columns: "amount", SQL: "CREATE INDEX idx_wallet_transactions_amount ON wallet_transactions(amount)"},
		
		// group_buys表优化索引
		{TableName: "group_buys", IndexName: "idx_group_buys_deadline_status", Columns: "deadline, status", SQL: "CREATE INDEX idx_group_buys_deadline_status ON group_buys(deadline, status)"},
		{TableName: "group_buys", IndexName: "idx_group_buys_uid_deadline", Columns: "uid, deadline", SQL: "CREATE INDEX idx_group_buys_uid_deadline ON group_buys(uid, deadline)"},
		{TableName: "group_buys", IndexName: "idx_group_buys_type_status_deadline", Columns: "group_buy_type, status, deadline", SQL: "CREATE INDEX idx_group_buys_type_status_deadline ON group_buys(group_buy_type, status, deadline)"},
		{TableName: "group_buys", IndexName: "idx_group_buys_creator_uid_status", Columns: "creator_uid, status", SQL: "CREATE INDEX idx_group_buys_creator_uid_status ON group_buys(creator_uid, status)"},
		
		// user_login_logs表优化索引
		{TableName: "user_login_logs", IndexName: "idx_user_login_logs_uid_login_time", Columns: "uid, login_time", SQL: "CREATE INDEX idx_user_login_logs_uid_login_time ON user_login_logs(uid, login_time)"},
		{TableName: "user_login_logs", IndexName: "idx_user_login_logs_uid_status_login_time", Columns: "uid, status, login_time", SQL: "CREATE INDEX idx_user_login_logs_uid_status_login_time ON user_login_logs(uid, status, login_time)"},
		{TableName: "user_login_logs", IndexName: "idx_user_login_logs_uid_login_ip", Columns: "uid, login_ip", SQL: "CREATE INDEX idx_user_login_logs_uid_login_ip ON user_login_logs(uid, login_ip)"},
		{TableName: "user_login_logs", IndexName: "idx_user_login_logs_uid_status", Columns: "uid, status", SQL: "CREATE INDEX idx_user_login_logs_uid_status ON user_login_logs(uid, status)"},
		{TableName: "user_login_logs", IndexName: "idx_user_login_logs_created_at", Columns: "created_at", SQL: "CREATE INDEX idx_user_login_logs_created_at ON user_login_logs(created_at)"},
		
		// lottery_periods表优化索引
		{TableName: "lottery_periods", IndexName: "idx_lottery_periods_status_order_end_time", Columns: "status, order_end_time", SQL: "CREATE INDEX idx_lottery_periods_status_order_end_time ON lottery_periods(status, order_end_time)"},
		{TableName: "lottery_periods", IndexName: "idx_lottery_periods_order_start_time_order_end_time", Columns: "order_start_time, order_end_time", SQL: "CREATE INDEX idx_lottery_periods_order_start_time_order_end_time ON lottery_periods(order_start_time, order_end_time)"},
		{TableName: "lottery_periods", IndexName: "idx_lottery_periods_lottery_result", Columns: "lottery_result", SQL: "CREATE INDEX idx_lottery_periods_lottery_result ON lottery_periods(lottery_result)"},
		
		// admin_users表优化索引
		{TableName: "admin_users", IndexName: "idx_admin_users_role_deleted_at", Columns: "role, deleted_at", SQL: "CREATE INDEX idx_admin_users_role_deleted_at ON admin_users(role, deleted_at)"},
		{TableName: "admin_users", IndexName: "idx_admin_users_status_deleted_at", Columns: "status, deleted_at", SQL: "CREATE INDEX idx_admin_users_status_deleted_at ON admin_users(status, deleted_at)"},
		{TableName: "admin_users", IndexName: "idx_admin_users_parent_id_deleted_at", Columns: "parent_id, deleted_at", SQL: "CREATE INDEX idx_admin_users_parent_id_deleted_at ON admin_users(parent_id, deleted_at)"},
		{TableName: "admin_users", IndexName: "idx_admin_users_created_at", Columns: "created_at", SQL: "CREATE INDEX idx_admin_users_created_at ON admin_users(created_at)"},
		
		// wallets表优化索引
		{TableName: "wallets", IndexName: "idx_wallets_status", Columns: "status", SQL: "CREATE INDEX idx_wallets_status ON wallets(status)"},
		{TableName: "wallets", IndexName: "idx_wallets_balance", Columns: "balance", SQL: "CREATE INDEX idx_wallets_balance ON wallets(balance)"},
		{TableName: "wallets", IndexName: "idx_wallets_last_active_at", Columns: "last_active_at", SQL: "CREATE INDEX idx_wallets_last_active_at ON wallets(last_active_at)"},
		{TableName: "wallets", IndexName: "idx_wallets_status_balance", Columns: "status, balance", SQL: "CREATE INDEX idx_wallets_status_balance ON wallets(status, balance)"},
		
		// amount_config表优化索引
		{TableName: "amount_config", IndexName: "idx_amount_config_type_is_active_sort", Columns: "type, is_active, sort_order", SQL: "CREATE INDEX idx_amount_config_type_is_active_sort ON amount_config(type, is_active, sort_order)"},
		{TableName: "amount_config", IndexName: "idx_amount_config_amount", Columns: "amount", SQL: "CREATE INDEX idx_amount_config_amount ON amount_config(amount)"},
		
		// announcements表优化索引
		{TableName: "announcements", IndexName: "idx_announcements_status", Columns: "status", SQL: "CREATE INDEX idx_announcements_status ON announcements(status)"},
		{TableName: "announcements", IndexName: "idx_announcements_tag", Columns: "tag", SQL: "CREATE INDEX idx_announcements_tag ON announcements(tag)"},
		{TableName: "announcements", IndexName: "idx_announcements_status_deleted_at_created_at", Columns: "status, deleted_at, created_at", SQL: "CREATE INDEX idx_announcements_status_deleted_at_created_at ON announcements(status, deleted_at, created_at)"},
		{TableName: "announcements", IndexName: "idx_announcements_created_at", Columns: "created_at", SQL: "CREATE INDEX idx_announcements_created_at ON announcements(created_at)"},
		
		// member_level表优化索引
		{TableName: "member_level", IndexName: "idx_member_level_level_deleted_at", Columns: "level, deleted_at", SQL: "CREATE INDEX idx_member_level_level_deleted_at ON member_level(level, deleted_at)"},
		{TableName: "member_level", IndexName: "idx_member_level_cashback_ratio", Columns: "cashback_ratio", SQL: "CREATE INDEX idx_member_level_cashback_ratio ON member_level(cashback_ratio)"},
		
		// announcement_banners表优化索引
		{TableName: "announcement_banners", IndexName: "idx_announcement_banners_announcement_id_deleted_at_sort", Columns: "announcement_id, deleted_at, sort", SQL: "CREATE INDEX idx_announcement_banners_announcement_id_deleted_at_sort ON announcement_banners(announcement_id, deleted_at, sort)"},
		{TableName: "announcement_banners", IndexName: "idx_announcement_banners_sort", Columns: "sort", SQL: "CREATE INDEX idx_announcement_banners_sort ON announcement_banners(sort)"},
	}

	// 检测并创建索引
	createdCount := 0
	skippedCount := 0
	failedCount := 0

	for _, indexDef := range indexDefs {
		// 检测索引是否存在
		exists, err := checkIndexExists(sqlDB, indexDef.TableName, indexDef.IndexName)
		if err != nil {
			log.Printf("检测索引 %s 失败: %v", indexDef.IndexName, err)
			failedCount++
			continue
		}

		if exists {
			log.Printf("索引 %s 已存在，跳过创建", indexDef.IndexName)
			skippedCount++
			continue
		}

		// 创建索引
		if _, err := sqlDB.Exec(indexDef.SQL); err != nil {
			log.Printf("创建索引 %s 失败: %v", indexDef.IndexName, err)
			failedCount++
		} else {
			log.Printf("✅ 索引创建成功: %s (%s)", indexDef.IndexName, indexDef.Columns)
			createdCount++
		}
	}

	log.Printf("📊 索引检测完成 - 创建: %d, 跳过: %d, 失败: %d", createdCount, skippedCount, failedCount)
	return nil
}

// checkIndexExists 检测索引是否存在
func checkIndexExists(sqlDB *sql.DB, tableName, indexName string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM information_schema.statistics 
		WHERE table_schema = DATABASE() 
		AND table_name = ? 
		AND index_name = ?
	`
	
	var count int
	err := sqlDB.QueryRow(query, tableName, indexName).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
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
			return utils.NewAppError(utils.CodeDBInstanceFailed, "获取数据库实例失败")
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
		return utils.NewAppError(utils.CodeDBNotInitialized, "数据库未初始化")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return utils.NewAppError(utils.CodeDBInstanceFailed, "获取数据库实例失败")
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

// WithCache 带缓存的查询（添加超时控制）
func (qc *QueryCache) WithCache(ctx context.Context, key string, expiration time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	// 设置查询超时
	queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 先尝试从Redis获取缓存
	if cached, err := RedisClient.Get(queryCtx, key).Result(); err == nil {
		return cached, nil
	}

	// 缓存未命中，执行查询
	result, err := fn()
	if err != nil {
		return nil, err
	}

	// 缓存结果
	RedisClient.Set(queryCtx, key, result, expiration)
	return result, nil
}

// InvalidateCache 清除缓存
func (qc *QueryCache) InvalidateCache(ctx context.Context, pattern string) error {
	// 设置操作超时
	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	keys, err := RedisClient.Keys(opCtx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return RedisClient.Del(opCtx, keys...).Err()
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

// CheckAndCreateIndexes 检测并创建缺失的索引
func CheckAndCreateIndexes() error {
	if DB == nil {
		return utils.NewAppError(utils.CodeDBNotInitialized, "数据库未初始化")
	}

	log.Println("🔍 开始检测数据库索引...")
	return createOptimizedIndexes()
}

// ShowAllIndexes 显示当前数据库的所有索引
func ShowAllIndexes() error {
	if DB == nil {
		return utils.NewAppError(utils.CodeDBNotInitialized, "数据库未初始化")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	query := `
		SELECT 
			table_name,
			index_name,
			column_name,
			seq_in_index,
			cardinality,
			non_unique
		FROM information_schema.statistics 
		WHERE table_schema = DATABASE()
		ORDER BY table_name, index_name, seq_in_index
	`

	rows, err := sqlDB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	log.Println("📋 当前数据库索引列表:")
	log.Println(strings.Repeat("=", 80))
	log.Printf("%-20s %-30s %-20s %-10s %-10s %-10s", "表名", "索引名", "字段名", "序号", "基数", "唯一性")
	log.Println(strings.Repeat("-", 80))

	var currentTable, currentIndex string
	var columns []string

	for rows.Next() {
		var tableName, indexName, columnName string
		var seqInIndex, cardinality, nonUnique int

		err := rows.Scan(&tableName, &indexName, &columnName, &seqInIndex, &cardinality, &nonUnique)
		if err != nil {
			log.Printf("读取索引信息失败: %v", err)
			continue
		}

		// 如果是新的表或索引，打印之前的信息
		if currentTable != tableName || currentIndex != indexName {
			if len(columns) > 0 {
				uniqueText := "唯一"
				if nonUnique == 1 {
					uniqueText = "非唯一"
				}
				log.Printf("%-20s %-30s %-20s %-10d %-10d %-10s", 
					currentTable, currentIndex, columns[0], 1, cardinality, uniqueText)
				
				// 打印复合索引的其他字段
				for i := 1; i < len(columns); i++ {
					log.Printf("%-20s %-30s %-20s %-10d %-10s %-10s", 
						"", "", columns[i], i+1, "", "")
				}
			}
			
			// 重置当前信息
			currentTable = tableName
			currentIndex = indexName
			columns = []string{columnName}
		} else {
			// 同一索引的其他字段
			columns = append(columns, columnName)
		}
	}

	// 打印最后一个索引
	if len(columns) > 0 {
		uniqueText := "唯一"
		log.Printf("%-20s %-30s %-20s %-10d %-10s %-10s", 
			currentTable, currentIndex, columns[0], 1, "", uniqueText)
		
		for i := 1; i < len(columns); i++ {
			log.Printf("%-20s %-30s %-20s %-10d %-10s %-10s", 
				"", "", columns[i], i+1, "", "")
		}
	}

	log.Println(strings.Repeat("=", 80))
	log.Println("✅ 索引信息显示完成")

	return nil
}
