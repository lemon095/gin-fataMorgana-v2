#!/bin/bash

# 测试操作员字段修复
# 验证充值、提现申请时OperatorUid为空，后台处理时由管理员设置

BASE_URL="http://localhost:8080"
API_PREFIX="/api/v1"

echo "=== 测试操作员字段修复 ==="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_transaction_operator() {
    local test_name="$1"
    local api_endpoint="$2"
    local request_data="$3"
    
    echo -e "\n${YELLOW}测试: $test_name${NC}"
    echo "API端点: $api_endpoint"
    echo "请求数据: $request_data"
    
    # 调用API
    response=$(curl -s -X POST "$BASE_URL$API_PREFIX$api_endpoint" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "$request_data")
    
    echo "响应: $response"
    
    # 提取交易流水号
    transaction_no=$(echo "$response" | grep -o '"transaction_no":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$transaction_no" ]; then
        echo -e "${GREEN}✓ 交易申请成功，流水号: $transaction_no${NC}"
        
        # 查询交易详情，检查OperatorUid是否为空
        echo -e "\n${YELLOW}查询交易详情...${NC}"
        detail_response=$(curl -s -X POST "$BASE_URL$API_PREFIX/wallet/transaction/detail" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $TOKEN" \
            -d "{\"transaction_no\": \"$transaction_no\"}")
        
        echo "交易详情: $detail_response"
        
        # 检查OperatorUid是否为空
        operator_uid=$(echo "$detail_response" | grep -o '"operator_uid":"[^"]*"' | cut -d'"' -f4)
        if [ "$operator_uid" = "" ]; then
            echo -e "${GREEN}✓ 操作员字段为空，符合预期${NC}"
            return 0
        else
            echo -e "${RED}✗ 操作员字段不为空: $operator_uid${NC}"
            return 1
        fi
    else
        echo -e "${RED}✗ 交易申请失败${NC}"
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

# 获取用户信息以获取uid
echo -e "\n${YELLOW}获取用户信息...${NC}"
USER_INFO_RESPONSE=$(curl -s -X GET "$BASE_URL$API_PREFIX/auth/profile" \
    -H "Authorization: Bearer $TOKEN")

USER_UID=$(echo "$USER_INFO_RESPONSE" | grep -o '"uid":"[^"]*"' | cut -d'"' -f4)

if [ -z "$USER_UID" ]; then
    echo -e "${RED}无法获取用户UID，跳过测试${NC}"
    exit 1
fi

echo -e "${GREEN}获取到用户UID: $USER_UID${NC}"

# 测试1: 充值申请
test_transaction_operator "充值申请" "/wallet/recharge" "{\"uid\": \"$USER_UID\", \"amount\": 100, \"description\": \"测试充值\"}"

# 测试2: 提现申请
test_transaction_operator "提现申请" "/wallet/withdraw" "{\"amount\": 10, \"password\": \"123456\"}"

echo -e "\n${GREEN}=== 测试完成 ===${NC}"
echo -e "${YELLOW}注意: 这些测试验证了申请时OperatorUid为空，后台处理时由管理员设置${NC}" 