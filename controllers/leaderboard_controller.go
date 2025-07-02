package controllers

import (
	"gin-fataMorgana/middleware"
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LeaderboardController 热榜控制器
type LeaderboardController struct {
	leaderboardService *services.LeaderboardService
}

// NewLeaderboardController 创建热榜控制器实例
func NewLeaderboardController() *LeaderboardController {
	return &LeaderboardController{
		leaderboardService: services.NewLeaderboardService(),
	}
}

// GetLeaderboard 获取任务热榜
func (c *LeaderboardController) GetLeaderboard(ctx *gin.Context) {
	// 获取当前用户ID
	userID := middleware.GetCurrentUser(ctx)
	if userID == 0 {
		utils.Unauthorized(ctx)
		return
	}

	uid := strconv.FormatUint(uint64(userID), 10)
	
	// 获取热榜数据
	response, err := c.leaderboardService.GetLeaderboard(uid)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeDatabaseError, err.Error())
		return
	}
	
	// 返回成功响应
	utils.Success(ctx, response)
} 