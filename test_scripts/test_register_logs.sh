#!/bin/bash

# 测试注册接口日志输出
echo "=== 测试注册接口日志输出 ==="

# 设置基础URL
BASE_URL="http://localhost:9001/api/v1"

echo "1. 测试正常注册"
echo "请求URL: $BASE_URL/auth/register"

# 正常注册请求
response=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "ABC123"
  }')

echo "响应:"
echo "$response" | jq '.'

echo ""
echo "2. 测试密码不一致"
echo "请求URL: $BASE_URL/auth/register"

# 密码不一致的注册请求
response2=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test2@example.com",
    "password": "123456",
    "confirm_password": "654321",
    "invite_code": "ABC123"
  }')

echo "响应:"
echo "$response2" | jq '.'

echo ""
echo "3. 测试重复注册"
echo "请求URL: $BASE_URL/auth/register"

# 重复注册请求
response3=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "ABC123"
  }')

echo "响应:"
echo "$response3" | jq '.'

echo ""
echo "4. 测试参数缺失"
echo "请求URL: $BASE_URL/auth/register"

# 参数缺失的注册请求
response4=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test3@example.com",
    "password": "123456"
  }')

echo "响应:"
echo "$response4" | jq '.'

echo ""
echo "5. 测试密码长度不足"
echo "请求URL: $BASE_URL/auth/register"

# 密码长度不足的注册请求
response5=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test4@example.com",
    "password": "123",
    "confirm_password": "123",
    "invite_code": "ABC123"
  }')

echo "响应:"
echo "$response5" | jq '.'

echo ""
echo "=== 测试完成 ==="
echo "请查看服务器日志以查看详细的注册过程日志输出" 