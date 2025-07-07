# 钱包并发优化方案

## 问题背景

在批量发奖等高频并发场景下，原有的分布式锁机制会直接返回"系统繁忙"错误，导致用户体验不佳。需要优化锁获取策略，支持等待和重试机制。

## 优化方案

### 1. 重试机制优化

#### 原有机制
```go
// 直接失败，不等待
func (s *WalletService) acquireLock(ctx context.Context, uid string, timeout time.Duration) (string, error) {
    success, err := database.GlobalRedisHelper.SetNX(ctx, lockKey, lockValue, timeout)
    if !success {
        return "", utils.NewAppError(utils.CodeSystemBusy, "系统繁忙，请稍后重试")
    }
    return lockValue, nil
}
```

#### 优化后机制
```go
// 支持重试和等待
func (s *WalletService) acquireLockWithRetry(ctx context.Context, uid string, timeout time.Duration, maxRetries int, retryDelay time.Duration) (string, error) {
    for attempt := 0; attempt <= maxRetries; attempt++ {
        success, err := database.GlobalRedisHelper.SetNX(ctx, lockKey, lockValue, timeout)
        if success {
            return lockValue, nil
        }
        
        if attempt < maxRetries {
            // 指数退避重试
            select {
            case <-ctx.Done():
                return "", utils.NewAppError(utils.CodeRequestTimeout, "请求超时")
            case <-time.After(retryDelay):
                retryDelay = time.Duration(float64(retryDelay) * 1.5)
            }
        }
    }
    return "", utils.NewAppError(utils.CodeSystemBusy, "系统繁忙，请稍后重试")
}
```

### 2. 不同场景的重试策略

#### 普通操作（默认）
- 重试次数：3次
- 初始延迟：100ms
- 指数退避：1.5倍递增
- 总等待时间：约1.5秒

#### 批量操作
- 重试次数：5次
- 初始延迟：200ms
- 指数退避：1.5倍递增
- 总等待时间：约3秒

#### 批量发奖（特殊优化）
- 重试次数：10次
- 初始延迟：500ms
- 指数退避：1.5倍递增
- 总等待时间：约15秒

### 3. 批量发奖专用方法

```go
// 批量发奖优化方法
func (s *WalletService) BatchAddBalanceForRewards(ctx context.Context, rewards []struct {
    UID     string  `json:"uid"`
    Amount  float64 `json:"amount"`
    Desc    string  `json:"description"`
}) error {
    // 按用户分组，避免重复获取锁
    userRewards := make(map[string][]Reward)
    
    // 并发处理不同用户
    var wg sync.WaitGroup
    errChan := make(chan error, len(userRewards))
    
    for uid, rewardList := range userRewards {
        wg.Add(1)
        go func(userUid string, rewards []Reward) {
            defer wg.Done()
            // 一次性处理用户的所有奖励
            s.executeUserRewards(ctx, userUid, rewards)
        }(uid, rewardList)
    }
    
    wg.Wait()
    // 检查错误...
}
```

### 4. 锁状态监控

#### 获取锁状态
```go
func (s *WalletService) GetLockStatus(ctx context.Context, uid string) (map[string]interface{}, error) {
    // 返回锁的详细信息：
    // - 是否存在
    // - 剩余过期时间
    // - 锁持有者标识
    // - 时间戳
}
```

#### 强制释放锁（紧急情况）
```go
func (s *WalletService) ForceReleaseLock(ctx context.Context, uid string) error {
    // 直接删除锁，用于紧急情况
    // 会记录警告日志
}
```

## 使用建议

### 1. 普通用户操作
```go
// 使用默认重试策略
err := walletService.WithdrawBalance(ctx, uid, amount, description)
```

### 2. 批量操作
```go
// 使用批量操作接口
operations := []BalanceOperation{
    {UID: "user1", Type: "add", Amount: 100},
    {UID: "user2", Type: "withdraw", Amount: 50},
}
err := walletService.BatchBalanceOperation(ctx, operations)
```

### 3. 批量发奖
```go
// 使用专门的批量发奖接口
rewards := []struct {
    UID     string  `json:"uid"`
    Amount  float64 `json:"amount"`
    Desc    string  `json:"description"`
}{
    {UID: "user1", Amount: 100, Desc: "活动奖励"},
    {UID: "user2", Amount: 200, Desc: "活动奖励"},
}
err := walletService.BatchAddBalanceForRewards(ctx, rewards)
```

### 4. 监控和调试
```go
// 检查锁状态
status, err := walletService.GetLockStatus(ctx, uid)
if err != nil {
    // 处理错误
}

// 紧急情况下强制释放锁
err = walletService.ForceReleaseLock(ctx, uid)
```

## 性能优化效果

### 1. 成功率提升
- 原有机制：批量发奖时大量失败
- 优化后：通过重试机制大幅提升成功率

### 2. 用户体验改善
- 原有机制：立即返回"系统繁忙"
- 优化后：自动等待和重试，用户无感知

### 3. 系统稳定性
- 保持数据一致性
- 避免死锁
- 支持监控和紧急处理

## 注意事项

### 1. 超时控制
- 所有操作都有Context超时控制
- 避免无限等待

### 2. 资源消耗
- 重试会增加Redis访问次数
- 需要监控Redis性能

### 3. 错误处理
- 重试失败后仍会返回错误
- 客户端需要处理最终失败情况

### 4. 监控告警
- 建议监控锁获取失败率
- 监控重试次数分布
- 设置合理的告警阈值

## 配置建议

### 1. 重试参数调优
根据实际业务场景调整：
- 高频场景：增加重试次数和延迟
- 低频场景：减少重试次数和延迟

### 2. 锁超时时间
- 普通操作：30秒
- 批量操作：60秒
- 根据操作复杂度调整

### 3. 监控指标
- 锁获取成功率
- 平均重试次数
- 锁持有时间分布
- 失败原因统计 