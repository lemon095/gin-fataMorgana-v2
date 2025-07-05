# 假订单生成系统实现文档

## 概述

本系统实现了自动生成假订单（水单）的功能，用于模拟真实的订单数据，提升用户体验。系统包含订单生成、数据清理、定时任务等核心功能。

## 系统架构

### 核心组件

1. **系统UID生成器** (`utils/system_uid_generator.go`)
   - 生成7位系统UID（区别于正常用户的8位UID）
   - 生成系统订单号和拼单号

2. **假订单生成服务** (`services/fake_order_service.go`)
   - 生成购买单和拼单
   - 计算真实订单金额和利润
   - 随机分配订单状态

3. **数据清理服务** (`services/data_cleanup_service.go`)
   - 定期清理旧的系统订单数据
   - 分批删除，避免锁表

4. **定时任务服务** (`services/cron_service.go`)
   - 管理订单生成和数据清理的定时任务
   - 支持手动触发和状态监控

## 功能特性

### 1. 系统UID生成

#### 生成规则
- **正常用户**: 8位数字（现有规则）
- **系统用户**: 7位数字（新规则）
- **组成**: 时间戳(3位) + 机器ID(2位) + 序列号(2位)

#### 订单号格式
- **正常订单**: `ORD` + 8位数字
- **系统订单**: `ORD` + 7位数字
- **系统拼单**: `GB` + 7位数字

### 2. 订单生成策略

#### 购买单生成
- **任务数量**: 100-2000随机
- **金额计算**: 基于真实价格配置
- **利润计算**: 基于随机用户等级
- **状态分布**:
  - 60% 进行中 (pending)
  - 30% 已完成 (success)
  - 10% 已关闭 (cancelled)

#### 拼单生成
- **人均金额**: 10.00-60.00随机
- **状态分布**:
  - 20% 待开始 (not_started)
  - 50% 进行中 (pending)
  - 30% 已完成 (success)

### 3. 时间窗口设计

#### 生成时间范围
- **时间窗口**: 当前时间前后各5分钟（总共10分钟）
- **目的**: 确保数据连续性，避免时间断层
- **实现**: 在时间窗口内随机分布订单创建时间

### 4. 数据清理策略

#### 清理规则
- **清理频率**: 每天凌晨2点
- **保留天数**: 2天（可配置）
- **清理范围**: 系统订单和拼单
- **清理方式**: 分批删除，避免锁表

## 配置说明

### 配置文件 (`config.yaml`)

```yaml
# 假订单生成配置
fake_data:
  enabled: true                    # 是否启用
  cron_expression: "*/5 * * * *"   # 订单生成定时表达式（每5分钟）
  cleanup_cron: "0 2 * * *"        # 数据清理定时表达式（每天凌晨2点）
  min_orders: 80                   # 最小订单数量
  max_orders: 100                  # 最大订单数量
  purchase_ratio: 0.7              # 购买单比例（70%）
  task_min_count: 100              # 任务数最小值
  task_max_count: 2000             # 任务数最大值
  retention_days: 2                # 数据保留天数
```

### 配置结构体 (`config/config.go`)

```go
type FakeDataConfig struct {
    Enabled         bool    `mapstructure:"enabled"`
    CronExpression  string  `mapstructure:"cron_expression"`
    CleanupCron     string  `mapstructure:"cleanup_cron"`
    MinOrders       int     `mapstructure:"min_orders"`
    MaxOrders       int     `mapstructure:"max_orders"`
    PurchaseRatio   float64 `mapstructure:"purchase_ratio"`
    TaskMinCount    int     `mapstructure:"task_min_count"`
    TaskMaxCount    int     `mapstructure:"task_max_count"`
    RetentionDays   int     `mapstructure:"retention_days"`
}
```

## 数据库设计

### 订单表新增字段

```sql
ALTER TABLE orders ADD COLUMN is_system_order tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否系统订单 0-否 1-是';
ALTER TABLE orders ADD INDEX idx_is_system_order (is_system_order);
```

### 字段说明

- `is_system_order`: 标识订单是否为系统生成
- `uid`: 7位系统UID（系统订单）或8位用户UID（正常订单）
- `order_no`: 系统订单号（7位）或正常订单号（8位）

## 使用方式

### 1. 启动定时任务

