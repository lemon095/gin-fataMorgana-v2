#!/bin/bash

# 钱包接口测试脚本

BASE_URL="http://localhost:9001"
TEST_UID="12345678"

echo "🧪 开始测试钱包接口..."
echo "=================================="

# 测试1: 创建钱包
echo "📊 测试1: 创建钱包"
curl -X POST "$BASE_URL/api/wallet/$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" | jq
echo ""

# 测试2: 获取钱包信息
echo "📊 测试2: 获取钱包信息"
curl -X GET "$BASE_URL/api/wallet/$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试3: 充值
echo "📊 测试3: 充值100元"
curl -X POST "$BASE_URL/api/wallet/recharge" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 100.00,
    \"description\": \"银行卡充值\"
  }" | jq
echo ""

# 测试4: 再次充值
echo "📊 测试4: 充值50元"
curl -X POST "$BASE_URL/api/wallet/recharge" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 50.00,
    \"description\": \"支付宝充值\"
  }" | jq
echo ""

# 测试5: 提现
echo "📊 测试5: 提现30元"
curl -X POST "$BASE_URL/api/wallet/withdraw" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": \"$TEST_UID\",
    \"amount\": 30.00,
    \"description\": \"提现到银行卡\"
  }" | jq
echo ""

# 测试6: 获取资金记录（默认分页）
echo "📊 测试6: 获取资金记录（默认分页）"
curl -X GET "$BASE_URL/api/wallet/transactions?uid=$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试7: 获取资金记录（自定义分页）
echo "📊 测试7: 获取资金记录（自定义分页）"
curl -X GET "$BASE_URL/api/wallet/transactions?uid=$TEST_UID&page=1&page_size=5" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

# 测试8: 获取钱包信息（查看余额变化）
echo "📊 测试8: 获取钱包信息（查看余额变化）"
curl -X GET "$BASE_URL/api/wallet/$TEST_UID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
echo ""

echo "✅ 钱包接口测试完成！"
echo "=================================="
echo "💡 测试说明："
echo "1. 请先替换YOUR_TOKEN_HERE为实际的访问令牌"
echo "2. 确保用户已注册并登录"
echo "3. 观察余额和交易记录的变化"
echo "4. 检查分页功能是否正常工作" 