# 钱包利润类型说明

## 概述

在原有的钱包流水类型基础上，新增了"利润"类型，用于记录用户获得的利润收入。

## 利润类型详情

### 基本信息
- **类型值**: `"profit"`
- **中文名称**: "利润"
- **说明**: 用户获得利润收入
- **金额显示**: 正数 (+)
- **默认状态**: 成功 (success)

### 特点
1. **直接成功**: 利润交易创建后直接为成功状态，无需审核
2. **增加余额**: 利润会直接增加用户钱包余额
3. **关联订单**: 可以关联到具体的订单号
4. **正数显示**: 在流水记录中显示为正数，表示收入

## 相关代码修改

### 1. 模型定义 (`models/wallet_transaction.go`)

```go
// 新增利润类型常量
const (
    TransactionTypeProfit = "profit"    // 利润
)

// 更新类型名称映射
typeNames := map[string]string{
    TransactionTypeProfit: "利润",
}

// 更新金额显示逻辑
case TransactionTypeRecharge, TransactionTypeProfit:
    return "+" + formatAmount(t.Amount)
```

### 2. 服务层 (`services/wallet_service.go`)

新增了两个方法：

#### AddProfit - 添加利润
```go
func (s *WalletService) AddProfit(ctx context.Context, uid string, amount float64, description string) error
```
- 验证利润金额必须大于0
- 检查钱包状态（不能是冻结状态）
- 使用原子操作增加余额

#### CreateProfitTransaction - 创建利润交易记录
```go
func (s *WalletService) CreateProfitTransaction(ctx context.Context, uid string, amount float64, description string, relatedOrderNo string) (string, error)
```
- 生成利润交易流水号
- 创建交易记录（状态为成功）
- 更新钱包余额
- 返回交易流水号

### 3. 控制器层 (`controllers/wallet_controller.go`)

新增接口：

#### AddProfit - 添加利润接口
```go
func (wc *WalletController) AddProfit(c *gin.Context)
```

**请求参数**:
```json
{
    "uid": "用户ID",
    "amount": 100.00,
    "description": "订单完成获得利润",
    "order_no": "关联订单号（可选）"
}
```

**响应**:
```json
{
    "code": 200,
    "message": "利润添加成功",
    "data": {
        "transaction_no": "PROFIT_20241201_001"
    }
}
```

## 使用场景

### 1. 订单完成获得利润
当用户完成订单任务后，系统自动为用户添加利润：

```go
// 示例：订单完成后添加利润
transactionNo, err := walletService.CreateProfitTransaction(
    ctx, 
    userUid, 
    order.ProfitAmount, 
    "订单完成获得利润", 
    order.OrderNo,
)
```

### 2. 活动奖励
系统活动或奖励发放：

```go
// 示例：活动奖励
transactionNo, err := walletService.CreateProfitTransaction(
    ctx, 
    userUid, 
    rewardAmount, 
    "活动奖励", 
    "",
)
```

### 3. 管理员手动添加
管理员为用户手动添加利润：

```go
// 示例：管理员手动添加
transactionNo, err := walletService.CreateProfitTransaction(
    ctx, 
    userUid, 
    amount, 
    "管理员手动添加利润", 
    "",
)
```

## 数据库记录示例

```sql
INSERT INTO wallet_transactions (
    transaction_no, uid, type, amount, 
    balance_before, balance_after, status, 
    description, related_order_no, created_at
) VALUES (
    'PROFIT_20241201_001', 'USER001', 'profit', 50.00,
    100.00, 150.00, 'success',
    '订单完成获得利润', 'ORDER_20241201_001', NOW()
);
```

## 流水记录显示

在用户的钱包流水记录中，利润类型会显示为：

```
+50.00  利润  订单完成获得利润  2024-12-01 10:30:00
```

## 注意事项

1. **权限控制**: 只有用户本人或管理员可以添加利润
2. **金额验证**: 利润金额必须大于0
3. **钱包状态**: 冻结状态的钱包无法添加利润
4. **并发安全**: 使用分布式锁确保并发安全
5. **缓存更新**: 利润添加后会自动更新缓存

## 错误码

- `CodeInvalidParams`: 利润金额必须大于0
- `CodeWalletFrozenRecharge`: 钱包已被冻结，无法添加利润
- `CodeWalletGetFailed`: 获取钱包信息失败
- `CodeTransactionCreateFailed`: 创建利润交易记录失败

## 扩展建议

1. **利润统计**: 可以添加利润统计功能，统计用户总利润
2. **利润分类**: 可以细分利润类型（如订单利润、活动利润、推荐利润等）
3. **利润规则**: 可以添加利润计算规则和限制
4. **利润报表**: 可以生成利润相关的报表和图表 