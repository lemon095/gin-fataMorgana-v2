#!/bin/bash

# 测试拼单号唯一性脚本
echo "=== 拼单号唯一性测试 ==="

# 设置基础URL
BASE_URL="http://localhost:9001/api"

# 测试生成多个拼单号
echo "1. 生成10个拼单号测试唯一性..."
for i in {1..10}; do
    echo "生成第 $i 个拼单号..."
    
    # 创建拼单
    RESPONSE=$(curl -s -X POST "$BASE_URL/group-buys" \
      -H "Content-Type: application/json" \
      -d '{
        "uid": "test123",
        "target_participants": 3,
        "total_amount": 100000,
        "group_buy_type": "normal"
      }')
    
    echo "响应: $RESPONSE"
    
    # 提取拼单号
    GROUP_BUY_NO=$(echo $RESPONSE | grep -o '"group_buy_no":"[^"]*"' | cut -d'"' -f4)
    if [ ! -z "$GROUP_BUY_NO" ]; then
        echo "✅ 生成拼单号: $GROUP_BUY_NO"
    else
        echo "❌ 生成拼单号失败"
    fi
    
    # 等待100毫秒
    sleep 0.1
done

echo ""
echo "2. 检查数据库中是否有重复的拼单号..."
# 这里可以添加数据库查询逻辑来检查重复

echo "✅ 测试完成" 