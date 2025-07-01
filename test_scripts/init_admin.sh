#!/bin/bash

# 初始化管理员用户脚本
# 通过API接口创建管理员用户

echo "=== 初始化管理员用户 ==="
echo

# 检查配置文件
CONFIG_FILE="config.yaml"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ 错误: 配置文件 $CONFIG_FILE 不存在"
    exit 1
fi

# 获取服务器配置
SERVER_HOST="localhost"
SERVER_PORT=$(grep "port:" "$CONFIG_FILE" | head -1 | awk '{print $2}')
SERVER_PORT=${SERVER_PORT:-9001}

echo "📋 服务器配置:"
echo "  主机: $SERVER_HOST"
echo "  端口: $SERVER_PORT"
echo

# 检查服务是否运行
echo "🔍 检查服务状态..."
if ! curl -s "http://$SERVER_HOST:$SERVER_PORT/health" > /dev/null 2>&1; then
    echo "❌ 服务未运行，请先启动服务:"
    echo "   ./dev.sh start (本地开发)"
    echo "   或者"
    echo "   ./prod.sh start (生产环境)"
    exit 1
fi
echo "✅ 服务运行正常"
echo

# 生成管理员信息
ADMIN_USERNAME="admin_$(date +%Y%m%d)"
ADMIN_PASSWORD="admin123"
ADMIN_EMAIL="admin_$(date +%Y%m%d)@example.com"

echo "👤 管理员信息:"
echo "  用户名: $ADMIN_USERNAME"
echo "  邮箱: $ADMIN_EMAIL"
echo "  密码: $ADMIN_PASSWORD"
echo

# 方法1: 尝试通过注册接口创建管理员（如果有管理员注册接口）
echo "📝 尝试创建管理员用户..."

# 首先尝试获取一个有效的邀请码（如果有的话）
echo "🔍 检查是否有可用的邀请码..."

# 这里可以添加获取邀请码的逻辑，或者使用默认邀请码
DEFAULT_INVITE_CODE="ADMIN123"

echo "📝 使用邀请码: $DEFAULT_INVITE_CODE"

# 通过注册接口创建管理员
ADMIN_RESPONSE=$(curl -s -X POST "http://$SERVER_HOST:$SERVER_PORT/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$ADMIN_EMAIL\",
    \"password\": \"$ADMIN_PASSWORD\",
    \"confirm_password\": \"$ADMIN_PASSWORD\",
    \"invite_code\": \"$DEFAULT_INVITE_CODE\"
  }")

echo "注册响应: $ADMIN_RESPONSE"

# 检查注册结果
if echo "$ADMIN_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 管理员用户创建成功！"
    ADMIN_UID=$(echo "$ADMIN_RESPONSE" | grep -o '"uid":"[^"]*"' | cut -d'"' -f4)
    ADMIN_INVITE_CODE=$(echo "$ADMIN_RESPONSE" | grep -o '"my_invite_code":"[^"]*"' | cut -d'"' -f4)
    echo "   UID: $ADMIN_UID"
    echo "   邀请码: $ADMIN_INVITE_CODE"
    
    # 测试登录
    echo "🔐 测试管理员登录..."
    LOGIN_RESPONSE=$(curl -s -X POST "http://$SERVER_HOST:$SERVER_PORT/auth/login" \
      -H "Content-Type: application/json" \
      -d "{
        \"email\": \"$ADMIN_EMAIL\",
        \"password\": \"$ADMIN_PASSWORD\"
      }")
    
    if echo "$LOGIN_RESPONSE" | grep -q '"code":200'; then
        echo "✅ 管理员登录成功"
        ADMIN_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
        echo "   访问令牌: ${ADMIN_TOKEN:0:20}..."
    else
        echo "❌ 管理员登录失败"
        echo "   错误信息: $LOGIN_RESPONSE"
    fi
    
else
    echo "❌ 管理员用户创建失败"
    echo "   错误信息: $ADMIN_RESPONSE"
    echo
    echo "💡 可能的解决方案:"
    echo "   1. 检查邀请码是否有效"
    echo "   2. 检查邮箱是否已被使用"
    echo "   3. 检查服务是否正常运行"
    echo "   4. 手动在数据库中创建管理员用户"
fi

echo
echo "📋 管理员信息总结:"
echo "  用户名: $ADMIN_USERNAME"
echo "  邮箱: $ADMIN_EMAIL"
echo "  密码: $ADMIN_PASSWORD"
echo "  邀请码: $DEFAULT_INVITE_CODE"
echo
echo "💡 提示:"
echo "   - 如果自动创建失败，请手动注册管理员用户"
echo "   - 确保邀请码有效"
echo "   - 可以修改脚本中的邀请码和用户信息"
echo
echo "=== 初始化完成 ===" 