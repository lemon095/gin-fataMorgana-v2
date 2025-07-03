#!/bin/bash

# 测试期号重复检查功能
BASE_URL="http://localhost:8080/api/v1"

echo "=== 测试期号重复检查功能 ==="

# 1. 用户登录获取token
echo "1. 用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')
UID=$(echo $LOGIN_RESPONSE | jq -r '.data.uid')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "登录失败: $LOGIN_RESPONSE"
    exit 1
fi

echo "登录成功，用户ID: $UID"

# 2. 创建第一个订单（期号: TEST001）
echo "2. 创建第一个订单（期号: TEST001）..."
ORDER1_RESPONSE=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "period_number": "TEST001",
    "amount": 10.0,
    "like_count": 1,
    "share_count": 0,
    "follow_count": 0,
    "favorite_count": 0
  }')

echo "第一个订单响应: $ORDER1_RESPONSE"

# 3. 尝试创建第二个订单（相同期号: TEST001）
echo "3. 尝试创建第二个订单（相同期号: TEST001）..."
ORDER2_RESPONSE=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "period_number": "TEST001",
    "amount": 20.0,
    "like_count": 2,
    "share_count": 0,
    "follow_count": 0,
    "favorite_count": 0
  }')

echo "第二个订单响应: $ORDER2_RESPONSE"

# 4. 创建第三个订单（不同期号: TEST002）
echo "4. 创建第三个订单（不同期号: TEST002）..."
ORDER3_RESPONSE=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "period_number": "TEST002",
    "amount": 15.0,
    "like_count": 1,
    "share_count": 1,
    "follow_count": 0,
    "favorite_count": 0
  }')

echo "第三个订单响应: $ORDER3_RESPONSE"

# 5. 获取订单列表验证
echo "5. 获取订单列表..."
ORDER_LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/orders?page=1&page_size=10&status=3" \
  -H "Authorization: Bearer $TOKEN")

echo "订单列表响应: $ORDER_LIST_RESPONSE"

echo "=== 测试完成 ===" 