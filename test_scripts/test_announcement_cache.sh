#!/bin/bash

# 公告缓存测试脚本
BASE_URL="http://localhost:9001/api/v1"

echo "=== 公告缓存测试 ==="

# 第一次请求（应该从数据库获取）
echo "1. 第一次请求（从数据库获取）"
START_TIME=$(date +%s%N)
RESPONSE1=$(curl -s -X POST "${BASE_URL}/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 1,
    "page_size": 10
  }')
END_TIME=$(date +%s%N)
FIRST_REQUEST_TIME=$((($END_TIME - $START_TIME) / 1000000))

echo "响应: $RESPONSE1"
echo "请求耗时: ${FIRST_REQUEST_TIME}ms"

# 检查响应状态
if echo "$RESPONSE1" | jq -e '.code == 0' > /dev/null 2>&1; then
    echo "✅ 第一次请求成功"
else
    echo "❌ 第一次请求失败"
    exit 1
fi

echo -e "\n"

# 等待1秒
echo "等待1秒..."
sleep 1

# 第二次请求（应该从缓存获取）
echo "2. 第二次请求（从缓存获取）"
START_TIME=$(date +%s%N)
RESPONSE2=$(curl -s -X POST "${BASE_URL}/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 1,
    "page_size": 10
  }')
END_TIME=$(date +%s%N)
SECOND_REQUEST_TIME=$((($END_TIME - $START_TIME) / 1000000))

echo "响应: $RESPONSE2"
echo "请求耗时: ${SECOND_REQUEST_TIME}ms"

# 检查响应状态
if echo "$RESPONSE2" | jq -e '.code == 0' > /dev/null 2>&1; then
    echo "✅ 第二次请求成功"
else
    echo "❌ 第二次请求失败"
    exit 1
fi

echo -e "\n"

# 比较两次请求的耗时
echo "3. 性能对比"
echo "第一次请求耗时: ${FIRST_REQUEST_TIME}ms"
echo "第二次请求耗时: ${SECOND_REQUEST_TIME}ms"

if [ $SECOND_REQUEST_TIME -lt $FIRST_REQUEST_TIME ]; then
    IMPROVEMENT=$((FIRST_REQUEST_TIME - SECOND_REQUEST_TIME))
    echo "✅ 缓存生效，性能提升: ${IMPROVEMENT}ms"
else
    echo "⚠️  缓存可能未生效或性能提升不明显"
fi

echo -e "\n"

# 比较两次请求的响应数据
echo "4. 数据一致性检查"
if [ "$RESPONSE1" = "$RESPONSE2" ]; then
    echo "✅ 两次请求返回的数据完全一致"
else
    echo "❌ 两次请求返回的数据不一致"
fi

echo -e "\n"

# 测试不同参数的缓存
echo "5. 测试不同参数的缓存"
RESPONSE3=$(curl -s -X POST "${BASE_URL}/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 2,
    "page_size": 5
  }')

echo "不同参数请求响应: $RESPONSE3"

if echo "$RESPONSE3" | jq -e '.code == 0' > /dev/null 2>&1; then
    echo "✅ 不同参数请求成功"
else
    echo "❌ 不同参数请求失败"
fi

echo -e "\n"

# 检查Redis缓存状态
echo "6. 检查Redis缓存状态"
if command -v redis-cli >/dev/null 2>&1; then
    CACHE_KEYS=$(redis-cli keys "announcement:list:*" 2>/dev/null | wc -l)
    echo "缓存键数量: $CACHE_KEYS"
    
    if [ $CACHE_KEYS -gt 0 ]; then
        echo "✅ Redis中存在公告缓存"
    else
        echo "⚠️  Redis中未找到公告缓存"
    fi
else
    echo "⚠️  redis-cli未安装，无法检查缓存状态"
fi

echo -e "\n=== 缓存测试完成 ==="
echo "💡 提示: 缓存时间为1分钟，相同参数的请求会直接返回缓存数据"
echo "💡 提示: 不同参数的请求会使用不同的缓存键" 