# API响应码说明

## 概述

所有API接口都使用统一的响应格式：

```json
{
  "code": 0,        // 响应码，0表示成功，其他值表示错误
  "message": "",    // 响应消息
  "data": {}        // 响应数据
}
```

## 响应码分类

### 成功响应码
- `0` - 操作成功

### 客户端错误码 (1000-1999)
- `1000` - 参数错误
- `1001` - 验证失败
- `1002` - 未授权
- `1003` - 禁止访问
- `1004` - 资源不存在
- `1005` - 方法不允许
- `1006` - 请求超时
- `1007` - 请求过于频繁
- `1008` - 资源冲突
- `1009` - 无法处理的实体

### 认证相关错误码 (2000-2099)
- `2000` - Token过期
- `2001` - Token无效
- `2002` - Token缺失
- `2003` - 登录失败
- `2004` - 用户不存在
- `2005` - 密码错误
- `2006` - 用户已存在
- `2007` - 邮箱已存在
- `2008` - 邀请码无效
- `2009` - 账户被锁定
- `2010` - 会话过期

### 业务逻辑错误码 (3000-3999)
- `3000` - 操作失败
- `3001` - 资源繁忙
- `3002` - 余额不足
- `3003` - 超出限制
- `3004` - 无效操作
- `3005` - 数据不一致
- `3006` - 违反业务规则

### 服务器错误码 (5000-5999)
- `5000` - 内部服务器错误
- `5001` - 数据库错误
- `5002` - Redis错误
- `5003` - 外部API错误
- `5004` - 服务不可用
- `5005` - 网关超时
- `5006` - 配置错误

## 使用示例

### 成功响应示例

```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "user": {
      "id": 1,
      "uid": "12345678",
      "username": "user_abc123",
      "email": "user@example.com"
    }
  }
}
```

### 错误响应示例

```json
{
  "code": 2007,
  "message": "邮箱已被注册",
  "data": null
}
```

### 带自定义消息的成功响应

```json
{
  "code": 0,
  "message": "用户注册成功",
  "data": {
    "user": {
      "id": 1,
      "uid": "12345678",
      "username": "user_abc123",
      "email": "user@example.com"
    }
  }
}
```

### 带数据的错误响应

```json
{
  "code": 1001,
  "message": "数据验证失败",
  "data": {
    "field": "email",
    "error": "邮箱格式不正确"
  }
}
```

## HTTP状态码映射

系统会根据业务错误码自动设置相应的HTTP状态码：

- `0` (成功) → `200 OK`
- `1000-1999` (客户端错误) → `400 Bad Request`
- `2000-2099` (认证错误) → `401 Unauthorized`
- `3000-3999` (业务错误) → `422 Unprocessable Entity`
- `5000-5999` (服务器错误) → `500 Internal Server Error`

## 常用响应函数

### 成功响应
```go
// 标准成功响应
utils.Success(c, data)

// 带自定义消息的成功响应
utils.SuccessWithMessage(c, "操作成功", data)
```

### 错误响应
```go
// 标准错误响应
utils.Error(c, code)

// 带自定义消息的错误响应
utils.ErrorWithMessage(c, code, "自定义错误消息")

// 带数据的错误响应
utils.ErrorWithData(c, code, errorData)
```

### 常用错误响应函数
```go
// 参数错误
utils.InvalidParams(c)
utils.InvalidParamsWithMessage(c, "具体错误信息")

// 认证相关
utils.Unauthorized(c)
utils.TokenExpired(c)
utils.TokenInvalid(c)
utils.LoginFailed(c)
utils.UserNotFound(c)
utils.UserAlreadyExists(c)
utils.EmailAlreadyExists(c)

// 服务器错误
utils.InternalError(c)
utils.DatabaseError(c)
utils.RedisError(c)
```

## 最佳实践

1. **统一格式**: 所有API接口都使用相同的响应格式
2. **明确错误码**: 使用具体的错误码而不是通用错误码
3. **友好消息**: 提供用户友好的错误消息
4. **详细数据**: 在data字段中提供详细的响应数据
5. **状态码一致**: 确保HTTP状态码与业务错误码一致

## 扩展错误码

如需添加新的错误码，请在 `utils/response.go` 文件中：

1. 在常量定义区域添加新的错误码
2. 在 `ResponseMessage` 映射中添加对应的消息
3. 根据需要添加便捷的响应函数

```go
// 添加新的错误码
const (
    CodeNewError = 4000 // 新错误类型
)

// 添加错误消息
var ResponseMessage = map[int]string{
    // ... 现有消息
    CodeNewError: "新错误消息",
}

// 添加便捷函数
func NewError(c *gin.Context) {
    Error(c, CodeNewError)
}
``` 