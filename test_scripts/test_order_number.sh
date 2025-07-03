#!/bin/bash

# 测试用户订单编号功能
echo "=== 测试用户订单编号功能 ==="

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 测试1: 创建第一个订单
echo "1. 创建第一个订单"
curl -X POST "${BASE_URL}/order/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "period_number": "20241201001",
    "amount": 14.8,
    "like_count": 10,
    "share_count": 5,
    "follow_count": 3,
    "favorite_count": 2
  }' \
  | jq '.'

echo ""

# 测试2: 创建第二个订单（同一个用户）
echo "2. 创建第二个订单（同一个用户）"
curl -X POST "${BASE_URL}/order/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "period_number": "20241201001",
    "amount": 20.0,
    "like_count": 15,
    "share_count": 8,
    "follow_count": 5,
    "favorite_count": 3
  }' \
  | jq '.'

echo ""

# 测试3: 查看订单列表，验证number字段
echo "3. 查看订单列表，验证number字段"
curl -X POST "${BASE_URL}/order/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 3
  }' \
  | jq '.'

echo ""
echo "=== 测试完成 ==="

# 预期结果：
# 第一个订单的number应该是: {uid}_1
# 第二个订单的number应该是: {uid}_2
# 每个用户的订单编号都是独立的，从1开始递增 