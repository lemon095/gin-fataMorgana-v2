# 钱包缓存基于用户登录的过期机制

## 概述

本系统实现了基于用户登录状态的钱包缓存过期机制，确保活跃用户的钱包数据始终可用，同时自动清理不活跃用户的缓存以节省内存。

## 设计原理

### 1. 缓存过期策略
- **活跃用户**：钱包缓存过期时间为30天
- **不活跃用户**：如果30天内没有登录活动，缓存会被自动清理
- **登录时延长**：用户每次登录都会重新设置30天的过期时间

### 2. 用户登录时间记录
- 使用Redis Key：`user:login_time:{uid}`
- 记录用户最后登录时间戳
- 过期时间：30天

### 3. 钱包缓存Key
- 钱包余额：`wallet:balance:{uid}`
- 空值缓存：`wallet:empty:{uid}`
- 过期时间：30天

## 核心功能

### 1. 用户登录时延长缓存过期时间

```go
// 用户登录时调用
func (s *WalletCacheService) ExtendWalletCacheOnLogin(ctx context.Context, uid string) error {
    // 1. 更新用户登录时间
    loginKey := s.generateUserLoginKey(uid)
    err := database.GlobalRedisHelper.Set(ctx, loginKey, time.Now().UTC().Unix(), 30*24*time.Hour)
    
    // 2. 获取当前钱包缓存
    wallet, err := s.GetCachedWalletBalance(ctx, uid)
    if err != nil {
        // 如果缓存不存在，从数据库获取并缓存
        walletRepo := database.NewWalletRepository()
        wallet, err = walletRepo.FindWalletByUid(ctx, uid)
    }
    
    // 3. 更新钱包最后活跃时间
    wallet.UpdateLastActive()
    
    // 4. 重新缓存钱包数据（延长过期时间）
    return s.CacheWalletBalance(ctx, wallet)
}
```

### 2. 检查用户登录状态

```go
// 检查用户是否在指定时间内有登录活动
func (s *WalletCacheService) HasRecentLogin(ctx context.Context, uid string, duration time.Duration) (bool, error) {
    loginKey := s.generateUserLoginKey(uid)
    loginTimeStr, err := database.GlobalRedisHelper.Get(ctx, loginKey)
    
    if loginTimeStr == "" {
        return false, nil // 没有登录记录
    }
    
    // 解析登录时间并检查是否在指定时间内
    var loginTimeUnix int64
    json.Unmarshal([]byte(loginTimeStr), &loginTimeUnix)
    loginTime := time.Unix(loginTimeUnix, 0)
    now := time.Now().UTC()
    
    return now.Sub(loginTime) <= duration, nil
}
```

### 3. 自动清理过期缓存

```go
// 清理过期钱包缓存
func (s *WalletCacheService) CleanupExpiredWalletCache(ctx context.Context) error {
    // 获取所有钱包缓存Key
    pattern := utils.RedisKeys.GetWalletKeyPattern()
    keys, err := database.Keys(ctx, pattern)
    
    expiredCount := 0
    for _, key := range keys {
        uid := s.extractUidFromKey(key)
        
        // 检查用户是否有最近的登录活动
        hasRecentLogin, err := s.HasRecentLogin(ctx, uid, 30*24*time.Hour)
        
        // 如果用户30天内没有登录，删除钱包缓存
        if !hasRecentLogin {
            err = database.GlobalRedisHelper.Del(ctx, key)
            if err == nil {
                expiredCount++
            }
        }
    }
    
    return nil
}
```

## API接口

### 1. 用户登录时延长缓存过期时间

```http
POST /api/wallet/extend-cache-on-login
Content-Type: application/json

{
    "uid": "user123"
}
```

**响应：**
```json
{
    "code": 200,
    "message": "钱包缓存过期时间已延长",
    "data": null
}
```

### 2. 清理过期钱包缓存（管理员接口）

```http
POST /api/wallet/cleanup-expired-cache
```

**响应：**
```json
{
    "code": 200,
    "message": "过期钱包缓存清理完成",
    "data": null
}
```

## 使用场景

