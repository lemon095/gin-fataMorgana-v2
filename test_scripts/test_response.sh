#!/bin/bash

# 测试统一响应格式的脚本

echo "🚀 开始测试统一响应格式..."
echo ""

# 测试基础接口
echo "📋 测试基础接口:"
echo "1. 首页接口:"
curl -s http://localhost:9001/ | jq .
echo ""

echo "2. 健康检查接口:"
curl -s http://localhost:9001/health | jq .
echo ""

# 测试注册接口
echo "📝 测试用户注册接口:"
echo "1. 正常注册:"
curl -s -X POST http://localhost:9001/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "INVITE123"
  }' | jq .
echo ""

echo "2. 重复邮箱注册:"
curl -s -X POST http://localhost:9001/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "INVITE123"
  }' | jq .
echo ""

echo "3. 参数错误注册:"
curl -s -X POST http://localhost:9001/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "password": "123",
    "confirm_password": "456",
    "invite_code": "INVITE123"
  }' | jq .
echo ""

# 测试登录接口
echo "🔐 测试用户登录接口:"
echo "1. 正常登录:"
curl -s -X POST http://localhost:9001/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456"
  }' | jq .
echo ""

echo "2. 错误密码登录:"
curl -s -X POST http://localhost:9001/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "wrongpassword"
  }' | jq .
echo ""

echo "3. 不存在的用户登录:"
curl -s -X POST http://localhost:9001/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "nonexistent@example.com",
    "password": "123456"
  }' | jq .
echo ""

# 测试会话接口
echo "💬 测试会话接口:"
echo "1. 检查登录状态:"
curl -s http://localhost:9001/session/status | jq .
echo ""

echo "2. 获取用户信息(未登录):"
curl -s http://localhost:9001/session/user | jq .
echo ""

# 测试公共接口
echo "🌐 测试公共接口:"
echo "1. 公共信息(未登录):"
curl -s http://localhost:9001/public/info | jq .
echo ""

echo "✅ 测试完成！"
echo ""
echo "📊 响应格式总结:"
echo "- 所有成功响应: code=0, message='操作成功' 或自定义消息"
echo "- 所有错误响应: code=错误码, message=错误消息, data=null"
echo "- HTTP状态码根据业务错误码自动映射"
echo ""
echo "🔍 错误码分类:"
echo "- 1000-1999: 客户端错误 (HTTP 400)"
echo "- 2000-2099: 认证错误 (HTTP 401)"
echo "- 3000-3999: 业务错误 (HTTP 422)"
echo "- 5000-5999: 服务器错误 (HTTP 500)" 