#!/bin/bash

# 超级管理员初始化脚本
# 通过API接口创建超级管理员用户

echo "=== 初始化超级管理员用户 ==="
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

# 生成超级管理员信息
SUPER_ADMIN_USERNAME="super_admin_$(date +%Y%m%d)"
SUPER_ADMIN_PASSWORD="SuperAdmin123!"
SUPER_ADMIN_EMAIL="super_admin_$(date +%Y%m%d)@example.com"

echo "👑 超级管理员信息:"
echo "  用户名: $SUPER_ADMIN_USERNAME"
echo "  邮箱: $SUPER_ADMIN_EMAIL"
echo "  密码: $SUPER_ADMIN_PASSWORD"
echo "  角色: 超级管理员 (role=1)"
echo

# 使用默认邀请码
DEFAULT_INVITE_CODE="SUPER123"

echo "📝 使用邀请码: $DEFAULT_INVITE_CODE"

# 通过注册接口创建超级管理员
echo "📝 创建超级管理员用户..."

SUPER_ADMIN_RESPONSE=$(curl -s -X POST "http://$SERVER_HOST:$SERVER_PORT/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$SUPER_ADMIN_EMAIL\",
    \"password\": \"$SUPER_ADMIN_PASSWORD\",
    \"confirm_password\": \"$SUPER_ADMIN_PASSWORD\",
    \"invite_code\": \"$DEFAULT_INVITE_CODE\"
  }")

echo "注册响应: $SUPER_ADMIN_RESPONSE"

# 检查注册结果
if echo "$SUPER_ADMIN_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 超级管理员用户创建成功！"
    SUPER_ADMIN_UID=$(echo "$SUPER_ADMIN_RESPONSE" | grep -o '"uid":"[^"]*"' | cut -d'"' -f4)
    SUPER_ADMIN_INVITE_CODE=$(echo "$SUPER_ADMIN_RESPONSE" | grep -o '"my_invite_code":"[^"]*"' | cut -d'"' -f4)
    echo "   UID: $SUPER_ADMIN_UID"
    echo "   邀请码: $SUPER_ADMIN_INVITE_CODE"
    
    # 测试登录
    echo "🔐 测试超级管理员登录..."
    LOGIN_RESPONSE=$(curl -s -X POST "http://$SERVER_HOST:$SERVER_PORT/auth/login" \
      -H "Content-Type: application/json" \
      -d "{
        \"email\": \"$SUPER_ADMIN_EMAIL\",
        \"password\": \"$SUPER_ADMIN_PASSWORD\"
      }")
    
    if echo "$LOGIN_RESPONSE" | grep -q '"code":200'; then
        echo "✅ 超级管理员登录成功"
        SUPER_ADMIN_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
        echo "   访问令牌: ${SUPER_ADMIN_TOKEN:0:20}..."
        
        # 获取用户信息验证角色
        echo "🔍 验证用户角色..."
        USER_INFO_RESPONSE=$(curl -s -X GET "http://$SERVER_HOST:$SERVER_PORT/auth/profile" \
          -H "Authorization: Bearer $SUPER_ADMIN_TOKEN")
        
        echo "用户信息: $USER_INFO_RESPONSE"
        
    else
        echo "❌ 超级管理员登录失败"
        echo "   错误信息: $LOGIN_RESPONSE"
    fi
    
else
    echo "❌ 超级管理员用户创建失败"
    echo "   错误信息: $SUPER_ADMIN_RESPONSE"
    echo
    echo "💡 可能的解决方案:"
    echo "   1. 检查邀请码是否有效"
    echo "   2. 检查邮箱是否已被使用"
    echo "   3. 检查服务是否正常运行"
    echo "   4. 手动在数据库中创建超级管理员用户"
    echo
    echo "🔧 手动创建超级管理员的SQL:"
    echo "INSERT INTO admin_users (admin_id, username, password, remark, status, role, my_invite_code, created_at, updated_at) VALUES (1, 'super_admin', 'hashed_password', '超级管理员', 1, 1, 'SUPER123', NOW(3), NOW(3));"
fi

echo
echo "📋 超级管理员信息总结:"
echo "  用户名: $SUPER_ADMIN_USERNAME"
echo "  邮箱: $SUPER_ADMIN_EMAIL"
echo "  密码: $SUPER_ADMIN_PASSWORD"
echo "  角色: 超级管理员 (role=1)"
echo "  邀请码: $DEFAULT_INVITE_CODE"
echo
echo "💡 提示:"
echo "   - 超级管理员拥有最高权限"
echo "   - 如果自动创建失败，请手动在数据库中创建"
echo "   - 确保邀请码有效"
echo "   - 可以修改脚本中的用户信息"
echo
echo "⚠️  安全提醒:"
echo "   - 请及时修改默认密码"
echo "   - 妥善保管超级管理员凭据"
echo "   - 定期更换密码"
echo
echo "=== 超级管理员初始化完成 ===" 