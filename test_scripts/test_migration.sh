#!/bin/bash

# 数据库迁移测试脚本
echo "=== 数据库迁移测试 ==="

# 检查配置文件
if [ ! -f "config/config.yaml" ]; then
    echo "❌ 配置文件不存在，请先复制 config.example.yaml 为 config.yaml"
    exit 1
fi

echo "📋 当前配置:"
echo "  配置文件: config/config.yaml"
echo "  数据库: $(grep 'dbname:' config/config.yaml | awk '{print $2}')"
echo "  主机: $(grep 'host:' config/config.yaml | head -1 | awk '{print $2}')"
echo "  端口: $(grep 'port:' config/config.yaml | head -1 | awk '{print $2}')"
echo

# 执行迁移
echo "🔄 执行数据库迁移..."
if make db-migrate; then
    echo "✅ 数据库迁移成功！"
else
    echo "❌ 数据库迁移失败！"
    exit 1
fi

echo
echo "📊 验证表结构..."

# 检查数据库连接
if command -v mysql &> /dev/null; then
    DB_HOST=$(grep "host:" config/config.yaml | head -1 | awk '{print $2}')
    DB_PORT=$(grep "port:" config/config.yaml | head -1 | awk '{print $2}')
    DB_NAME=$(grep "dbname:" config/config.yaml | head -1 | awk '{print $2}')
    DB_USER=$(grep "username:" config/config.yaml | head -1 | awk '{print $2}')
    DB_PASS=$(grep "password:" config/config.yaml | head -1 | awk '{print $2}')
    
    # 设置默认值
    DB_HOST=${DB_HOST:-"localhost"}
    DB_PORT=${DB_PORT:-3306}
    DB_NAME=${DB_NAME:-"gin_fataMorgana"}
    DB_USER=${DB_USER:-"root"}
    DB_PASS=${DB_PASS:-""}
    
    echo "🔍 检查表结构..."
    
    # 检查表是否存在
    TABLES=("users" "wallets" "wallet_transactions" "admin_users" "user_login_logs")
    
    for table in "${TABLES[@]}"; do
        if [ -z "$DB_PASS" ]; then
            if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" "$DB_NAME" -e "DESCRIBE $table;" >/dev/null 2>&1; then
                echo "   ✅ $table 表存在"
            else
                echo "   ❌ $table 表不存在"
            fi
        else
            if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "DESCRIBE $table;" >/dev/null 2>&1; then
                echo "   ✅ $table 表存在"
            else
                echo "   ❌ $table 表不存在"
            fi
        fi
    done
    
    echo
    echo "📋 表结构详情:"
    if [ -z "$DB_PASS" ]; then
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" "$DB_NAME" -e "SHOW TABLES;"
    else
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "SHOW TABLES;"
    fi
else
    echo "⚠️  MySQL客户端未安装，跳过表结构验证"
fi

echo
echo "🎉 迁移测试完成！"
echo "💡 提示:"
echo "   - 使用 'make db-seed' 初始化管理员账户"
echo "   - 使用 'make run' 启动应用"
echo "   - 使用 'make health' 检查应用状态" 