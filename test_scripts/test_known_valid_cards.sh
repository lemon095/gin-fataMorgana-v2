#!/bin/bash

# 测试已知有效的银行卡号

echo "=== 测试已知有效的银行卡号 ==="

# 一些已知有效的银行卡号（用于测试）
VALID_CARDS=(
    "6222021234567890123"  # 工商银行借记卡
    "6227601234567890123"  # 中国银行信用卡
    "6225801234567890123"  # 招商银行信用卡
    "6212861234567890123"  # 招商银行借记卡
    "6216601234567890123"  # 中国银行借记卡
)

luhn_check() {
    local card_number="$1"
    local sum=0
    local alternate=false
    
    echo "验证卡号: $card_number"
    
    # 从右到左遍历
    for ((i=${#card_number}-1; i>=0; i--)); do
        local digit=${card_number:$i:1}
        
        if [ "$alternate" = true ]; then
            digit=$((digit * 2))
            if [ $digit -gt 9 ]; then
                digit=$((digit % 10 + digit / 10))
            fi
        fi
        
        sum=$((sum + digit))
        alternate=$([ "$alternate" = true ] && echo false || echo true)
    done
    
    echo "校验和: $sum"
    
    if [ $((sum % 10)) -eq 0 ]; then
        echo "✓ Luhn算法验证通过"
        return 0
    else
        echo "✗ Luhn算法验证失败"
        return 1
    fi
}

# 测试每个卡号
for card in "${VALID_CARDS[@]}"; do
    echo ""
    luhn_check "$card"
done

echo ""
echo "=== 生成一个简单的有效卡号 ==="

# 使用一个简单的BIN生成有效卡号
# 招商银行借记卡BIN: 621286
# 让我们手动计算一个有效的卡号

# 卡号: 6212860000000000000
# 手动计算校验位
CARD_TEST="621286000000000000"
echo "测试卡号（不含校验位）: $CARD_TEST"

# 手动计算Luhn算法
sum=0
alternate=false
for ((i=${#CARD_TEST}-1; i>=0; i--)); do
    digit=${CARD_TEST:$i:1}
    
    if [ "$alternate" = true ]; then
        digit=$((digit * 2))
        if [ $digit -gt 9 ]; then
            digit=$((digit % 10 + digit / 10))
        fi
    fi
    
    sum=$((sum + digit))
    alternate=$([ "$alternate" = true ] && echo false || echo true)
done

echo "校验和: $sum"
checksum=$((10 - (sum % 10)))
if [ $checksum -eq 10 ]; then
    checksum=0
fi

VALID_CARD="${CARD_TEST}${checksum}"
echo "计算的校验位: $checksum"
echo "完整有效卡号: $VALID_CARD"

# 验证
echo ""
luhn_check "$VALID_CARD" 