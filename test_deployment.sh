#!/bin/bash

# 部署验证脚本
echo "🔍 验证部署状态..."

# 设置基础URL
BASE_URL="http://localhost:8080"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 测试函数
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    local token=$5
    
    echo -e "\n${YELLOW}测试: $description${NC}"
    
    if [ "$method" = "GET" ]; then
        if [ -n "$token" ]; then
            response=$(curl -s -w "\n%{http_code}" -H "Authorization: Bearer $token" "$BASE_URL$endpoint")
        else
            response=$(curl -s -w "\n%{http_code}" "$BASE_URL$endpoint")
        fi
    else
        if [ -n "$token" ]; then
            response=$(curl -s -w "\n%{http_code}" -X "$method" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $token" \
                -d "$data" \
                "$BASE_URL$endpoint")
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" \
                -H "Content-Type: application/json" \
                -d "$data" \
                "$BASE_URL$endpoint")
        fi
    fi
    
    # 分离响应体和状态码
    body=$(echo "$response" | head -n -1)
    status_code=$(echo "$response" | tail -n 1)
    
    if [ "$status_code" = "200" ] || [ "$status_code" = "201" ]; then
        echo -e "${GREEN}✅ 成功 (HTTP $status_code)${NC}"
        echo "响应: $body" | jq '.' 2>/dev/null || echo "响应: $body"
        return 0
    else
        echo -e "${RED}❌ 失败 (HTTP $status_code)${NC}"
        echo "响应: $body" | jq '.' 2>/dev/null || echo "响应: $body"
        return 1
    fi
}

# 检查Docker服务状态
check_docker_services() {
    echo -e "\n${BLUE}=== 检查Docker服务状态 ===${NC}"
    
    if command -v docker-compose &> /dev/null; then
        echo "Docker Compose服务状态:"
        docker-compose ps
        
        # 检查容器健康状态
        echo -e "\n容器健康状态:"
        for service in app mysql redis nginx; do
            if docker-compose ps | grep -q "$service.*Up"; then
                echo -e "${GREEN}✅ $service 运行正常${NC}"
            else
                echo -e "${RED}❌ $service 未运行或异常${NC}"
            fi
        done
    else
        echo -e "${RED}❌ Docker Compose未安装${NC}"
        return 1
    fi
}

# 检查网络连接
check_network() {
    echo -e "\n${BLUE}=== 检查网络连接 ===${NC}"
    
    # 检查端口是否开放
    for port in 80 8080 3306 6379; do
        if netstat -tuln 2>/dev/null | grep -q ":$port "; then
            echo -e "${GREEN}✅ 端口 $port 已开放${NC}"
        else
            echo -e "${RED}❌ 端口 $port 未开放${NC}"
        fi
    done
}

# 检查应用健康状态
check_application_health() {
    echo -e "\n${BLUE}=== 检查应用健康状态 ===${NC}"
    
    # 基础健康检查
    test_endpoint "GET" "/health" "" "基础健康检查"
    
    # 系统健康检查
    test_endpoint "GET" "/health/check" "" "系统健康检查"
    
    # 数据库健康检查
    test_endpoint "GET" "/health/database" "" "数据库健康检查"
    
    # Redis健康检查
    test_endpoint "GET" "/health/redis" "" "Redis健康检查"
    
    # 数据库统计
    test_endpoint "GET" "/health/stats" "" "数据库统计"
    
    # 查询统计
    test_endpoint "GET" "/health/query-stats" "" "查询统计"
}

# 测试用户注册和登录
test_user_authentication() {
    echo -e "\n${BLUE}=== 测试用户认证 ===${NC}"
    
    # 注册测试用户
    test_endpoint "POST" "/auth/register" '{
        "email": "deployment_test@example.com",
        "password": "123456",
        "confirm_password": "123456",
        "invite_code": "ADMIN1"
    }' "用户注册"
    
    # 用户登录
    login_response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "email": "deployment_test@example.com",
            "password": "123456"
        }' \
        "$BASE_URL/auth/login")
    
    token=$(echo "$login_response" | jq -r '.data.access_token')
    
    if [ "$token" != "null" ] && [ "$token" != "" ]; then
        echo -e "${GREEN}✅ 登录成功，获取到token${NC}"
        
        # 获取用户信息
        test_endpoint "GET" "/auth/profile" "" "获取用户信息" "$token"
        
        # 获取会话信息
        test_endpoint "GET" "/session/status" "" "检查登录状态" "$token"
        
        return 0
    else
        echo -e "${RED}❌ 登录失败${NC}"
        return 1
    fi
}

