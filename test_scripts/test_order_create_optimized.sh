#!/bin/bash

# 测试优化后的订单创建API
# 验证uid从token中获取，不需要在请求中传递

BASE_URL="http://localhost:8080"
API_VERSION="v1"

echo "=== 测试优化后的订单创建API ==="

# 1. 用户登录获取token
echo "1. 用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/${API_VERSION}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "登录失败，无法获取token"
    exit 1
fi

echo "登录成功，获取到token"

# 2. 创建订单（不传递uid参数）
echo "2. 创建订单（不传递uid参数）..."
CREATE_ORDER_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/${API_VERSION}/order/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "amount": 100.00,
    "profit_amount": 10.00,
    "like_count": 5,
    "share_count": 2,
    "follow_count": 1,
    "favorite_count": 3
  }')

echo "创建订单响应: $CREATE_ORDER_RESPONSE"

# 检查订单创建是否成功
if echo "$CREATE_ORDER_RESPONSE" | grep -q '"code":200'; then
    echo "✓ 订单创建成功"
    
    # 提取订单号
    ORDER_NO=$(echo "$CREATE_ORDER_RESPONSE" | grep -o '"order_no":"[^"]*"' | cut -d'"' -f4)
    echo "订单号: $ORDER_NO"
else
    echo "✗ 订单创建失败"
    echo "$CREATE_ORDER_RESPONSE"
fi

# 3. 测试传递uid参数的情况（应该被忽略）
echo "3. 测试传递uid参数的情况（应该被忽略）..."
CREATE_ORDER_WITH_UID_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/${API_VERSION}/order/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "uid": "fake_uid_123",
    "amount": 50.00,
    "profit_amount": 5.00,
    "like_count": 3,
    "share_count": 1,
    "follow_count": 0,
    "favorite_count": 2
  }')

echo "传递uid参数的创建订单响应: $CREATE_ORDER_WITH_UID_RESPONSE"

# 检查是否成功（应该成功，因为uid会被忽略）
if echo "$CREATE_ORDER_WITH_UID_RESPONSE" | grep -q '"code":200'; then
    echo "✓ 传递uid参数的订单创建成功（uid被忽略）"
    
    # 提取订单号
    ORDER_NO_2=$(echo "$CREATE_ORDER_WITH_UID_RESPONSE" | grep -o '"order_no":"[^"]*"' | cut -d'"' -f4)
    echo "订单号: $ORDER_NO_2"
else
    echo "✗ 传递uid参数的订单创建失败"
    echo "$CREATE_ORDER_WITH_UID_RESPONSE"
fi

# 4. 测试无token的情况
echo "4. 测试无token的情况..."
CREATE_ORDER_NO_TOKEN_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/${API_VERSION}/order/create" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.00,
    "profit_amount": 10.00,
    "like_count": 5,
    "share_count": 2,
    "follow_count": 1,
    "favorite_count": 3
  }')

echo "无token的创建订单响应: $CREATE_ORDER_NO_TOKEN_RESPONSE"

# 检查是否返回401
if echo "$CREATE_ORDER_NO_TOKEN_RESPONSE" | grep -q '"code":401'; then
    echo "✓ 无token时正确返回401未授权"
else
    echo "✗ 无token时没有正确返回401"
fi

echo "=== 测试完成 ===" 