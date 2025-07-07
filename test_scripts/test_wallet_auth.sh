#!/bin/bash

# 测试钱包接口认证问题
echo "=== 钱包接口认证测试 ==="

# 设置基础URL
BASE_URL="http://localhost:9001/api"

# 1. 先登录获取token
echo "1. 用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test@example.com",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$TOKEN" ]; then
    echo "❌ 登录失败，无法获取token"
    exit 1
fi

echo "✅ 登录成功，获取到token: ${TOKEN:0:20}..."

# 2. 测试钱包信息接口
echo ""
echo "2. 测试钱包信息接口..."
WALLET_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/wallet/info" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "钱包信息响应: $WALLET_RESPONSE"

# 检查响应
if echo "$WALLET_RESPONSE" | grep -q '"code":0'; then
    echo "✅ 钱包信息接口调用成功"
else
    echo "❌ 钱包信息接口调用失败"
    echo "错误详情: $WALLET_RESPONSE"
fi

# 3. 测试不带token的情况
echo ""
echo "3. 测试不带token的情况..."
NO_TOKEN_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/wallet/info" \
  -H "Content-Type: application/json" \
  -d '{}')

echo "无token响应: $NO_TOKEN_RESPONSE"

# 4. 测试无效token的情况
echo ""
echo "4. 测试无效token的情况..."
INVALID_TOKEN_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/wallet/info" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer invalid_token_123" \
  -d '{}')

echo "无效token响应: $INVALID_TOKEN_RESPONSE"

echo ""
echo "✅ 测试完成" 