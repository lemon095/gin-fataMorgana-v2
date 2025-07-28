# Gin-FataMorgana API 文档

## 概述

Gin-FataMorgana 是一个基于 Gin 框架的 Go Web 服务，提供用户认证、钱包管理、健康监控等功能。所有业务接口都使用 POST 请求，统一前缀为 `/api/v1`。

## 基础信息

- **基础 URL**: `http://localhost:9001`
- **API 版本**: v1
- **统一前缀**: `/api/v1`
- **请求方式**: 主要使用 POST
- **数据格式**: JSON
- **认证方式**: Bearer Token

## 通用响应格式

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {},
  "timestamp": 1751365370123
}
```

### 响应码说明

| 状态码 | 说明         |
| ------ | ------------ |
| 0      | 成功         |
| 401    | 认证失败     |
| 403    | 禁止访问     |
| 404    | 资源不存在   |
| 422    | 数据验证错误 |
| 500    | 服务器错误   |

### 时间戳说明

- **timestamp**: 服务器响应时间戳（毫秒级 Unix 时间戳）
- **示例**: `1751365370123` 表示 2025 年 7 月 4 日 08:56:10.123

## 1. 健康检查接口

### 1.1 系统健康检查

- **接口**: `GET /health`
- **说明**: 基础健康检查，用于监控
- **请求参数**: 无
- **返回示例**:

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "services": {
      "database": "healthy",
      "redis": "healthy"
    },
    "status": "healthy"
  },
  "timestamp": 1751365370
}
```

### 1.2 详细健康检查

- **接口**: `GET /api/v1/health/check`
- **说明**: 系统详细健康检查
- **请求参数**: 无
- **返回示例**: 同上

### 1.3 数据库健康检查

- **接口**: `GET /api/v1/health/database`
- **说明**: 数据库连接状态检查
- **请求参数**: 无
- **返回示例**:

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "status": "healthy",
    "connection": "ok"
  },
  "timestamp": 1751365370
}
```

### 1.4 Redis 健康检查

- **接口**: `GET /api/v1/health/redis`
- **说明**: Redis 连接状态检查
- **请求参数**: 无
- **返回示例**:

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "status": "healthy",
    "connection": "ok"
  },
  "timestamp": 1751365370
}
```

## 2. 假数据接口

### 2.1 实时动态假数据

- **接口**: `POST /api/v1/fake/activities`
- **说明**: 获取实时动态假数据，用于前端展示
- **请求参数**:

```json
{
  "count": 10
}
```

- **参数说明**:
  - `count` (可选): 返回数据条数，默认 10 条，最大 50 条
- **返回示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "uid": "12***567",
      "time": "14:30",
      "amount": 456.78,
      "type": "点赞"
    },
    {
      "uid": "34***890",
      "time": "14:29",
      "amount": 789.12,
      "type": "关注"
    }
  ]
}
```

## 3. 认证相关接口

### 3.1 用户注册

- **接口**: `POST /api/v1/auth/register`
- **说明**: 用户注册
- **请求参数**:

```json
{
  "email": "user@example.com",
  "password": "123456",
  "confirm_password": "123456",
  "invite_code": "INVITE123"
}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "用户注册成功",
  "data": {
    "user": {
      "id": 1,
      "uid": "12345678",
      "username": "user123",
      "email": "u***@example.com",
      "phone": "",
      "bank_card_info": "",
      "experience": 0,
      "credit_score": 100,
      "status": 1,
      "invited_by": "7TRABJ",
      "has_group_buy_qualification": false,
      "created_at": "2025-01-01T00:00:00Z"
    }
  },
  "timestamp": 1751365370
}
```

### 3.2 用户登录

- **接口**: `POST /api/v1/auth/login`
- **说明**: 用户登录
- **请求参数**:

```json
{
  "email": "user@example.com",
  "password": "123456"
}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "token_type": "Bearer",
      "expires_in": 3600
    }
  },
  "timestamp": 1751365370
}
```

### 3.3 刷新令牌

- **接口**: `POST /api/v1/auth/refresh`
- **说明**: 刷新访问令牌
- **请求参数**:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "令牌刷新成功",
  "data": {
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "token_type": "Bearer",
      "expires_in": 3600
    }
  },
  "timestamp": 1751365370
}
```

