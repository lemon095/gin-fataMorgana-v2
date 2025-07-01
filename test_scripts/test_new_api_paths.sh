#!/bin/bash

# 测试新的API路径结构
BASE_URL="http://localhost:9001"

echo "🧪 测试新的API路径结构 (/api/v1)"
echo "=================================="
echo ""

# 1. 测试根路径
echo "1. 测试根路径..."
curl -s "$BASE_URL/" | jq .
echo ""

# 2. 测试健康检查（保持原有路径）
echo "2. 测试健康检查（原有路径）..."
curl -s "$BASE_URL/health" | jq .
echo ""

# 3. 测试新的健康检查路径
echo "3. 测试新的健康检查路径..."
curl -s "$BASE_URL/api/v1/health/check" | jq .
echo ""

# 4. 测试数据库健康检查
echo "4. 测试数据库健康检查..."
curl -s "$BASE_URL/api/v1/health/database" | jq .
echo ""

# 5. 测试Redis健康检查
echo "5. 测试Redis健康检查..."
curl -s "$BASE_URL/api/v1/health/redis" | jq .
echo ""

# 6. 测试会话状态（新路径）
echo "6. 测试会话状态（新路径）..."
curl -s "$BASE_URL/api/v1/session/status" | jq .
echo ""

# 7. 测试注册接口（新路径）
echo "7. 测试注册接口（新路径）..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_api_path@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "INVITE123"
  }')

echo "$REGISTER_RESPONSE" | jq .
echo ""

# 8. 测试登录接口（新路径）
echo "8. 测试登录接口（新路径）..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_api_path@example.com",
    "password": "123456"
  }')

echo "$LOGIN_RESPONSE" | jq .
echo ""

# 提取访问令牌
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.tokens.access_token')

if [ "$ACCESS_TOKEN" != "null" ] && [ "$ACCESS_TOKEN" != "" ]; then
    echo "✅ 登录成功，测试需要认证的接口"
    echo ""

    # 9. 测试获取用户信息（新路径）
    echo "9. 测试获取用户信息（新路径）..."
    curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/api/v1/auth/profile" | jq .
    echo ""

    # 10. 测试获取钱包信息（新路径）
    echo "10. 测试获取钱包信息（新路径）..."
    curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/api/v1/wallet/info" | jq .
    echo ""

    # 11. 测试获取交易记录（新路径）
    echo "11. 测试获取交易记录（新路径）..."
    curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/api/v1/wallet/transactions" | jq .
    echo ""

    # 12. 测试会话用户信息（新路径）
    echo "12. 测试会话用户信息（新路径）..."
    curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "$BASE_URL/api/v1/session/user" | jq .
    echo ""

else
    echo "❌ 登录失败，跳过需要认证的接口测试"
fi

echo "🎉 新API路径测试完成！"
echo ""
echo "📋 新的API路径结构总结："
echo "   - 根路径: /"
echo "   - 健康检查: /health (保持原有)"
echo "   - API v1: /api/v1"
echo "   - 认证接口: /api/v1/auth/*"
echo "   - 会话接口: /api/v1/session/*"
echo "   - 钱包接口: /api/v1/wallet/*"
echo "   - 管理员接口: /api/v1/admin/*"
echo "   - 健康检查: /api/v1/health/*"
echo ""
echo "✅ 所有接口现在都有统一的 /api/v1 前缀！" 