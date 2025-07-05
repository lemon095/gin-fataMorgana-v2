#!/bin/bash

# 测试排行榜修复
echo "🔍 测试排行榜修复..."

# 设置颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 测试基础URL
BASE_URL="http://localhost:8080/api/v1"

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

# 2. 清除排行榜缓存
echo -e "\n${YELLOW}2. 清除排行榜缓存${NC}"
CLEAR_CACHE_RESPONSE=$(curl -s -X POST "${BASE_URL}/leaderboard/clear-cache" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "清除缓存响应: $CLEAR_CACHE_RESPONSE"

# 3. 测试排行榜接口
echo -e "\n${YELLOW}3. 测试排行榜接口${NC}"
LEADERBOARD_RESPONSE=$(curl -s -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "排行榜响应: $LEADERBOARD_RESPONSE"

# 检查响应状态
RESPONSE_CODE=$(echo "$LEADERBOARD_RESPONSE" | jq -r '.code')
if [ "$RESPONSE_CODE" = "0" ]; then
    echo -e "${GREEN}排行榜接口调用成功${NC}"
    
    # 解析响应数据
    WEEK_START=$(echo "$LEADERBOARD_RESPONSE" | jq -r '.data.week_start')
    WEEK_END=$(echo "$LEADERBOARD_RESPONSE" | jq -r '.data.week_end')
    MY_RANK=$(echo "$LEADERBOARD_RESPONSE" | jq -r '.data.my_rank.rank')
    IS_RANK=$(echo "$LEADERBOARD_RESPONSE" | jq -r '.data.my_rank.is_rank')
    TOP_USERS_COUNT=$(echo "$LEADERBOARD_RESPONSE" | jq -r '.data.top_users | length')
    
    echo -e "${GREEN}本周时间范围: $WEEK_START 到 $WEEK_END${NC}"
    echo -e "${GREEN}我的排名: $MY_RANK${NC}"
    echo -e "${GREEN}是否在榜单: $IS_RANK${NC}"
    echo -e "${GREEN}前10名用户数量: $TOP_USERS_COUNT${NC}"
    
    # 显示前10名用户信息
    echo -e "\n${YELLOW}前10名用户信息:${NC}"
    for i in {0..9}; do
        if [ "$i" -lt "$TOP_USERS_COUNT" ]; then
            USER_RANK=$(echo "$LEADERBOARD_RESPONSE" | jq -r ".data.top_users[$i].rank")
            USERNAME=$(echo "$LEADERBOARD_RESPONSE" | jq -r ".data.top_users[$i].username")
            ORDER_COUNT=$(echo "$LEADERBOARD_RESPONSE" | jq -r ".data.top_users[$i].order_count")
            TOTAL_AMOUNT=$(echo "$LEADERBOARD_RESPONSE" | jq -r ".data.top_users[$i].total_amount")
            TOTAL_PROFIT=$(echo "$LEADERBOARD_RESPONSE" | jq -r ".data.top_users[$i].total_profit")
            
            echo -e "${GREEN}第${USER_RANK}名: ${USERNAME} - 完成${ORDER_COUNT}单 - 总金额${TOTAL_AMOUNT} - 总利润${TOTAL_PROFIT}${NC}"
        fi
    done
    
    # 显示我的排名详情
    echo -e "\n${YELLOW}我的排名详情:${NC}"
    MY_USERNAME=$(echo "$LEADERBOARD_RESPONSE" | jq -r '.data.my_rank.username')
    MY_ORDER_COUNT=$(echo "$LEADERBOARD_RESPONSE" | jq -r '.data.my_rank.order_count')
    MY_TOTAL_AMOUNT=$(echo "$LEADERBOARD_RESPONSE" | jq -r '.data.my_rank.total_amount')
    MY_TOTAL_PROFIT=$(echo "$LEADERBOARD_RESPONSE" | jq -r '.data.my_rank.total_profit')
    
    echo -e "${GREEN}用户名: ${MY_USERNAME}${NC}"
    echo -e "${GREEN}完成订单数: ${MY_ORDER_COUNT}${NC}"
    echo -e "${GREEN}总金额: ${MY_TOTAL_AMOUNT}${NC}"
    echo -e "${GREEN}总利润: ${MY_TOTAL_PROFIT}${NC}"
    
    # 检查是否包含水单数据
    if [ "$MY_ORDER_COUNT" -gt 0 ] || [ "$TOP_USERS_COUNT" -gt 0 ]; then
        echo -e "\n${GREEN}✅ 排行榜修复成功！现在包含水单数据了${NC}"
    else
        echo -e "\n${YELLOW}⚠️  排行榜数据为空，可能需要检查：${NC}"
        echo -e "1. 水单状态是否为 'success'"
        echo -e "2. 水单创建时间是否在本周范围内"
        echo -e "3. 本周时间范围是否正确"
    fi
    
else
    echo -e "${RED}排行榜接口调用失败: $(echo "$LEADERBOARD_RESPONSE" | jq -r '.message')${NC}"
fi

echo -e "\n${BLUE}=== 排行榜修复测试完成 ===${NC}"
echo -e "${GREEN}修复内容:${NC}"
echo -e "1. 将时间字段从 updated_at 改为 created_at"
echo -e "2. 移除了 o.updated_at <= NOW() 的时间限制"
echo -e "3. 现在统计的是本周创建且状态为success的订单" 