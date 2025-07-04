package utils

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// IdempotencyKey 幂等键结构
type IdempotencyKey struct {
	Key       string    `json:"key"`
	UserID    string    `json:"user_id"`
	Operation string    `json:"operation"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IdempotencyManager 幂等性管理器
type IdempotencyManager struct {
	keys  map[string]*IdempotencyKey
	mutex sync.RWMutex
}

// NewIdempotencyManager 创建新的幂等性管理器
func NewIdempotencyManager() *IdempotencyManager {
	return &IdempotencyManager{
		keys: make(map[string]*IdempotencyKey),
	}
}

// GenerateIdempotencyKey 生成幂等键
func (im *IdempotencyManager) GenerateIdempotencyKey(userID, operation string) string {
	// 使用时间戳+用户ID+操作类型+随机字符串生成唯一键
	timestamp := time.Now().UTC().UnixNano()
	randomStr := RandomString(8)
	return fmt.Sprintf("%d_%s_%s_%s", timestamp, userID, operation, randomStr)
}

// CheckAndSetKey 检查并设置幂等键
func (im *IdempotencyManager) CheckAndSetKey(key, userID, operation string, expireMinutes int) (bool, error) {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	// 检查键是否已存在
	if existingKey, exists := im.keys[key]; exists {
		// 检查是否过期
		if time.Now().Before(existingKey.ExpiresAt) {
			return false, fmt.Errorf("重复请求，幂等键已存在: %s", key)
		}
		// 已过期，删除旧记录
		delete(im.keys, key)
	}

	// 设置新键
	expiresAt := time.Now().Add(time.Duration(expireMinutes) * time.Minute)
	im.keys[key] = &IdempotencyKey{
		Key:       key,
		UserID:    userID,
		Operation: operation,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	return true, nil
}

// RemoveKey 移除幂等键
func (im *IdempotencyManager) RemoveKey(key string) {
	im.mutex.Lock()
	defer im.mutex.Unlock()
	delete(im.keys, key)
}

// CleanExpiredKeys 清理过期的键
func (im *IdempotencyManager) CleanExpiredKeys() {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	now := time.Now()
	for key, idempotencyKey := range im.keys {
		if now.After(idempotencyKey.ExpiresAt) {
			delete(im.keys, key)
		}
	}
}

// 全局幂等性管理器实例
var globalIdempotencyManager = NewIdempotencyManager()

// GenerateTransactionNo 生成交易流水号
func GenerateTransactionNo(operation string) string {
	// 格式：TX + 年月日时分秒 + 4位随机数 + 操作类型
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := RandomString(4)
	return fmt.Sprintf("TX%s%s%s", timestamp, random, operation[:3])
}

// GenerateWithdrawNo 生成提现流水号
func GenerateWithdrawNo() string {
	return GenerateTransactionNo("WITHDRAW")
}

// GenerateRechargeNo 生成充值流水号
func GenerateRechargeNo() string {
	return GenerateTransactionNo("RECHARGE")
}

// GenerateTransferNo 生成转账流水号
func GenerateTransferNo() string {
	return GenerateTransactionNo("TRANSFER")
}

// CheckIdempotency 检查幂等性（简化版，用于内存存储）
func CheckIdempotency(key, userID, operation string) (bool, error) {
	return globalIdempotencyManager.CheckAndSetKey(key, userID, operation, 30) // 30分钟过期
}

// RemoveIdempotencyKey 移除幂等键
func RemoveIdempotencyKey(key string) {
	globalIdempotencyManager.RemoveKey(key)
}

// StartIdempotencyCleaner 启动幂等键清理器
func StartIdempotencyCleaner(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟清理一次
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			globalIdempotencyManager.CleanExpiredKeys()
		}
	}
}
