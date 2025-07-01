#!/bin/bash

# 热榜接口测试脚本
BASE_URL="http://localhost:9001/api/v1"

echo "=== 任务热榜接口测试 ==="

# 测试获取热榜数据（用户ID: 1001 - 在榜上）
echo "1. 测试获取热榜数据（用户ID: 1001 - 在榜上）"
curl -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1001
  }' | jq '.'

echo -e "\n"

# 测试获取热榜数据（用户ID: 9999 - 不在榜上）
echo "2. 测试获取热榜数据（用户ID: 9999 - 不在榜上）"
curl -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 9999
  }' | jq '.'

echo -e "\n"

# 测试无效用户ID
echo "3. 测试无效用户ID"
curl -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 0
  }' | jq '.'

echo -e "\n"

# 测试缺少用户ID
echo "4. 测试缺少用户ID"
curl -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -d '{}' | jq '.'

echo -e "\n"

# 测试无效JSON
echo "5. 测试无效JSON"
curl -X POST "${BASE_URL}/leaderboard/ranking" \
  -H "Content-Type: application/json" \
  -d '{invalid json}' | jq '.'

echo -e "\n=== 测试完成 ===" 