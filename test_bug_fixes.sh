#!/bin/bash

# 测试Bug修复效果
echo "🔧 测试Bug修复效果..."

# 设置基础URL
BASE_URL="http://localhost:8080"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "\n${YELLOW}测试: $description${NC}"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$BASE_URL$endpoint")
    fi
    
    # 分离响应体和状态码
    body=$(echo "$response" | head -n -1)
    status_code=$(echo "$response" | tail -n 1)
    
    if [ "$status_code" = "200" ] || [ "$status_code" = "201" ]; then
        echo -e "${GREEN}✅ 成功 (HTTP $status_code)${NC}"
        echo "响应: $body" | jq '.' 2>/dev/null || echo "响应: $body"
    else
        echo -e "${RED}❌ 失败 (HTTP $status_code)${NC}"
        echo "响应: $body" | jq '.' 2>/dev/null || echo "响应: $body"
    fi
}

# 1. 测试配置验证
echo -e "\n${YELLOW}=== 1. 测试配置验证 ===${NC}"
test_endpoint "GET" "/health/check" "" "系统健康检查"

# 2. 测试数据库连接
echo -e "\n${YELLOW}=== 2. 测试数据库连接 ===${NC}"
test_endpoint "GET" "/health/database" "" "数据库健康检查"

# 3. 测试Redis连接
echo -e "\n${YELLOW}=== 3. 测试Redis连接 ===${NC}"
test_endpoint "GET" "/health/redis" "" "Redis健康检查"

# 4. 测试用户注册（参数验证）
echo -e "\n${YELLOW}=== 4. 测试用户注册参数验证 ===${NC}"

# 测试空邮箱
test_endpoint "POST" "/auth/register" '{
    "email": "",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "TEST01"
}' "空邮箱验证"

# 测试短密码
test_endpoint "POST" "/auth/register" '{
    "email": "test@example.com",
    "password": "123",
    "confirm_password": "123",
    "invite_code": "TEST01"
}' "短密码验证"

# 测试密码不匹配
test_endpoint "POST" "/auth/register" '{
    "email": "test@example.com",
    "password": "123456",
    "confirm_password": "654321",
    "invite_code": "TEST01"
}' "密码不匹配验证"

# 5. 测试雪花算法
echo -e "\n${YELLOW}=== 5. 测试雪花算法 ===${NC}"
test_endpoint "POST" "/auth/register" '{
    "email": "snowflake_test@example.com",
    "password": "123456",
    "password_confirm": "123456",
    "invite_code": "TEST01"
}' "雪花算法UID生成"

# 6. 测试银行卡验证
echo -e "\n${YELLOW}=== 6. 测试银行卡验证 ===${NC}"

# 先注册一个用户
echo "注册测试用户..."
register_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "email": "bankcard_test@example.com",
        "password": "123456",
        "confirm_password": "123456",
        "invite_code": "TEST01"
    }' \
    "$BASE_URL/auth/register")

# 提取用户信息
user_data=$(echo "$register_response" | jq -r '.data.user')
uid=$(echo "$user_data" | jq -r '.uid')

if [ "$uid" != "null" ] && [ "$uid" != "" ]; then
    echo "用户注册成功，UID: $uid"
    
    # 获取登录token
    login_response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "email": "bankcard_test@example.com",
            "password": "123456"
        }' \
        "$BASE_URL/auth/login")
    
    token=$(echo "$login_response" | jq -r '.data.access_token')
    
    if [ "$token" != "null" ] && [ "$token" != "" ]; then
        echo "登录成功，获取到token"
        
        # 测试无效银行卡号
        test_endpoint "POST" "/auth/bind-bank-card" "{
            \"uid\": \"$uid\",
            \"card_number\": \"1234567890123456\",
            \"card_type\": \"借记卡\",
            \"bank_name\": \"测试银行\",
            \"card_holder\": \"张三\"
        }" "无效银行卡号验证" "$token"
        
        # 测试有效银行卡号
        test_endpoint "POST" "/auth/bind-bank-card" "{
            \"uid\": \"$uid\",
            \"card_number\": \"6225881234567890\",
            \"card_type\": \"借记卡\",
            \"bank_name\": \"招商银行\",
            \"card_holder\": \"张三\"
        }" "有效银行卡号验证" "$token"
    else
        echo "登录失败，无法获取token"
    fi
else
    echo "用户注册失败"
fi

# 7. 测试错误处理
echo -e "\n${YELLOW}=== 7. 测试错误处理 ===${NC}"

# 测试不存在的用户
test_endpoint "GET" "/wallet/info?uid=NONEXISTENT" "" "不存在的用户查询"

# 测试无效的token
test_endpoint "GET" "/auth/profile" "" "无效token验证" "Bearer invalid_token"

# 8. 测试并发安全性
echo -e "\n${YELLOW}=== 8. 测试并发安全性 ===${NC}"
echo "注意：并发测试需要在实际运行环境中进行压力测试"

# 9. 测试资源清理
echo -e "\n${YELLOW}=== 9. 测试资源清理 ===${NC}"
test_endpoint "GET" "/health/stats" "" "数据库连接池统计"

echo -e "\n${GREEN}✅ Bug修复测试完成${NC}"
echo -e "\n${YELLOW}建议：${NC}"
echo "1. 检查日志文件中的错误信息"
echo "2. 监控数据库连接池状态"
echo "3. 验证雪花算法的唯一性"
echo "4. 测试优雅关闭功能"
echo "5. 进行压力测试验证并发安全性" 