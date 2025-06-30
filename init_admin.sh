#!/bin/bash

# 初始化管理员用户脚本
# 用于创建管理员用户并获取邀请码

echo "=== 初始化管理员用户 ==="
echo

# 检查是否安装了mysql客户端
if ! command -v mysql &> /dev/null; then
    echo "❌ 错误: 未找到mysql客户端，请先安装MySQL客户端"
    exit 1
fi

# 从配置文件读取数据库信息
CONFIG_FILE="config/config.yaml"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ 错误: 配置文件 $CONFIG_FILE 不存在"
    exit 1
fi

# 提取数据库配置（简单的方式）
DB_HOST=$(grep "host:" "$CONFIG_FILE" | head -1 | awk '{print $2}')
DB_PORT=$(grep "port:" "$CONFIG_FILE" | head -1 | awk '{print $2}')
DB_NAME=$(grep "dbname:" "$CONFIG_FILE" | head -1 | awk '{print $2}')
DB_USER=$(grep "username:" "$CONFIG_FILE" | head -1 | awk '{print $2}')
DB_PASS=$(grep "password:" "$CONFIG_FILE" | head -1 | awk '{print $2}')

# 设置默认值
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-3306}
DB_NAME=${DB_NAME:-"gin_fatamorgana"}
DB_USER=${DB_USER:-"root"}
DB_PASS=${DB_PASS:-""}

echo "📊 数据库配置:"
echo "  主机: $DB_HOST"
echo "  端口: $DB_PORT"
echo "  数据库: $DB_NAME"
echo "  用户: $DB_USER"
echo

# 生成管理员信息
ADMIN_ID=$(date +%s | tail -c 8)
ADMIN_USERNAME="admin_$(date +%Y%m%d)"
ADMIN_PASSWORD="admin123"
ADMIN_INVITE_CODE=$(cat /dev/urandom | tr -dc 'A-Z0-9' | fold -w 6 | head -n 1)

echo "👤 管理员信息:"
echo "  用户名: $ADMIN_USERNAME"
echo "  密码: $ADMIN_PASSWORD"
echo "  邀请码: $ADMIN_INVITE_CODE"
echo

# 构建SQL语句
SQL="INSERT INTO admin_users (admin_id, username, password, remark, status, role, my_invite_code, created_at, updated_at) VALUES ($ADMIN_ID, '$ADMIN_USERNAME', '$(echo -n "$ADMIN_PASSWORD" | openssl dgst -sha256 | cut -d' ' -f2)', '系统管理员', 1, 4, '$ADMIN_INVITE_CODE', NOW(3), NOW(3));"

echo "📝 执行SQL语句:"
echo "$SQL"
echo

# 执行SQL
if [ -z "$DB_PASS" ]; then
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" "$DB_NAME" -e "$SQL"
else
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "$SQL"
fi

if [ $? -eq 0 ]; then
    echo "✅ 管理员用户创建成功！"
    echo
    echo "📋 请记录以下信息:"
    echo "  用户名: $ADMIN_USERNAME"
    echo "  密码: $ADMIN_PASSWORD"
    echo "  邀请码: $ADMIN_INVITE_CODE"
    echo
    echo "💡 用户注册时需要使用此邀请码"
    echo "⚠️  注意：密码已使用SHA256加密存储，实际使用时需要使用bcrypt"
else
    echo "❌ 管理员用户创建失败！"
    echo "请检查数据库连接和表结构"
fi

echo
echo "=== 初始化完成 ===" 