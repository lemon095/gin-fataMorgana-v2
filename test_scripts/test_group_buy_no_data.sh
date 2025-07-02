#!/bin/bash

# 测试拼单详情接口 - 无数据情况
echo "=== 测试拼单详情接口 - 无数据情况 ==="

# 设置API基础URL
API_BASE="http://localhost:8080/api/v1"

# 测试获取活跃拼单详情（按时间最近）- 可能没有数据
echo "1. 测试获取活跃拼单详情（按时间最近）..."
ACTIVE_DETAIL_RESPONSE=$(curl -s -X POST "$API_BASE/groupBuy/active-detail" \
  -H "Content-Type: application/json" \
  -d '{}')

echo "活跃拼单详情响应: $ACTIVE_DETAIL_RESPONSE"

# 检查是否有数据
HAS_DATA=$(echo $ACTIVE_DETAIL_RESPONSE | grep -o '"has_data":[^,]*' | cut -d':' -f2)
echo "是否有数据: $HAS_DATA"

if [ "$HAS_DATA" = "true" ]; then
    # 提取拼单编号用于后续测试
    GROUP_BUY_NO=$(echo $ACTIVE_DETAIL_RESPONSE | grep -o '"group_buy_no":"[^"]*"' | cut -d'"' -f4)
    echo "拼单编号: $GROUP_BUY_NO"
    
    # 测试用户登录
    echo -e "\n2. 测试用户登录..."
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
        echo -e "\n3. 测试确认参与拼单..."
        JOIN_RESPONSE=$(curl -s -X POST "$API_BASE/groupBuy/join" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $TOKEN" \
          -d "{
            \"group_buy_no\": \"$GROUP_BUY_NO\"
          }")

        echo "确认参与拼单响应: $JOIN_RESPONSE"
    else
        echo "跳过确认参与拼单测试（需要有效的token和拼单编号）"
    fi
else
    echo "当前没有符合条件的活跃拼单数据"
    echo "查询条件：AutoStart=true, deadline>当前时间, complete=cancelled"
fi

# 测试获取活跃拼单详情（随机返回）
echo -e "\n4. 测试获取活跃拼单详情（随机返回）..."
RANDOM_DETAIL_RESPONSE=$(curl -s -X POST "$API_BASE/groupBuy/active-detail?random=true" \
  -H "Content-Type: application/json" \
  -d '{}')

echo "随机拼单详情响应: $RANDOM_DETAIL_RESPONSE"

echo -e "\n=== 测试完成 ==="
echo "接口说明："
echo "1. 当没有符合条件的拼单数据时，接口返回："
echo "   - has_data: false"
echo "   - 其他字段为空值或默认值"
echo "2. 当有符合条件的拼单数据时，接口返回："
echo "   - has_data: true"
echo "   - 完整的拼单详情信息"
echo "3. 查询条件：AutoStart=true, deadline>当前时间, complete=cancelled" 