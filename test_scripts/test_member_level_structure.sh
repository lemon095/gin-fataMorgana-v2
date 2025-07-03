#!/bin/bash

# 测试member_level表结构
echo "=== 测试member_level表结构 ==="

# 数据库连接信息（请根据实际情况修改）
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="gin_fata_morgana"
DB_USER="root"
DB_PASS="123456"

echo "1. 查看当前表结构"
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASS $DB_NAME -e "
DESCRIBE member_level;
"

echo ""
echo "2. 查看完整建表语句"
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASS $DB_NAME -e "
SHOW CREATE TABLE member_level;
"

echo ""
echo "3. 查看表索引"
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASS $DB_NAME -e "
SHOW INDEX FROM member_level;
"

echo ""
echo "4. 查看表数据（如果有的话）"
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASS $DB_NAME -e "
SELECT * FROM member_level ORDER BY level;
"

echo ""
echo "=== 表结构验证完成 ==="
echo ""
echo "预期字段列表："
echo "- id (bigint unsigned, 主键)"
echo "- level (bigint, 等级数值, 唯一索引)"
echo "- name (varchar(20), 等级名称)"
echo "- logo (varchar(255), 等级logo)"
echo "- remark (varchar(255), 备注)"
echo "- cashback_ratio (decimal(5,2), 返现比例)"
echo "- single_amount (int, 单数字额, 默认值1)"
echo "- created_at (datetime(3), 创建时间)"
echo "- updated_at (datetime(3), 更新时间)"
echo "- deleted_at (datetime(3), 软删除时间)"
echo ""
echo "预期索引："
echo "- PRIMARY KEY (id)"
echo "- UNIQUE KEY uniq_level (level)"
echo "- KEY idx_member_level_deleted_at (deleted_at)" 