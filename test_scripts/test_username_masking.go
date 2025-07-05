package main

import (
	"fmt"
	"gin-fataMorgana/models"
)

func main() {
	fmt.Println("=== 用户名脱敏功能测试（统一版本）===")
	
	// 测试用例
	testCases := []struct {
		input    string
		expected string
	}{
		{"a", "a"},                         // 单字符，不脱敏
		{"张三", "张**三"},                   // 长度=2，统一格式
		{"张三丰", "张**丰"},                 // 长度=3，统一格式
		{"张三丰李", "张**李"},               // 长度=4，统一格式
		{"张三丰李四", "张**四"},             // 长度=5，统一格式
		{"张三丰李四王", "张**王"},           // 长度=6，统一格式
		{"ab", "a**b"},                     // 英文，长度=2
		{"test", "t**t"},                   // 英文，长度=4
		{"test_user", "t**r"},              // 英文，长度=9
		{"test_user_123", "t**3"},          // 英文+数字，长度=13
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
	fmt.Println("2. 用户名长度 ≥ 2: 统一格式：首位 + ** + 末位")
	
	fmt.Println("\n=== 统一脱敏效果 ===")
	fmt.Println("优化前示例:")
	fmt.Println("  '张三' → '张*三'（2位特殊处理）")
	fmt.Println("  '张三丰' → '张*丰'（3位特殊处理）")
	fmt.Println("  '张三丰李四王' → '张****王'（复杂处理）")
	fmt.Println("  'test_user_123' → 't*********3'（复杂处理）")
	fmt.Println("")
	fmt.Println("统一后示例:")
	fmt.Println("  '张三' → '张**三'")
	fmt.Println("  '张三丰' → '张**丰'")
	fmt.Println("  '张三丰李四王' → '张**王'")
	fmt.Println("  'test_user_123' → 't**3'")
	fmt.Println("")
	fmt.Println("统一效果:")
	fmt.Println("  - 所有用户名使用相同的脱敏规则")
	fmt.Println("  - 脱敏长度统一为3个字符")
	fmt.Println("  - 代码逻辑更简单，易于维护")
	fmt.Println("  - 用户体验更一致")
	fmt.Println("  - 仍然保护了用户隐私")
} 