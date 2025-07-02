#!/bin/bash

# 公告接口测试脚本
BASE_URL="http://localhost:9001/api/v1"

echo "=== 公告接口测试 ==="

# 测试获取公告列表
echo "1. 测试获取公告列表（第一页）"
RESPONSE1=$(curl -s -X POST "${BASE_URL}/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 1,
    "page_size": 10
  }')

echo "响应: $RESPONSE1"

# 检查响应状态
if echo "$RESPONSE1" | jq -e '.code == 0' > /dev/null 2>&1; then
    echo "✅ 获取公告列表成功"
else
    echo "❌ 获取公告列表失败"
fi

echo -e "\n"

# 测试获取公告列表（第二页）
echo "2. 测试获取公告列表（第二页）"
RESPONSE2=$(curl -s -X POST "${BASE_URL}/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 2,
    "page_size": 5
  }')

echo "响应: $RESPONSE2"

# 检查响应状态
if echo "$RESPONSE2" | jq -e '.code == 0' > /dev/null 2>&1; then
    echo "✅ 获取第二页公告列表成功"
else
    echo "❌ 获取第二页公告列表失败"
fi

echo -e "\n"

# 测试参数验证
echo "3. 测试参数验证（无效页码）"
RESPONSE3=$(curl -s -X POST "${BASE_URL}/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 0,
    "page_size": 10
  }')

echo "响应: $RESPONSE3"

# 检查响应状态
if echo "$RESPONSE3" | jq -e '.code != 0' > /dev/null 2>&1; then
    echo "✅ 参数验证正确"
else
    echo "❌ 参数验证失败"
fi

echo -e "\n"

# 测试默认参数
echo "4. 测试默认参数（不传参数）"
RESPONSE4=$(curl -s -X POST "${BASE_URL}/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{}')

echo "响应: $RESPONSE4"

# 检查响应状态
if echo "$RESPONSE4" | jq -e '.code == 0' > /dev/null 2>&1; then
    echo "✅ 默认参数处理正确"
else
    echo "❌ 默认参数处理失败"
fi

echo -e "\n"

# 统计返回的公告数量
echo "5. 统计返回的公告数量"
ANNOUNCEMENT_COUNT=$(echo "$RESPONSE1" | jq '.data.announcements | length' 2>/dev/null || echo "0")
TOTAL_COUNT=$(echo "$RESPONSE1" | jq '.data.pagination.total' 2>/dev/null || echo "0")

echo "当前页公告数量: $ANNOUNCEMENT_COUNT"
echo "总公告数量: $TOTAL_COUNT"

echo -e "\n"

# 检查公告数据结构
echo "6. 检查公告数据结构"
if echo "$RESPONSE1" | jq -e '.data.announcements[0] | has("id", "title", "content", "tag", "created_at", "banners")' > /dev/null 2>&1; then
    echo "✅ 公告数据结构正确"
else
    echo "❌ 公告数据结构不正确"
fi

# 检查banners数据结构
echo "7. 检查banners数据结构"
if echo "$RESPONSE1" | jq -e '.data.announcements[0].banners[0] | type == "string"' > /dev/null 2>&1; then
    echo "✅ banners数据结构正确（字符串数组）"
else
    echo "❌ banners数据结构不正确"
fi

echo -e "\n=== 测试完成 ==="
echo "💡 提示: 如果数据库中没有公告数据，返回的列表可能为空" 