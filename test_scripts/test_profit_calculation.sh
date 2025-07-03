#!/bin/bash

# 测试利润计算功能
# 需要先确保数据库中有 member_level 表的数据和用户数据

BASE_URL="http://localhost:8080/api"

echo "=== 测试利润计算功能 ==="

# 1. 先获取用户token（需要先注册或登录一个用户）
echo "1. 用户登录获取token"
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
echo "Token: $TOKEN"

if [ -z "$TOKEN" ]; then
    echo "获取token失败，请先注册用户或检查登录信息"
    exit 1
fi

# 2. 创建订单（测试利润计算）
echo "2. 创建订单测试利润计算"
ORDER_RESPONSE=$(curl -s -X POST "${BASE_URL}/orders/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "amount": 100.00,
    "like_count": 10,
    "share_count": 5,
    "follow_count": 3,
    "favorite_count": 2
  }')

echo "创建订单响应: $ORDER_RESPONSE"

# 3. 获取订单详情查看利润金额
echo "3. 获取订单详情查看利润金额"
ORDER_NO=$(echo $ORDER_RESPONSE | grep -o '"order_no":"[^"]*"' | cut -d'"' -f4)
echo "订单号: $ORDER_NO"

if [ ! -z "$ORDER_NO" ]; then
    ORDER_DETAIL_RESPONSE=$(curl -s -X GET "${BASE_URL}/orders/detail" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{
        \"order_no\": \"$ORDER_NO\"
      }")
    
    echo "订单详情响应: $ORDER_DETAIL_RESPONSE"
else
    echo "获取订单号失败"
fi

# 4. 测试拼单参与（测试利润计算）
echo "4. 获取拼单详情"
GROUP_BUY_RESPONSE=$(curl -s -X GET "${BASE_URL}/group-buy/detail" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN")

echo "拼单详情响应: $GROUP_BUY_RESPONSE"

# 提取拼单号
GROUP_BUY_NO=$(echo $GROUP_BUY_RESPONSE | grep -o '"group_buy_no":"[^"]*"' | cut -d'"' -f4)
echo "拼单号: $GROUP_BUY_NO"

if [ ! -z "$GROUP_BUY_NO" ]; then
    echo "5. 参与拼单测试利润计算"
    JOIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/group-buy/join" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{
        \"group_buy_no\": \"$GROUP_BUY_NO\"
      }")
    
    echo "参与拼单响应: $JOIN_RESPONSE"
    
    # 获取生成的订单ID
    ORDER_ID=$(echo $JOIN_RESPONSE | grep -o '"order_id":[0-9]*' | cut -d':' -f2)
    echo "生成的订单ID: $ORDER_ID"
else
    echo "获取拼单号失败，跳过拼单测试"
fi

# 6. 测试等级配置查询
echo "6. 测试等级配置查询"
LEVEL_RESPONSE=$(curl -s -X GET "${BASE_URL}/member-levels" \
  -H "Content-Type: application/json")

echo "等级配置响应: $LEVEL_RESPONSE"

# 7. 测试利润计算API
echo "7. 测试利润计算API"
PROFIT_RESPONSE=$(curl -s -X GET "${BASE_URL}/member-levels/calculate-cashback?experience=1&amount=100" \
  -H "Content-Type: application/json")

echo "利润计算响应: $PROFIT_RESPONSE"

echo "=== 测试完成 ==="
echo ""
echo "注意事项："
echo "1. 确保数据库中有 member_level 表的数据"
echo "2. 确保有测试用户数据"
echo "3. 检查订单表中的 profit_amount 字段是否正确计算"
echo "4. 不同等级的用户应该有不同的利润比例" 