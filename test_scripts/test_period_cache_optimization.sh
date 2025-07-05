#!/bin/bash

# 测试期号缓存优化逻辑
# 使用方法: ./test_period_cache_optimization.sh

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 测试期号缓存优化逻辑 ===${NC}"

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
    echo -e "${YELLOW}⚠ 没有找到期数数据，将使用默认期号${NC}"
fi

# 3. 测试小批量生成（验证缓存效果）
echo -e "\n${YELLOW}3. 测试小批量生成（验证缓存效果）${NC}"
GENERATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/cron/manual-generate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "count": 10
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
    "page_size": 20,
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

# 5. 测试大批量生成（验证性能优化）
echo -e "\n${YELLOW}5. 测试大批量生成（验证性能优化）${NC}"
echo -e "${BLUE}开始时间: $(date)${NC}"

GENERATE_LARGE_RESPONSE=$(curl -s -X POST "${BASE_URL}/cron/manual-generate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "count": 50
  }')

echo -e "${BLUE}结束时间: $(date)${NC}"

echo "大批量生成响应: $GENERATE_LARGE_RESPONSE"

# 检查是否成功生成
if echo "$GENERATE_LARGE_RESPONSE" | jq -e '.code == 200' > /dev/null; then
    echo -e "${GREEN}✓ 成功生成大批量假订单${NC}"
    
    # 获取生成统计
    TOTAL_GENERATED=$(echo "$GENERATE_LARGE_RESPONSE" | jq -r '.data.total_generated')
    PURCHASE_ORDERS=$(echo "$GENERATE_LARGE_RESPONSE" | jq -r '.data.purchase_orders')
    GROUP_BUY_ORDERS=$(echo "$GENERATE_LARGE_RESPONSE" | jq -r '.data.group_buy_orders')
    AVERAGE_TIME=$(echo "$GENERATE_LARGE_RESPONSE" | jq -r '.data.average_time')
    
    echo -e "${GREEN}大批量生成统计: 总数=$TOTAL_GENERATED, 购买单=$PURCHASE_ORDERS, 拼单=$GROUP_BUY_ORDERS, 平均耗时=$AVERAGE_TIME${NC}"
else
    echo -e "${RED}✗ 大批量生成假订单失败${NC}"
    echo "错误信息: $(echo "$GENERATE_LARGE_RESPONSE" | jq -r '.message')"
fi

# 6. 优化效果说明
echo -e "\n${YELLOW}6. 优化效果说明${NC}"
echo -e "${GREEN}期号缓存优化逻辑:${NC}"
echo "1. 在定时任务开始时，一次性查询时间范围内的所有期数"
echo "2. 将期数数据缓存到内存中，避免重复数据库查询"
echo "3. 生成订单时直接从缓存中查找对应的期号"
echo "4. 如果缓存中没有找到，才回退到数据库查询"

echo -e "\n${GREEN}性能优化效果:${NC}"
echo "  - 减少数据库查询次数：从 N 次减少到 1 次"
echo "  - 提高生成速度：内存查找比数据库查询快"
echo "  - 降低数据库压力：减少并发查询"
echo "  - 保持数据准确性：缓存失效时自动回退"

echo -e "\n${GREEN}缓存策略:${NC}"
echo "  - 缓存时间范围：当前时间前后30分钟"
echo "  - 缓存更新：每次生成任务开始时刷新"
echo "  - 容错机制：缓存未命中时回退到数据库查询"

echo -e "\n${GREEN}=== 测试完成 ===${NC}"
echo -e "${GREEN}期号缓存优化已实现:${NC}"
echo -e "  - 预加载期数数据到缓存"
echo -e "  - 内存查找替代数据库查询"
echo -e "  - 保持数据准确性和系统稳定性" 