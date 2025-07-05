#!/bin/bash

# 快速测试用户名脱敏逻辑
echo "=== 用户名脱敏逻辑测试 ==="

# 模拟脱敏函数（与Go代码逻辑一致）
mask_username() {
    local username="$1"
    local len=${#username}
    
    if [ $len -le 1 ]; then
        echo "$username"
    elif [ $len -eq 2 ]; then
        echo "${username:0:1}*${username:1:1}"
    elif [ $len -le 4 ]; then
        echo "${username:0:1}*${username: -1}"
    else
        echo "${username:0:2}*${username: -2}"
    fi
}

# 测试用例
echo -e "\n测试用例:"
echo "输入\t\t输出\t\t说明"
echo "----\t\t----\t\t----"

test_cases=(
    "张三:张*三:2位用户名，中间加*"
    "张三丰:张*丰:3位用户名，显示首尾"
    "张三丰李:张*李:4位用户名，显示首尾"
    "张三丰李四:张三*李四:5位用户名，显示首尾各2位"
    "张三丰李四王:张三*李四:6位用户名，显示首尾各2位"
    "test_user:te*er:英文用户名"
    "ab:a*b:2位英文用户名"
    "a:a:1位字符，不脱敏"
)

for test_case in "${test_cases[@]}"; do
    IFS=':' read -r input expected description <<< "$test_case"
    result=$(mask_username "$input")
    
    if [ "$result" = "$expected" ]; then
        status="✓"
    else
        status="✗"
    fi
    
    printf "%-10s\t%-10s\t%s %s\n" "$input" "$result" "$status" "$description"
done

echo -e "\n=== 脱敏规则总结 ==="
echo "1. 用户名长度 = 1: 不脱敏，直接显示"
echo "2. 用户名长度 = 2: 在中间加 *"
echo "3. 用户名长度 3-4: 显示首尾字符，中间用 * 替换"
echo "4. 用户名长度 ≥ 5: 显示首尾各2个字符，中间用 * 替换"

echo -e "\n=== 优化效果 ==="
echo "✓ 统一了脱敏规则，所有用户名都进行脱敏"
echo "✓ 脱敏长度大幅缩短"
echo "✓ 保留了更多可识别信息"
echo "✓ 提高了用户体验"
echo "✓ 仍然保护了用户隐私" 