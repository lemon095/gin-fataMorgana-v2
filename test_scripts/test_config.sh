#!/bin/bash

# 配置测试脚本

echo "=== 运行配置测试 ==="

# 编译配置测试程序
echo "🔨 编译配置测试程序..."
go build -o test_config cmd/test_config/main.go

if [ $? -eq 0 ]; then
    echo "✅ 编译成功"
    
    # 运行配置测试
    echo "🧪 运行配置测试..."
    ./test_config
    
    # 清理
    rm -f test_config
else
    echo "❌ 编译失败"
    exit 1
fi

echo "=== 配置测试完成 ===" 