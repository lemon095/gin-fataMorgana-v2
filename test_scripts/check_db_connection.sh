#!/bin/bash

# 数据库连接检查脚本
echo "=== 数据库连接检查 ==="

# 检查配置文件
if [ ! -f "config.yaml" ]; then
    echo "❌ 配置文件不存在，请先创建 config.yaml"
    exit 1
fi

# 获取数据库配置
DB_HOST=$(grep "host:" config.yaml | head -1 | awk '{print $2}')
DB_PORT=$(grep "port:" config.yaml | head -1 | awk '{print $2}')
DB_USER=$(grep "username:" config.yaml | head -1 | awk '{print $2}')
DB_PASS=$(grep "password:" config.yaml | head -1 | awk '{print $2}')
DB_NAME=$(grep "dbname:" config.yaml | head -1 | awk '{print $2}')

# 设置默认值
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-3306}
DB_USER=${DB_USER:-"root"}
DB_PASS=${DB_PASS:-""}
DB_NAME=${DB_NAME:-"future"}

echo "📋 数据库配置:"
echo "  主机: $DB_HOST"
echo "  端口: $DB_PORT"
echo "  用户: $DB_USER"
echo "  数据库: $DB_NAME"
echo

# 检查MySQL服务是否运行
echo "🔍 检查MySQL服务状态..."
if command -v mysql &> /dev/null; then
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "SELECT 1;" 2>/dev/null; then
        echo "✅ MySQL服务运行正常"
    else
        echo "❌ MySQL服务连接失败"
        echo "   请检查MySQL服务是否启动，以及连接参数是否正确"
        exit 1
    fi
else
    echo "⚠️  mysql客户端未安装，跳过直接连接测试"
fi
echo

# 检查数据库是否存在
echo "🔍 检查数据库是否存在..."
if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "USE $DB_NAME;" 2>/dev/null; then
    echo "✅ 数据库 $DB_NAME 存在"
else
    echo "❌ 数据库 $DB_NAME 不存在"
    echo "   请先创建数据库: CREATE DATABASE $DB_NAME;"
    exit 1
fi
echo

# 检查表是否存在
echo "🔍 检查核心表是否存在..."
TABLES=("users" "wallets" "wallet_transactions" "admin_users" "user_login_logs" "orders" "amount_config")

for table in "${TABLES[@]}"; do
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -D"$DB_NAME" -e "DESCRIBE $table;" 2>/dev/null | grep -q "Field"; then
        echo "✅ 表 $table 存在"
    else
        echo "❌ 表 $table 不存在"
    fi
done
echo

# 检查连接池配置
echo "🔍 检查连接池配置..."
MAX_IDLE_CONNS=$(grep "max_idle_conns:" config.yaml | head -1 | awk '{print $2}')
MAX_OPEN_CONNS=$(grep "max_open_conns:" config.yaml | head -1 | awk '{print $2}')
CONN_MAX_LIFETIME=$(grep "conn_max_lifetime:" config.yaml | head -1 | awk '{print $2}')
CONN_MAX_IDLE_TIME=$(grep "conn_max_idle_time:" config.yaml | head -1 | awk '{print $2}')

echo "📋 连接池配置:"
echo "  最大空闲连接: ${MAX_IDLE_CONNS:-10}"
echo "  最大连接数: ${MAX_OPEN_CONNS:-100}"
echo "  连接最大生存时间: ${CONN_MAX_LIFETIME:-3600}秒"
echo "  连接空闲超时: ${CONN_MAX_IDLE_TIME:-1800}秒"
echo

# 检查应用服务状态
echo "🔍 检查应用服务状态..."
SERVER_HOST="localhost"
SERVER_PORT=$(grep "port:" config.yaml | head -1 | awk '{print $2}')
SERVER_PORT=${SERVER_PORT:-9001}

if curl -s "http://$SERVER_HOST:$SERVER_PORT/health" > /dev/null; then
    echo "✅ 应用服务运行正常"
    
    # 检查数据库健康状态
    echo "🔍 检查应用数据库连接..."
    DB_HEALTH_RESPONSE=$(curl -s "http://$SERVER_HOST:$SERVER_PORT/api/v1/health/database")
    if echo "$DB_HEALTH_RESPONSE" | grep -q '"status":"healthy"'; then
        echo "✅ 应用数据库连接正常"
    else
        echo "❌ 应用数据库连接异常"
        echo "   响应: $DB_HEALTH_RESPONSE"
    fi
else
    echo "❌ 应用服务未运行"
    echo "   请先启动服务: ./dev.sh start"
fi
echo

echo "🎉 数据库连接检查完成！" 