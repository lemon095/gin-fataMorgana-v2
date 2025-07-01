package controllers

import (
	"gin-fataMorgana/models"
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"

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

// GetLeaderboard godoc
// @Summary 获取任务热榜
// @Description 获取任务热榜排行榜列表和当前用户数据
// @Tags 热榜
// @Accept json
// @Produce json
// @Param request body models.LeaderboardRequest true "热榜请求参数"
// @Success 200 {object} models.LeaderboardResponse "成功返回热榜数据"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/leaderboard [post]
func (c *LeaderboardController) GetLeaderboard(ctx *gin.Context) {
	var request models.LeaderboardRequest
	
	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.InvalidParamsWithMessage(ctx, "请求参数错误: "+err.Error())
		return
	}
	
	// 验证用户ID
	if request.UserID <= 0 {
		utils.InvalidParamsWithMessage(ctx, "用户ID无效，用户ID必须大于0")
		return
	}
	
	// 获取热榜数据
	response, err := c.leaderboardService.GetLeaderboard(request.UserID)
	if err != nil {
		utils.InternalError(ctx)
		return
	}
	
	// 返回成功响应
	utils.SuccessWithMessage(ctx, "获取热榜数据成功", response)
} 