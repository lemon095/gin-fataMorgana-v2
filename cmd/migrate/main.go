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
		help = flag.Bool("help", false, "显示帮助信息")
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
}
