#!/bin/bash

# 测试优化后的提现接口
# 验证简化参数、增强校验等功能

BASE_URL="http://localhost:8080"
API_PREFIX="/api/v1"

echo "=== 测试优化后的提现接口 ==="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_withdraw() {
    local test_name="$1"
    local request_data="$2"
    local expected_error="$3"
    
    echo -e "\n${YELLOW}测试: $test_name${NC}"
    echo "请求数据: $request_data"
    
    # 调用提现接口
    response=$(curl -s -X POST "$BASE_URL$API_PREFIX/wallet/withdraw" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "$request_data")
    
    echo "响应: $response"
    
    # 检查是否包含预期的错误信息
    if [ -n "$expected_error" ]; then
        if echo "$response" | grep -q "$expected_error"; then
            echo -e "${GREEN}✓ 测试通过: 正确返回错误信息${NC}"
            return 0
        else
            echo -e "${RED}✗ 测试失败: 未返回预期的错误信息${NC}"
            return 1
        fi
    else
        # 检查是否成功
        if echo "$response" | grep -q '"code":200'; then
            echo -e "${GREEN}✓ 测试通过: 提现申请成功${NC}"
            return 0
        else
            echo -e "${RED}✗ 测试失败: 提现申请失败${NC}"
            return 1
        fi
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

# 测试1: 参数缺失测试
test_withdraw "缺少金额参数" '{"password": "123456"}' "请求参数错误"

# 测试2: 缺少密码参数
test_withdraw "缺少密码参数" '{"amount": 100}' "请求参数错误"

# 测试3: 金额为0
test_withdraw "金额为0" '{"amount": 0, "password": "123456"}' "提现金额必须大于0"

# 测试4: 金额为负数
test_withdraw "金额为负数" '{"amount": -100, "password": "123456"}' "提现金额必须大于0"

# 测试5: 金额超过限额
test_withdraw "金额超过限额" '{"amount": 2000000, "password": "123456"}' "单笔提现金额不能超过100万元"

# 测试6: 密码错误
test_withdraw "密码错误" '{"amount": 100, "password": "wrong_password"}' "登录密码错误"

# 测试7: 余额不足（假设用户余额为0）
test_withdraw "余额不足" '{"amount": 1000, "password": "123456"}' "余额不足"

# 测试8: 正常提现（小额测试）
test_withdraw "正常提现" '{"amount": 10, "password": "123456"}' ""

echo -e "\n${GREEN}=== 测试完成 ===${NC}"
echo -e "${YELLOW}注意: 这些测试验证了优化后的提现接口的简化参数和增强校验功能${NC}" 