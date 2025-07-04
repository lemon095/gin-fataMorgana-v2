#!/bin/bash

# 测试用户待审核功能
echo "=== 测试用户待审核功能 ==="

BASE_URL="http://localhost:8080/api/v1"

# 测试1: 用户注册（应该默认为待审核状态）
echo "测试1: 用户注册（应该默认为待审核状态）"
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_approval@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "TEST123"
  }')

echo "注册响应:"
echo "$REGISTER_RESPONSE" | jq '.'

# 提取用户信息
USER_STATUS=$(echo "$REGISTER_RESPONSE" | jq -r '.data.user.status')
echo "用户状态: $USER_STATUS"

if [ "$USER_STATUS" = "2" ]; then
    echo "✅ 用户注册成功，状态为待审核(2)"
else
    echo "❌ 用户注册失败，状态不是待审核"
fi

echo -e "\n"

# 测试2: 待审核用户尝试登录（应该失败）
echo "测试2: 待审核用户尝试登录（应该失败）"
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_approval@example.com",
    "password": "123456"
  }')

echo "登录响应:"
echo "$LOGIN_RESPONSE" | jq '.'

LOGIN_MESSAGE=$(echo "$LOGIN_RESPONSE" | jq -r '.message')
if [[ "$LOGIN_MESSAGE" == *"待审核"* ]]; then
    echo "✅ 待审核用户登录被正确拒绝"
else
    echo "❌ 待审核用户登录没有被正确拒绝"
fi

echo -e "\n"

# 测试3: 模拟管理员审核通过（将状态改为1）
echo "测试3: 模拟管理员审核通过（将状态改为1）"
echo "请在数据库中执行以下SQL来模拟审核通过："
echo "UPDATE users SET status = 1 WHERE email = 'test_approval@example.com';"
echo ""

# 测试4: 审核通过后用户登录（应该成功）
echo "测试4: 审核通过后用户登录（应该成功）"
echo "请先执行上面的SQL，然后按回车继续测试..."
read -p "按回车继续..."

LOGIN_RESPONSE_APPROVED=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_approval@example.com",
    "password": "123456"
  }')

echo "审核后登录响应:"
echo "$LOGIN_RESPONSE_APPROVED" | jq '.'

LOGIN_CODE=$(echo "$LOGIN_RESPONSE_APPROVED" | jq -r '.code')
if [ "$LOGIN_CODE" = "0" ]; then
    echo "✅ 审核通过后用户登录成功"
else
    echo "❌ 审核通过后用户登录失败"
fi

echo -e "\n"

# 测试5: 测试禁用用户登录（应该失败）
echo "测试5: 测试禁用用户登录（应该失败）"
echo "请在数据库中执行以下SQL来模拟禁用用户："
echo "UPDATE users SET status = 0 WHERE email = 'test_approval@example.com';"
echo ""

echo "请先执行上面的SQL，然后按回车继续测试..."
read -p "按回车继续..."

LOGIN_RESPONSE_DISABLED=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_approval@example.com",
    "password": "123456"
  }')

echo "禁用后登录响应:"
echo "$LOGIN_RESPONSE_DISABLED" | jq '.'

LOGIN_MESSAGE_DISABLED=$(echo "$LOGIN_RESPONSE_DISABLED" | jq -r '.message')
if [[ "$LOGIN_MESSAGE_DISABLED" == *"禁用"* ]]; then
    echo "✅ 禁用用户登录被正确拒绝"
else
    echo "❌ 禁用用户登录没有被正确拒绝"
fi

echo -e "\n"

echo "=== 测试完成 ==="
echo ""
echo "用户状态说明："
echo "0: 禁用"
echo "1: 正常"
echo "2: 待审核" 