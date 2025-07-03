#!/bin/bash

# 测试订单创建时期数校验功能
echo "=== 测试订单创建时期数校验功能 ==="

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 测试1: 使用正确的期数创建订单
echo "1. 测试使用正确的期数创建订单"
curl -X POST "${BASE_URL}/order/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "period_number": "20241201001",
    "amount": 100.00,
    "like_count": 10,
    "share_count": 5,
    "follow_count": 3,
    "favorite_count": 2
  }' \
  | jq '.'

echo ""

# 测试2: 使用不存在的期数创建订单
echo "2. 测试使用不存在的期数创建订单"
curl -X POST "${BASE_URL}/order/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "period_number": "INVALID_PERIOD",
    "amount": 100.00,
    "like_count": 10,
    "share_count": 5,
    "follow_count": 3,
    "favorite_count": 2
  }' \
  | jq '.'

echo ""

# 测试3: 缺少期数参数
echo "3. 测试缺少期数参数"
curl -X POST "${BASE_URL}/order/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "amount": 100.00,
    "like_count": 10,
    "share_count": 5,
    "follow_count": 3,
    "favorite_count": 2
  }' \
  | jq '.'

echo ""
echo "=== 测试完成 ===" 