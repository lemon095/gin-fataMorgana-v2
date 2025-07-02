#!/bin/bash

# 测试拼单接口错误码
echo "=== 测试拼单接口错误码 ==="

# 设置API基础URL
API_BASE="http://localhost:8080/api/v1"

# 测试用户登录
echo "1. 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser1",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token: $TOKEN"

if [ -n "$TOKEN" ]; then
    # 测试参与不存在的拼单
    echo -e "\n2. 测试参与不存在的拼单 (错误码: 1010)..."
    JOIN_NOTEXIST_RESPONSE=$(curl -s -X POST "$API_BASE/groupBuy/join" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "group_buy_no": "NOTEXIST123"
      }')

    echo "参与不存在拼单响应: $JOIN_NOTEXIST_RESPONSE"

    # 测试参与已过期的拼单（需要先创建一个过期的拼单）
    echo -e "\n3. 测试参与已过期的拼单 (错误码: 1011)..."
    # 这里需要数据库中有过期的拼单数据才能测试
    
    # 测试参与已被占用的拼单（需要先创建一个已被占用的拼单）
    echo -e "\n4. 测试参与已被占用的拼单 (错误码: 1012)..."
    # 这里需要数据库中有已被占用的拼单数据才能测试

    # 测试参数错误
    echo -e "\n5. 测试参数错误 (错误码: 1003)..."
    PARAM_ERROR_RESPONSE=$(curl -s -X POST "$API_BASE/groupBuy/join" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "invalid_param": "test"
      }')

    echo "参数错误响应: $PARAM_ERROR_RESPONSE"

    # 测试缺少参数
    echo -e "\n6. 测试缺少参数 (错误码: 1003)..."
    MISSING_PARAM_RESPONSE=$(curl -s -X POST "$API_BASE/groupBuy/join" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{}')

    echo "缺少参数响应: $MISSING_PARAM_RESPONSE"

else
    echo "跳过需要认证的测试（登录失败）"
fi

# 测试未认证访问
echo -e "\n7. 测试未认证访问 (错误码: 401)..."
UNAUTH_RESPONSE=$(curl -s -X POST "$API_BASE/groupBuy/join" \
  -H "Content-Type: application/json" \
  -d '{
    "group_buy_no": "TEST123"
  }')

echo "未认证访问响应: $UNAUTH_RESPONSE"

echo -e "\n=== 错误码说明 ==="
echo "1010 - 拼单不存在或已被删除"
echo "1011 - 拼单已超过截止时间"
echo "1012 - 拼单已被其他用户参与"
echo "1001 - 数据库操作失败"
echo "1003 - 参数错误"
echo "401  - 认证失败"
echo "0    - 操作成功" 