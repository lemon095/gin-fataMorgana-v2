#!/bin/bash

# 测试排行榜实时查询功能
echo "🔍 测试排行榜实时查询功能..."

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

# 2. 测试排行榜接口（第一次调用）
echo -e "\n${YELLOW}2. 第一次调用排行榜接口${NC}"
FIRST_RESPONSE=$(curl -s -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "第一次响应: $FIRST_RESPONSE"

# 检查响应状态
RESPONSE_CODE=$(echo "$FIRST_RESPONSE" | jq -r '.code')
if [ "$RESPONSE_CODE" = "0" ]; then
    echo -e "${GREEN}第一次调用成功${NC}"
    
    # 解析响应数据
    WEEK_START=$(echo "$FIRST_RESPONSE" | jq -r '.data.week_start')
    WEEK_END=$(echo "$FIRST_RESPONSE" | jq -r '.data.week_end')
    MY_RANK=$(echo "$FIRST_RESPONSE" | jq -r '.data.my_rank.rank')
    IS_RANK=$(echo "$FIRST_RESPONSE" | jq -r '.data.my_rank.is_rank')
    TOP_USERS_COUNT=$(echo "$FIRST_RESPONSE" | jq -r '.data.top_users | length')
    
    echo -e "${GREEN}本周时间范围: $WEEK_START 到 $WEEK_END${NC}"
    echo -e "${GREEN}我的排名: $MY_RANK${NC}"
    echo -e "${GREEN}是否在榜单: $IS_RANK${NC}"
    echo -e "${GREEN}前10名用户数量: $TOP_USERS_COUNT${NC}"
    
    # 检查是否还有cache_expire字段
    CACHE_EXPIRE=$(echo "$FIRST_RESPONSE" | jq -r '.data.cache_expire')
    if [ "$CACHE_EXPIRE" = "null" ]; then
        echo -e "${GREEN}✓ 已移除缓存过期时间字段${NC}"
    else
        echo -e "${RED}✗ 缓存过期时间字段仍然存在${NC}"
    fi
    
else
    echo -e "${RED}第一次调用失败: $(echo "$FIRST_RESPONSE" | jq -r '.message')${NC}"
    exit 1
fi

# 3. 等待1秒后再次调用（测试实时性）
echo -e "\n${YELLOW}3. 等待1秒后再次调用排行榜接口${NC}"
sleep 1

SECOND_RESPONSE=$(curl -s -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "第二次响应: $SECOND_RESPONSE"

# 检查响应状态
RESPONSE_CODE=$(echo "$SECOND_RESPONSE" | jq -r '.code')
if [ "$RESPONSE_CODE" = "0" ]; then
    echo -e "${GREEN}第二次调用成功${NC}"
    
    # 比较两次调用的响应时间
    FIRST_TIME=$(echo "$FIRST_RESPONSE" | jq -r '.data.week_start')
    SECOND_TIME=$(echo "$SECOND_RESPONSE" | jq -r '.data.week_start')
    
    if [ "$FIRST_TIME" = "$SECOND_TIME" ]; then
        echo -e "${GREEN}✓ 两次调用返回相同的时间范围（正常）${NC}"
    else
        echo -e "${YELLOW}⚠ 两次调用返回不同的时间范围${NC}"
    fi
    
else
    echo -e "${RED}第二次调用失败: $(echo "$SECOND_RESPONSE" | jq -r '.message')${NC}"
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
    REQUIRED_FIELDS=("week_start" "week_end" "my_rank" "top_users")
    for field in "${REQUIRED_FIELDS[@]}"; do
        if echo "$FIRST_RESPONSE" | jq -e ".data.$field" > /dev/null; then
            echo -e "${GREEN}✓ 字段 $field 存在${NC}"
        else
            echo -e "${RED}✗ 字段 $field 缺失${NC}"
        fi
    done
    
    # 检查my_rank字段
    MY_RANK_FIELDS=("id" "uid" "username" "completed_at" "order_count" "total_amount" "total_profit" "rank" "is_rank")
    for field in "${MY_RANK_FIELDS[@]}"; do
        if echo "$FIRST_RESPONSE" | jq -e ".data.my_rank.$field" > /dev/null; then
            echo -e "${GREEN}✓ my_rank.$field 存在${NC}"
        else
            echo -e "${RED}✗ my_rank.$field 缺失${NC}"
        fi
    done
    
    # 确认cache_expire字段已移除
    if echo "$FIRST_RESPONSE" | jq -e ".data.cache_expire" > /dev/null; then
        echo -e "${RED}✗ cache_expire 字段仍然存在${NC}"
    else
        echo -e "${GREEN}✓ cache_expire 字段已成功移除${NC}"
    fi
fi

# 6. 测试清除缓存接口是否已移除
echo -e "\n${YELLOW}6. 测试清除缓存接口是否已移除${NC}"
CLEAR_CACHE_RESPONSE=$(curl -s -X POST "${BASE_URL}/leaderboard/clear-cache" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "清除缓存响应: $CLEAR_CACHE_RESPONSE"

# 检查是否返回404或其他错误
if echo "$CLEAR_CACHE_RESPONSE" | grep -q "404\|Not Found\|Method Not Allowed"; then
    echo -e "${GREEN}✓ 清除缓存接口已成功移除${NC}"
else
    echo -e "${YELLOW}⚠ 清除缓存接口可能仍然存在${NC}"
fi

echo -e "\n${BLUE}=== 排行榜实时查询功能测试完成 ===${NC}"
echo -e "${GREEN}主要改进:${NC}"
echo -e "1. 移除了Redis缓存，改为实时查询数据库"
echo -e "2. 移除了cache_expire字段"
echo -e "3. 移除了清除缓存接口"
echo -e "4. 每次请求都从数据库获取最新数据"
echo -e "5. 确保数据实时性和准确性" 