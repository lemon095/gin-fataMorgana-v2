#!/bin/bash

# 测试拼单资格字段功能
echo "🧪 测试拼单资格字段功能"
echo "================================"

BASE_URL="http://localhost:9001"

# 测试用户注册
echo "1. 测试用户注册..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "groupbuy_test@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "ADMIN1"
  }')

echo "注册响应: $REGISTER_RESPONSE"

# 提取用户ID
USER_ID=$(echo $REGISTER_RESPONSE | grep -o '"uid":"[^"]*"' | cut -d'"' -f4)
echo "用户ID: $USER_ID"

# 测试用户登录
echo "2. 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "groupbuy_test@example.com",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取访问令牌
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
echo "访问令牌: $ACCESS_TOKEN"

# 测试获取用户信息（检查拼单资格字段）
echo "3. 测试获取用户信息..."
PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "用户信息响应: $PROFILE_RESPONSE"

# 检查拼单资格字段是否存在
if echo "$PROFILE_RESPONSE" | grep -q "has_group_buy_qualification"; then
    echo "✅ 拼单资格字段已成功添加到用户信息中"
else
    echo "❌ 拼单资格字段未在用户信息中找到"
fi

# 检查拼单资格默认值
if echo "$PROFILE_RESPONSE" | grep -q '"has_group_buy_qualification":false'; then
    echo "✅ 拼单资格默认值为false"
else
    echo "❌ 拼单资格默认值不是false"
fi

echo "================================"
echo "测试完成！" 