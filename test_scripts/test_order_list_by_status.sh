#!/bin/bash

# 测试订单列表按状态类型查询
# 使用方法: ./test_order_list_by_status.sh

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== 测试订单列表按状态类型查询 ===${NC}"

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
    
    # 显示前3个订单信息
    if [ "$ORDERS_COUNT" -gt 0 ]; then
        echo -e "\n${YELLOW}前3个进行中订单:${NC}"
        for i in {0..2}; do
            if [ "$i" -lt "$ORDERS_COUNT" ]; then
                ORDER_NO=$(echo "$IN_PROGRESS_RESPONSE" | jq -r ".data.orders[$i].order_no")
                STATUS=$(echo "$IN_PROGRESS_RESPONSE" | jq -r ".data.orders[$i].status")
                AMOUNT=$(echo "$IN_PROGRESS_RESPONSE" | jq -r ".data.orders[$i].amount")
                CREATED_AT=$(echo "$IN_PROGRESS_RESPONSE" | jq -r ".data.orders[$i].created_at")
                
                echo -e "${GREEN}订单号: ${ORDER_NO} - 状态: ${STATUS} - 金额: ${AMOUNT} - 创建时间: ${CREATED_AT}${NC}"
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
    
    # 显示前3个订单信息
    if [ "$ORDERS_COUNT" -gt 0 ]; then
        echo -e "\n${YELLOW}前3个已完成订单:${NC}"
        for i in {0..2}; do
            if [ "$i" -lt "$ORDERS_COUNT" ]; then
                ORDER_NO=$(echo "$COMPLETED_RESPONSE" | jq -r ".data.orders[$i].order_no")
                STATUS=$(echo "$COMPLETED_RESPONSE" | jq -r ".data.orders[$i].status")
                AMOUNT=$(echo "$COMPLETED_RESPONSE" | jq -r ".data.orders[$i].amount")
                CREATED_AT=$(echo "$COMPLETED_RESPONSE" | jq -r ".data.orders[$i].created_at")
                
                echo -e "${GREEN}订单号: ${ORDER_NO} - 状态: ${STATUS} - 金额: ${AMOUNT} - 创建时间: ${CREATED_AT}${NC}"
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
    
    # 显示前3个订单信息
    if [ "$ORDERS_COUNT" -gt 0 ]; then
        echo -e "\n${YELLOW}前3个全部订单:${NC}"
        for i in {0..2}; do
            if [ "$i" -lt "$ORDERS_COUNT" ]; then
                ORDER_NO=$(echo "$ALL_RESPONSE" | jq -r ".data.orders[$i].order_no")
                STATUS=$(echo "$ALL_RESPONSE" | jq -r ".data.orders[$i].status")
                AMOUNT=$(echo "$ALL_RESPONSE" | jq -r ".data.orders[$i].amount")
                CREATED_AT=$(echo "$ALL_RESPONSE" | jq -r ".data.orders[$i].created_at")
                
                echo -e "${GREEN}订单号: ${ORDER_NO} - 状态: ${STATUS} - 金额: ${AMOUNT} - 创建时间: ${CREATED_AT}${NC}"
            fi
        done
    fi
else
    echo -e "${RED}查询全部订单失败: $(echo "$ALL_RESPONSE" | jq -r '.message')${NC}"
fi

# 5. 测试无效状态类型
echo -e "\n${YELLOW}5. 测试无效状态类型 (status: 4)${NC}"
INVALID_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 4
  }')

echo "无效状态类型响应: $INVALID_RESPONSE"

# 检查响应状态
RESPONSE_CODE=$(echo "$INVALID_RESPONSE" | jq -r '.code')
if [ "$RESPONSE_CODE" != "0" ]; then
    echo -e "${GREEN}无效状态类型验证成功，正确返回错误信息${NC}"
else
    echo -e "${RED}无效状态类型验证失败，应该返回错误信息${NC}"
fi

echo -e "\n${GREEN}=== 订单列表按状态类型查询测试完成 ===${NC}" 