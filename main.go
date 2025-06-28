package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-fataMorgana/config"
	"gin-fataMorgana/controllers"
	"gin-fataMorgana/database"
	"gin-fataMorgana/middleware"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Printf("加载配置失败: %v", err)
		os.Exit(1)
	}

	// 验证配置
	if err := config.ValidateConfig(); err != nil {
		log.Printf("配置验证失败: %v", err)
		os.Exit(1)
	}

	// 初始化JWT配置
	utils.InitJWT()

	// 初始化雪花算法
	utils.InitSnowflake()

	// 初始化MySQL数据库
	if err := database.InitMySQL(); err != nil {
		log.Printf("初始化MySQL失败: %v", err)
		os.Exit(1)
	}

	// 自动迁移数据库表
	if err := database.AutoMigrate(); err != nil {
		log.Printf("数据库迁移失败: %v", err)
		os.Exit(1)
	}

	// 初始化Redis
	if err := database.InitRedis(); err != nil {
		log.Printf("初始化Redis失败: %v", err)
		os.Exit(1)
	}

	// 注册自定义验证器
	utils.RegisterCustomValidators()

	// 设置Gin模式
	gin.SetMode(config.GlobalConfig.Server.Mode)

	// 创建默认的gin引擎
	r := gin.Default()

	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.SessionMiddleware()) // 全局会话管理中间件

	// 创建控制器实例
	authController := controllers.NewAuthController()
	sessionController := controllers.NewSessionController()
	healthController := controllers.NewHealthController()
	walletController := controllers.NewWalletController()

	// 定义路由
	r.GET("/", func(c *gin.Context) {
		utils.Success(c, gin.H{
			"message": "欢迎使用 Gin-FataMorgana 服务!",
			"status":  "running",
			"version": "1.0.0",
		})
	}) // 首页

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		utils.Success(c, gin.H{
			"status":  "healthy",
			"service": "gin-fataMorgana",
		})
	}) // 健康检查

	// 系统健康检查路由组
	health := r.Group("/health")
	{
		health.GET("/check", healthController.HealthCheck)           // 系统健康检查
		health.GET("/database", healthController.DatabaseHealth)     // 数据库健康检查
		health.GET("/redis", healthController.RedisHealth)           // Redis健康检查
		health.GET("/stats", healthController.DatabaseStats)         // 数据库统计信息
		health.GET("/query-stats", healthController.QueryStats)      // 查询统计信息
		health.GET("/optimization", healthController.PerformanceOptimization) // 性能优化建议
	}

	// 认证相关路由组
	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)           // 用户注册
		auth.POST("/login", authController.Login)                 // 用户登录
		auth.POST("/refresh", authController.RefreshToken)        // 刷新令牌
		auth.POST("/logout", authController.Logout)               // 用户登出
		auth.GET("/profile", middleware.AuthMiddleware(), authController.GetProfile) // 获取用户信息
		auth.POST("/bind-bank-card", middleware.AuthMiddleware(), authController.BindBankCard) // 绑定银行卡
		auth.GET("/bank-card", middleware.AuthMiddleware(), authController.GetBankCardInfo) // 获取银行卡信息
	}

	// 会话管理路由组
	session := r.Group("/session")
	{
		session.GET("/status", sessionController.CheckLoginStatus) // 检查登录状态
		session.GET("/user", sessionController.GetCurrentUserInfo) // 获取当前用户信息
		session.POST("/logout", sessionController.Logout)          // 用户登出
		session.POST("/refresh", sessionController.RefreshSession) // 刷新会话
	}

	// 钱包相关路由组
	wallet := r.Group("/wallet")
	{
		wallet.Use(middleware.AuthMiddleware()) // 需要认证
		wallet.GET("/info", walletController.GetWallet)                    // 获取钱包信息
		wallet.GET("/transactions", walletController.GetUserTransactions)      // 获取资金记录
		wallet.POST("/withdraw", walletController.RequestWithdraw)             // 申请提现
		wallet.GET("/withdraw-summary", walletController.GetWithdrawSummary)   // 获取提现汇总
	}

	// 管理员路由组
	admin := r.Group("/admin")
	{
		admin.Use(middleware.AuthMiddleware()) // 需要认证
		admin.POST("/withdraw/confirm", walletController.ConfirmWithdraw) // 确认提现
		admin.POST("/withdraw/cancel", walletController.CancelWithdraw)   // 取消提现
	}

	// 创建HTTP服务器
	addr := fmt.Sprintf("%s:%d", config.GlobalConfig.Server.Host, config.GlobalConfig.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 启动服务器
	go func() {
		log.Printf("服务器启动在 %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("服务器启动失败: %v", err)
			os.Exit(1)
		}
	}()

	// 设置优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 设置关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("服务器关闭失败: %v", err)
	}

	// 关闭数据库连接
	if err := database.CloseDB(); err != nil {
		log.Printf("关闭数据库连接失败: %v", err)
	}

	// 关闭Redis连接
	if err := database.CloseRedis(); err != nil {
		log.Printf("关闭Redis连接失败: %v", err)
	}

	log.Println("服务器已安全关闭")
}
