# Gin-FataMorgana API 文档

## 概述

Gin-FataMorgana 是一个基于Gin框架的Go Web服务，提供用户认证、钱包管理、健康监控等功能。所有业务接口都使用POST请求，统一前缀为 `/api/v1`。

## 基础信息

- **基础URL**: `http://localhost:9001`
- **API版本**: v1
- **统一前缀**: `/api/v1`
- **请求方式**: 主要使用POST
- **数据格式**: JSON
- **认证方式**: Bearer Token

## 通用响应格式

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {},
  "timestamp": 1751365370
}
```

### 响应码说明

| 状态码 | 说明 |
|--------|------|
| 0 | 成功 |
| 401 | 认证失败 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 422 | 数据验证错误 |
| 500 | 服务器错误 |

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

### 1.4 Redis健康检查
- **接口**: `GET /api/v1/health/redis`
- **说明**: Redis连接状态检查
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

## 2. 认证相关接口

### 2.1 用户注册
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
      "invited_by": "INVITE123",
      "has_group_buy_qualification": false,
      "created_at": "2025-01-01T00:00:00Z"
    }
  },
  "timestamp": 1751365370
}
```

### 2.2 用户登录
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

### 2.3 刷新令牌
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

### 2.4 用户登出
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

### 2.5 获取用户信息
- **接口**: `POST /api/v1/auth/profile`
- **说明**: 获取当前用户信息
- **认证**: 需要Bearer Token
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
      "invited_by": "INVITE123",
      "has_group_buy_qualification": false,
      "created_at": "2025-01-01T00:00:00Z"
    }
  },
  "timestamp": 1751365370
}
```

### 2.6 绑定银行卡
- **接口**: `POST /api/v1/auth/bind-bank-card`
- **说明**: 绑定银行卡
- **认证**: 需要Bearer Token
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
      "invited_by": "INVITE123",
      "has_group_buy_qualification": false,
      "created_at": "2025-01-01T00:00:00Z"
    }
  },
  "timestamp": 1751365370
}
```

### 2.7 获取银行卡信息
- **接口**: `POST /api/v1/auth/bank-card`
- **说明**: 获取当前用户银行卡信息
- **认证**: 需要Bearer Token
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

## 3. 会话管理接口

### 3.1 检查登录状态
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

### 3.2 获取当前用户信息
- **接口**: `POST /api/v1/session/user`
- **说明**: 获取当前会话用户信息
- **认证**: 需要Bearer Token
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

### 3.3 用户登出
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

### 3.4 刷新会话
- **接口**: `POST /api/v1/session/refresh`
- **说明**: 刷新会话
- **认证**: 需要Bearer Token
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

## 4. 钱包相关接口

### 4.1 获取钱包信息
- **接口**: `POST /api/v1/wallet/info`
- **说明**: 获取当前用户钱包信息
- **认证**: 需要Bearer Token
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
    "balance": 1000.00,
    "frozen_amount": 0.00,
    "total_recharge": 2000.00,
    "total_withdraw": 1000.00,
    "status": 1,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  },
  "timestamp": 1751365370
}
```

### 4.2 获取资金记录
- **接口**: `POST /api/v1/wallet/transactions`
- **说明**: 获取用户资金记录
- **认证**: 需要Bearer Token
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
        "amount": 1000.00,
        "balance": 1000.00,
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

### 4.3 申请提现
- **接口**: `POST /api/v1/wallet/withdraw`
- **说明**: 申请提现
- **认证**: 需要Bearer Token
- **请求参数**:
```json
{
  "amount": 100.00,
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
    "amount": 100.00,
    "status": "pending",
    "created_at": "2025-01-01T00:00:00Z"
  },
  "timestamp": 1751365370
}
```

### 4.4 获取提现汇总
- **接口**: `POST /api/v1/wallet/withdraw-summary`
- **说明**: 获取提现汇总信息
- **认证**: 需要Bearer Token
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
    "total_withdraw": 1000.00,
    "pending_withdraw": 100.00,
    "completed_withdraw": 900.00,
    "withdraw_count": 5
  },
  "timestamp": 1751365370
}
```

### 4.5 充值申请
- **接口**: `POST /api/v1/wallet/recharge-apply`
- **说明**: 充值申请
- **认证**: 需要Bearer Token
- **请求参数**:
```json
{
  "uid": "12345678",
  "amount": 500.00,
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

### 4.6 充值确认
- **接口**: `POST /api/v1/wallet/recharge-confirm`
- **说明**: 充值确认
- **认证**: 需要Bearer Token
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

## 5. 管理员接口

### 5.1 确认提现
- **接口**: `POST /api/v1/admin/withdraw/confirm`
- **说明**: 管理员确认提现
- **认证**: 需要Bearer Token
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

### 5.2 取消提现
- **接口**: `POST /api/v1/admin/withdraw/cancel`
- **说明**: 管理员取消提现
- **认证**: 需要Bearer Token
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

## 6. 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 401 | 认证失败 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 422 | 数据验证错误 |
| 500 | 服务器错误 |
| 1001 | 数据库错误 |
| 1002 | Redis错误 |
| 1003 | 参数错误 |
| 1004 | 操作失败 |
| 1005 | 用户不存在 |
| 1006 | 用户已存在 |
| 1007 | 验证失败 |
| 1008 | 账户锁定 |
| 1009 | 注册关闭 |

## 7. 使用示例

### 7.1 用户注册和登录流程

```bash
# 1. 用户注册
curl -X POST http://localhost:9001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "INVITE123"
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

### 7.2 钱包操作流程

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

## 8. 注意事项

1. **认证**: 需要认证的接口必须在请求头中包含 `Authorization: Bearer <token>`
2. **请求格式**: 所有请求都使用JSON格式
3. **响应格式**: 所有响应都遵循统一的JSON格式
4. **错误处理**: 请根据返回的code字段判断请求是否成功
5. **分页**: 支持分页的接口使用page和page_size参数
6. **限流**: 部分接口有访问频率限制，请合理使用

## 9. 更新日志

- **v1.0.0**: 初始版本，所有业务接口改为POST请求
- 统一API前缀为 `/api/v1`
- 健康检查接口保持GET请求以便监控
- 完善错误码和响应格式 