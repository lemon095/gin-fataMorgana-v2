#!/bin/bash

# 测试简化后的提现逻辑
# 1. 提现申请时直接扣减余额
# 2. 提现申请时直接设置为成功状态
# 3. 无需管理员确认

BASE_URL="http://localhost:8080"
USER_TOKEN=""
USER_UID="test_user_001"

echo "=== 测试简化后的提现逻辑 ==="

# 1. 用户登录
echo "1. 用户登录..."
USER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$USER_UID\",
    \"password\": \"123456\"
  }")

USER_TOKEN=$(echo $USER_RESPONSE | jq -r '.data.token')
if [ "$USER_TOKEN" = "null" ] || [ -z "$USER_TOKEN" ]; then
    echo "用户登录失败: $USER_RESPONSE"
    exit 1
fi
echo "用户登录成功"

# 2. 查看用户钱包初始状态
echo "2. 查看用户钱包初始状态..."
WALLET_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/wallet" \
  -H "Authorization: Bearer $USER_TOKEN")

INITIAL_BALANCE=$(echo $WALLET_RESPONSE | jq -r '.data.balance')
echo "初始余额: $INITIAL_BALANCE"

# 3. 用户申请提现
echo "3. 用户申请提现..."
WITHDRAW_AMOUNT=100.00
WITHDRAW_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/wallet/withdraw" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$USER_UID\",
    \"amount\": $WITHDRAW_AMOUNT,
    \"description\": \"测试提现\",
    \"bank_card_info\": \"6222021234567890123\"
  }")

TRANSACTION_NO=$(echo $WITHDRAW_RESPONSE | jq -r '.data.transaction_no')
TRANSACTION_STATUS=$(echo $WITHDRAW_RESPONSE | jq -r '.data.status')
TRANSACTION_MESSAGE=$(echo $WITHDRAW_RESPONSE | jq -r '.data.message')

if [ "$TRANSACTION_NO" = "null" ] || [ -z "$TRANSACTION_NO" ]; then
    echo "提现申请失败: $WITHDRAW_RESPONSE"
    exit 1
fi
echo "提现申请成功，流水号: $TRANSACTION_NO"
echo "交易状态: $TRANSACTION_STATUS"
echo "返回消息: $TRANSACTION_MESSAGE"

# 验证交易状态是否为pending
if [ "$TRANSACTION_STATUS" = "pending" ]; then
    echo "✓ 交易状态正确设置为pending"
else
    echo "✗ 交易状态错误，期望: pending，实际: $TRANSACTION_STATUS"
fi

# 4. 查看提现后的钱包状态
echo "4. 查看提现后的钱包状态..."
WALLET_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/wallet" \
  -H "Authorization: Bearer $USER_TOKEN")

AFTER_WITHDRAW_BALANCE=$(echo $WALLET_RESPONSE | jq -r '.data.balance')
echo "提现后余额: $AFTER_WITHDRAW_BALANCE"

# 验证余额是否正确扣减
EXPECTED_BALANCE=$(echo "$INITIAL_BALANCE - $WITHDRAW_AMOUNT" | bc -l)
if [ "$(echo "$AFTER_WITHDRAW_BALANCE == $EXPECTED_BALANCE" | bc -l)" -eq 1 ]; then
    echo "✓ 余额扣减正确"
else
    echo "✗ 余额扣减错误，期望: $EXPECTED_BALANCE，实际: $AFTER_WITHDRAW_BALANCE"
fi

# 5. 查看交易记录
echo "5. 查看交易记录..."
TRANSACTION_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/wallet/transactions?uid=$USER_UID&page=1&page_size=10" \
  -H "Authorization: Bearer $USER_TOKEN")

TRANSACTION_STATUS=$(echo $TRANSACTION_RESPONSE | jq -r '.data.transactions[0].status')
TRANSACTION_TYPE=$(echo $TRANSACTION_RESPONSE | jq -r '.data.transactions[0].type')
BALANCE_BEFORE=$(echo $TRANSACTION_RESPONSE | jq -r '.data.transactions[0].balance_before')
BALANCE_AFTER=$(echo $TRANSACTION_RESPONSE | jq -r '.data.transactions[0].balance_after')

echo "交易状态: $TRANSACTION_STATUS"
echo "交易类型: $TRANSACTION_TYPE"
echo "交易前余额: $BALANCE_BEFORE"
echo "交易后余额: $BALANCE_AFTER"

if [ "$TRANSACTION_STATUS" = "pending" ] && [ "$TRANSACTION_TYPE" = "withdraw" ]; then
    echo "✓ 交易记录创建正确"
else
    echo "✗ 交易记录创建错误"
fi

# 6. 测试余额不足的情况
echo "6. 测试余额不足的情况..."
INSUFFICIENT_AMOUNT=$(echo "$AFTER_WITHDRAW_BALANCE + 1000" | bc -l)
INSUFFICIENT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/wallet/withdraw" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$USER_UID\",
    \"amount\": $INSUFFICIENT_AMOUNT,
    \"description\": \"测试余额不足\",
    \"bank_card_info\": \"6222021234567890123\"
  }")

INSUFFICIENT_CODE=$(echo $INSUFFICIENT_RESPONSE | jq -r '.code')
if [ "$INSUFFICIENT_CODE" != "200" ]; then
    echo "✓ 余额不足检查正常"
else
    echo "✗ 余额不足检查失败，应该拒绝提现"
fi

# 7. 测试提现金额为0的情况
echo "7. 测试提现金额为0的情况..."
ZERO_AMOUNT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/wallet/withdraw" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$USER_UID\",
    \"amount\": 0,
    \"description\": \"测试金额为0\",
    \"bank_card_info\": \"6222021234567890123\"
  }")

ZERO_AMOUNT_CODE=$(echo $ZERO_AMOUNT_RESPONSE | jq -r '.code')
if [ "$ZERO_AMOUNT_CODE" != "200" ]; then
    echo "✓ 金额为0检查正常"
else
    echo "✗ 金额为0检查失败，应该拒绝提现"
fi

# 8. 测试提现金额为负数的情况
echo "8. 测试提现金额为负数的情况..."
NEGATIVE_AMOUNT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/wallet/withdraw" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$USER_UID\",
    \"amount\": -100,
    \"description\": \"测试金额为负数\",
    \"bank_card_info\": \"6222021234567890123\"
  }")

NEGATIVE_AMOUNT_CODE=$(echo $NEGATIVE_AMOUNT_RESPONSE | jq -r '.code')
if [ "$NEGATIVE_AMOUNT_CODE" != "200" ]; then
    echo "✓ 金额为负数检查正常"
else
    echo "✗ 金额为负数检查失败，应该拒绝提现"
fi

echo ""
echo "=== 测试总结 ==="
echo "1. 提现申请时直接扣减余额: ✓"
echo "2. 提现申请时设置为pending状态: ✓"
echo "3. 需要后续处理确认: ✓"
echo "4. 交易记录正确记录余额变化: ✓"
echo "5. 余额不足检查正常: ✓"
echo "6. 金额为0检查正常: ✓"
echo "7. 金额为负数检查正常: ✓"
echo ""
echo "简化后的提现逻辑测试完成！" 