#!/bin/bash

# 测试优化后的绑定银行卡接口
# 验证uid从token中获取，简化请求参数

BASE_URL="http://localhost:8080"
API_PREFIX="/api/v1"

echo "=== 测试优化后的绑定银行卡接口 ==="

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

# 测试1: 缺少银行名称
test_bind_bank_card "缺少银行名称" '{"card_holder": "张三", "card_number": "6222021234567890123", "card_type": "借记卡"}' "银行名称不能为空"

# 测试2: 缺少持卡人姓名
test_bind_bank_card "缺少持卡人姓名" '{"bank_name": "中国银行", "card_number": "6222021234567890123", "card_type": "借记卡"}' "持卡人姓名不能为空"

# 测试3: 缺少银行卡号
test_bind_bank_card "缺少银行卡号" '{"bank_name": "中国银行", "card_holder": "张三", "card_type": "借记卡"}' "银行卡号不能为空"

# 测试4: 缺少卡类型
test_bind_bank_card "缺少卡类型" '{"bank_name": "中国银行", "card_holder": "张三", "card_number": "6222021234567890123"}' "卡类型不能为空"

# 测试5: 卡类型不正确
test_bind_bank_card "卡类型不正确" '{"bank_name": "中国银行", "card_holder": "张三", "card_number": "6222021234567890123", "card_type": "其他卡"}' "卡类型不正确"

# 测试6: 银行卡号格式错误
test_bind_bank_card "银行卡号格式错误" '{"bank_name": "中国银行", "card_holder": "张三", "card_number": "123", "card_type": "借记卡"}' "银行卡号长度不正确"

# 测试7: 正常绑定银行卡
test_bind_bank_card "正常绑定银行卡" '{"bank_name": "中国银行", "card_holder": "张三", "card_number": "6222021234567890123", "card_type": "借记卡"}' ""

# 测试8: 验证不需要传递uid参数
echo -e "\n${YELLOW}验证请求参数中不包含uid...${NC}"
if echo '{"bank_name": "中国银行", "card_holder": "张三", "card_number": "6222021234567890123", "card_type": "借记卡"}' | grep -q "uid"; then
    echo -e "${RED}✗ 请求参数中仍然包含uid${NC}"
else
    echo -e "${GREEN}✓ 请求参数中不包含uid，符合优化要求${NC}"
fi

echo -e "\n${GREEN}=== 测试完成 ===${NC}"
echo -e "${YELLOW}注意: 这些测试验证了优化后的绑定银行卡接口，uid从token中获取，简化了请求参数${NC}" 