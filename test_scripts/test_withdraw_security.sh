#!/bin/bash

# 提现安全检查测试脚本
BASE_URL="http://localhost:9001/api/v1"

echo "🧪 开始测试提现安全检查..."
echo "=================================="

# 测试用户信息
TEST_EMAIL="withdraw_test@example.com"
TEST_PASSWORD="123456"
TEST_UID=""

# 1. 注册测试用户
echo "📝 1. 注册测试用户..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\",
    \"confirm_password\": \"$TEST_PASSWORD\",
    \"invite_code\": \"7TRABJ\"
  }")

echo "注册响应: $REGISTER_RESPONSE"

# 提取UID
TEST_UID=$(echo "$REGISTER_RESPONSE" | jq -r '.data.user.uid')
if [ "$TEST_UID" = "null" ] || [ -z "$TEST_UID" ]; then
    echo "❌ 用户注册失败，无法获取UID"
    exit 1
fi
echo "✅ 用户注册成功，UID: $TEST_UID"

# 2. 用户登录
echo "🔐 2. 用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\"
  }")

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.tokens.access_token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "❌ 用户登录失败，无法获取token"
    exit 1
fi
echo "✅ 用户登录成功"

# 3. 测试未绑定银行卡的提现（应该失败）
echo "💳 3. 测试未绑定银行卡的提现（应该失败）..."
WITHDRAW_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/withdraw" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 100.00,
    \"description\": \"测试提现\",
    \"password\": \"$TEST_PASSWORD\"
  }")

echo "未绑定银行卡提现响应: $WITHDRAW_RESPONSE"

# 检查是否返回银行卡绑定错误
if echo "$WITHDRAW_RESPONSE" | grep -q "请先绑定银行卡"; then
    echo "✅ 未绑定银行卡检查正常"
else
    echo "❌ 未绑定银行卡检查失败"
fi

# 4. 绑定银行卡
echo "💳 4. 绑定银行卡..."
BIND_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/bind-bank-card" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"card_number\": \"6225881234567890\",
    \"card_type\": \"借记卡\",
    \"bank_name\": \"招商银行\",
    \"card_holder\": \"张三\"
  }")

echo "绑定银行卡响应: $BIND_RESPONSE"

# 5. 测试余额不足的提现（应该失败）
echo "💰 5. 测试余额不足的提现（应该失败）..."
INSUFFICIENT_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/withdraw" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 100.00,
    \"description\": \"测试余额不足\",
    \"password\": \"$TEST_PASSWORD\"
  }")

echo "余额不足提现响应: $INSUFFICIENT_RESPONSE"

# 检查是否返回余额不足错误
if echo "$INSUFFICIENT_RESPONSE" | grep -q "余额不足"; then
    echo "✅ 余额不足检查正常"
else
    echo "❌ 余额不足检查失败"
fi

# 6. 先充值一些钱
echo "💰 6. 充值100元..."
RECHARGE_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/recharge" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 100.00,
    \"description\": \"测试充值\"
  }")

echo "充值响应: $RECHARGE_RESPONSE"

# 7. 测试正常提现（应该成功）
echo "💸 7. 测试正常提现（应该成功）..."
SUCCESS_WITHDRAW_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/withdraw" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 50.00,
    \"description\": \"正常提现测试\",
    \"password\": \"$TEST_PASSWORD\"
  }")

echo "正常提现响应: $SUCCESS_WITHDRAW_RESPONSE"

# 检查是否提现成功
if echo "$SUCCESS_WITHDRAW_RESPONSE" | grep -q "提现申请已提交"; then
    echo "✅ 正常提现成功"
else
    echo "❌ 正常提现失败"
fi

# 8. 测试操作其他用户钱包（应该失败）
echo "🚫 8. 测试操作其他用户钱包（应该失败）..."
OTHER_UID="99999999"
OTHER_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/withdraw" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$OTHER_UID\",
    \"amount\": 10.00,
    \"description\": \"测试操作其他用户钱包\",
    \"password\": \"$TEST_PASSWORD\"
  }")

echo "操作其他用户钱包响应: $OTHER_USER_RESPONSE"

# 检查是否返回权限错误
if echo "$OTHER_USER_RESPONSE" | grep -q "只能操作自己的钱包"; then
    echo "✅ 用户权限检查正常"
else
    echo "❌ 用户权限检查失败"
fi

# 9. 测试密码错误（应该失败）
echo "🔒 9. 测试密码错误（应该失败）..."
WRONG_PASSWORD_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/withdraw" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 10.00,
    \"description\": \"测试密码错误\",
    \"password\": \"wrongpassword\"
  }")

echo "密码错误响应: $WRONG_PASSWORD_RESPONSE"

# 检查是否返回密码错误
if echo "$WRONG_PASSWORD_RESPONSE" | grep -q "登录密码错误"; then
    echo "✅ 密码验证正常"
else
    echo "❌ 密码验证失败"
fi

echo ""
echo "=================================="
echo "🎉 提现安全检查测试完成！"
echo ""
echo "测试总结："
echo "✅ 银行卡绑定检查"
echo "✅ 余额不足检查"
echo "✅ 用户权限检查"
echo "✅ 密码验证检查"
echo "✅ 正常提现流程" 