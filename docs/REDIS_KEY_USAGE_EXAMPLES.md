# Redis Key 使用示例

## 概述

本文档展示了项目中各种Redis Key的使用示例，帮助开发者快速理解和使用Redis Key管理工具。

## 钱包相关Key

### 1. 钱包分布式锁
```go
// 生成钱包锁Key
lockKey := utils.RedisKeys.GenerateWalletLockKey("user123")
// 结果: wallet:lock:user123

// 使用示例：在分布式锁服务中
func (s *WalletDistributedLockService) acquireLock(ctx context.Context, uid string) (string, error) {
    lockKey := utils.RedisKeys.GenerateWalletLockKey(uid)
    // 使用Redis SET NX EX命令获取锁
    success, err := database.RedisClient.SetNX(ctx, lockKey, lockValue, 30*time.Second).Result()
    // ...
}
```

### 2. 钱包余额缓存
```go
// 生成钱包余额缓存Key
balanceKey := utils.RedisKeys.GenerateWalletBalanceKey("user123")
// 结果: wallet:balance:user123

// 使用示例：缓存钱包余额
func CacheWalletBalance(ctx context.Context, wallet *models.Wallet) error {
    cacheKey := utils.RedisKeys.GenerateWalletBalanceKey(wallet.Uid)
    walletJSON, _ := json.Marshal(wallet)
    return database.SetKey(ctx, cacheKey, string(walletJSON), 0) // 不过期
}
```

### 3. 钱包空值缓存
```go
// 生成钱包空值缓存Key
emptyKey := utils.RedisKeys.GenerateWalletEmptyKey("user123")
// 结果: wallet:empty:user123

// 使用示例：防止缓存穿透
func CacheEmptyWallet(ctx context.Context, uid string) error {
    emptyKey := utils.RedisKeys.GenerateWalletEmptyKey(uid)
    return database.SetKey(ctx, emptyKey, "empty", 10*time.Minute)
}
```

## 订单相关Key

### 1. 订单缓存
```go
// 生成订单缓存Key
cacheKey := utils.RedisKeys.GenerateOrderCacheKey("order_12345")
// 结果: order:cache:order_12345

// 使用示例：缓存订单数据
func CacheOrder(ctx context.Context, order *models.Order) error {
    cacheKey := utils.RedisKeys.GenerateOrderCacheKey(order.OrderID)
    orderJSON, _ := json.Marshal(order)
    return database.SetKey(ctx, cacheKey, string(orderJSON), 1*time.Hour)
}
```

### 2. 订单锁
```go
// 生成订单锁Key
lockKey := utils.RedisKeys.GenerateOrderLockKey("order_12345")
// 结果: order:lock:order_12345

// 使用示例：订单操作锁
func ProcessOrder(ctx context.Context, orderID string) error {
    lockKey := utils.RedisKeys.GenerateOrderLockKey(orderID)
    // 获取订单锁，防止并发处理
    success, err := database.RedisClient.SetNX(ctx, lockKey, "locked", 30*time.Second).Result()
    // ...
}
```

## 用户相关Key

### 1. 邮箱存在检查
```go
// 生成邮箱存在检查Key
emailKey := utils.RedisKeys.GenerateEmailExistsKey("test@example.com")
// 结果: email:test@example.com:exists

// 使用示例：缓存邮箱存在检查结果
func CheckEmailExists(ctx context.Context, email string) (bool, error) {
    emailKey := utils.RedisKeys.GenerateEmailExistsKey(email)
    exists, err := database.GetKey(ctx, emailKey)
    if err == nil && exists == "1" {
        return true, nil
    }
    // 查询数据库并缓存结果
    // ...
}
```

### 2. 用户名存在检查
```go
// 生成用户名存在检查Key
usernameKey := utils.RedisKeys.GenerateUsernameExistsKey("john_doe")
// 结果: username:john_doe:exists

// 使用示例：缓存用户名存在检查结果
func CheckUsernameExists(ctx context.Context, username string) (bool, error) {
    usernameKey := utils.RedisKeys.GenerateUsernameExistsKey(username)
    exists, err := database.GetKey(ctx, usernameKey)
    if err == nil && exists == "1" {
        return true, nil
    }
    // 查询数据库并缓存结果
    // ...
}
```

## 排行榜相关Key

### 1. 排行榜数据
```go
// 生成排行榜Key
leaderboardKey := utils.RedisKeys.GenerateLeaderboardKey("daily")
// 结果: leaderboard:daily

// 使用示例：使用Redis有序集合存储排行榜
func UpdateLeaderboard(ctx context.Context, leaderboardType string, userID string, score float64) error {
    leaderboardKey := utils.RedisKeys.GenerateLeaderboardKey(leaderboardType)
    return database.RedisClient.ZAdd(ctx, leaderboardKey, redis.Z{
        Score:  score,
        Member: userID,
    }).Err()
}
```

