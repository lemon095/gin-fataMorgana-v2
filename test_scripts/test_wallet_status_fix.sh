#!/bin/bash

# 测试钱包状态检查修复
echo "=== 测试钱包状态检查修复 ==="

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

# 测试获取钱包信息
echo "2. 获取钱包信息..."
WALLET_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/info" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "钱包信息响应: $WALLET_RESPONSE"

# 检查钱包状态
WALLET_STATUS=$(echo $WALLET_RESPONSE | grep -o '"status":[0-9]*' | cut -d':' -f2)
echo "钱包状态: $WALLET_STATUS"

if [ "$WALLET_STATUS" = "1" ]; then
    echo "✅ 钱包状态正常 (1)"
else
    echo "❌ 钱包状态异常: $WALLET_STATUS"
fi

# 测试充值申请
echo "3. 测试充值申请..."
RECHARGE_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/recharge" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "amount": 100.00,
    "description": "测试充值"
  }')

echo "充值响应: $RECHARGE_RESPONSE"

# 检查充值是否成功
if echo $RECHARGE_RESPONSE | grep -q '"code":200'; then
    echo "✅ 充值申请成功"
else
    echo "❌ 充值申请失败"
    echo "错误信息: $RECHARGE_RESPONSE"
fi

# 测试提现申请
echo "4. 测试提现申请..."
WITHDRAW_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/withdraw" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "amount": 50.00,
    "password": "123456"
  }')

echo "提现响应: $WITHDRAW_RESPONSE"

# 检查提现是否成功（可能会因为余额不足或未绑定银行卡而失败，这是正常的）
if echo $WITHDRAW_RESPONSE | grep -q '"code":200'; then
    echo "✅ 提现申请成功"
elif echo $WITHDRAW_RESPONSE | grep -q "余额不足\|银行卡\|密码错误"; then
    echo "✅ 提现申请被正确拒绝（余额不足/未绑定银行卡/密码错误）"
else
    echo "❌ 提现申请出现异常错误"
    echo "错误信息: $WITHDRAW_RESPONSE"
fi

echo "=== 测试完成 ===" 