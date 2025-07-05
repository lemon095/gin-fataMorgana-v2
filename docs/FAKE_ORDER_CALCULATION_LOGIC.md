# 假订单计算逻辑详细文档

## 概述

本文档详细说明假订单生成系统中的计算逻辑，包括拼单水单计算、购买单状态处理、任务状态管理等核心逻辑。

## 拼单水单计算逻辑

### 1. 拼单人均金额计算

#### 计算原理
拼单的人均金额不是随机生成的，而是基于当前配置的价格来计算，确保数据的真实性和合理性。

#### 计算步骤

1. **获取价格配置**
   ```go
   purchaseConfig := s.getPurchaseConfig()
   // 包含：LikeAmount, ShareAmount, ForwardAmount, FavoriteAmount
   ```

2. **选择任务类型组合**
   - 随机选择2-4种任务类型（点赞、转发、关注、收藏）
   - 确保任务类型不重复

3. **计算每种任务类型的最大数量**
   - 最大任务数量：3-7个（每种任务类型）
   - 每种任务类型的实际数量：1到最大任务数量之间

4. **计算人均金额**
   ```go
   var totalAmount float64
   for _, taskType := range selectedTasks {
       taskCount := rand.Intn(baseTaskCount) + 1
       switch taskType {
       case "like":
           totalAmount += float64(taskCount) * purchaseConfig.LikeAmount
       case "share":
           totalAmount += float64(taskCount) * purchaseConfig.ShareAmount
       case "follow":
           totalAmount += float64(taskCount) * purchaseConfig.ForwardAmount
       case "favorite":
           totalAmount += float64(taskCount) * purchaseConfig.FavoriteAmount
       }
   }
   ```

5. **金额范围限制**
   - 最小金额：5.00元
   - 最大金额：50.00元
   - 确保金额在合理范围内

#### 示例计算

假设当前价格配置为：
- 点赞：0.1元/次
- 转发：0.2元/次
- 关注：0.3元/次
- 收藏：0.4元/次

随机选择任务类型：点赞、转发、关注
每种任务类型的最大数量：5个

计算过程：
- 点赞：3次 × 0.1元 = 0.3元
- 转发：2次 × 0.2元 = 0.4元
- 关注：4次 × 0.3元 = 1.2元

人均金额：0.3 + 0.4 + 1.2 = 1.9元
（总计：9次任务，分布在3种任务类型中）

### 2. 拼单人数处理

#### 参与人数和目标人数

1. **当前参与人数** (`current_participants`)
   - 范围：1-3人
   - 表示当前已参与拼单的人数

2. **目标参与人数** (`target_participants`)
   - 范围：3-7人
   - 表示拼单成功需要的人数

3. **已付款金额** (`paid_amount`)
   ```go
   paidAmount := perPersonAmount * float64(currentParticipants)
   ```

4. **总金额** (`total_amount`)
   ```go
   totalAmount := perPersonAmount * float64(targetParticipants)
   ```

#### 示例数据

```go
groupBuy := &models.GroupBuy{
    CurrentParticipants: 2,  // 当前2人参与
    TargetParticipants:  5,  // 目标5人
    PerPersonAmount:     15.50, // 人均15.50元
    TotalAmount:         77.50,  // 总金额77.50元
    PaidAmount:          31.00,  // 已付款31.00元
}
```

## 购买单状态处理逻辑

### 1. 订单状态与任务状态的关系

#### 状态对应关系

| 订单状态 | 任务状态处理逻辑 |
|---------|----------------|
| `pending` (进行中) | 根据概率设置：30%已完成，70%待完成 |
| `success` (已完成) | 所有任务状态都设置为已完成 |
| `cancelled` (已关闭) | 所有任务状态都设置为已关闭 |

#### 实现逻辑

```go
func (s *FakeOrderService) getTaskStatus(count int, orderStatus string) string {
    if count == 0 {
        return models.TaskStatusSuccess // 任务数为0时直接完成
    }
    
    // 如果订单状态是已完成，任务状态也应该是已完成
    if orderStatus == models.OrderStatusSuccess {
        return models.TaskStatusSuccess
    }
    
    // 如果订单状态是已关闭，任务状态也应该是已关闭
    if orderStatus == models.OrderStatusCancelled {
        return models.TaskStatusCancelled
    }
    
    // 如果订单状态是进行中，根据概率设置任务状态
    randNum := rand.Float64()
    if randNum < 0.3 {
        return models.TaskStatusSuccess // 30% 已完成
    } else {
        return models.TaskStatusPending // 70% 待完成
    }
}
```

### 2. 任务状态设置

#### 购买单生成时的任务状态设置

