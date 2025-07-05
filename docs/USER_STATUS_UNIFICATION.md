# 用户状态统一规范

## 概述

本文档定义了项目中用户状态的统一规范，确保所有相关代码使用一致的状态值。

## 状态定义

### 用户状态常量
```go
const (
    UserStatusDisabled = 0 // 禁用
    UserStatusActive   = 1 // 正常
    UserStatusPending  = 2 // 待审核
)
```

### 状态说明
- **0 (UserStatusDisabled)**: 禁用 - 用户无法登录和进行任何操作
- **1 (UserStatusActive)**: 正常 - 用户可以正常使用所有功能
- **2 (UserStatusPending)**: 待审核 - 用户注册后等待管理员审核，无法登录

## 数据库字段定义

### users表status字段
```sql
`status` int DEFAULT 2 COMMENT '用户状态 0:禁用 1:正常 2:待审核'
```

### 默认值
- 新注册用户默认状态为 `2` (待审核)
- 需要管理员审核通过后状态变为 `1` (正常)
- 管理员可以设置用户状态为 `0` (禁用)

## 代码实现

### 1. 模型定义 (models/user.go)
```go
// 用户状态常量
const (
    UserStatusDisabled = 0 // 禁用
    UserStatusActive   = 1 // 正常
    UserStatusPending  = 2 // 待审核
)

type User struct {
    // ... 其他字段
    Status int `json:"status" gorm:"default:2;comment:用户状态 0:禁用 1:正常 2:待审核"`
    // ... 其他字段
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
    return u.Status == UserStatusActive
}
```

### 2. 登录状态检查 (services/user_service.go)
```go
// 新增：校验用户状态
if user.DeletedAt != nil {
    s.recordFailedLogin(ctx, req.Account, loginIP, userAgent, "账户已被删除")
    return nil, utils.NewAppError(utils.CodeUserDeletedLogin, "账户已被删除，无法登录")
}
if user.Status == UserStatusDisabled {
    s.recordFailedLogin(ctx, req.Account, loginIP, userAgent, "账户已被禁用")
    return nil, utils.NewAppError(utils.CodeUserDisabledLogin, "账户已被禁用，无法登录")
}
if user.Status == UserStatusPending {
    s.recordFailedLogin(ctx, req.Account, loginIP, userAgent, "账户待审核")
    return nil, utils.NewAppError(utils.CodeUserPendingApproval, "账户待审核，无法登录")
}
```

### 3. Token刷新状态检查
```go
// 检查用户状态
if user.Status == UserStatusDisabled {
    return nil, utils.NewAppError(utils.CodeUserDisabledRefresh, "账户已被禁用，无法刷新令牌")
}
if user.Status == UserStatusPending {
    return nil, utils.NewAppError(utils.CodeUserPendingRefresh, "账户待审核，无法刷新令牌")
}
```

### 4. 密码修改状态检查
```go
// 检查用户状态
if user.Status == UserStatusDisabled {
    return utils.NewAppError(utils.CodeUserDisabledChangePwd, "账户已被禁用，无法修改密码")
}
if user.Status == UserStatusPending {
    return utils.NewAppError(utils.CodeUserPendingApproval, "账户待审核，无法修改密码")
}
```

## 状态检查规则

### 登录阶段
- ✅ 状态为 `0` (禁用) - 无法登录
- ✅ 状态为 `2` (待审核) - 无法登录
- ✅ 状态为 `1` (正常) - 可以登录

### 操作阶段
- ✅ 状态为 `0` (禁用) - 无法进行任何操作
- ✅ 状态为 `2` (待审核) - 无法进行任何操作
- ✅ 状态为 `1` (正常) - 可以正常操作

### 需要修复的接口
以下接口需要添加用户状态检查：

1. **拼单控制器** (`controllers/group_buy_controller.go`)
   - `GetActiveGroupBuyDetail`
   - `JoinGroupBuy`

2. **排行榜控制器** (`controllers/leaderboard_controller.go`)
   - `GetLeaderboard`

3. **会话控制器** (`controllers/session_controller.go`)
   - `GetCurrentUserInfo`
   - `RefreshSession`

## 测试验证

### 测试脚本
运行 `test_scripts/test_user_status_unification.sh` 来验证状态统一性：

```bash
./test_scripts/test_user_status_unification.sh
```

### 测试用例
1. **注册用户** - 验证默认状态为待审核
2. **待审核用户登录** - 验证被正确拦截
3. **禁用用户登录** - 验证被正确拦截
4. **正常用户登录** - 验证可以正常登录

## 错误码对应

### 用户状态相关错误码
- `CodeUserDisabledLogin` - 账户已被禁用，无法登录
- `CodeUserPendingApproval` - 账户待审核，无法登录
- `CodeUserDisabledRefresh` - 账户已被禁用，无法刷新令牌
- `CodeUserPendingRefresh` - 账户待审核，无法刷新令牌
- `CodeUserDisabledChangePwd` - 账户已被禁用，无法修改密码

## 最佳实践

### 1. 使用常量
```go
// ✅ 正确 - 使用常量
if user.Status == UserStatusDisabled {
    return error
}

// ❌ 错误 - 硬编码数字
if user.Status == 0 {
    return error
}
```

### 2. 状态检查顺序
```go
// 推荐的状态检查顺序
if user.DeletedAt != nil {
    // 1. 先检查是否已删除
}
if user.Status == UserStatusDisabled {
    // 2. 再检查是否禁用
}
if user.Status == UserStatusPending {
    // 3. 最后检查是否待审核
}
```

### 3. 错误消息
- 使用统一的错误消息格式
- 错误消息应该清晰说明当前状态和操作限制

## 后续优化

### 1. 添加状态检查中间件
考虑创建一个中间件来统一处理用户状态检查，避免在每个控制器中重复代码。

### 2. 状态变更日志
记录用户状态变更的日志，便于审计和问题排查。

### 3. 状态变更通知
当用户状态变更时，可以考虑发送通知给用户。

## 总结

通过统一用户状态定义和检查逻辑，确保了：

1. **一致性** - 所有代码使用相同的状态值
2. **安全性** - 非正常状态用户无法进行敏感操作
3. **可维护性** - 使用常量避免硬编码，便于维护
4. **可扩展性** - 状态定义清晰，便于后续扩展

建议在后续开发中严格遵循此规范，确保用户状态管理的一致性和安全性。 