package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"gin-fataMorgana/utils"
)

// Config 简化后的配置结构体
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
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
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`  // 单位秒
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
	Secret             string `mapstructure:"secret"`
	AccessTokenExpire  int    `mapstructure:"access_token_expire"`
	RefreshTokenExpire int    `mapstructure:"refresh_token_expire"`
}

// SnowflakeConfig 雪花算法配置
type SnowflakeConfig struct {
	WorkerID     int64 `mapstructure:"worker_id"`
	DatacenterID int64 `mapstructure:"datacenter_id"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// LoadConfig 加载配置，支持-c/--config参数和环境变量
func LoadConfig() error {
	// 解析命令行参数
	var configFile string
	flag.StringVar(&configFile, "c", "", "配置文件路径")
	flag.StringVar(&configFile, "config", "", "配置文件路径")
	flag.Parse()

	// 设置默认配置文件
	if configFile == "" {
		configFile = "config.yaml"
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Printf("配置文件不存在: %s，使用默认配置", configFile)
		// 如果配置文件不存在，使用默认配置
		GlobalConfig = &Config{}
		setDefaults()
		overrideWithEnvVars()
		log.Printf("使用默认配置和环境变量")
		return nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		return utils.NewAppError(utils.CodeConfigReadFailed, "读取配置文件失败")
	}

	// 解析配置文件
	if err := yaml.Unmarshal(data, &GlobalConfig); err != nil {
		return utils.NewAppError(utils.CodeConfigParseFailed, "解析配置文件失败")
	}

	// 设置默认值
	setDefaults()

	// 使用环境变量覆盖配置
	overrideWithEnvVars()

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
		GlobalConfig.Server.Mode = "debug"
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
		GlobalConfig.JWT.AccessTokenExpire = 86400 // 1天
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

// overrideWithEnvVars 使用环境变量覆盖配置
func overrideWithEnvVars() {
	// 服务器配置
	if env := os.Getenv("GIN_MODE"); env != "" {
		GlobalConfig.Server.Mode = env
	}
	if env := os.Getenv("SERVER_MODE"); env != "" {
		GlobalConfig.Server.Mode = env
	}

	// 数据库配置
	if env := os.Getenv("DATABASE_HOST"); env != "" {
		GlobalConfig.Database.Host = env
	}
	if env := os.Getenv("MYSQL_HOST"); env != "" {
		GlobalConfig.Database.Host = env
	}
	if env := os.Getenv("DATABASE_PORT"); env != "" {
		if port := parsePort(env); port > 0 {
			GlobalConfig.Database.Port = port
		}
	}
	if env := os.Getenv("MYSQL_PORT"); env != "" {
		if port := parsePort(env); port > 0 {
			GlobalConfig.Database.Port = port
		}
	}
	if env := os.Getenv("DATABASE_USERNAME"); env != "" {
		GlobalConfig.Database.Username = env
	}
	if env := os.Getenv("MYSQL_USERNAME"); env != "" {
		GlobalConfig.Database.Username = env
	}
	if env := os.Getenv("DATABASE_PASSWORD"); env != "" {
		GlobalConfig.Database.Password = env
	}
	if env := os.Getenv("MYSQL_PASSWORD"); env != "" {
		GlobalConfig.Database.Password = env
	}
	if env := os.Getenv("DATABASE_DBNAME"); env != "" {
		GlobalConfig.Database.DBName = env
	}
	if env := os.Getenv("MYSQL_DATABASE"); env != "" {
		GlobalConfig.Database.DBName = env
	}

	// Redis配置
	if env := os.Getenv("REDIS_HOST"); env != "" {
		GlobalConfig.Redis.Host = env
	}
	if env := os.Getenv("REDIS_PORT"); env != "" {
		if port := parsePort(env); port > 0 {
			GlobalConfig.Redis.Port = port
		}
	}
	if env := os.Getenv("REDIS_PASSWORD"); env != "" {
		GlobalConfig.Redis.Password = env
	}
}

// parsePort 解析端口号
func parsePort(s string) int {
	var port int
	fmt.Sscanf(s, "%d", &port)
	return port
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
	// 检查配置是否已加载
	if GlobalConfig == nil {
		return utils.NewAppError(utils.CodeConfigNotLoaded, "配置未加载")
	}

	// 验证数据库配置
	if GlobalConfig.Database.Host == "" {
		return utils.NewAppError(utils.CodeDBHostEmpty, "数据库主机地址不能为空")
	}
	if GlobalConfig.Database.Username == "" {
		return utils.NewAppError(utils.CodeDBUserEmpty, "数据库用户名不能为空")
	}
	if GlobalConfig.Database.DBName == "" {
		return utils.NewAppError(utils.CodeDBNameEmpty, "数据库名称不能为空")
	}

	// 验证Redis配置
	if GlobalConfig.Redis.Host == "" {
		return utils.NewAppError(utils.CodeRedisHostEmpty, "Redis主机地址不能为空")
	}

	// 验证JWT配置
	if GlobalConfig.JWT.Secret == "" {
		return utils.NewAppError(utils.CodeJWTSecretEmpty, "JWT密钥不能为空")
	}

	return nil
}
