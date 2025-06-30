#!/bin/bash

# 余额不足和多笔提现测试脚本

BASE_URL="http://localhost:9001"
TEST_UID="12345678"

echo "🧪 开始测试余额不足和多笔提现情况..."
echo "=================================="

# 测试1: 创建钱包
echo "📊 测试1: 创建钱包"
curl -X POST "$BASE_URL/api/wallet/$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" | jq
echo ""

# 测试2: 充值100元
echo "📊 测试2: 充值100元"
curl -X POST "$BASE_URL/api/wallet/recharge" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 100.00,
    \"description\": \"银行卡充值\"
  }" | jq
echo ""

# 测试3: 查看钱包余额
echo "📊 测试3: 查看钱包余额"
curl -X GET "$BASE_URL/api/wallet/$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试4: 申请提现50元（应该成功）
echo "📊 测试4: 申请提现50元（应该成功）"
RESPONSE1=$(curl -s -X POST "$BASE_URL/api/wallet/request-withdraw" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 50.00,
    \"description\": \"提现到银行卡1\",
    \"bank_card_info\": \"招商银行 6225****1234 张三\"
  }")

echo "$RESPONSE1" | jq
echo ""

# 测试5: 申请提现30元（应该成功）
echo "📊 测试5: 申请提现30元（应该成功）"
RESPONSE2=$(curl -s -X POST "$BASE_URL/api/wallet/request-withdraw" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 30.00,
    \"description\": \"提现到银行卡2\",
    \"bank_card_info\": \"工商银行 6222****5678 李四\"
  }")

echo "$RESPONSE2" | jq
echo ""

# 测试6: 查看钱包余额（应该显示冻结金额）
echo "📊 测试6: 查看钱包余额（应该显示冻结金额）"
curl -X GET "$BASE_URL/api/wallet/$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试7: 申请提现30元（应该失败，余额不足）
echo "📊 测试7: 申请提现30元（应该失败，余额不足）"
curl -X POST "$BASE_URL/api/wallet/request-withdraw" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 30.00,
    \"description\": \"提现到银行卡3\",
    \"bank_card_info\": \"建设银行 6227****9999 王五\"
  }" | jq
echo ""

# 测试8: 申请提现200元（应该失败，总余额不足）
echo "📊 测试8: 申请提现200元（应该失败，总余额不足）"
curl -X POST "$BASE_URL/api/wallet/request-withdraw" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 200.00,
    \"description\": \"提现到银行卡4\",
    \"bank_card_info\": \"农业银行 6228****1111 赵六\"
  }" | jq
echo ""

# 测试9: 申请提现0元（应该失败，金额无效）
echo "📊 测试9: 申请提现0元（应该失败，金额无效）"
curl -X POST "$BASE_URL/api/wallet/request-withdraw" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 0.00,
    \"description\": \"提现到银行卡5\",
    \"bank_card_info\": \"交通银行 6222****2222 孙七\"
  }" | jq
echo ""

# 测试10: 申请提现-10元（应该失败，金额无效）
echo "📊 测试10: 申请提现-10元（应该失败，金额无效）"
curl -X POST "$BASE_URL/api/wallet/request-withdraw" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": -10.00,
    \"description\": \"提现到银行卡6\",
    \"bank_card_info\": \"中信银行 6226****3333 周八\"
  }" | jq
echo ""

# 测试11: 获取提现汇总信息
echo "📊 测试11: 获取提现汇总信息"
curl -X GET "$BASE_URL/api/wallet/withdraw-summary?uid=$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试12: 查看交易记录
echo "📊 测试12: 查看交易记录"
curl -X GET "$BASE_URL/api/wallet/transactions?uid=$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

echo "✅ 余额不足和多笔提现测试完成！"
echo ""
echo "📝 测试结果说明："
echo "1. 用户有100元余额"
echo "2. 申请提现50元和30元都成功（冻结80元）"
echo "3. 再申请30元失败（可用余额只有20元）"
echo "4. 申请200元失败（总余额不足）"
echo "5. 申请0元和负数金额失败（金额无效）"
echo "6. 汇总信息显示当前钱包状态和提现统计" 