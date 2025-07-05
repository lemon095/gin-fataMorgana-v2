#!/bin/bash

# 测试拼单号唯一性
echo "=== 测试拼单号唯一性 ==="

# 设置基础URL
BASE_URL="http://localhost:9001/api/v1"

echo "1. 测试拼单号生成器"
echo "生成10个拼单号进行唯一性测试..."

# 创建一个Go测试程序
cat > /tmp/test_group_buy_no.go << 'EOF'
package main

import (
	"fmt"
	"gin-fataMorgana/utils"
)

func main() {
	// 初始化系统UID生成器
	utils.InitSystemUIDGenerator(1)
	
	// 生成10个拼单号
	groupBuyNos := make(map[string]bool)
	
	fmt.Println("生成的拼单号:")
	for i := 0; i < 10; i++ {
		groupBuyNo := utils.GenerateSystemGroupBuyNo()
		fmt.Printf("  %d: %s\n", i+1, groupBuyNo)
		
		if groupBuyNos[groupBuyNo] {
			fmt.Printf("❌ 重复的拼单号: %s\n", groupBuyNo)
		} else {
			groupBuyNos[groupBuyNo] = true
		}
	}
	
	fmt.Printf("\n✅ 生成了 %d 个唯一的拼单号\n", len(groupBuyNos))
}
EOF

# 运行测试
cd /tmp && go run test_group_buy_no.go

echo ""
echo "2. 测试高并发拼单号生成"
echo "模拟高并发场景，快速生成100个拼单号..."

cat > /tmp/test_concurrent_group_buy_no.go << 'EOF'
package main

import (
	"fmt"
	"gin-fataMorgana/utils"
	"sync"
)

func main() {
	// 初始化系统UID生成器
	utils.InitSystemUIDGenerator(1)
	
	// 并发生成拼单号
	var wg sync.WaitGroup
	groupBuyNos := make(map[string]bool)
	var mutex sync.Mutex
	
	// 启动10个goroutine，每个生成10个拼单号
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				groupBuyNo := utils.GenerateSystemGroupBuyNo()
				
				mutex.Lock()
				if groupBuyNos[groupBuyNo] {
					fmt.Printf("❌ 线程%d发现重复拼单号: %s\n", threadID, groupBuyNo)
				} else {
					groupBuyNos[groupBuyNo] = true
				}
				mutex.Unlock()
			}
		}(i)
	}
	
	wg.Wait()
	
	fmt.Printf("✅ 并发测试完成，生成了 %d 个唯一的拼单号\n", len(groupBuyNos))
}
EOF

# 运行并发测试
cd /tmp && go run test_concurrent_group_buy_no.go

echo ""
echo "3. 清理临时文件"
rm -f /tmp/test_group_buy_no.go /tmp/test_concurrent_group_buy_no.go

echo ""
echo "=== 测试完成 ===" 