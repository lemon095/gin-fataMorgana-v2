#!/bin/bash

# 公告关联查询测试脚本
BASE_URL="http://localhost:9001/api/v1"

echo "=== 公告关联查询测试 ==="

# 测试获取公告列表
echo "1. 测试获取公告列表"
RESPONSE=$(curl -s -X POST "${BASE_URL}/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 1,
    "page_size": 5
  }')

echo "响应: $RESPONSE"

# 检查响应状态
if echo "$RESPONSE" | jq -e '.code == 0' > /dev/null 2>&1; then
    echo "✅ 获取公告列表成功"
else
    echo "❌ 获取公告列表失败"
    exit 1
fi

echo -e "\n"

# 检查公告数据结构
echo "2. 检查公告数据结构"
if echo "$RESPONSE" | jq -e '.data.announcements[0] | has("id", "title", "content", "tag", "created_at", "banners")' > /dev/null 2>&1; then
    echo "✅ 公告数据结构正确"
else
    echo "❌ 公告数据结构不正确"
fi

echo -e "\n"

# 检查banners数据结构
echo "3. 检查banners数据结构"
if echo "$RESPONSE" | jq -e '.data.announcements[0].banners | type == "array"' > /dev/null 2>&1; then
    echo "✅ banners是数组格式"
else
    echo "❌ banners不是数组格式"
fi

echo -e "\n"

# 检查banners数组元素
echo "4. 检查banners数组元素"
if echo "$RESPONSE" | jq -e '.data.announcements[0].banners[0] | type == "string"' > /dev/null 2>&1; then
    echo "✅ banners数组元素是字符串（图片URL）"
else
    echo "❌ banners数组元素不是字符串"
fi

echo -e "\n"

# 统计公告和图片数量
echo "5. 统计公告和图片数量"
ANNOUNCEMENT_COUNT=$(echo "$RESPONSE" | jq '.data.announcements | length' 2>/dev/null || echo "0")
TOTAL_BANNERS=$(echo "$RESPONSE" | jq '[.data.announcements[].banners[]] | length' 2>/dev/null || echo "0")

echo "公告数量: $ANNOUNCEMENT_COUNT"
echo "总图片数量: $TOTAL_BANNERS"

echo -e "\n"

# 检查每个公告的图片关联
echo "6. 检查每个公告的图片关联"
for i in $(seq 0 $(($ANNOUNCEMENT_COUNT - 1))); do
    ANNOUNCEMENT_ID=$(echo "$RESPONSE" | jq -r ".data.announcements[$i].id" 2>/dev/null)
    BANNER_COUNT=$(echo "$RESPONSE" | jq ".data.announcements[$i].banners | length" 2>/dev/null || echo "0")
    echo "公告ID: $ANNOUNCEMENT_ID, 图片数量: $BANNER_COUNT"
done

echo -e "\n"

# 显示第一个公告的详细信息
echo "7. 第一个公告的详细信息"
FIRST_ANNOUNCEMENT=$(echo "$RESPONSE" | jq '.data.announcements[0]' 2>/dev/null)
if [ "$FIRST_ANNOUNCEMENT" != "null" ]; then
    echo "公告ID: $(echo "$FIRST_ANNOUNCEMENT" | jq -r '.id')"
    echo "标题: $(echo "$FIRST_ANNOUNCEMENT" | jq -r '.title')"
    echo "标签: $(echo "$FIRST_ANNOUNCEMENT" | jq -r '.tag')"
    echo "图片URLs: $(echo "$FIRST_ANNOUNCEMENT" | jq -r '.banners | join(", ")')"
else
    echo "没有找到公告数据"
fi

echo -e "\n=== 测试完成 ==="
echo "💡 提示: 如果数据库中没有公告数据，返回的列表可能为空"
echo "💡 提示: 每个公告的banners字段应该包含该公告对应的图片URL数组" 