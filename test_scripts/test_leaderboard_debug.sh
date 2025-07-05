#!/bin/bash

# 调试排行榜数据
echo "🔍 调试排行榜数据..."

# 设置颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试基础URL
BASE_URL="http://localhost:8080/api/v1"

# 1. 检查排行榜接口
echo -e "\n${YELLOW}1. 检查排行榜接口${NC}"
LEADERBOARD_RESPONSE=$(curl -s -X GET "${BASE_URL}/leaderboard/ranking" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE")

echo "排行榜响应: $LEADERBOARD_RESPONSE"

# 2. 检查订单状态分布
echo -e "\n${YELLOW}2. 检查订单状态分布${NC}"
echo "需要检查数据库中订单的状态分布..."

# 3. 检查本周时间范围
echo -e "\n${YELLOW}3. 检查本周时间范围${NC}"
echo "排行榜只统计本周（周一到周日）的数据"

# 4. 检查Redis缓存
echo -e "\n${YELLOW}4. 检查Redis缓存${NC}"
echo "排行榜数据有5分钟缓存，可能需要清除缓存"

# 5. 手动清除缓存
echo -e "\n${YELLOW}5. 手动清除排行榜缓存${NC}"
echo "可以通过以下方式清除缓存："
echo "1. 重启应用"
echo "2. 等待5分钟缓存过期"
echo "3. 调用清除缓存接口（如果有的话）"

echo -e "\n${GREEN}🎉 调试完成！${NC}"
echo -e "\n${YELLOW}可能的问题：${NC}"
echo "1. 水单状态不是 'success'"
echo "2. 水单创建时间不在本周范围内"
echo "3. Redis缓存了旧数据"
echo "4. 数据库中确实只有3个用户有完成订单" 