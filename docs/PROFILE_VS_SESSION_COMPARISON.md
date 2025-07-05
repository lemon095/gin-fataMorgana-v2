# Profile接口 vs Session接口对比

## 概述

项目中提供了两个获取用户信息的接口，它们有不同的用途和返回数据：

## 1. Profile接口 (`POST /api/v1/auth/profile`)

### 功能描述
获取当前用户的**完整资料信息**，包括用户的所有基本信息和业务数据。

### 实现逻辑
```go
// AuthController.GetProfile
func (ac *AuthController) GetProfile(c *gin.Context) {
    userID := middleware.GetCurrentUser(c)
    user, err := ac.userService.GetUserByID(userID)
    // 返回完整的用户信息
    utils.Success(c, gin.H{"user": user})
}
```

### 数据来源
- **数据库查询**: 通过 `UserService.GetUserByID()` 查询数据库
- **Redis查询**: 从Redis获取用户等级进度 (`rate`)
- **数据脱敏**: 对邮箱、手机号等敏感信息进行脱敏处理

### 返回数据
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
      "created_at": "2025-01-01T00:00:00Z"
    }
  }
}
```

### 特点
- ✅ **数据完整**: 包含用户的所有业务信息
- ✅ **实时数据**: 从数据库和Redis获取最新数据
- ✅ **数据脱敏**: 自动对敏感信息进行脱敏
- ✅ **业务逻辑**: 包含等级进度等业务数据
- ⚠️ **性能开销**: 需要查询数据库和Redis

---

## 2. Session接口 (`POST /api/v1/session/user`)

### 功能描述
获取当前**会话中的用户基本信息**，主要用于验证登录状态和获取基本会话信息。

### 实现逻辑
```go
// SessionController.GetCurrentUserInfo
func (sc *SessionController) GetCurrentUserInfo(c *gin.Context) {
    userID := middleware.GetCurrentUser(c)
    username := middleware.GetCurrentUsername(c)
    // 返回会话中的基本信息
    utils.SuccessWithMessage(c, "获取用户信息成功", gin.H{
        "user_id":    userID,
        "username":   username,
        "login_time": time.Now().Unix(),
    })
}
```

### 数据来源
- **内存数据**: 从JWT token解析出的用户信息（存储在中间件中）
- **无数据库查询**: 直接从token中获取信息

### 返回数据
```json
{
  "code": 0,
  "message": "获取用户信息成功",
  "data": {
    "user_id": 1,
    "username": "test_xxxxxx",
    "login_time": 1751365370
  }
}
```

### 特点
- ✅ **性能高效**: 无需查询数据库，直接从token获取
- ✅ **响应快速**: 毫秒级响应时间
- ✅ **轻量级**: 只返回基本会话信息
- ❌ **数据有限**: 只包含token中的基本信息
- ❌ **非实时**: 数据来自token，可能不是最新的

---

## 3. 详细对比表

| 特性 | Profile接口 | Session接口 |
|------|-------------|-------------|
| **接口路径** | `POST /api/v1/auth/profile` | `POST /api/v1/session/user` |
| **数据来源** | 数据库 + Redis | JWT Token |
| **数据完整性** | 完整用户资料 | 基本会话信息 |
| **响应速度** | 较慢（需要查询） | 很快（内存获取） |
| **数据实时性** | 实时最新数据 | Token中的快照数据 |
| **包含字段** | 所有用户字段 | 仅user_id、username、login_time |
| **数据脱敏** | 自动脱敏 | 无脱敏处理 |
| **业务数据** | 包含等级进度等 | 无业务数据 |
| **使用场景** | 用户资料页面 | 登录状态检查 |

---

## 4. 使用建议

### 何时使用Profile接口
- 📱 **用户资料页面**: 显示完整的用户信息
- 💳 **个人中心**: 需要显示经验值、信用分等业务数据
- 🔧 **设置页面**: 需要编辑用户信息时
- 📊 **数据展示**: 需要显示用户等级进度等

### 何时使用Session接口
- 🔍 **登录状态检查**: 快速验证用户是否已登录
- 🚀 **页面初始化**: 需要快速获取用户基本信息
- ⚡ **性能敏感场景**: 需要毫秒级响应
- 🔄 **频繁调用**: 不需要完整用户数据的场景

---

## 5. 示例代码

### 前端使用示例

```javascript
// 获取完整用户信息（用于个人资料页面）
async function getUserProfile() {
  const response = await fetch('/api/v1/auth/profile', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({})
  });
  return response.json();
}

// 获取会话信息（用于登录状态检查）
async function getSessionInfo() {
  const response = await fetch('/api/v1/session/user', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({})
  });
  return response.json();
}
```

---

## 6. 总结

两个接口各有优势，建议根据具体使用场景选择合适的接口：

- **Profile接口**: 适用于需要完整用户信息的场景
- **Session接口**: 适用于需要快速验证登录状态的场景

这样的设计既保证了数据的完整性，又优化了性能，是一个很好的架构设计。 