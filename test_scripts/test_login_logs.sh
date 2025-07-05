#!/bin/bash

# 测试登录接口日志输出
echo "=== 测试登录接口日志输出 ==="

# 设置基础URL
BASE_URL="http://localhost:9001/api/v1"

echo "1. 测试正常登录"
echo "请求URL: $BASE_URL/auth/login"

# 正常登录请求
response=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test@example.com",
    "password": "123456"
  }')

echo "响应:"
echo "$response" | jq '.'

echo ""
echo "2. 测试密码错误"
echo "请求URL: $BASE_URL/auth/login"

# 密码错误的登录请求
response2=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test@example.com",
    "password": "wrongpassword"
  }')

echo "响应:"
echo "$response2" | jq '.'

echo ""
echo "3. 测试账号不存在"
echo "请求URL: $BASE_URL/auth/login"

# 账号不存在的登录请求
response3=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "nonexistent@example.com",
    "password": "123456"
  }')

echo "响应:"
echo "$response3" | jq '.'

echo ""
echo "4. 测试参数缺失"
echo "请求URL: $BASE_URL/auth/login"

# 参数缺失的登录请求
response4=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test@example.com"
  }')

echo "响应:"
echo "$response4" | jq '.'

echo ""
echo "5. 测试密码长度不足"
echo "请求URL: $BASE_URL/auth/login"

# 密码长度不足的登录请求
response5=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test@example.com",
    "password": "123"
  }')

echo "响应:"
echo "$response5" | jq '.'

echo ""
echo "=== 测试完成 ==="
echo "请查看服务器日志以查看详细的登录过程日志输出" 