### 2. 排行榜锁
```go
// 生成排行榜锁Key
lockKey := utils.RedisKeys.GenerateLeaderboardLockKey("daily")
// 结果: leaderboard:lock:daily

// 使用示例：排行榜更新锁
func UpdateLeaderboardWithLock(ctx context.Context, leaderboardType string) error {
    lockKey := utils.RedisKeys.GenerateLeaderboardLockKey(leaderboardType)
    // 获取排行榜更新锁
    success, err := database.RedisClient.SetNX(ctx, lockKey, "updating", 60*time.Second).Result()
    // ...
}
```

## 限流相关Key

### 1. API限流
```go
// 生成限流Key
rateLimitKey := utils.RedisKeys.GenerateRateLimitKey("192.168.1.1", "1m")
// 结果: rate_limit:192.168.1.1:1m

// 使用示例：实现API限流
func CheckRateLimit(ctx context.Context, identifier string, window string, limit int) (bool, error) {
    rateLimitKey := utils.RedisKeys.GenerateRateLimitKey(identifier, window)
    
    // 获取当前时间窗口内的请求次数
    count, err := database.RedisClient.Get(ctx, rateLimitKey).Int()
    if err != nil && err != redis.Nil {
        return false, err
    }
    
    if count >= limit {
        return false, nil // 超过限制
    }
    
    // 增加计数
    pipe := database.RedisClient.Pipeline()
    pipe.Incr(ctx, rateLimitKey)
    pipe.Expire(ctx, rateLimitKey, 1*time.Minute)
    _, err = pipe.Exec(ctx)
    
    return true, err
}
```

## 全局Key

### 1. 全局锁
```go
// 生成全局锁Key
globalLockKey := utils.RedisKeys.GenerateGlobalLockKey("system_maintenance")
// 结果: global:lock:system_maintenance

// 使用示例：系统维护锁
func StartSystemMaintenance(ctx context.Context) error {
    lockKey := utils.RedisKeys.GenerateGlobalLockKey("system_maintenance")
    success, err := database.RedisClient.SetNX(ctx, lockKey, "maintenance", 30*time.Minute).Result()
    if !success {
        return errors.New("系统维护已在进行中")
    }
    // 执行系统维护
    // ...
}
```

### 2. 全局计数器
```go
// 生成全局计数器Key
counterKey := utils.RedisKeys.GenerateGlobalCounterKey("total_users")
// 结果: global:counter:total_users

// 使用示例：用户总数统计
func IncrementUserCount(ctx context.Context) error {
    counterKey := utils.RedisKeys.GenerateGlobalCounterKey("total_users")
    return database.RedisClient.Incr(ctx, counterKey).Err()
}
```

## 批量操作

### 1. 批量清理缓存
```go
// 获取钱包相关Key模式
pattern := utils.RedisKeys.GetWalletKeyPattern()
// 结果: wallet:*

// 使用示例：批量清理钱包缓存
func CleanupWalletCache(ctx context.Context) error {
    pattern := utils.RedisKeys.GetWalletKeyPattern()
    var cursor uint64
    for {
        var keys []string
        keys, cursor, err := database.RedisClient.Scan(ctx, cursor, pattern, 100).Result()
        if err != nil {
            return err
        }
        
        if len(keys) > 0 {
            err := database.RedisClient.Del(ctx, keys...).Err()
            if err != nil {
                return err
            }
        }
        
        if cursor == 0 {
            break
        }
    }
    return nil
}
```

## 最佳实践

### 1. Key命名规范
- 使用冒号分隔不同层级
- 保持命名的一致性和可读性
- 避免使用特殊字符

### 2. 过期时间设置
- 根据数据更新频率设置合适的过期时间
- 重要数据可以设置较长的过期时间或不设置过期时间
- 临时数据设置较短的过期时间

### 3. 错误处理
- 缓存操作失败不应影响主业务流程
- 记录缓存操作的错误日志
- 提供降级策略

### 4. 监控和维护
- 定期清理过期的Key
- 监控Redis的内存使用情况
- 设置Key的TTL告警

## 总结

通过统一的Redis Key管理工具，我们可以：
1. 确保Key命名的一致性
2. 提高代码的可维护性
3. 方便批量操作和监控
4. 减少Key冲突的可能性

建议在项目中统一使用`utils.RedisKeys`来生成所有Redis Key，避免硬编码Key字符串。 