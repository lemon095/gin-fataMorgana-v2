#!/bin/bash

# 简化的Token测试脚本
echo "=== Token修复测试 ==="

# 设置基础URL
BASE_URL="http://localhost:9001/api"

# 测试用户登录
echo "1. 用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
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

# 立即测试token有效性
echo "2. 立即测试token有效性..."
IMMEDIATE_TEST=$(curl -s -X POST "$BASE_URL/auth/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "立即测试响应: $IMMEDIATE_TEST"

# 等待3秒后再次测试
echo "3. 等待3秒后再次测试..."
sleep 3

DELAYED_TEST=$(curl -s -X POST "$BASE_URL/auth/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "延迟测试响应: $DELAYED_TEST"

echo "=== 测试完成 ===" 