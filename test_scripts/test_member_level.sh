#!/bin/bash

# 测试用户等级配置功能
# 需要先确保数据库中有 member_level 表的数据

BASE_URL="http://localhost:8080/api"

echo "=== 测试用户等级配置功能 ==="

# 1. 获取所有等级配置
echo "1. 获取所有等级配置"
curl -X GET "${BASE_URL}/member-levels" \
  -H "Content-Type: application/json" \
  -w "\n状态码: %{http_code}\n\n"

# 2. 根据等级获取配置
echo "2. 根据等级获取配置 (等级1)"
curl -X GET "${BASE_URL}/member-levels/1" \
  -H "Content-Type: application/json" \
  -w "\n状态码: %{http_code}\n\n"

# 3. 获取用户等级信息 (经验值1)
echo "3. 获取用户等级信息 (经验值1)"
curl -X GET "${BASE_URL}/member-levels/user-info?experience=1" \
  -H "Content-Type: application/json" \
  -w "\n状态码: %{http_code}\n\n"

# 4. 获取用户等级信息 (经验值50)
echo "4. 获取用户等级信息 (经验值50)"
curl -X GET "${BASE_URL}/member-levels/user-info?experience=50" \
  -H "Content-Type: application/json" \
  -w "\n状态码: %{http_code}\n\n"

# 5. 计算返现金额 (经验值1, 金额100)
echo "5. 计算返现金额 (经验值1, 金额100)"
curl -X GET "${BASE_URL}/member-levels/calculate-cashback?experience=1&amount=100" \
  -H "Content-Type: application/json" \
  -w "\n状态码: %{http_code}\n\n"

# 6. 计算返现金额 (经验值50, 金额1000)
echo "6. 计算返现金额 (经验值50, 金额1000)"
curl -X GET "${BASE_URL}/member-levels/calculate-cashback?experience=50&amount=1000" \
  -H "Content-Type: application/json" \
  -w "\n状态码: %{http_code}\n\n"

# 7. 测试无效参数
echo "7. 测试无效参数 - 等级参数无效"
curl -X GET "${BASE_URL}/member-levels/invalid" \
  -H "Content-Type: application/json" \
  -w "\n状态码: %{http_code}\n\n"

echo "8. 测试无效参数 - 经验值参数无效"
curl -X GET "${BASE_URL}/member-levels/user-info?experience=invalid" \
  -H "Content-Type: application/json" \
  -w "\n状态码: %{http_code}\n\n"

echo "9. 测试无效参数 - 金额参数无效"
curl -X GET "${BASE_URL}/member-levels/calculate-cashback?experience=1&amount=invalid" \
  -H "Content-Type: application/json" \
  -w "\n状态码: %{http_code}\n\n"

echo "=== 测试完成 ===" 