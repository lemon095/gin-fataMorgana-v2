#!/bin/bash

# 测试参数验证错误处理
echo "=== 测试参数验证错误处理 ==="

BASE_URL="http://localhost:8080/api/v1"

# 测试1: 邮箱格式错误
echo "测试1: 邮箱格式错误"
curl -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "password": "123456"
  }' | jq '.'

echo -e "\n"

# 测试2: 密码为空
echo "测试2: 密码为空"
curl -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": ""
  }' | jq '.'

echo -e "\n"

# 测试3: 邮箱为空
echo "测试3: 邮箱为空"
curl -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "",
    "password": "123456"
  }' | jq '.'

echo -e "\n"

# 测试4: 注册时密码长度不足
echo "测试4: 注册时密码长度不足"
curl -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123",
    "confirm_password": "123",
    "invite_code": "TEST123"
  }' | jq '.'

echo -e "\n"

# 测试5: 注册时两次密码不一致
echo "测试5: 注册时两次密码不一致"
curl -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456",
    "confirm_password": "654321",
    "invite_code": "TEST123"
  }' | jq '.'

echo -e "\n"

# 测试6: 修改密码时新密码长度不足
echo "测试6: 修改密码时新密码长度不足"
curl -X POST "${BASE_URL}/auth/change-password" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "old_password": "123456",
    "new_password": "123"
  }' | jq '.'

echo -e "\n"

# 测试7: 绑定银行卡时参数缺失
echo "测试7: 绑定银行卡时参数缺失"
curl -X POST "${BASE_URL}/auth/bind-bank-card" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "bank_name": "中国银行",
    "card_holder": "张三"
  }' | jq '.'

echo -e "\n"

# 测试8: 获取用户信息时参数错误
echo "测试8: 获取用户信息时参数错误"
curl -X POST "${BASE_URL}/auth/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "invalid_field": "invalid_value"
  }' | jq '.'

echo -e "\n"

echo "=== 测试完成 ===" 