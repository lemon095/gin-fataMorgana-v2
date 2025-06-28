# 用户状态管理功能说明

## 概述

本系统实现了完整的用户状态管理功能，包括用户注册、登录、禁用、删除等状态检查。

## 用户状态定义

### 用户状态字段
```go
type User struct {
    Status    int        `json:"status" gorm:"default:1;comment:用户状态 1:正常 0:禁用"`
    DeletedAt *time.Time `json:"-" gorm:"index;comment:软删除时间"`
}
```

### 状态说明
- **Status = 1**: 用户正常，可以登录
- **Status = 0**: 用户被禁用，无法登录
- **DeletedAt != nil**: 用户已被软删除，无法登录和注册

## 功能实现

### 1. 用户注册逻辑

#### 邮箱重复检查
- 检查邮箱是否已存在（包括已删除的用户）
- 如果邮箱存在且用户已被删除，返回"该邮箱已被删除，无法重新注册"
- 如果邮箱存在且用户未被删除，返回"邮箱已被注册"

```go
// 检查邮箱是否已存在（包括已删除的用户）
emailExists, err := s.userRepo.EmailExistsIncludeDeleted(ctx, req.Email)
if emailExists {
    // 检查是否是被删除的用户
    isDeleted, err := s.userRepo.IsUserDeleted(ctx, req.Email)
    if isDeleted {
        return nil, errors.New("该邮箱已被删除，无法重新注册")
    } else {
        return nil, errors.New("邮箱已被注册")
    }
}
```

### 2. 用户登录逻辑

#### 登录状态检查
1. **用户存在性检查**: 检查用户是否存在（包括已删除的）
2. **删除状态检查**: 如果用户已被删除，返回"账户已被删除，无法登录"
3. **禁用状态检查**: 如果用户被禁用，返回"账户已被禁用，无法登录"
4. **密码验证**: 验证密码是否正确

```go
// 首先检查用户是否存在（包括已删除的）
user, err := s.userRepo.FindByEmailIncludeDeleted(ctx, req.Email)

// 检查用户是否已被删除
if user.DeletedAt != nil {
    return nil, errors.New("账户已被删除，无法登录")
}

// 检查用户是否被禁用
if user.Status == 0 {
    return nil, errors.New("账户已被禁用，无法登录")
}
```

### 3. Token刷新逻辑

#### 刷新状态检查
- 验证刷新令牌的有效性
- 检查用户是否存在
- 检查用户是否已被删除
- 检查用户是否被禁用

```go
// 查找用户
user, err := s.userRepo.FindByUsername(ctx, claims.Username)

// 检查用户是否已被删除
if user.DeletedAt != nil {
    return nil, errors.New("账户已被删除，无法刷新令牌")
}

// 检查用户是否被禁用
if user.Status == 0 {
    return nil, errors.New("账户已被禁用，无法刷新令牌")
}
```

### 4. 用户信息获取

#### 信息获取状态检查
- 检查用户是否存在
- 检查用户是否已被删除

```go
// 检查用户是否已被删除
if user.DeletedAt != nil {
    return nil, errors.New("用户已被删除")
}
```

## 错误码映射

### 注册相关错误
- `2007` - 邮箱已被注册
- `2006` - 该邮箱已被删除，无法重新注册
- `1001` - 两次输入的密码不一致

### 登录相关错误
- `2003` - 邮箱或密码错误
- `2004` - 账户已被删除，无法登录
- `2009` - 账户已被禁用，无法登录

### Token刷新相关错误
- `2001` - 无效的刷新令牌
- `2004` - 账户已被删除，无法刷新令牌
- `2009` - 账户已被禁用，无法刷新令牌

## 数据库操作

### 用户禁用
```sql
-- 禁用用户
UPDATE users SET status = 0 WHERE email = 'user@example.com';

-- 启用用户
UPDATE users SET status = 1 WHERE email = 'user@example.com';
```

