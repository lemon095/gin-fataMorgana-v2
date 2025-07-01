#!/bin/bash

# 用户注册脚本
echo "=== 用户注册脚本 ==="

# 检查配置文件
if [ ! -f "config.yaml" ]; then
    echo "❌ 配置文件不存在，请先创建 config.yaml"
    exit 1
fi

# 获取服务器地址
SERVER_HOST="localhost"
SERVER_PORT=$(grep "port:" config.yaml | head -1 | awk '{print $2}')

# 设置默认值
SERVER_PORT=${SERVER_PORT:-9001}

echo "📋 服务器配置:"
echo "  主机: $SERVER_HOST"
echo "  端口: $SERVER_PORT"
echo

# 邀请码和密码
INVITE_CODE="7TRABJ"
PASSWORD="123456"

# 用户信息
USER1_EMAIL="user1@example.com"
USER1_USERNAME="user1"

USER2_EMAIL="user2@example.com"
USER2_USERNAME="user2"

echo "👥 准备注册用户:"
echo "  邀请码: $INVITE_CODE"
echo "  密码: $PASSWORD"
echo "  用户1: $USER1_USERNAME ($USER1_EMAIL)"
echo "  用户2: $USER2_USERNAME ($USER2_EMAIL)"
echo

# 检查服务是否运行
echo "🔍 检查服务状态..."
if ! curl -s "http://$SERVER_HOST:$SERVER_PORT/health" > /dev/null; then
    echo "❌ 服务未运行，请先启动服务:"
    echo "   ./dev.sh start"
    echo "   或者"
    echo "   ./prod.sh start"
    exit 1
fi
echo "✅ 服务运行正常"
echo

# 注册第一个用户
echo "📝 注册第一个用户: $USER1_USERNAME"
USER1_RESPONSE=$(curl -s -X POST "http://$SERVER_HOST:$SERVER_PORT/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$USER1_EMAIL\",
    \"password\": \"$PASSWORD\",
    \"confirm_password\": \"$PASSWORD\",
    \"invite_code\": \"$INVITE_CODE\"
  }")

echo "响应: $USER1_RESPONSE"

# 检查第一个用户注册结果
if echo "$USER1_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 用户1注册成功"
    USER1_UID=$(echo "$USER1_RESPONSE" | grep -o '"uid":"[^"]*"' | cut -d'"' -f4)
    USER1_INVITE_CODE=$(echo "$USER1_RESPONSE" | grep -o '"my_invite_code":"[^"]*"' | cut -d'"' -f4)
    echo "   UID: $USER1_UID"
    echo "   邀请码: $USER1_INVITE_CODE"
else
    echo "❌ 用户1注册失败"
    echo "   错误信息: $USER1_RESPONSE"
fi
echo

# 注册第二个用户
echo "📝 注册第二个用户: $USER2_USERNAME"
USER2_RESPONSE=$(curl -s -X POST "http://$SERVER_HOST:$SERVER_PORT/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$USER2_EMAIL\",
    \"password\": \"$PASSWORD\",
    \"confirm_password\": \"$PASSWORD\",
    \"invite_code\": \"$INVITE_CODE\"
  }")

echo "响应: $USER2_RESPONSE"

# 检查第二个用户注册结果
if echo "$USER2_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 用户2注册成功"
    USER2_UID=$(echo "$USER2_RESPONSE" | grep -o '"uid":"[^"]*"' | cut -d'"' -f4)
    USER2_INVITE_CODE=$(echo "$USER2_RESPONSE" | grep -o '"my_invite_code":"[^"]*"' | cut -d'"' -f4)
    echo "   UID: $USER2_UID"
    echo "   邀请码: $USER2_INVITE_CODE"
else
    echo "❌ 用户2注册失败"
    echo "   错误信息: $USER2_RESPONSE"
fi
echo

# 测试登录
echo "🔐 测试用户登录..."

# 测试第一个用户登录
echo "📝 测试用户1登录: $USER1_USERNAME"
LOGIN1_RESPONSE=$(curl -s -X POST "http://$SERVER_HOST:$SERVER_PORT/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$USER1_EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

if echo "$LOGIN1_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 用户1登录成功"
    USER1_TOKEN=$(echo "$LOGIN1_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    echo "   访问令牌: ${USER1_TOKEN:0:20}..."
else
    echo "❌ 用户1登录失败"
    echo "   错误信息: $LOGIN1_RESPONSE"
fi
echo

# 测试第二个用户登录
echo "📝 测试用户2登录: $USER2_USERNAME"
LOGIN2_RESPONSE=$(curl -s -X POST "http://$SERVER_HOST:$SERVER_PORT/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$USER2_EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

if echo "$LOGIN2_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 用户2登录成功"
    USER2_TOKEN=$(echo "$LOGIN2_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    echo "   访问令牌: ${USER2_TOKEN:0:20}..."
else
    echo "❌ 用户2登录失败"
    echo "   错误信息: $LOGIN2_RESPONSE"
fi
echo

# 查询钱包信息
echo "💰 查询用户钱包信息..."

if [ ! -z "$USER1_TOKEN" ]; then
    echo "📝 查询用户1钱包: $USER1_USERNAME"
    WALLET1_RESPONSE=$(curl -s -X GET "http://$SERVER_HOST:$SERVER_PORT/wallet/info" \
      -H "Authorization: Bearer $USER1_TOKEN")
    
    if echo "$WALLET1_RESPONSE" | grep -q '"code":200'; then
        echo "✅ 用户1钱包查询成功"
        BALANCE1=$(echo "$WALLET1_RESPONSE" | grep -o '"balance":"[^"]*"' | cut -d'"' -f4)
        echo "   余额: $BALANCE1"
    else
        echo "❌ 用户1钱包查询失败"
        echo "   错误信息: $WALLET1_RESPONSE"
    fi
    echo
fi

if [ ! -z "$USER2_TOKEN" ]; then
    echo "📝 查询用户2钱包: $USER2_USERNAME"
    WALLET2_RESPONSE=$(curl -s -X GET "http://$SERVER_HOST:$SERVER_PORT/wallet/info" \
      -H "Authorization: Bearer $USER2_TOKEN")
    
    if echo "$WALLET2_RESPONSE" | grep -q '"code":200'; then
        echo "✅ 用户2钱包查询成功"
        BALANCE2=$(echo "$WALLET2_RESPONSE" | grep -o '"balance":"[^"]*"' | cut -d'"' -f4)
        echo "   余额: $BALANCE2"
    else
        echo "❌ 用户2钱包查询失败"
        echo "   错误信息: $WALLET2_RESPONSE"
    fi
    echo
fi

echo "🎉 用户注册测试完成！"
echo
echo "📋 注册结果总结:"
echo "  邀请码: $INVITE_CODE"
echo "  用户1: $USER1_USERNAME ($USER1_EMAIL) - 密码: $PASSWORD"
echo "  用户2: $USER2_USERNAME ($USER2_EMAIL) - 密码: $PASSWORD"
echo
echo "💡 提示:"
echo "   - 使用 './dev.sh start' 启动本地开发服务"
echo "   - 使用 './prod.sh start' 启动生产服务"
echo "   - 使用 'curl' 或 Postman 测试API"
echo "   - 查看日志: tail -f logs/app.log" 