# 钱包并发安全操作指南

## 🚨 并发问题场景

在钱包系统中，如果同时有多个进程对同一用户的余额进行操作，可能会导致数据不一致：

### 问题示例
```go
// 时序问题：
// T1: 进程A查询余额 -> 1000元
// T2: 进程B查询余额 -> 1000元  
// T3: 进程A扣钱500元 -> 数据库更新为500元
// T4: 进程B加钱300元 -> 数据库更新为1300元（错误！应该是800元）
// 结果：最终余额是1300元，而不是正确的800元
```

## 🛠️ 解决方案

### 1. 用户级别互斥锁
- 为每个用户创建独立的互斥锁
- 确保同一用户的余额操作串行化执行
- 不同用户的操作可以并行执行

### 2. 原子性操作
- 所有余额操作都在锁保护下执行
- 操作前重新获取最新数据
- 操作后立即更新缓存

### 3. 死锁预防
- 转账操作按UID排序加锁
- 避免循环等待

## 📖 使用示例

### 基本余额操作

```go
// 初始化服务
walletService := NewWalletConcurrentService()

// 扣减余额
err := walletService.WithdrawBalance(ctx, "user123", 100.0, "购买商品")
if err != nil {
    // 处理错误
}

// 增加余额
err = walletService.AddBalance(ctx, "user123", 50.0, "充值")
if err != nil {
    // 处理错误
}

// 查询余额
wallet, err := walletService.GetWalletBalance(ctx, "user123")
if err != nil {
    // 处理错误
}
fmt.Printf("当前余额: %.2f\n", wallet.Balance)
```

### 转账操作

```go
// 用户A向用户B转账
err := walletService.TransferBalance(ctx, "userA", "userB", 200.0, "转账")
if err != nil {
    // 处理错误
}
```

### 批量操作

```go
// 批量余额操作
operations := []BalanceOperation{
    {UID: "user1", Type: "withdraw", Amount: 100.0, Description: "扣款"},
    {UID: "user1", Type: "add", Amount: 50.0, Description: "退款"},
    {UID: "user2", Type: "add", Amount: 200.0, Description: "充值"},
}

err := walletService.BatchBalanceOperation(ctx, operations)
if err != nil {
    // 处理错误
}
```

### 余额检查

```go
// 检查余额是否足够
sufficient, err := walletService.CheckBalanceSufficient(ctx, "user123", 500.0)
if err != nil {
    // 处理错误
}

if sufficient {
    // 执行扣款操作
    err = walletService.WithdrawBalance(ctx, "user123", 500.0, "大额购买")
} else {
    // 余额不足
}
```

## 🔒 并发安全特性

### 1. 用户级别隔离
- 不同用户的余额操作完全独立
- 不会相互影响

### 2. 原子性保证
- 每个操作要么完全成功，要么完全失败
- 不会出现部分更新的情况

### 3. 数据一致性
- 操作前重新获取最新数据
- 操作后立即更新缓存
- 确保数据库和缓存的一致性

### 4. 死锁预防
- 转账操作按UID排序加锁
- 避免循环等待导致的死锁

## ⚡ 性能优化

### 1. 缓存优先
- 查询操作优先从缓存获取
- 减少数据库访问

### 2. 批量操作
- 支持批量余额操作
- 减少锁的获取和释放次数

### 3. 异步缓存更新
- 缓存更新失败不影响主流程
- 只记录警告日志

## 🚨 注意事项

### 1. 锁的粒度
- 使用用户级别的锁，不是全局锁
- 不同用户的操作可以并行执行

### 2. 错误处理
- 所有操作都可能返回错误
- 需要正确处理各种错误情况

### 3. 事务边界
- 每个操作都是独立的事务
- 如果需要跨操作的事务，需要额外的处理

### 4. 内存管理
- 互斥锁会占用内存
- 需要定期清理长时间未使用的锁

## 🔧 集成到现有系统

### 1. 替换现有钱包服务
```go
// 在控制器中使用
func (c *WalletController) Withdraw(ctx *gin.Context) {
    walletService := NewWalletConcurrentService()
    
    uid := ctx.GetString("uid")
    amount := ctx.GetFloat64("amount")
    
    err := walletService.WithdrawBalance(ctx, uid, amount, "提现")
    if err != nil {
        utils.ErrorWithMessage(ctx, err.Code, err.Message)
        return
    }
    
    utils.Success(ctx, "提现成功")
}
```

### 2. 在订单创建中使用
```go
// 创建订单时扣减余额
func (s *OrderService) CreateOrder(ctx context.Context, uid string, amount float64) error {
    walletService := NewWalletConcurrentService()
    
    // 原子性扣减余额
    err := walletService.WithdrawBalance(ctx, uid, amount, "创建订单")
    if err != nil {
        return err
    }
    
    // 创建订单逻辑...
    return nil
}
```

## 📊 监控和日志

### 1. 操作日志
- 记录每次余额操作的前后余额
- 便于审计和问题排查

### 2. 性能监控
- 监控锁的等待时间
- 监控并发操作的性能

### 3. 错误监控
- 监控余额不足等业务错误
- 监控系统错误

## 🎯 最佳实践

### 1. 合理使用批量操作
- 对于多个相关操作，使用批量操作
- 减少锁的获取和释放次数

### 2. 及时处理错误
- 余额不足等错误需要及时处理
- 避免重试导致的问题

### 3. 定期清理
- 定期清理长时间未使用的互斥锁
- 防止内存泄漏

### 4. 监控告警
- 设置合理的监控告警
- 及时发现和处理问题 