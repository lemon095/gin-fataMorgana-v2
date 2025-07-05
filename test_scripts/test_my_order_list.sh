#!/bin/bash

# 测试获取我的订单列表API
echo "=== 测试获取我的订单列表API ==="

# 设置基础URL
BASE_URL="http://localhost:8080"
API_VERSION="v1"

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

# 2. 测试获取我的订单列表（进行中的订单）
echo "2. 测试获取我的订单列表（进行中的订单）..."
MY_ORDER_LIST_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/${API_VERSION}/order/my-list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 1
  }')

echo "获取我的订单列表响应: $MY_ORDER_LIST_RESPONSE"

# 检查是否成功
if echo "$MY_ORDER_LIST_RESPONSE" | grep -q '"code":200'; then
    echo "✓ 获取我的订单列表成功"
else
    echo "✗ 获取我的订单列表失败"
    echo "$MY_ORDER_LIST_RESPONSE"
fi

# 3. 测试获取我的订单列表（已完成的订单）
echo "3. 测试获取我的订单列表（已完成的订单）..."
MY_ORDER_LIST_COMPLETED_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/${API_VERSION}/order/my-list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 2
  }')

echo "获取我的已完成订单列表响应: $MY_ORDER_LIST_COMPLETED_RESPONSE"

# 检查是否成功
if echo "$MY_ORDER_LIST_COMPLETED_RESPONSE" | grep -q '"code":200'; then
    echo "✓ 获取我的已完成订单列表成功"
else
    echo "✗ 获取我的已完成订单列表失败"
    echo "$MY_ORDER_LIST_COMPLETED_RESPONSE"
fi

# 4. 测试获取我的订单列表（全部订单）
echo "4. 测试获取我的订单列表（全部订单）..."
MY_ORDER_LIST_ALL_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/${API_VERSION}/order/my-list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 3
  }')

echo "获取我的全部订单列表响应: $MY_ORDER_LIST_ALL_RESPONSE"

# 检查是否成功
if echo "$MY_ORDER_LIST_ALL_RESPONSE" | grep -q '"code":200'; then
    echo "✓ 获取我的全部订单列表成功"
else
    echo "✗ 获取我的全部订单列表失败"
    echo "$MY_ORDER_LIST_ALL_RESPONSE"
fi

# 5. 测试无token的情况
echo "5. 测试无token的情况..."
MY_ORDER_LIST_NO_TOKEN_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/${API_VERSION}/order/my-list" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 1
  }')

echo "无token的获取我的订单列表响应: $MY_ORDER_LIST_NO_TOKEN_RESPONSE"

# 检查是否返回401
if echo "$MY_ORDER_LIST_NO_TOKEN_RESPONSE" | grep -q '"code":401'; then
    echo "✓ 无token时正确返回401未授权"
else
    echo "✗ 无token时没有正确返回401"
fi

# 6. 对比原有接口和新接口的返回结果
echo "6. 对比原有接口和新接口的返回结果..."
ORIGINAL_ORDER_LIST_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/${API_VERSION}/order/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 1
  }')

echo "原有接口响应: $ORIGINAL_ORDER_LIST_RESPONSE"

# 比较两个接口的返回结果是否相同
if [ "$MY_ORDER_LIST_RESPONSE" = "$ORIGINAL_ORDER_LIST_RESPONSE" ]; then
    echo "✓ 新接口和原有接口返回结果相同"
else
    echo "✗ 新接口和原有接口返回结果不同"
fi

echo "=== 测试完成 ===" 