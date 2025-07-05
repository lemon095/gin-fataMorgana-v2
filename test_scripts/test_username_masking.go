package main

import (
	"fmt"
	"gin-fataMorgana/models"
)

func main() {
	fmt.Println("=== 用户名脱敏功能测试（更新版本）===")
	
	// 测试用例
	testCases := []struct {
		input    string
		expected string
	}{
		{"a", "a"},                         // 单字符，不脱敏
		{"张三", "张*三"},                    // 长度=2，在中间加*
		{"张三丰", "张*丰"},                  // 长度=3，显示首尾
		{"张三丰李", "张*李"},                // 长度=4，显示首尾
		{"张三丰李四", "张三*李四"},           // 长度=5，显示首尾各2个字符
		{"张三丰李四王", "张三*李四"},         // 长度=6，显示首尾各2个字符
		{"ab", "a*b"},                      // 英文，长度=2
		{"test", "t*t"},                    // 英文，长度=4
		{"test_user", "te*er"},             // 英文，长度=9
		{"test_user_123", "te*23"},         // 英文+数字，长度=13
		{"", ""},                           // 空字符串
	}
	
	fmt.Println("\n测试结果:")
	fmt.Println("输入用户名\t\t脱敏结果\t\t预期结果\t\t是否通过")
	fmt.Println("--------\t\t--------\t\t--------\t\t--------")
	
	allPassed := true
	for _, tc := range testCases {
		result := models.MaskUsername(tc.input)
		passed := result == tc.expected
		if !passed {
			allPassed = false
		}
		
		status := "✓"
		if !passed {
			status = "✗"
		}
		
		fmt.Printf("%-15s\t\t%-15s\t\t%-15s\t\t%s\n", 
			tc.input, result, tc.expected, status)
	}
	
	fmt.Println("\n=== 测试总结 ===")
	if allPassed {
		fmt.Println("✓ 所有测试用例通过")
	} else {
		fmt.Println("✗ 部分测试用例失败")
	}
	
	fmt.Println("\n=== 脱敏规则说明 ===")
	fmt.Println("1. 用户名长度 = 1: 不脱敏，直接显示")
	fmt.Println("2. 用户名长度 = 2: 在中间加 *")
	fmt.Println("3. 用户名长度 3-4: 显示首尾字符，中间用 * 替换")
	fmt.Println("4. 用户名长度 ≥ 5: 显示首尾各2个字符，中间用 * 替换")
	
	fmt.Println("\n=== 优化效果 ===")
	fmt.Println("优化前示例:")
	fmt.Println("  '张三' → '张三'（不脱敏）")
	fmt.Println("  '张三丰李四王' → '张****王'")
	fmt.Println("  'test_user_123' → 't*********3'")
	fmt.Println("")
	fmt.Println("优化后示例:")
	fmt.Println("  '张三' → '张*三'（在中间加*）")
	fmt.Println("  '张三丰李四王' → '张三*李四'")
	fmt.Println("  'test_user_123' → 'te*23'")
	fmt.Println("")
	fmt.Println("优化效果:")
	fmt.Println("  - 统一了脱敏规则，所有用户名都进行脱敏")
	fmt.Println("  - 脱敏长度大幅缩短")
	fmt.Println("  - 保留了更多可识别信息")
	fmt.Println("  - 提高了用户体验")
	fmt.Println("  - 仍然保护了用户隐私")
} 