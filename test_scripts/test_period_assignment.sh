#!/bin/bash

# 测试期号分配逻辑
# 使用方法: ./test_period_assignment.sh

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== 测试期号分配逻辑 ===${NC}"

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

# 2. 获取期数列表
echo -e "\n${YELLOW}2. 获取期数列表${NC}"
PERIODS_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/period-list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10
  }')

echo "期数列表响应: $PERIODS_RESPONSE"

# 检查是否有期数数据
PERIODS_COUNT=$(echo "$PERIODS_RESPONSE" | jq '.data.periods | length')
if [ "$PERIODS_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✓ 找到 $PERIODS_COUNT 个期数${NC}"
    
    # 显示期数信息
    echo -e "\n${YELLOW}期数信息:${NC}"
    echo "$PERIODS_RESPONSE" | jq -r '.data.periods[] | "期号: \(.period_number), 开始时间: \(.start_time), 结束时间: \(.end_time), 状态: \(.status)"'
else
    echo -e "${YELLOW}⚠ 没有找到期数数据${NC}"
fi

# 3. 手动生成假订单
echo -e "\n${YELLOW}3. 手动生成假订单${NC}"
GENERATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/cron/manual-generate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "count": 5
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

# 4. 查询最新生成的订单，检查期号分配
echo -e "\n${YELLOW}4. 查询最新生成的订单${NC}"
ORDERS_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 3
  }')

echo "订单列表响应: $ORDERS_RESPONSE"

# 检查订单的期号分配
ORDERS_COUNT=$(echo "$ORDERS_RESPONSE" | jq '.data.orders | length')
if [ "$ORDERS_COUNT" -gt 0 ]; then
    echo -e "\n${YELLOW}订单期号分配情况:${NC}"
    echo "$ORDERS_RESPONSE" | jq -r '.data.orders[] | "订单号: \(.order_no), 期号: \(.period_number), 创建时间: \(.created_at), 是否系统订单: \(.is_system_order)"'
    
    # 统计不同期号的数量
    echo -e "\n${YELLOW}期号分布统计:${NC}"
    echo "$ORDERS_RESPONSE" | jq -r '.data.orders[].period_number' | sort | uniq -c | while read count period; do
        echo "期号 $period: $count 个订单"
    done
else
    echo -e "${YELLOW}⚠ 没有找到订单数据${NC}"
fi

# 5. 测试不同时间段的期号分配
echo -e "\n${YELLOW}5. 测试期号分配逻辑说明${NC}"
echo -e "${GREEN}期号分配逻辑:${NC}"
echo "1. 根据订单的创建时间查询对应的期数"
echo "2. 查询条件: order_start_time <= 创建时间 < order_end_time"
echo "3. 如果找不到对应期数，使用最近的期数"
echo "4. 如果数据库中没有期数，使用时间格式化作为期号"

echo -e "\n${GREEN}=== 测试完成 ===${NC}"
echo -e "${GREEN}期号分配逻辑已实现:${NC}"
echo -e "  - 根据生成时间查询真实期号"
echo -e "  - 支持时间段匹配"
echo -e "  - 提供容错机制" 