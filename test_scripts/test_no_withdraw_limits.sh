#!/bin/bash

# 测试提现限额已被移除
echo "🧪 测试提现限额已被移除..."
echo "=================================="

# 设置颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 测试基础URL
BASE_URL="http://localhost:8080/api/v1"

# 1. 先登录获取token
echo -e "\n${YELLOW}1. 用户登录${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test_user",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo -e "${RED}登录失败，无法获取token${NC}"
    exit 1
fi

echo -e "${GREEN}登录成功，获取到token${NC}"

# 2. 获取提现汇总信息，检查限额
echo -e "\n${YELLOW}2. 获取提现汇总信息${NC}"
SUMMARY_RESPONSE=$(curl -s -X POST "${BASE_URL}/wallet/withdraw-summary" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "汇总响应: $SUMMARY_RESPONSE"

# 检查是否还有限额信息
LIMITS=$(echo "$SUMMARY_RESPONSE" | jq -r '.data.limits')
if [ "$LIMITS" = "{}" ] || [ "$LIMITS" = "null" ]; then
    echo -e "${GREEN}✓ 测试通过: 限额信息已移除${NC}"
else
    echo -e "${RED}✗ 测试失败: 仍有限额信息 $LIMITS${NC}"
fi

# 3. 测试大额提现申请（应该不会因为限额而失败）
echo -e "\n${YELLOW}3. 测试大额提现申请${NC}"
LARGE_AMOUNT=999999999.99
WITHDRAW_RESPONSE=$(curl -s -X POST "${BASE_URL}/wallet/withdraw" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"amount\": $LARGE_AMOUNT,
    \"password\": \"123456\"
  }")

echo "大额提现响应: $WITHDRAW_RESPONSE"

# 检查是否因为限额而失败
ERROR_CODE=$(echo "$WITHDRAW_RESPONSE" | jq -r '.code')
if [ "$ERROR_CODE" = "3018" ] || [ "$ERROR_CODE" = "3021" ]; then
    echo -e "${RED}✗ 测试失败: 仍然有提现限额限制${NC}"
else
    echo -e "${GREEN}✓ 测试通过: 大额提现不会因为限额而失败${NC}"
    echo -e "${BLUE}  注意: 可能会因为余额不足而失败，这是正常的${NC}"
fi

echo -e "\n${BLUE}测试完成！${NC}" 