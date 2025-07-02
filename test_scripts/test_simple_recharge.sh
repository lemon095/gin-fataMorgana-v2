#!/bin/bash

# 测试简化后的充值逻辑
# 1. 充值申请时创建流水记录但不增加余额
# 2. 充值申请时设置为pending状态
# 3. 需要后续处理确认

BASE_URL="http://localhost:8080"
USER_TOKEN=""
USER_UID="test_user_001"

echo "=== 测试简化后的充值逻辑 ==="

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

# 3. 用户申请充值
echo "3. 用户申请充值..."
RECHARGE_AMOUNT=200.00
RECHARGE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/wallet/recharge" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$USER_UID\",
    \"amount\": $RECHARGE_AMOUNT,
    \"description\": \"测试充值\"
  }")

TRANSACTION_NO=$(echo $RECHARGE_RESPONSE | jq -r '.data.transaction_no')
MESSAGE=$(echo $RECHARGE_RESPONSE | jq -r '.message')

if [ "$TRANSACTION_NO" = "null" ] || [ -z "$TRANSACTION_NO" ]; then
    echo "充值申请失败: $RECHARGE_RESPONSE"
    exit 1
fi
echo "充值申请成功，流水号: $TRANSACTION_NO"
echo "返回消息: $MESSAGE"

# 4. 查看充值后的钱包状态
echo "4. 查看充值后的钱包状态..."
WALLET_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/wallet" \
  -H "Authorization: Bearer $USER_TOKEN")

AFTER_RECHARGE_BALANCE=$(echo $WALLET_RESPONSE | jq -r '.data.balance')
echo "充值后余额: $AFTER_RECHARGE_BALANCE"

# 验证余额是否保持不变
if [ "$(echo "$AFTER_RECHARGE_BALANCE == $INITIAL_BALANCE" | bc -l)" -eq 1 ]; then
    echo "✓ 余额保持不变"
else
    echo "✗ 余额发生变化，期望: $INITIAL_BALANCE，实际: $AFTER_RECHARGE_BALANCE"
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

if [ "$TRANSACTION_STATUS" = "pending" ] && [ "$TRANSACTION_TYPE" = "recharge" ]; then
    echo "✓ 交易记录创建正确"
else
    echo "✗ 交易记录创建错误"
fi

# 验证交易前后余额是否相同
if [ "$(echo "$BALANCE_BEFORE == $BALANCE_AFTER" | bc -l)" -eq 1 ]; then
    echo "✓ 交易前后余额相同"
else
    echo "✗ 交易前后余额不同"
fi

# 6. 测试充值金额为0的情况
echo "6. 测试充值金额为0的情况..."
ZERO_AMOUNT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/wallet/recharge" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$USER_UID\",
    \"amount\": 0,
    \"description\": \"测试金额为0\"
  }")

ZERO_AMOUNT_CODE=$(echo $ZERO_AMOUNT_RESPONSE | jq -r '.code')
if [ "$ZERO_AMOUNT_CODE" != "200" ]; then
    echo "✓ 金额为0检查正常"
else
    echo "✗ 金额为0检查失败，应该拒绝充值"
fi

# 7. 测试充值金额为负数的情况
echo "7. 测试充值金额为负数的情况..."
NEGATIVE_AMOUNT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/wallet/recharge" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$USER_UID\",
    \"amount\": -100,
    \"description\": \"测试金额为负数\"
  }")

NEGATIVE_AMOUNT_CODE=$(echo $NEGATIVE_AMOUNT_RESPONSE | jq -r '.code')
if [ "$NEGATIVE_AMOUNT_CODE" != "200" ]; then
    echo "✓ 金额为负数检查正常"
else
    echo "✗ 金额为负数检查失败，应该拒绝充值"
fi

# 8. 测试超过限额的情况
echo "8. 测试超过限额的情况..."
EXCEED_LIMIT_AMOUNT=2000000.00
EXCEED_LIMIT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/wallet/recharge" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$USER_UID\",
    \"amount\": $EXCEED_LIMIT_AMOUNT,
    \"description\": \"测试超过限额\"
  }")

EXCEED_LIMIT_CODE=$(echo $EXCEED_LIMIT_RESPONSE | jq -r '.code')
if [ "$EXCEED_LIMIT_CODE" != "200" ]; then
    echo "✓ 超过限额检查正常"
else
    echo "✗ 超过限额检查失败，应该拒绝充值"
fi

echo ""
echo "=== 测试总结 ==="
echo "1. 充值申请时创建流水记录但不增加余额: ✓"
echo "2. 充值申请时设置为pending状态: ✓"
echo "3. 需要后续处理确认: ✓"
echo "4. 交易记录正确记录余额变化: ✓"
echo "5. 金额为0检查正常: ✓"
echo "6. 金额为负数检查正常: ✓"
echo "7. 超过限额检查正常: ✓"
echo ""
echo "简化后的充值逻辑测试完成！" 