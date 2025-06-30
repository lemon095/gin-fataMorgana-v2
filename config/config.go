package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 简化后的配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Snowflake SnowflakeConfig `mapstructure:"snowflake"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // 单位秒
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"` // 单位秒
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret            string `mapstructure:"secret"`
	AccessTokenExpire int    `mapstructure:"access_token_expire"`
	RefreshTokenExpire int    `mapstructure:"refresh_token_expire"`
}

// SnowflakeConfig 雪花算法配置
type SnowflakeConfig struct {
	WorkerID      int64 `mapstructure:"worker_id"`
	DatacenterID  int64 `mapstructure:"datacenter_id"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// LoadConfig 加载配置，支持-c/--config参数
func LoadConfig() error {
	// 解析命令行参数
	var configFile string
	flag.StringVar(&configFile, "c", "", "配置文件路径")
	flag.StringVar(&configFile, "config", "", "配置文件路径")
	flag.Parse()

	if configFile == "" {
		configFile = filepath.Join("config", "config.yaml")
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Printf("配置文件不存在: %s", configFile)
		return fmt.Errorf("配置文件不存在: %s", configFile)
	}

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	viper.AutomaticEnv()
	setDefaults()

	log.Printf("配置加载成功，使用文件: %s", configFile)
	return nil
}

// setDefaults 设置默认值
func setDefaults() {
	if GlobalConfig.Server.Host == "" {
		GlobalConfig.Server.Host = "0.0.0.0"
	}
	if GlobalConfig.Server.Port == 0 {
		GlobalConfig.Server.Port = 9001
	}
	if GlobalConfig.Server.Mode == "" {
		GlobalConfig.Server.Mode = "release"
	}
	if GlobalConfig.Database.Port == 0 {
		GlobalConfig.Database.Port = 3306
	}
	if GlobalConfig.Redis.Port == 0 {
		GlobalConfig.Redis.Port = 6379
	}
	if GlobalConfig.Redis.DB == 0 {
		GlobalConfig.Redis.DB = 0
	}
	if GlobalConfig.JWT.AccessTokenExpire == 0 {
		GlobalConfig.JWT.AccessTokenExpire = 3600 // 1小时
	}
	if GlobalConfig.JWT.RefreshTokenExpire == 0 {
		GlobalConfig.JWT.RefreshTokenExpire = 604800 // 7天
	}
	if GlobalConfig.Snowflake.WorkerID == 0 {
		GlobalConfig.Snowflake.WorkerID = 1
	}
	if GlobalConfig.Snowflake.DatacenterID == 0 {
		GlobalConfig.Snowflake.DatacenterID = 1
	}
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.DBName)
}

// GetRedisAddr 获取Redis地址
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// ValidateConfig 简化后的配置验证
func ValidateConfig() error {
	if GlobalConfig == nil {
		return fmt.Errorf("配置未加载")
	}

	// 基本验证
	if GlobalConfig.Database.Host == "" {
		return fmt.Errorf("数据库主机地址不能为空")
	}
	if GlobalConfig.Database.Username == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}
	if GlobalConfig.Database.DBName == "" {
		return fmt.Errorf("数据库名称不能为空")
	}
	if GlobalConfig.Redis.Host == "" {
		return fmt.Errorf("Redis主机地址不能为空")
	}
	if GlobalConfig.JWT.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}

	return nil
}
