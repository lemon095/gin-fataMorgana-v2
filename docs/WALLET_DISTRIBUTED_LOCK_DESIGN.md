# 钱包分布式锁服务设计文档

## 概述

`WalletDistributedLockService` 是基于Redis的分布式锁钱包服务，支持跨进程的并发安全操作。该服务确保同一用户的余额操作在所有进程间串行化执行，避免数据不一致问题。

## 核心特性

### 1. 分布式锁机制
- **锁的Key**: `wallet_lock_{uid}` - 基于用户UID生成
- **锁的值**: `{timestamp}_{random_string}` - 用于标识锁的持有者
- **锁超时**: 30秒自动过期，防止死锁
- **原子性**: 使用Redis的SET NX EX命令确保原子性获取锁

### 2. 锁的获取和释放
```go
// 获取锁
success, err := database.RedisClient.SetNX(ctx, lockKey, lockValue, timeout).Result()

// 释放锁（使用Lua脚本确保原子性）
luaScript := `
    if redis.call("get", KEYS[1]) == ARGV[1] then
        return redis.call("del", KEYS[1])
    else
        return 0
    end
`
```

### 3. 并发安全操作
- **扣减余额**: `WithdrawBalance()`
- **增加余额**: `AddBalance()`
- **转账操作**: `TransferBalance()`
- **批量操作**: `BatchBalanceOperation()`

## 使用场景

### 1. 单用户操作
```go
// 扣减余额
err := walletService.WithdrawBalance(ctx, "user123", 100.0, "购买商品")

// 增加余额
err := walletService.AddBalance(ctx, "user123", 50.0, "退款")
```

### 2. 转账操作
```go
// 用户A向用户B转账
err := walletService.TransferBalance(ctx, "userA", "userB", 100.0, "转账")
```

### 3. 批量操作
```go
operations := []BalanceOperation{
    {UID: "user1", Type: "withdraw", Amount: 50.0},
    {UID: "user2", Type: "add", Amount: 30.0},
}
err := walletService.BatchBalanceOperation(ctx, operations)
```

## 技术实现

### 1. 锁的获取策略
- 按用户UID排序获取锁，避免死锁
- 转账时先获取UID较小的锁，再获取UID较大的锁
- 获取锁失败时立即释放已获取的锁

### 2. 数据一致性保证
1. **获取分布式锁**
2. **从数据库获取最新数据**
3. **执行余额操作**
4. **更新数据库**
5. **更新缓存**
6. **释放分布式锁**

### 3. 错误处理
- 锁获取失败：返回"系统繁忙"错误
- 锁释放失败：记录警告日志，不影响主流程
- 缓存更新失败：记录警告日志，不影响主流程

## 性能考虑

### 1. 锁的粒度
- 按用户级别加锁，不同用户的操作可以并行
- 同一用户的多个操作串行执行

### 2. 锁的超时时间
- 默认30秒超时，平衡性能和安全性
- 可根据业务复杂度调整超时时间

### 3. 缓存策略
- 优先从缓存读取，减少数据库压力
- 缓存更新失败不影响主流程

## 监控和运维

### 1. 锁的监控
```go
// 清理过期的锁
func (s *WalletDistributedLockService) CleanupExpiredLocks(ctx context.Context) error {
    // 实现清理逻辑
    return nil
}
```

### 2. 日志记录
- 记录所有余额操作的详细信息
- 记录锁的获取和释放状态
- 记录缓存更新失败的情况

### 3. 错误码
- `CodeRedisError`: Redis操作失败
- `CodeSystemBusy`: 系统繁忙，锁获取失败
- `CodeBalanceInsufficient`: 余额不足
- `CodeWalletFrozenWithdraw`: 钱包冻结，无法扣减
- `CodeWalletFrozenRecharge`: 钱包冻结，无法充值

## 与现有服务的对比

### 1. 内存锁 vs 分布式锁
| 特性 | 内存锁 | 分布式锁 |
|------|--------|----------|
| 并发范围 | 单进程 | 跨进程 |
| 性能 | 更高 | 稍低 |
| 可靠性 | 进程重启丢失 | 持久化 |
| 适用场景 | 单机部署 | 分布式部署 |

### 2. 迁移建议
- 单机环境：继续使用 `WalletConcurrentService`
- 分布式环境：使用 `WalletDistributedLockService`
- 混合环境：根据部署方式选择对应服务

## 测试用例

### 1. 并发测试
```bash
# 测试多个进程同时操作同一用户钱包
./test_wallet_distributed_lock.sh
```

### 2. 死锁测试
```bash
# 测试转账操作的死锁预防
./test_wallet_transfer_deadlock.sh
```

### 3. 故障恢复测试
```bash
# 测试Redis故障时的行为
./test_wallet_redis_failure.sh
```

## 注意事项

1. **Redis可用性**: 分布式锁依赖Redis，需要确保Redis的高可用
2. **网络延迟**: 跨进程锁操作会有网络延迟，影响性能
3. **锁超时**: 需要合理设置锁超时时间，避免死锁
4. **监控告警**: 建议监控锁的获取失败率和超时情况
5. **降级策略**: 在Redis不可用时，可以考虑降级到数据库级别的锁

## 总结

`WalletDistributedLockService` 提供了跨进程的并发安全保证，适用于分布式部署环境。通过Redis分布式锁机制，确保同一用户的余额操作在所有进程间串行化执行，有效避免了数据不一致问题。 