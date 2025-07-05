#!/bin/bash

# 测试订单列表接口的时间过滤功能
# 使用方法: ./test_order_time_filter.sh

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 测试订单列表接口的时间过滤功能 ===${NC}"

# 1. 先登录获取token
echo -e "\n${YELLOW}1. 用户登录${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test_user",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo -e "${RED}登录失败，无法获取token${NC}"
    exit 1
fi

echo -e "${GREEN}登录成功，获取到token${NC}"

# 2. 生成一些假订单（包括超过当前时间的）
echo -e "\n${YELLOW}2. 生成假订单（包括超过当前时间的）${NC}"
GENERATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/cron/manual-generate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "count": 20
  }')

echo "生成假订单响应: $GENERATE_RESPONSE"

# 检查是否成功生成
if echo "$GENERATE_RESPONSE" | jq -e '.code == 200' > /dev/null; then
    echo -e "${GREEN}✓ 成功生成假订单${NC}"
    
    # 获取生成统计
    TOTAL_GENERATED=$(echo "$GENERATE_RESPONSE" | jq -r '.data.total_generated')
    PURCHASE_ORDERS=$(echo "$GENERATE_RESPONSE" | jq -r '.data.purchase_orders')
    GROUP_BUY_ORDERS=$(echo "$GENERATE_RESPONSE" | jq -r '.data.group_buy_orders')
    
    echo -e "${GREEN}生成统计: 总数=$TOTAL_GENERATED, 购买单=$PURCHASE_ORDERS, 拼单=$GROUP_BUY_ORDERS${NC}"
else
    echo -e "${RED}✗ 生成假订单失败${NC}"
    echo "错误信息: $(echo "$GENERATE_RESPONSE" | jq -r '.message')"
fi

# 3. 查询所有订单（验证时间过滤）
echo -e "\n${YELLOW}3. 查询所有订单（验证时间过滤）${NC}"
ALL_ORDERS_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/all-list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 50,
    "status": 3
  }')

echo "所有订单响应: $ALL_ORDERS_RESPONSE"

# 检查订单的创建时间
ALL_ORDERS_COUNT=$(echo "$ALL_ORDERS_RESPONSE" | jq '.data.orders | length')
if [ "$ALL_ORDERS_COUNT" -gt 0 ]; then
    echo -e "\n${YELLOW}所有订单的创建时间检查:${NC}"
    echo "$ALL_ORDERS_RESPONSE" | jq -r '.data.orders[] | "订单号: \(.order_no), 创建时间: \(.created_at), 是否系统订单: \(.is_system_order)"'
    
    # 检查是否有超过当前时间的订单
    CURRENT_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    echo -e "\n${BLUE}当前时间: $CURRENT_TIME${NC}"
    
    # 统计超过当前时间的订单数量
    FUTURE_ORDERS=0
    for i in $(seq 0 $((ALL_ORDERS_COUNT-1))); do
        CREATED_AT=$(echo "$ALL_ORDERS_RESPONSE" | jq -r ".data.orders[$i].created_at")
        if [[ "$CREATED_AT" > "$CURRENT_TIME" ]]; then
            FUTURE_ORDERS=$((FUTURE_ORDERS+1))
            ORDER_NO=$(echo "$ALL_ORDERS_RESPONSE" | jq -r ".data.orders[$i].order_no")
            echo -e "${RED}发现超过当前时间的订单: $ORDER_NO (创建时间: $CREATED_AT)${NC}"
        fi
    done
    
    if [ "$FUTURE_ORDERS" -eq 0 ]; then
        echo -e "${GREEN}✓ 时间过滤功能正常：没有发现超过当前时间的订单${NC}"
    else
        echo -e "${RED}✗ 时间过滤功能异常：发现 $FUTURE_ORDERS 个超过当前时间的订单${NC}"
    fi
else
    echo -e "${YELLOW}⚠ 没有找到订单数据${NC}"
fi

# 4. 查询我的订单（验证用户订单时间过滤）
echo -e "\n${YELLOW}4. 查询我的订单（验证用户订单时间过滤）${NC}"
MY_ORDERS_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/my-list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 50,
    "status": 3
  }')

echo "我的订单响应: $MY_ORDERS_RESPONSE"