### 1. 用户登录流程
```go
// 在用户登录成功后调用
func (auth *AuthController) Login(c *gin.Context) {
    // ... 登录验证逻辑 ...
    
    // 登录成功后延长钱包缓存过期时间
    walletService := services.NewWalletService()
    err := walletService.ExtendWalletCacheOnLogin(c, user.Uid)
    if err != nil {
        utils.LogWarn(c, "延长钱包缓存过期时间失败: %v", err)
    }
    
    // ... 返回登录成功响应 ...
}
```

### 2. 定时清理任务
```go
// 在定时任务中调用
func (cron *CronService) CleanupExpiredWalletCache() {
    ctx := context.Background()
    walletService := services.NewWalletService()
    
    err := walletService.CleanupExpiredWalletCache(ctx)
    if err != nil {
        utils.LogError(ctx, "清理过期钱包缓存失败: %v", err)
    } else {
        utils.LogInfo(ctx, "过期钱包缓存清理完成")
    }
}
```

## 配置说明

### Redis Key配置
```go
// utils/redis_key.go
const (
    WALLET_PREFIX = "wallet"  // 钱包相关前缀
    USER_PREFIX   = "user"    // 用户相关前缀
)

// 生成用户登录时间Key
func (r *RedisKeyManager) GenerateUserLoginTimeKey(uid string) string {
    return fmt.Sprintf("%s:login_time:%s", USER_PREFIX, uid)
}

// 生成钱包余额Key
func (r *RedisKeyManager) GenerateWalletBalanceKey(uid string) string {
    return fmt.Sprintf("%s:balance:%s", WALLET_PREFIX, uid)
}
```

### 过期时间配置
```go
const (
    WALLET_CACHE_EXPIRY = 30 * 24 * time.Hour  // 钱包缓存过期时间：30天
    LOGIN_TIME_EXPIRY   = 30 * 24 * time.Hour  // 登录时间记录过期时间：30天
    EMPTY_CACHE_EXPIRY  = 10 * time.Minute     // 空值缓存过期时间：10分钟
)
```

## 优势

### 1. 内存优化
- 自动清理不活跃用户的缓存
- 减少Redis内存占用
- 提高系统整体性能

### 2. 用户体验
- 活跃用户的钱包数据始终可用
- 登录时自动刷新缓存过期时间
- 减少缓存未命中的情况

### 3. 系统稳定性
- 防止缓存无限增长
- 自动化的缓存管理
- 可配置的清理策略

## 注意事项

### 1. 性能考虑
- 清理任务应该在低峰期执行
- 避免频繁的清理操作
- 监控清理任务的执行时间

### 2. 数据一致性
- 清理缓存不影响数据库数据
- 用户重新登录时会重新加载数据
- 缓存更新采用事件驱动机制

### 3. 监控告警
- 监控缓存命中率
- 监控清理任务执行情况
- 设置适当的告警阈值

## 扩展功能

### 1. 自定义过期时间
可以根据用户等级或VIP状态设置不同的缓存过期时间：

```go
func (s *WalletCacheService) getCacheExpiryByUserLevel(uid string) time.Duration {
    // 根据用户等级返回不同的过期时间
    switch userLevel {
    case "vip":
        return 60 * 24 * time.Hour  // VIP用户60天
    case "premium":
        return 45 * 24 * time.Hour  // 高级用户45天
    default:
        return 30 * 24 * time.Hour  // 普通用户30天
    }
}
```

### 2. 批量操作优化
对于批量发奖等场景，可以优化锁的获取策略：

```go
func (s *WalletService) BatchAddBalanceForRewards(ctx context.Context, rewards []Reward) error {
    // 使用更长的重试时间和延迟
    return s.AtomicBalanceOperationWithRetry(ctx, uid, operation, 5, 200*time.Millisecond)
}
```

### 3. 缓存预热
在系统启动时预热活跃用户的缓存：

```go
func (s *WalletService) WarmupActiveUserCache(ctx context.Context) error {
    // 获取最近7天有登录的用户
    activeUsers := s.getActiveUsers(ctx, 7*24*time.Hour)
    
    for _, user := range activeUsers {
        s.GetWalletBalanceWithCache(ctx, user.Uid)
    }
    
    return nil
}
``` 