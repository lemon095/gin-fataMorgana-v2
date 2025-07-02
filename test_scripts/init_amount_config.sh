#!/bin/bash

# 初始化金额配置数据脚本
echo "=== 初始化金额配置数据 ==="

# 检查配置文件
if [ ! -f "config.yaml" ]; then
    echo "❌ 配置文件不存在，请先创建 config.yaml"
    exit 1
fi

# 获取服务器地址
SERVER_HOST="localhost"
SERVER_PORT=$(grep "port:" config.yaml | head -1 | awk '{print $2}')

# 设置默认值
SERVER_PORT=${SERVER_PORT:-9001}

echo "📋 服务器配置:"
echo "  主机: $SERVER_HOST"
echo "  端口: $SERVER_PORT"
echo

# 检查服务是否运行
echo "🔍 检查服务状态..."
if ! curl -s "http://$SERVER_HOST:$SERVER_PORT/health" > /dev/null; then
    echo "❌ 服务未运行，请先启动服务:"
    echo "   ./dev.sh start"
    echo "   或者"
    echo "   ./prod.sh start"
    exit 1
fi
echo "✅ 服务运行正常"
echo

# 初始化充值金额配置
echo "💰 初始化充值金额配置..."

# 充值金额配置数据
RECHARGE_AMOUNTS=(
    '{"type": "recharge", "amount": 100.00, "description": "充值100元", "is_active": true, "sort_order": 1}'
    '{"type": "recharge", "amount": 200.00, "description": "充值200元", "is_active": true, "sort_order": 2}'
    '{"type": "recharge", "amount": 500.00, "description": "充值500元", "is_active": true, "sort_order": 3}'
    '{"type": "recharge", "amount": 1000.00, "description": "充值1000元", "is_active": true, "sort_order": 4}'
    '{"type": "recharge", "amount": 2000.00, "description": "充值2000元", "is_active": true, "sort_order": 5}'
    '{"type": "recharge", "amount": 5000.00, "description": "充值5000元", "is_active": true, "sort_order": 6}'
)

for config in "${RECHARGE_AMOUNTS[@]}"; do
    echo "  添加充值配置: $config"
    # 这里需要先实现创建接口，暂时跳过
    # curl -X POST "http://$SERVER_HOST:$SERVER_PORT/api/v1/amount-config/create" \
    #   -H "Content-Type: application/json" \
    #   -d "$config"
done

# 初始化提现金额配置
echo "💸 初始化提现金额配置..."

# 提现金额配置数据
WITHDRAW_AMOUNTS=(
    '{"type": "withdraw", "amount": 50.00, "description": "提现50元", "is_active": true, "sort_order": 1}'
    '{"type": "withdraw", "amount": 100.00, "description": "提现100元", "is_active": true, "sort_order": 2}'
    '{"type": "withdraw", "amount": 200.00, "description": "提现200元", "is_active": true, "sort_order": 3}'
    '{"type": "withdraw", "amount": 500.00, "description": "提现500元", "is_active": true, "sort_order": 4}'
    '{"type": "withdraw", "amount": 1000.00, "description": "提现1000元", "is_active": true, "sort_order": 5}'
)

for config in "${WITHDRAW_AMOUNTS[@]}"; do
    echo "  添加提现配置: $config"
    # 这里需要先实现创建接口，暂时跳过
    # curl -X POST "http://$SERVER_HOST:$SERVER_PORT/api/v1/amount-config/create" \
    #   -H "Content-Type: application/json" \
    #   -d "$config"
done

echo
echo "📝 注意: 由于尚未实现创建接口，请手动在数据库中插入配置数据"
echo
echo "💡 SQL示例:"
echo "INSERT INTO amount_config (type, amount, description, is_active, sort_order) VALUES"
echo "('recharge', 100.00, '充值100元', 1, 1),"
echo "('recharge', 200.00, '充值200元', 1, 2),"
echo "('recharge', 500.00, '充值500元', 1, 3),"
echo "('recharge', 1000.00, '充值1000元', 1, 4),"
echo "('recharge', 2000.00, '充值2000元', 1, 5),"
echo "('recharge', 5000.00, '充值5000元', 1, 6),"
echo "('withdraw', 50.00, '提现50元', 1, 1),"
echo "('withdraw', 100.00, '提现100元', 1, 2),"
echo "('withdraw', 200.00, '提现200元', 1, 3),"
echo "('withdraw', 500.00, '提现500元', 1, 4),"
echo "('withdraw', 1000.00, '提现1000元', 1, 5);"
echo
echo "🎉 初始化脚本完成！" 