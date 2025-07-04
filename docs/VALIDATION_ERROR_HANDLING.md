# 参数验证错误处理优化

## 概述

为了提供更好的用户体验，我们对参数验证错误处理进行了优化，将原来技术性的错误信息转换为用户友好的中文提示。

## 主要改进

### 1. 错误信息本地化

- 将字段名转换为中文（如 `Email` → `邮箱`）
- 将验证标签转换为友好的中文提示（如 `required` → `不能为空`）
- 支持多种验证规则的错误信息

### 2. 结构化错误响应

新的错误响应格式包含：
- `code`: 错误码（422）
- `message`: 主要错误信息
- `errors`: 详细的错误列表
- `timestamp`: 时间戳

### 3. 支持的验证规则

| 验证标签 | 中文提示 | 示例 |
|---------|---------|------|
| required | 不能为空 | 邮箱不能为空 |
| email | 格式不正确 | 邮箱格式不正确 |
| min | 长度不能少于X位 | 密码长度不能少于6位 |
| max | 长度不能超过X位 | 密码长度不能超过50位 |
| len | 长度必须为X位 | 验证码长度必须为6位 |
| oneof | 必须是以下值之一：X | 类型必须是以下值之一：充值、提现 |

## 使用方式

### 在控制器中使用

```go
// 旧方式
if err := c.ShouldBindJSON(&req); err != nil {
    utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
    return
}

// 新方式
if err := c.ShouldBindJSON(&req); err != nil {
    utils.HandleValidationError(c, err)
    return
}
```

### 手动验证结构体

```go
// 验证结构体
if !utils.ValidateAndHandleError(c, &req) {
    return
}
```

## 错误响应示例

### 单个错误

```json
{
  "code": 422,
  "message": "邮箱格式不正确",
  "errors": [
    {
      "field": "邮箱",
      "tag": "email",
      "value": "invalid-email",
      "message": "邮箱格式不正确"
    }
  ],
  "timestamp": 1751623240988
}
```

### 多个错误

```json
{
  "code": 422,
  "message": "邮箱不能为空；密码长度不能少于6位",
  "errors": [
    {
      "field": "邮箱",
      "tag": "required",
      "value": "",
      "message": "邮箱不能为空"
    },
    {
      "field": "密码",
      "tag": "min",
      "value": "123",
      "message": "密码长度不能少于6位"
    }
  ],
  "timestamp": 1751623240988
}
```

## 字段名称映射

系统自动将以下字段名转换为中文：

| 英文字段名 | 中文字段名 |
|-----------|-----------|
| Email | 邮箱 |
| Password | 密码 |
| ConfirmPassword | 确认密码 |
| Username | 用户名 |
| InviteCode | 邀请码 |
| OldPassword | 当前密码 |
| NewPassword | 新密码 |
| BankName | 银行名称 |
| CardHolder | 持卡人 |
| CardNumber | 银行卡号 |
| CardType | 卡类型 |
| Page | 页码 |
| PageSize | 每页大小 |
| Type | 类型 |
| Amount | 金额 |
| PeriodNumber | 期号 |
| Uid | 用户ID |
| OrderNo | 订单号 |
| TransactionNo | 交易流水号 |
| WithdrawAmount | 提现金额 |
| RechargeAmount | 充值金额 |
| BankCardInfo | 银行卡信息 |

## 测试

运行测试脚本验证错误处理功能：

```bash
chmod +x test_scripts/test_validation_errors.sh
./test_scripts/test_validation_errors.sh
```

## 扩展

### 添加新的字段映射

在 `utils/validator.go` 中的 `fieldNameMap` 添加新的映射：

```go
var fieldNameMap = map[string]string{
    // 现有映射...
    "NewField": "新字段",
}
```

### 添加新的验证规则

在 `tagErrorMap` 中添加新的验证规则：

```go
var tagErrorMap = map[string]string{
    // 现有规则...
    "custom_rule": "自定义错误信息",
}
```

## 注意事项

1. 错误信息会自动合并，多个错误用分号分隔
2. 字段名映射是大小写敏感的
3. 如果字段名没有映射，会直接使用原字段名
4. 时间戳使用毫秒级精度

## 兼容性

新的错误处理方式完全向后兼容，不会影响现有的业务逻辑，只是改进了错误信息的展示方式。 