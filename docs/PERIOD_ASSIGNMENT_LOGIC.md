# 期号分配逻辑

## 概述

在假订单生成系统中，期号分配是一个重要的逻辑，需要确保生成的订单分配到正确的期数中。本文档详细说明了期号分配的机制和实现。

## 问题背景

之前的实现中，所有假订单都使用当前期号，这不符合实际情况。真实的订单应该根据其创建时间分配到对应的期数中。

## 解决方案

### 1. 期号缓存优化逻辑

为了优化性能，我们实现了期号缓存机制：

#### 1.1 预加载期数数据

在定时任务开始时，一次性查询时间范围内的所有期数并缓存：

```go
// preloadPeriodData 预加载期数数据到缓存
func (s *FakeOrderService) preloadPeriodData() error {
    ctx := context.Background()
    periodRepo := database.NewLotteryPeriodRepository()
    
    // 清空缓存
    s.periodCache = make(map[string]string)
    
    // 获取当前时间前后30分钟的时间范围
    now := time.Now()
    startTime := now.Add(-30 * time.Minute)
    endTime := now.Add(30 * time.Minute)
    
    // 查询这个时间范围内的所有期数
    periods, err := periodRepo.GetPeriodsByTimeRange(ctx, startTime, endTime)
    if err != nil {
        return err
    }
    
    // 将期数数据缓存到内存中
    for _, period := range periods {
        key := fmt.Sprintf("%s_%s", period.OrderStartTime.Format("2006-01-02 15:04:05"), 
            period.OrderEndTime.Format("2006-01-02 15:04:05"))
        s.periodCache[key] = period.PeriodNumber
    }
    
    return nil
}
```

#### 1.2 缓存查找逻辑

生成订单时优先从缓存中查找期号：

```go
// getPeriodNumberByTime 根据时间获取对应的期号（使用缓存）
func (s *FakeOrderService) getPeriodNumberByTime(targetTime time.Time) string {
    // 首先尝试从缓存中查找
    for key, periodNumber := range s.periodCache {
        parts := strings.Split(key, "_")
        if len(parts) == 2 {
            startTime, _ := time.Parse("2006-01-02 15:04:05", parts[0])
            endTime, _ := time.Parse("2006-01-02 15:04:05", parts[1])
            
            // 检查目标时间是否在这个范围内
            if targetTime.After(startTime) && targetTime.Before(endTime) {
                return periodNumber
            }
        }
    }
    
    // 如果缓存中没有找到，回退到数据库查询
    ctx := context.Background()
    periodRepo := database.NewLotteryPeriodRepository()
    
    period, err := periodRepo.GetPeriodByTime(ctx, targetTime)
    if err != nil {
        return targetTime.Format("20240101")
    }
    
    return period.PeriodNumber
}
```

### 2. 期数查询方法

#### 2.1 时间范围查询

新增了 `GetPeriodsByTimeRange` 方法，支持批量查询期数：

```go
// GetPeriodsByTimeRange 根据时间范围获取期数列表
func (r *LotteryPeriodRepository) GetPeriodsByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*models.LotteryPeriod, error) {
    var periods []*models.LotteryPeriod

    // 查询与指定时间范围有重叠的期数
    err := r.db.WithContext(ctx).
        Where("(order_start_time <= ? AND order_end_time > ?) OR (order_start_time < ? AND order_end_time >= ?) OR (order_start_time >= ? AND order_end_time <= ?)",
            endTime, startTime, endTime, startTime, startTime, endTime).
        Order("order_start_time ASC").
        Find(&periods).Error

    if err != nil {
        return nil, err
    }

    return periods, nil
}
```

#### 2.2 单时间点查询

保留了原有的 `GetPeriodByTime` 方法作为回退机制：

```go
// GetPeriodByTime 根据时间获取对应的期数
func (r *LotteryPeriodRepository) GetPeriodByTime(ctx context.Context, targetTime time.Time) (*models.LotteryPeriod, error) {
    var period models.LotteryPeriod

    // 查询在指定时间范围内的期数
    err := r.db.WithContext(ctx).
        Where("order_start_time <= ? AND order_end_time > ?", targetTime, targetTime).
        First(&period).Error

    if err != nil {
        // 如果没有找到对应时间的期数，返回最近的期数
        err = r.db.WithContext(ctx).
            Order("created_at DESC").
            First(&period).Error
        
        if err != nil {
            return nil, err
        }
    }

    return &period, nil
}
```