### 3.4 用户登出

- **接口**: `POST /api/v1/auth/logout`
- **说明**: 用户登出
- **请求参数**: 无
- **返回示例**:

```json
{
  "code": 0,
  "message": "登出成功",
  "data": null,
  "timestamp": 1751365370
}
```

### 3.5 获取用户信息

- **接口**: `POST /api/v1/auth/profile`
- **说明**: 获取当前用户信息
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "user": {
      "id": 1,
      "uid": "12345678",
      "username": "user123",
      "email": "u***@example.com",
      "phone": "",
      "bank_card_info": "",
      "experience": 0,
      "credit_score": 100,
      "status": 1,
      "invited_by": "7TRABJ",
      "has_group_buy_qualification": false,
      "created_at": "2025-01-01T00:00:00Z"
    }
  },
  "timestamp": 1751365370
}
```

### 3.6 绑定银行卡

- **接口**: `POST /api/v1/auth/bind-bank-card`
- **说明**: 绑定银行卡
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{
  "bank_name": "中国工商银行",
  "card_holder": "张三",
  "card_number": "6222021234567890123",
  "card_type": "借记卡"
}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "银行卡绑定成功",
  "data": {
    "user": {
      "id": 1,
      "uid": "12345678",
      "username": "user123",
      "email": "u***@example.com",
      "phone": "",
      "bank_card_info": "{\"card_number\":\"6222021234567890123\",\"card_type\":\"借记卡\",\"bank_name\":\"中国工商银行\",\"card_holder\":\"张三\"}",
      "experience": 0,
      "credit_score": 100,
      "status": 1,
      "invited_by": "7TRABJ",
      "has_group_buy_qualification": false,
      "created_at": "2025-01-01T00:00:00Z"
    }
  },
  "timestamp": 1751365370
}
```

### 3.7 获取银行卡信息

- **接口**: `POST /api/v1/auth/bank-card`
- **说明**: 获取当前用户银行卡信息
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "bank_card_info": {
      "card_number": "6222021234567890123",
      "card_type": "借记卡",
      "bank_name": "中国工商银行",
      "card_holder": "张三"
    }
  },
  "timestamp": 1751365370
}
```

## 4. 会话管理接口

### 4.1 检查登录状态

- **接口**: `POST /api/v1/session/status`
- **说明**: 检查当前登录状态
- **请求参数**:

```json
{}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "获取登录状态成功",
  "data": {
    "is_authenticated": true,
    "user_id": 1,
    "username": "user123",
    "timestamp": 1751365370
  },
  "timestamp": 1751365370
}
```

### 4.2 获取当前用户信息

- **接口**: `POST /api/v1/session/user`
- **说明**: 获取当前会话用户信息
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "获取用户信息成功",
  "data": {
    "user_id": 1,
    "username": "user123",
    "login_time": 1751365370
  },
  "timestamp": 1751365370
}
```

### 4.3 用户登出

- **接口**: `POST /api/v1/session/logout`
- **说明**: 用户登出
- **请求参数**: 无
- **返回示例**:

```json
{
  "code": 0,
  "message": "登出成功",
  "data": {
    "logout_time": 1751365370
  },
  "timestamp": 1751365370
}
```

### 4.4 刷新会话

- **接口**: `POST /api/v1/session/refresh`
- **说明**: 刷新会话
- **认证**: 需要 Bearer Token
- **请求参数**: 无
- **返回示例**:

```json
{
  "code": 0,
  "message": "会话刷新成功",
  "data": {
    "refresh_time": 1751365370,
    "user_id": 1,
    "username": "user123"
  },
  "timestamp": 1751365370
}
```

## 5. 钱包相关接口

### 5.1 获取钱包信息

