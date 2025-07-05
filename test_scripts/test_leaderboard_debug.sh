#!/bin/bash

# 调试排行榜查询
echo "🔍 调试排行榜查询..."

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

# 2. 调用排行榜接口（这会触发详细的日志输出）
echo -e "\n${YELLOW}2. 调用排行榜接口（查看容器日志）${NC}"
echo -e "${BLUE}请查看容器日志，会显示详细的查询过程${NC}"

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
    
    # 显示前3名用户信息
    echo -e "\n${YELLOW}前3名用户信息:${NC}"
    for i in {0..2}; do
        if [ "$i" -lt "$TOP_USERS_COUNT" ]; then
            USER_RANK=$(echo "$LEADERBOARD_RESPONSE" | jq -r ".data.top_users[$i].rank")
            USERNAME=$(echo "$LEADERBOARD_RESPONSE" | jq -r ".data.top_users[$i].username")
            ORDER_COUNT=$(echo "$LEADERBOARD_RESPONSE" | jq -r ".data.top_users[$i].order_count")
            TOTAL_AMOUNT=$(echo "$LEADERBOARD_RESPONSE" | jq -r ".data.top_users[$i].total_amount")
            TOTAL_PROFIT=$(echo "$LEADERBOARD_RESPONSE" | jq -r ".data.top_users[$i].total_profit")
            
            echo -e "${GREEN}第${USER_RANK}名: ${USERNAME} - 完成${ORDER_COUNT}单 - 总金额${TOTAL_AMOUNT} - 总利润${TOTAL_PROFIT}${NC}"
        fi
    done
    
else
    echo -e "${RED}排行榜接口调用失败: $(echo "$LEADERBOARD_RESPONSE" | jq -r '.message')${NC}"
fi

echo -e "\n${BLUE}=== 排行榜调试完成 ===${NC}"
echo -e "${YELLOW}请查看容器日志，会显示以下信息：${NC}"
echo -e "1. 查询的时间范围"
echo -e "2. SQL查询语句和参数"
echo -e "3. 查询结果的数量和详细信息"
echo -e "4. 用户排名计算过程"
echo -e "5. 任何错误或警告信息" 