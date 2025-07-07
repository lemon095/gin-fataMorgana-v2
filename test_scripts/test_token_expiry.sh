#!/bin/bash

# 测试Token过期问题诊断脚本
echo "=== Token过期问题诊断 ==="

# 设置基础URL
BASE_URL="http://localhost:9001/api"

# 测试用户登录
echo "1. 用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test@example.com",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$TOKEN" ]; then
    echo "❌ 登录失败，无法获取token"
    exit 1
fi

echo "✅ 登录成功，获取到token: ${TOKEN:0:20}..."

# 立即测试token有效性
echo "2. 立即测试token有效性..."
IMMEDIATE_TEST=$(curl -s -X POST "$BASE_URL/auth/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "立即测试响应: $IMMEDIATE_TEST"

# 等待5秒后再次测试
echo "3. 等待5秒后再次测试..."
sleep 5

DELAYED_TEST=$(curl -s -X POST "$BASE_URL/auth/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "延迟测试响应: $DELAYED_TEST"

# 检查Redis中的token信息
echo "4. 检查Redis中的token信息..."
echo "请手动检查Redis中的以下key:"
echo "- user:active_token:*"
echo "- token:blacklist:*"

# 解析JWT token（需要安装jq）
echo "5. 解析JWT token信息..."
if command -v jq &> /dev/null; then
    # 分割JWT token
    IFS='.' read -ra JWT_PARTS <<< "$TOKEN"
    if [ ${#JWT_PARTS[@]} -eq 3 ]; then
        # 解码payload部分
        PAYLOAD=$(echo "${JWT_PARTS[1]}" | base64 -d 2>/dev/null | jq . 2>/dev/null)
        if [ $? -eq 0 ]; then
            echo "JWT Payload: $PAYLOAD"
        else
            echo "无法解析JWT payload"
        fi
    else
        echo "JWT token格式不正确"
    fi
else
    echo "请安装jq来解析JWT token"
fi

echo "=== 诊断完成 ===" 