package main

import (
	"fmt"
	"sync"
	"time"

	"gin-fataMorgana/utils"
)

func main() {
	// 初始化系统UID生成器
	utils.InitSystemUIDGenerator(1)

	fmt.Println("🧪 开始测试订单号唯一性...")
	fmt.Println("📊 测试参数: 1000个并发订单号生成")

	// 使用map来检查重复
	orderNos := make(map[string]bool)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// 并发生成1000个订单号
	concurrency := 1000
	wg.Add(concurrency)

	startTime := time.Now()

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			
			// 生成订单号
			orderNo := utils.GenerateSystemOrderNo()
			
			// 检查重复
			mutex.Lock()
			if orderNos[orderNo] {
				fmt.Printf("❌ 发现重复订单号: %s (协程ID: %d)\n", orderNo, id)
			} else {
				orderNos[orderNo] = true
			}
			mutex.Unlock()
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	// 统计结果
	fmt.Printf("\n📈 测试结果:\n")
	fmt.Printf("   生成订单号数量: %d\n", len(orderNos))
	fmt.Printf("   期望数量: %d\n", concurrency)
	fmt.Printf("   耗时: %v\n", duration)
	
	if len(orderNos) == concurrency {
		fmt.Println("✅ 所有订单号都是唯一的！")
	} else {
		fmt.Printf("❌ 发现重复订单号！实际生成: %d, 期望: %d\n", len(orderNos), concurrency)
	}

	// 测试拼单号唯一性
	fmt.Println("\n🧪 开始测试拼单号唯一性...")
	
	groupBuyNos := make(map[string]bool)
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			
			// 生成拼单号
			groupBuyNo := utils.GenerateSystemGroupBuyNo()
			
			// 检查重复
			mutex.Lock()
			if groupBuyNos[groupBuyNo] {
				fmt.Printf("❌ 发现重复拼单号: %s (协程ID: %d)\n", groupBuyNo, id)
			} else {
				groupBuyNos[groupBuyNo] = true
			}
			mutex.Unlock()
		}(i)
	}

	wg.Wait()

	fmt.Printf("\n📈 拼单号测试结果:\n")
	fmt.Printf("   生成拼单号数量: %d\n", len(groupBuyNos))
	fmt.Printf("   期望数量: %d\n", concurrency)
	
	if len(groupBuyNos) == concurrency {
		fmt.Println("✅ 所有拼单号都是唯一的！")
	} else {
		fmt.Printf("❌ 发现重复拼单号！实际生成: %d, 期望: %d\n", len(groupBuyNos), concurrency)
	}

	// 显示一些示例订单号
	fmt.Println("\n📝 示例订单号:")
	count := 0
	for orderNo := range orderNos {
		if count < 5 {
			fmt.Printf("   %s\n", orderNo)
			count++
		} else {
			break
		}
	}

	fmt.Println("\n📝 示例拼单号:")
	count = 0
	for groupBuyNo := range groupBuyNos {
		if count < 5 {
			fmt.Printf("   %s\n", groupBuyNo)
			count++
		} else {
			break
		}
	}
} 