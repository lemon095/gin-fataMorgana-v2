package controllers

import (
	"gin-fataMorgana/middleware"
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// CronController 定时任务控制器
type CronController struct {
	cronService *services.CronService
}

// NewCronController 创建定时任务控制器实例
func NewCronController() *CronController {
	return &CronController{
		cronService: nil, // 将在main.go中注入
	}
}

// SetCronService 设置定时任务服务（用于依赖注入）
func (cc *CronController) SetCronService(cronService *services.CronService) {
	cc.cronService = cronService
}

// ManualGenerateOrders 手动生成订单
func (cc *CronController) ManualGenerateOrders(c *gin.Context) {
	// 检查是否有定时任务服务
	if cc.cronService == nil {
		utils.ErrorWithMessage(c, utils.CodeOperationFailed, "定时任务服务未初始化")
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	// 检查用户权限（这里可以添加管理员权限检查）
	// TODO: 添加管理员权限验证

	var req struct {
		Count int `json:"count" binding:"required,min=1,max=1000"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 手动生成订单
	stats, err := cc.cronService.ManualGenerateOrders(req.Count)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeOperationFailed, "手动生成订单失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "手动生成订单成功", gin.H{
		"total_generated":    stats.TotalGenerated,
		"purchase_orders":    stats.PurchaseOrders,
		"group_buy_orders":   stats.GroupBuyOrders,
		"total_amount":       stats.TotalAmount,
		"total_profit":       stats.TotalProfit,
		"last_generation":    stats.LastGeneration,
		"average_time":       stats.AverageTime.String(),
	})
}

// ManualCleanup 手动清理数据
func (cc *CronController) ManualCleanup(c *gin.Context) {
	// 检查是否有定时任务服务
	if cc.cronService == nil {
		utils.ErrorWithMessage(c, utils.CodeOperationFailed, "定时任务服务未初始化")
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	// 检查用户权限（这里可以添加管理员权限检查）
	// TODO: 添加管理员权限验证

	// 手动清理数据
	stats, err := cc.cronService.ManualCleanup()
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeOperationFailed, "手动清理数据失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "手动清理数据成功", gin.H{
		"deleted_orders":     stats.DeletedOrders,
		"deleted_group_buys": stats.DeletedGroupBuys,
		"last_cleanup":       stats.LastCleanup,
		"cleanup_time":       stats.CleanupTime.String(),
	})
}

// GetCronStatus 获取定时任务状态
func (cc *CronController) GetCronStatus(c *gin.Context) {
	// 检查是否有定时任务服务
	if cc.cronService == nil {
		utils.ErrorWithMessage(c, utils.CodeOperationFailed, "定时任务服务未初始化")
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	// 检查用户权限（这里可以添加管理员权限检查）
	// TODO: 添加管理员权限验证

	// 获取定时任务状态
	status := cc.cronService.GetCronStatus()

	utils.Success(c, gin.H{
		"cron_status": status,
	})
} 

// ManualUpdateLeaderboardCache 手动更新热榜缓存
func (cc *CronController) ManualUpdateLeaderboardCache(c *gin.Context) {
	// 检查是否有定时任务服务
	if cc.cronService == nil {
		utils.ErrorWithMessage(c, utils.CodeOperationFailed, "定时任务服务未初始化")
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	// 检查用户权限（这里可以添加管理员权限检查）
	// TODO: 添加管理员权限验证

	// 手动更新热榜缓存
	err := cc.cronService.ManualUpdateLeaderboardCache()
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeOperationFailed, "手动更新热榜缓存失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "手动更新热榜缓存成功", gin.H{
		"update_time": "now",
	})
} 