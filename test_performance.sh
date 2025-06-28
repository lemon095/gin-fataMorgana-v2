#!/bin/bash

# 数据库性能测试脚本
# 用于测试数据库优化效果

BASE_URL="http://localhost:8080"
TEST_COUNT=100

echo "🚀 开始数据库性能测试..."
echo "=================================="

# 测试1: 健康检查
echo "📊 测试1: 系统健康检查"
time for i in $(seq 1 $TEST_COUNT); do
    curl -s "$BASE_URL/health/system" > /dev/null
done
echo "完成 $TEST_COUNT 次健康检查"
echo ""

# 测试2: 数据库统计信息
echo "📊 测试2: 数据库统计信息"
time for i in $(seq 1 $TEST_COUNT); do
    curl -s "$BASE_URL/health/db-stats" > /dev/null
done
echo "完成 $TEST_COUNT 次数据库统计查询"
echo ""

# 测试3: 查询统计信息
echo "📊 测试3: 查询统计信息"
time for i in $(seq 1 $TEST_COUNT); do
    curl -s "$BASE_URL/health/query-stats" > /dev/null
done
echo "完成 $TEST_COUNT 次查询统计查询"
echo ""

# 测试4: 性能优化建议
echo "📊 测试4: 性能优化建议"
time for i in $(seq 1 $TEST_COUNT); do
    curl -s "$BASE_URL/health/optimization" > /dev/null
done
echo "完成 $TEST_COUNT 次性能优化建议查询"
echo ""

# 测试5: 用户注册（测试数据库写入性能）
echo "📊 测试5: 用户注册性能测试"
time for i in $(seq 1 10); do
    curl -s -X POST "$BASE_URL/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\": \"test$i@example.com\",
            \"password\": \"123456\",
            \"confirm_password\": \"123456\",
            \"invite_code\": \"TEST$i\"
        }" > /dev/null
done
echo "完成 10 次用户注册测试"
echo ""

# 测试6: 并发测试
echo "📊 测试6: 并发健康检查测试"
time for i in $(seq 1 50); do
    curl -s "$BASE_URL/health/system" > /dev/null &
done
wait
echo "完成 50 次并发健康检查"
echo ""

# 测试7: 缓存效果测试
echo "📊 测试7: 缓存效果测试"
echo "第一次查询（缓存未命中）:"
time curl -s "$BASE_URL/health/db-stats" | jq '.data.database_stats' > /dev/null

echo "第二次查询（缓存命中）:"
time curl -s "$BASE_URL/health/db-stats" | jq '.data.database_stats' > /dev/null
echo ""

# 显示当前系统状态
echo "📊 当前系统状态:"
echo "数据库统计信息:"
curl -s "$BASE_URL/health/db-stats" | jq '.data.database_stats'

echo ""
echo "查询统计信息:"
curl -s "$BASE_URL/health/query-stats" | jq '.data.query_stats'

echo ""
echo "性能优化建议:"
curl -s "$BASE_URL/health/optimization" | jq '.data.optimization_recommendations'

echo ""
echo "✅ 性能测试完成！"
echo "=================================="
echo "💡 优化建议:"
echo "1. 观察连接池使用情况"
echo "2. 检查缓存命中率"
echo "3. 监控查询响应时间"
echo "4. 根据负载调整连接池参数" 