#!/bin/bash

# 测试订单列表查询修复
# 使用方法: ./test_order_list_fix.sh

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== 测试订单列表查询修复 ===${NC}"

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

# 2. 测试查询进行中的订单 (status: 1)
echo -e "\n${YELLOW}2. 测试查询进行中的订单 (status: 1)${NC}"
IN_PROGRESS_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 1
  }')

echo "进行中订单响应: $IN_PROGRESS_RESPONSE"

# 检查响应状态
RESPONSE_CODE=$(echo "$IN_PROGRESS_RESPONSE" | jq -r '.code')
if [ "$RESPONSE_CODE" = "0" ]; then
    echo -e "${GREEN}查询进行中订单成功${NC}"
    
    # 解析响应数据
    TOTAL=$(echo "$IN_PROGRESS_RESPONSE" | jq -r '.data.pagination.total')
    ORDERS_COUNT=$(echo "$IN_PROGRESS_RESPONSE" | jq -r '.data.orders | length')
    
    echo -e "${GREEN}进行中订单总数: $TOTAL${NC}"
    echo -e "${GREEN}当前页订单数量: $ORDERS_COUNT${NC}"
    
    # 显示订单信息
    if [ "$ORDERS_COUNT" -gt 0 ]; then
        echo -e "\n${YELLOW}进行中订单详情:${NC}"
        for i in {0..2}; do
            if [ "$i" -lt "$ORDERS_COUNT" ]; then
                ORDER_NO=$(echo "$IN_PROGRESS_RESPONSE" | jq -r ".data.orders[$i].order_no")
                STATUS=$(echo "$IN_PROGRESS_RESPONSE" | jq -r ".data.orders[$i].status")
                AMOUNT=$(echo "$IN_PROGRESS_RESPONSE" | jq -r ".data.orders[$i].amount")
                UID=$(echo "$IN_PROGRESS_RESPONSE" | jq -r ".data.orders[$i].uid")
                CREATED_AT=$(echo "$IN_PROGRESS_RESPONSE" | jq -r ".data.orders[$i].created_at")
                
                echo -e "${GREEN}订单号: ${ORDER_NO} - 状态: ${STATUS} - 金额: ${AMOUNT} - UID: ${UID} - 创建时间: ${CREATED_AT}${NC}"
            fi
        done
    fi
else
    echo -e "${RED}查询进行中订单失败: $(echo "$IN_PROGRESS_RESPONSE" | jq -r '.message')${NC}"
fi

# 3. 测试查询已完成的订单 (status: 2)
echo -e "\n${YELLOW}3. 测试查询已完成的订单 (status: 2)${NC}"
COMPLETED_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 2
  }')

echo "已完成订单响应: $COMPLETED_RESPONSE"

# 检查响应状态
RESPONSE_CODE=$(echo "$COMPLETED_RESPONSE" | jq -r '.code')
if [ "$RESPONSE_CODE" = "0" ]; then
    echo -e "${GREEN}查询已完成订单成功${NC}"
    
    # 解析响应数据
    TOTAL=$(echo "$COMPLETED_RESPONSE" | jq -r '.data.pagination.total')
    ORDERS_COUNT=$(echo "$COMPLETED_RESPONSE" | jq -r '.data.orders | length')
    
    echo -e "${GREEN}已完成订单总数: $TOTAL${NC}"
    echo -e "${GREEN}当前页订单数量: $ORDERS_COUNT${NC}"
    
    # 显示订单信息
    if [ "$ORDERS_COUNT" -gt 0 ]; then
        echo -e "\n${YELLOW}已完成订单详情:${NC}"
        for i in {0..2}; do
            if [ "$i" -lt "$ORDERS_COUNT" ]; then
                ORDER_NO=$(echo "$COMPLETED_RESPONSE" | jq -r ".data.orders[$i].order_no")
                STATUS=$(echo "$COMPLETED_RESPONSE" | jq -r ".data.orders[$i].status")
                AMOUNT=$(echo "$COMPLETED_RESPONSE" | jq -r ".data.orders[$i].amount")
                UID=$(echo "$COMPLETED_RESPONSE" | jq -r ".data.orders[$i].uid")
                CREATED_AT=$(echo "$COMPLETED_RESPONSE" | jq -r ".data.orders[$i].created_at")
                
                echo -e "${GREEN}订单号: ${ORDER_NO} - 状态: ${STATUS} - 金额: ${AMOUNT} - UID: ${UID} - 创建时间: ${CREATED_AT}${NC}"
            fi
        done
    fi
