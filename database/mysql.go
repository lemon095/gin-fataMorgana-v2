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

// InitMySQL 初始化MySQL连接
func InitMySQL() error {
	cfg := config.GlobalConfig.Database

	// 优化的DSN连接字符串，包含更多缓存和性能参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=true&cachePrepStmts=true&prepStmtCacheSize=256&prepStmtCacheSqlLimit=2048&rejectReadOnly=true&timeout=10s&readTimeout=30s&writeTimeout=30s&multiStatements=true&autocommit=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// 数据库优化配置
		DisableForeignKeyConstraintWhenMigrating: true,  // 禁用外键约束
		PrepareStmt:                              true,  // 启用预处理语句缓存
		SkipDefaultTransaction:                   true,  // 跳过默认事务，提升性能
		DryRun:                                   false, // 生产环境禁用DryRun
	})
	if err != nil {
		return fmt.Errorf("连接MySQL失败: %w", err)
	}

	// 获取底层的sql.DB对象
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 优化的连接池参数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)                                  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)                                  // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Hour) // 连接最大生存时间
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Hour) // 空闲连接最大生存时间

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("MySQL连接测试失败: %w", err)
	}

	log.Printf("MySQL连接成功 - 数据库: %s, 主机: %s:%d", cfg.DBName, cfg.Host, cfg.Port)
	log.Printf("连接池配置 - 最大连接数: %d, 最大空闲连接数: %d", cfg.MaxOpenConns, cfg.MaxIdleConns)
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
		&models.UserLoginLog{},
		&models.Wallet{},
		&models.WalletTransaction{},
		&models.AdminUser{},
		// 在这里添加其他模型
	)
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	log.Println("数据库表迁移成功")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("关闭数据库连接失败: %w", err)
	}

	log.Println("数据库连接已关闭")
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("数据库健康检查失败: %w", err)
	}

	return nil
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
