# 钱包余额缓存设计方案

## 概述

用户钱包余额缓存方案，使用Redis缓存高频查询的钱包数据，提升查询性能。

## 缓存结构

### Key格式
- **钱包余额**: `wallet_balance_{uid}`
- **钱包交易**: `wallet_transactions_{uid}`
- **数据类型**: String (JSON格式)

### 示例
```
wallet_balance_user123: {"uid":"user123","balance":1000.50,"status":1,"currency":"PHP",...}
wallet_transactions_user123: [{"transaction_no":"TX202412011200001",...}]
```

## 功能特性

### 1. 自动缓存
- 用户查询钱包时自动缓存钱包数据
- 钱包余额更新时自动同步缓存
- 缓存失败不影响主业务流程

### 2. 缓存管理
- **钱包余额过期时间**: 30分钟
- **交易记录过期时间**: 15分钟
- **数据格式**: JSON序列化存储
- **错误处理**: 缓存操作失败不影响主流程

### 3. 查询功能
- 优先从缓存获取钱包余额
- 缓存未命中时从数据库获取并缓存
- 支持钱包余额和交易记录的缓存

## 使用方法

### 1. 获取钱包余额（自动缓存）
```go
// 在 services/wallet_service.go 的 GetWallet 方法中
// 优先从缓存获取，未命中时从数据库获取并缓存
wallet, err := s.GetWallet(uid)
```

### 2. 手动缓存钱包余额
```go
cacheService := NewWalletCacheService()
err := cacheService.CacheWalletBalance(ctx, wallet)
```

### 3. 获取缓存的钱包余额
```go
cacheService := NewWalletCacheService()
wallet, err := cacheService.GetCachedWalletBalance(ctx, uid)
```

### 4. 更新钱包余额缓存
```go
cacheService := NewWalletCacheService()
err := cacheService.UpdateWalletBalance(ctx, uid, newBalance)
```

### 5. 删除钱包余额缓存
```go
cacheService := NewWalletCacheService()
err := cacheService.DeleteWalletBalance(ctx, uid)
```

### 6. 检查钱包是否已缓存
```go
cacheService := NewWalletCacheService()
exists, err := cacheService.IsWalletCached(ctx, uid)
```

### 7. 缓存钱包交易记录
```go
cacheService := NewWalletCacheService()
err := cacheService.CacheWalletTransactions(ctx, uid, transactions)
```

### 8. 获取缓存的钱包交易记录
```go
cacheService := NewWalletCacheService()
transactions, err := cacheService.GetCachedWalletTransactions(ctx, uid)
```

### 9. 批量缓存钱包余额
```go
cacheService := NewWalletCacheService()
err := cacheService.BatchCacheWalletBalances(ctx, wallets)
```

### 10. 预热钱包缓存
```go
cacheService := NewWalletCacheService()
err := cacheService.WarmUpWalletCache(ctx, []string{"user1", "user2", "user3"})
```

## 缓存数据字段

### 钱包余额缓存数据
```json
{
    "uid": "user123",
    "balance": 1000.50,
    "status": 1,
    "currency": "PHP",
    "created_at": "2024-12-01T10:00:00Z",
    "updated_at": "2024-12-01T10:30:00Z"
}
```

### 钱包交易记录缓存数据
```json
[
    {
        "transaction_no": "TX202412011200001",
        "uid": "user123",
        "type": "recharge",
        "amount": 100.00,
        "balance_before": 900.50,
        "balance_after": 1000.50,
        "status": "success",
        "description": "充值",
        "created_at": "2024-12-01T12:00:00Z"
    }
]
```

## 性能优化

### 1. 缓存策略
- **LRU淘汰**: Redis自动管理内存
- **过期时间**: 避免数据过期问题
- **批量操作**: 支持批量缓存和查询

### 2. 错误处理
- 缓存操作失败不影响主业务流程
- 序列化失败时跳过该数据，继续处理其他数据
- 记录详细的错误日志

### 3. 数据一致性
- 钱包余额更新时同步更新缓存
- 支持手动删除缓存，强制从数据库获取最新数据

## 监控和维护

### 1. 缓存命中率监控
```go
// 检查缓存是否存在
exists, err := cacheService.IsWalletCached(ctx, uid)
if !exists {
    // 缓存未命中，从数据库查询
    wallet, err := walletRepo.FindWalletByUid(ctx, uid)
    // 然后缓存到Redis
    cacheService.CacheWalletBalance(ctx, wallet)
}
```

### 2. 缓存清理
- 定期清理过期数据
- 手动清理指定用户的缓存
- 支持批量清理操作

### 3. 数据一致性保证
- 钱包余额更新时同步更新缓存
- 支持缓存失效策略
- 提供手动刷新缓存接口

## 扩展功能

### 1. 缓存统计
```go
// 获取缓存统计信息
func GetWalletCacheStats(ctx context.Context) map[string]interface{} {
    cacheService := NewWalletCacheService()
    return cacheService.GetCacheStats(ctx)
}
```

### 2. 缓存预热
```go
// 系统启动时预热活跃用户的钱包数据
func WarmUpActiveUserWallets(ctx context.Context) error {
    // 获取活跃用户列表
    activeUsers := getActiveUsers()
    
    // 预热钱包缓存
    cacheService := NewWalletCacheService()
    return cacheService.WarmUpWalletCache(ctx, activeUsers)
}
```

### 3. 缓存监控
```go
// 监控缓存性能
func MonitorWalletCache(ctx context.Context) {
    // 监控缓存命中率
    // 监控缓存大小
    // 监控缓存响应时间
}
```

## 使用场景

### 1. 高频查询场景
- 用户查看钱包余额
- 订单创建时检查余额
- 提现申请时验证余额

### 2. 实时更新场景
- 充值成功后更新余额
- 提现成功后更新余额
- 订单支付后更新余额

### 3. 批量操作场景
- 批量查询用户钱包余额
- 批量更新钱包状态
- 批量清理过期缓存

## 注意事项

1. **数据一致性**: 缓存数据可能与数据库数据不一致，需要定期同步
2. **内存使用**: 大量用户钱包数据可能占用较多内存，需要合理设置过期时间
3. **网络延迟**: Redis操作可能增加网络延迟，需要权衡缓存收益
4. **错误处理**: 缓存操作失败不应影响主业务流程
5. **监控告警**: 需要监控缓存服务的可用性和性能指标
6. **数据安全**: 钱包余额是敏感数据，需要确保缓存安全性

## 最佳实践

1. **合理设置过期时间**: 根据业务需求设置合适的过期时间
2. **监控缓存命中率**: 定期监控缓存效果，优化缓存策略
3. **处理缓存穿透**: 对于不存在的用户，可以缓存空值避免穿透
4. **处理缓存雪崩**: 设置随机过期时间，避免同时过期
5. **定期清理**: 定期清理过期和无效的缓存数据 