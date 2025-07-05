#!/bin/bash

# 完整的用户信息获取测试脚本
# 包括登录和获取profile信息

BASE_URL="http://localhost:9001/api/v1"

echo "=== 完整的用户信息获取测试 ==="
echo ""

# 1. 用户登录
echo "🔐 步骤1: 用户登录"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test@example.com",
    "password": "123456"
  }')

echo "登录响应:"
echo "$LOGIN_RESPONSE" | jq '.' 2>/dev/null || echo "$LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.tokens.access_token' 2>/dev/null)

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo ""
    echo "❌ 登录失败，无法获取token"
    echo "请检查用户账号密码是否正确"
    exit 1
fi

echo ""
echo "✅ 登录成功，获取到token: ${TOKEN:0:20}..."
echo ""

# 2. 获取用户信息
echo "👤 步骤2: 获取用户信息"
PROFILE_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "用户信息响应:"
echo "$PROFILE_RESPONSE" | jq '.' 2>/dev/null || echo "$PROFILE_RESPONSE"

echo ""
echo "=== 测试完成 ==="

# 3. 对比session接口
echo ""
echo "🔄 步骤3: 对比session接口"
SESSION_RESPONSE=$(curl -s -X POST "$BASE_URL/session/user" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "Session用户信息响应:"
echo "$SESSION_RESPONSE" | jq '.' 2>/dev/null || echo "$SESSION_RESPONSE"

echo ""
echo "📊 接口对比说明:"
echo "- /auth/profile: 返回完整的用户资料信息（邮箱、经验值、信用分等）"
echo "- /session/user: 返回基本的会话信息（user_id、username、login_time）" 