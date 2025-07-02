#!/bin/bash

# 测试钱包创建时的用户存在性检查
BASE_URL="http://localhost:8080"
API_PREFIX="/api/v1"

echo "=== 测试钱包创建时的用户存在性检查 ==="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试充值接口（会自动创建钱包）
test_recharge_with_invalid_uid() {
    echo -e "\n${YELLOW}测试1: 使用不存在的UID进行充值${NC}"
    
    # 先登录获取token
    LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL$API_PREFIX/auth/login" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "test@example.com",
            "password": "123456"
        }')
    
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$TOKEN" ]; then
        echo -e "${RED}登录失败，无法获取token${NC}"
        return 1
    fi
    
    echo "使用不存在的UID: 99999999"
    
    # 尝试充值（会自动创建钱包）
    RESPONSE=$(curl -s -X POST "$BASE_URL$API_PREFIX/wallet/recharge" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "uid": "99999999",
            "amount": 100.00,
            "description": "测试充值"
        }')
    
    echo "响应: $RESPONSE"
    
    # 检查是否返回用户不存在的错误
    if echo "$RESPONSE" | grep -q "用户不存在"; then
        echo -e "${GREEN}✓ 测试通过: 正确返回用户不存在错误${NC}"
        return 0
    else
        echo -e "${RED}✗ 测试失败: 未返回预期的错误信息${NC}"
        return 1
    fi
}

# 测试提现接口（会自动创建钱包）
test_withdraw_with_invalid_uid() {
    echo -e "\n${YELLOW}测试2: 使用不存在的UID进行提现${NC}"
    
    # 先登录获取token
    LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL$API_PREFIX/auth/login" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "test@example.com",
            "password": "123456"
        }')
    
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$TOKEN" ]; then
        echo -e "${RED}登录失败，无法获取token${NC}"
        return 1
    fi
    
    echo "使用不存在的UID: 99999999"
    
    # 尝试提现（会自动创建钱包）
    RESPONSE=$(curl -s -X POST "$BASE_URL$API_PREFIX/wallet/withdraw" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "uid": "99999999",
            "amount": 50.00,
            "description": "测试提现",
            "password": "123456"
        }')
    
    echo "响应: $RESPONSE"
    
    # 检查是否返回用户不存在的错误
    if echo "$RESPONSE" | grep -q "用户不存在"; then
        echo -e "${GREEN}✓ 测试通过: 正确返回用户不存在错误${NC}"
        return 0
    else
        echo -e "${RED}✗ 测试失败: 未返回预期的错误信息${NC}"
        return 1
    fi
}

# 测试获取钱包信息接口（会自动创建钱包）
test_wallet_info_with_invalid_uid() {
    echo -e "\n${YELLOW}测试3: 使用不存在的UID获取钱包信息${NC}"
    
    # 先登录获取token
    LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL$API_PREFIX/auth/login" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "test@example.com",
            "password": "123456"
        }')
    
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$TOKEN" ]; then
        echo -e "${RED}登录失败，无法获取token${NC}"
        return 1
    fi
    
    echo "使用不存在的UID: 99999999"
    
    # 尝试获取钱包信息（会自动创建钱包）
    RESPONSE=$(curl -s -X POST "$BASE_URL$API_PREFIX/wallet/info" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "uid": "99999999"
        }')
    
    echo "响应: $RESPONSE"
    
    # 检查是否返回用户不存在的错误
    if echo "$RESPONSE" | grep -q "用户不存在"; then
        echo -e "${GREEN}✓ 测试通过: 正确返回用户不存在错误${NC}"
        return 0
    else
        echo -e "${RED}✗ 测试失败: 未返回预期的错误信息${NC}"
        return 1
    fi
}

# 运行测试
echo "开始运行测试..."

test_recharge_with_invalid_uid
test_withdraw_with_invalid_uid
test_wallet_info_with_invalid_uid

echo -e "\n${GREEN}=== 所有测试完成 ===${NC}" 