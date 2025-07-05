#!/bin/bash

# 定时任务服务测试脚本
# 测试手动生成订单、清理数据和获取状态功能

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
CRON_CLEANUP_URL="${BASE_URL}/api/v1/cron/manual-cleanup"
CRON_STATUS_URL="${BASE_URL}/api/v1/cron/status"

# 测试账号（需要先注册）
TEST_EMAIL="test_cron@example.com"
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

# 测试手动生成订单
test_manual_generate() {
    print_message "测试手动生成订单功能..."
    
    GENERATE_RESPONSE=$(curl -s -X POST "$CRON_GENERATE_URL" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "count": 10
        }')
    
    echo "生成订单响应: $GENERATE_RESPONSE"
    
    # 检查是否成功生成
    if echo "$GENERATE_RESPONSE" | jq -e '.code == 0' > /dev/null; then
        print_success "手动生成订单成功"
        
        # 获取生成统计
        TOTAL_GENERATED=$(echo "$GENERATE_RESPONSE" | jq -r '.data.total_generated')
        PURCHASE_ORDERS=$(echo "$GENERATE_RESPONSE" | jq -r '.data.purchase_orders')
        GROUP_BUY_ORDERS=$(echo "$GENERATE_RESPONSE" | jq -r '.data.group_buy_orders')
        TOTAL_AMOUNT=$(echo "$GENERATE_RESPONSE" | jq -r '.data.total_amount')
        TOTAL_PROFIT=$(echo "$GENERATE_RESPONSE" | jq -r '.data.total_profit')
        
        echo -e "${GREEN}生成统计:${NC}"
        echo "  总数: $TOTAL_GENERATED"
        echo "  购买单: $PURCHASE_ORDERS"
        echo "  拼单: $GROUP_BUY_ORDERS"
        echo "  总金额: $TOTAL_AMOUNT"
        echo "  总利润: $TOTAL_PROFIT"
    else
        print_error "手动生成订单失败"
        echo "错误信息: $(echo "$GENERATE_RESPONSE" | jq -r '.message')"
    fi
}

# 测试获取定时任务状态
test_cron_status() {
    print_message "测试获取定时任务状态..."
    
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

# 测试手动清理数据
test_manual_cleanup() {
    print_message "测试手动清理数据功能..."
    
    CLEANUP_RESPONSE=$(curl -s -X POST "$CRON_CLEANUP_URL" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN")
    
    echo "清理数据响应: $CLEANUP_RESPONSE"
    
    # 检查是否成功清理
    if echo "$CLEANUP_RESPONSE" | jq -e '.code == 0' > /dev/null; then
        print_success "手动清理数据成功"
        
        # 获取清理统计
        DELETED_ORDERS=$(echo "$CLEANUP_RESPONSE" | jq -r '.data.deleted_orders')
        DELETED_GROUP_BUYS=$(echo "$CLEANUP_RESPONSE" | jq -r '.data.deleted_group_buys')
        CLEANUP_TIME=$(echo "$CLEANUP_RESPONSE" | jq -r '.data.cleanup_time')
        
        echo -e "${GREEN}清理统计:${NC}"
        echo "  删除订单: $DELETED_ORDERS"
        echo "  删除拼单: $DELETED_GROUP_BUYS"
        echo "  清理耗时: $CLEANUP_TIME"
    else
        print_error "手动清理数据失败"
        echo "错误信息: $(echo "$CLEANUP_RESPONSE" | jq -r '.message')"
    fi
}

# 测试定时任务配置
test_cron_config() {
    print_message "检查定时任务配置..."
    
    # 检查配置文件中的定时任务配置
    if [ -f "config.yaml" ]; then
        echo -e "${GREEN}配置文件中的定时任务配置:${NC}"
        echo "启用状态: $(grep -A 10 "fake_data:" config.yaml | grep "enabled:" | awk '{print $2}')"
        echo "订单生成表达式: $(grep -A 10 "fake_data:" config.yaml | grep "cron_expression:" | awk '{print $2}')"
        echo "数据清理表达式: $(grep -A 10 "fake_data:" config.yaml | grep "cleanup_cron:" | awk '{print $2}')"
        echo "最小订单数: $(grep -A 10 "fake_data:" config.yaml | grep "min_orders:" | awk '{print $2}')"
        echo "最大订单数: $(grep -A 10 "fake_data:" config.yaml | grep "max_orders:" | awk '{print $2}')"
        echo "购买单比例: $(grep -A 10 "fake_data:" config.yaml | grep "purchase_ratio:" | awk '{print $2}')"
        echo "保留天数: $(grep -A 10 "fake_data:" config.yaml | grep "retention_days:" | awk '{print $2}')"
    else
        print_warning "配置文件不存在"
    fi
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

# 主函数
main() {
    print_message "开始测试定时任务服务功能..."
    
    # 检查依赖
    if ! command -v jq &> /dev/null; then
        print_error "jq 命令未找到，请安装 jq"
        exit 1
    fi
    
    # 检查服务状态
    check_service_status
    
    # 测试定时任务配置
    test_cron_config
    
    # 登录
    login
    
    # 测试获取定时任务状态
    test_cron_status
    
    # 测试手动生成订单
    test_manual_generate
    
    # 等待一下，让数据生成完成
    print_message "等待5秒让数据生成完成..."
    sleep 5
    
    # 再次测试获取状态
    test_cron_status
    
    # 测试手动清理数据（可选，谨慎使用）
    read -p "是否要测试手动清理数据功能？(y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        test_manual_cleanup
    else
        print_warning "跳过手动清理数据测试"
    fi
    
    print_message "定时任务服务测试完成！"
    print_message "如果看到生成统计和状态信息，说明定时任务功能正常工作。"
}

# 运行主函数
main 