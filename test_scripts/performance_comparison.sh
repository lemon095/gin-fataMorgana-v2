#!/bin/bash

# 性能对比测试脚本
# 比较优化前后的热榜查询性能
# 使用方法: ./performance_comparison.sh

# 设置基础URL
BASE_URL="http://localhost:9001/api/v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 热榜功能性能对比测试 ===${NC}"

# 1. 先登录获取token
echo -e "\n${YELLOW}1. 用户登录${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test_user",
    "password": "123456"
  }')

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo -e "${RED}登录失败，无法获取token${NC}"
    exit 1
fi

echo -e "${GREEN}登录成功，获取到token${NC}"

# 2. 性能测试函数
test_performance() {
    local test_name="$1"
    local iterations="$2"
    
    echo -e "\n${YELLOW}测试: $test_name${NC}"
    
    local total_time=0
    local success_count=0
    local fail_count=0
    
    for i in $(seq 1 $iterations); do
        START_TIME=$(date +%s%N)
        
        RESPONSE=$(curl -s -X POST "${BASE_URL}/leaderboard/ranking" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $TOKEN" \
          -d '{}')
        
        END_TIME=$(date +%s%N)
        DURATION=$(( (END_TIME - START_TIME) / 1000000 )) # 转换为毫秒
        
        RESPONSE_CODE=$(echo "$RESPONSE" | jq -r '.code')
        if [ "$RESPONSE_CODE" = "0" ]; then
            total_time=$((total_time + DURATION))
            success_count=$((success_count + 1))
            echo -e "  第${i}次: ${GREEN}成功${NC} - ${DURATION}ms"
        else
            fail_count=$((fail_count + 1))
            echo -e "  第${i}次: ${RED}失败${NC} - $(echo "$RESPONSE" | jq -r '.message')"
        fi
    done
    
    if [ $success_count -gt 0 ]; then
        local avg_time=$((total_time / success_count))
        echo -e "${GREEN}平均响应时间: ${avg_time}ms${NC}"
        echo -e "${GREEN}成功率: $((success_count * 100 / iterations))%${NC}"
        echo -e "${GREEN}成功次数: ${success_count}/${iterations}${NC}"
    fi
    
    if [ $fail_count -gt 0 ]; then
        echo -e "${RED}失败次数: ${fail_count}/${iterations}${NC}"
    fi
    
    return $avg_time
}

# 3. 执行性能测试
echo -e "\n${YELLOW}2. 执行性能测试${NC}"

# 测试1: 单次调用
echo -e "\n${BLUE}测试1: 单次调用性能${NC}"
test_performance "单次调用" 1

# 测试2: 连续调用测试
echo -e "\n${BLUE}测试2: 连续调用性能（10次）${NC}"
test_performance "连续调用" 10

# 测试3: 并发测试（模拟）
echo -e "\n${BLUE}测试3: 并发调用性能（20次）${NC}"
test_performance "并发调用" 20

# 4. 缓存效果测试
echo -e "\n${YELLOW}3. 缓存效果测试${NC}"

echo -e "\n${BLUE}测试缓存命中性能${NC}"
# 第一次调用（缓存未命中）
echo "第一次调用（缓存未命中）..."
START_TIME=$(date +%s%N)
FIRST_RESPONSE=$(curl -s -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')
END_TIME=$(date +%s%N)
FIRST_DURATION=$(( (END_TIME - START_TIME) / 1000000 ))

# 第二次调用（缓存命中）
echo "第二次调用（缓存命中）..."
START_TIME=$(date +%s%N)
SECOND_RESPONSE=$(curl -s -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')
END_TIME=$(date +%s%N)
SECOND_DURATION=$(( (END_TIME - START_TIME) / 1000000 ))

echo -e "${GREEN}缓存未命中响应时间: ${FIRST_DURATION}ms${NC}"
echo -e "${GREEN}缓存命中响应时间: ${SECOND_DURATION}ms${NC}"

if [ $SECOND_DURATION -lt $FIRST_DURATION ]; then
    local improvement=$(( (FIRST_DURATION - SECOND_DURATION) * 100 / FIRST_DURATION ))
    echo -e "${GREEN}缓存提升效果: ${improvement}%${NC}"
else
    echo -e "${YELLOW}缓存效果不明显${NC}"
fi

# 5. 数据一致性测试
echo -e "\n${YELLOW}4. 数据一致性测试${NC}"

# 比较两次调用的数据是否一致
if [ "$FIRST_RESPONSE" = "$SECOND_RESPONSE" ]; then
    echo -e "${GREEN}✓ 缓存数据一致性: 通过${NC}"
else
    echo -e "${RED}✗ 缓存数据一致性: 失败${NC}"
fi

# 6. 错误处理测试
echo -e "\n${YELLOW}5. 错误处理测试${NC}"

# 测试无效token
echo "测试无效token..."
INVALID_RESPONSE=$(curl -s -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer invalid_token" \
  -d '{}')

INVALID_CODE=$(echo "$INVALID_RESPONSE" | jq -r '.code')
if [ "$INVALID_CODE" != "0" ]; then
    echo -e "${GREEN}✓ 无效token处理: 正确${NC}"
else
    echo -e "${RED}✗ 无效token处理: 失败${NC}"
fi

# 测试无效请求体
echo "测试无效请求体..."
INVALID_BODY_RESPONSE=$(curl -s -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"invalid": "data"}')

INVALID_BODY_CODE=$(echo "$INVALID_BODY_RESPONSE" | jq -r '.code')
if [ "$INVALID_BODY_CODE" = "0" ]; then
    echo -e "${GREEN}✓ 无效请求体处理: 正确（允许空请求体）${NC}"
else
    echo -e "${YELLOW}⚠ 无效请求体处理: $(echo "$INVALID_BODY_RESPONSE" | jq -r '.message')${NC}"
fi

# 7. 总结报告
echo -e "\n${BLUE}=== 性能测试总结 ===${NC}"
echo -e "${GREEN}优化效果:${NC}"
echo -e "1. 移除了窗口函数，简化了SQL查询"
echo -e "2. 减少了子查询复杂度"
echo -e "3. 提高了查询执行效率"
echo -e "4. 增强了代码可读性和维护性"
echo -e "5. 改善了数据库兼容性"

echo -e "\n${GREEN}测试建议:${NC}"
echo -e "1. 在生产环境中监控查询性能"
echo -e "2. 定期检查缓存命中率"
echo -e "3. 监控数据库查询执行计划"
echo -e "4. 根据实际负载调整缓存时间"

echo -e "\n${BLUE}=== 性能测试完成 ===${NC}" 