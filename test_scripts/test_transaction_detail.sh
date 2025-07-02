#!/bin/bash

# 测试交易详情接口
# 使用方法: ./test_transaction_detail.sh

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== 测试交易详情接口 ===${NC}"

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

# 2. 先获取交易列表，获取一个交易流水号
echo -e "\n${YELLOW}2. 获取交易列表${NC}"
TRANSACTIONS_RESPONSE=$(curl -s -X POST "${BASE_URL}/wallet/transactions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10
  }')

echo "交易列表响应: $TRANSACTIONS_RESPONSE"

# 提取第一个交易流水号
TRANSACTION_NO=$(echo "$TRANSACTIONS_RESPONSE" | jq -r '.data.transactions[0].transaction_no')
if [ "$TRANSACTION_NO" = "null" ] || [ -z "$TRANSACTION_NO" ]; then
    echo -e "${RED}没有找到交易记录，无法测试交易详情接口${NC}"
    exit 1
fi

echo -e "${GREEN}找到交易流水号: $TRANSACTION_NO${NC}"

# 3. 测试交易详情接口
echo -e "\n${YELLOW}3. 测试交易详情接口${NC}"
TRANSACTION_DETAIL_RESPONSE=$(curl -s -X POST "${BASE_URL}/wallet/transaction-detail" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"transaction_no\": \"$TRANSACTION_NO\"
  }")

echo "交易详情响应: $TRANSACTION_DETAIL_RESPONSE"

# 检查响应
CODE=$(echo "$TRANSACTION_DETAIL_RESPONSE" | jq -r '.code')
if [ "$CODE" = "0" ]; then
    echo -e "${GREEN}交易详情接口测试成功！${NC}"
    
    # 显示交易详情
    echo -e "\n${YELLOW}交易详情:${NC}"
    echo "$TRANSACTION_DETAIL_RESPONSE" | jq '.data | {
        "交易流水号": .transaction_no,
        "交易类型": .type_name,
        "交易金额": .amount,
        "交易状态": .status_name,
        "交易描述": .description,
        "创建时间": .created_at
    }'
else
    echo -e "${RED}交易详情接口测试失败！${NC}"
    echo "错误信息: $(echo "$TRANSACTION_DETAIL_RESPONSE" | jq -r '.message')"
fi

# 4. 测试无效的交易流水号
echo -e "\n${YELLOW}4. 测试无效的交易流水号${NC}"
INVALID_RESPONSE=$(curl -s -X POST "${BASE_URL}/wallet/transaction-detail" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "transaction_no": "INVALID_TRANSACTION_NO"
  }')

echo "无效交易流水号响应: $INVALID_RESPONSE"

INVALID_CODE=$(echo "$INVALID_RESPONSE" | jq -r '.code')
if [ "$INVALID_CODE" != "0" ]; then
    echo -e "${GREEN}无效交易流水号测试通过（正确返回错误）${NC}"
else
    echo -e "${RED}无效交易流水号测试失败（应该返回错误）${NC}"
fi

echo -e "\n${YELLOW}=== 测试完成 ===${NC}" 