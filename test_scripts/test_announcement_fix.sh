#!/bin/bash

# 公告状态修复测试脚本
echo "=== 公告状态修复测试 ==="

# 1. 执行SQL修复脚本
echo "1. 执行公告状态修复..."
mysql -u root -p123456 gin_fata_morgana < test_scripts/fix_announcement_status.sql

if [ $? -eq 0 ]; then
    echo "✅ SQL执行成功"
else
    echo "❌ SQL执行失败"
    exit 1
fi

echo ""

# 2. 测试公告接口
echo "2. 测试公告接口..."
BASE_URL="http://localhost:9001/api/v1"
RESPONSE=$(curl -s -X POST "${BASE_URL}/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 1,
    "page_size": 10
  }')

echo "接口响应: $RESPONSE"

# 检查响应状态
if echo "$RESPONSE" | jq -e '.code == 0' > /dev/null 2>&1; then
    echo "✅ 公告接口返回成功"
    
    # 检查是否有数据
    ANNOUNCEMENT_COUNT=$(echo "$RESPONSE" | jq '.data.announcements | length' 2>/dev/null || echo "0")
    TOTAL_COUNT=$(echo "$RESPONSE" | jq '.data.pagination.total' 2>/dev/null || echo "0")
    
    echo "公告数量: $ANNOUNCEMENT_COUNT"
    echo "总数量: $TOTAL_COUNT"
    
    if [ "$ANNOUNCEMENT_COUNT" -gt 0 ]; then
        echo "✅ 公告数据已正常返回"
    else
        echo "⚠️  公告数据仍为空，可能需要检查数据库连接或创建测试数据"
    fi
else
    echo "❌ 公告接口返回失败"
fi

echo ""

# 3. 清除Redis缓存（如果需要）
echo "3. 清除Redis缓存..."
if command -v redis-cli >/dev/null 2>&1; then
    redis-cli del "announcement:list:page:1:size:10" 2>/dev/null
    echo "✅ Redis缓存已清除"
else
    echo "⚠️  redis-cli未安装，无法清除缓存"
fi

echo ""

# 4. 再次测试接口（清除缓存后）
echo "4. 再次测试接口（清除缓存后）..."
RESPONSE2=$(curl -s -X POST "${BASE_URL}/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 1,
    "page_size": 10
  }')

echo "清除缓存后的响应: $RESPONSE2"

echo ""

echo "=== 修复测试完成 ===" 