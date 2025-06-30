#!/bin/bash

# 登录记录功能测试脚本
BASE_URL="http://localhost:9001"

echo "=== 登录记录功能测试 ==="
echo

# 1. 注册测试用户
echo "1. 注册测试用户..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test_login_log@example.com",
    "password": "password123",
    "confirm_password": "password123",
    "invite_code": ""
  }')

echo "注册响应: $REGISTER_RESPONSE"
echo

# 提取用户UID
USER_UID=$(echo $REGISTER_RESPONSE | grep -o '"uid":"[^"]*"' | cut -d'"' -f4)
echo "用户UID: $USER_UID"
echo

# 2. 测试登录（记录登录信息）
echo "2. 测试登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36" \
  -H "X-Forwarded-For: 192.168.1.100" \
  -d '{
    "email": "test_login_log@example.com",
    "password": "password123"
  }')

echo "登录响应: $LOGIN_RESPONSE"
echo

# 3. 测试失败登录（记录失败信息）
echo "3. 测试失败登录..."
FAILED_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -H "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36" \
  -H "X-Forwarded-For: 192.168.1.101" \
  -d '{
    "email": "test_login_log@example.com",
    "password": "wrongpassword"
  }')

echo "失败登录响应: $FAILED_LOGIN_RESPONSE"
echo

# 4. 再次成功登录
echo "4. 再次成功登录..."
LOGIN_RESPONSE2=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -H "User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15" \
  -H "X-Forwarded-For: 192.168.1.102" \
  -d '{
    "email": "test_login_log@example.com",
    "password": "password123"
  }')

echo "第二次登录响应: $LOGIN_RESPONSE2"
echo

# 5. 获取用户登录历史
echo "5. 获取用户登录历史..."
LOGIN_HISTORY_RESPONSE=$(curl -s -X GET "$BASE_URL/api/users/$USER_UID/login-history?page=1&size=10")

echo "登录历史响应: $LOGIN_HISTORY_RESPONSE"
echo

# 6. 获取用户最后登录信息
echo "6. 获取用户最后登录信息..."
LAST_LOGIN_RESPONSE=$(curl -s -X GET "$BASE_URL/api/users/$USER_UID/last-login")

echo "最后登录信息响应: $LAST_LOGIN_RESPONSE"
echo

# 7. 获取登录统计信息
echo "7. 获取登录统计信息..."
LOGIN_STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/users/$USER_UID/login-stats")

echo "登录统计信息响应: $LOGIN_STATS_RESPONSE"
echo

# 8. 按IP地址查询登录记录
echo "8. 按IP地址查询登录记录..."
LOGIN_BY_IP_RESPONSE=$(curl -s -X GET "$BASE_URL/api/users/$USER_UID/login-by-ip?ip=192.168.1.100")

echo "按IP查询响应: $LOGIN_BY_IP_RESPONSE"
echo

# 9. 按时间范围查询登录记录
echo "9. 按时间范围查询登录记录..."
START_TIME=$(date -d "1 hour ago" "+%Y-%m-%d %H:%M:%S")
END_TIME=$(date "+%Y-%m-%d %H:%M:%S")
LOGIN_BY_TIME_RESPONSE=$(curl -s -X GET "$BASE_URL/api/users/$USER_UID/login-by-time?start_time=$START_TIME&end_time=$END_TIME")

echo "按时间范围查询响应: $LOGIN_BY_TIME_RESPONSE"
echo

echo "=== 登录记录功能测试完成 ===" 