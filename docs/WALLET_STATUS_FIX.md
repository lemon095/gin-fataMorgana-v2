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
./test_scripts/test_wallet_status_fix.sh
```

测试内容包括：
- 获取钱包信息，确认状态
- 测试充值申请
- 测试提现申请
- 验证正常状态钱包可以正常操作

## 影响范围

### 修复前的影响
- 所有正常状态的钱包无法充值
- 所有正常状态的钱包无法提现
- 冻结状态的钱包反而可以操作（安全风险）

### 修复后的效果
- 正常状态的钱包可以正常充值
- 正常状态的钱包可以正常提现
- 冻结状态的钱包被正确阻止操作

## 相关文件

- `services/wallet_service.go` - 主要修复文件
- `models/wallet.go` - 钱包状态定义
- `services/order_service.go` - 订单服务中的正确实现
- `test_scripts/test_wallet_status_fix.sh` - 测试脚本

## 注意事项

1. 此修复只影响钱包状态检查逻辑，不涉及数据库结构变更
2. 修复后需要重启服务才能生效
3. 建议在生产环境部署前进行充分测试
4. 检查数据库中是否有钱包状态异常的数据需要修正 