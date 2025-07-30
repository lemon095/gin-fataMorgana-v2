package services

import (
	"context"
	"encoding/json"
	"gin-fataMorgana/database"
	"gin-fataMorgana/utils"

	"github.com/redis/go-redis/v9"
)

// CurrencyConfig 货币配置结构
type CurrencyConfig struct {
	Symbol string `json:"symbol"`
}

// CurrencyService 货币服务
type CurrencyService struct {
	redisClient *redis.Client
}

// NewCurrencyService 创建货币服务实例
func NewCurrencyService() *CurrencyService {
	return &CurrencyService{
		redisClient: database.RedisClient,
	}
}

// GetCurrentCurrency 获取当前货币配置
func (s *CurrencyService) GetCurrentCurrency(ctx context.Context) (*CurrencyConfig, error) {
	// Redis key
	key := "currency_config:current"

	// 从Redis获取货币配置数据
	data, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			// 如果key不存在，返回默认配置
			return &CurrencyConfig{
				Symbol: "COP",
			}, nil
		}
		return nil, utils.NewAppError(utils.CodeRedisError, "获取货币配置失败")
	}

	// 解析JSON数据
	var currencyConfig CurrencyConfig
	if err := json.Unmarshal([]byte(data), &currencyConfig); err != nil {
		return nil, utils.NewAppError(utils.CodeServer, "解析货币配置数据失败")
	}

	return &currencyConfig, nil
}

// SetCurrentCurrency 设置当前货币配置
func (s *CurrencyService) SetCurrentCurrency(ctx context.Context, config *CurrencyConfig) error {
	// Redis key
	key := "currency_config:current"

	// 序列化配置数据
	data, err := json.Marshal(config)
	if err != nil {
		return utils.NewAppError(utils.CodeServer, "序列化货币配置数据失败")
	}

	// 存储到Redis
	err = s.redisClient.Set(ctx, key, data, 0).Err()
	if err != nil {
		return utils.NewAppError(utils.CodeRedisError, "保存货币配置失败")
	}

	return nil
} 