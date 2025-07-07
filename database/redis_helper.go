package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisHelper Redis操作助手
type RedisHelper struct {
	client *redis.Client
}

// NewRedisHelper 创建Redis助手
func NewRedisHelper() *RedisHelper {
	return &RedisHelper{
		client: RedisClient,
	}
}

// Get 获取值
func (r *RedisHelper) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set 设置值
func (r *RedisHelper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// SetNX 设置值（如果不存在）
func (r *RedisHelper) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

// Del 删除键
func (r *RedisHelper) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisHelper) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func (r *RedisHelper) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func (r *RedisHelper) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// Incr 递增
func (r *RedisHelper) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// IncrBy 按指定值递增
func (r *RedisHelper) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

// Decr 递减
func (r *RedisHelper) Decr(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

// DecrBy 按指定值递减
func (r *RedisHelper) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.DecrBy(ctx, key, value).Result()
}

// HGet 获取哈希字段值
func (r *RedisHelper) HGet(ctx context.Context, key, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

// HSet 设置哈希字段值
func (r *RedisHelper) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.client.HSet(ctx, key, values...).Err()
}

// HGetAll 获取所有哈希字段
func (r *RedisHelper) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func (r *RedisHelper) HDel(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

// HLen 获取哈希长度
func (r *RedisHelper) HLen(ctx context.Context, key string) (int64, error) {
	return r.client.HLen(ctx, key).Result()
}

// HExists 检查哈希字段是否存在
func (r *RedisHelper) HExists(ctx context.Context, key, field string) (bool, error) {
	return r.client.HExists(ctx, key, field).Result()
}

// LPush 左推入列表
func (r *RedisHelper) LPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.LPush(ctx, key, values...).Err()
}

// RPop 右弹出列表
func (r *RedisHelper) RPop(ctx context.Context, key string) (string, error) {
	return r.client.RPop(ctx, key).Result()
}

// LLen 获取列表长度
func (r *RedisHelper) LLen(ctx context.Context, key string) (int64, error) {
	return r.client.LLen(ctx, key).Result()
}

// ZAdd 添加有序集合成员
func (r *RedisHelper) ZAdd(ctx context.Context, key string, score float64, member string) error {
	return r.client.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err()
}

// ZRange 获取有序集合范围
func (r *RedisHelper) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(ctx, key, start, stop).Result()
}

// ZRevRange 获取有序集合倒序范围
func (r *RedisHelper) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRevRange(ctx, key, start, stop).Result()
}

// ZScore 获取有序集合成员分数
func (r *RedisHelper) ZScore(ctx context.Context, key, member string) (float64, error) {
	return r.client.ZScore(ctx, key, member).Result()
}

// ZRank 获取有序集合成员排名
func (r *RedisHelper) ZRank(ctx context.Context, key, member string) (int64, error) {
	return r.client.ZRank(ctx, key, member).Result()
}

// ZRevRank 获取有序集合成员倒序排名
func (r *RedisHelper) ZRevRank(ctx context.Context, key, member string) (int64, error) {
	return r.client.ZRevRank(ctx, key, member).Result()
}

// SAdd 添加集合成员
func (r *RedisHelper) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SAdd(ctx, key, members...).Err()
}

// SRem 删除集合成员
func (r *RedisHelper) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SRem(ctx, key, members...).Err()
}

// SIsMember 检查集合成员
func (r *RedisHelper) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return r.client.SIsMember(ctx, key, member).Result()
}

// SMembers 获取集合所有成员
func (r *RedisHelper) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

// SCard 获取集合成员数量
func (r *RedisHelper) SCard(ctx context.Context, key string) (int64, error) {
	return r.client.SCard(ctx, key).Result()
}

// GetJSON 获取JSON值
func (r *RedisHelper) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// SetJSON 设置JSON值
func (r *RedisHelper) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Set(ctx, key, string(data), expiration)
}

// GetOrSet 获取值，如果不存在则设置默认值
func (r *RedisHelper) GetOrSet(ctx context.Context, key string, defaultValue interface{}, expiration time.Duration) (string, error) {
	value, err := r.Get(ctx, key)
	if err == redis.Nil {
		// 键不存在，设置默认值
		if err := r.Set(ctx, key, defaultValue, expiration); err != nil {
			return "", err
		}
		return fmt.Sprintf("%v", defaultValue), nil
	}
	return value, err
}

// GetOrSetJSON 获取JSON值，如果不存在则设置默认值
func (r *RedisHelper) GetOrSetJSON(ctx context.Context, key string, defaultValue interface{}, dest interface{}, expiration time.Duration) error {
	err := r.GetJSON(ctx, key, dest)
	if err == redis.Nil {
		// 键不存在，设置默认值
		if err := r.SetJSON(ctx, key, defaultValue, expiration); err != nil {
			return err
		}
		// 将默认值复制到dest
		data, _ := json.Marshal(defaultValue)
		return json.Unmarshal(data, dest)
	}
	return err
}

// Pipeline 执行管道操作
func (r *RedisHelper) Pipeline(ctx context.Context, fn func(redis.Pipeliner) error) error {
	pipe := r.client.Pipeline()
	if err := fn(pipe); err != nil {
		return err
	}
	_, err := pipe.Exec(ctx)
	return err
}

// Transaction 执行事务
func (r *RedisHelper) Transaction(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	return r.client.Watch(ctx, fn, keys...)
}

// Lock 分布式锁
func (r *RedisHelper) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.SetNX(ctx, key, "1", expiration)
}

// Unlock 释放分布式锁
func (r *RedisHelper) Unlock(ctx context.Context, key string) error {
	return r.Del(ctx, key)
}

// 全局Redis助手实例（延迟初始化）
var GlobalRedisHelper *RedisHelper

// GetGlobalRedisHelper 获取全局Redis助手实例（延迟初始化）
func GetGlobalRedisHelper() *RedisHelper {
	if GlobalRedisHelper == nil {
		if RedisClient == nil {
			panic("Redis客户端未初始化，请先调用 InitRedis()")
		}
		GlobalRedisHelper = &RedisHelper{client: RedisClient}
	}
	return GlobalRedisHelper
} 