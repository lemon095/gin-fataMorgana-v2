#!/bin/bash

# 生成有效的银行卡号

echo "=== 生成有效的银行卡号 ==="

# 招商银行借记卡BIN: 621286
BIN="621286"

# 生成中间数字（除最后一位校验位）
MIDDLE="123456789012"

# 计算校验位
calculate_checksum() {
    local card_without_checksum="$1"
    local sum=0
    local alternate=false
    
    # 从右到左遍历（不包括校验位）
    for ((i=${#card_without_checksum}-1; i>=0; i--)); do
        local digit=${card_without_checksum:$i:1}
        
        if [ "$alternate" = true ]; then
            digit=$((digit * 2))
            if [ $digit -gt 9 ]; then
                digit=$((digit % 10 + digit / 10))
            fi
        fi
        
        sum=$((sum + digit))
        alternate=$([ "$alternate" = true ] && echo false || echo true)
    done
    
    # 计算校验位
    local checksum=$((10 - (sum % 10)))
    if [ $checksum -eq 10 ]; then
        checksum=0
    fi
    
    echo $checksum
}

# 生成卡号（不含校验位）
CARD_WITHOUT_CHECKSUM="${BIN}${MIDDLE}"

# 计算校验位
CHECKSUM=$(calculate_checksum "$CARD_WITHOUT_CHECKSUM")

# 完整卡号
VALID_CARD="${CARD_WITHOUT_CHECKSUM}${CHECKSUM}"

echo "招商银行借记卡BIN: $BIN"
echo "中间数字: $MIDDLE"
echo "计算的校验位: $CHECKSUM"
echo "完整有效卡号: $VALID_CARD"

# 验证生成的卡号
echo ""
echo "=== 验证生成的卡号 ==="

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

luhn_check "$VALID_CARD"

echo ""
echo "=== 测试数据 ==="
echo "请使用以下有效的银行卡号进行测试："
echo "招商银行借记卡: $VALID_CARD" 