package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 配置结构体
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Log       LogConfig       `mapstructure:"log"`
	Snowflake SnowflakeConfig `mapstructure:"snowflake"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	Domain string `mapstructure:"domain"`
	Mode   string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	Charset         string `mapstructure:"charset"`
	ParseTime       bool   `mapstructure:"parse_time"`
	Loc             string `mapstructure:"loc"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
	MaxRetries   int    `mapstructure:"max_retries"`
	DialTimeout  int    `mapstructure:"dial_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret             string `mapstructure:"secret"`
	AccessTokenExpire  int64  `mapstructure:"access_token_expire"`
	RefreshTokenExpire int64  `mapstructure:"refresh_token_expire"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	SessionTimeout   int `mapstructure:"session_timeout"`
	MaxLoginAttempts int `mapstructure:"max_login_attempts"`
	LockoutDuration  int `mapstructure:"lockout_duration"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string     `mapstructure:"level"`
	Format string     `mapstructure:"format"`
	Output string     `mapstructure:"output"`
	File   FileConfig `mapstructure:"file"`
}

// FileConfig 文件配置
type FileConfig struct {
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

// SnowflakeConfig 雪花算法配置
type SnowflakeConfig struct {
	WorkerID     int64 `mapstructure:"worker_id"`
	DatacenterID int64 `mapstructure:"datacenter_id"`
}

var GlobalConfig *Config

// LoadConfig 加载配置
func LoadConfig() error {
	// 设置配置文件路径
	configPath := "config"
	configName := "config"
	configType := "yaml"

	// 检查配置文件是否存在
	configFile := filepath.Join(configPath, configName+"."+configType)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Printf("配置文件不存在: %s", configFile)
		return fmt.Errorf("配置文件不存在: %s", configFile)
	}

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置到结构体
	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 设置环境变量支持
	viper.AutomaticEnv()

	log.Println("配置加载成功")
	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Username, c.Password, c.Host, c.Port, c.DBName, c.Charset, c.ParseTime, c.Loc)
}

// GetRedisAddr 获取Redis地址
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// ValidateConfig 验证配置
func ValidateConfig() error {
	if GlobalConfig == nil {
		return fmt.Errorf("配置未加载")
	}

	// 验证服务器配置
	if err := validateServerConfig(&GlobalConfig.Server); err != nil {
		return fmt.Errorf("服务器配置错误: %w", err)
	}

	// 验证数据库配置
	if err := validateDatabaseConfig(&GlobalConfig.Database); err != nil {
		return fmt.Errorf("数据库配置错误: %w", err)
	}

	// 验证Redis配置
	if err := validateRedisConfig(&GlobalConfig.Redis); err != nil {
		return fmt.Errorf("Redis配置错误: %w", err)
	}

	// 验证JWT配置
	if err := validateJWTConfig(&GlobalConfig.JWT); err != nil {
		return fmt.Errorf("JWT配置错误: %w", err)
	}

	// 验证雪花算法配置
	if err := validateSnowflakeConfig(&GlobalConfig.Snowflake); err != nil {
		return fmt.Errorf("雪花算法配置错误: %w", err)
	}

	return nil
}

// validateServerConfig 验证服务器配置
func validateServerConfig(cfg *ServerConfig) error {
	if cfg.Host == "" {
		return fmt.Errorf("服务器主机地址不能为空")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("服务器端口必须在1-65535之间")
	}
	if cfg.Mode == "" {
		cfg.Mode = "release" // 设置默认模式
	}
	return nil
}

// validateDatabaseConfig 验证数据库配置
func validateDatabaseConfig(cfg *DatabaseConfig) error {
	if cfg.Host == "" {
		return fmt.Errorf("数据库主机地址不能为空")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("数据库端口必须在1-65535之间")
	}
	if cfg.Username == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}
	if cfg.DBName == "" {
		return fmt.Errorf("数据库名称不能为空")
	}
	if cfg.Charset == "" {
		cfg.Charset = "utf8mb4" // 设置默认字符集
	}
	if cfg.Loc == "" {
		cfg.Loc = "Local" // 设置默认时区
	}
	
	// 设置连接池默认值
	if cfg.MaxIdleConns <= 0 {
		cfg.MaxIdleConns = 10
	}
	if cfg.MaxOpenConns <= 0 {
		cfg.MaxOpenConns = 100
	}
	if cfg.ConnMaxLifetime <= 0 {
		cfg.ConnMaxLifetime = 1 // 1小时
	}
	if cfg.ConnMaxIdleTime <= 0 {
		cfg.ConnMaxIdleTime = 1 // 1小时
	}
	
	return nil
}

// validateRedisConfig 验证Redis配置
func validateRedisConfig(cfg *RedisConfig) error {
	if cfg.Host == "" {
		return fmt.Errorf("Redis主机地址不能为空")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("Redis端口必须在1-65535之间")
	}
	if cfg.PoolSize <= 0 {
		cfg.PoolSize = 10 // 设置默认连接池大小
	}
	if cfg.MinIdleConns <= 0 {
		cfg.MinIdleConns = 5 // 设置默认最小空闲连接数
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 3 // 设置默认重试次数
	}
	if cfg.DialTimeout <= 0 {
		cfg.DialTimeout = 5 // 设置默认连接超时时间
	}
	if cfg.ReadTimeout <= 0 {
		cfg.ReadTimeout = 3 // 设置默认读取超时时间
	}
	if cfg.WriteTimeout <= 0 {
		cfg.WriteTimeout = 3 // 设置默认写入超时时间
	}
	return nil
}

// validateJWTConfig 验证JWT配置
func validateJWTConfig(cfg *JWTConfig) error {
	if cfg.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}
	if len(cfg.Secret) < 32 {
		return fmt.Errorf("JWT密钥长度不能少于32位")
	}
	if cfg.AccessTokenExpire <= 0 {
		cfg.AccessTokenExpire = 3600 // 设置默认访问令牌过期时间（1小时）
	}
	if cfg.RefreshTokenExpire <= 0 {
		cfg.RefreshTokenExpire = 604800 // 设置默认刷新令牌过期时间（7天）
	}
	return nil
}

// validateSnowflakeConfig 验证雪花算法配置
func validateSnowflakeConfig(cfg *SnowflakeConfig) error {
	if cfg.WorkerID < 0 || cfg.WorkerID > 99 {
		return fmt.Errorf("WorkerID必须在0-99之间")
	}
	if cfg.DatacenterID < 0 || cfg.DatacenterID > 9 {
		return fmt.Errorf("DatacenterID必须在0-9之间")
	}
	return nil
}
