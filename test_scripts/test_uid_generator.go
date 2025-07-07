package main

import (
	"fmt"
	"gin-fataMorgana/utils"
	"sync"
)

func main() {
	fmt.Println("=== UID生成器唯一性测试 ===")
	
	// 初始化系统UID生成器
	utils.InitSystemUIDGenerator(1)
	
	// 测试并发生成拼单号
	fmt.Println("1. 测试并发生成拼单号...")
	
	var wg sync.WaitGroup
	groupBuyNos := make([]string, 100)
	
	// 启动100个goroutine并发生成拼单号
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			groupBuyNos[index] = utils.GenerateSystemGroupBuyNo()
		}(i)
	}
	
	wg.Wait()
	
	// 检查是否有重复
	seen := make(map[string]bool)
	duplicates := 0
	
	for i, no := range groupBuyNos {
		if seen[no] {
			fmt.Printf("❌ 发现重复拼单号: %s (索引: %d)\n", no, i)
			duplicates++
		} else {
			seen[no] = true
		}
	}
	
	if duplicates == 0 {
		fmt.Printf("✅ 成功生成 %d 个唯一拼单号\n", len(groupBuyNos))
	} else {
		fmt.Printf("❌ 发现 %d 个重复拼单号\n", duplicates)
	}
	
	// 测试订单号生成
	fmt.Println("\n2. 测试并发生成订单号...")
	
	orderNos := make([]string, 100)
	
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			orderNos[index] = utils.GenerateSystemOrderNo()
		}(i)
	}
	
	wg.Wait()
	
	// 检查订单号是否有重复
	seen = make(map[string]bool)
	duplicates = 0
	
	for i, no := range orderNos {
		if seen[no] {
			fmt.Printf("❌ 发现重复订单号: %s (索引: %d)\n", no, i)
			duplicates++
		} else {
			seen[no] = true
		}
	}
	
	if duplicates == 0 {
		fmt.Printf("✅ 成功生成 %d 个唯一订单号\n", len(orderNos))
	} else {
		fmt.Printf("❌ 发现 %d 个重复订单号\n", duplicates)
	}
	
	// 显示一些生成的号码示例
	fmt.Println("\n3. 生成的号码示例:")
	fmt.Println("拼单号示例:")
	for i := 0; i < 5; i++ {
		fmt.Printf("  %s\n", groupBuyNos[i])
	}
	
	fmt.Println("订单号示例:")
	for i := 0; i < 5; i++ {
		fmt.Printf("  %s\n", orderNos[i])
	}
	
	fmt.Println("\n✅ 测试完成")
} 