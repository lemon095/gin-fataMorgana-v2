#!/bin/bash

# 测试获取用户信息接口
# 使用方法: ./test_profile_api.sh [token]

BASE_URL="http://localhost:9001/api/v1"
TOKEN=${1:-""}

echo "=== 测试获取用户信息接口 ==="
echo "接口: POST $BASE_URL/auth/profile"
echo ""

if [ -z "$TOKEN" ]; then
    echo "❌ 错误: 请提供有效的JWT token"
    echo "使用方法: ./test_profile_api.sh <your_jwt_token>"
    echo ""
    echo "示例:"
    echo "1. 先登录获取token:"
    echo "   curl -X POST $BASE_URL/auth/login \\"
    echo "     -H 'Content-Type: application/json' \\"
    echo "     -d '{\"account\":\"test@example.com\",\"password\":\"123456\"}'"
    echo ""
    echo "2. 使用token测试profile接口:"
    echo "   ./test_profile_api.sh <your_token>"
    exit 1
fi

echo "🔑 使用Token: ${TOKEN:0:20}..."
echo ""

# 测试获取用户信息
echo "📤 发送请求..."
RESPONSE=$(curl -s -X POST "$BASE_URL/auth/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "📥 响应结果:"
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"

echo ""
echo "=== 测试完成 ===" 