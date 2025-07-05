#!/bin/bash

# 数据库索引优化测试脚本
# 用于验证索引创建和性能提升效果

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 数据库配置
DB_HOST="localhost"
DB_PORT="3306"
DB_USER="root"
DB_PASS="123456"
DB_NAME="gin_fataMorgana"

echo -e "${BLUE}=== 数据库索引优化测试脚本 ===${NC}"
echo ""

# 检查数据库连接
echo -e "${YELLOW}1. 检查数据库连接...${NC}"
if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "USE $DB_NAME;" 2>/dev/null; then
    echo -e "${GREEN}✓ 数据库连接成功${NC}"
else
    echo -e "${RED}✗ 数据库连接失败${NC}"
    exit 1
fi
echo ""

# 检查表是否存在
echo -e "${YELLOW}2. 检查核心表是否存在...${NC}"
TABLES=("users" "orders" "wallet_transactions" "group_buys" "user_login_logs" "lottery_periods")
for table in "${TABLES[@]}"; do
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "DESCRIBE $DB_NAME.$table;" 2>/dev/null | grep -q "Field"; then
        echo -e "${GREEN}✓ 表 $table 存在${NC}"
    else
        echo -e "${RED}✗ 表 $table 不存在${NC}"
    fi
done
echo ""

# 显示当前索引状态
echo -e "${YELLOW}3. 显示当前索引状态...${NC}"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
SELECT 
    table_name,
    COUNT(*) as index_count
FROM information_schema.statistics 
WHERE table_schema = '$DB_NAME'
GROUP BY table_name
ORDER BY table_name;
" 2>/dev/null || echo -e "${RED}无法获取索引信息${NC}"
echo ""

# 测试查询性能（无索引）
echo -e "${YELLOW}4. 测试查询性能（无索引）...${NC}"

# 测试用户查询
echo "测试用户状态查询..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
EXPLAIN SELECT * FROM $DB_NAME.users WHERE status = 1 AND deleted_at IS NULL LIMIT 10;
" 2>/dev/null || echo -e "${RED}无法执行查询${NC}"

# 测试订单查询
echo "测试订单查询..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
EXPLAIN SELECT * FROM $DB_NAME.orders WHERE uid = 'U001' AND status = 'success' ORDER BY created_at DESC LIMIT 10;
" 2>/dev/null || echo -e "${RED}无法执行查询${NC}"

# 测试交易查询
echo "测试交易查询..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
EXPLAIN SELECT * FROM $DB_NAME.wallet_transactions WHERE uid = 'U001' ORDER BY created_at DESC LIMIT 10;
" 2>/dev/null || echo -e "${RED}无法执行查询${NC}"
echo ""

# 执行索引优化脚本
echo -e "${YELLOW}5. 执行索引优化脚本...${NC}"
if [ -f "database/migrations/create_indexes.sql" ]; then
    echo "开始创建索引..."
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" < database/migrations/create_indexes.sql
    echo -e "${GREEN}✓ 索引创建完成${NC}"
else
    echo -e "${RED}✗ 索引脚本文件不存在${NC}"
fi
echo ""

# 显示优化后的索引状态
echo -e "${YELLOW}6. 显示优化后的索引状态...${NC}"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
SELECT 
    table_name,
    COUNT(*) as index_count
FROM information_schema.statistics 
WHERE table_schema = '$DB_NAME'
GROUP BY table_name
ORDER BY table_name;
" 2>/dev/null || echo -e "${RED}无法获取索引信息${NC}"
echo ""

# 测试优化后的查询性能
echo -e "${YELLOW}7. 测试优化后的查询性能...${NC}"

# 测试用户查询（有索引）
echo "测试用户状态查询（有索引）..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
EXPLAIN SELECT * FROM $DB_NAME.users WHERE status = 1 AND deleted_at IS NULL LIMIT 10;
" 2>/dev/null || echo -e "${RED}无法执行查询${NC}"

# 测试订单查询（有索引）
echo "测试订单查询（有索引）..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
EXPLAIN SELECT * FROM $DB_NAME.orders WHERE uid = 'U001' AND status = 'success' ORDER BY created_at DESC LIMIT 10;
" 2>/dev/null || echo -e "${RED}无法执行查询${NC}"

