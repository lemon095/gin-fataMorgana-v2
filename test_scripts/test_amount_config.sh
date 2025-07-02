#!/bin/bash

# 金额配置接口测试脚本
BASE_URL="http://localhost:9001/api/v1"

# 获取访问令牌（需要先登录）
echo "🔐 获取访问令牌..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user1@example.com",
    "password": "123456"
  }')

# 提取访问令牌
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$ACCESS_TOKEN" ]; then
    echo "❌ 登录失败，无法获取访问令牌"
    echo "响应: $LOGIN_RESPONSE"
    exit 1
fi

echo "✅ 获取访问令牌成功: ${ACCESS_TOKEN:0:20}..."
echo

echo "=== 金额配置接口测试 ==="

# 测试获取充值金额配置列表
echo "1. 测试获取充值金额配置列表"
curl -X POST "${BASE_URL}/amount-config/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "type": "recharge"
  }' | jq '.'

echo -e "\n"

# 测试获取提现金额配置列表
echo "2. 测试获取提现金额配置列表"
curl -X POST "${BASE_URL}/amount-config/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "type": "withdraw"
  }' | jq '.'

echo -e "\n"

# 测试无效类型
echo "3. 测试无效类型"
curl -X POST "${BASE_URL}/amount-config/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "type": "invalid"
  }' | jq '.'

echo -e "\n"

# 测试缺少类型参数
echo "4. 测试缺少类型参数"
curl -X POST "${BASE_URL}/amount-config/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{}' | jq '.'

echo -e "\n"

# 测试获取配置详情（如果存在ID为1的配置）
echo "5. 测试获取配置详情 (ID: 1)"
curl -X GET "${BASE_URL}/amount-config/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq '.'

echo -e "\n"

# 测试获取不存在的配置详情
echo "6. 测试获取不存在的配置详情 (ID: 999)"
curl -X GET "${BASE_URL}/amount-config/999" \
  -H "Content-Type: application/json" | jq '.'

echo -e "\n"

# 测试无效ID格式
echo "7. 测试无效ID格式"
curl -X GET "${BASE_URL}/amount-config/abc" \
  -H "Content-Type: application/json" | jq '.'

echo -e "\n=== 测试完成 ===" 