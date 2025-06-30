#!/bin/bash

# 数据库迁移测试脚本
echo "=== 数据库迁移测试 ==="

# 检查配置文件
if [ ! -f "config.yaml" ]; then
    echo "❌ 配置文件不存在，请先创建 config.yaml"
    exit 1
fi

echo "📋 当前配置:"
echo "  配置文件: config.yaml"
echo "  数据库: $(grep 'dbname:' config.yaml | awk '{print $2}')"
echo "  主机: $(grep 'host:' config.yaml | head -1 | awk '{print $2}')"
echo "  端口: $(grep 'port:' config.yaml | head -1 | awk '{print $2}')"
echo

# 检查MySQL客户端
if ! command -v mysql &> /dev/null; then
    echo "❌ 错误: 未找到mysql客户端，请先安装MySQL客户端"
    exit 1
fi

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到Go环境，请先安装Go"
    exit 1
fi

echo "✅ 环境检查通过"
echo

# 提取数据库配置
DB_HOST=$(grep "host:" config.yaml | head -1 | awk '{print $2}')
DB_PORT=$(grep "port:" config.yaml | head -1 | awk '{print $2}')
DB_NAME=$(grep "dbname:" config.yaml | head -1 | awk '{print $2}')
DB_USER=$(grep "username:" config.yaml | head -1 | awk '{print $2}')
DB_PASS=$(grep "password:" config.yaml | head -1 | awk '{print $2}')

# 设置默认值
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-3306}
DB_NAME=${DB_NAME:-"future"}
DB_USER=${DB_USER:-"root"}
DB_PASS=${DB_PASS:-""}

echo "🔍 检查数据库连接..."
if [ -z "$DB_PASS" ]; then
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -e "SELECT 1;" > /dev/null 2>&1; then
        echo "✅ 数据库连接成功"
    else
        echo "❌ 数据库连接失败"
        exit 1
    fi
else
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "SELECT 1;" > /dev/null 2>&1; then
        echo "✅ 数据库连接成功"
    else
        echo "❌ 数据库连接失败"
        exit 1
    fi
fi

echo "🔍 检查数据库是否存在..."
if [ -z "$DB_PASS" ]; then
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -e "USE $DB_NAME;" > /dev/null 2>&1; then
        echo "✅ 数据库 $DB_NAME 存在"
    else
        echo "⚠️  数据库 $DB_NAME 不存在，正在创建..."
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
        echo "✅ 数据库创建成功"
    fi
else
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "USE $DB_NAME;" > /dev/null 2>&1; then
        echo "✅ 数据库 $DB_NAME 存在"
    else
        echo "⚠️  数据库 $DB_NAME 不存在，正在创建..."
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
        echo "✅ 数据库创建成功"
    fi
fi

echo "🔍 检查数据表..."
TABLES=("users" "wallets" "wallet_transactions" "admin_users" "user_login_logs")
EXISTING_TABLES=()

for table in "${TABLES[@]}"; do
    if [ -z "$DB_PASS" ]; then
        if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" "$DB_NAME" -e "SHOW TABLES LIKE '$table';" 2>/dev/null | grep -q "$table"; then
            EXISTING_TABLES+=("$table")
        fi
    else
        if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "SHOW TABLES LIKE '$table';" 2>/dev/null | grep -q "$table"; then
            EXISTING_TABLES+=("$table")
        fi
    fi
done

echo "📊 现有数据表: ${EXISTING_TABLES[*]}"
echo "📊 需要的数据表: ${TABLES[*]}"

if [ ${#EXISTING_TABLES[@]} -eq ${#TABLES[@]} ]; then
    echo "✅ 所有数据表已存在"
else
    echo "⚠️  部分数据表缺失，需要运行迁移"
fi

echo
echo "🚀 运行数据库迁移..."
if [ -f "cmd/migrate/main.go" ]; then
    go run cmd/migrate/main.go
    if [ $? -eq 0 ]; then
        echo "✅ 迁移完成"
    else
        echo "❌ 迁移失败"
        exit 1
    fi
else
    echo "❌ 迁移文件不存在: cmd/migrate/main.go"
    exit 1
fi

echo
echo "🔍 验证迁移结果..."
for table in "${TABLES[@]}"; do
    if [ -z "$DB_PASS" ]; then
        if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" "$DB_NAME" -e "SHOW TABLES LIKE '$table';" 2>/dev/null | grep -q "$table"; then
            echo "✅ 表 $table 存在"
        else
            echo "❌ 表 $table 不存在"
        fi
    else
        if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "SHOW TABLES LIKE '$table';" 2>/dev/null | grep -q "$table"; then
            echo "✅ 表 $table 存在"
        else
            echo "❌ 表 $table 不存在"
        fi
    fi
done

echo
echo "🎉 数据库迁移测试完成！" 