#!/bin/bash

# 公告数据排查脚本
echo "=== 公告数据排查 ==="

# 1. 检查数据库连接
echo "1. 检查数据库连接..."
if mysql -u root -p123456 -e "USE gin_fata_morgana; SELECT 1;" >/dev/null 2>&1; then
    echo "✅ 数据库连接正常"
else
    echo "❌ 数据库连接失败"
    exit 1
fi

echo ""

# 2. 检查公告表是否存在
echo "2. 检查公告表结构..."
mysql -u root -p123456 -e "USE gin_fata_morgana; DESCRIBE announcements;" 2>/dev/null
if [ $? -eq 0 ]; then
    echo "✅ 公告表存在"
else
    echo "❌ 公告表不存在"
    exit 1
fi

echo ""

# 3. 统计公告数据
echo "3. 统计公告数据..."
echo "总公告数量:"
mysql -u root -p123456 -e "USE gin_fata_morgana; SELECT COUNT(*) as total FROM announcements;" 2>/dev/null

echo "已发布公告数量:"
mysql -u root -p123456 -e "USE gin_fata_morgana; SELECT COUNT(*) as published FROM announcements WHERE status = 1 AND is_publish = 1;" 2>/dev/null

echo "草稿公告数量:"
mysql -u root -p123456 -e "USE gin_fata_morgana; SELECT COUNT(*) as draft FROM announcements WHERE status = 0 OR is_publish = 0;" 2>/dev/null

echo ""

# 4. 查看公告详情
echo "4. 查看公告详情（前5条）:"
mysql -u root -p123456 -e "USE gin_fata_morgana; SELECT id, title, status, is_publish, created_at FROM announcements ORDER BY created_at DESC LIMIT 5;" 2>/dev/null

echo ""

# 5. 测试接口
echo "5. 测试公告接口..."
BASE_URL="http://localhost:9001/api/v1"
RESPONSE=$(curl -s -X POST "${BASE_URL}/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 1,
    "page_size": 10
  }')

echo "接口响应: $RESPONSE"

echo ""

# 6. 检查Redis缓存
echo "6. 检查Redis缓存..."
if command -v redis-cli >/dev/null 2>&1; then
    CACHE_KEYS=$(redis-cli keys "announcement:list:*" 2>/dev/null | wc -l)
    echo "公告缓存键数量: $CACHE_KEYS"
    
    if [ $CACHE_KEYS -gt 0 ]; then
        echo "缓存键列表:"
        redis-cli keys "announcement:list:*" 2>/dev/null
    fi
else
    echo "⚠️  redis-cli未安装"
fi

echo ""

# 7. 创建测试数据（如果没有数据）
echo "7. 检查是否需要创建测试数据..."
TOTAL_COUNT=$(mysql -u root -p123456 -e "USE gin_fata_morgana; SELECT COUNT(*) FROM announcements;" 2>/dev/null | tail -n 1)

if [ "$TOTAL_COUNT" = "0" ]; then
    echo "📝 数据库中没有公告数据，建议创建测试数据"
    echo "可以执行以下SQL创建测试公告:"
    echo ""
    echo "INSERT INTO announcements (title, content, tag, status, is_publish, created_at) VALUES"
    echo "('系统维护通知', '系统将于今晚22:00-24:00进行维护升级，期间可能影响部分功能使用。', '系统通知', 1, 1, NOW()),"
    echo "('新功能上线', '我们新增了用户等级功能，完成任务可以获得经验值提升等级。', '功能更新', 1, 1, NOW()),"
    echo "('活动公告', '参与拼单活动可以获得额外奖励，快来参与吧！', '活动', 1, 1, NOW());"
else
    echo "✅ 数据库中有 $TOTAL_COUNT 条公告数据"
fi

echo ""
echo "=== 排查完成 ===" 