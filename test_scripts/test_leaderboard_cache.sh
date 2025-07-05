#!/bin/bash

# 测试热榜缓存功能
echo "=== 测试热榜缓存功能 ==="

# 设置基础URL
BASE_URL="http://localhost:9001/api/v1"

# 测试用户UID
TEST_UID="1000001"

echo "1. 测试热榜接口（从缓存读取）"
echo "请求URL: $BASE_URL/leaderboard"
echo "用户UID: $TEST_UID"

# 发送请求
response=$(curl -s -X POST "$BASE_URL/leaderboard" \
  -H "Content-Type: application/json" \
  -d "{\"uid\": \"$TEST_UID\"}")

echo "响应:"
echo "$response" | jq '.'

echo ""
echo "2. 检查缓存数据"
echo "请求Redis缓存..."

# 检查Redis中的缓存数据
redis_key="leaderboard:weekly:$(date +%Y-%m-%d)"
echo "缓存键: $redis_key"

# 如果有redis-cli，可以检查缓存
if command -v redis-cli &> /dev/null; then
    echo "Redis缓存内容:"
    redis-cli get "$redis_key" | jq '.' 2>/dev/null || echo "缓存不存在或格式错误"
else
    echo "redis-cli 未安装，无法直接检查缓存"
fi

echo ""
echo "3. 手动触发缓存更新"
echo "请求URL: $BASE_URL/cron/update-leaderboard-cache"

# 手动触发缓存更新
cache_response=$(curl -s -X POST "$BASE_URL/cron/update-leaderboard-cache" \
  -H "Content-Type: application/json")

echo "缓存更新响应:"
echo "$cache_response" | jq '.'

echo ""
echo "4. 再次测试热榜接口（验证缓存更新）"
echo "请求URL: $BASE_URL/leaderboard"

# 再次发送请求
response2=$(curl -s -X POST "$BASE_URL/leaderboard" \
  -H "Content-Type: application/json" \
  -d "{\"uid\": \"$TEST_UID\"}")

echo "响应:"
echo "$response2" | jq '.'

echo ""
echo "=== 测试完成 ===" 