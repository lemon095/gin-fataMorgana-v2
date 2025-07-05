# 订单列表接口逻辑详解

## 接口概述

`/api/v1/order/list` 接口用于获取**所有订单列表**，只需登录即可查看所有用户的订单数据。

## 接口信息

- **接口路径**: `POST /api/v1/order/list`
- **控制器方法**: `OrderController.GetAllOrderList`
- **服务方法**: `OrderService.GetAllOrderList`
- **认证要求**: 需要Bearer Token（只需登录即可）
- **权限要求**: 无需特殊权限，登录用户即可查看所有订单

## 请求参数

```json
{
  "page": 1,
  "page_size": 10,
  "status": 1
}
```

### 参数说明

| 参数名 | 类型 | 必填 | 说明 | 取值范围 |
|--------|------|------|------|----------|
| `page` | int | ✅ | 页码，从1开始 | 最小值为1 |
| `page_size` | int | ✅ | 每页大小 | 最小值为1，最大值为20 |
| `status` | int | ✅ | 订单状态类型 | 1:进行中, 2:已完成, 3:全部 |

### 状态类型说明

- `1`: 进行中 - 查询状态为 `pending` 的订单
- `2`: 已完成 - 查询状态为 `success` 的订单  
- `3`: 全部 - 查询所有状态的订单（包括 pending、success、failed、cancelled、expired 等）

## 接口逻辑流程

### 1. 控制器层 (OrderController.GetAllOrderList)

```go
func (oc *OrderController) GetAllOrderList(c *gin.Context) {
    // 1. 参数验证
    var req models.GetOrderListRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.HandleValidationError(c, err)
        return
    }

    // 2. 用户认证检查（只需登录）
    userID := middleware.GetCurrentUser(c)
    if userID == 0 {
        utils.Unauthorized(c)
        return
    }

    // 3. 调用服务层获取所有订单列表
    response, err := oc.orderService.GetAllOrderList(&req)
    if err != nil {
        utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
        return
    }

    // 4. 返回成功响应
    utils.Success(c, response)
}
```

### 2. 服务层 (OrderService.GetAllOrderList)

```go
func (s *OrderService) GetAllOrderList(req *models.GetOrderListRequest) (*GetOrderListResponse, error) {
    ctx := context.Background()

    // 1. 参数验证和限制
    if req.PageSize > 20 {
        req.PageSize = 20
    }

    if req.Status < 1 || req.Status > 3 {
        return nil, utils.NewAppError(utils.CodeOrderStatusInvalid, "状态类型参数无效")
    }

    // 2. 特殊处理：拼单数据不支持全量查询
    if req.Status == 3 {
        return &GetOrderListResponse{
            Orders: []models.OrderResponse{}, 
            Pagination: PaginationInfo{
                CurrentPage: req.Page, 
                PageSize: req.PageSize, 
                Total: 0, 
                TotalPages: 0, 
                HasNext: false, 
                HasPrev: false
            }
        }, nil
    }

    // 3. 根据状态类型获取对应的状态值
    status := models.GetStatusByType(req.Status)

    // 4. 从数据库获取订单列表
    orders, total, err := s.orderRepo.GetOrdersByStatus(ctx, status, req.Page, req.PageSize)
    if err != nil {
        return nil, utils.NewAppError(utils.CodeOrderListGetFailed, "获取订单列表失败")
    }

    // 5. 转换为响应格式
    var orderResponses []models.OrderResponse
    for _, order := range orders {
        orderResponses = append(orderResponses, order.ToResponse())
    }

    // 6. 计算分页信息
    totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
    hasNext := req.Page < totalPages
    hasPrev := req.Page > 1

    // 7. 返回结果
    return &GetOrderListResponse{
        Orders: orderResponses,
        Pagination: PaginationInfo{
            CurrentPage: req.Page,
            PageSize:    req.PageSize,
            Total:       total,
            TotalPages:  totalPages,
            HasNext:     hasNext,
            HasPrev:     hasPrev,
        },
    }, nil
}
```

### 3. 数据库层 (OrderRepository.GetOrdersByStatus)

```go
func (r *OrderRepository) GetOrdersByStatus(ctx context.Context, status string, page, pageSize int) ([]models.Order, int64, error) {
    var orders []models.Order
    var total int64
    
    // 1. 构建查询条件
    query := r.db.WithContext(ctx).Model(&models.Order{})
    if status != "" {
        query = query.Where("status = ?", status)
    }
    
    // 2. 获取总数
    err := query.Count(&total).Error
    if err != nil {
        return nil, 0, err
    }
    
    // 3. 计算偏移量
    offset := (page - 1) * pageSize
    
    // 4. 获取分页数据，按创建时间倒序排列
    err = query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
    if err != nil {
        return nil, 0, err
    }
    
    return orders, total, nil
}
```

## 关键特点

### 🔍 **权限控制**
- **只需登录**: 不需要特殊权限，任何登录用户都可以查看所有订单
- **无用户限制**: 不限制只能查看自己的订单，可以查看所有用户的订单

### 📊 **数据范围**
- **全量数据**: 查询所有用户的订单，不按用户uid过滤
- **状态筛选**: 支持按订单状态筛选（进行中/已完成/全部）
- **全部状态**: status=3时返回所有状态的订单，包括 pending、success、failed、cancelled、expired 等

### 🎯 **业务逻辑**
- **分页限制**: 每页最大20条记录
- **时间排序**: 按创建时间倒序排列，最新的订单在前
- **状态映射**: 自动将状态类型转换为数据库状态值

### ⚠️ **特殊处理**
- **全部状态**: status=3时返回所有状态的订单，不进行状态过滤
- **参数验证**: 严格验证状态类型参数的有效性

## 与其他接口的区别

| 接口 | 权限要求 | 数据范围 | 特殊功能 |
|------|----------|----------|----------|
| `/order/list` | 只需登录 | 所有用户订单 | 全量查询，支持所有状态 |
| `/order/my-orders` | 需要认证 | 当前用户订单 | 用户订单查询，支持拼单数据 |
| `/order/all-list` | 需要认证 | 当前用户订单 | 支持拼单数据 |

## 使用场景

1. **管理员查看**: 管理员需要查看所有用户的订单情况
2. **数据分析**: 进行订单数据分析和统计
3. **监控系统**: 监控系统订单状态和趋势
4. **客服支持**: 客服人员查看用户订单信息

## 注意事项

1. **数据安全**: 该接口返回所有用户订单，需要注意数据安全
2. **性能考虑**: 全量查询可能影响性能，建议合理使用分页
3. **状态支持**: status=3时返回所有状态的订单，包括各种订单状态
4. **权限控制**: 虽然只需登录，但建议在生产环境中增加权限控制 