# 订单价格计算逻辑修改

## 修改概述

订单创建接口的价格计算逻辑进行了重大调整，`amount` 字段从总价改为单价。

## 修改内容

### 1. 价格计算逻辑

**修改前**：
- `amount` 表示总价
- 需要根据价格配置校验金额匹配

**修改后**：
- `amount` 表示单价
- 总价 = 单价 × 有值的类型数量
- 有值的类型数量 = 1，没有值的类型数量 = 0

### 2. 类型数量计算

```go
// 计算有值的类型数量
typeCount := 0
if req.LikeCount > 0 {
    typeCount++
}
if req.ShareCount > 0 {
    typeCount++
}
if req.FollowCount > 0 {
    typeCount++
}
if req.FavoriteCount > 0 {
    typeCount++
}

// 计算总价 = 单价 × 类型数量
totalAmount := req.Amount * float64(typeCount)
```

### 3. 验证逻辑简化

**修改前**：
- 根据价格配置计算总价
- 校验请求金额与计算金额是否匹配

**修改后**：
- 只验证至少选择一种任务类型
- 移除具体金额校验

```go
// 验证订单金额（已简化，不再校验具体金额）
func (s *OrderService) validateOrderAmount(ctx context.Context, req *CreateOrderRequest) error {
    // 新的逻辑：amount作为单价，不再需要校验具体金额
    // 只需要确保至少选择了一种类型
    if req.LikeCount == 0 && req.ShareCount == 0 && req.FollowCount == 0 && req.FavoriteCount == 0 {
        return utils.NewAppError(utils.CodeOrderAmountMismatch, "请至少选择一种任务类型")
    }
    
    return nil
}
```

## 示例场景

### 场景1：只选择点赞
```json
{
  "amount": 10.00,        // 单价10元
  "like_count": 1,        // 有值，类型数量=1
  "share_count": 0,       // 无值，类型数量=0
  "follow_count": 0,      // 无值，类型数量=0
  "favorite_count": 0     // 无值，类型数量=0
}
```
**总价计算**：10.00 × 1 = 10.00元

### 场景2：选择点赞和分享
```json
{
  "amount": 15.00,        // 单价15元
  "like_count": 1,        // 有值，类型数量=1
  "share_count": 1,       // 有值，类型数量=1
  "follow_count": 0,      // 无值，类型数量=0
  "favorite_count": 0     // 无值，类型数量=0
}
```
**总价计算**：15.00 × 2 = 30.00元

### 场景3：选择所有类型
```json
{
  "amount": 20.00,        // 单价20元
  "like_count": 1,        // 有值，类型数量=1
  "share_count": 1,       // 有值，类型数量=1
  "follow_count": 1,      // 有值，类型数量=1
  "favorite_count": 1     // 有值，类型数量=1
}
```
**总价计算**：20.00 × 4 = 80.00元

## 影响范围

### 修改的文件
- `services/order_service.go` - 主要逻辑修改
- `test_scripts/test_order_price_logic.sh` - 新增测试脚本

### 影响的功能
- 订单创建接口
- 钱包余额扣减
- 交易记录创建
- 利润计算

## 测试验证

运行测试脚本验证修改效果：
```bash
./test_scripts/test_order_price_logic.sh
```

测试内容包括：
- 单类型订单价格计算
- 多类型订单价格计算
- 全类型订单价格计算
- 无类型选择错误处理 