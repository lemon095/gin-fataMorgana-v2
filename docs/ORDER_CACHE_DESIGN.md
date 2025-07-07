# 订单缓存设计方案

## 概述

基于期数的用户购买数据缓存方案，使用Redis Hash结构存储订单数据。

## 缓存结构

### Key格式
- **大Key**: `fataMorgana_{期号}`
- **小Key**: 订单号 (OrderNo)
- **数据类型**: Hash结构

### 示例
```
fataMorgana_20241201_001
├── ORDER_20241201_001_001: {订单数据JSON}
├── ORDER_20241201_001_002: {订单数据JSON}
├── ORDER_20241201_001_003: {订单数据JSON}
└── ...
```

## 功能特性

### 1. 自动缓存
- 用户通过 `/create` 接口创建订单时，自动缓存订单数据
- 缓存失败不影响主业务流程，只记录日志

### 2. 缓存管理
- **过期时间**: 24小时自动过期
- **数据格式**: JSON序列化存储
- **错误处理**: 缓存操作失败不影响主流程

### 3. 查询功能
- 按期数查询所有订单
- 查询指定订单
- 获取订单数量
- 检查订单是否已缓存

## 使用方法

### 1. 创建订单时自动缓存
```go
// 在 services/order_service.go 的 CreateOrder 方法中
// 订单创建成功后自动缓存
if err := s.cacheOrderData(ctx, order); err != nil {
    // 缓存失败不影响主流程，只记录日志
    fmt.Printf("缓存订单数据失败: %v\n", err)
}
```

### 2. 手动缓存订单
```go
cacheService := NewOrderCacheService()
err := cacheService.CacheOrder(ctx, order)
```

### 3. 查询期数下的所有订单
```go
cacheService := NewOrderCacheService()
orders, err := cacheService.GetOrdersByPeriod(ctx, "20241201_001")
```

### 4. 查询指定订单
```go
cacheService := NewOrderCacheService()
order, err := cacheService.GetOrder(ctx, "20241201_001", "ORDER_20241201_001_001")
```

### 5. 获取订单数量
```go
cacheService := NewOrderCacheService()
count, err := cacheService.GetOrderCount(ctx, "20241201_001")
```

### 6. 检查订单是否已缓存
```go
cacheService := NewOrderCacheService()
exists, err := cacheService.IsOrderCached(ctx, "20241201_001", "ORDER_20241201_001_001")
```

### 7. 删除订单缓存
```go
cacheService := NewOrderCacheService()
// 删除指定订单
err := cacheService.DeleteOrder(ctx, "20241201_001", "ORDER_20241201_001_001")

// 删除期数下的所有订单
err := cacheService.DeletePeriodOrders(ctx, "20241201_001")
```

### 8. 批量缓存订单
```go
cacheService := NewOrderCacheService()
err := cacheService.BatchCacheOrders(ctx, orders)
```

## 缓存数据字段

缓存的订单数据包含以下字段：
```json
{
    "order_no": "ORDER_20241201_001_001",
    "uid": "user123",
    "period_number": "20241201_001",
    "amount": 100.00,
    "profit_amount": 5.00,
    "like_count": 10,
    "share_count": 5,
    "follow_count": 3,
    "favorite_count": 2,
    "status": "pending",
    "expire_time": "2024-12-01T10:05:00Z",
    "created_at": "2024-12-01T10:00:00Z",
    "updated_at": "2024-12-01T10:00:00Z"
}
```

## 性能优化

### 1. 批量操作
- 使用 `HMSet` 进行批量设置
- 使用 `HGetAll` 进行批量获取

### 2. 过期策略
- 设置24小时过期时间，避免内存泄漏
- 过期时间设置失败不影响主流程

### 3. 错误处理
- 缓存操作失败不影响主业务流程
- 序列化失败时跳过该订单，继续处理其他订单

## 监控和维护

### 1. 缓存命中率监控
```go
// 检查缓存是否存在
exists, err := cacheService.IsOrderCached(ctx, periodNumber, orderNo)
if !exists {
    // 缓存未命中，从数据库查询
    order, err := orderRepo.FindOrderByOrderNo(ctx, orderNo)
    // 然后缓存到Redis
    cacheService.CacheOrder(ctx, order)
}
```

### 2. 缓存清理
- 定期清理过期数据
- 手动清理指定期数的缓存

### 3. 数据一致性
- 订单状态更新时同步更新缓存
- 订单删除时同步删除缓存

## 扩展功能

### 1. 缓存预热
```go
// 系统启动时预热热门期数的订单数据
func WarmUpCache(ctx context.Context, periodNumber string) error {
    // 从数据库查询期数下的所有订单
    orders, err := orderRepo.GetOrdersByPeriod(ctx, periodNumber)
    if err != nil {
        return err
    }
    
    // 批量缓存到Redis
    cacheService := NewOrderCacheService()
    return cacheService.BatchCacheOrders(ctx, orders)
}
```

### 2. 缓存统计
```go
// 获取缓存统计信息
func GetCacheStats(ctx context.Context) map[string]interface{} {
    stats := make(map[string]interface{})
    
    // 获取所有期数缓存的数量
    // 获取缓存命中率
    // 获取缓存大小
    
    return stats
}
```

## 注意事项

1. **数据一致性**: 缓存数据可能与数据库数据不一致，需要定期同步
2. **内存使用**: 大量订单数据可能占用较多内存，需要合理设置过期时间
3. **网络延迟**: Redis操作可能增加网络延迟，需要权衡缓存收益
4. **错误处理**: 缓存操作失败不应影响主业务流程
5. **监控告警**: 需要监控缓存服务的可用性和性能指标 