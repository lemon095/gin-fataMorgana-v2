#!/bin/bash

# 测试登录态校验功能
BASE_URL="http://localhost:9001"

echo "🧪 开始测试登录态校验功能..."
echo ""

# 1. 检查服务是否启动
echo "1. 检查服务状态..."
curl -s "$BASE_URL/health" | jq .
echo ""

# 2. 检查未登录状态
echo "2. 检查未登录状态..."
curl -s "$BASE_URL/session/status" | jq .
echo ""

# 3. 注册用户
echo "3. 注册新用户..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "INVITE123"
  }')

echo "$REGISTER_RESPONSE" | jq .
echo ""

# 4. 用户登录
echo "4. 用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456"
  }')

echo "$LOGIN_RESPONSE" | jq .
echo ""

# 提取访问令牌
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.tokens.access_token')

if [ "$ACCESS_TOKEN" != "null" ] && [ "$ACCESS_TOKEN" != "" ]; then
    echo "✅ 登录成功，获取到访问令牌"
    echo ""

    # 5. 检查登录状态
    echo "5. 检查登录状态..."
    curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/session/status" | jq .
    echo ""

    # 6. 获取用户信息
    echo "6. 获取用户信息..."
    curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/session/user" | jq .
    echo ""

    # 7. 访问需要认证的接口
    echo "7. 访问需要认证的接口..."
    curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/api/profile" | jq .
    echo ""

    # 8. 访问可选认证的接口
    echo "8. 访问可选认证的接口..."
    curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/public/info" | jq .
    echo ""

    # 9. 测试无token访问需要认证的接口
    echo "9. 测试无token访问需要认证的接口..."
    curl -s "$BASE_URL/api/profile" | jq .
    echo ""

    # 10. 测试无token访问可选认证的接口
    echo "10. 测试无token访问可选认证的接口..."
    curl -s "$BASE_URL/public/info" | jq .
    echo ""

    # 11. 刷新会话
    echo "11. 刷新会话..."
    curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/session/refresh" | jq .
    echo ""

    # 12. 用户登出
    echo "12. 用户登出..."
    curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/session/logout" | jq .
    echo ""

    # 13. 登出后再次检查状态
    echo "13. 登出后再次检查状态..."
    curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/session/status" | jq .
    echo ""

else
    echo "❌ 登录失败，无法获取访问令牌"
fi

echo "🎉 测试完成！" 