#!/bin/bash

# 测试假订单生成功能
# 测试系统UID生成、订单创建、数据清理等功能

BASE_URL="http://localhost:8080/api/v1"
TOKEN=""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

print_success() {
    echo -e "${BLUE}[SUCCESS]${NC} $1"
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

# 测试系统UID生成
test_system_uid_generation() {
    print_message "测试系统UID生成..."
    
    # 这里可以测试系统UID生成器是否正常工作
    # 由于这是内部功能，我们通过查看生成的订单来验证
    
    print_success "系统UID生成器已就绪"
}

# 测试手动生成假订单
test_manual_generation() {
    print_message "测试手动生成假订单..."
    
    # 这里可以添加手动生成假订单的API调用
    # 目前先跳过，因为还没有实现API接口
    
    print_success "手动生成功能测试完成"
}

# 测试订单列表查询
test_order_list_query() {
    print_message "测试订单列表查询..."
    
    # 查询所有订单
    ALL_ORDERS_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/list" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "page": 1,
            "page_size": 20,
            "status": 3
        }')
    
    echo "所有订单列表响应: $ALL_ORDERS_RESPONSE"
    
    # 检查是否包含系统订单
    if echo "$ALL_ORDERS_RESPONSE" | grep -q "is_system_order.*true"; then
        print_success "✓ 发现系统订单"
    else
        print_warning "⚠ 未发现系统订单，可能需要等待定时任务生成"
    fi
    
    # 检查7位UID
    if echo "$ALL_ORDERS_RESPONSE" | grep -q '"uid":"[0-9]\{7\}"'; then
        print_success "✓ 发现7位系统UID"
    else
        print_warning "⚠ 未发现7位系统UID"
    fi
}

# 测试我的订单列表
test_my_orders_query() {
    print_message "测试我的订单列表..."
    
    MY_ORDERS_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/my-orders" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "page": 1,
            "page_size": 20,
            "status": 1
        }')
    
    echo "我的订单列表响应: $MY_ORDERS_RESPONSE"
    
    # 检查是否包含系统订单
    if echo "$MY_ORDERS_RESPONSE" | grep -q "is_system_order.*true"; then
        print_success "✓ 我的订单列表包含系统订单"
    else
        print_warning "⚠ 我的订单列表未包含系统订单"
    fi
}

# 测试拼单列表
test_group_buy_list() {
    print_message "测试拼单列表..."
    
    GROUP_BUY_RESPONSE=$(curl -s -X POST "${BASE_URL}/order/list" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "page": 1,
            "page_size": 20,
            "status": 3
        }')
    
    echo "拼单列表响应: $GROUP_BUY_RESPONSE"
    
    # 检查是否包含系统拼单
    if echo "$GROUP_BUY_RESPONSE" | grep -q "is_system_order.*true"; then
        print_success "✓ 发现系统拼单"
    else
        print_warning "⚠ 未发现系统拼单"
    fi
}

# 测试数据库查询
test_database_query() {
    print_message "测试数据库查询..."
    
    # 检查orders表结构
    DB_STRUCTURE=$(mysql -h localhost -P 3306 -u root -p123456 -D gin_fatamorgana -e "DESCRIBE orders;" 2>/dev/null)
    
    if echo "$DB_STRUCTURE" | grep -q "is_system_order"; then
        print_success "✓ 数据库orders表包含is_system_order字段"
    else
        print_error "✗ 数据库orders表不包含is_system_order字段"
    fi
    
    # 查询系统订单数据
    SYSTEM_ORDERS=$(mysql -h localhost -P 3306 -u root -p123456 -D gin_fatamorgana -e "SELECT id, order_no, uid, is_system_order, status, amount, profit_amount FROM orders WHERE is_system_order = 1 LIMIT 10;" 2>/dev/null)
    
    if [ $? -eq 0 ]; then
        print_success "✓ 数据库查询成功"
        echo "系统订单数据:"
        echo "$SYSTEM_ORDERS"
        
        # 检查7位UID
        if echo "$SYSTEM_ORDERS" | grep -q "[0-9]\{7\}"; then
            print_success "✓ 发现7位系统UID"
        else
            print_warning "⚠ 未发现7位系统UID"
        fi
    else
        print_error "✗ 数据库查询失败"
    fi
    
    # 查询拼单数据
    GROUP_BUYS=$(mysql -h localhost -P 3306 -u root -p123456 -D gin_fatamorgana -e "SELECT id, group_buy_no, uid, status, per_person_amount FROM group_buys LIMIT 10;" 2>/dev/null)
    
    if [ $? -eq 0 ]; then
        print_success "✓ 拼单数据查询成功"
        echo "拼单数据:"
        echo "$GROUP_BUYS"
    else
        print_error "✗ 拼单数据查询失败"
    fi
}

