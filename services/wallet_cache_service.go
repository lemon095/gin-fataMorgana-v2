package services

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"
)

// WalletCacheService 统一的钱包缓存服务
type WalletCacheService struct {
	// 用于防止缓存击穿的互斥锁
	mutexMap sync.Map
}

func NewWalletCacheService() *WalletCacheService {
	return &WalletCacheService{}
}

// 生成钱包缓存Key
func (s *WalletCacheService) generateWalletKey(uid string) string {
	return utils.RedisKeys.GenerateWalletBalanceKey(uid)
}

// 生成空值缓存Key（防止缓存穿透）
func (s *WalletCacheService) generateEmptyKey(uid string) string {
	return utils.RedisKeys.GenerateWalletEmptyKey(uid)
}

// 生成用户登录时间Key
func (s *WalletCacheService) generateUserLoginKey(uid string) string {
	return utils.RedisKeys.GenerateUserLoginTimeKey(uid)
}

// 获取用户级别的互斥锁（防止缓存击穿）
func (s *WalletCacheService) getUserMutex(uid string) *sync.Mutex {
	value, _ := s.mutexMap.LoadOrStore(uid, &sync.Mutex{})
	return value.(*sync.Mutex)
}

// 缓存钱包余额（带过期时间，基于用户登录状态）
func (s *WalletCacheService) CacheWalletBalance(ctx context.Context, wallet *models.Wallet) error {
	if wallet == nil || wallet.Uid == "" {
		return utils.NewAppError(utils.CodeInvalidParams, "钱包数据无效")
	}

	// 生成缓存Key
	cacheKey := s.generateWalletKey(wallet.Uid)

	// 将钱包数据转换为JSON
	walletJSON, err := json.Marshal(wallet)
	if err != nil {
		return utils.NewAppError(utils.CodeInvalidParams, "钱包数据序列化失败")
	}

	// 计算过期时间（30天）
	expiration := 30 * 24 * time.Hour

	// 缓存钱包数据（带过期时间）
	err = database.GetGlobalRedisHelper().Set(ctx, cacheKey, string(walletJSON), expiration)
	if err != nil {
		return utils.NewAppError(utils.CodeRedisError, "缓存钱包余额失败")
	}

	return nil
}

// 缓存空值（短期过期，防止缓存穿透）
func (s *WalletCacheService) CacheEmptyWallet(ctx context.Context, uid string) error {
	if uid == "" {
		return utils.NewAppError(utils.CodeInvalidParams, "用户ID不能为空")
	}

	// 生成空值缓存Key
	emptyKey := s.generateEmptyKey(uid)

	// 缓存空值，过期时间较短（10分钟）
	err := database.GetGlobalRedisHelper().Set(ctx, emptyKey, "empty", 10*time.Minute)
	if err != nil {
		return utils.NewAppError(utils.CodeRedisError, "缓存空值失败")
	}

	return nil
}

// 检查是否为空值缓存
func (s *WalletCacheService) IsEmptyCached(ctx context.Context, uid string) (bool, error) {
	if uid == "" {
		return false, utils.NewAppError(utils.CodeInvalidParams, "用户ID不能为空")
	}

	// 生成空值缓存Key
	emptyKey := s.generateEmptyKey(uid)

	// 检查空值缓存是否存在
	exists, err := database.GetGlobalRedisHelper().Exists(ctx, emptyKey)
	if err != nil {
		return false, utils.NewAppError(utils.CodeRedisError, "检查空值缓存失败")
	}

	return exists > 0, nil
}

// 获取缓存的钱包余额
func (s *WalletCacheService) GetCachedWalletBalance(ctx context.Context, uid string) (*models.Wallet, error) {
	if uid == "" {
		return nil, utils.NewAppError(utils.CodeInvalidParams, "用户ID不能为空")
	}

	// 生成缓存Key
	cacheKey := s.generateWalletKey(uid)

	// 获取缓存数据
	walletJSON, err := database.GetGlobalRedisHelper().Get(ctx, cacheKey)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeRedisError, "获取缓存钱包数据失败")
	}

	if walletJSON == "" {
		return nil, utils.NewAppError(utils.CodeWalletGetFailed, "钱包数据不存在")
	}

	// 反序列化钱包数据
	var wallet models.Wallet
	err = json.Unmarshal([]byte(walletJSON), &wallet)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeInvalidParams, "钱包数据反序列化失败")
	}

	return &wallet, nil
}

