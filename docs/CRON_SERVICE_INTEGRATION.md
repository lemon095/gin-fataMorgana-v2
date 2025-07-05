# 定时任务服务集成文档

## 概述

定时任务服务已经成功集成到主应用程序中，用于自动生成假订单（水单）和定期清理旧数据。

## 调用位置

### 1. 启动位置

定时任务服务在 `main.go` 中的启动顺序如下：

```go
// 1. 初始化数据库连接
if err := database.InitMySQL(); err != nil {
    log.Printf("初始化MySQL失败: %v", err)
    os.Exit(1)
}

// 2. 自动迁移数据库表
if err := database.AutoMigrate(); err != nil {
    log.Printf("数据库迁移失败: %v", err)
    os.Exit(1)
}

// 3. 初始化Redis
if err := database.InitRedis(); err != nil {
    log.Printf("初始化Redis失败: %v", err)
    os.Exit(1)
}

// 4. 初始化定时任务控制器
cronController := controllers.NewCronController()

// 5. 启动定时任务服务
var cronService *services.CronService
if config.GlobalConfig.FakeData.Enabled {
    log.Println("启动定时任务服务...")
    
    // 创建定时任务配置
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
    
    // 创建并启动定时任务服务
    cronService = services.NewCronService(cronConfig)
    if err := cronService.Start(); err != nil {
        log.Printf("启动定时任务失败: %v", err)
    } else {
        log.Println("定时任务服务启动成功")
    }
    
    // 注入定时任务服务到控制器
    cronController.SetCronService(cronService)
    
    // 优雅关闭时停止定时任务
    defer func() {
        if cronService != nil {
            cronService.Stop()
            log.Println("定时任务服务已停止")
        }
    }()
} else {
    log.Println("定时任务服务已禁用")
}
```

### 2. 启动顺序的重要性

定时任务服务必须在以下组件初始化之后启动：

1. **数据库连接** - 定时任务需要访问数据库表
2. **数据库迁移** - 确保所需的表结构已创建
3. **Redis连接** - 定时任务可能使用Redis缓存
4. **系统UID生成器** - 生成系统订单需要UID生成器

这样可以确保定时任务在启动时能够正常访问所有必要的资源。

## API接口

### 1. 手动生成订单

**接口**: `POST /api/v1/cron/manual-generate`

**认证**: 需要Bearer Token

**请求参数**:
```json
{
    "count": 10
}
```

**参数说明**:
- `count`: 要生成的订单数量，范围1-1000

**响应示例**:
```json
{
    "code": 0,
    "message": "手动生成订单成功",
    "data": {
        "total_generated": 10,
        "purchase_orders": 7,
        "group_buy_orders": 3,
        "total_amount": 1234.56,
        "total_profit": 67.89,
        "last_generation": "2024-01-01T12:00:00Z",
        "average_time": "1.2s"
    },
    "timestamp": 1751365370
}
```

### 2. 手动清理数据

**接口**: `POST /api/v1/cron/manual-cleanup`

**认证**: 需要Bearer Token

**请求参数**: 无

**响应示例**:
```json
{
    "code": 0,
    "message": "手动清理数据成功",
    "data": {
        "deleted_orders": 150,
        "deleted_group_buys": 25,
        "last_cleanup": "2024-01-01T12:00:00Z",
        "cleanup_time": "2.5s"
    },
    "timestamp": 1751365370
}
```

### 3. 获取定时任务状态

**接口**: `GET /api/v1/cron/status`

**认证**: 需要Bearer Token

**请求参数**: 无

**响应示例**:
```json
{
    "code": 0,
    "message": "操作成功",
    "data": {
        "cron_status": {
            "task_0": {
                "next_run": "2024-01-01T12:05:00Z",
                "prev_run": "2024-01-01T12:00:00Z",
                "schedule": "*/5 * * * *"
            },
            "task_1": {
                "next_run": "2024-01-02T02:00:00Z",
                "prev_run": "2024-01-01T02:00:00Z",
                "schedule": "0 2 * * *"
            }
        }
    },
    "timestamp": 1751365370
}
```

## 配置说明

### 配置文件 (config.yaml)

