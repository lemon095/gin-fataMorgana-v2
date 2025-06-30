#!/bin/bash

# 后台提现管理测试脚本

BASE_URL="http://localhost:9001"
TEST_UID="12345678"

echo "🧪 开始测试后台提现管理..."
echo "=================================="

# 测试1: 查看待处理的提现申请
echo "📊 测试1: 查看待处理的提现申请"
curl -X GET "$BASE_URL/api/wallet/transactions?uid=$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试2: 确认第一个提现申请
echo "📊 测试2: 确认第一个提现申请"
# 这里需要手动输入交易流水号，或者从上面的查询结果中获取
echo "请输入要确认的交易流水号:"
read TRANSACTION_NO

curl -X POST "$BASE_URL/api/wallet/confirm-withdraw" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"transaction_no\": \"$TRANSACTION_NO\"
  }" | jq
echo ""

# 测试3: 查看钱包余额（确认提现后）
echo "📊 测试3: 查看钱包余额（确认提现后）"
curl -X GET "$BASE_URL/api/wallet/$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试4: 查看交易记录（确认提现后）
echo "📊 测试4: 查看交易记录（确认提现后）"
curl -X GET "$BASE_URL/api/wallet/transactions?uid=$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试5: 取消第二个提现申请
echo "📊 测试5: 取消第二个提现申请"
echo "请输入要取消的交易流水号:"
read TRANSACTION_NO2

curl -X POST "$BASE_URL/api/wallet/cancel-withdraw" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"transaction_no\": \"$TRANSACTION_NO2\",
    \"reason\": \"后台审核不通过\"
  }" | jq
echo ""

# 测试6: 查看钱包余额（取消提现后）
echo "📊 测试6: 查看钱包余额（取消提现后）"
curl -X GET "$BASE_URL/api/wallet/$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试7: 查看交易记录（取消提现后）
echo "📊 测试7: 查看交易记录（取消提现后）"
curl -X GET "$BASE_URL/api/wallet/transactions?uid=$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

echo "✅ 后台提现管理测试完成！"
echo ""
echo "📝 说明："
echo "1. 后台可以查看所有待处理的提现申请"
echo "2. 后台可以确认提现申请，完成提现"
echo "3. 后台可以取消提现申请，解冻金额"
echo "4. 确认提现后，冻结金额会被扣减，余额减少"
echo "5. 取消提现后，冻结金额会被解冻，余额恢复" 