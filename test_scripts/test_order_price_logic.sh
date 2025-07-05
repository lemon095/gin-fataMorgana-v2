#!/bin/bash

# 测试订单价格计算逻辑
echo "=== 测试订单价格计算逻辑 ==="

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 测试用户登录
echo "1. 用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$TOKEN" ]; then
    echo "❌ 登录失败，无法获取token"
    exit 1
fi

echo "✅ 登录成功，获取到token: ${TOKEN:0:20}..."

# 测试场景1：只选择点赞（单价10元，1个类型）
echo "2. 测试场景1：只选择点赞（单价10元，1个类型）..."
RESPONSE1=$(curl -s -X POST "$BASE_URL/order/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "period_number": "20250101001",
    "amount": 10.00,
    "like_count": 1,
    "share_count": 0,
    "follow_count": 0,
    "favorite_count": 0
  }')

echo "场景1响应: $RESPONSE1"

# 检查响应中的金额
AMOUNT1=$(echo $RESPONSE1 | grep -o '"amount":[0-9]*\.?[0-9]*' | cut -d':' -f2)
echo "场景1订单金额: $AMOUNT1"

if [ "$AMOUNT1" = "10" ]; then
    echo "✅ 场景1正确：单价10元 × 1个类型 = 10元"
else
    echo "❌ 场景1错误：期望10元，实际$AMOUNT1元"
fi

# 测试场景2：选择点赞和分享（单价15元，2个类型）
echo "3. 测试场景2：选择点赞和分享（单价15元，2个类型）..."
RESPONSE2=$(curl -s -X POST "$BASE_URL/order/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "period_number": "20250101002",
    "amount": 15.00,
    "like_count": 1,
    "share_count": 1,
    "follow_count": 0,
    "favorite_count": 0
  }')

echo "场景2响应: $RESPONSE2"

# 检查响应中的金额
AMOUNT2=$(echo $RESPONSE2 | grep -o '"amount":[0-9]*\.?[0-9]*' | cut -d':' -f2)
echo "场景2订单金额: $AMOUNT2"

if [ "$AMOUNT2" = "30" ]; then
    echo "✅ 场景2正确：单价15元 × 2个类型 = 30元"
else
    echo "❌ 场景2错误：期望30元，实际$AMOUNT2元"
fi

# 测试场景3：选择所有类型（单价20元，4个类型）
echo "4. 测试场景3：选择所有类型（单价20元，4个类型）..."
RESPONSE3=$(curl -s -X POST "$BASE_URL/order/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "period_number": "20250101003",
    "amount": 20.00,
    "like_count": 1,
    "share_count": 1,
    "follow_count": 1,
    "favorite_count": 1
  }')

echo "场景3响应: $RESPONSE3"

# 检查响应中的金额
AMOUNT3=$(echo $RESPONSE3 | grep -o '"amount":[0-9]*\.?[0-9]*' | cut -d':' -f2)
echo "场景3订单金额: $AMOUNT3"

if [ "$AMOUNT3" = "80" ]; then
    echo "✅ 场景3正确：单价20元 × 4个类型 = 80元"
else
    echo "❌ 场景3错误：期望80元，实际$AMOUNT3元"
fi

# 测试场景4：没有选择任何类型（应该失败）
echo "5. 测试场景4：没有选择任何类型（应该失败）..."
RESPONSE4=$(curl -s -X POST "$BASE_URL/order/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "period_number": "20250101004",
    "amount": 10.00,
    "like_count": 0,
    "share_count": 0,
    "follow_count": 0,
    "favorite_count": 0
  }')

echo "场景4响应: $RESPONSE4"

# 检查是否返回错误
if echo $RESPONSE4 | grep -q "请至少选择一种任务类型"; then
    echo "✅ 场景4正确：没有选择类型时返回错误"
else
    echo "❌ 场景4错误：应该返回错误但没有"
fi

echo "=== 测试完成 ===" 