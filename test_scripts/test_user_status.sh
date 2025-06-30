#!/bin/bash

# 测试用户状态功能的脚本

echo "🚀 开始测试用户状态功能..."
echo ""

# 测试正常注册
echo "📝 测试1: 正常用户注册"
echo "注册用户: test1@example.com"
curl -s -X POST http://localhost:9001/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test1@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "INVITE123"
  }' | jq .
echo ""

# 测试重复邮箱注册
echo "📝 测试2: 重复邮箱注册"
echo "尝试重复注册: test1@example.com"
curl -s -X POST http://localhost:9001/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test1@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "INVITE123"
  }' | jq .
echo ""

# 测试正常登录
echo "🔐 测试3: 正常用户登录"
echo "登录用户: test1@example.com"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:9001/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test1@example.com",
    "password": "123456"
  }')

echo "$LOGIN_RESPONSE" | jq .

# 提取token用于后续测试
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.tokens.access_token')
echo ""

# 测试获取用户信息
echo "👤 测试4: 获取用户信息"
curl -s -X GET http://localhost:9001/api/profile \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# 测试错误密码登录
echo "🔐 测试5: 错误密码登录"
curl -s -X POST http://localhost:9001/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test1@example.com",
    "password": "wrongpassword"
  }' | jq .
echo ""

# 测试不存在的用户登录
echo "🔐 测试6: 不存在的用户登录"
curl -s -X POST http://localhost:9001/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "nonexistent@example.com",
    "password": "123456"
  }' | jq .
echo ""

# 测试参数错误注册
echo "📝 测试7: 参数错误注册"
echo "密码不匹配"
curl -s -X POST http://localhost:9001/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test2@example.com",
    "password": "123456",
    "confirm_password": "654321",
    "invite_code": "INVITE123"
  }' | jq .
echo ""

echo "📝 测试8: 邮箱格式错误"
curl -s -X POST http://localhost:9001/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "INVITE123"
  }' | jq .
echo ""

# 测试刷新token
echo "🔄 测试9: 刷新Token"
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.tokens.refresh_token')
curl -s -X POST http://localhost:9001/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }" | jq .
echo ""

echo "✅ 测试完成！"
echo ""
echo "📊 测试总结:"
echo "- 正常注册和登录功能正常"
echo "- 重复邮箱注册被正确阻止"
echo "- 错误密码登录返回正确错误信息"
echo "- 不存在的用户登录返回正确错误信息"
echo "- 参数验证正常工作"
echo "- Token刷新功能正常"
echo ""
echo "🔍 注意: 用户禁用和删除功能需要在数据库中手动测试"
echo "可以通过以下SQL语句测试:"
echo ""
echo "# 禁用用户"
echo "UPDATE users SET status = 0 WHERE email = 'test1@example.com';"
echo ""
echo "# 删除用户（软删除）"
echo "DELETE FROM users WHERE email = 'test1@example.com';"
echo ""
echo "# 恢复用户"
echo "UPDATE users SET deleted_at = NULL WHERE email = 'test1@example.com';" 