#!/bin/bash

# 测试默认货币已改为PHP
echo "🧪 测试默认货币已改为PHP..."
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

# 2. 获取钱包信息（会自动创建钱包）
echo -e "\n${YELLOW}2. 获取钱包信息${NC}"
WALLET_RESPONSE=$(curl -s -X POST "${BASE_URL}/wallet/get" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "uid": "test_user"
  }')

echo "钱包响应: $WALLET_RESPONSE"

# 检查货币类型
CURRENCY=$(echo "$WALLET_RESPONSE" | jq -r '.data.currency')
if [ "$CURRENCY" = "PHP" ]; then
    echo -e "${GREEN}✓ 测试通过: 默认货币已改为PHP${NC}"
else
    echo -e "${RED}✗ 测试失败: 默认货币仍然是 $CURRENCY，应该是PHP${NC}"
fi

# 3. 检查数据库中的默认值
echo -e "\n${YELLOW}3. 检查数据库默认值${NC}"
DB_CURRENCY=$(mysql -u root -p123456 -e "SELECT currency FROM wallets WHERE uid = 'test_user' LIMIT 1;" 2>/dev/null | tail -n 1)
if [ "$DB_CURRENCY" = "PHP" ]; then
    echo -e "${GREEN}✓ 数据库中的货币类型也是PHP${NC}"
else
    echo -e "${RED}✗ 数据库中的货币类型是 $DB_CURRENCY，应该是PHP${NC}"
fi

echo -e "\n${BLUE}测试完成！${NC}" 