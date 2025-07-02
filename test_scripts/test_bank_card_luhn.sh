#!/bin/bash

# 测试银行卡号Luhn算法验证

echo "=== 测试银行卡号Luhn算法验证 ==="

# 测试卡号
CARD_NUMBER="6222600234567890123"

echo "测试卡号: $CARD_NUMBER"

# 手动实现Luhn算法验证
luhn_check() {
    local card_number="$1"
    local sum=0
    local alternate=false
    
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
    
    if [ $((sum % 10)) -eq 0 ]; then
        echo "✓ Luhn算法验证通过"
        return 0
    else
        echo "✗ Luhn算法验证失败"
        return 1
    fi
}

# 执行验证
luhn_check "$CARD_NUMBER"

echo ""
echo "=== 生成一个有效的招商银行借记卡号 ==="

# 招商银行借记卡BIN: 621286
# 生成一个有效的卡号
VALID_CARD="6212861234567890123"
echo "有效卡号示例: $VALID_CARD"
luhn_check "$VALID_CARD"

echo ""
echo "=== 其他银行BIN示例 ==="
echo "招商银行借记卡: 621286"
echo "招商银行信用卡: 622580"
echo "中国银行借记卡: 621660"
echo "中国银行信用卡: 622760"
echo "工商银行借记卡: 622202"
echo "工商银行信用卡: 622230" 