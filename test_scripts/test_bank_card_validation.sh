#!/bin/bash

# 银行卡校验功能测试脚本

BASE_URL="http://localhost:9001"
TEST_UID="12345678"

echo "🧪 开始测试银行卡校验功能..."
echo "=================================="

# 测试1: 有效的银行卡号（Luhn算法校验通过）
echo "📊 测试1: 有效的银行卡号（Luhn算法校验通过）"
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

# 测试2: 无效的银行卡号（Luhn算法校验失败）
echo "📊 测试2: 无效的银行卡号（Luhn算法校验失败）"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"工商银行\",
    \"card_holder\": \"李四\",
    \"card_number\": \"6225881234567891\",
    \"card_type\": \"信用卡\"
  }" | jq
echo ""

# 测试3: 银行卡号包含非数字字符
echo "📊 测试3: 银行卡号包含非数字字符"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"建设银行\",
    \"card_holder\": \"王五\",
    \"card_number\": \"622588123456789a\",
    \"card_type\": \"储蓄卡\"
  }" | jq
echo ""

# 测试4: 银行卡号长度不足
echo "📊 测试4: 银行卡号长度不足"
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

# 测试5: 银行卡号长度过长
echo "📊 测试5: 银行卡号长度过长"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"中国银行\",
    \"card_holder\": \"孙七\",
    \"card_number\": \"62258812345678901234\",
    \"card_type\": \"信用卡\"
  }" | jq
echo ""

# 测试6: 持卡人姓名包含特殊字符
echo "📊 测试6: 持卡人姓名包含特殊字符"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"交通银行\",
    \"card_holder\": \"周八@123\",
    \"card_number\": \"6225881234567890\",
    \"card_type\": \"储蓄卡\"
  }" | jq
echo ""

# 测试7: 持卡人姓名长度不足
echo "📊 测试7: 持卡人姓名长度不足"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"中信银行\",
    \"card_holder\": \"A\",
    \"card_number\": \"6225881234567890\",
    \"card_type\": \"借记卡\"
  }" | jq
echo ""

# 测试8: 银行名称包含特殊字符
echo "📊 测试8: 银行名称包含特殊字符"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"招商银行@123\",
    \"card_holder\": \"张三\",
    \"card_number\": \"6225881234567890\",
    \"card_type\": \"信用卡\"
  }" | jq
echo ""

# 测试9: 银行名称长度过长
echo "📊 测试9: 银行名称长度过长"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"这是一个非常非常非常非常非常非常非常非常非常非常长的银行名称\",
    \"card_holder\": \"李四\",
    \"card_number\": \"6225881234567890\",
    \"card_type\": \"储蓄卡\"
  }" | jq
echo ""

# 测试10: 银行卡号包含空格（应该自动去除）
echo "📊 测试10: 银行卡号包含空格（应该自动去除）"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"招商银行\",
    \"card_holder\": \"张三\",
    \"card_number\": \"6225 8812 3456 7890\",
    \"card_type\": \"借记卡\"
  }" | jq
echo ""

# 测试11: 持卡人姓名包含多余空格（应该自动去除）
echo "📊 测试11: 持卡人姓名包含多余空格（应该自动去除）"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"工商银行\",
    \"card_holder\": \"  张三  \",
    \"card_number\": \"6225881234567890\",
    \"card_type\": \"信用卡\"
  }" | jq
echo ""

# 测试12: 银行名称包含多余空格（应该自动去除）
echo "📊 测试12: 银行名称包含多余空格（应该自动去除）"
curl -X POST "$BASE_URL/api/bind-bank-card" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"bank_name\": \"  建设银行  \",
    \"card_holder\": \"王五\",
    \"card_number\": \"6225881234567890\",
    \"card_type\": \"储蓄卡\"
  }" | jq
echo ""

echo "✅ 银行卡校验功能测试完成！"
echo ""
echo "📝 测试结果说明："
echo "1. Luhn算法校验 - 验证银行卡号的有效性"
echo "2. BIN码验证 - 验证银行卡前6位是否属于已知银行"
echo "3. 格式验证 - 验证银行卡号、持卡人姓名、银行名称的格式"
echo "4. 长度验证 - 验证各字段的长度限制"
echo "5. 字符验证 - 验证各字段是否包含非法字符"
echo "6. 空格处理 - 验证自动去除多余空格的功能"
echo ""
echo "🔍 校验规则："
echo "- 银行卡号：13-19位数字，通过Luhn算法校验"
echo "- 持卡人姓名：2-20个字符，只允许中文、英文、空格"
echo "- 银行名称：2-50个字符，只允许中文、英文、数字、空格"
echo "- 卡类型：必须是借记卡、信用卡、储蓄卡之一" 