- **接口**: `POST /api/v1/wallet/info`
- **说明**: 获取当前用户钱包信息
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "id": 1,
    "uid": "12345678",
    "balance": 1000.0,
    "frozen_amount": 0.0,
    "total_recharge": 2000.0,
    "total_withdraw": 1000.0,
    "status": 1,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  },
  "timestamp": 1751365370
}
```

### 5.2 获取资金记录

- **接口**: `POST /api/v1/wallet/transactions`
- **说明**: 获取用户资金记录
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{
  "page": 1,
  "page_size": 10
}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "transactions": [
      {
        "id": 1,
        "transaction_no": "TX202501010001",
        "uid": "12345678",
        "type": "recharge",
        "amount": 1000.0,
        "balance": 1000.0,
        "description": "充值",
        "status": "completed",
        "operator_uid": "system",
        "created_at": "2025-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  },
  "timestamp": 1751365370
}
```

### 5.3 申请提现

- **接口**: `POST /api/v1/wallet/withdraw`
- **说明**: 申请提现
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{
  "amount": 100.0,
  "bank_card_info": {
    "card_number": "6222021234567890123",
    "card_type": "借记卡",
    "bank_name": "中国工商银行",
    "card_holder": "张三"
  },
  "description": "提现申请"
}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "transaction_no": "TX202501010002",
    "amount": 100.0,
    "status": "pending",
    "created_at": "2025-01-01T00:00:00Z"
  },
  "timestamp": 1751365370
}
```

### 5.4 获取提现汇总

- **接口**: `POST /api/v1/wallet/withdraw-summary`
- **说明**: 获取提现汇总信息
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "total_withdraw": 1000.0,
    "pending_withdraw": 100.0,
    "completed_withdraw": 900.0,
    "withdraw_count": 5
  },
  "timestamp": 1751365370
}
```

### 5.5 充值申请

- **接口**: `POST /api/v1/wallet/recharge-apply`
- **说明**: 充值申请
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{
  "uid": "12345678",
  "amount": 500.0,
  "description": "充值申请"
}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "充值申请已提交",
  "data": {
    "transaction_no": "TX202501010003"
  },
  "timestamp": 1751365370
}
```

### 5.6 充值确认

- **接口**: `POST /api/v1/wallet/recharge-confirm`
- **说明**: 充值确认
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{
  "transaction_no": "TX202501010003"
}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "充值已到账",
  "data": null,
  "timestamp": 1751365370
}
```

## 6. 订单相关接口

### 6.1 获取订单列表

- **接口**: `POST /api/v1/order/list`
- **说明**: 获取用户订单列表
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{
  "page": 1,
  "page_size": 10
}
```

- **参数说明**:
  - `page`: 页码，从 1 开始
  - `page_size`: 每页大小，最大 100
- **返回示例**:

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "orders": [
      {
        "id": 1,
        "order_no": "ORD12345678",
        "uid": "12345678",
        "buy_amount": 1000.0,
        "profit_amount": 150.0,
        "status": "success",
        "status_name": "成功",
        "description": "股票买入",
        "remark": "",
        "created_at": "2025-01-01T12:00:00Z"
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

### 6.2 创建订单

- **接口**: `POST /api/v1/order/create`
- **说明**: 创建新订单（uid 从 token 中自动获取，无需在请求中传递）
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{
  "amount": 100.0,
  "profit_amount": 10.0,
  "like_count": 5,
  "share_count": 2,
  "follow_count": 1,
  "favorite_count": 3
}
```

- **参数说明**:
  - `amount`: 订单金额，必须大于 0
  - `profit_amount`: 预期利润金额，必须大于等于 0
  - `like_count`: 点赞数量，必须大于等于 0
  - `share_count`: 分享数量，必须大于等于 0
  - `follow_count`: 关注数量，必须大于等于 0
  - `favorite_count`: 收藏数量，必须大于等于 0
  - `uid`: 用户 ID（从 token 中自动获取，无需传递）
- **返回示例**:

