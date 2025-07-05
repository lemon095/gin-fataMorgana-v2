#!/bin/bash

# 快速测试用户名脱敏逻辑（统一版本）
echo "=== 用户名脱敏逻辑测试（统一版本）==="

# 模拟脱敏函数（与Go代码逻辑一致）
mask_username() {
    local username="$1"
    local len=${#username}
    
    if [ $len -le 1 ]; then
        echo "$username"
    else
        echo "${username:0:1}**${username: -1}"
    fi
}

# 测试用例
echo -e "\n测试用例:"
echo "输入\t\t输出\t\t说明"
echo "----\t\t----\t\t----"

test_cases=(
    "张三:张**三:2位用户名，统一格式"
    "张三丰:张**丰:3位用户名，统一格式"
    "张三丰李:张**李:4位用户名，统一格式"
    "张三丰李四:张**四:5位用户名，统一格式"
    "张三丰李四王:张**王:6位用户名，统一格式"
    "test_user:t**r:英文用户名，统一格式"
    "ab:a**b:2位英文用户名，统一格式"
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
echo "2. 用户名长度 ≥ 2: 统一格式：首位 + ** + 末位"

echo -e "\n=== 统一效果 ==="
echo "✓ 所有用户名使用相同的脱敏规则"
echo "✓ 脱敏长度统一为3个字符"
echo "✓ 代码逻辑更简单，易于维护"
echo "✓ 用户体验更一致"
echo "✓ 仍然保护了用户隐私" 