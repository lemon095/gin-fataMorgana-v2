# 订单表新增 is_system_order 字段修改文档

## 概述

本次修改为订单表新增了 `is_system_order` 字段，用于标识订单是否为系统订单。

## 修改内容

### 1. 数据库迁移文件

**文件**: `database/migrations/create_orders_table.sql`

```sql
`is_system_order` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否系统订单 0-否 1-是'
```

- 字段类型: `tinyint(1)` (布尔值)
- 默认值: `0` (false)
- 注释: 是否系统订单 0-否 1-是
- 索引: 已添加索引 `idx_is_system_order`

### 2. 模型结构体更新

**文件**: `models/order.go`

#### Order 结构体
```go
type Order struct {
    // ... 其他字段 ...
    IsSystemOrder  bool      `json:"is_system_order" gorm:"default:false;comment:是否系统订单"`
    // ... 其他字段 ...
}
```

#### OrderResponse 结构体
```go
type OrderResponse struct {
    // ... 其他字段 ...
    IsSystemOrder      bool      `json:"is_system_order"`
    // ... 其他字段 ...
}
```

#### ToResponse 方法更新
```go
func (o *Order) ToResponse() OrderResponse {
    return OrderResponse{
        // ... 其他字段 ...
        IsSystemOrder:      o.IsSystemOrder,
        // ... 其他字段 ...
    }
}
```

### 3. 服务层更新

**文件**: `services/order_service.go`

#### 订单创建逻辑
```go
// 创建订单对象
order := &models.Order{
    // ... 其他字段 ...
    IsSystemOrder: false, // 默认为用户订单，不是系统订单
}
```

#### 拼单列表转换
```go
orderResponse := models.OrderResponse{
    // ... 其他字段 ...
    IsSystemOrder: false, // 拼单不是系统订单
    // ... 其他字段 ...
}
```

**文件**: `services/group_buy_service.go`

#### 拼单订单创建
```go
order := &models.Order{
    // ... 其他字段 ...
    IsSystemOrder:  false, // 拼单订单也是用户订单，不是系统订单
    // ... 其他字段 ...
}
```

## 字段说明

### is_system_order 字段含义

- `false` (0): 用户订单 - 由用户主动创建的订单
- `true` (1): 系统订单 - 由系统自动生成的订单

### 使用场景

1. **用户订单** (is_system_order = false)
   - 用户通过正常流程创建的订单
   - 拼单参与生成的订单
   - 用户主动购买的任务订单

2. **系统订单** (is_system_order = true)
   - 系统自动生成的订单
   - 活动奖励订单
   - 补偿订单
   - 系统测试订单

## 影响范围

### 1. API 接口

以下接口的响应中现在包含 `is_system_order` 字段：

- `POST /api/v1/order/create` - 创建订单
- `GET /api/v1/order/list` - 获取订单列表
- `GET /api/v1/order/my-orders` - 获取我的订单列表
- `GET /api/v1/order/detail` - 获取订单详情
- `GET /api/v1/order/list?status=3` - 获取拼单列表

### 2. 数据库查询

- 所有订单查询都会包含该字段
- 可以通过该字段筛选系统订单或用户订单
- 已添加数据库索引，查询性能良好

### 3. 前端展示

前端可以根据 `is_system_order` 字段：
- 区分显示不同类型的订单
- 应用不同的样式或图标
- 控制不同的操作权限

## 测试验证

### 测试脚本

**文件**: `test_scripts/test_is_system_order.sh`

测试内容包括：
1. 创建订单时字段设置
2. 订单列表查询字段返回
3. 拼单列表字段返回
4. 数据库字段存在性验证

### 运行测试

```bash
chmod +x test_scripts/test_is_system_order.sh
./test_scripts/test_is_system_order.sh
```

## 注意事项

### 1. 向后兼容性

- 新增字段有默认值，不会影响现有数据
- 现有订单的 `is_system_order` 字段默认为 `false`
- API 响应格式向后兼容

### 2. 数据迁移

如果数据库中已有订单数据，需要执行以下 SQL 更新默认值：

```sql
UPDATE orders SET is_system_order = 0 WHERE is_system_order IS NULL;
```

### 3. 权限控制

- 系统订单可能需要特殊的权限控制
- 建议在相关业务逻辑中添加权限检查

## 后续扩展

### 1. 系统订单创建

可以添加创建系统订单的接口：

```go
func (s *OrderService) CreateSystemOrder(req *CreateSystemOrderRequest) (*CreateOrderResponse, error) {
    // 创建系统订单的逻辑
    order := &models.Order{
        // ... 其他字段 ...
        IsSystemOrder: true, // 标记为系统订单
    }
}
```

### 2. 订单筛选

可以添加按订单类型筛选的接口：

```go
func (s *OrderService) GetOrdersByType(req *GetOrdersByTypeRequest) (*GetOrderListResponse, error) {
    // 根据 is_system_order 字段筛选订单
}
```

## 总结

本次修改成功为订单表新增了 `is_system_order` 字段，实现了订单类型的区分功能。修改涉及数据库、模型、服务层等多个层面，确保了数据的一致性和 API 的完整性。所有相关代码都已更新，并提供了测试脚本进行验证。 