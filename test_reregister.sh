#!/bin/bash

# 测试删除状态用户重新注册功能
BASE_URL="http://localhost:9001"

echo "=== 测试删除状态用户重新注册功能 ==="
echo

# 1. 注册第一个用户
echo "1. 注册第一个用户..."
FIRST_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_reregister@example.com",
    "password": "password123",
    "confirm_password": "password123",
    "invite_code": "INVITE123"
  }')

echo "注册响应: $FIRST_USER_RESPONSE"
echo

# 提取第一个用户的ID
FIRST_USER_ID=$(echo $FIRST_USER_RESPONSE | grep -o '"uid":"[^"]*"' | cut -d'"' -f4)
echo "第一个用户ID: $FIRST_USER_ID"
echo

# 2. 删除第一个用户
echo "2. 删除第一个用户..."
DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/users/$FIRST_USER_ID" \
  -H "Content-Type: application/json")

echo "删除响应: $DELETE_RESPONSE"
echo

# 3. 尝试用相同邮箱重新注册
echo "3. 用相同邮箱重新注册..."
SECOND_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_reregister@example.com",
    "password": "newpassword123",
    "confirm_password": "newpassword123",
    "invite_code": "INVITE123"
  }')

echo "重新注册响应: $SECOND_USER_RESPONSE"
echo

# 4. 验证新用户是否可以登录
echo "4. 验证新用户是否可以登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_reregister@example.com",
    "password": "newpassword123"
  }')

echo "登录响应: $LOGIN_RESPONSE"
echo

# 5. 尝试用旧密码登录（应该失败）
echo "5. 尝试用旧密码登录（应该失败）..."
OLD_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_reregister@example.com",
    "password": "password123"
  }')

echo "旧密码登录响应: $OLD_LOGIN_RESPONSE"
echo

echo "=== 测试完成 ===" 