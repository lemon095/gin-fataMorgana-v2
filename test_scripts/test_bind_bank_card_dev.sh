#!/bin/bash

# 测试开发阶段的银行卡绑定功能（已禁用Luhn校验）

BASE_URL="http://localhost:8080"
API_PREFIX="/api/v1"

echo "=== 测试开发阶段银行卡绑定功能 ==="
echo "注意: 已禁用Luhn算法校验，方便开发测试"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_bind_bank_card() {
    local test_name="$1"
    local request_data="$2"
    local expected_error="$3"
    
    echo -e "\n${YELLOW}测试: $test_name${NC}"
    echo "请求数据: $request_data"
    
    # 调用绑定银行卡接口
    response=$(curl -s -X POST "$BASE_URL$API_PREFIX/auth/bind-bank-card" \
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
            echo -e "${GREEN}✓ 测试通过: 银行卡绑定成功${NC}"
            return 0
        else
            echo -e "${RED}✗ 测试失败: 银行卡绑定失败${NC}"
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

# 测试1: 使用您之前提供的卡号（现在应该能通过）
test_bind_bank_card "使用之前失败的卡号" '{"bank_name": "招商银行", "card_holder": "张三", "card_number": "6222600234567890123", "card_type": "借记卡"}' ""

# 测试2: 使用其他银行的卡号
test_bind_bank_card "工商银行借记卡" '{"bank_name": "工商银行", "card_holder": "李四", "card_number": "6222021234567890123", "card_type": "借记卡"}' ""

# 测试3: 使用中国银行信用卡
test_bind_bank_card "中国银行信用卡" '{"bank_name": "中国银行", "card_holder": "王五", "card_number": "6227601234567890123", "card_type": "信用卡"}' ""

# 测试4: 使用储蓄卡
test_bind_bank_card "储蓄卡" '{"bank_name": "建设银行", "card_holder": "赵六", "card_number": "6217001234567890123", "card_type": "储蓄卡"}' ""

# 测试5: 验证参数验证仍然有效
test_bind_bank_card "缺少银行名称" '{"card_holder": "张三", "card_number": "6222600234567890123", "card_type": "借记卡"}' "银行名称不能为空"

# 测试6: 验证卡号长度检查仍然有效
test_bind_bank_card "卡号长度错误" '{"bank_name": "招商银行", "card_holder": "张三", "card_number": "123", "card_type": "借记卡"}' "银行卡号长度不正确"

echo -e "\n${GREEN}=== 测试完成 ===${NC}"
echo -e "${YELLOW}注意: 开发阶段已禁用Luhn算法校验，生产环境需要重新启用${NC}" 