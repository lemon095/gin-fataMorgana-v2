# gin-fataMorgana 项目接口文档

## 1. 用户认证与管理

### 1.1 用户注册
- **接口**：`POST /auth/register`
- **参数**（JSON）：
  ```json
  {
    "email": "test@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "7TRABJ"
  }
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "用户注册成功",
    "data": {
      "user": {
        "id": 1,
        "uid": "12345678",
        "username": "test_xxxxxx",
        "email": "t***@example.com",
        "phone": "",
        "status": 1,
        "created_at": "2025-06-29T18:44:49.114+08:00"
      }
    }
  }
  ```

### 1.2 用户登录
- **接口**：`POST /auth/login`
- **参数**（JSON）：
  ```json
  {
    "email": "test@example.com",
    "password": "123456"
  }
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "登录成功",
    "data": {
      "tokens": {
        "access_token": "xxx",
        "refresh_token": "xxx",
        "token_type": "Bearer",
        "expires_in": 3600
      }
    }
  }
  ```

### 1.3 获取用户信息
- **接口**：`POST /auth/profile`
- **Header**：`Authorization: Bearer <token>`
- **参数**（JSON）：
  ```json
  {}
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "user": {
        "id": 1,
        "uid": "12345678",
        "username": "test_xxxxxx",
        "email": "t***@example.com",
        "phone": "",
        "bank_card_info": "",
        "experience": 0,
        "credit_score": 100,
        "status": 1,
        "invited_by": "7TRABJ",
        "has_group_buy_qualification": false,
        "rate": 0,
        "created_at": "2025-06-29T18:44:49.114+08:00"
      }
    }
  }
  ```

### 1.4 刷新Token
- **接口**：`POST /auth/refresh`
- **参数**（JSON）：
  ```json
  {
    "refresh_token": "xxx"
  }
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "令牌刷新成功",
    "data": {
      "tokens": {
        "access_token": "xxx",
        "refresh_token": "xxx",
        "token_type": "Bearer",
        "expires_in": 3600
      }
    }
  }
  ```

### 1.5 登出
- **接口**：`POST /auth/logout`
- **Header**：`Authorization: Bearer <token>`
- **返回**：
  ```json
  {
    "code": 0,
    "message": "登出成功",
    "data": null
  }
  ```

### 1.6 获取当前会话用户信息
- **接口**：`GET /session/user`
- **Header**：`Authorization: Bearer <token>`
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "user": {
        "id": 1,
        "uid": "12345678",
        "username": "test_xxxxxx",
        "email": "t***@example.com",
        "phone": "",
        "status": 1,
        "created_at": "2025-06-29T18:44:49.114+08:00"
      }
    }
  }
  ```

### 1.7 获取用户状态
- **接口**：`GET /session/status`
- **Header**：`Authorization: Bearer <token>`
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "is_logged_in": true
    }
  }
  ```

---

## 2. 钱包与资金管理

### 2.1 获取钱包信息
- **接口**：`POST /api/v1/wallet/info`
- **Header**：`Authorization: Bearer <token>`
- **参数**（JSON）：
  ```json
  {}
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "id": 1,
      "uid": "12345678",
      "balance": 1000.00,
      "status": 1,
      "currency": "CNY",
      "last_active_at": "2025-01-01T12:00:00Z",
      "created_at": "2025-01-01T12:00:00Z",
      "updated_at": "2025-01-01T12:00:00Z"
    }
  }
  ```

### 2.2 资金流水分页查询
- **接口**：`POST /api/v1/wallet/transactions`
- **Header**：`Authorization: Bearer <token>`
- **参数**（JSON）：
  ```json
  {
    "page": 1,
    "page_size": 10
  }
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "transactions": [
        {
          "id": 1,
          "transaction_no": "TX202501011200001234",
          "uid": "12345678",
          "type": "recharge",
          "type_name": "充值",
          "amount": 1000.00,
          "balance_before": 0.00,
          "balance_after": 1000.00,
          "status": "pending",
          "status_name": "待处理",
          "description": "银行卡充值",
          "created_at": "2025-01-01T12:00:00Z"
        }
      ],
      "pagination": {
        "current_page": 1,
        "page_size": 10,
        "total": 1,
        "total_pages": 1,
        "has_next": false,
        "has_prev": false
      }
    }
  }
  ```

### 2.3 获取交易详情
- **接口**：`POST /api/v1/wallet/transaction-detail`
- **Header**：`Authorization: Bearer <token>`
- **参数**（JSON）：
  ```json
  {
    "transaction_no": "TX202501011200001234"
  }
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "id": 1,
      "transaction_no": "TX202501011200001234",
      "uid": "12345678",
      "type": "recharge",
      "type_name": "充值",
      "amount": 1000.00,
      "balance_before": 0.00,
      "balance_after": 1000.00,
      "status": "pending",
      "status_name": "待处理",
      "description": "银行卡充值",
      "remark": "",
      "related_order_no": "",
      
      "operator_uid": "system",
      "ip_address": "",
      "user_agent": "",
      "created_at": "2025-01-01T12:00:00Z"
    }
  }
  ```

### 2.4 申请提现
- **接口**：`POST /api/v1/wallet/withdraw`
- **Header**：`Authorization: Bearer <token>`
- **参数**（JSON）：
  ```json
  {
    "uid": "12345678",
    "amount": 100.00,
    "description": "提现到银行卡",
    "bank_card_info": "招商银行 6225****8888",
    "password": "123456"
  }
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "transaction_no": "TX202501011200001234",
      "amount": 100.00,
      "status": "pending",
      "message": "提现申请已提交，等待处理"
    }
  }
  ```