else
    echo -e "${RED}查询已完成订单失败: $(echo "$COMPLETED_RESPONSE" | jq -r '.message')${NC}"
fi

# 4. 测试查询全部订单 (status: 3)
echo -e "\n${YELLOW}4. 测试查询全部订单 (status: 3)${NC}"
ALL_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 3
  }')

echo "全部订单响应: $ALL_RESPONSE"

# 检查响应状态
RESPONSE_CODE=$(echo "$ALL_RESPONSE" | jq -r '.code')
if [ "$RESPONSE_CODE" = "0" ]; then
    echo -e "${GREEN}查询全部订单成功${NC}"
    
    # 解析响应数据
    TOTAL=$(echo "$ALL_RESPONSE" | jq -r '.data.pagination.total')
    ORDERS_COUNT=$(echo "$ALL_RESPONSE" | jq -r '.data.orders | length')
    
    echo -e "${GREEN}全部订单总数: $TOTAL${NC}"
    echo -e "${GREEN}当前页订单数量: $ORDERS_COUNT${NC}"
    
    # 显示订单信息
    if [ "$ORDERS_COUNT" -gt 0 ]; then
        echo -e "\n${YELLOW}全部订单详情:${NC}"
        for i in {0..2}; do
            if [ "$i" -lt "$ORDERS_COUNT" ]; then
                ORDER_NO=$(echo "$ALL_RESPONSE" | jq -r ".data.orders[$i].order_no")
                STATUS=$(echo "$ALL_RESPONSE" | jq -r ".data.orders[$i].status")
                AMOUNT=$(echo "$ALL_RESPONSE" | jq -r ".data.orders[$i].amount")
                UID=$(echo "$ALL_RESPONSE" | jq -r ".data.orders[$i].uid")
                CREATED_AT=$(echo "$ALL_RESPONSE" | jq -r ".data.orders[$i].created_at")
                
                echo -e "${GREEN}订单号: ${ORDER_NO} - 状态: ${STATUS} - 金额: ${AMOUNT} - UID: ${UID} - 创建时间: ${CREATED_AT}${NC}"
            fi
        done
    fi
else
    echo -e "${RED}查询全部订单失败: $(echo "$ALL_RESPONSE" | jq -r '.message')${NC}"
fi

# 5. 测试订单统计
echo -e "\n${YELLOW}5. 测试订单统计${NC}"
STATS_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/stats" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "订单统计响应: $STATS_RESPONSE"

# 检查响应状态
RESPONSE_CODE=$(echo "$STATS_RESPONSE" | jq -r '.code')
if [ "$RESPONSE_CODE" = "0" ]; then
    echo -e "${GREEN}查询订单统计成功${NC}"
    
    # 解析统计数据
    TOTAL_ORDERS=$(echo "$STATS_RESPONSE" | jq -r '.data.stats.total_orders')
    PENDING_ORDERS=$(echo "$STATS_RESPONSE" | jq -r '.data.stats.pending_orders')
    SUCCESS_ORDERS=$(echo "$STATS_RESPONSE" | jq -r '.data.stats.success_orders')
    FAILED_ORDERS=$(echo "$STATS_RESPONSE" | jq -r '.data.stats.failed_orders')
    TOTAL_AMOUNT=$(echo "$STATS_RESPONSE" | jq -r '.data.stats.total_amount')
    TOTAL_PROFIT=$(echo "$STATS_RESPONSE" | jq -r '.data.stats.total_profit')
    
    echo -e "${GREEN}总订单数: $TOTAL_ORDERS${NC}"
    echo -e "${GREEN}待处理订单数: $PENDING_ORDERS${NC}"
    echo -e "${GREEN}成功订单数: $SUCCESS_ORDERS${NC}"
    echo -e "${GREEN}失败订单数: $FAILED_ORDERS${NC}"
    echo -e "${GREEN}总金额: $TOTAL_AMOUNT${NC}"
    echo -e "${GREEN}总利润: $TOTAL_PROFIT${NC}"
else
    echo -e "${RED}查询订单统计失败: $(echo "$STATS_RESPONSE" | jq -r '.message')${NC}"
fi

echo -e "\n${GREEN}=== 订单列表查询修复测试完成 ===${NC}" 