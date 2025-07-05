#!/bin/bash

# 测试优化后的热榜功能（不使用窗口函数）
# 使用方法: ./test_leaderboard_optimized.sh

# 设置基础URL
BASE_URL="http://localhost:9001/api/v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 测试优化后的热榜功能（不使用窗口函数）===${NC}"

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

# 2. 测试热榜接口
echo -e "\n${YELLOW}2. 测试热榜接口${NC}"
LEADERBOARD_RESPONSE=$(curl -s -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "热榜响应: $LEADERBOARD_RESPONSE"

# 检查响应状态
RESPONSE_CODE=$(echo "$LEADERBOARD_RESPONSE" | jq -r '.code')
if [ "$RESPONSE_CODE" = "0" ]; then
    echo -e "${GREEN}热榜接口调用成功${NC}"
    
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
    
    # 显示前5名用户信息
    echo -e "\n${YELLOW}前5名用户信息:${NC}"
    for i in {0..4}; do
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
    
else
    echo -e "${RED}热榜接口调用失败: $(echo "$LEADERBOARD_RESPONSE" | jq -r '.message')${NC}"
fi

# 3. 测试缓存功能（连续调用两次）
echo -e "\n${YELLOW}3. 测试缓存功能${NC}"
echo "第一次调用..."
FIRST_CALL=$(curl -s -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "第二次调用..."
SECOND_CALL=$(curl -s -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

# 比较两次调用的响应时间（简单比较）
if [ "$FIRST_CALL" = "$SECOND_CALL" ]; then
    echo -e "${GREEN}缓存功能正常，两次调用响应一致${NC}"
else
    echo -e "${YELLOW}缓存功能可能有问题，两次调用响应不一致${NC}"
fi

# 4. 测试性能（连续调用多次）
echo -e "\n${YELLOW}4. 测试性能（连续调用5次）${NC}"
for i in {1..5}; do
    echo "第${i}次调用..."
    START_TIME=$(date +%s%N)
    
    RESPONSE=$(curl -s -X POST "${BASE_URL}/leaderboard/ranking" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{}')
    
    END_TIME=$(date +%s%N)
    DURATION=$(( (END_TIME - START_TIME) / 1000000 )) # 转换为毫秒
    
    RESPONSE_CODE=$(echo "$RESPONSE" | jq -r '.code')
    if [ "$RESPONSE_CODE" = "0" ]; then
        echo -e "${GREEN}第${i}次调用成功，耗时: ${DURATION}ms${NC}"
    else
        echo -e "${RED}第${i}次调用失败${NC}"
    fi
done

# 5. 验证数据结构
echo -e "\n${YELLOW}5. 验证数据结构${NC}"
if [ "$RESPONSE_CODE" = "0" ]; then
    # 检查必要字段是否存在
    REQUIRED_FIELDS=("week_start" "week_end" "my_rank" "top_users" "cache_expire")
    for field in "${REQUIRED_FIELDS[@]}"; do
        if echo "$LEADERBOARD_RESPONSE" | jq -e ".data.$field" > /dev/null; then
            echo -e "${GREEN}✓ 字段 $field 存在${NC}"
        else
            echo -e "${RED}✗ 字段 $field 缺失${NC}"
        fi
    done
    
    # 检查my_rank字段
    MY_RANK_FIELDS=("id" "uid" "username" "completed_at" "order_count" "total_amount" "total_profit" "rank" "is_rank")
    for field in "${MY_RANK_FIELDS[@]}"; do
        if echo "$LEADERBOARD_RESPONSE" | jq -e ".data.my_rank.$field" > /dev/null; then
            echo -e "${GREEN}✓ my_rank.$field 存在${NC}"
        else
            echo -e "${RED}✗ my_rank.$field 缺失${NC}"
        fi
    done
fi

echo -e "\n${BLUE}=== 优化后的热榜功能测试完成 ===${NC}"
echo -e "${GREEN}主要优化点:${NC}"
echo -e "1. 移除了窗口函数，使用简单的GROUP BY和ORDER BY"
echo -e "2. 简化了排名计算逻辑"
echo -e "3. 优化了缓存处理"
echo -e "4. 改进了错误处理" 