# 测试钱包功能
test_wallet_functionality() {
    echo -e "\n${BLUE}=== 测试钱包功能 ===${NC}"
    
    # 获取登录token
    login_response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "email": "deployment_test@example.com",
            "password": "123456"
        }' \
        "$BASE_URL/auth/login")
    
    token=$(echo "$login_response" | jq -r '.data.access_token')
    
    if [ "$token" != "null" ] && [ "$token" != "" ]; then
        # 获取钱包信息
        test_endpoint "GET" "/wallet/info" "" "获取钱包信息" "$token"
        
        # 获取交易记录
        test_endpoint "GET" "/wallet/transactions" "" "获取交易记录" "$token"
        
        # 获取提现汇总
        test_endpoint "GET" "/wallet/withdraw-summary" "" "获取提现汇总" "$token"
    else
        echo -e "${RED}❌ 无法获取token，跳过钱包测试${NC}"
    fi
}

# 测试银行卡功能
test_bank_card_functionality() {
    echo -e "\n${BLUE}=== 测试银行卡功能 ===${NC}"
    
    # 获取登录token
    login_response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "email": "deployment_test@example.com",
            "password": "123456"
        }' \
        "$BASE_URL/auth/login")
    
    token=$(echo "$login_response" | jq -r '.data.access_token')
    
    if [ "$token" != "null" ] && [ "$token" != "" ]; then
        # 绑定银行卡
        test_endpoint "POST" "/auth/bind-bank-card" '{
            "uid": "12345678",
            "card_number": "6225881234567890",
            "card_type": "借记卡",
            "bank_name": "招商银行",
            "card_holder": "张三"
        }' "绑定银行卡" "$token"
        
        # 获取银行卡信息
        test_endpoint "GET" "/auth/bank-card" "" "获取银行卡信息" "$token"
    else
        echo -e "${RED}❌ 无法获取token，跳过银行卡测试${NC}"
    fi
}

# 性能测试
performance_test() {
    echo -e "\n${BLUE}=== 性能测试 ===${NC}"
    
    # 测试健康检查接口的响应时间
    echo "测试健康检查接口响应时间..."
    start_time=$(date +%s%N)
    curl -s -o /dev/null "$BASE_URL/health"
    end_time=$(date +%s%N)
    response_time=$(( (end_time - start_time) / 1000000 ))
    echo -e "${GREEN}响应时间: ${response_time}ms${NC}"
    
    # 并发测试（简单版本）
    echo "执行简单并发测试..."
    for i in {1..10}; do
        curl -s -o /dev/null "$BASE_URL/health" &
    done
    wait
    echo -e "${GREEN}并发测试完成${NC}"
}

# 安全检查
security_check() {
    echo -e "\n${BLUE}=== 安全检查 ===${NC}"
    
    # 测试未认证访问
    test_endpoint "GET" "/auth/profile" "" "未认证访问测试"
    
    # 测试无效token
    test_endpoint "GET" "/auth/profile" "" "无效token测试" "invalid_token"
    
    # 测试SQL注入防护
    test_endpoint "POST" "/auth/register" '{
        "email": "test@example.com\"; DROP TABLE users; --",
        "password": "123456",
        "confirm_password": "123456",
        "invite_code": "ADMIN1"
    }' "SQL注入防护测试"
}

# 主函数
main() {
    echo -e "${GREEN}==========================================${NC}"
    echo -e "${GREEN}        Gin-FataMorgana 部署验证${NC}"
    echo -e "${GREEN}==========================================${NC}"
    echo ""
    
    # 检查Docker服务
    check_docker_services
    
    # 检查网络连接
    check_network
    
    # 等待服务启动
    echo -e "\n${YELLOW}等待服务启动...${NC}"
    sleep 10
    
    # 检查应用健康状态
    check_application_health
    
    # 测试用户认证
    test_user_authentication
    
    # 测试钱包功能
    test_wallet_functionality
    
    # 测试银行卡功能
    test_bank_card_functionality
    
    # 性能测试
    performance_test
    
    # 安全检查
    security_check
    
    echo -e "\n${GREEN}==========================================${NC}"
    echo -e "${GREEN}           验证完成！${NC}"
    echo -e "${GREEN}==========================================${NC}"
    echo ""
    echo -e "${YELLOW}建议：${NC}"
    echo "1. 检查所有测试结果"
    echo "2. 查看应用日志: docker-compose logs app"
    echo "3. 监控系统资源使用情况"
    echo "4. 定期执行此验证脚本"
}

# 执行主函数
main "$@" 