```json
{
  "code": 0,
  "message": "订单创建成功",
  "data": {
    "order_no": "ORD12345678",
    "amount": 100.0,
    "status": "pending",
    "message": "订单创建成功"
  },
  "timestamp": 1751365370
}
```

### 6.3 获取订单详情

- **接口**: `POST /api/v1/order/detail`
- **说明**: 获取订单详情
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{
  "order_no": "ORD12345678"
}
```

- **参数说明**:
  - `order_no`: 订单编号
- **返回示例**: 同订单列表中的单个订单数据

### 6.4 获取订单统计

- **接口**: `POST /api/v1/order/stats`
- **说明**: 获取用户订单统计信息
- **认证**: 需要 Bearer Token
- **请求参数**: 无
- **返回示例**:

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "total_orders": 25,
    "success_orders": 20,
    "pending_orders": 3,
    "failed_orders": 2,
    "total_buy_amount": 50000.0,
    "total_profit_amount": 7500.0
  },
  "timestamp": 1751365370
}
```

### 6.5 根据状态获取订单

- **接口**: `POST /api/v1/order/by-status`
- **说明**: 根据状态筛选订单
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{
  "page": 1,
  "page_size": 10,
  "status": "success"
}
```

- **参数说明**:
  - `page`: 页码，从 1 开始
  - `page_size`: 每页大小，最大 100
  - `status`: 订单状态（pending/success/failed/cancelled）
- **返回示例**: 同订单列表

### 6.6 根据日期范围获取订单

- **接口**: `POST /api/v1/order/by-date`
- **说明**: 根据日期范围筛选订单
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{
  "page": 1,
  "page_size": 10,
  "start_date": "2025-01-01",
  "end_date": "2025-01-31"
}
```

- **参数说明**:
  - `page`: 页码，从 1 开始
  - `page_size`: 每页大小，最大 100
  - `start_date`: 开始日期（YYYY-MM-DD）
  - `end_date`: 结束日期（YYYY-MM-DD）
- **返回示例**: 同订单列表

## 7. 热榜接口

### 7.1 获取任务热榜

- **接口**: `POST /api/v1/leaderboard/ranking`
- **说明**: 获取任务热榜排行榜列表和当前用户数据（优化版本，不使用窗口函数）
- **认证**: 需要 Bearer Token
- **请求参数**: 无（空 JSON 对象）

```json
{}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "获取热榜数据成功",
  "data": {
    "week_start": "2025-01-06T00:00:00Z",
    "week_end": "2025-01-12T23:59:59Z",
    "my_rank": {
      "id": 1,
      "uid": "12345678",
      "username": "张*三",
      "completed_at": "2025-01-10T15:30:00Z",
      "order_count": 25,
      "total_amount": 2500.0,
      "total_profit": 500.0,
      "rank": 5,
      "is_rank": true
    },
    "top_users": [
      {
        "id": 1,
        "uid": "87654321",
        "username": "李*四",
        "completed_at": "2025-01-10T16:45:00Z",
        "order_count": 50,
        "total_amount": 5000.0,
        "total_profit": 1000.0,
        "rank": 1,
        "is_rank": true
      }
    ],
    "cache_expire": "2025-01-10T16:50:00Z"
  }
}
```

**数据结构说明**:

- `week_start`: 本周开始时间（周一）
- `week_end`: 本周结束时间（周日）
- `my_rank`: 当前用户排名信息
  - `id`: 用户 ID
  - `uid`: 用户唯一标识
  - `username`: 脱敏后的用户名
  - `completed_at`: 最新完成订单时间
  - `order_count`: 完成订单数量
  - `total_amount`: 总金额
  - `total_profit`: 总利润
  - `rank`: 排名（999 表示未上榜）
  - `is_rank`: 是否在前 10 名榜单上
- `top_users`: 前 10 名用户列表
- `cache_expire`: 缓存过期时间（5 分钟缓存）

**排名规则**:

1. 完成订单数量（降序）
2. 总金额（降序）
3. 最新完成时间（升序）

**如果用户不在排行榜上**:

