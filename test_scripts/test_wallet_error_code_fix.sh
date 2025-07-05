#!/bin/bash

# 测试钱包错误码修复
echo "=== 测试钱包错误码修复 ==="

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 测试用户登录
echo "1. 用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$TOKEN" ]; then
    echo "❌ 登录失败，无法获取token"
    exit 1
fi

echo "✅ 登录成功，获取到token: ${TOKEN:0:20}..."

# 测试充值申请
echo "2. 测试充值申请..."
RECHARGE_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/recharge" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "uid": "21580100",
    "amount": 100.00,
    "description": "测试充值"
  }')

echo "充值响应: $RECHARGE_RESPONSE"

# 检查错误码
ERROR_CODE=$(echo $RECHARGE_RESPONSE | grep -o '"code":[0-9]*' | cut -d':' -f2)
ERROR_MESSAGE=$(echo $RECHARGE_RESPONSE | grep -o '"message":"[^"]*"' | cut -d'"' -f4)

echo "错误码: $ERROR_CODE"
echo "错误消息: $ERROR_MESSAGE"

# 验证错误码是否正确
if [ "$ERROR_CODE" = "3013" ]; then
    echo "✅ 错误码正确 (3013 - 钱包已被冻结，无法充值)"
elif [ "$ERROR_CODE" = "200" ]; then
    echo "✅ 充值申请成功"
else
    echo "❌ 错误码不正确: $ERROR_CODE"
fi

# 测试提现申请
echo "3. 测试提现申请..."
WITHDRAW_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/withdraw" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "amount": 50.00,
    "password": "123456"
  }')

echo "提现响应: $WITHDRAW_RESPONSE"

# 检查提现错误码
WITHDRAW_ERROR_CODE=$(echo $WITHDRAW_RESPONSE | grep -o '"code":[0-9]*' | cut -d':' -f2)
WITHDRAW_ERROR_MESSAGE=$(echo $WITHDRAW_RESPONSE | grep -o '"message":"[^"]*"' | cut -d'"' -f4)

echo "提现错误码: $WITHDRAW_ERROR_CODE"
echo "提现错误消息: $WITHDRAW_ERROR_MESSAGE"

# 验证提现错误码
if [ "$WITHDRAW_ERROR_CODE" = "3016" ]; then
    echo "✅ 提现错误码正确 (3016 - 钱包已被冻结，无法提现)"
elif [ "$WITHDRAW_ERROR_CODE" = "200" ]; then
    echo "✅ 提现申请成功"
elif [ "$WITHDRAW_ERROR_CODE" = "3017" ]; then
    echo "✅ 提现被正确拒绝 (3017 - 余额不足)"
elif [ "$WITHDRAW_ERROR_CODE" = "3018" ]; then
    echo "✅ 提现被正确拒绝 (3018 - 请先绑定银行卡)"
elif [ "$WITHDRAW_ERROR_CODE" = "3019" ]; then
    echo "✅ 提现被正确拒绝 (3019 - 登录密码错误)"
else
    echo "❌ 提现错误码不正确: $WITHDRAW_ERROR_CODE"
fi

echo "=== 测试完成 ===" 