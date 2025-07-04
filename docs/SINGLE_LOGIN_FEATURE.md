# 单点登录功能实现文档

## 功能概述

本项目已实现单点登录（Single Sign-On）功能，确保同一用户只能在一个设备上保持登录状态。当用户在新设备登录时，会自动踢出旧设备的登录会话。

## 核心特性

### 1. 自动踢出机制
- 新设备登录时，自动使旧token失效
- 旧设备继续使用token时，返回401错误
- 支持设备信息记录（IP、User-Agent等）

### 2. Token黑名单
- 使用Redis存储已撤销的token
- 黑名单过期时间比token过期时间稍长
- 自动清理过期黑名单记录

### 3. 活跃会话管理
- 每个用户只维护一个活跃token
- 新登录时自动更新活跃token
- 支持会话信息查询

## 技术实现

### 1. 架构设计

```
用户登录 → 验证凭据 → 检查旧会话 → 撤销旧token → 生成新token → 设置活跃token
    ↓
API请求 → 验证token → 检查黑名单 → 检查活跃token → 返回结果
```

### 2. Redis存储结构

```redis
# 用户活跃token
user:active_token:{uid} -> {
  "token_hash": "sha256哈希",
  "login_time": "登录时间",
  "device_info": "设备信息",
  "login_ip": "登录IP",
  "user_agent": "User-Agent"
}

# Token黑名单
token:blacklist:{token_hash} -> 加入黑名单的时间戳
```

### 3. 核心组件

#### TokenService
- `SetUserActiveToken()`: 设置用户活跃token
- `AddTokenToBlacklist()`: 将token加入黑名单
- `IsTokenBlacklisted()`: 检查token是否在黑名单中
- `IsActiveToken()`: 检查是否为当前活跃token
- `ValidateTokenWithBlacklist()`: 完整的token验证流程

#### 认证中间件
- 使用`ValidateTokenWithBlacklist()`进行token验证
- 返回详细的错误信息
- 支持不同的错误类型处理

## 使用流程

### 1. 用户登录流程

```go
// 1. 验证用户凭据
user := validateCredentials(req)

// 2. 检查并撤销旧会话
tokenService := NewTokenService()
activeToken, err := tokenService.GetUserActiveToken(ctx, user.Uid)
if err == nil && activeToken != nil {
    tokenService.AddTokenToBlacklist(ctx, activeToken.TokenHash)
}

// 3. 生成新token
accessToken := utils.GenerateAccessToken(user.ID, user.Uid, user.Username)

// 4. 设置新的活跃token
tokenService.SetUserActiveToken(ctx, user.Uid, accessToken, deviceInfo, loginIP, userAgent)
```

### 2. API请求验证流程

```go
// 1. 解析Authorization头
tokenString := extractTokenFromHeader(c)

// 2. 验证token（包含黑名单检查）
claims, err := tokenService.ValidateTokenWithBlacklist(ctx, tokenString)
if err != nil {
    return errorResponse(err)
}

// 3. 设置用户信息到上下文
c.Set("user_id", claims.UserID)
c.Set("uid", claims.Uid)
c.Set("username", claims.Username)
```

### 3. 用户登出流程

```go
// 1. 获取当前token
tokenString := getCurrentToken(c)

// 2. 将token加入黑名单
tokenService.AddTokenToBlacklist(ctx, tokenString)

// 3. 撤销用户会话
tokenService.RevokeUserSession(ctx, uid)
```

## 错误处理

### 1. 错误类型

| 错误类型 | 错误码 | 错误信息 | 说明 |
|---------|--------|----------|------|
| TOKEN_REVOKED | 401 | "您的账号已在其他设备登录，请重新登录" | token被踢出 |
| TOKEN_EXPIRED | 401 | "令牌已过期，请重新登录" | token自然过期 |
| INVALID_TOKEN | 401 | "无效的认证令牌" | token格式错误 |
| MISSING_TOKEN | 401 | "缺少认证令牌" | 未提供token |

### 2. 错误响应示例

```json
{
  "code": 401,
  "message": "您的账号已在其他设备登录，请重新登录",
  "error": "TOKEN_REVOKED",
  "timestamp": 1751622617
}
```

## 配置说明

### 1. Token过期时间

在`config.yaml`中配置：

```yaml
jwt:
  access_token_expire: 86400  # 访问令牌有效期（秒）
  refresh_token_expire: 604800  # 刷新令牌有效期（秒）
```

### 2. 黑名单过期时间

黑名单过期时间自动设置为：`access_token_expire + 300秒`

## 测试验证

### 1. 测试脚本

使用`test_scripts/test_single_login.sh`进行功能测试：

```bash
./test_scripts/test_single_login.sh
```

### 2. 测试场景

1. **设备A登录** → 获得token，可以正常访问
2. **设备B登录** → 设备A的token被踢出，设备B获得新token
3. **设备A使用旧token** → 返回401错误，提示"已在其他设备登录"
4. **设备B使用新token** → 可以正常访问
5. **设备A重新登录** → 设备B的token被踢出，设备A获得新token

### 3. 手动测试

```bash
# 设备A登录
curl -X POST "http://localhost:9001/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'

# 设备B登录
curl -X POST "http://localhost:9001/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'

# 设备A使用旧token访问（应该被拒绝）
curl -X POST "http://localhost:9001/api/v1/auth/profile" \
  -H "Authorization: Bearer OLD_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}'
```

## 性能考虑

### 1. Redis性能
- 使用Redis存储token信息，响应快速
- 自动过期机制，无需手动清理
- 支持集群部署

### 2. 内存使用
- token哈希存储，节省内存
- 定期清理过期数据
- 合理的过期时间设置

### 3. 网络开销
- 每次API请求需要2次Redis查询
- 可以考虑本地缓存优化
- 支持Redis连接池

## 安全考虑

### 1. Token安全
- 使用SHA256哈希存储token
- 不直接存储明文token
- 支持token撤销机制

### 2. 会话安全
- 单点登录防止会话劫持
- 设备信息记录便于审计
- 支持强制登出功能

### 3. 错误信息
- 不泄露敏感信息
- 统一的错误处理
- 详细的日志记录

## 扩展功能

### 1. 设备管理
- 支持查看当前活跃设备
- 支持手动踢出指定设备
- 支持设备白名单

### 2. 会话统计
- 统计用户登录历史
- 分析登录设备分布
- 异常登录检测

### 3. 多因素认证
- 支持短信验证码
- 支持邮箱验证码
- 支持硬件密钥

## 故障排查

### 1. 常见问题

**Q: 用户反映频繁被踢出登录**
A: 检查是否有多个客户端同时使用，或者网络环境导致重复登录

**Q: Redis连接失败**
A: 检查Redis服务状态，确保连接配置正确

**Q: Token验证失败**
A: 检查JWT密钥配置，确保token格式正确

### 2. 日志分析

```bash
# 查看应用日志
tail -f logs/app.log | grep "token"

# 查看Redis日志
redis-cli monitor | grep "token"
```

### 3. 监控指标

- 登录成功率
- Token验证失败率
- Redis连接状态
- 黑名单大小

## 总结

单点登录功能已成功实现，具备以下特点：

1. **安全性高**：确保同一用户只能在一个设备登录
2. **性能好**：基于Redis，响应快速
3. **易维护**：自动过期，无需手动清理
4. **可扩展**：支持设备信息记录和会话管理
5. **用户体验好**：明确的错误提示和操作指引

该功能已通过完整测试，可以投入生产使用。 