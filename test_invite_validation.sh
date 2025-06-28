#!/bin/bash

# 测试邀请码校验功能
BASE_URL="http://localhost:8080"

echo "=== 测试邀请码校验功能 ==="
echo

# 1. 测试无效邀请码注册
echo "1. 测试无效邀请码注册..."
INVALID_INVITE_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_invalid@example.com",
    "password": "password123",
    "confirm_password": "password123",
    "invite_code": "INVALID"
  }')

echo "无效邀请码注册响应: $INVALID_INVITE_RESPONSE"
echo

# 2. 测试空邀请码注册
echo "2. 测试空邀请码注册..."
EMPTY_INVITE_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_empty@example.com",
    "password": "password123",
    "confirm_password": "password123",
    "invite_code": ""
  }')

echo "空邀请码注册响应: $EMPTY_INVITE_RESPONSE"
echo

# 3. 测试有效邀请码注册（需要先在数据库中创建管理员用户）
echo "3. 测试有效邀请码注册..."
echo "请先在数据库中创建管理员用户，然后使用其邀请码进行测试"
echo "可以使用以下SQL查询获取邀请码："
echo "SELECT my_invite_code FROM admin_users WHERE status = 1 LIMIT 1;"
echo

# 4. 测试重复邮箱注册
echo "4. 测试重复邮箱注册..."
DUPLICATE_EMAIL_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_invalid@example.com",
    "password": "password123",
    "confirm_password": "password123",
    "invite_code": "INVALID"
  }')

echo "重复邮箱注册响应: $DUPLICATE_EMAIL_RESPONSE"
echo

# 5. 测试密码不匹配
echo "5. 测试密码不匹配..."
PASSWORD_MISMATCH_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_mismatch@example.com",
    "password": "password123",
    "confirm_password": "password456",
    "invite_code": "INVALID"
  }')

echo "密码不匹配注册响应: $PASSWORD_MISMATCH_RESPONSE"
echo

echo "=== 测试完成 ==="
echo
echo "📊 测试总结:"
echo "- 无效邀请码应该被拒绝"
echo "- 空邀请码应该被拒绝"
echo "- 重复邮箱应该被拒绝"
echo "- 密码不匹配应该被拒绝"
echo
echo "💡 要测试有效邀请码，请先："
echo "1. 运行 ./init_admin.sh 创建管理员用户"
echo "2. 从数据库查询邀请码"
echo "3. 使用有效邀请码进行注册测试"
echo
echo "🔍 数据库查询邀请码的SQL:"
echo "SELECT username, my_invite_code, status FROM admin_users WHERE status = 1;" 