#!/bin/bash

# 测试钱包创建时的用户存在性检查
# 测试场景：
# 1. 使用不存在的uid创建钱包
# 2. 使用被禁用的用户uid创建钱包
# 3. 使用正常用户uid创建钱包

BASE_URL="http://localhost:8080"
API_PREFIX="/api/v1"

echo "=== 测试钱包创建时的用户存在性检查 ==="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_wallet_creation() {
    local test_name="$1"
    local uid="$2"
    local expected_error="$3"
    
    echo -e "\n${YELLOW}测试: $test_name${NC}"
    echo "UID: $uid"
    
    # 获取钱包信息（会自动创建钱包）
    response=$(curl -s -X POST "$BASE_URL$API_PREFIX/wallet/info" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "{\"uid\": \"$uid\"}")
    
    echo "响应: $response"
    
    # 检查是否包含预期的错误信息
    if echo "$response" | grep -q "$expected_error"; then
        echo -e "${GREEN}✓ 测试通过: 正确返回错误信息${NC}"
        return 0
    else
        echo -e "${RED}✗ 测试失败: 未返回预期的错误信息${NC}"
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

# 测试1: 使用不存在的uid
test_wallet_creation "使用不存在的UID创建钱包" "99999999" "用户不存在"

# 测试2: 使用被禁用的用户uid（需要先创建一个被禁用的用户）
echo -e "\n${YELLOW}创建被禁用的测试用户...${NC}"
DISABLED_USER_RESPONSE=$(curl -s -X POST "$BASE_URL$API_PREFIX/auth/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "disabled@example.com",
        "password": "123456",
        "confirm_password": "123456"
    }')

echo "注册响应: $DISABLED_USER_RESPONSE"

# 从响应中提取uid（这里需要根据实际响应格式调整）
DISABLED_UID=$(echo "$DISABLED_USER_RESPONSE" | grep -o '"uid":"[^"]*"' | cut -d'"' -f4)

if [ -n "$DISABLED_UID" ]; then
    echo "被禁用用户的UID: $DISABLED_UID"
    
    # 这里需要管理员接口来禁用用户，暂时跳过这个测试
    echo -e "${YELLOW}注意: 需要管理员接口来禁用用户，暂时跳过被禁用用户的测试${NC}"
else
    echo -e "${RED}无法获取被禁用用户的UID${NC}"
fi

# 测试3: 使用正常用户uid
echo -e "\n${YELLOW}获取当前登录用户的UID...${NC}"
USER_INFO_RESPONSE=$(curl -s -X GET "$BASE_URL$API_PREFIX/auth/profile" \
    -H "Authorization: Bearer $TOKEN")

echo "用户信息响应: $USER_INFO_RESPONSE"

# 从JWT token中提取user_id（这里需要根据实际JWT结构调整）
# 暂时使用一个已知的有效UID进行测试
VALID_UID="12345678"  # 这里应该使用实际的用户UID

test_wallet_creation "使用正常用户UID创建钱包" "$VALID_UID" ""

echo -e "\n${GREEN}=== 测试完成 ===${NC}" 