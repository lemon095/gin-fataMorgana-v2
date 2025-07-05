#!/bin/bash

# 测试新的水单生成逻辑
echo "=== 测试新的水单生成逻辑 ==="

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 测试管理员登录
echo "1. 管理员登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$TOKEN" ]; then
    echo "❌ 登录失败，无法获取token"
    exit 1
fi

echo "✅ 登录成功，获取到token: ${TOKEN:0:20}..."

# 测试生成水单
echo "2. 测试生成水单..."
GENERATE_RESPONSE=$(curl -s -X POST "$BASE_URL/cron/generate-fake-orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "count": 5
  }')

echo "生成水单响应: $GENERATE_RESPONSE"

# 检查是否成功
if echo $GENERATE_RESPONSE | grep -q '"code":200'; then
    echo "✅ 水单生成成功"
    
    # 提取统计信息
    TOTAL_GENERATED=$(echo $GENERATE_RESPONSE | grep -o '"total_generated":[0-9]*' | cut -d':' -f2)
    PURCHASE_ORDERS=$(echo $GENERATE_RESPONSE | grep -o '"purchase_orders":[0-9]*' | cut -d':' -f2)
    GROUP_BUY_ORDERS=$(echo $GENERATE_RESPONSE | grep -o '"group_buy_orders":[0-9]*' | cut -d':' -f2)
    TOTAL_AMOUNT=$(echo $GENERATE_RESPONSE | grep -o '"total_amount":[0-9]*\.?[0-9]*' | cut -d':' -f2)
    
    echo "📊 生成统计:"
    echo "  总生成数量: $TOTAL_GENERATED"
    echo "  购买单数量: $PURCHASE_ORDERS"
    echo "  拼单数量: $GROUP_BUY_ORDERS"
    echo "  总金额: $TOTAL_AMOUNT"
    
    # 验证金额范围
    if [ ! -z "$TOTAL_AMOUNT" ] && [ "$TOTAL_AMOUNT" -gt 0 ]; then
        echo "✅ 总金额生成正常"
        
        # 验证拼单逻辑
        if [ "$GROUP_BUY_ORDERS" -gt 0 ]; then
            echo "📊 拼单验证:"
            echo "  拼单数量: $GROUP_BUY_ORDERS"
            echo "  拼单总金额: $TOTAL_AMOUNT"
            echo "  💡 拼单逻辑: 单价×任务数量=总金额，总金额÷目标人数=人均金额"
        fi
    else
        echo "❌ 总金额异常: $TOTAL_AMOUNT"
    fi
else
    echo "❌ 水单生成失败"
    echo "错误信息: $GENERATE_RESPONSE"
fi

# 测试查询水单列表
echo "3. 测试查询水单列表..."
ORDER_LIST_RESPONSE=$(curl -s -X POST "$BASE_URL/order/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10,
    "status": 1
  }')

echo "订单列表响应: $ORDER_LIST_RESPONSE"

# 检查是否有水单数据
if echo $ORDER_LIST_RESPONSE | grep -q '"code":200'; then
    echo "✅ 订单列表查询成功"
    
    # 检查是否有系统订单（水单）
    if echo $ORDER_LIST_RESPONSE | grep -q '"is_system_order":true'; then
        echo "✅ 发现系统订单（水单）"
    else
        echo "⚠️  未发现系统订单，可能需要等待生成"
    fi
else
    echo "❌ 订单列表查询失败"
fi

echo "=== 测试完成 ===" 