```go
order := &models.Order{
    // ... 其他字段 ...
    LikeStatus:     s.getTaskStatus(likeCount, status),
    ShareStatus:    s.getTaskStatus(shareCount, status),
    FollowStatus:   s.getTaskStatus(followCount, status),
    FavoriteStatus: s.getTaskStatus(favoriteCount, status),
}
```

#### 状态分布示例

**订单状态：进行中 (pending)**
```json
{
    "status": "pending",
    "like_status": "success",      // 30%概率
    "share_status": "pending",     // 70%概率
    "follow_status": "success",    // 30%概率
    "favorite_status": "pending"   // 70%概率
}
```

**订单状态：已完成 (success)**
```json
{
    "status": "success",
    "like_status": "success",      // 100%概率
    "share_status": "success",     // 100%概率
    "follow_status": "success",    // 100%概率
    "favorite_status": "success"   // 100%概率
}
```

**订单状态：已关闭 (cancelled)**
```json
{
    "status": "cancelled",
    "like_status": "cancelled",      // 100%概率
    "share_status": "cancelled",     // 100%概率
    "follow_status": "cancelled",    // 100%概率
    "favorite_status": "cancelled"   // 100%概率
}
```

## 数据一致性保证

### 1. 订单状态与任务状态的一致性

- **已完成订单**：所有任务必须显示为已完成
- **已关闭订单**：所有任务显示为已关闭（表示任务被取消）
- **进行中订单**：任务状态根据概率分布，保持真实感

### 2. 金额计算的一致性

- **购买单金额**：基于真实价格配置计算
- **拼单金额**：基于真实价格配置计算
- **利润计算**：基于真实用户等级计算

### 3. 时间分布的一致性

- **创建时间**：在10分钟时间窗口内随机分布
- **过期时间**：根据订单状态合理设置
- **截止时间**：根据拼单状态合理设置

## 配置参数说明

### 价格配置

```yaml
# 从Redis获取的价格配置
purchase_config:
  like_amount: 0.1      # 点赞价格
  share_amount: 0.2     # 转发价格
  forward_amount: 0.3   # 关注价格
  favorite_amount: 0.4  # 收藏价格
```

### 生成参数

```yaml
fake_data:
  task_min_count: 100   # 任务数最小值
  task_max_count: 2000  # 任务数最大值
  purchase_ratio: 0.7   # 购买单比例
```

## 测试验证

### 1. 拼单金额验证

```sql
-- 验证拼单金额是否在合理范围内
SELECT 
    group_buy_no,
    per_person_amount,
    total_amount,
    target_participants,
    CASE 
        WHEN per_person_amount >= 5.0 AND per_person_amount <= 50.0 
        THEN 'OK' 
        ELSE 'INVALID' 
    END as amount_check
FROM group_buys 
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
LIMIT 10;
```

### 2. 任务状态一致性验证

```sql
-- 验证订单状态与任务状态的一致性
SELECT 
    order_no,
    status,
    like_status,
    share_status,
    follow_status,
    favorite_status,
    CASE 
        WHEN status = 'success' AND (
            like_status = 'success' AND 
            share_status = 'success' AND 
            follow_status = 'success' AND 
            favorite_status = 'success'
        ) THEN 'CONSISTENT'
        WHEN status = 'cancelled' AND (
            like_status = 'cancelled' AND 
            share_status = 'cancelled' AND 
            follow_status = 'cancelled' AND 
            favorite_status = 'cancelled'
        ) THEN 'CONSISTENT'
        ELSE 'INCONSISTENT'
    END as consistency_check
FROM orders 
WHERE is_system_order = 1 
AND created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
LIMIT 10;
```

### 3. 金额计算验证

```sql
-- 验证购买单金额计算是否正确
SELECT 
    order_no,
    amount,
    like_count,
    share_count,
    follow_count,
    favorite_count,
    ROUND(
        like_count * 0.1 + 
        share_count * 0.2 + 
        follow_count * 0.3 + 
        favorite_count * 0.4, 
        2
    ) as calculated_amount,
    CASE 
        WHEN ABS(amount - (like_count * 0.1 + share_count * 0.2 + follow_count * 0.3 + favorite_count * 0.4)) < 0.01 
        THEN 'CORRECT' 
        ELSE 'INCORRECT' 
    END as calculation_check
FROM orders 
WHERE is_system_order = 1 
AND created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
LIMIT 10;
```

## 总结

假订单生成系统的计算逻辑确保了：

1. **数据真实性**：基于真实价格配置计算金额
2. **状态一致性**：订单状态与任务状态保持逻辑一致
3. **分布合理性**：任务数量、状态分布符合真实场景
4. **时间连续性**：时间窗口设计避免数据断层
5. **可配置性**：支持灵活的参数配置和调整

通过这些逻辑，系统生成的假订单数据既保持了真实性，又提供了良好的用户体验。 