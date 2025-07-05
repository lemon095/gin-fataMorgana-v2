#!/bin/bash

# 测试用户拼单逻辑
echo "=== 测试用户拼单逻辑 ==="

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

# 测试用户注册
echo "2. 注册测试用户..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser_groupbuy",
    "password": "test123456",
    "email": "test_groupbuy@example.com"
  }')

echo "注册响应: $REGISTER_RESPONSE"

# 提取用户token
USER_TOKEN=$(echo $REGISTER_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$USER_TOKEN" ]; then
    echo "❌ 用户注册失败，无法获取token"
    exit 1
fi

echo "✅ 用户注册成功，获取到token: ${USER_TOKEN:0:20}..."

# 测试获取活跃拼单详情
echo "3. 获取活跃拼单详情..."
GROUP_BUY_DETAIL_RESPONSE=$(curl -s -X POST "$BASE_URL/group-buy/active-detail" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN")

echo "拼单详情响应: $GROUP_BUY_DETAIL_RESPONSE"

# 检查是否有拼单数据
if echo $GROUP_BUY_DETAIL_RESPONSE | grep -q '"has_data":true'; then
    echo "✅ 发现活跃拼单"
    
    # 提取拼单编号
    GROUP_BUY_NO=$(echo $GROUP_BUY_DETAIL_RESPONSE | grep -o '"group_buy_no":"[^"]*"' | cut -d'"' -f4)
    if [ ! -z "$GROUP_BUY_NO" ]; then
        echo "📋 拼单编号: $GROUP_BUY_NO"
        
        # 测试参与拼单
        echo "4. 参与拼单..."
        JOIN_RESPONSE=$(curl -s -X POST "$BASE_URL/group-buy/join" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $USER_TOKEN" \
          -d "{
            \"group_buy_no\": \"$GROUP_BUY_NO\"
          }")
        
        echo "参与拼单响应: $JOIN_RESPONSE"
        
        # 检查是否成功参与
        if echo $JOIN_RESPONSE | grep -q '"code":200'; then
            echo "✅ 成功参与拼单"
            
            # 提取订单ID
            ORDER_ID=$(echo $JOIN_RESPONSE | grep -o '"order_id":[0-9]*' | cut -d':' -f2)
            if [ ! -z "$ORDER_ID" ]; then
                echo "📋 生成的订单ID: $ORDER_ID"
                
                # 测试查询订单详情
                echo "5. 查询订单详情..."
                ORDER_DETAIL_RESPONSE=$(curl -s -X POST "$BASE_URL/order/detail" \
                  -H "Content-Type: application/json" \
                  -H "Authorization: Bearer $USER_TOKEN" \
                  -d "{
                    \"order_no\": \"$ORDER_ID\"
                  }")
                
                echo "订单详情响应: $ORDER_DETAIL_RESPONSE"
                
                # 检查任务数量
                if echo $ORDER_DETAIL_RESPONSE | grep -q '"code":200'; then
                    echo "✅ 订单详情查询成功"
                    
                    # 提取任务数量
                    LIKE_COUNT=$(echo $ORDER_DETAIL_RESPONSE | grep -o '"like_count":[0-9]*' | cut -d':' -f2)
                    SHARE_COUNT=$(echo $ORDER_DETAIL_RESPONSE | grep -o '"share_count":[0-9]*' | cut -d':' -f2)
                    FOLLOW_COUNT=$(echo $ORDER_DETAIL_RESPONSE | grep -o '"follow_count":[0-9]*' | cut -d':' -f2)
                    FAVORITE_COUNT=$(echo $ORDER_DETAIL_RESPONSE | grep -o '"favorite_count":[0-9]*' | cut -d':' -f2)
                    
                    echo "📊 任务数量统计:"
                    echo "  点赞数量: $LIKE_COUNT"
                    echo "  分享数量: $SHARE_COUNT"
                    echo "  关注数量: $FOLLOW_COUNT"
                    echo "  收藏数量: $FAVORITE_COUNT"
                    
                    # 验证任务数量逻辑
                    TOTAL_TASKS=$((LIKE_COUNT + SHARE_COUNT + FOLLOW_COUNT + FAVORITE_COUNT))
                    echo "  总任务数量: $TOTAL_TASKS"
                    
                    if [ "$TOTAL_TASKS" -ge 1 ] && [ "$TOTAL_TASKS" -le 4 ]; then
                        echo "✅ 任务数量符合逻辑（1-4个）"
                    else
                        echo "❌ 任务数量不符合逻辑: $TOTAL_TASKS"
                    fi
                    
                    # 检查每个类型的数量是否为0或1
                    if [ "$LIKE_COUNT" -le 1 ] && [ "$SHARE_COUNT" -le 1 ] && [ "$FOLLOW_COUNT" -le 1 ] && [ "$FAVORITE_COUNT" -le 1 ]; then
                        echo "✅ 每个类型数量正确（0或1）"
                    else
                        echo "❌ 类型数量异常"
                    fi
                else
                    echo "❌ 订单详情查询失败"
                fi
            else
                echo "❌ 无法获取订单ID"
            fi
        else
            echo "❌ 参与拼单失败"
            echo "错误信息: $JOIN_RESPONSE"
        fi
    else
        echo "❌ 无法获取拼单编号"
    fi
else
    echo "⚠️  没有活跃拼单，可能需要先生成水单"
fi

echo "=== 测试完成 ===" 