```yaml
# 假订单生成配置
fake_data:
  enabled: true                    # 是否启用定时任务
  cron_expression: "*/5 * * * *"   # 订单生成定时表达式（每5分钟）
  cleanup_cron: "0 2 * * *"        # 数据清理定时表达式（每天凌晨2点）
  min_orders: 80                   # 最小订单数量
  max_orders: 100                  # 最大订单数量
  purchase_ratio: 0.7              # 购买单比例（70%）
  task_min_count: 100              # 任务数最小值
  task_max_count: 2000             # 任务数最大值
  retention_days: 2                # 数据保留天数
```

### 配置参数说明

| 参数名 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| `enabled` | bool | 是否启用定时任务 | true |
| `cron_expression` | string | 订单生成定时表达式 | "*/5 * * * *" |
| `cleanup_cron` | string | 数据清理定时表达式 | "0 2 * * *" |
| `min_orders` | int | 每次生成的最小订单数 | 80 |
| `max_orders` | int | 每次生成的最大订单数 | 100 |
| `purchase_ratio` | float64 | 购买单比例（0-1） | 0.7 |
| `task_min_count` | int | 任务数最小值 | 100 |
| `task_max_count` | int | 任务数最大值 | 2000 |
| `retention_days` | int | 数据保留天数 | 2 |

## 定时任务说明

### 1. 订单生成任务

- **执行频率**: 每5分钟（可配置）
- **生成数量**: 80-100条随机订单（可配置）
- **订单类型**: 购买单（70%）和拼单（30%）
- **时间窗口**: 当前时间前后各5分钟
- **状态分布**: 合理的订单状态分布

### 2. 数据清理任务

- **执行频率**: 每天凌晨2点（可配置）
- **清理范围**: 系统生成的订单和拼单
- **保留天数**: 2天（可配置）
- **清理方式**: 分批删除，避免锁表

## 测试方法

### 1. 使用测试脚本

```bash
# 运行定时任务测试脚本
./test_scripts/test_cron_service.sh
```

### 2. 手动测试API

```bash
# 1. 登录获取token
curl -X POST http://localhost:9001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"account": "test@example.com", "password": "password"}'

# 2. 手动生成订单
curl -X POST http://localhost:9001/api/v1/cron/manual-generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"count": 10}'

# 3. 获取定时任务状态
curl -X GET http://localhost:9001/api/v1/cron/status \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 监控和日志

### 1. 启动日志

```
[INFO] 启动定时任务服务...
[INFO] 假订单生成定时任务已启动，表达式: */5 * * * *
[INFO] 数据清理定时任务已启动，表达式: 0 2 * * *
[INFO] 定时任务服务启动成功
```

### 2. 执行日志

```
[INFO] 开始执行假订单生成定时任务...
[INFO] 开始生成 85 条假订单
[INFO] 成功插入 60 条购买单
[INFO] 成功插入 25 条拼单
[INFO] 假订单生成定时任务完成: 总数=85, 购买单=60, 拼单=25, 总金额=1234.56, 总利润=67.89, 耗时=1.2s
```

### 3. 清理日志

```
[INFO] 开始执行数据清理定时任务...
[INFO] 开始清理 2 天前的系统订单数据，清理时间点: 2024-01-01 12:00:00
[INFO] 删除系统订单批次: 1000 条
[INFO] 删除系统拼单批次: 250 条
[INFO] 数据清理定时任务完成: 删除订单=1500, 删除拼单=250, 耗时=2.5s
```

## 注意事项

### 1. 性能考虑

- 定时任务在后台异步执行，不会阻塞主服务
- 数据清理采用分批删除，避免锁表
- 订单生成使用缓存优化期号分配

### 2. 安全考虑

- API接口需要认证，防止未授权访问
- 可以添加管理员权限验证
- 手动清理功能需要谨慎使用

### 3. 配置管理

- 支持环境变量覆盖配置
- 配置文件热更新（需要重启服务）
- 可以通过API动态调整部分参数

### 4. 错误处理

- 定时任务异常不会影响主服务
- 自动重试机制
- 详细的错误日志记录

## 总结

定时任务服务已经成功集成到主应用程序中，提供了完整的假订单生成和数据清理功能。通过合理的启动顺序、完善的API接口和详细的监控日志，确保了系统的稳定性和可维护性。 