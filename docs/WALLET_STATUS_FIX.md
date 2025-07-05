# 钱包状态检查问题修复

## 问题描述

用户反馈钱包被冻结无法充值，经过代码检查发现钱包状态检查逻辑存在错误。

## 问题原因

在 `services/wallet_service.go` 文件中，充值(`Recharge`)和提现(`RequestWithdraw`)方法中的钱包状态检查逻辑写反了：

### 原始错误代码

```go
// 检查钱包状态
if wallet.Status == 1 { // 已冻结
    return "", utils.NewAppError(utils.CodeWalletFrozenRecharge, "钱包已被冻结，无法充值")
}
```

### 钱包状态定义

根据 `models/wallet.go` 中的定义：
- `Status = 1`: 正常状态
- `Status = 0`: 冻结状态

### 问题分析

原始代码中 `wallet.Status == 1` 表示检查钱包状态是否为1（正常），但注释写的是"已冻结"，逻辑完全相反。这导致：
- 正常状态的钱包（Status=1）被误判为冻结
- 冻结状态的钱包（Status=0）反而可以通过检查

## 修复方案

### 修复后的代码

```go
// 检查钱包状态
if wallet.Status == 0 { // 已冻结
    return "", utils.NewAppError(utils.CodeWalletFrozenRecharge, "钱包已被冻结，无法充值")
}
```

### 修复位置

1. **充值方法** (`services/wallet_service.go:244`)
   - 修复前：`if wallet.Status == 1`
   - 修复后：`if wallet.Status == 0`

2. **提现方法** (`services/wallet_service.go:340`)
   - 修复前：`if wallet.Status == 1`
   - 修复后：`if wallet.Status == 0`

## 错误码修复

### 问题描述

修复状态检查逻辑后，发现控制器中的错误处理有问题。当服务层返回AppError时，控制器使用了固定的错误码1001（数据库错误），导致错误码和错误消息不匹配。

### 原始错误代码

```go
transactionNo, err := wc.walletService.Recharge(req.Uid, req.Amount, req.Description)
if err != nil {
    utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error()) // 错误码1001
    return
}
```

### 修复后的代码

```go
transactionNo, err := wc.walletService.Recharge(req.Uid, req.Amount, req.Description)
if err != nil {
    // 检查是否是AppError类型
    if appErr, ok := err.(*utils.AppError); ok {
        utils.ErrorWithMessage(c, appErr.Code, appErr.Message) // 使用服务层的错误码
    } else {
        utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
    }
    return
}
```

### 错误码对应关系

- **3013**: 钱包已被冻结，无法充值
- **3016**: 钱包已被冻结，无法提现
- **3012**: 钱包已被冻结，无法创建订单

## 验证方法

### 1. 代码逻辑验证

订单服务中的钱包状态检查是正确的：
```go
// 检查钱包状态
if !wallet.IsActive() {
    return nil, utils.NewAppError(utils.CodeWalletFrozenOrder, "钱包已被冻结，无法创建订单")
}
```

其中 `IsActive()` 方法：
```go
func (w *Wallet) IsActive() bool {
    return w.Status == 1
}
```

### 2. 测试验证

运行测试脚本验证修复效果：
```bash
# 测试状态检查修复
./test_scripts/test_wallet_status_fix.sh

# 测试错误码修复
./test_scripts/test_wallet_error_code_fix.sh
```

测试内容包括：
- 获取钱包信息，确认状态
- 测试充值申请，验证错误码3013
- 测试提现申请，验证错误码3016
- 验证正常状态钱包可以正常操作

## 影响范围

### 修复前的影响
- 所有正常状态的钱包无法充值
- 所有正常状态的钱包无法提现
- 冻结状态的钱包反而可以操作（安全风险）
- 错误码显示为1001，但消息是钱包冻结相关

### 修复后的效果
- 正常状态的钱包可以正常充值
- 正常状态的钱包可以正常提现
- 冻结状态的钱包被正确阻止操作
- 错误码正确显示为3013/3016，与错误消息匹配

## 相关文件

- `services/wallet_service.go` - 主要修复文件
- `controllers/wallet_controller.go` - 控制器错误处理修复
- `models/wallet.go` - 钱包状态定义
- `services/order_service.go` - 订单服务中的正确实现
- `utils/apperror.go` - 应用错误类型定义
- `test_scripts/test_wallet_status_fix.sh` - 状态检查测试脚本
- `test_scripts/test_wallet_error_code_fix.sh` - 错误码测试脚本

## 注意事项

1. 此修复只影响钱包状态检查逻辑和错误处理，不涉及数据库结构变更
2. 修复后需要重启服务才能生效
3. 建议在生产环境部署前进行充分测试
4. 检查数据库中是否有钱包状态异常的数据需要修正
5. 确保所有相关接口的错误码都正确传递 