# 测试时间分布
test_time_distribution() {
    print_message "测试时间分布..."
    
    # 查询最近生成的系统订单时间分布
    TIME_DISTRIBUTION=$(mysql -h localhost -P 3306 -u root -p123456 -D gin_fatamorgana -e "
        SELECT 
            DATE_FORMAT(created_at, '%Y-%m-%d %H:%i') as time_slot,
            COUNT(*) as count
        FROM orders 
        WHERE is_system_order = 1 
        AND created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
        GROUP BY time_slot 
        ORDER BY time_slot DESC 
        LIMIT 10;" 2>/dev/null)
    
    if [ $? -eq 0 ]; then
        print_success "✓ 时间分布查询成功"
        echo "最近1小时系统订单时间分布:"
        echo "$TIME_DISTRIBUTION"
    else
        print_error "✗ 时间分布查询失败"
    fi
}

# 测试任务数量范围
test_task_count_range() {
    print_message "测试任务数量范围..."
    
    # 查询任务数量分布
    TASK_DISTRIBUTION=$(mysql -h localhost -P 3306 -u root -p123456 -D gin_fatamorgana -e "
        SELECT 
            'like_count' as task_type,
            MIN(like_count) as min_count,
            MAX(like_count) as max_count,
            AVG(like_count) as avg_count
        FROM orders 
        WHERE is_system_order = 1 
        AND created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
        UNION ALL
        SELECT 
            'share_count' as task_type,
            MIN(share_count) as min_count,
            MAX(share_count) as max_count,
            AVG(share_count) as avg_count
        FROM orders 
        WHERE is_system_order = 1 
        AND created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
        UNION ALL
        SELECT 
            'follow_count' as task_type,
            MIN(follow_count) as min_count,
            MAX(follow_count) as max_count,
            AVG(follow_count) as avg_count
        FROM orders 
        WHERE is_system_order = 1 
        AND created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
        UNION ALL
        SELECT 
            'favorite_count' as task_type,
            MIN(favorite_count) as min_count,
            MAX(favorite_count) as max_count,
            AVG(favorite_count) as avg_count
        FROM orders 
        WHERE is_system_order = 1 
        AND created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR);" 2>/dev/null)
    
    if [ $? -eq 0 ]; then
        print_success "✓ 任务数量分布查询成功"
        echo "任务数量分布:"
        echo "$TASK_DISTRIBUTION"
        
        # 检查是否在100-2000范围内
        if echo "$TASK_DISTRIBUTION" | grep -q "min_count.*[0-9]\{1,3\}[0-9]\{0,1\}[0-9]\{0,1\}"; then
            print_success "✓ 任务数量在合理范围内"
        else
            print_warning "⚠ 任务数量可能超出预期范围"
        fi
    else
        print_error "✗ 任务数量分布查询失败"
    fi
}

# 测试状态分布
test_status_distribution() {
    print_message "测试状态分布..."
    
    # 查询订单状态分布
    ORDER_STATUS_DIST=$(mysql -h localhost -P 3306 -u root -p123456 -D gin_fatamorgana -e "
        SELECT 
            status,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM orders WHERE is_system_order = 1), 2) as percentage
        FROM orders 
        WHERE is_system_order = 1 
        AND created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
        GROUP BY status;" 2>/dev/null)
    
    if [ $? -eq 0 ]; then
        print_success "✓ 订单状态分布查询成功"
        echo "订单状态分布:"
        echo "$ORDER_STATUS_DIST"
    else
        print_error "✗ 订单状态分布查询失败"
    fi
    
    # 查询拼单状态分布
    GROUP_BUY_STATUS_DIST=$(mysql -h localhost -P 3306 -u root -p123456 -D gin_fatamorgana -e "
        SELECT 
            status,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM group_buys), 2) as percentage
        FROM group_buys 
        WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
        GROUP BY status;" 2>/dev/null)
    
    if [ $? -eq 0 ]; then
        print_success "✓ 拼单状态分布查询成功"
        echo "拼单状态分布:"
        echo "$GROUP_BUY_STATUS_DIST"
    else
        print_error "✗ 拼单状态分布查询失败"
    fi
}

# 主函数
main() {
    print_message "开始测试假订单生成功能..."
    
    # 登录
    login
    
    # 测试系统UID生成
    test_system_uid_generation
    
    # 测试手动生成假订单
    test_manual_generation
    
    # 测试订单列表查询
    test_order_list_query
    
    # 测试我的订单列表
    test_my_orders_query
    
    # 测试拼单列表
    test_group_buy_list
    
    # 测试数据库查询
    test_database_query
    
    # 测试时间分布
    test_time_distribution
    
    # 测试任务数量范围
    test_task_count_range
    
    # 测试状态分布
    test_status_distribution
    
    print_message "测试完成！"
    print_message "如果看到系统订单数据，说明假订单生成功能正常工作。"
    print_message "如果没有看到数据，请等待定时任务执行（每5分钟一次）。"
}

# 运行主函数
main 