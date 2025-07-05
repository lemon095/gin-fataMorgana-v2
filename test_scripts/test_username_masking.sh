#!/bin/bash

# 测试用户名脱敏功能
# 使用方法: ./test_username_masking.sh

# 设置基础URL
BASE_URL="http://localhost:9001/api/v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 测试用户名脱敏功能（优化版本）===${NC}"

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

# 2. 测试热榜接口，查看用户名脱敏效果
echo -e "\n${YELLOW}2. 测试热榜接口，查看用户名脱敏效果${NC}"
LEADERBOARD_RESPONSE=$(curl -s -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "热榜响应: $LEADERBOARD_RESPONSE"

# 检查响应状态
RESPONSE_CODE=$(echo "$LEADERBOARD_RESPONSE" | jq -r '.code')
if [ "$RESPONSE_CODE" = "0" ]; then
    echo -e "${GREEN}热榜接口调用成功${NC}"
    
    # 显示我的排名信息中的用户名
    MY_USERNAME=$(echo "$LEADERBOARD_RESPONSE" | jq -r '.data.my_rank.username')
    echo -e "${GREEN}我的用户名（脱敏后）: ${MY_USERNAME}${NC}"
    
    # 显示前5名用户的用户名
    echo -e "\n${YELLOW}前5名用户用户名（脱敏后）:${NC}"
    TOP_USERS_COUNT=$(echo "$LEADERBOARD_RESPONSE" | jq -r '.data.top_users | length')
    
    for i in {0..4}; do
        if [ "$i" -lt "$TOP_USERS_COUNT" ]; then
            USER_RANK=$(echo "$LEADERBOARD_RESPONSE" | jq -r ".data.top_users[$i].rank")
            USERNAME=$(echo "$LEADERBOARD_RESPONSE" | jq -r ".data.top_users[$i].username")
            echo -e "${GREEN}第${USER_RANK}名: ${USERNAME}${NC}"
        fi
    done
    
else
    echo -e "${RED}热榜接口调用失败: $(echo "$LEADERBOARD_RESPONSE" | jq -r '.message')${NC}"
fi

# 3. 脱敏效果说明
echo -e "\n${YELLOW}3. 脱敏效果说明${NC}"
echo -e "${GREEN}统一的脱敏规则:${NC}"
echo -e "  - 用户名长度 = 1: 不脱敏，直接显示"
echo -e "  - 用户名长度 ≥ 2: 统一格式：首位 + ** + 末位"
echo -e ""
echo -e "${GREEN}示例:${NC}"
echo -e "  - '张三' → '张**三'"
echo -e "  - '张三丰' → '张**丰'"
echo -e "  - '张三丰李' → '张**李'"
echo -e "  - '张三丰李四' → '张**四'"
echo -e "  - '张三丰李四王' → '张**王'"
echo -e "  - 'test_user_123' → 't**3'"

# 4. 对比优化前后的效果
echo -e "\n${YELLOW}4. 优化前后对比${NC}"
echo -e "${GREEN}优化前:${NC}"
echo -e "  - '张三' → '张*三'（2位特殊处理）"
echo -e "  - '张三丰' → '张*丰'（3位特殊处理）"
echo -e "  - '张三丰李四王' → '张****王'（复杂处理）"
echo -e "  - 'test_user_123' → 't*********3'（复杂处理）"
echo -e ""
echo -e "${GREEN}统一后:${NC}"
echo -e "  - '张三' → '张**三'"
echo -e "  - '张三丰' → '张**丰'"
echo -e "  - '张三丰李四王' → '张**王'"
echo -e "  - 'test_user_123' → 't**3'"
echo -e ""
echo -e "${GREEN}统一效果:${NC}"
echo -e "  - 所有用户名使用相同的脱敏规则"
echo -e "  - 脱敏长度统一为3个字符"
echo -e "  - 代码逻辑更简单，易于维护"
echo -e "  - 用户体验更一致"
echo -e "  - 仍然保护了用户隐私"

echo -e "\n${BLUE}=== 用户名脱敏测试完成 ===${NC}" 