// 获取钱包余额（基于用户登录状态的缓存策略）
func (s *WalletCacheService) GetWalletBalanceWithCache(ctx context.Context, uid string) (*models.Wallet, error) {
	if uid == "" {
		return nil, utils.NewAppError(utils.CodeInvalidParams, "用户ID不能为空")
	}

	// 1. 先尝试从缓存获取
	wallet, err := s.GetCachedWalletBalance(ctx, uid)
	if err == nil {
		return wallet, nil
	}

	// 2. 检查是否为空值缓存（防止缓存穿透）
	isEmpty, err := s.IsEmptyCached(ctx, uid)
	if err == nil && isEmpty {
		return nil, utils.NewAppError(utils.CodeWalletGetFailed, "钱包不存在")
	}

	// 3. 获取用户级别的互斥锁（防止缓存击穿）
	mutex := s.getUserMutex(uid)
	mutex.Lock()
	defer mutex.Unlock()

	// 4. 双重检查：再次尝试从缓存获取
	wallet, err = s.GetCachedWalletBalance(ctx, uid)
	if err == nil {
		return wallet, nil
	}

	// 5. 从数据库获取
	walletRepo := database.NewWalletRepository()
	wallet, err = walletRepo.FindWalletByUid(ctx, uid)
	if err != nil {
		// 6. 如果钱包不存在，缓存空值防止穿透
		if cacheErr := s.CacheEmptyWallet(ctx, uid); cacheErr != nil {
			utils.LogWarn(nil, "缓存空值失败: %v", cacheErr)
		}
		return nil, utils.NewAppError(utils.CodeWalletGetFailed, "获取钱包数据失败")
	}

	// 7. 缓存到Redis（带过期时间）
	if cacheErr := s.CacheWalletBalance(ctx, wallet); cacheErr != nil {
		// 缓存失败不影响主流程，只记录日志
		utils.LogWarn(nil, "缓存钱包余额失败: %v", cacheErr)
	}

	return wallet, nil
}

// 事件驱动更新钱包余额（余额变化时调用）
func (s *WalletCacheService) UpdateWalletBalanceOnEvent(ctx context.Context, uid string, newBalance float64) error {
	if uid == "" {
		return utils.NewAppError(utils.CodeInvalidParams, "用户ID不能为空")
	}

	// 1. 获取现有缓存数据
	wallet, err := s.GetCachedWalletBalance(ctx, uid)
	if err != nil {
		// 缓存不存在，从数据库获取
		walletRepo := database.NewWalletRepository()
		wallet, err = walletRepo.FindWalletByUid(ctx, uid)
		if err != nil {
			return utils.NewAppError(utils.CodeWalletGetFailed, "获取钱包数据失败")
		}
	}

	// 2. 更新余额
	wallet.Balance = newBalance
	wallet.UpdatedAt = time.Now()

	// 3. 更新缓存（保持原有过期时间）
	return s.CacheWalletBalance(ctx, wallet)
}

// 用户登录时延长钱包缓存过期时间
func (s *WalletCacheService) ExtendWalletCacheOnLogin(ctx context.Context, uid string) error {
	if uid == "" {
		return utils.NewAppError(utils.CodeInvalidParams, "用户ID不能为空")
	}

	// 1. 更新用户登录时间
	loginKey := s.generateUserLoginKey(uid)
	err := database.GetGlobalRedisHelper().Set(ctx, loginKey, time.Now().Unix(), 30*24*time.Hour)
	if err != nil {
		return utils.NewAppError(utils.CodeRedisError, "更新用户登录时间失败")
	}

	// 2. 获取当前钱包缓存
	wallet, err := s.GetCachedWalletBalance(ctx, uid)
	if err != nil {
		// 如果缓存不存在，从数据库获取并缓存
		walletRepo := database.NewWalletRepository()
		wallet, err = walletRepo.FindWalletByUid(ctx, uid)
		if err != nil {
			return utils.NewAppError(utils.CodeWalletGetFailed, "获取钱包数据失败")
		}
	}

	// 3. 更新钱包最后活跃时间
	wallet.UpdateLastActive()

	// 4. 重新缓存钱包数据（延长过期时间）
	return s.CacheWalletBalance(ctx, wallet)
}