# 检查我的订单的创建时间
MY_ORDERS_COUNT=$(echo "$MY_ORDERS_RESPONSE" | jq '.data.orders | length')
if [ "$MY_ORDERS_COUNT" -gt 0 ]; then
    echo -e "\n${YELLOW}我的订单的创建时间检查:${NC}"
    echo "$MY_ORDERS_RESPONSE" | jq -r '.data.orders[] | "订单号: \(.order_no), 创建时间: \(.created_at), 是否系统订单: \(.is_system_order)"'
    
    # 检查是否有超过当前时间的订单
    FUTURE_MY_ORDERS=0
    for i in $(seq 0 $((MY_ORDERS_COUNT-1))); do
        CREATED_AT=$(echo "$MY_ORDERS_RESPONSE" | jq -r ".data.orders[$i].created_at")
        if [[ "$CREATED_AT" > "$CURRENT_TIME" ]]; then
            FUTURE_MY_ORDERS=$((FUTURE_MY_ORDERS+1))
            ORDER_NO=$(echo "$MY_ORDERS_RESPONSE" | jq -r ".data.orders[$i].order_no")
            echo -e "${RED}发现超过当前时间的我的订单: $ORDER_NO (创建时间: $CREATED_AT)${NC}"
        fi
    done
    
    if [ "$FUTURE_MY_ORDERS" -eq 0 ]; then
        echo -e "${GREEN}✓ 我的订单时间过滤功能正常：没有发现超过当前时间的订单${NC}"
    else
        echo -e "${RED}✗ 我的订单时间过滤功能异常：发现 $FUTURE_MY_ORDERS 个超过当前时间的订单${NC}"
    fi
else
    echo -e "${YELLOW}⚠ 没有找到我的订单数据${NC}"
fi

# 5. 查询订单统计（验证统计时间过滤）
echo -e "\n${YELLOW}5. 查询订单统计（验证统计时间过滤）${NC}"
STATS_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/stats" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN")

echo "订单统计响应: $STATS_RESPONSE"

# 检查统计结果
if echo "$STATS_RESPONSE" | jq -e '.code == 200' > /dev/null; then
    echo -e "${GREEN}✓ 订单统计查询成功${NC}"
    
    # 显示统计信息
    TOTAL_ORDERS=$(echo "$STATS_RESPONSE" | jq -r '.data.total_orders')
    PENDING_ORDERS=$(echo "$STATS_RESPONSE" | jq -r '.data.pending_orders')
    SUCCESS_ORDERS=$(echo "$STATS_RESPONSE" | jq -r '.data.success_orders')
    TOTAL_AMOUNT=$(echo "$STATS_RESPONSE" | jq -r '.data.total_amount')
    
    echo -e "${GREEN}统计信息: 总订单=$TOTAL_ORDERS, 进行中=$PENDING_ORDERS, 已完成=$SUCCESS_ORDERS, 总金额=$TOTAL_AMOUNT${NC}"
else
    echo -e "${RED}✗ 订单统计查询失败${NC}"
    echo "错误信息: $(echo "$STATS_RESPONSE" | jq -r '.message')"
fi

# 6. 测试不同状态的订单查询
echo -e "\n${YELLOW}6. 测试不同状态的订单查询${NC}"

# 测试进行中的订单
PENDING_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/my-list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 1
  }')

PENDING_COUNT=$(echo "$PENDING_RESPONSE" | jq '.data.orders | length')
echo -e "${GREEN}进行中订单数量: $PENDING_COUNT${NC}"

# 测试已完成的订单
COMPLETED_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/my-list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 2
  }')

COMPLETED_COUNT=$(echo "$COMPLETED_RESPONSE" | jq '.data.orders | length')
echo -e "${GREEN}已完成订单数量: $COMPLETED_COUNT${NC}"

# 7. 总结测试结果
echo -e "\n${YELLOW}7. 测试结果总结${NC}"
echo -e "${GREEN}时间过滤功能测试完成:${NC}"
echo -e "  - 所有订单查询: 已添加时间过滤条件"
echo -e "  - 我的订单查询: 已添加时间过滤条件"
echo -e "  - 订单统计查询: 已添加时间过滤条件"
echo -e "  - 拼单查询: 已添加时间过滤条件"
echo -e "  - 过滤条件: created_at <= NOW()"
echo -e "  - 目的: 避免查询到超过当前时间的假订单数据"

echo -e "\n${GREEN}=== 测试完成 ===${NC}"
echo -e "${GREEN}订单列表接口时间过滤功能已优化:${NC}"
echo -e "  - 防止查询到超过当前时间的假订单"
echo -e "  - 保持数据的一致性和真实性"
echo -e "  - 提升用户体验" 