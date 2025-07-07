#!/bin/bash

# 测试拼单号唯一性脚本
echo "=== 拼单号唯一性测试 ==="

# 设置基础URL
BASE_URL="http://localhost:9001/api"

# 测试快速生成多个拼单
echo "1. 快速生成10个拼单测试唯一性..."

# 创建临时文件存储拼单号
TEMP_FILE=$(mktemp)

# 快速生成10个拼单
for i in {1..10}; do
    echo "生成第 $i 个拼单..."
    
    # 调用创建拼单接口
    RESPONSE=$(curl -s -X POST "$BASE_URL/group-buys" \
      -H "Content-Type: application/json" \
      -d '{
        "uid": "test123",
        "target_participants": 5,
        "per_person_amount": 1000.00,
        "group_buy_type": "normal"
      }')
    
    # 提取拼单号
    GROUP_BUY_NO=$(echo $RESPONSE | grep -o '"group_buy_no":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$GROUP_BUY_NO" ]; then
        echo "  拼单号: $GROUP_BUY_NO"
        echo "$GROUP_BUY_NO" >> $TEMP_FILE
    else
        echo "  ❌ 创建拼单失败: $RESPONSE"
    fi
    
    # 短暂延迟
    sleep 0.1
done

echo ""
echo "2. 检查拼单号唯一性..."

# 统计拼单号数量
TOTAL_COUNT=$(wc -l < $TEMP_FILE)
UNIQUE_COUNT=$(sort $TEMP_FILE | uniq | wc -l)

echo "总生成数量: $TOTAL_COUNT"
echo "唯一数量: $UNIQUE_COUNT"

if [ "$TOTAL_COUNT" -eq "$UNIQUE_COUNT" ]; then
    echo "✅ 所有拼单号都是唯一的！"
else
    echo "❌ 发现重复的拼单号！"
    echo "重复的拼单号:"
    sort $TEMP_FILE | uniq -d
fi

echo ""
echo "3. 所有生成的拼单号:"
sort $TEMP_FILE | nl

# 清理临时文件
rm $TEMP_FILE

echo ""
echo "=== 测试完成 ===" 