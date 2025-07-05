#!/bin/bash

# 测试订单号唯一性
echo "🧪 测试订单号唯一性..."

# 编译测试程序
echo "🔨 编译测试程序..."
go build -o test_order_no test_order_no_uniqueness.go

if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi

# 运行测试
echo "🚀 运行订单号唯一性测试..."
./test_order_no

# 清理
rm -f test_order_no

echo "✅ 测试完成" 