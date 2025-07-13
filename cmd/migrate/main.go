package main

import (
	"flag"
	"log"
	"os"

	"gin-fataMorgana/config"
	"gin-fataMorgana/database"
)

func main() {
	// 解析命令行参数
	var (
		help       = flag.Bool("help", false, "显示帮助信息")
		checkIndex = flag.Bool("check-index", false, "检测并创建缺失的索引")
		showIndex  = flag.Bool("show-index", false, "显示当前数据库的所有索引")
	)
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	log.Println("=== 数据库迁移工具 ===")

	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化数据库连接
	if err := database.InitMySQL(); err != nil {
		log.Fatalf("初始化数据库连接失败: %v", err)
	}
	defer database.CloseDB()

	// 根据参数执行不同操作
	if *checkIndex {
		log.Println("🔍 开始检测和创建索引...")
		if err := database.CheckAndCreateIndexes(); err != nil {
			log.Fatalf("索引检测和创建失败: %v", err)
		}
		log.Println("✅ 索引检测和创建完成！")
		return
	}

	if *showIndex {
		log.Println("📋 显示当前数据库索引...")
		if err := database.ShowAllIndexes(); err != nil {
			log.Fatalf("显示索引失败: %v", err)
		}
		return
	}

	// 默认执行完整迁移
	log.Println("🚀 开始数据库迁移...")

	// 执行迁移
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("✅ 数据库迁移完成！")
	log.Println("📋 已创建的表:")
	log.Println("   - users (用户表)")
	log.Println("   - wallets (钱包表)")
	log.Println("   - wallet_transactions (钱包交易流水表)")
	log.Println("   - admin_users (邀请码管理表)")
	log.Println("   - user_login_logs (用户登录日志表)")
	log.Println("   - orders (订单表)")
	log.Println("   - group_buys (拼单表)")
	log.Println("   - amount_config (金额配置表)")
	log.Println("   - announcements (公告表)")
	log.Println("   - announcement_banners (公告图片表)")
	log.Println("   - member_level (用户等级配置表)")
	log.Println("   - lottery_periods (游戏期数表)")
	log.Println("")
	log.Println("🔍 如需检测索引，请运行: go run cmd/migrate/main.go -check-index")
	log.Println("📋 如需查看索引，请运行: go run cmd/migrate/main.go -show-index")
}
