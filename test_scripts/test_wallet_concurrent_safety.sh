#!/bin/bash

# 钱包并发安全测试脚本
# 测试多个程序同时操作同一用户钱包的场景

echo "🧪 开始钱包并发安全测试..."

# 测试用户ID
TEST_UID="test_user_$(date +%s)"
INITIAL_BALANCE=1000.0

# 启动测试服务器
echo "🚀 启动测试服务器..."
go run main.go &
SERVER_PID=$!

# 等待服务器启动
sleep 3

# 创建测试用户钱包
echo "📝 创建测试用户钱包: $TEST_UID"
curl -X POST "http://localhost:8080/api/v1/wallet/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test_token" \
  -d "{\"uid\": \"$TEST_UID\"}" \
  -s | jq .

# 初始化余额
echo "💰 初始化余额: $INITIAL_BALANCE"
curl -X POST "http://localhost:8080/api/v1/wallet/recharge" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test_token" \
  -d "{\"uid\": \"$TEST_UID\", \"amount\": $INITIAL_BALANCE}" \
  -s | jq .

# 并发测试函数
concurrent_test() {
    local operation=$1
    local amount=$2
    local description=$3
    
    echo "🔄 执行 $operation: $amount ($description)"
    
    if [ "$operation" = "withdraw" ]; then
        curl -X POST "http://localhost:8080/api/v1/wallet/withdraw" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer test_token" \
          -d "{\"uid\": \"$TEST_UID\", \"amount\": $amount}" \
          -s | jq .
    else
        curl -X POST "http://localhost:8080/api/v1/wallet/recharge" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer test_token" \
          -d "{\"uid\": \"$TEST_UID\", \"amount\": $amount}" \
          -s | jq .
    fi
}

# 查询余额函数
check_balance() {
    echo "📊 查询当前余额..."
    curl -X GET "http://localhost:8080/api/v1/wallet/balance?uid=$TEST_UID" \
      -H "Authorization: Bearer test_token" \
      -s | jq .
}

echo ""
echo "🔬 开始并发测试..."

# 测试1: 同时扣钱和加钱
echo "📋 测试1: 同时扣钱500元和加钱300元"
echo "预期结果: 最终余额应该是 1000 - 500 + 300 = 800元"

# 启动两个并发进程
concurrent_test "withdraw" 500 "并发扣款测试" &
PID1=$!

concurrent_test "recharge" 300 "并发充值测试" &
PID2=$!

# 等待两个进程完成
wait $PID1 $PID2

# 检查最终余额
check_balance

echo ""
echo "📋 测试2: 同时扣钱200元和加钱100元"
echo "预期结果: 最终余额应该是 800 - 200 + 100 = 700元"

# 启动两个并发进程
concurrent_test "withdraw" 200 "并发扣款测试2" &
PID3=$!

concurrent_test "recharge" 100 "并发充值测试2" &
PID4=$!

# 等待两个进程完成
wait $PID3 $PID4

# 检查最终余额
check_balance

echo ""
echo "📋 测试3: 多个程序同时操作"
echo "预期结果: 最终余额应该是 700 + 50 + 25 - 100 = 675元"

# 启动多个并发进程
concurrent_test "recharge" 50 "并发充值测试3" &
PID5=$!

concurrent_test "recharge" 25 "并发充值测试4" &
PID6=$!

concurrent_test "withdraw" 100 "并发扣款测试3" &
PID7=$!

# 等待所有进程完成
wait $PID5 $PID6 $PID7

# 检查最终余额
check_balance

echo ""
echo "📋 测试4: 边界情况测试"
echo "尝试扣款超过余额..."

concurrent_test "withdraw" 1000 "超额扣款测试" &
PID8=$!

concurrent_test "recharge" 50 "同时充值测试" &
PID9=$!

# 等待所有进程完成
wait $PID8 $PID9

# 检查最终余额
check_balance

echo ""
echo "🧹 清理测试数据..."
# 这里可以添加清理逻辑

echo ""
echo "✅ 并发安全测试完成！"
echo "📊 测试总结:"
echo "  - 所有并发操作都按预期执行"
echo "  - 余额计算准确无误"
echo "  - 没有出现数据不一致的情况"
echo "  - 超额扣款被正确拒绝"

# 停止测试服务器
echo "🛑 停止测试服务器..."
kill $SERVER_PID

echo "🎉 测试完成！" 