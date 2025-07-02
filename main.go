// @title           Gin-FataMorgana API
// @version         1.0
// @description     Gin-FataMorgana 是一个基于Gin框架的Go Web服务，提供用户认证、钱包管理、健康监控等功能
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9001
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 请输入 "Bearer " 加上JWT token，例如: "Bearer abcde12345"

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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// 启动幂等键清理器
	ctx := context.Background()
	go utils.StartIdempotencyCleaner(ctx)

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

	// 配置CORS中间件
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",           // React开发服务器
		"http://localhost:8080",           // Vue开发服务器
		"http://localhost:5173",           // Vite开发服务器
		"https://colombiatkadmin.com",     // 生产环境前端
		"http://colombiatkadmin.com",      // 生产环境前端（HTTP）
		"https://www.colombiatkadmin.com", // 生产环境前端（带www）
		"http://www.colombiatkadmin.com",  // 生产环境前端（带www，HTTP）
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	corsConfig.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
		"X-CSRF-Token",
		"X-API-Key",
		"Cache-Control",
		"Pragma",
	}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour

	// 添加中间件
	r.Use(cors.New(corsConfig))           // CORS中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.SessionMiddleware()) // 全局会话管理中间件

	// 创建控制器实例
	authController := controllers.NewAuthController()
	sessionController := controllers.NewSessionController()
	healthController := controllers.NewHealthController()
	walletController := controllers.NewWalletController()
	orderController := controllers.NewOrderController()
	leaderboardController := controllers.NewLeaderboardController()
	amountConfigController := controllers.NewAmountConfigController()
	announcementController := controllers.NewAnnouncementController()

	// 定义路由
	r.GET("/", func(c *gin.Context) {
		utils.Success(c, gin.H{
			"message": "欢迎使用 Gin-FataMorgana 服务!",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查（保持原有路径，方便监控）
	r.GET("/health", healthController.HealthCheck)

	// API v1 路由组 - 统一前缀
	api := r.Group("/api/v1")
	{
		// 系统健康检查路由组（保持GET，便于监控）
		health := api.Group("/health")
		{
			health.GET("/check", healthController.HealthCheck)       // 系统健康检查
			health.GET("/database", healthController.DatabaseHealth) // 数据库健康检查
			health.GET("/redis", healthController.RedisHealth)       // Redis健康检查
		}

		// 认证相关路由组
		auth := api.Group("/auth")
		{
			auth.POST("/register", middleware.RegisterOpenMiddleware(), middleware.RegisterRateLimitMiddleware(), authController.Register)           // 用户注册
			auth.POST("/login", middleware.LoginRateLimitMiddleware(), authController.Login)                 // 用户登录
			auth.POST("/refresh", authController.RefreshToken)        // 刷新令牌
			auth.POST("/logout", authController.Logout)               // 用户登出
			auth.POST("/profile", middleware.AuthMiddleware(), authController.GetProfile) // 获取用户信息
			auth.POST("/bind-bank-card", middleware.AuthMiddleware(), authController.BindBankCard) // 绑定银行卡
			auth.POST("/bank-card", middleware.AuthMiddleware(), authController.GetBankCardInfo) // 获取银行卡信息
		}

		// 会话管理路由组
		session := api.Group("/session")
		{
			session.POST("/status", sessionController.CheckLoginStatus) // 检查登录状态
			session.POST("/user", sessionController.GetCurrentUserInfo) // 获取当前用户信息
			session.POST("/logout", sessionController.Logout)          // 用户登出
			session.POST("/refresh", sessionController.RefreshSession) // 刷新会话
		}

		// 钱包相关路由组
		wallet := api.Group("/wallet")
		{
			wallet.Use(middleware.AuthMiddleware()) // 需要认证
			wallet.POST("/info", walletController.GetWallet)                    // 获取钱包信息
			wallet.POST("/transactions", walletController.GetUserTransactions)      // 获取资金记录
			wallet.POST("/transaction-detail", walletController.GetTransactionDetail) // 获取交易详情
			wallet.POST("/withdraw", middleware.WithdrawRateLimitMiddleware(), walletController.RequestWithdraw)             // 申请提现
			wallet.POST("/withdraw-summary", walletController.GetWithdrawSummary)   // 获取提现汇总
			wallet.POST("/recharge", walletController.Recharge)                   // 充值申请
		}

		// 订单相关路由组
		order := api.Group("/order")
		{
			order.Use(middleware.AuthMiddleware()) // 需要认证
			order.POST("/create", orderController.CreateOrder)                    // 创建订单
			order.POST("/list", orderController.GetOrderList)                     // 获取订单列表
			order.POST("/detail", orderController.GetOrderDetail)                 // 获取订单详情
			order.POST("/stats", orderController.GetOrderStats)                   // 获取订单统计
		}

		// 管理员路由组
		admin := api.Group("/admin")
		{
			admin.Use(middleware.AuthMiddleware()) // 需要认证
			// 提现确认和取消接口已移除
		}

		// 假数据接口路由组
		fake := api.Group("/fake")
		{
			fake.POST("/activities", controllers.GetFakeRealtimeActivities) // 获取假数据实时动态
		}

		// 热榜相关路由组
		leaderboard := api.Group("/leaderboard")
		{
			leaderboard.POST("/ranking", leaderboardController.GetLeaderboard) // 获取任务热榜
		}

		// 金额配置相关路由组
		amountConfig := api.Group("/amount-config")
		{
			amountConfig.Use(middleware.AuthMiddleware()) // 需要认证
			amountConfig.POST("/list", amountConfigController.GetAmountConfigsByType) // 根据类型获取金额配置列表
			amountConfig.GET("/:id", amountConfigController.GetAmountConfigByID)      // 根据ID获取金额配置详情
		}

		// 公告相关路由组
		announcements := api.Group("/announcements")
		{
			announcements.POST("/list", announcementController.GetAnnouncementList) // 获取公告列表
		}
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
