#!/bin/bash

# 测试响应格式统一性
echo "=== 响应格式统一性测试 ==="

# 设置基础URL
BASE_URL="http://localhost:9001/api"

# 1. 测试健康检查接口（不需要认证）
echo "1. 测试健康检查接口..."
HEALTH_RESPONSE=$(curl -s -X GET "$BASE_URL/v1/health/system")

echo "健康检查响应: $HEALTH_RESPONSE"

# 检查响应格式
if echo "$HEALTH_RESPONSE" | grep -q '"code":0'; then
    echo "✅ 健康检查接口响应格式正确"
else
    echo "❌ 健康检查接口响应格式错误"
    echo "期望: code:0, 实际: $HEALTH_RESPONSE"
fi

# 2. 测试登录接口
echo ""
echo "2. 测试登录接口..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test@example.com",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
if [ ! -z "$TOKEN" ]; then
    echo "✅ 登录成功，获取到token"
    
    # 3. 测试钱包信息接口
    echo ""
    echo "3. 测试钱包信息接口..."
    WALLET_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/wallet/info" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{}')
    
    echo "钱包信息响应: $WALLET_RESPONSE"
    
    # 检查响应格式
    if echo "$WALLET_RESPONSE" | grep -q '"code":0'; then
        echo "✅ 钱包信息接口响应格式正确"
    else
        echo "❌ 钱包信息接口响应格式错误"
        echo "期望: code:0, 实际: $WALLET_RESPONSE"
    fi
else
    echo "❌ 登录失败，无法测试需要认证的接口"
fi

# 4. 测试错误响应格式
echo ""
echo "4. 测试错误响应格式..."
ERROR_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/wallet/info" \
  -H "Content-Type: application/json" \
  -d '{}')

echo "错误响应: $ERROR_RESPONSE"

# 检查错误响应格式
if echo "$ERROR_RESPONSE" | grep -q '"code":401'; then
    echo "✅ 错误响应格式正确"
else
    echo "❌ 错误响应格式错误"
    echo "期望: code:401, 实际: $ERROR_RESPONSE"
fi

echo ""
echo "✅ 响应格式测试完成" 