### 2. 查询条件

期号查询使用以下条件：
- `order_start_time <= 目标时间`
- `order_end_time > 目标时间`

这确保了订单被分配到正确的时间段内。

### 3. 容错机制

如果找不到对应时间的期数，系统会：
1. 首先尝试获取最近的期数
2. 如果数据库中没有期数，使用时间格式化作为期号（如：20240101）

## 实现细节

### 1. 假订单生成服务修改

在 `FakeOrderService` 中：

```go
// getPeriodNumberByTime 根据时间获取对应的期号
func (s *FakeOrderService) getPeriodNumberByTime(targetTime time.Time) string {
    ctx := context.Background()
    periodRepo := database.NewLotteryPeriodRepository()
    
    period, err := periodRepo.GetPeriodByTime(ctx, targetTime)
    if err != nil {
        // 如果获取失败，使用目标时间生成期号
        return targetTime.Format("20240101")
    }
    
    return period.PeriodNumber
}
```

### 2. 订单生成时使用

在生成购买单时：

```go
order := &models.Order{
    OrderNo:        utils.GenerateSystemOrderNo(),
    Uid:            utils.GenerateSystemUID(),
    PeriodNumber:   s.getPeriodNumberByTime(createdAt), // 使用生成时间对应的期号
    // ... 其他字段
}
```

## 期数表结构

期数表（lottery_periods）包含以下关键字段：

```sql
CREATE TABLE `lottery_periods` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `period_number` varchar(20) NOT NULL COMMENT '期数编号',
  `order_start_time` datetime NOT NULL COMMENT '订单开始时间',
  `order_end_time` datetime NOT NULL COMMENT '订单结束时间',
  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '期数状态',
  -- ... 其他字段
);
```

## 测试验证

### 1. 测试脚本

使用 `test_scripts/test_period_assignment.sh` 进行测试：

```bash
./test_period_assignment.sh
```

### 2. 测试内容

- 验证期数查询功能
- 检查订单期号分配
- 统计不同期号的订单数量
- 验证时间段匹配逻辑

### 3. 预期结果

- 订单应该根据创建时间分配到正确的期数
- 同一时间段的订单应该使用相同的期号
- 不同时间段的订单可能使用不同的期号

## 优势

### 1. 数据真实性

- 订单期号与真实时间段对应
- 符合业务逻辑要求
- 便于数据分析和统计

### 2. 性能优化

- **减少数据库查询**：从 N 次减少到 1 次
- **提高生成速度**：内存查找比数据库查询快
- **降低数据库压力**：减少并发查询
- **保持数据准确性**：缓存失效时自动回退

### 3. 灵活性

- 支持多个期数同时存在
- 自动处理期数切换
- 容错机制确保系统稳定

### 4. 可维护性

- 逻辑清晰，易于理解
- 代码复用性好
- 便于后续扩展

### 5. 缓存策略

- **缓存时间范围**：当前时间前后30分钟
- **缓存更新**：每次生成任务开始时刷新
- **容错机制**：缓存未命中时回退到数据库查询

## 注意事项

### 1. 期数管理

- 确保期数的时间段不重叠
- 定期清理过期的期数数据
- 监控期数状态变化

### 2. 性能考虑

- 期数查询会涉及数据库操作
- 考虑添加缓存机制
- 批量生成时优化查询频率

### 3. 数据一致性

- 确保期数数据的准确性
- 定期验证期号分配的正确性
- 监控异常情况

## 总结

通过实现基于时间的期号分配逻辑，假订单生成系统现在能够：

1. **准确分配期号**：根据订单创建时间查询对应的期数
2. **支持多期数**：同时处理多个活跃期数
3. **提供容错机制**：确保系统在各种情况下都能正常工作
4. **保持数据一致性**：订单期号与业务逻辑保持一致

这个改进使得假订单数据更加真实和有用，为业务分析和系统测试提供了更好的支持。 