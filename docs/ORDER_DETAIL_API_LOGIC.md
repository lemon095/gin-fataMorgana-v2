# 订单详情接口逻辑详解

## 接口概述

订单详情接口 (`POST /api/v1/order/detail`) 用于获取指定订单的详细信息，包括订单状态、任务完成情况、时间信息等。

## 接口信息

- **接口路径**: `POST /api/v1/order/detail`
- **认证方式**: 需要Bearer Token
- **Content-Type**: `application/json`

## 请求参数

```json
{
  "order_no": "ORD202501011200001234"
}
```

### 参数说明

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `order_no` | string | ✅ | 订单编号 |

## 接口逻辑流程

### 1. 控制器层 (OrderController.GetOrderDetail)

```go
func (oc *OrderController) GetOrderDetail(c *gin.Context) {
    // 1. 参数验证
    var req models.GetOrderDetailRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.HandleValidationError(c, err)
        return
    }

    // 2. 用户认证检查
    userID := middleware.GetCurrentUser(c)
    if userID == 0 {
        utils.Unauthorized(c)
        return
    }

    // 3. 获取用户信息
    userRepo := database.NewUserRepository()
    var user models.User
    err := userRepo.FindByID(context.Background(), userID, &user)
    if err != nil {
        utils.ErrorWithMessage(c, utils.CodeDatabaseError, "获取用户信息失败")
        return
    }

    // 4. 用户状态检查
    if user.DeletedAt != nil {
        utils.ErrorWithMessage(c, utils.CodeUserNotFound, "账户已被删除，无法查询订单")
        return
    }
    if user.Status == 0 {
        utils.ErrorWithMessage(c, utils.CodeAccountLocked, "账户已被禁用，无法查询订单")
        return
    }

    // 5. 调用服务层获取订单详情
    response, err := oc.orderService.GetOrderDetail(&req, user.Uid)
    if err != nil {
        utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
        return
    }

    // 6. 返回成功响应
    utils.Success(c, response)
}
```

### 2. 服务层 (OrderService.GetOrderDetail)

```go
func (s *OrderService) GetOrderDetail(req *models.GetOrderDetailRequest, uid string) (*models.OrderResponse, error) {
    ctx := context.Background()

    // 1. 根据订单号查询订单
    order, err := s.orderRepo.FindOrderByOrderNo(ctx, req.OrderNo)
    if err != nil {
        return nil, utils.NewAppError(utils.CodeOrderDetailGetFailed, "获取订单详情失败")
    }

    // 2. 权限检查：确保订单属于当前用户
    if order.Uid != uid {
        return nil, utils.NewAppError(utils.CodeOrderAccessDenied, "无权访问此订单")
    }

    // 3. 转换为响应格式
    response := order.ToResponse()

    return &response, nil
}
```

### 3. 数据库层 (OrderRepository.FindOrderByOrderNo)

```go
func (r *OrderRepository) FindOrderByOrderNo(ctx context.Context, orderNo string) (*models.Order, error) {
    var order models.Order
    err := r.FindByCondition(ctx, map[string]interface{}{"order_no": orderNo}, &order)
    if err != nil {
        return nil, err
    }
    return &order, nil
}
```

## 数据转换逻辑

### Order.ToResponse() 方法

```go
func (o *Order) ToResponse() OrderResponse {
    return OrderResponse{
        ID:                 o.ID,
        OrderNo:            o.OrderNo,
        Uid:                o.Uid,
        Number:             o.PeriodNumber,
        Amount:             o.Amount,
        ProfitAmount:       o.ProfitAmount,
        Status:             o.Status,
        StatusName:         o.GetStatusName(),
        ExpireTime:         o.ExpireTime,
        LikeCount:          o.LikeCount,
        ShareCount:         o.ShareCount,
        FollowCount:        o.FollowCount,
        FavoriteCount:      o.FavoriteCount,
        LikeStatus:         o.LikeStatus,
        LikeStatusName:     o.GetTaskStatusName(o.LikeStatus),
        ShareStatus:        o.ShareStatus,
        ShareStatusName:    o.GetTaskStatusName(o.ShareStatus),
        FollowStatus:       o.FollowStatus,
        FollowStatusName:   o.GetTaskStatusName(o.FollowStatus),
        FavoriteStatus:     o.FavoriteStatus,
        FavoriteStatusName: o.GetTaskStatusName(o.FavoriteStatus),
        AuditorUid:         o.AuditorUid,
        CreatedAt:          o.CreatedAt,
        UpdatedAt:          o.UpdatedAt,
        IsExpired:          o.IsExpired(),
        RemainingTime:      o.GetRemainingTime(),
    }
}
```

### 状态名称转换

```go
// 订单状态名称
func (o *Order) GetStatusName() string {
    statusNames := map[string]string{
        OrderStatusPending:   "待处理",
        OrderStatusSuccess:   "成功",
        OrderStatusFailed:    "失败",
        OrderStatusCancelled: "已取消",
        OrderStatusExpired:   "已过期",
    }
    return statusNames[o.Status]
}

// 任务状态名称
func (o *Order) GetTaskStatusName(status string) string {
    statusNames := map[string]string{
        TaskStatusPending: "待完成",
        TaskStatusSuccess: "已完成",
    }
    return statusNames[status]
}
```

