#!/bin/bash

# 测试邀请码功能
BASE_URL="http://localhost:9001"

echo "=== 测试邀请码功能 ==="
echo

# 1. 注册第一个用户（不使用邀请码）
echo "1. 注册第一个用户（不使用邀请码）..."
FIRST_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user1@example.com",
    "password": "password123",
    "confirm_password": "password123",
    "invite_code": ""
  }')

echo "第一个用户注册响应: $FIRST_USER_RESPONSE"
echo

# 提取第一个用户的邀请码
FIRST_USER_INVITE_CODE=$(echo $FIRST_USER_RESPONSE | grep -o '"my_invite_code":"[^"]*"' | cut -d'"' -f4)
echo "第一个用户的邀请码: $FIRST_USER_INVITE_CODE"
echo

# 2. 使用第一个用户的邀请码注册第二个用户
echo "2. 使用第一个用户的邀请码注册第二个用户..."
SECOND_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"user2@example.com\",
    \"password\": \"password123\",
    \"confirm_password\": \"password123\",
    \"invite_code\": \"$FIRST_USER_INVITE_CODE\"
  }")

echo "第二个用户注册响应: $SECOND_USER_RESPONSE"
echo

# 提取第二个用户的邀请码和邀请关系
SECOND_USER_INVITE_CODE=$(echo $SECOND_USER_RESPONSE | grep -o '"my_invite_code":"[^"]*"' | cut -d'"' -f4)
SECOND_USER_INVITED_BY=$(echo $SECOND_USER_RESPONSE | grep -o '"invited_by":"[^"]*"' | cut -d'"' -f4)
echo "第二个用户的邀请码: $SECOND_USER_INVITE_CODE"
echo "第二个用户被邀请码: $SECOND_USER_INVITED_BY"
echo

# 3. 验证邀请码唯一性 - 尝试使用已存在的邀请码注册
echo "3. 验证邀请码唯一性 - 尝试使用已存在的邀请码注册..."
DUPLICATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"user3@example.com\",
    \"password\": \"password123\",
    \"confirm_password\": \"password123\",
    \"invite_code\": \"$FIRST_USER_INVITE_CODE\"
  }")

echo "重复邀请码注册响应: $DUPLICATE_RESPONSE"
echo

# 4. 验证无效邀请码
echo "4. 验证无效邀请码..."
INVALID_INVITE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user4@example.com",
    "password": "password123",
    "confirm_password": "password123",
    "invite_code": "INVALID"
  }')

echo "无效邀请码注册响应: $INVALID_INVITE_RESPONSE"
echo

# 5. 验证两个用户都可以登录
echo "5. 验证两个用户都可以登录..."
LOGIN1_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user1@example.com",
    "password": "password123"
  }')

LOGIN2_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user2@example.com",
    "password": "password123"
  }')

echo "用户1登录响应: $LOGIN1_RESPONSE"
echo "用户2登录响应: $LOGIN2_RESPONSE"
echo

echo "=== 测试完成 ===" 