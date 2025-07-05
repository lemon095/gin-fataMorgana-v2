#!/bin/bash

# 水单生成调试脚本
# 用于测试和调试水单生成功能

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
BASE_URL="http://localhost:9001"
LOGIN_URL="${BASE_URL}/api/v1/auth/login"
CRON_GENERATE_URL="${BASE_URL}/api/v1/cron/manual-generate"
CRON_STATUS_URL="${BASE_URL}/api/v1/cron/status"

# 测试账号（需要先注册）
TEST_EMAIL="test_water@example.com"
TEST_PASSWORD="test123456"

# 全局变量
TOKEN=""

# 打印消息函数
print_message() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# 检查服务状态
check_service_status() {
    print_message "检查服务状态..."
    
    HEALTH_RESPONSE=$(curl -s -X GET "${BASE_URL}/health")
    
    if echo "$HEALTH_RESPONSE" | jq -e '.code == 0' > /dev/null; then
        print_success "服务运行正常"
    else
        print_error "服务未运行或异常"
        echo "请确保服务已启动: go run main.go"
        exit 1
    fi
}

# 注册测试账号
register_test_user() {
    print_message "注册测试账号..."
    
    REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"account\": \"$TEST_EMAIL\",
            \"password\": \"$TEST_PASSWORD\",
            \"invite_code\": \"TEST123\"
        }")
    
    echo "注册响应: $REGISTER_RESPONSE"
    
    if echo "$REGISTER_RESPONSE" | jq -e '.code == 0' > /dev/null; then
        print_success "注册成功"
    else
        print_warning "注册失败或用户已存在"
    fi
}

# 登录函数
login() {
    print_message "开始登录..."
    
    LOGIN_RESPONSE=$(curl -s -X POST "$LOGIN_URL" \
        -H "Content-Type: application/json" \
        -d "{
            \"account\": \"$TEST_EMAIL\",
            \"password\": \"$TEST_PASSWORD\"
        }")
    
    echo "登录响应: $LOGIN_RESPONSE"
    
    # 检查登录是否成功
    if echo "$LOGIN_RESPONSE" | jq -e '.code == 0' > /dev/null; then
        TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')
        print_success "登录成功，获取到token"
        echo "Token: ${TOKEN:0:20}..."
    else
        print_error "登录失败"
        echo "错误信息: $(echo "$LOGIN_RESPONSE" | jq -r '.message')"
        exit 1
    fi
}

# 测试手动生成水单
test_manual_generate_water_orders() {
    print_message "测试手动生成水单功能..."
    
    # 生成少量水单进行测试
    GENERATE_RESPONSE=$(curl -s -X POST "$CRON_GENERATE_URL" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "count": 5
        }')
    
    echo "生成水单响应: $GENERATE_RESPONSE"
    
    # 检查是否成功生成
    if echo "$GENERATE_RESPONSE" | jq -e '.code == 0' > /dev/null; then
        print_success "手动生成水单成功"
        
        # 获取生成统计
        TOTAL_GENERATED=$(echo "$GENERATE_RESPONSE" | jq -r '.data.total_generated')
        PURCHASE_ORDERS=$(echo "$GENERATE_RESPONSE" | jq -r '.data.purchase_orders')
        GROUP_BUY_ORDERS=$(echo "$GENERATE_RESPONSE" | jq -r '.data.group_buy_orders')
        TOTAL_AMOUNT=$(echo "$GENERATE_RESPONSE" | jq -r '.data.total_amount')
        TOTAL_PROFIT=$(echo "$GENERATE_RESPONSE" | jq -r '.data.total_profit')
        
        echo -e "${GREEN}水单生成统计:${NC}"
        echo "  总数: $TOTAL_GENERATED"
        echo "  购买单: $PURCHASE_ORDERS"
        echo "  拼单: $GROUP_BUY_ORDERS"
        echo "  总金额: $TOTAL_AMOUNT"
        echo "  总利润: $TOTAL_PROFIT"
    else
        print_error "手动生成水单失败"
        echo "错误信息: $(echo "$GENERATE_RESPONSE" | jq -r '.message')"
    fi
}