// 检查用户是否在指定时间内有登录活动
func (s *WalletCacheService) HasRecentLogin(ctx context.Context, uid string, duration time.Duration) (bool, error) {
	if uid == "" {
		return false, utils.NewAppError(utils.CodeInvalidParams, "用户ID不能为空")
	}

	loginKey := s.generateUserLoginKey(uid)
	loginTimeStr, err := database.GetGlobalRedisHelper().Get(ctx, loginKey)
	if err != nil {
		return false, nil // 没有登录记录，返回false
	}

	if loginTimeStr == "" {
		return false, nil
	}

	// 解析登录时间
	var loginTimeUnix int64
	if err := json.Unmarshal([]byte(loginTimeStr), &loginTimeUnix); err != nil {
		return false, nil
	}

	loginTime := time.Unix(loginTimeUnix, 0)
	now := time.Now()

	return now.Sub(loginTime) <= duration, nil
}

// 清理过期钱包缓存
func (s *WalletCacheService) CleanupExpiredWalletCache(ctx context.Context) error {
	// 获取所有钱包缓存Key
	pattern := utils.RedisKeys.GetWalletKeyPattern()
	keys, err := database.Keys(ctx, pattern)
	if err != nil {
		return utils.NewAppError(utils.CodeRedisError, "获取钱包缓存Key失败")
	}

	expiredCount := 0

	for _, key := range keys {
		// 提取UID
		uid := s.extractUidFromKey(key)
		if uid == "" {
			continue
		}

		// 检查用户是否有最近的登录活动
		hasRecentLogin, err := s.HasRecentLogin(ctx, uid, 30*24*time.Hour)
		if err != nil {
			utils.LogWarn(nil, "检查用户登录状态失败: %v", err)
			continue
		}

		// 如果用户30天内没有登录，删除钱包缓存
		if !hasRecentLogin {
			err = database.GetGlobalRedisHelper().Del(ctx, key)
			if err != nil {
				utils.LogWarn(nil, "删除过期钱包缓存失败: %v", err)
			} else {
				expiredCount++
			}
		}
	}

	utils.LogInfo(nil, "清理过期钱包缓存完成，删除 %d 个缓存", expiredCount)
	return nil
}

// 从缓存Key中提取UID
func (s *WalletCacheService) extractUidFromKey(key string) string {
	// Key格式为 wallet:balance:uid 或 wallet:empty:uid
	// 使用冒号分割，取最后一部分作为UID
	parts := strings.Split(key, ":")
	if len(parts) >= 3 {
		return parts[len(parts)-1]
	}
	return ""
}

// 删除钱包缓存
func (s *WalletCacheService) DeleteWalletBalance(ctx context.Context, uid string) error {
	if uid == "" {
		return utils.NewAppError(utils.CodeInvalidParams, "用户ID不能为空")
	}

	// 生成缓存Key
	cacheKey := s.generateWalletKey(uid)
	emptyKey := s.generateEmptyKey(uid)

	// 删除钱包缓存
	if err := database.DelKey(ctx, cacheKey); err != nil {
		utils.LogWarn(nil, "删除钱包缓存失败: %v", err)
	}

	// 删除空值缓存
	if err := database.DelKey(ctx, emptyKey); err != nil {
		utils.LogWarn(nil, "删除空值缓存失败: %v", err)
	}

	return nil
}

// 清理互斥锁映射
func (s *WalletCacheService) CleanupMutexMap() {
	s.mutexMap = sync.Map{}
}
