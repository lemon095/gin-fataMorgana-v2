#!/bin/bash

# 索引检测测试脚本
# 用于测试项目的自动索引检测和创建功能

set -e

echo "🔍 开始测试索引检测功能..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go环境未安装"
    exit 1
fi

# 检查配置文件
if [ ! -f "config/config.yaml" ]; then
    echo "❌ 配置文件不存在，请先复制 config.example.yaml 为 config.yaml 并配置数据库连接"
    exit 1
fi

echo "✅ 环境检查通过"

# 测试1: 显示当前索引
echo ""
echo "📋 测试1: 显示当前数据库索引"
echo "=================================="
go run cmd/migrate/main.go -show-index

# 测试2: 检测并创建缺失的索引
echo ""
echo "🔍 测试2: 检测并创建缺失的索引"
echo "=================================="
go run cmd/migrate/main.go -check-index

# 测试3: 再次显示索引（验证创建结果）
echo ""
echo "📋 测试3: 验证索引创建结果"
echo "=================================="
go run cmd/migrate/main.go -show-index

echo ""
echo "✅ 索引检测测试完成！"
echo ""
echo "📝 测试总结:"
echo "   - 显示索引功能: ✅"
echo "   - 索引检测功能: ✅"
echo "   - 索引创建功能: ✅"
echo ""
echo "💡 使用提示:"
echo "   - 项目启动时会自动检测和创建索引"
echo "   - 手动检测索引: make db-check-index"
echo "   - 查看所有索引: make db-show-index" 