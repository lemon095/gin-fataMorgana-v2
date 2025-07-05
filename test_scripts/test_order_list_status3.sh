#!/bin/bash

# 测试订单列表接口 status=3 的逻辑
BASE_URL="http://localhost:9001/api/v1"

echo "=== 测试订单列表接口 status=3 逻辑 ==="

# 1. 用户登录获取token
echo "1. 用户登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test@example.com",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.tokens.access_token' 2>/dev/null)

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "❌ 登录失败，无法获取token"
    echo "请检查用户账号密码是否正确"
    exit 1
fi

echo "✅ 登录成功，获取到token: ${TOKEN:0:20}..."
echo ""

# 2. 测试 status=1 (进行中)
echo "2. 测试 status=1 (进行中)..."
RESPONSE1=$(curl -s -X POST "${BASE_URL}/order/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 1
  }')

echo "status=1 响应: $RESPONSE1"

# 检查响应状态
if echo "$RESPONSE1" | jq -e '.code == 0' > /dev/null 2>&1; then
    echo "✅ status=1 请求成功"
    ORDERS_COUNT1=$(echo "$RESPONSE1" | jq '.data.orders | length' 2>/dev/null || echo "0")
    TOTAL1=$(echo "$RESPONSE1" | jq '.data.pagination.total' 2>/dev/null || echo "0")
    echo "   - 订单数量: $ORDERS_COUNT1"
    echo "   - 总数量: $TOTAL1"
else
    echo "❌ status=1 请求失败"
fi

echo ""

# 3. 测试 status=2 (已完成)
echo "3. 测试 status=2 (已完成)..."
RESPONSE2=$(curl -s -X POST "${BASE_URL}/order/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 2
  }')

echo "status=2 响应: $RESPONSE2"

# 检查响应状态
if echo "$RESPONSE2" | jq -e '.code == 0' > /dev/null 2>&1; then
    echo "✅ status=2 请求成功"
    ORDERS_COUNT2=$(echo "$RESPONSE2" | jq '.data.orders | length' 2>/dev/null || echo "0")
    TOTAL2=$(echo "$RESPONSE2" | jq '.data.pagination.total' 2>/dev/null || echo "0")
    echo "   - 订单数量: $ORDERS_COUNT2"
    echo "   - 总数量: $TOTAL2"
else
    echo "❌ status=2 请求失败"
fi

echo ""

# 4. 测试 status=3 (全部)
echo "4. 测试 status=3 (全部)..."
RESPONSE3=$(curl -s -X POST "${BASE_URL}/order/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 3
  }')

echo "status=3 响应: $RESPONSE3"

# 检查响应状态
if echo "$RESPONSE3" | jq -e '.code == 0' > /dev/null 2>&1; then
    echo "✅ status=3 请求成功"
    ORDERS_COUNT3=$(echo "$RESPONSE3" | jq '.data.orders | length' 2>/dev/null || echo "0")
    TOTAL3=$(echo "$RESPONSE3" | jq '.data.pagination.total' 2>/dev/null || echo "0")
    echo "   - 订单数量: $ORDERS_COUNT3"
    echo "   - 总数量: $TOTAL3"
    
    # 验证status=3返回的数据是否包含所有状态的订单
    if [ "$TOTAL3" -gt 0 ]; then
        echo "✅ status=3 成功返回订单数据"
        
        # 检查是否包含不同状态的订单
        PENDING_COUNT=$(echo "$RESPONSE3" | jq '[.data.orders[] | select(.status == "pending")] | length' 2>/dev/null || echo "0")
        SUCCESS_COUNT=$(echo "$RESPONSE3" | jq '[.data.orders[] | select(.status == "success")] | length' 2>/dev/null || echo "0")
        FAILED_COUNT=$(echo "$RESPONSE3" | jq '[.data.orders[] | select(.status == "failed")] | length' 2>/dev/null || echo "0")
        
        echo "   - 进行中订单: $PENDING_COUNT"
        echo "   - 已完成订单: $SUCCESS_COUNT"
        echo "   - 失败订单: $FAILED_COUNT"
    else
        echo "⚠️  status=3 返回空数据（可能是数据库中没有订单）"
    fi
else
    echo "❌ status=3 请求失败"
fi

echo ""

# 5. 数据对比分析
echo "5. 数据对比分析..."
echo "status=1 总数量: $TOTAL1"
echo "status=2 总数量: $TOTAL2"
echo "status=3 总数量: $TOTAL3"

# 验证status=3的总数应该大于等于status=1和status=2的总数
if [ "$TOTAL3" -ge "$TOTAL1" ] && [ "$TOTAL3" -ge "$TOTAL2" ]; then
    echo "✅ status=3 的数据量符合预期（包含所有状态的订单）"
else
    echo "⚠️  status=3 的数据量可能不符合预期"
fi

echo ""

# 6. 测试无效状态
echo "6. 测试无效状态 (status=4)..."
RESPONSE4=$(curl -s -X POST "${BASE_URL}/order/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 4
  }')

echo "status=4 响应: $RESPONSE4"

# 检查响应状态
if echo "$RESPONSE4" | jq -e '.code != 0' > /dev/null 2>&1; then
    echo "✅ 无效状态参数验证正确"
else
    echo "❌ 无效状态参数验证失败"
fi

echo ""
echo "=== 测试完成 ===" 