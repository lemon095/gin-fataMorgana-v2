#!/bin/bash

# 测试订单金额校验功能
echo "=== 测试订单金额校验功能 ==="

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 首先设置Redis中的价格配置
echo "1. 设置Redis价格配置"
redis-cli set purchase_config '{
  "like_amount": 0.5,
  "share_amount": 1.0,
  "forward_amount": 0.8,
  "favorite_amount": 1.2,
  "created_at": "2025-07-03T20:18:13.463409+08:00",
  "updated_at": "2025-07-03T20:18:13.46341+08:00"
}'

echo ""

# 测试1: 金额匹配的情况
# 计算：10*0.5 + 5*1.0 + 3*0.8 + 2*1.2 = 5 + 5 + 2.4 + 2.4 = 14.8
echo "2. 测试金额匹配的情况（计算金额: 14.8）"
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

# 测试2: 金额不匹配的情况
echo "3. 测试金额不匹配的情况（请求金额: 15.0，计算金额: 14.8）"
curl -X POST "${BASE_URL}/order/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "period_number": "20241201001",
    "amount": 15.0,
    "like_count": 10,
    "share_count": 5,
    "follow_count": 3,
    "favorite_count": 2
  }' \
  | jq '.'

echo ""

# 测试3: 只有点赞的情况
# 计算：20*0.5 = 10.0
echo "4. 测试只有点赞的情况（计算金额: 10.0）"
curl -X POST "${BASE_URL}/order/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "period_number": "20241201001",
    "amount": 10.0,
    "like_count": 20,
    "share_count": 0,
    "follow_count": 0,
    "favorite_count": 0
  }' \
  | jq '.'

echo ""
echo "=== 测试完成 ===" 