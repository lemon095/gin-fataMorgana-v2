# 金额配置接口文档

## 接口概述

金额配置接口用于获取充值、提现等操作的金额配置列表，支持按类型查询和排序。**只返回激活状态（is_active = true）的配置数据**。

## 1. 获取金额配置列表

### 接口信息
- **接口路径**: `POST /api/v1/amount-config/list`
- **请求方式**: POST
- **认证要求**: 需要Bearer Token登录校验
- **Content-Type**: `application/json`

### 请求参数

#### 请求头
```
Authorization: Bearer <your_access_token>
Content-Type: application/json
```

#### 请求体
```json
{
  "type": "recharge"
}
```

#### 参数说明
| 参数名 | 类型 | 必填 | 说明 | 示例值 |
|--------|------|------|------|--------|
| type | string | 是 | 配置类型 | recharge / withdraw |

**type参数可选值**:
- `recharge`: 充值配置
- `withdraw`: 提现配置

### 返回数据

#### 成功响应 (200)
```json
{
  "code": 0,
  "message": "获取金额配置列表成功",
  "data": [
    {
      "id": 1,
      "type": "recharge",
      "amount": 100.00,
      "description": "充值100元",
      "is_active": true,
      "sort_order": 1,
      "created_at": "2024-01-15 14:30:00",
      "updated_at": "2024-01-15 14:30:00"
    },
    {
      "id": 2,
      "type": "recharge",
      "amount": 200.00,
      "description": "充值200元",
      "is_active": true,
      "sort_order": 2,
      "created_at": "2024-01-15 14:30:00",
      "updated_at": "2024-01-15 14:30:00"
    },
    {
      "id": 3,
      "type": "recharge",
      "amount": 500.00,
      "description": "充值500元",
      "is_active": true,
      "sort_order": 3,
      "created_at": "2024-01-15 14:30:00",
      "updated_at": "2024-01-15 14:30:00"
    }
  ],
  "timestamp": 1705312800
}
```

#### 数据结构说明
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int64 | 配置ID |
| type | string | 配置类型(recharge/withdraw) |
| amount | float64 | 金额 |
| description | string | 描述信息 |
| is_active | bool | 是否激活（接口只返回true的数据） |
| sort_order | int | 排序值 |
| created_at | string | 创建时间 |
| updated_at | string | 更新时间 |

### 错误响应

#### 未登录 (401)
```json
{
  "code": 401,
  "message": "认证失败",
  "timestamp": 1705312800
}
```

#### 参数错误 (422)
```json
{
  "code": 422,
  "message": "请求参数错误: Key: 'AmountConfigRequest.Type' Error:Field validation for 'Type' failed on the 'oneof' tag",
  "timestamp": 1705312800
}
```

#### 缺少参数 (422)
```json
{
  "code": 422,
  "message": "请求参数错误: Key: 'AmountConfigRequest.Type' Error:Field validation for 'Type' failed on the 'required' tag",
  "timestamp": 1705312800
}
```

#### 服务器错误 (500)
```json
{
  "code": 500,
  "message": "服务器内部错误",
  "timestamp": 1705312800
}
```

## 2. 获取金额配置详情

### 接口信息
- **接口路径**: `GET /api/v1/amount-config/{id}`
- **请求方式**: GET
- **认证要求**: 需要Bearer Token登录校验
- **Content-Type**: `application/json`

### 请求参数

#### 路径参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 配置ID |

#### 请求头
```
Authorization: Bearer <your_access_token>
Content-Type: application/json
```

### 返回数据

#### 成功响应 (200)
```json
{
  "code": 0,
  "message": "获取金额配置详情成功",
  "data": {
    "id": 1,
    "type": "recharge",
    "amount": 100.00,
    "description": "充值100元",
    "is_active": true,
    "sort_order": 1,
    "created_at": "2024-01-15 14:30:00",
    "updated_at": "2024-01-15 14:30:00"
  },
  "timestamp": 1705312800
}
```

## 使用示例

### cURL示例

#### 获取充值配置列表
```bash
curl -X POST "http://localhost:9001/api/v1/amount-config/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "type": "recharge"
  }'
```

#### 获取提现配置列表
```bash
curl -X POST "http://localhost:9001/api/v1/amount-config/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "type": "withdraw"
  }'
```

#### 获取配置详情
```bash
curl -X GET "http://localhost:9001/api/v1/amount-config/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### JavaScript示例

#### 获取充值配置列表
```javascript
const response = await fetch('http://localhost:9001/api/v1/amount-config/list', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer YOUR_ACCESS_TOKEN'
  },
  body: JSON.stringify({
    type: 'recharge'
  })
});

const data = await response.json();
console.log(data);
```

## 注意事项

1. **认证要求**: 所有接口都需要有效的Bearer Token
2. **参数验证**: type参数只能是 `recharge` 或 `withdraw`
3. **排序规则**: 返回结果按 `sort_order` 升序，然后按 `amount` 升序排列
4. **激活状态**: 只返回 `is_active = true` 的配置，未激活的配置不会返回
5. **错误处理**: 请根据返回的 `code` 字段判断请求是否成功

## 测试脚本

项目提供了测试脚本：
```bash
# 运行金额配置接口测试
./test_scripts/test_amount_config.sh

# 运行数据库连接检查
./test_scripts/check_db_connection.sh
``` 