### 用户删除（软删除）
```sql
-- 删除用户（软删除）
DELETE FROM users WHERE email = 'user@example.com';

-- 恢复用户
UPDATE users SET deleted_at = NULL WHERE email = 'user@example.com';
```

### 查询用户状态
```sql
-- 查询所有用户（包括已删除的）
SELECT * FROM users;

-- 查询正常用户
SELECT * FROM users WHERE status = 1 AND deleted_at IS NULL;

-- 查询禁用用户
SELECT * FROM users WHERE status = 0 AND deleted_at IS NULL;

-- 查询已删除用户
SELECT * FROM users WHERE deleted_at IS NOT NULL;
```

## API响应示例

### 正常注册响应
```json
{
  "code": 0,
  "message": "用户注册成功",
  "data": {
    "user": {
      "id": 1,
      "uid": "12345678",
      "username": "user_abc123",
      "email": "user@example.com",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 邮箱已存在响应
```json
{
  "code": 2007,
  "message": "邮箱已被注册",
  "data": null
}
```

### 邮箱已被删除响应
```json
{
  "code": 2006,
  "message": "该邮箱已被删除，无法重新注册",
  "data": null
}
```

### 账户被禁用响应
```json
{
  "code": 2009,
  "message": "账户已被禁用，无法登录",
  "data": null
}
```

### 账户已删除响应
```json
{
  "code": 2004,
  "message": "账户已被删除，无法登录",
  "data": null
}
```

## 测试方法

### 1. 使用测试脚本
```bash
# 运行用户状态测试
./test_user_status.sh
```

### 2. 手动测试
```bash
# 1. 注册用户
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "INVITE123"
  }'

# 2. 登录用户
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456"
  }'

# 3. 在数据库中禁用用户后再次登录
# 4. 在数据库中删除用户后再次登录
```

### 3. 数据库测试
```sql
-- 1. 查看用户状态
SELECT id, email, status, deleted_at FROM users WHERE email = 'test@example.com';

-- 2. 禁用用户
UPDATE users SET status = 0 WHERE email = 'test@example.com';

-- 3. 删除用户
DELETE FROM users WHERE email = 'test@example.com';

-- 4. 恢复用户
UPDATE users SET deleted_at = NULL WHERE email = 'test@example.com';
```

## 最佳实践

### 1. 状态管理
- 使用软删除而不是硬删除，保留数据完整性
- 提供用户状态恢复功能
- 记录状态变更日志

### 2. 安全考虑
- 禁用用户时，立即使其当前会话失效
- 删除用户时，清理相关数据
- 提供管理员操作日志

### 3. 用户体验
- 提供清晰的错误提示信息
- 支持用户申诉和恢复流程
- 提供状态变更通知

### 4. 性能优化
- 为状态字段创建索引
- 使用复合索引优化查询
- 定期清理过期数据

## 扩展功能建议

### 1. 用户状态历史
```sql
CREATE TABLE user_status_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    old_status INT,
    new_status INT,
    reason VARCHAR(255),
    operator_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 2. 批量操作
```go
// 批量禁用用户
func (s *UserService) BatchDisableUsers(userIDs []uint, reason string) error

// 批量删除用户
func (s *UserService) BatchDeleteUsers(userIDs []uint, reason string) error
```

### 3. 状态变更通知
```go
// 发送状态变更通知
func (s *UserService) SendStatusChangeNotification(userID uint, newStatus int, reason string) error
```

## 总结

用户状态管理功能提供了完整的用户生命周期管理，包括：

1. **注册控制**: 防止重复注册和已删除用户重新注册
2. **登录控制**: 阻止禁用和删除用户登录
3. **Token控制**: 确保Token刷新时也进行状态检查
4. **信息保护**: 防止获取已删除用户信息
5. **错误处理**: 提供清晰的错误信息和状态码

这些功能确保了系统的安全性和数据完整性，同时提供了良好的用户体验。 