# 获取定时任务状态
test_cron_status() {
    print_message "获取定时任务状态..."
    
    STATUS_RESPONSE=$(curl -s -X GET "$CRON_STATUS_URL" \
        -H "Authorization: Bearer $TOKEN")
    
    echo "状态响应: $STATUS_RESPONSE"
    
    # 检查是否成功获取状态
    if echo "$STATUS_RESPONSE" | jq -e '.code == 0' > /dev/null; then
        print_success "获取定时任务状态成功"
        
        # 显示状态信息
        echo -e "${GREEN}定时任务状态:${NC}"
        echo "$STATUS_RESPONSE" | jq -r '.data.cron_status | to_entries[] | "任务: \(.key), 下次执行: \(.value.next_run), 上次执行: \(.value.prev_run)"'
    else
        print_error "获取定时任务状态失败"
        echo "错误信息: $(echo "$STATUS_RESPONSE" | jq -r '.message')"
    fi
}

# 检查数据库中的水单数据
check_water_orders_in_db() {
    print_message "检查数据库中的水单数据..."
    
    # 检查系统订单数量
    SYSTEM_ORDERS_COUNT=$(mysql -h 127.0.0.1 -u root -proot -D future -e "SELECT COUNT(*) as count FROM orders WHERE is_system_order = 1;" -s -N 2>/dev/null)
    
    if [ $? -eq 0 ]; then
        print_success "数据库连接正常"
        echo "系统订单数量: $SYSTEM_ORDERS_COUNT"
        
        if [ "$SYSTEM_ORDERS_COUNT" -gt 0 ]; then
            print_success "发现 $SYSTEM_ORDERS_COUNT 条系统订单"
            
            # 显示最近的几条系统订单
            echo -e "${GREEN}最近的系统订单:${NC}"
            mysql -h 127.0.0.1 -u root -proot -D future -e "
                SELECT 
                    order_no, 
                    uid, 
                    amount, 
                    status, 
                    created_at,
                    is_system_order
                FROM orders 
                WHERE is_system_order = 1 
                ORDER BY created_at DESC 
                LIMIT 5;
            " 2>/dev/null
        else
            print_warning "数据库中没有系统订单"
        fi
    else
        print_error "数据库连接失败"
    fi
}

# 检查配置文件
check_config() {
    print_message "检查配置文件..."
    
    if [ -f "config.yaml" ]; then
        echo -e "${GREEN}配置文件中的水单配置:${NC}"
        echo "启用状态: $(grep -A 10 "fake_data:" config.yaml | grep "enabled:" | awk '{print $2}')"
        echo "订单生成表达式: $(grep -A 10 "fake_data:" config.yaml | grep "cron_expression:" | awk '{print $2}')"
        echo "最小订单数: $(grep -A 10 "fake_data:" config.yaml | grep "min_orders:" | awk '{print $2}')"
        echo "最大订单数: $(grep -A 10 "fake_data:" config.yaml | grep "max_orders:" | awk '{print $2}')"
        echo "购买单比例: $(grep -A 10 "fake_data:" config.yaml | grep "purchase_ratio:" | awk '{print $2}')"
    else
        print_warning "配置文件不存在"
    fi
}

# 主函数
main() {
    print_message "开始水单生成调试..."
    
    # 检查依赖
    if ! command -v jq &> /dev/null; then
        print_error "jq 命令未找到，请安装 jq"
        exit 1
    fi
    
    # 检查服务状态
    check_service_status
    
    # 检查配置文件
    check_config
    
    # 注册测试账号
    register_test_user
    
    # 登录
    login
    
    # 获取定时任务状态
    test_cron_status
    
    # 测试手动生成水单
    test_manual_generate_water_orders
    
    # 等待一下，让数据生成完成
    print_message "等待3秒让数据生成完成..."
    sleep 3
    
    # 检查数据库中的水单数据
    check_water_orders_in_db
    
    # 再次获取状态
    test_cron_status
    
    print_message "水单生成调试完成！"
    print_message "如果看到生成统计和数据库中有系统订单，说明水单功能正常工作。"
    print_message "如果没有数据，请检查服务日志中的详细错误信息。"
}

# 运行主函数
main 