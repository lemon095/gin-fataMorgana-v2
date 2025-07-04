# 用户待审核功能

## 概述

为了加强用户管理，系统新增了用户待审核功能。新注册的用户默认为待审核状态，需要管理员审核通过后才能正常使用系统。

## 用户状态说明

| 状态值 | 状态名称 | 说明 |
|--------|----------|------|
| 0 | 禁用 | 用户被管理员禁用，无法登录和使用系统 |
| 1 | 正常 | 用户已通过审核，可以正常登录和使用系统 |
| 2 | 待审核 | 用户刚注册，等待管理员审核 |

## 功能特性

### 1. 注册时默认状态
- 新用户注册时，状态默认为 `2`（待审核）
- 注册成功后，用户无法立即登录
- 需要管理员将状态改为 `1`（正常）后才能登录

### 2. 登录限制
- 状态为 `0`（禁用）的用户无法登录
- 状态为 `2`（待审核）的用户无法登录
- 只有状态为 `1`（正常）的用户才能登录

### 3. 错误提示
- 待审核用户登录时，返回错误信息："账户待审核，请等待管理员审核后登录"
- 禁用用户登录时，返回错误信息："账户已被禁用，无法登录"

## 数据库变更

### 用户表状态字段更新

```sql
-- 更新状态字段注释和默认值
ALTER TABLE `users` 
MODIFY COLUMN `status` bigint NOT NULL DEFAULT 2 
COMMENT '用户状态 0:禁用 1:正常 2:待审核';
```

### 现有数据迁移

```sql
-- 将现有正常用户的状态设置为1（正常）
UPDATE users SET status = 1 WHERE status = 1 OR status IS NULL;

-- 将新注册用户的状态设置为2（待审核）
-- 这个会在新用户注册时自动设置
```

## 代码变更

### 1. 用户模型更新

```go
// User 用户模型
type User struct {
    // ... 其他字段
    Status int `json:"status" gorm:"default:2;comment:用户状态 0:禁用 1:正常 2:待审核"`
    // ... 其他字段
}
```

### 2. 用户服务更新

#### 注册逻辑
```go
// 创建新用户
user := &models.User{
    // ... 其他字段
    Status: 2, // 默认待审核
    // ... 其他字段
}
```

#### 登录逻辑
```go
// 检查用户是否被禁用
if user.Status == 0 {
    return nil, errors.New("账户已被禁用，无法登录")
}

// 检查用户是否待审核
if user.Status == 2 {
    return nil, errors.New("账户待审核，请等待管理员审核后登录")
}
```

### 3. 刷新令牌逻辑
```go
// 检查用户是否被禁用
if user.Status == 0 {
    return nil, errors.New("账户已被禁用，无法刷新令牌")
}

// 检查用户是否待审核
if user.Status == 2 {
    return nil, errors.New("账户待审核，无法刷新令牌")
}
```

## 管理员操作

### 审核用户

管理员可以通过以下SQL语句来审核用户：

```sql
-- 审核通过用户
UPDATE users SET status = 1 WHERE email = 'user@example.com';

-- 禁用用户
UPDATE users SET status = 0 WHERE email = 'user@example.com';

-- 查看待审核用户
SELECT id, uid, username, email, status, created_at 
FROM users 
WHERE status = 2 
ORDER BY created_at DESC;
```

### 批量审核

```sql
-- 批量审核通过所有待审核用户
UPDATE users SET status = 1 WHERE status = 2;

-- 批量审核通过指定时间范围内的用户
UPDATE users 
SET status = 1 
WHERE status = 2 
AND created_at >= '2024-01-01 00:00:00' 
AND created_at <= '2024-01-31 23:59:59';
```

## 测试

### 运行测试脚本

```bash
chmod +x test_scripts/test_user_approval.sh
./test_scripts/test_user_approval.sh
```

### 测试场景

1. **用户注册测试**
   - 注册新用户，验证状态为2（待审核）
   - 尝试登录，验证被拒绝

2. **审核通过测试**
   - 管理员将用户状态改为1（正常）
   - 用户登录，验证成功

3. **禁用用户测试**
   - 管理员将用户状态改为0（禁用）
   - 用户登录，验证被拒绝

## 前端集成

### 登录错误处理

前端需要处理新的错误信息：

```javascript
// 登录请求处理
async function login(email, password) {
    try {
        const response = await fetch('/api/v1/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })
        });
        
        const data = await response.json();
        
        if (data.code === 0) {
            // 登录成功
            console.log('登录成功');
        } else {
            // 处理错误
            switch (data.message) {
                case '账户待审核，请等待管理员审核后登录':
                    alert('您的账户正在审核中，请耐心等待');
                    break;
                case '账户已被禁用，无法登录':
                    alert('您的账户已被禁用，请联系管理员');
                    break;
                default:
                    alert(data.message);
            }
        }
    } catch (error) {
        console.error('登录失败:', error);
    }
}
```

### 注册成功提示

```javascript
// 注册成功后的提示
function showRegistrationSuccess() {
    alert('注册成功！您的账户正在审核中，审核通过后即可登录使用。');
}
```

## 注意事项

1. **现有用户处理**：现有正常用户的状态会自动设置为1（正常）
2. **数据一致性**：确保所有相关服务都正确处理新的状态值
3. **错误信息**：前端需要更新错误处理逻辑，支持新的错误信息
4. **管理员权限**：只有管理员可以修改用户状态
5. **日志记录**：建议记录用户状态变更的日志

## 扩展功能

### 自动审核

可以考虑添加自动审核功能：

```go
// 自动审核条件
func (s *UserService) autoApproveUser(user *models.User) bool {
    // 例如：特定邮箱域名的用户自动审核通过
    if strings.HasSuffix(user.Email, "@company.com") {
        return true
    }
    
    // 例如：特定邀请码的用户自动审核通过
    if user.InvitedBy == "AUTO_APPROVE" {
        return true
    }
    
    return false
}
```

### 审核通知

可以添加审核结果通知功能：

```go
// 发送审核结果通知
func (s *UserService) sendApprovalNotification(user *models.User, approved bool) {
    if approved {
        // 发送审核通过通知
        // 可以通过邮件、短信等方式通知用户
    } else {
        // 发送审核拒绝通知
    }
}
``` 