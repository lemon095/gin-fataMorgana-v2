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

	"runtime/debug"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 设置时区为 UTC
	os.Setenv("TZ", "UTC")
	time.LoadLocation("UTC")

	// 全局panic恢复机制
	defer func() {
		if r := recover(); r != nil {
			log.Printf("程序发生严重错误，正在恢复: %v", r)
			// 记录堆栈信息
			debug.PrintStack()
			// 优雅关闭
			os.Exit(1)
		}
	}()

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

	// 初始化JWT
	utils.InitJWT()

	// 初始化雪花算法
	utils.InitSnowflake()

	// 启动幂等性清理器
	ctx := context.Background()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("幂等性清理器发生panic: %v", r)
			}
		}()
		utils.StartIdempotencyCleaner(ctx)
	}()

	// 初始化MySQL
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

	// 创建Gin引擎
	r := gin.Default()

	// 添加全局recover中间件
	r.Use(gin.Recovery())

	// 添加自定义recover中间件
	r.Use(func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("HTTP请求发生panic: %v, URL: %s, Method: %s", r, c.Request.URL.Path, c.Request.Method)
				debug.PrintStack()
				
				// 返回500错误
				c.JSON(500, gin.H{
					"code":    500,
					"message": "服务器内部错误",
					"data":    nil,
				})
				c.Abort()
			}
		}()
		c.Next()
	})

	// CORS配置
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"http://localhost:8080",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:8080",
		"https://yourdomain.com",
	}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	corsConfig.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"X-Requested-With",
		"Accept",
		"X-CSRF-Token",
		"X-API-Key",
	}

	// 根据环境设置CORS
	if config.GlobalConfig.Server.Mode == "debug" {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.MaxAge = 12 * time.Hour
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
		corsConfig.AllowHeaders = []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"X-CSRF-Token",
			"X-API-Key",
		}
	}

	// 应用中间件
	r.Use(middleware.CORSMiddleware()) // 自定义CORS中间件
	r.Use(gin.Logger())
	r.Use(middleware.SessionMiddleware()) // 全局会话管理中间件

	// 初始化控制器
	authController := controllers.NewAuthController()
	sessionController := controllers.NewSessionController()
	healthController := controllers.NewHealthController()
	walletController := controllers.NewWalletController()
	orderController := controllers.NewOrderController()
	leaderboardController := controllers.NewLeaderboardController()
	amountConfigController := controllers.NewAmountConfigController()
	announcementController := controllers.NewAnnouncementController()
	groupBuyController := controllers.NewGroupBuyController()
	shareController := controllers.NewShareController()

	// 根路径
	r.GET("/", func(c *gin.Context) {
		utils.Success(c, gin.H{
			"message": "Gin-FataMorgana API Server",
			"version": "1.0.0",
			"time":    time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查
	r.GET("/health", healthController.HealthCheck)

	// API路由组
	api := r.Group("/api/v1")

	// 健康检查路由
	health := api.Group("/health")
	{
		health.GET("/check", healthController.HealthCheck)       // 系统整体健康检查
		health.GET("/database", healthController.DatabaseHealth) // 数据库连接健康检查
		health.GET("/redis", healthController.RedisHealth)       // Redis连接健康检查
	}

	// 认证相关路由
	auth := api.Group("/auth")
	{
		// 注册路由支持OPTIONS
		auth.OPTIONS("/register", func(c *gin.Context) {
			c.Status(200)
		})

		auth.POST("/register", middleware.RegisterOpenMiddleware(), middleware.RegisterRateLimitMiddleware(), authController.Register) // 用户注册 - 创建新用户账户
		auth.POST("/login", middleware.LoginRateLimitMiddleware(), authController.Login)                                               // 用户登录 - 验证用户身份并生成令牌
		auth.POST("/refresh", authController.RefreshToken)                                                                             // 刷新令牌 - 延长用户登录状态
		auth.POST("/logout", authController.Logout)                                                                                    // 用户登出 - 清除用户登录状态
		auth.POST("/profile", middleware.AuthMiddleware(), authController.GetProfile)                                                  // 获取用户信息 - 获取当前用户详细资料
		auth.POST("/change-password", middleware.AuthMiddleware(), authController.ChangePassword)                                      // 修改密码 - 用户修改登录密码
		auth.POST("/bind-bank-card", middleware.AuthMiddleware(), authController.BindBankCard)                                         // 绑定银行卡 - 用户绑定提现银行卡
		auth.POST("/bank-card", middleware.AuthMiddleware(), authController.GetBankCardInfo)                                           // 获取银行卡信息 - 查询用户绑定的银行卡
	}

	// 会话管理路由
	session := api.Group("/session")
	{
		session.POST("/status", sessionController.CheckLoginStatus) // 检查登录状态 - 验证用户是否已登录
		session.POST("/user", sessionController.GetCurrentUserInfo) // 获取当前用户信息 - 获取会话中的用户信息
		session.POST("/logout", sessionController.Logout)           // 用户登出 - 清除用户会话
		session.POST("/refresh", sessionController.RefreshSession)  // 刷新会话 - 延长会话有效期
	}

	// 钱包相关路由
	wallet := api.Group("/wallet")
	{
		wallet.Use(middleware.AuthMiddleware())                                                             // 需要认证
		wallet.POST("/info", walletController.GetWallet)                                                    // 获取钱包信息 - 查询用户余额和钱包状态
		wallet.POST("/transactions", walletController.GetUserTransactions)                                  // 获取资金记录 - 查询用户交易流水历史
		wallet.POST("/transaction-detail", walletController.GetTransactionDetail)                           // 获取交易详情 - 根据流水号查询具体交易信息
		wallet.POST("/withdraw", middleware.GeneralRateLimitMiddleware(), walletController.RequestWithdraw) // 申请提现 - 用户申请从钱包提现到银行卡
		wallet.POST("/withdraw-summary", walletController.GetWithdrawSummary)                               // 获取提现汇总 - 查询用户提现统计信息
		wallet.POST("/recharge", walletController.Recharge)                                                 // 充值申请 - 用户申请从银行卡充值到钱包
	}

	// 订单相关路由
	order := api.Group("/order")
	{
		order.Use(middleware.AuthMiddleware())                // 需要认证
		order.POST("/create", orderController.CreateOrder)    // 创建订单 - 用户创建新任务订单
		order.POST("/list", orderController.GetOrderList)     // 获取订单列表 - 查询用户订单历史（支持状态筛选）
		order.POST("/detail", orderController.GetOrderDetail) // 获取订单详情 - 查询具体订单的详细信息
		order.POST("/stats", orderController.GetOrderStats)   // 获取订单统计 - 查询用户订单统计数据
		order.POST("/period", orderController.GetPeriodList)  // 获取期数列表 - 获取当前活跃期数和价格配置
	}

	// 管理员路由
	admin := api.Group("/admin")
	{
		admin.Use(middleware.AuthMiddleware()) // 需要认证
		// 这里可以添加管理员相关的路由
	}

	// 假数据路由
	fake := api.Group("/fake")
	{
		fake.POST("/activities", controllers.GetFakeRealtimeActivities) // 获取假数据实时动态 - 生成模拟活动数据用于前端测试
	}

	// 排行榜路由
	leaderboard := api.Group("/leaderboard")
	{
		leaderboard.POST("/ranking", leaderboardController.GetLeaderboard) // 获取任务热榜 - 查询周度任务完成排行榜
	}

	// 金额配置路由
	amountConfig := api.Group("/amountConfig")
	{
		amountConfig.Use(middleware.AuthMiddleware())                             // 需要认证
		amountConfig.POST("/list", amountConfigController.GetAmountConfigsByType) // 获取金额配置列表 - 根据类型获取金额配置
		amountConfig.GET("/:id", amountConfigController.GetAmountConfigByID)      // 获取金额配置详情 - 根据ID获取具体配置
	}

	// 公告路由
	announcements := api.Group("/announcements")
	{
		announcements.POST("/list", announcementController.GetAnnouncementList) // 获取公告列表 - 分页获取公告信息
	}

	// 拼单路由
	groupBuy := api.Group("/groupBuy")
	{
		groupBuy.Use(middleware.AuthMiddleware())                                    // 需要认证
		groupBuy.POST("/active-detail", groupBuyController.GetActiveGroupBuyDetail) // 获取活跃拼单详情 - 获取当前可参与的拼单信息
		groupBuy.POST("/join", groupBuyController.JoinGroupBuy)                     // 参与拼单 - 用户参与拼单活动
	}

	// 分享链接接口 - 获取分享链接
	api.POST("/shareLink", shareController.GetShareLink)

	// 启动服务器
	port := fmt.Sprintf("%d", config.GlobalConfig.Server.Port)
	if port == "0" {
		port = "9001"
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// 优雅关闭服务器
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("HTTP服务器goroutine发生panic: %v", r)
			}
		}()
		
		log.Printf("服务器启动在 0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 优雅关闭
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("服务器强制关闭: %v", err)
	}

	log.Println("服务器已关闭")
}
