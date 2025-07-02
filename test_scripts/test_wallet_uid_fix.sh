#!/bin/bash

# 测试钱包接口中uid获取的修复
# 验证修复后的接口能正确获取用户uid而不是错误的user_id转换

BASE_URL="http://localhost:8080"
API_PREFIX="/api/v1"

echo "=== 测试钱包接口中uid获取的修复 ==="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_wallet_api() {
    local test_name="$1"
    local api_endpoint="$2"
    local expected_uid="$3"
    
    echo -e "\n${YELLOW}测试: $test_name${NC}"
    echo "API端点: $api_endpoint"
    
    # 调用API
    response=$(curl -s -X POST "$BASE_URL$API_PREFIX$api_endpoint" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "{}")
    
    echo "响应: $response"
    
    # 检查是否包含正确的uid
    if echo "$response" | grep -q "\"uid\":\"$expected_uid\""; then
        echo -e "${GREEN}✓ 测试通过: 返回了正确的uid${NC}"
        return 0
    else
        echo -e "${RED}✗ 测试失败: 未返回正确的uid${NC}"
        return 1
    fi
}

# 登录获取token
echo "正在登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL$API_PREFIX/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "test@example.com",
        "password": "123456"
    }')

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}登录失败，无法获取token${NC}"
    echo "登录响应: $LOGIN_RESPONSE"
    exit 1
fi

echo -e "${GREEN}登录成功，获取到token${NC}"

# 获取用户信息以获取正确的uid
echo -e "\n${YELLOW}获取用户信息...${NC}"
USER_INFO_RESPONSE=$(curl -s -X GET "$BASE_URL$API_PREFIX/auth/profile" \
    -H "Authorization: Bearer $TOKEN")

echo "用户信息响应: $USER_INFO_RESPONSE"

# 从响应中提取uid
CORRECT_UID=$(echo "$USER_INFO_RESPONSE" | grep -o '"uid":"[^"]*"' | cut -d'"' -f4)

if [ -z "$CORRECT_UID" ]; then
    echo -e "${RED}无法获取用户UID，跳过测试${NC}"
    exit 1
fi

echo -e "${GREEN}获取到正确的UID: $CORRECT_UID${NC}"

# 测试1: 获取钱包信息
test_wallet_api "获取钱包信息" "/wallet/info" "$CORRECT_UID"

# 测试2: 获取交易记录
test_wallet_api "获取交易记录" "/wallet/transactions" "$CORRECT_UID" <<< '{"page": 1, "page_size": 10}'

# 测试3: 获取提现汇总
test_wallet_api "获取提现汇总" "/wallet/withdraw/summary" "$CORRECT_UID"

echo -e "\n${GREEN}=== 测试完成 ===${NC}"
echo -e "${YELLOW}注意: 这些测试验证了修复后的接口能正确获取用户uid，而不是错误地将user_id转换为uid${NC}" 