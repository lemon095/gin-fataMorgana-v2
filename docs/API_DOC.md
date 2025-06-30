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
    "invite_code": "RMB3IX"
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
- **接口**：`GET /auth/profile`
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
- **接口**：`GET /wallet/info`
- **Header**：`Authorization: Bearer <token>`
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "wallet": {
        "uid": "12345678",
        "balance": 1000,
        "frozen": 0,
        "status": 1,
        "updated_at": "2025-06-29T18:44:49.114+08:00"
      }
    }
  }
  ```

### 2.2 资金流水分页查询
- **接口**：`GET /wallet/transactions?page=1&size=20`
- **Header**：`Authorization: Bearer <token>`
- **返回**：
  ```json
  {
    "code": 0,
    "message": "操作成功",
    "data": {
      "transactions": [
        {
          "id": 1,
          "uid": "12345678",
          "type": "recharge",
          "amount": 1000,
          "status": 1,
          "created_at": "2025-06-29T18:44:49.114+08:00"
        }
      ],
      "total": 1,
      "page": 1,
      "size": 20
    }
  }
  ```

### 2.3 申请提现
- **接口**：`POST /wallet/withdraw`
- **Header**：`Authorization: Bearer <token>`
- **参数**（JSON）：
  ```json
  {
    "amount": 100
  }
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "提现申请成功",
    "data": null
  }
  ```

### 2.4 充值申请/确认
- **接口**：`POST /wallet/recharge-apply`、`POST /wallet/recharge-confirm`
- **Header**：`Authorization: Bearer <token>`
- **参数**（JSON）：
  ```json
  {
    "amount": 1000
  }
  ```
- **返回**：
  ```json
  {
    "code": 0,
    "message": "充值申请成功",
    "data": null
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
          "my_invite_code": "RMB3IX",
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