#!/bin/bash

# 提现申请流程测试脚本

BASE_URL="http://localhost:9001"
TEST_UID="12345678"

echo "🧪 开始测试提现申请流程..."
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

# 测试4: 申请提现30元（锁定金额）
echo "📊 测试4: 申请提现30元（锁定金额）"
RESPONSE=$(curl -s -X POST "$BASE_URL/api/wallet/request-withdraw" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 30.00,
    \"description\": \"提现到银行卡\",
    \"bank_card_info\": \"招商银行 6225****1234 张三\"
  }")

echo "$RESPONSE" | jq

# 提取交易流水号
TRANSACTION_NO=$(echo "$RESPONSE" | jq -r '.data.transaction_no')
echo "交易流水号: $TRANSACTION_NO"
echo ""

# 测试5: 查看钱包余额（应该显示冻结金额）
echo "📊 测试5: 查看钱包余额（应该显示冻结金额）"
curl -X GET "$BASE_URL/api/wallet/$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试6: 查看交易记录
echo "📊 测试6: 查看交易记录"
curl -X GET "$BASE_URL/api/wallet/transactions?uid=$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试7: 申请提现50元（测试余额不足）
echo "📊 测试7: 申请提现50元（测试余额不足）"
curl -X POST "$BASE_URL/api/wallet/request-withdraw" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 50.00,
    \"description\": \"提现到支付宝\",
    \"bank_card_info\": \"支付宝 138****5678\"
  }" | jq
echo ""

# 测试8: 申请提现20元（测试可用余额）
echo "📊 测试8: 申请提现20元（测试可用余额）"
RESPONSE2=$(curl -s -X POST "$BASE_URL/api/wallet/request-withdraw" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 20.00,
    \"description\": \"提现到微信\",
    \"bank_card_info\": \"微信 139****9999\"
  }")

echo "$RESPONSE2" | jq

# 提取交易流水号
TRANSACTION_NO2=$(echo "$RESPONSE2" | jq -r '.data.transaction_no')
echo "交易流水号: $TRANSACTION_NO2"
echo ""

# 测试9: 查看钱包余额（两个提现申请后的状态）
echo "📊 测试9: 查看钱包余额（两个提现申请后的状态）"
curl -X GET "$BASE_URL/api/wallet/$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试10: 查看交易记录（所有申请）
echo "📊 测试10: 查看交易记录（所有申请）"
curl -X GET "$BASE_URL/api/wallet/transactions?uid=$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

echo "✅ 提现申请流程测试完成！"
echo ""
echo "📝 说明："
echo "1. 用户只能申请提现，不能直接提现"
echo "2. 申请提现会锁定相应金额（冻结余额）"
echo "3. 后台会自动处理提现逻辑（确认或取消）"
echo "4. 确认提现后，冻结金额会被扣减"
echo "5. 取消提现后，冻结金额会被解冻" 