```go
// main.go
func main() {
    // 初始化系统UID生成器
    utils.InitSystemUIDGenerator(config.GlobalConfig.Snowflake.WorkerID)
    
    // 创建定时任务服务
    cronConfig := &services.CronConfig{
        Enabled:           config.GlobalConfig.FakeData.Enabled,
        OrderCronExpr:     config.GlobalConfig.FakeData.CronExpression,
        CleanupCronExpr:   config.GlobalConfig.FakeData.CleanupCron,
        MinOrders:         config.GlobalConfig.FakeData.MinOrders,
        MaxOrders:         config.GlobalConfig.FakeData.MaxOrders,
        PurchaseRatio:     config.GlobalConfig.FakeData.PurchaseRatio,
        TaskMinCount:      config.GlobalConfig.FakeData.TaskMinCount,
        TaskMaxCount:      config.GlobalConfig.FakeData.TaskMaxCount,
        RetentionDays:     config.GlobalConfig.FakeData.RetentionDays,
    }
    
    cronService := services.NewCronService(cronConfig)
    
    // 启动定时任务
    if err := cronService.Start(); err != nil {
        log.Fatalf("启动定时任务失败: %v", err)
    }
    
    // 优雅关闭
    defer cronService.Stop()
}
```

### 2. 手动生成订单

```go
// 手动生成指定数量的假订单
stats, err := cronService.ManualGenerateOrders(50)
if err != nil {
    log.Printf("手动生成订单失败: %v", err)
} else {
    log.Printf("生成完成: 总数=%d, 购买单=%d, 拼单=%d", 
        stats.TotalGenerated, stats.PurchaseOrders, stats.GroupBuyOrders)
}
```

### 3. 手动清理数据

```go
// 手动清理旧数据
stats, err := cronService.ManualCleanup()
if err != nil {
    log.Printf("手动清理失败: %v", err)
} else {
    log.Printf("清理完成: 删除订单=%d, 删除拼单=%d", 
        stats.DeletedOrders, stats.DeletedGroupBuys)
}
```

## 测试验证

### 测试脚本

使用 `test_scripts/test_fake_order_generation.sh` 进行功能测试：

```bash
# 运行测试脚本
./test_scripts/test_fake_order_generation.sh
```

### 测试内容

1. **系统UID生成验证**
   - 检查7位UID生成
   - 验证订单号格式

2. **订单数据验证**
   - 检查系统订单标识
   - 验证任务数量范围（100-2000）
   - 验证时间分布

3. **状态分布验证**
   - 检查订单状态分布
   - 验证拼单状态分布

4. **数据库查询验证**
   - 检查表结构
   - 验证数据完整性

## 监控和日志

### 日志输出

系统会输出详细的生成和清理日志：

```
[INFO] 开始生成 85 条假订单
[INFO] 成功插入 60 条购买单
[INFO] 成功插入 25 条拼单
[INFO] 假订单生成完成: 总数=85, 购买单=60, 拼单=25, 总金额=1234.56, 总利润=67.89, 耗时=1.2s
```

### 状态监控

可以通过以下方式监控系统状态：

```go
// 获取定时任务状态
status := cronService.GetCronStatus()
for name, info := range status {
    log.Printf("任务: %s, 下次执行: %v", name, info["next_run"])
}
```

## 性能优化

### 1. 批量操作

- 使用分批删除避免锁表
- 批量插入提高性能

### 2. 数据库优化

- 添加 `is_system_order` 索引
- 定期清理旧数据减少表大小

### 3. 内存管理

- 使用对象池减少GC压力
- 合理设置批次大小

## 注意事项

### 1. 数据一致性

- 系统订单不影响真实用户数据
- 通过 `is_system_order` 字段区分

### 2. 性能影响

- 定时任务在低峰期执行
- 数据清理分批进行

### 3. 配置管理

- 支持环境变量覆盖
- 配置文件热更新

### 4. 错误处理

- 定时任务异常恢复
- 数据库连接重试

## 扩展功能

### 1. 可扩展的生成策略

- 支持自定义订单类型
- 可配置的状态分布
- 灵活的时间窗口

### 2. 监控告警

- 生成失败告警
- 数据量异常告警
- 性能监控

### 3. 管理界面

- 手动生成订单
- 查看生成统计
- 配置管理

## 总结

假订单生成系统提供了完整的模拟数据生成方案，通过合理的配置和优化，可以有效地提升用户体验，同时保证系统的稳定性和性能。系统设计考虑了数据一致性、性能优化和可扩展性，为后续功能扩展提供了良好的基础。 