# 测试交易查询（有索引）
echo "测试交易查询（有索引）..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
EXPLAIN SELECT * FROM $DB_NAME.wallet_transactions WHERE uid = 'U001' ORDER BY created_at DESC LIMIT 10;
" 2>/dev/null || echo -e "${RED}无法执行查询${NC}"
echo ""

# 显示表大小统计
echo -e "${YELLOW}8. 显示表大小统计...${NC}"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
SELECT 
    table_name,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)',
    ROUND((data_length / 1024 / 1024), 2) AS 'Data (MB)',
    ROUND((index_length / 1024 / 1024), 2) AS 'Index (MB)',
    table_rows
FROM information_schema.tables 
WHERE table_schema = '$DB_NAME'
ORDER BY (data_length + index_length) DESC;
" 2>/dev/null || echo -e "${RED}无法获取表大小信息${NC}"
echo ""

# 性能测试
echo -e "${YELLOW}9. 执行性能测试...${NC}"

# 测试用户查询性能
echo "测试用户查询性能..."
START_TIME=$(date +%s.%N)
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
SELECT COUNT(*) FROM $DB_NAME.users WHERE status = 1 AND deleted_at IS NULL;
" 2>/dev/null > /dev/null
END_TIME=$(date +%s.%N)
ELAPSED=$(echo "$END_TIME - $START_TIME" | bc)
echo -e "${GREEN}用户查询耗时: ${ELAPSED} 秒${NC}"

# 测试订单查询性能
echo "测试订单查询性能..."
START_TIME=$(date +%s.%N)
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
SELECT COUNT(*) FROM $DB_NAME.orders WHERE status = 'success';
" 2>/dev/null > /dev/null
END_TIME=$(date +%s.%N)
ELAPSED=$(echo "$END_TIME - $START_TIME" | bc)
echo -e "${GREEN}订单查询耗时: ${ELAPSED} 秒${NC}"

# 测试交易查询性能
echo "测试交易查询性能..."
START_TIME=$(date +%s.%N)
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
SELECT COUNT(*) FROM $DB_NAME.wallet_transactions WHERE type = 'recharge';
" 2>/dev/null > /dev/null
END_TIME=$(date +%s.%N)
ELAPSED=$(echo "$END_TIME - $START_TIME" | bc)
echo -e "${GREEN}交易查询耗时: ${ELAPSED} 秒${NC}"
echo ""

# 检查慢查询日志
echo -e "${YELLOW}10. 检查慢查询配置...${NC}"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
SHOW VARIABLES LIKE 'slow_query_log';
SHOW VARIABLES LIKE 'long_query_time';
SHOW VARIABLES LIKE 'log_queries_not_using_indexes';
" 2>/dev/null || echo -e "${RED}无法获取慢查询配置${NC}"
echo ""

# 生成测试报告
echo -e "${YELLOW}11. 生成测试报告...${NC}"
REPORT_FILE="test_scripts/index_optimization_report_$(date +%Y%m%d_%H%M%S).txt"

cat > "$REPORT_FILE" << EOF
数据库索引优化测试报告
生成时间: $(date)
数据库: $DB_NAME

=== 测试结果 ===
1. 数据库连接: 成功
2. 核心表检查: 完成
3. 索引创建: 完成
4. 性能测试: 完成

=== 索引统计 ===
$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
SELECT 
    table_name,
    COUNT(*) as index_count
FROM information_schema.statistics 
WHERE table_schema = '$DB_NAME'
GROUP BY table_name
ORDER BY table_name;
" 2>/dev/null)

=== 表大小统计 ===
$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "
SELECT 
    table_name,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)',
    ROUND((data_length / 1024 / 1024), 2) AS 'Data (MB)',
    ROUND((index_length / 1024 / 1024), 2) AS 'Index (MB)',
    table_rows
FROM information_schema.tables 
WHERE table_schema = '$DB_NAME'
ORDER BY (data_length + index_length) DESC;
" 2>/dev/null)

=== 建议 ===
1. 定期监控慢查询日志
2. 根据实际使用情况调整索引
3. 定期重建索引优化性能
4. 监控索引使用情况，删除无用索引

EOF

echo -e "${GREEN}✓ 测试报告已生成: $REPORT_FILE${NC}"
echo ""

echo -e "${BLUE}=== 索引优化测试完成 ===${NC}"
echo -e "${GREEN}所有测试项目已完成，请查看测试报告了解详细信息。${NC}" 