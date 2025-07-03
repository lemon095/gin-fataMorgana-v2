#!/bin/bash

# 测试数据库迁移功能
echo "=== 测试数据库迁移功能 ==="

# 1. 检查数据库连接
echo "1. 检查数据库连接"
if ! mysql -u root -p -e "SELECT 1;" > /dev/null 2>&1; then
    echo "数据库连接失败，请检查配置"
    exit 1
fi

# 2. 检查新表是否已创建
echo "2. 检查新表是否已创建"

# 检查 lottery_periods 表
if mysql -u root -p -e "USE gin_fataMorgana; DESCRIBE lottery_periods;" > /dev/null 2>&1; then
    echo "✓ lottery_periods 表已存在"
else
    echo "✗ lottery_periods 表不存在"
fi

# 检查 member_level 表
if mysql -u root -p -e "USE gin_fataMorgana; DESCRIBE member_level;" > /dev/null 2>&1; then
    echo "✓ member_level 表已存在"
else
    echo "✗ member_level 表不存在"
fi

# 3. 检查表结构
echo "3. 检查表结构"

echo "lottery_periods 表结构："
mysql -u root -p -e "USE gin_fataMorgana; DESCRIBE lottery_periods;" 2>/dev/null || echo "表不存在"

echo ""
echo "member_level 表结构："
mysql -u root -p -e "USE gin_fataMorgana; DESCRIBE member_level;" 2>/dev/null || echo "表不存在"

# 4. 检查索引
echo "4. 检查索引"

echo "lottery_periods 表索引："
mysql -u root -p -e "USE gin_fataMorgana; SHOW INDEX FROM lottery_periods;" 2>/dev/null || echo "表不存在"

echo ""
echo "member_level 表索引："
mysql -u root -p -e "USE gin_fataMorgana; SHOW INDEX FROM member_level;" 2>/dev/null || echo "表不存在"

# 5. 检查默认数据
echo "5. 检查默认数据"

echo "member_level 表数据："
mysql -u root -p -e "USE gin_fataMorgana; SELECT level, name, min_experience, max_experience, cashback_ratio FROM member_level ORDER BY level;" 2>/dev/null || echo "表不存在或无数据"

# 6. 检查所有表
echo "6. 检查所有表"
echo "数据库中的所有表："
mysql -u root -p -e "USE gin_fataMorgana; SHOW TABLES;" 2>/dev/null || echo "无法获取表列表"

echo "=== 迁移测试完成 ==="
echo ""
echo "如果表不存在，请运行以下命令手动创建："
echo "mysql -u root -p gin_fataMorgana < database/migrations/create_lottery_periods_table.sql"
echo ""
echo "或者重启应用，让GORM自动迁移创建表" 