```json
{
  "code": 0,
  "message": "获取热榜数据成功",
  "data": {
    "ranking_list": [
      // ... 排行榜数据（共10条）
    ],
    "my_data": {
      "rank": 0,
      "user_id": 9999,
      "is_ranked": false,
      "entry": null
    }
  },
  "timestamp": 1751365370
}
```

## 8. 管理员接口

### 8.1 确认提现

- **接口**: `POST /api/v1/admin/withdraw/confirm`
- **说明**: 管理员确认提现
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{
  "transaction_no": "TX202501010002"
}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "提现确认成功",
  "data": null,
  "timestamp": 1751365370
}
```

### 8.2 取消提现

- **接口**: `POST /api/v1/admin/withdraw/cancel`
- **说明**: 管理员取消提现
- **认证**: 需要 Bearer Token
- **请求参数**:

```json
{
  "transaction_no": "TX202501010002",
  "reason": "银行卡信息有误"
}
```

- **返回示例**:

```json
{
  "code": 0,
  "message": "提现取消成功",
  "data": null,
  "timestamp": 1751365370
}
```

## 9. 错误码说明

| 错误码 | 说明         |
| ------ | ------------ |
| 0      | 成功         |
| 401    | 认证失败     |
| 403    | 禁止访问     |
| 404    | 资源不存在   |
| 422    | 数据验证错误 |
| 500    | 服务器错误   |
| 1001   | 数据库错误   |
| 1002   | Redis 错误   |
| 1003   | 参数错误     |
| 1004   | 操作失败     |
| 1005   | 用户不存在   |
| 1006   | 用户已存在   |
| 1007   | 验证失败     |
| 1008   | 账户锁定     |
| 1009   | 注册关闭     |

## 10. 使用示例

### 9.1 用户注册和登录流程

```bash
# 1. 用户注册
curl -X POST http://localhost:9001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "7TRABJ"
  }'

# 2. 用户登录
curl -X POST http://localhost:9001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "123456"
  }'

# 3. 获取用户信息（需要token）
curl -X POST http://localhost:9001/api/v1/auth/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{}'
```

### 10.2 钱包操作流程

```bash
# 1. 获取钱包信息
curl -X POST http://localhost:9001/api/v1/wallet/info \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{}'

# 2. 申请提现
curl -X POST http://localhost:9001/api/v1/wallet/withdraw \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "amount": 100.00,
    "bank_card_info": {
      "card_number": "6222021234567890123",
      "card_type": "借记卡",
      "bank_name": "中国工商银行",
      "card_holder": "张三"
    },
    "description": "提现申请"
  }'
```

### 10.3 热榜接口使用示例

```bash
# 获取任务热榜（用户ID: 1001 - 在榜上）
curl -X POST http://localhost:9001/api/v1/leaderboard/ranking \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1001
  }'

# 获取任务热榜（用户ID: 9999 - 不在榜上）
curl -X POST http://localhost:9001/api/v1/leaderboard/ranking \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 9999
  }'

# 测试无效用户ID
curl -X POST http://localhost:9001/api/v1/leaderboard/ranking \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 0
  }'
```

## 11. 注意事项

1. **认证**: 需要认证的接口必须在请求头中包含 `Authorization: Bearer <token>`
2. **请求格式**: 所有请求都使用 JSON 格式
3. **响应格式**: 所有响应都遵循统一的 JSON 格式
4. **错误处理**: 请根据返回的 code 字段判断请求是否成功
5. **分页**: 支持分页的接口使用 page 和 page_size 参数
6. **限流**: 部分接口有访问频率限制，请合理使用

## 12. 更新日志

- **v1.0.0**: 初始版本，所有业务接口改为 POST 请求
- 统一 API 前缀为 `/api/v1`
- 健康检查接口保持 GET 请求以便监控
- 完善错误码和响应格式
- **v1.1.0**: 新增任务热榜接口
  - 添加 `/api/v1/leaderboard/ranking` 接口
  - 支持获取排行榜列表和当前用户数据
  - 包含排名、用户名、金额、完成单数、利润金额、时间等字段