### 2.5 充值申请
- **接口**：`POST /api/v1/wallet/recharge`
- **Header**：`Authorization: Bearer <token>`
- **参数**（JSON）：
  ```json
  {
    "uid": "12345678",
    "amount": 1000.00,
    "description": "银行卡充值"
  }
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "充值申请已提交",
    "data": {
      "transaction_no": "TX202501011200001234"
    }
  }
  ```

---

## 3. 订单管理

### 3.1 创建订单
- **接口**：`POST /api/v1/order/create`
- **Header**：`Authorization: Bearer <token>`
- **参数**（JSON）：
  ```json
  {
    "amount": 100.00,
    "profit_amount": 20.00,
    "like_count": 5,
    "share_count": 2,
    "follow_count": 1,
    "favorite_count": 0
  }
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "订单创建成功",
    "data": {
      "order_no": "ORD202501011200001234",
      "amount": 100.00,
      "status": "pending",
      "message": "订单创建成功"
    }
  }
  ```

### 3.2 订单列表查询
- **接口**：`POST /api/v1/order/list`
- **Header**：`Authorization: Bearer <token>`
- **参数**（JSON）：
  ```json
  {
    "page": 1,
    "page_size": 10,
    "status": 1
  }
  ```
- **参数说明**：
  - `status`: 订单状态类型
    - `1`: 进行中（pending状态）
    - `2`: 已完成（success状态）
    - `3`: 全部订单
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "orders": [
        {
          "id": 1,
          "order_no": "ORD202501011200001234",
          "uid": "12345678",
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
        "total": 1,
        "total_pages": 1,
        "has_next": false,
        "has_prev": false
      }
    }
  }
  ```

### 3.3 订单详情查询
- **接口**：`POST /api/v1/order/detail`
- **Header**：`Authorization: Bearer <token>`
- **参数**（JSON）：
  ```json
  {
    "order_no": "ORD202501011200001234"
  }
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "id": 1,
      "order_no": "ORD202501011200001234",
      "uid": "12345678",
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

### 3.4 订单统计
- **接口**：`POST /api/v1/order/stats`
- **Header**：`Authorization: Bearer <token>`
- **参数**（JSON）：
  ```json
  {}
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "stats": {
        "total_orders": 10,
        "pending_orders": 3,
        "success_orders": 7,
        "failed_orders": 0,
        "total_amount": 1000.00,
        "total_profit": 200.00
      }
    }
  }
  ```

---

## 3. 银行卡管理

### 3.1 绑定银行卡
- **接口**：`POST /auth/bind-bank-card`
- **Header**：`Authorization: Bearer <token>`
- **参数**（JSON）：
  ```json
  {
    "bank_name": "招商银行",
    "card_holder": "张三",
    "card_number": "6225888888888888",
    "card_type": "借记卡"
  }
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "银行卡绑定成功",
    "data": {
      "user": { ... }
    }
  }
  ```

### 3.2 获取银行卡信息
- **接口**：`GET /auth/bank-card`
- **Header**：`Authorization: Bearer <token>`
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "bank_card_info": {
        "bank_name": "招商银行",
        "card_holder": "张三",
        "card_number": "6225888888888888",
        "card_type": "借记卡"
      }
    }
  }
  ```

---

## 4. 用户登录日志

### 4.1 获取用户登录历史
- **接口**：`GET /login-log/list?uid=12345678&page=1&size=20`
- **Header**：`Authorization: Bearer <token>`
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "logs": [
        {
          "id": 1,
          "uid": "12345678",
          "login_ip": "127.0.0.1",
          "login_time": "2025-06-29T18:44:49.114+08:00",
          "status": 1
        }
      ],
      "total": 1,
      "page": 1,
      "size": 20
    }
  }
  ```

---

## 5. 管理员/邀请码管理

### 5.1 获取邀请码列表
- **接口**：`GET /admin/invite-codes`
- **Header**：`Authorization: Bearer <token>`
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "list": [
        {
          "id": 1,
          "username": "admin",
          "my_invite_code": "7TRABJ",
          "role": 1,
          "status": 1,
          "created_at": "2025-06-29T18:44:49.114+08:00"
        }
      ]
    }
  }
  ```

---

## 6. 健康检查与系统信息

### 6.1 健康检查
- **接口**：`GET /health`
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "service": "gin-fataMorgana",
      "status": "healthy"
    }
  }
  ```

---

## 7. 任务热榜

### 7.1 获取任务热榜
- **接口**：`POST /api/v1/leaderboard/ranking`
- **Header**：`Authorization: Bearer <token>`
- **参数**（JSON）：
  ```json
  {}
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "week_start": "2025-01-06T00:00:00Z",
      "week_end": "2025-01-12T23:59:59Z",
      "my_rank": {
        "id": 1,
        "uid": "12345678",
        "username": "张*三",
        "completed_at": "2025-01-10T15:30:00Z",
        "order_count": 25,
        "total_amount": 2500.00,
        "total_profit": 500.00,
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
          "total_amount": 5000.00,
          "total_profit": 1000.00,
          "rank": 1,
          "is_rank": true
        }
      ],
      "cache_expire": "2025-01-10T16:50:00Z"
    }
  }
  ```

**数据结构说明**：
- `week_start`: 本周开始时间（周一）
- `week_end`: 本周结束时间（周日）
- `my_rank`: 当前用户排名信息
  - `rank`: 排名（999表示未上榜）
  - `is_rank`: 是否在前10名榜单上
  - `username`: 脱敏后的用户名
  - `completed_at`: 最新完成订单时间
  - `order_count`: 完成订单数量
  - `total_amount`: 总金额
  - `total_profit`: 总利润
- `top_users`: 前10名用户列表
- `cache_expire`: 缓存过期时间（5分钟缓存） 