package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-fataMorgana/database"
	"math"

	"gin-fataMorgana/utils"

	"github.com/redis/go-redis/v9"
)

// UserLevelInfo Redis中存储的用户等级信息结构
type UserLevelInfo struct {
	CurrentLevel         int     `json:"current_level"`
	CurrentLevelName     string  `json:"current_level_name"`
	NextLevel            int     `json:"next_level"`
	NextLevelName        string  `json:"next_level_name"`
	Progress             float64 `json:"progress"`
	OrderCount           int     `json:"order_count"`
	NextLevelRequirement int     `json:"next_level_requirement"`
}

// UserLevelService 用户等级服务
type UserLevelService struct {
	redisClient *redis.Client
}

// NewUserLevelService 创建用户等级服务实例
func NewUserLevelService() *UserLevelService {
	return &UserLevelService{
		redisClient: database.RedisClient,
	}
}

// GetUserLevelInfo 从Redis获取用户等级信息
func (s *UserLevelService) GetUserLevelInfo(ctx context.Context, uid string) (*UserLevelInfo, error) {
	// 构建Redis key
	key := fmt.Sprintf("user:level:%s", uid)

	// 从Redis获取数据
	data, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		// 如果key不存在，返回默认值
		if err.Error() == "redis: nil" {
			return &UserLevelInfo{
				CurrentLevel:         1,
				CurrentLevelName:     "等级1",
				NextLevel:            2,
				NextLevelName:        "等级2",
				Progress:             0,
				OrderCount:           0,
				NextLevelRequirement: 1,
			}, nil
		}
		return nil, utils.NewAppError(utils.CodeUserLevelGetFailed, "获取用户等级信息失败")
	}

	// 解析JSON数据
	var levelInfo UserLevelInfo
	if err := json.Unmarshal([]byte(data), &levelInfo); err != nil {
		return nil, utils.NewAppError(utils.CodeUserLevelParseFailed, "解析用户等级信息失败")
	}

	return &levelInfo, nil
}

// GetUserLevelRate 获取用户等级进度（整数）
func (s *UserLevelService) GetUserLevelRate(ctx context.Context, uid string) (int, error) {
	levelInfo, err := s.GetUserLevelInfo(ctx, uid)
	if err != nil {
		return 0, err
	}

	// 将进度转换为整数（向下取整）
	return int(math.Floor(levelInfo.Progress)), nil
}

// SetUserLevelInfo 设置用户等级信息到Redis
func (s *UserLevelService) SetUserLevelInfo(ctx context.Context, uid string, levelInfo *UserLevelInfo) error {
	// 构建Redis key
	key := fmt.Sprintf("user:level:%s", uid)

	// 序列化为JSON
	data, err := json.Marshal(levelInfo)
	if err != nil {
		return utils.NewAppError(utils.CodeUserLevelSerializeFailed, "序列化用户等级信息失败")
	}

	// 存储到Redis（设置过期时间为24小时）
	err = s.redisClient.Set(ctx, key, data, 24*60*60).Err()
	if err != nil {
		return utils.NewAppError(utils.CodeUserLevelStoreFailed, "存储用户等级信息失败")
	}

	return nil
}
