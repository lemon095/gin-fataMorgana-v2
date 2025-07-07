package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"gin-fataMorgana/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() error {
	cfg := config.GlobalConfig.Redis

	// 创建Redis客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         cfg.GetRedisAddr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     10,              // 连接池大小
		MinIdleConns: 5,               // 最小空闲连接数
		MaxRetries:   3,               // 最大重试次数
		DialTimeout:  5 * time.Second, // 连接超时
		ReadTimeout:  3 * time.Second, // 读取超时
		WriteTimeout: 3 * time.Second, // 写入超时
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis连接测试失败: %w", err)
	}

	log.Println("Redis连接成功")
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// SetKey 设置键值对
func SetKey(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// GetKey 获取键值
func GetKey(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// GetKeyOrDefault 获取键值，如果键不存在则返回默认值
func GetKeyOrDefault(ctx context.Context, key string, defaultValue string) (string, error) {
	value, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return defaultValue, nil // 键不存在时返回默认值
		}
		return "", err // 其他错误仍然返回错误
	}
	return value, nil
}

// DelKey 删除键
func DelKey(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// ExistsKey 检查键是否存在
func ExistsKey(ctx context.Context, key string) (bool, error) {
	result, err := RedisClient.Exists(ctx, key).Result()
	return result > 0, err
}

// Keys 获取匹配模式的键
func Keys(ctx context.Context, pattern string) ([]string, error) {
	return RedisClient.Keys(ctx, pattern).Result()
}

// SetExpire 设置过期时间
func SetExpire(ctx context.Context, key string, expiration time.Duration) error {
	return RedisClient.Expire(ctx, key, expiration).Err()
}

// CacheEmailExists 缓存邮箱存在检查结果
func CacheEmailExists(ctx context.Context, email string, exists bool) error {
	key := fmt.Sprintf("email_exists:%s", email)
	value := "1"
	if !exists {
		value = "0"
	}

	// 缓存5分钟
	return RedisClient.Set(ctx, key, value, 5*time.Minute).Err()
}

// GetCachedEmailExists 获取缓存的邮箱存在检查结果
func GetCachedEmailExists(ctx context.Context, email string) (bool, bool, error) {
	key := fmt.Sprintf("email_exists:%s", email)
	value, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, false, nil // 缓存不存在
		}
		return false, false, err
	}

	exists := value == "1"
	return exists, true, nil // 第二个返回值表示缓存命中
}

// InvalidateEmailCache 用户创建后清除邮箱缓存
func InvalidateEmailCache(ctx context.Context, email string) error {
	key := fmt.Sprintf("email_exists:%s", email)
	return RedisClient.Del(ctx, key).Err()
}

// CacheUsernameExists 缓存用户名存在检查结果
func CacheUsernameExists(ctx context.Context, username string, exists bool) error {
	key := fmt.Sprintf("username_exists:%s", username)
	value := "1"
	if !exists {
		value = "0"
	}

	// 缓存5分钟
	return RedisClient.Set(ctx, key, value, 5*time.Minute).Err()
}

// GetCachedUsernameExists 获取缓存的用户名存在检查结果
func GetCachedUsernameExists(ctx context.Context, username string) (bool, bool, error) {
	key := fmt.Sprintf("username_exists:%s", username)
	value, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, false, nil // 缓存不存在
		}
		return false, false, err
	}

	exists := value == "1"
	return exists, true, nil // 第二个返回值表示缓存命中
}

// CacheInviteCodeExists 缓存邀请码存在检查结果
func CacheInviteCodeExists(ctx context.Context, inviteCode string, exists bool) error {
	key := fmt.Sprintf("invite_code_exists:%s", inviteCode)
	value := "1"
	if !exists {
		value = "0"
	}

	// 缓存10分钟（邀请码变化较少）
	return RedisClient.Set(ctx, key, value, 10*time.Minute).Err()
}

// GetCachedInviteCodeExists 获取缓存的邀请码存在检查结果
func GetCachedInviteCodeExists(ctx context.Context, inviteCode string) (bool, bool, error) {
	key := fmt.Sprintf("invite_code_exists:%s", inviteCode)
	value, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, false, nil // 缓存不存在
		}
		return false, false, err
	}

	exists := value == "1"
	return exists, true, nil // 第二个返回值表示缓存命中
}

// InvalidateUsernameCache 用户创建后清除用户名缓存
func InvalidateUsernameCache(ctx context.Context, username string) error {
	key := fmt.Sprintf("username_exists:%s", username)
	return RedisClient.Del(ctx, key).Err()
}

// InvalidateInviteCodeCache 用户创建后清除邀请码缓存
func InvalidateInviteCodeCache(ctx context.Context, inviteCode string) error {
	key := fmt.Sprintf("invite_code_exists:%s", inviteCode)
	return RedisClient.Del(ctx, key).Err()
}
