#!/bin/bash

# 测试银行卡信息脱敏功能
# 1. 先绑定银行卡
# 2. 然后获取银行卡信息，验证是否脱敏

BASE_URL="http://localhost:8080"
API_VERSION="v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 测试银行卡信息脱敏功能 ===${NC}"

# 1. 用户登录获取token
echo -e "${YELLOW}1. 用户登录...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/${API_VERSION}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}登录失败，无法获取token${NC}"
    exit 1
fi

echo -e "${GREEN}登录成功，获取到token: ${TOKEN:0:20}...${NC}"

# 2. 绑定银行卡
echo -e "${YELLOW}2. 绑定银行卡...${NC}"
BIND_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/${API_VERSION}/auth/bank-card" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "bank_name": "中国工商银行",
    "card_holder": "张三",
    "card_number": "6222021234567890123",
    "card_type": "借记卡"
  }')

echo "绑定银行卡响应: $BIND_RESPONSE"

# 检查绑定是否成功
if echo "$BIND_RESPONSE" | grep -q '"code":200'; then
    echo -e "${GREEN}银行卡绑定成功${NC}"
else
    echo -e "${RED}银行卡绑定失败${NC}"
    echo "$BIND_RESPONSE"
    exit 1
fi

# 3. 获取银行卡信息（应该显示脱敏后的信息）
echo -e "${YELLOW}3. 获取银行卡信息（脱敏显示）...${NC}"
GET_CARD_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/${API_VERSION}/auth/bank-card" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "获取银行卡信息响应: $GET_CARD_RESPONSE"

# 检查返回的银行卡信息是否脱敏
if echo "$GET_CARD_RESPONSE" | grep -q '"card_number":"6222 \*\*\*\* \*\*\*\* 0123"'; then
    echo -e "${GREEN}✓ 银行卡号脱敏正确${NC}"
else
    echo -e "${RED}✗ 银行卡号脱敏失败${NC}"
    echo "期望格式: 6222 **** **** 0123"
    echo "实际返回: $(echo "$GET_CARD_RESPONSE" | grep -o '"card_number":"[^"]*"')"
fi

if echo "$GET_CARD_RESPONSE" | grep -q '"card_holder":"张\*\*"'; then
    echo -e "${GREEN}✓ 持卡人姓名脱敏正确${NC}"
else
    echo -e "${RED}✗ 持卡人姓名脱敏失败${NC}"
    echo "期望格式: 张**"
    echo "实际返回: $(echo "$GET_CARD_RESPONSE" | grep -o '"card_holder":"[^"]*"')"
fi

# 4. 验证银行名称和卡类型没有脱敏
if echo "$GET_CARD_RESPONSE" | grep -q '"bank_name":"中国工商银行"'; then
    echo -e "${GREEN}✓ 银行名称未脱敏（正确）${NC}"
else
    echo -e "${RED}✗ 银行名称显示异常${NC}"
fi

if echo "$GET_CARD_RESPONSE" | grep -q '"card_type":"借记卡"'; then
    echo -e "${GREEN}✓ 卡类型未脱敏（正确）${NC}"
else
    echo -e "${RED}✗ 卡类型显示异常${NC}"
fi

echo -e "${BLUE}=== 测试完成 ===${NC}" 