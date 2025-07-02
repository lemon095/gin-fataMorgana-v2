#!/bin/bash

# 测试拼单详情接口
echo "=== 测试拼单详情接口 ==="

# 设置API基础URL
API_BASE="http://localhost:8080/api/v1"

# 测试获取活跃拼单详情（按时间最近）
echo "1. 测试获取活跃拼单详情（按时间最近）..."
ACTIVE_DETAIL_RESPONSE=$(curl -s -X POST "$API_BASE/groupBuy/active-detail" \
  -H "Content-Type: application/json" \
  -d '{}')

echo "活跃拼单详情响应: $ACTIVE_DETAIL_RESPONSE"

# 提取拼单编号用于后续测试
GROUP_BUY_NO=$(echo $ACTIVE_DETAIL_RESPONSE | grep -o '"group_buy_no":"[^"]*"' | cut -d'"' -f4)
echo "拼单编号: $GROUP_BUY_NO"

# 测试获取活跃拼单详情（随机返回）
echo -e "\n2. 测试获取活跃拼单详情（随机返回）..."
RANDOM_DETAIL_RESPONSE=$(curl -s -X POST "$API_BASE/groupBuy/active-detail?random=true" \
  -H "Content-Type: application/json" \
  -d '{}')

echo "随机拼单详情响应: $RANDOM_DETAIL_RESPONSE"

# 测试空请求体
echo -e "\n3. 测试空请求体..."
EMPTY_RESPONSE=$(curl -s -X POST "$API_BASE/groupBuy/active-detail" \
  -H "Content-Type: application/json")

echo "空请求体响应: $EMPTY_RESPONSE"

# 测试用户登录
echo -e "\n4. 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser1",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token: $TOKEN"

if [ -n "$TOKEN" ] && [ -n "$GROUP_BUY_NO" ]; then
    # 测试确认参与拼单
    echo -e "\n5. 测试确认参与拼单..."
    JOIN_RESPONSE=$(curl -s -X POST "$API_BASE/groupBuy/join" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{
        \"group_buy_no\": \"$GROUP_BUY_NO\"
      }")

    echo "确认参与拼单响应: $JOIN_RESPONSE"

    # 测试参与不存在的拼单
    echo -e "\n6. 测试参与不存在的拼单..."
    JOIN_INVALID_RESPONSE=$(curl -s -X POST "$API_BASE/groupBuy/join" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "group_buy_no": "NOTEXIST"
      }')

    echo "参与不存在拼单响应: $JOIN_INVALID_RESPONSE"
else
    echo "跳过确认参与拼单测试（需要有效的token和拼单编号）"
fi

echo -e "\n=== 测试完成 ==="
echo "接口功能说明："
echo "1. POST /api/v1/groupBuy/active-detail - 获取活跃拼单详情"
echo "   - 查询条件：AutoStart=true, deadline>当前时间, complete=cancelled"
echo "   - 支持random参数：true=随机返回，false=按时间最近返回"
echo "   - 返回字段：拼单编号、类型、总金额、当前参与人数、最大参与人数、已付款金额、每人付款金额、还需付款金额、截止时间"
echo "2. POST /api/v1/groupBuy/join - 确认参与拼单"
echo "   - 需要认证：Authorization: Bearer {token}"
echo "   - 请求参数：group_buy_no (拼单编号)"
echo "   - 业务逻辑：检查拼单状态，创建订单，更新拼单信息"
echo "   - 返回字段：order_id (订单ID)" 