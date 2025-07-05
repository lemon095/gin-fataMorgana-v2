#!/bin/bash

# 测试新增的is_system_order字段功能
# 测试订单创建和查询时是否正确包含该字段

BASE_URL="http://localhost:8080/api/v1"
TOKEN=""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 登录获取token
login() {
    print_message "正在登录获取token..."
    
    LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "testuser",
            "password": "123456"
        }')
    
    TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$TOKEN" ]; then
        print_error "登录失败，无法获取token"
        echo "登录响应: $LOGIN_RESPONSE"
        exit 1
    fi
    
    print_message "登录成功，获取到token: ${TOKEN:0:20}..."
}

# 测试创建订单
test_create_order() {
    print_message "测试创建订单..."
    
    # 获取期数列表
    PERIOD_RESPONSE=$(curl -s -X GET "${BASE_URL}/order/period-list" \
        -H "Authorization: Bearer $TOKEN")
    
    PERIOD_NUMBER=$(echo $PERIOD_RESPONSE | grep -o '"period_number":"[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$PERIOD_NUMBER" ]; then
        print_error "无法获取期数信息"
        echo "期数响应: $PERIOD_RESPONSE"
        return 1
    fi
    
    print_message "使用期数: $PERIOD_NUMBER"
    
    # 创建订单
    CREATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/create" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"period_number\": \"$PERIOD_NUMBER\",
            \"amount\": 10.00,
            \"like_count\": 5,
            \"share_count\": 3,
            \"follow_count\": 2,
            \"favorite_count\": 1
        }")
    
    echo "创建订单响应: $CREATE_RESPONSE"
    
    # 检查是否包含is_system_order字段
    if echo "$CREATE_RESPONSE" | grep -q "is_system_order"; then
        print_message "✓ 创建订单响应包含is_system_order字段"
    else
        print_warning "⚠ 创建订单响应中未找到is_system_order字段"
    fi
}

# 测试获取订单列表
test_get_order_list() {
    print_message "测试获取订单列表..."
    
    LIST_RESPONSE=$(curl -s -X GET "${BASE_URL}/order/list?page=1&page_size=10&status=1" \
        -H "Authorization: Bearer $TOKEN")
    
    echo "订单列表响应: $LIST_RESPONSE"
    
    # 检查是否包含is_system_order字段
    if echo "$LIST_RESPONSE" | grep -q "is_system_order"; then
        print_message "✓ 订单列表响应包含is_system_order字段"
    else
        print_warning "⚠ 订单列表响应中未找到is_system_order字段"
    fi
}

# 测试获取我的订单列表
test_get_my_orders() {
    print_message "测试获取我的订单列表..."
    
    MY_ORDERS_RESPONSE=$(curl -s -X GET "${BASE_URL}/order/my-orders?page=1&page_size=10&status=1" \
        -H "Authorization: Bearer $TOKEN")
    
    echo "我的订单列表响应: $MY_ORDERS_RESPONSE"
    
    # 检查是否包含is_system_order字段
    if echo "$MY_ORDERS_RESPONSE" | grep -q "is_system_order"; then
        print_message "✓ 我的订单列表响应包含is_system_order字段"
    else
        print_warning "⚠ 我的订单列表响应中未找到is_system_order字段"
    fi
}

# 测试拼单列表
test_group_buy_list() {
    print_message "测试获取拼单列表..."
    
    GROUP_BUY_RESPONSE=$(curl -s -X GET "${BASE_URL}/order/list?page=1&page_size=10&status=3" \
        -H "Authorization: Bearer $TOKEN")
    
    echo "拼单列表响应: $GROUP_BUY_RESPONSE"
    
    # 检查是否包含is_system_order字段
    if echo "$GROUP_BUY_RESPONSE" | grep -q "is_system_order"; then
        print_message "✓ 拼单列表响应包含is_system_order字段"
    else
        print_warning "⚠ 拼单列表响应中未找到is_system_order字段"
    fi
}

# 测试数据库查询
test_database_query() {
    print_message "测试数据库查询..."
    
    # 检查orders表结构
    DB_STRUCTURE=$(mysql -h localhost -P 3306 -u root -p123456 -D gin_fatamorgana -e "DESCRIBE orders;" 2>/dev/null)
    
    if echo "$DB_STRUCTURE" | grep -q "is_system_order"; then
        print_message "✓ 数据库orders表包含is_system_order字段"
    else
        print_error "✗ 数据库orders表不包含is_system_order字段"
    fi
    
    # 查询订单数据
    ORDER_DATA=$(mysql -h localhost -P 3306 -u root -p123456 -D gin_fatamorgana -e "SELECT id, order_no, uid, is_system_order FROM orders LIMIT 5;" 2>/dev/null)
    
    if [ $? -eq 0 ]; then
        print_message "✓ 数据库查询成功"
        echo "订单数据:"
        echo "$ORDER_DATA"
    else
        print_error "✗ 数据库查询失败"
    fi
}

# 主函数
main() {
    print_message "开始测试is_system_order字段功能..."
    
    # 登录
    login
    
    # 测试创建订单
    test_create_order
    
    # 测试获取订单列表
    test_get_order_list
    
    # 测试获取我的订单列表
    test_get_my_orders
    
    # 测试拼单列表
    test_group_buy_list
    
    # 测试数据库查询
    test_database_query
    
    print_message "测试完成！"
}

# 运行主函数
main 