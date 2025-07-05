#!/bin/bash

# 检查配置状态脚本

echo "=== 检查配置状态 ==="

# 检查配置文件
echo "📁 配置文件检查:"
if [ -f "config.yaml" ]; then
    echo "✅ config.yaml 存在"
    echo "假数据配置:"
    grep -A 10 "fake_data:" config.yaml
else
    echo "❌ config.yaml 不存在"
fi

echo ""

# 检查环境变量
echo "🌍 环境变量检查:"
echo "FAKE_DATA_ENABLED: ${FAKE_DATA_ENABLED:-未设置}"
echo "FAKE_DATA_CRON_EXPRESSION: ${FAKE_DATA_CRON_EXPRESSION:-未设置}"
echo "FAKE_DATA_MIN_ORDERS: ${FAKE_DATA_MIN_ORDERS:-未设置}"
echo "FAKE_DATA_MAX_ORDERS: ${FAKE_DATA_MAX_ORDERS:-未设置}"

echo ""

# 检查服务是否运行
echo "🔍 服务状态检查:"
if pgrep -f "gin-fataMorgana" > /dev/null; then
    echo "✅ 服务正在运行"
    ps aux | grep "gin-fataMorgana" | grep -v grep
else
    echo "❌ 服务未运行"
fi

echo ""

# 检查数据库连接
echo "🗄️  数据库连接检查:"
if command -v mysql &> /dev/null; then
    if mysql -h 127.0.0.1 -u root -proot -D future -e "SELECT 1;" &> /dev/null; then
        echo "✅ 数据库连接正常"
        
        # 检查系统订单数量
        SYSTEM_ORDERS=$(mysql -h 127.0.0.1 -u root -proot -D future -e "SELECT COUNT(*) FROM orders WHERE is_system_order = 1;" -s -N 2>/dev/null)
        echo "系统订单数量: ${SYSTEM_ORDERS:-0}"
    else
        echo "❌ 数据库连接失败"
    fi
else
    echo "⚠️  mysql 命令未找到"
fi

echo ""

echo "=== 配置检查完成 ===" 