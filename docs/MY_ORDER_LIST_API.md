# 获取我的订单列表API文档

## 接口概述

获取我的订单列表接口用于查询当前登录用户的订单历史记录，支持按状态筛选和分页查询。

## 接口信息

- **接口路径**: `POST /api/v1/order/my-list`
- **认证方式**: 需要Bearer Token
- **Content-Type**: `application/json`

## 请求参数

### 请求体参数

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
- `3`: 全部 - 查询所有订单（包括拼单数据）

## 响应格式

### 成功响应

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "orders": [
      {
        "id": 1,
        "order_no": "ORD202501011200001234",
        "uid": "12345678",
        "period_number": "20241201001",
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
    ],
    "pagination": {
      "current_page": 1,
      "page_size": 10,
      "total": 25,
      "total_pages": 3,
      "has_next": true,
      "has_prev": false
    }
  },
  "timestamp": 1751365370
}
```

### 错误响应

```json
{
  "code": 401,
  "message": "未授权访问",
  "data": null,
  "timestamp": 1751365370
}
```

## 业务逻辑

1. **用户认证**: 验证JWT token的有效性
2. **用户状态检查**: 检查用户是否被删除或禁用
3. **数据查询**: 根据用户uid查询对应的订单数据
4. **状态筛选**: 根据status参数筛选不同状态的订单
5. **分页处理**: 支持分页查询，限制每页最大数量为20
6. **数据转换**: 将数据库记录转换为前端友好的响应格式

## 使用示例

### cURL示例

```bash
curl -X POST "http://localhost:8080/api/v1/order/my-list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 1
  }'
```

### JavaScript示例

```javascript
const response = await fetch('/api/v1/order/my-list', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    page: 1,
    page_size: 10,
    status: 1
  })
});

const data = await response.json();
console.log(data);
```

## 注意事项

1. **权限控制**: 该接口只能获取当前登录用户的订单，无法获取其他用户的订单
2. **分页限制**: 每页最大数量限制为20条，超出时自动调整为20
3. **状态筛选**: status=3时会查询所有订单，包括拼单数据
4. **时间格式**: 所有时间字段使用UTC时间格式
5. **订单过期**: 系统会自动计算订单是否过期和剩余时间

## 与原有接口的区别

该接口与原有的 `/api/v1/order/list` 接口功能完全相同，只是接口名称更加明确，表示只获取当前用户的订单列表。两个接口的入参、返回结构和业务逻辑完全一致。

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 操作成功 |
| 401 | 未授权访问（token无效或缺失） |
| 400 | 请求参数错误 |
| 500 | 服务器内部错误 