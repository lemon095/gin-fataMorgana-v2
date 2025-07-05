package main

import (
	"fmt"
	"log"
	"os"

	"gin-fataMorgana/config"
)

func main() {
	fmt.Println("=== 测试配置加载 ===")
	
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}
	
	fmt.Println("✅ 配置加载成功")
	
	// 检查假数据配置
	fakeData := config.GlobalConfig.FakeData
	fmt.Printf("📋 假数据配置:\n")
	fmt.Printf("  启用: %v\n", fakeData.Enabled)
	fmt.Printf("  表达式: %s\n", fakeData.CronExpression)
	fmt.Printf("  清理表达式: %s\n", fakeData.CleanupCron)
	fmt.Printf("  最小订单: %d\n", fakeData.MinOrders)
	fmt.Printf("  最大订单: %d\n", fakeData.MaxOrders)
	fmt.Printf("  购买单比例: %.2f\n", fakeData.PurchaseRatio)
	fmt.Printf("  任务数最小值: %d\n", fakeData.TaskMinCount)
	fmt.Printf("  任务数最大值: %d\n", fakeData.TaskMaxCount)
	fmt.Printf("  保留天数: %d\n", fakeData.RetentionDays)
	
	// 检查配置文件路径
	configFile := "config.yaml"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}
	
	fmt.Printf("📁 使用的配置文件: %s\n", configFile)
	
	if _, err := os.Stat(configFile); err == nil {
		fmt.Println("✅ 配置文件存在")
	} else {
		fmt.Printf("❌ 配置文件不存在: %v\n", err)
	}
	
	fmt.Println("=== 配置测试完成 ===")
} 