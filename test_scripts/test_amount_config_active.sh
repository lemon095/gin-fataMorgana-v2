#!/bin/bash

# 金额配置激活状态测试脚本
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

echo "=== 金额配置激活状态测试 ==="

# 测试获取充值配置列表（应该只返回激活的）
echo "1. 测试获取充值配置列表（只返回激活状态）"
RECHARGE_RESPONSE=$(curl -s -X POST "${BASE_URL}/amount-config/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "type": "recharge"
  }')

echo "响应: $RECHARGE_RESPONSE"

# 检查是否所有返回的配置都是激活状态
echo "2. 验证返回的配置都是激活状态"
if echo "$RECHARGE_RESPONSE" | jq -r '.data[]?.is_active' 2>/dev/null | grep -q "false"; then
    echo "❌ 发现未激活的配置被返回"
else
    echo "✅ 所有返回的配置都是激活状态"
fi

echo -e "\n"

# 测试获取提现配置列表（应该只返回激活的）
echo "3. 测试获取提现配置列表（只返回激活状态）"
WITHDRAW_RESPONSE=$(curl -s -X POST "${BASE_URL}/amount-config/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "type": "withdraw"
  }')

echo "响应: $WITHDRAW_RESPONSE"

# 检查是否所有返回的配置都是激活状态
echo "4. 验证返回的配置都是激活状态"
if echo "$WITHDRAW_RESPONSE" | jq -r '.data[]?.is_active' 2>/dev/null | grep -q "false"; then
    echo "❌ 发现未激活的配置被返回"
else
    echo "✅ 所有返回的配置都是激活状态"
fi

echo -e "\n"

# 测试获取配置详情（应该只返回激活的）
echo "5. 测试获取配置详情（只返回激活状态）"
DETAIL_RESPONSE=$(curl -s -X GET "${BASE_URL}/amount-config/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "响应: $DETAIL_RESPONSE"

# 检查返回的配置是否是激活状态
echo "6. 验证返回的配置是激活状态"
if echo "$DETAIL_RESPONSE" | jq -r '.data?.is_active' 2>/dev/null | grep -q "false"; then
    echo "❌ 返回了未激活的配置"
else
    echo "✅ 返回的配置是激活状态"
fi

echo -e "\n"

# 统计返回的配置数量
echo "7. 统计返回的配置数量"
RECHARGE_COUNT=$(echo "$RECHARGE_RESPONSE" | jq '.data | length' 2>/dev/null || echo "0")
WITHDRAW_COUNT=$(echo "$WITHDRAW_RESPONSE" | jq '.data | length' 2>/dev/null || echo "0")

echo "充值配置数量: $RECHARGE_COUNT"
echo "提现配置数量: $WITHDRAW_COUNT"

echo -e "\n=== 测试完成 ==="
echo "💡 提示: 如果数据库中有未激活的配置，它们不会在接口返回结果中出现" 