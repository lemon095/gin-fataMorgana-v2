package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-fataMorgana/database"
	"time"

	"gin-fataMorgana/utils"

	"github.com/redis/go-redis/v9"
)

// UserLevelRule 等级规则
type UserLevelRule struct {
	Level           int    `json:"level"`
	Name            string `json:"name"`
	Logo            string `json:"logo"`
	Requirement     int64  `json:"requirement"`
	RequirementType string `json:"requirement_type"`
	Remark          string `json:"remark"`
}

// UserLevelConfig 用户等级配置
type UserLevelConfig struct {
	Uid        string          `json:"uid"`
	LevelRules []UserLevelRule `json:"level_rules"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	CreatedBy  string          `json:"created_by"`
	UpdatedBy  string          `json:"updated_by"`
}

// UserLevelInfo 用户等级信息结构
type UserLevelInfo struct {
	CurrentLevel         int     `json:"current_level"`
	CurrentLevelName     string  `json:"current_level_name"`
	NextLevel            int     `json:"next_level"`
	NextLevelName        string  `json:"next_level_name"`
	Progress             float64 `json:"progress"`
	Balance              float64 `json:"balance"`
	NextLevelRequirement int     `json:"next_level_requirement"`
}

// UserLevelService 用户等级服务
type UserLevelService struct {
	redisClient   *redis.Client
	walletService *WalletService
}

// NewUserLevelService 创建用户等级服务实例
func NewUserLevelService() *UserLevelService {
	return &UserLevelService{
		redisClient:   database.RedisClient,
		walletService: NewWalletService(),
	}
}

// GetUserLevelInfo 实时计算用户等级信息（不缓存）
func (s *UserLevelService) GetUserLevelInfo(ctx context.Context, uid string) (*UserLevelInfo, error) {
	// 1. 获取用户等级配置
	levelConfig, err := s.getUserLevelConfig(ctx, uid)
	if err != nil {
		return nil, err
	}

	// 2. 获取用户钱包余额
	wallet, err := s.walletService.GetWallet(uid)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeUserLevelGetFailed, "获取用户钱包失败")
	}

	balance := wallet.Balance

	// 3. 根据余额实时计算当前等级和进度
	levelInfo := s.calculateUserLevel(levelConfig, balance)

	return levelInfo, nil
}

// getUserLevelConfig 获取用户等级配置
func (s *UserLevelService) getUserLevelConfig(ctx context.Context, uid string) (*UserLevelConfig, error) {
	// 构建Redis key
	key := fmt.Sprintf("user:level:config:%s", uid)

	// 从Redis获取配置数据
	data, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		// 如果key不存在，尝试从Redis获取默认配置
		if err.Error() == "redis: nil" {
			// 尝试从Redis获取默认配置
			defaultKey := "user:level:config:default"
			defaultData, defaultErr := s.redisClient.Get(ctx, defaultKey).Result()
			if defaultErr == nil {
				// 解析默认配置
				var defaultConfig UserLevelConfig
				if err := json.Unmarshal([]byte(defaultData), &defaultConfig); err != nil {
					// 解析失败，使用硬编码默认配置
					fmt.Printf("JSON解析失败: %v, 数据: %s\n", err, defaultData)
					return s.getDefaultLevelConfig(uid), nil
				}
				// 将默认配置的uid设置为当前用户uid
				defaultConfig.Uid = uid
				fmt.Printf("成功解析默认配置: %+v\n", defaultConfig)
				return &defaultConfig, nil
			}
			// 默认配置也不存在，使用硬编码默认配置
			return s.getDefaultLevelConfig(uid), nil
		}
		return nil, utils.NewAppError(utils.CodeUserLevelGetFailed, "获取用户等级配置失败")
	}

	// 解析JSON数据
	var levelConfig UserLevelConfig
	if err := json.Unmarshal([]byte(data), &levelConfig); err != nil {
		return nil, utils.NewAppError(utils.CodeUserLevelParseFailed, "解析用户等级配置失败")
	}

	return &levelConfig, nil
}

// getDefaultLevelConfig 获取默认等级配置
func (s *UserLevelService) getDefaultLevelConfig(uid string) *UserLevelConfig {
	now := time.Now()
	return &UserLevelConfig{
		Uid: uid,
		LevelRules: []UserLevelRule{
			{Level: 1, Name: "青铜会员", Logo: "", Requirement: 70000, RequirementType: "balance", Remark: "余额达到7万升级"},
			{Level: 2, Name: "白银会员", Logo: "", Requirement: 100000, RequirementType: "balance", Remark: "余额达到10万升级"},
			{Level: 3, Name: "黄金会员", Logo: "", Requirement: 200000, RequirementType: "balance", Remark: "余额达到20万升级"},
		},
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: "system",
		UpdatedBy: "system",
	}
}

// calculateUserLevel 根据余额计算用户等级
func (s *UserLevelService) calculateUserLevel(config *UserLevelConfig, balance float64) *UserLevelInfo {
	// 按等级排序规则
	rules := config.LevelRules
	if len(rules) == 0 {
		return &UserLevelInfo{
			CurrentLevel:         1,
			CurrentLevelName:     "默认等级",
			NextLevel:            1,
			NextLevelName:        "默认等级",
			Progress:             0,
			Balance:              balance,
			NextLevelRequirement: 0,
		}
	}

	// 找到当前等级
	currentLevel := 1
	currentLevelName := "默认等级"
	nextLevel := 1
	nextLevelName := "默认等级"
	nextLevelRequirement := 0
	progress := 0.0

	// 遍历等级规则，找到当前等级和下一等级
	for i, rule := range rules {
		if balance >= float64(rule.Requirement) {
			currentLevel = rule.Level
			currentLevelName = rule.Name
		} else {
			// 找到第一个不满足的等级，这就是下一等级
			if i < len(rules) {
				nextLevel = rule.Level
				nextLevelName = rule.Name
				nextLevelRequirement = int(rule.Requirement)
			}
			break
		}
	}

	// 计算进度百分比
	if nextLevel > currentLevel {
		// 找到当前等级和下一等级的要求
		var currentRequirement, nextRequirement int
		for _, rule := range rules {
			if rule.Level == currentLevel {
				currentRequirement = int(rule.Requirement)
			}
			if rule.Level == nextLevel {
				nextRequirement = int(rule.Requirement)
			}
		}

		if nextRequirement > currentRequirement {
			progress = float64(balance-float64(currentRequirement)) / float64(nextRequirement-currentRequirement) * 100
			if progress > 100 {
				progress = 100
			}
		}
	} else {
		// 已经是最高等级
		progress = 100
	}

	return &UserLevelInfo{
		CurrentLevel:         currentLevel,
		CurrentLevelName:     currentLevelName,
		NextLevel:            nextLevel,
		NextLevelName:        nextLevelName,
		Progress:             progress,
		Balance:              balance,
		NextLevelRequirement: nextLevelRequirement,
	}
}

// GetUserLevel 获取用户当前等级（整数）
func (s *UserLevelService) GetUserLevel(ctx context.Context, uid string) (int, error) {
	// 1. 获取用户等级配置
	levelConfig, err := s.getUserLevelConfig(ctx, uid)
	if err != nil {
		return 1, err
	}

	// 2. 如果没有配置 LevelRules，返回默认等级1
	if len(levelConfig.LevelRules) == 0 {
		return 1, nil
	}

	// 3. 获取用户钱包余额
	wallet, err := s.walletService.GetWallet(uid)
	if err != nil {
		return 1, utils.NewAppError(utils.CodeUserLevelGetFailed, "获取用户钱包失败")
	}

	balance := wallet.Balance

	// 4. 找到用户当前等级
	currentLevel := 1 // 默认等级为1

	// 遍历等级规则，找到用户当前等级
	for _, rule := range levelConfig.LevelRules {
		if balance >= float64(rule.Requirement) {
			currentLevel = rule.Level
		} else {
			// 用户余额不满足这个等级要求，停止遍历
			break
		}
	}

	return currentLevel, nil
}

// GetUserLevelRate 获取用户下一个等级对应的 Requirement 值
func (s *UserLevelService) GetUserLevelRate(ctx context.Context, uid string) (int, error) {
	// 1. 获取用户等级配置
	levelConfig, err := s.getUserLevelConfig(ctx, uid)
	if err != nil {
		return 0, err
	}

	// 2. 如果没有配置 LevelRules，返回 0
	if len(levelConfig.LevelRules) == 0 {
		return 0, nil
	}

	// 3. 获取用户钱包余额
	wallet, err := s.walletService.GetWallet(uid)
	if err != nil {
		return 0, utils.NewAppError(utils.CodeUserLevelGetFailed, "获取用户钱包失败")
	}

	balance := wallet.Balance

	// 4. 找到用户下一个等级对应的 Requirement
	nextRequirement := int64(levelConfig.LevelRules[0].Requirement) // 默认使用 level1 的 Requirement

	// 遍历等级规则，找到用户下一个等级对应的 Requirement
	for i, rule := range levelConfig.LevelRules {
		if balance >= float64(rule.Requirement) {
			// 用户余额满足这个等级要求，检查是否有下一级
			if i+1 < len(levelConfig.LevelRules) {
				// 有下一级，使用下一级的 Requirement
				nextRequirement = levelConfig.LevelRules[i+1].Requirement
			} else {
				// 没有下一级，使用当前等级的 Requirement（最高等级）
				nextRequirement = rule.Requirement
			}
		} else {
			// 用户余额不满足这个等级要求，这个等级就是下一个等级
			nextRequirement = rule.Requirement
			break
		}
	}

	return int(nextRequirement), nil
}

// SetUserLevelConfig 设置用户等级配置到Redis
func (s *UserLevelService) SetUserLevelConfig(ctx context.Context, uid string, config *UserLevelConfig) error {
	// 构建Redis key
	key := fmt.Sprintf("user:level:config:%s", uid)

	// 序列化为JSON
	data, err := json.Marshal(config)
	if err != nil {
		return utils.NewAppError(utils.CodeUserLevelSerializeFailed, "序列化用户等级配置失败")
	}

	// 存储到Redis（设置过期时间为24小时）
	err = s.redisClient.Set(ctx, key, data, 24*60*60).Err()
	if err != nil {
		return utils.NewAppError(utils.CodeUserLevelStoreFailed, "存储用户等级配置失败")
	}

	return nil
}