### 时间计算逻辑

```go
// 检查订单是否已过期
func (o *Order) IsExpired() bool {
    return time.Now().UTC().After(o.ExpireTime)
}

// 获取剩余时间（秒）
func (o *Order) GetRemainingTime() int64 {
    if o.IsExpired() {
        return 0
    }
    return int64(time.Until(o.ExpireTime).Seconds())
}
```

## 响应数据结构

### 成功响应

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "id": 1,
    "order_no": "ORD202501011200001234",
    "uid": "12345678",
    "period_number": "20250101001",
    "amount": 100.00,
    "profit_amount": 20.00,
    "status": "pending",
    "status_name": "待处理",
    "expire_time": "2025-01-01T12:05:00Z",
    "like_count": 5,
    "share_count": 2,
    "follow_count": 1,
    "favorite_count": 0,
    "like_status": "pending",
    "like_status_name": "待完成",
    "share_status": "pending",
    "share_status_name": "待完成",
    "follow_status": "pending",
    "follow_status_name": "待完成",
    "favorite_status": "success",
    "favorite_status_name": "已完成",
    "auditor_uid": "system",
    "created_at": "2025-01-01T12:00:00Z",
    "updated_at": "2025-01-01T12:00:00Z",
    "is_expired": false,
    "remaining_time": 300
  }
}
```

### 字段说明

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `id` | uint | 订单ID |
| `order_no` | string | 订单编号 |
| `uid` | string | 用户唯一ID |
| `period_number` | string | 期号 |
| `amount` | float64 | 订单金额 |
| `profit_amount` | float64 | 利润金额 |
| `status` | string | 订单状态 |
| `status_name` | string | 订单状态名称 |
| `expire_time` | string | 过期时间 |
| `like_count` | int | 点赞数量 |
| `share_count` | int | 分享数量 |
| `follow_count` | int | 关注数量 |
| `favorite_count` | int | 收藏数量 |
| `like_status` | string | 点赞完成状态 |
| `like_status_name` | string | 点赞状态名称 |
| `share_status` | string | 分享完成状态 |
| `share_status_name` | string | 分享状态名称 |
| `follow_status` | string | 关注完成状态 |
| `follow_status_name` | string | 关注状态名称 |
| `favorite_status` | string | 收藏完成状态 |
| `favorite_status_name` | string | 收藏状态名称 |
| `auditor_uid` | string | 审核员ID |
| `created_at` | string | 创建时间 |
| `updated_at` | string | 更新时间 |
| `is_expired` | bool | 是否已过期 |
| `remaining_time` | int64 | 剩余时间（秒） |

## 错误处理

### 常见错误码

| 错误码 | 说明 | 处理方式 |
|--------|------|----------|
| `401` | 未授权 | 需要登录 |
| `400` | 参数错误 | 检查请求参数 |
| `404` | 订单不存在 | 检查订单号是否正确 |
| `403` | 无权访问 | 只能查看自己的订单 |
| `500` | 服务器错误 | 联系技术支持 |

### 错误响应示例

```json
{
  "code": 404,
  "message": "获取订单详情失败",
  "data": null,
  "timestamp": 1751365370
}
```

## 安全机制

### 1. 用户认证
- 必须提供有效的JWT Token
- Token中必须包含有效的用户信息

### 2. 权限控制
- 只能查看自己的订单
- 通过比较订单的 `uid` 和当前用户的 `uid` 进行验证

### 3. 用户状态检查
- 检查用户是否已被删除
- 检查用户是否已被禁用

### 4. 参数验证
- 订单号不能为空
- 订单号格式验证

## 性能优化

### 1. 数据库查询优化
- 使用索引查询订单号
- 只查询必要的字段

### 2. 缓存策略
- 可以考虑对热门订单进行缓存
- 缓存订单状态和任务完成情况

### 3. 响应优化
- 实时计算剩余时间
- 动态生成状态名称

## 使用示例

### 前端调用示例

```javascript
// 获取订单详情
async function getOrderDetail(orderNo) {
  try {
    const response = await fetch('/api/v1/order/detail', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        order_no: orderNo
      })
    });
    
    const result = await response.json();
    
    if (result.code === 0) {
      return result.data;
    } else {
      throw new Error(result.message);
    }
  } catch (error) {
    console.error('获取订单详情失败:', error);
    throw error;
  }
}

// 使用示例
const orderDetail = await getOrderDetail('ORD202501011200001234');
console.log('订单详情:', orderDetail);
```

## 总结

订单详情接口是一个典型的RESTful API，具有以下特点：

1. **安全性**: 严格的用户认证和权限控制
2. **完整性**: 返回订单的所有相关信息
3. **实时性**: 动态计算过期时间和剩余时间
4. **可扩展性**: 支持多种订单状态和任务类型
5. **易用性**: 清晰的错误处理和响应格式

这个接口为前端提供了完整的订单信息，支持订单详情页面的展示和状态管理。 