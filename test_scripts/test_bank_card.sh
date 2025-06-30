#!/bin/bash

# 银行卡绑定测试脚本

BASE_URL="http://localhost:9001"
TEST_UID="12345678"

echo "🧪 开始测试银行卡绑定功能..."
echo "=================================="

# 测试1: 绑定银行卡（成功）
echo "📊 测试1: 绑定银行卡（成功）"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"招商银行\",
    \"card_holder\": \"张三\",
    \"card_number\": \"6225881234567890\",
    \"card_type\": \"借记卡\"
  }" | jq
echo ""

# 测试2: 获取银行卡信息（成功）
echo "📊 测试2: 获取银行卡信息（成功）"
curl -X GET "$BASE_URL/api/bank-card-info?uid=$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试3: 绑定银行卡（银行名称为空）
echo "📊 测试3: 绑定银行卡（银行名称为空）"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"\",
    \"card_holder\": \"李四\",
    \"card_number\": \"6225881234567891\",
    \"card_type\": \"信用卡\"
  }" | jq
echo ""

# 测试4: 绑定银行卡（持卡人姓名为空）
echo "📊 测试4: 绑定银行卡（持卡人姓名为空）"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"工商银行\",
    \"card_holder\": \"\",
    \"card_number\": \"6225881234567892\",
    \"card_type\": \"储蓄卡\"
  }" | jq
echo ""

# 测试5: 绑定银行卡（卡号为空）
echo "📊 测试5: 绑定银行卡（卡号为空）"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"建设银行\",
    \"card_holder\": \"王五\",
    \"card_number\": \"\",
    \"card_type\": \"借记卡\"
  }" | jq
echo ""

# 测试6: 绑定银行卡（卡号长度不正确）
echo "📊 测试6: 绑定银行卡（卡号长度不正确）"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"农业银行\",
    \"card_holder\": \"赵六\",
    \"card_number\": \"123456789\",
    \"card_type\": \"借记卡\"
  }" | jq
echo ""

# 测试7: 绑定银行卡（卡类型不正确）
echo "📊 测试7: 绑定银行卡（卡类型不正确）"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"交通银行\",
    \"card_holder\": \"孙七\",
    \"card_number\": \"6225881234567893\",
    \"card_type\": \"会员卡\"
  }" | jq
echo ""

# 测试8: 绑定银行卡（用户不存在）
echo "📊 测试8: 绑定银行卡（用户不存在）"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"99999999\",
    \"bank_name\": \"中信银行\",
    \"card_holder\": \"周八\",
    \"card_number\": \"6225881234567894\",
    \"card_type\": \"信用卡\"
  }" | jq
echo ""

# 测试9: 获取银行卡信息（用户不存在）
echo "📊 测试9: 获取银行卡信息（用户不存在）"
curl -X GET "$BASE_URL/api/bank-card-info?uid=99999999" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试10: 获取银行卡信息（参数为空）
echo "📊 测试10: 获取银行卡信息（参数为空）"
curl -X GET "$BASE_URL/api/bank-card-info" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试11: 更新银行卡信息（覆盖原有信息）
echo "📊 测试11: 更新银行卡信息（覆盖原有信息）"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"中国银行\",
    \"card_holder\": \"张三\",
    \"card_number\": \"6225881234567895\",
    \"card_type\": \"储蓄卡\"
  }" | jq
echo ""

# 测试12: 获取更新后的银行卡信息
echo "📊 测试12: 获取更新后的银行卡信息"
curl -X GET "$BASE_URL/api/bank-card-info?uid=$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

echo "✅ 银行卡绑定功能测试完成！"
echo ""
echo "📝 测试结果说明："
echo "1. 成功绑定银行卡，信息存储在用户表的bank_card_info字段"
echo "2. 验证各种参数错误情况（空值、格式错误等）"
echo "3. 验证用户不存在的情况"
echo "4. 验证银行卡信息更新功能"
echo "5. 支持借记卡、信用卡、储蓄卡三种类型" 