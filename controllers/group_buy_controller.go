package controllers

import (
	"context"

	"gin-fataMorgana/database"
	"gin-fataMorgana/middleware"
	"gin-fataMorgana/models"
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// GroupBuyController 拼单控制器
type GroupBuyController struct {
	groupBuyService *services.GroupBuyService
}

// NewGroupBuyController 创建拼单控制器实例
func NewGroupBuyController() *GroupBuyController {
	return &GroupBuyController{
		groupBuyService: services.NewGroupBuyService(),
	}
}

// GetActiveGroupBuyDetail 获取活跃拼单详情
// @Summary 获取活跃拼单详情
// @Description 根据用户ID获取未截止的拼单详情
// @Tags 拼单
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=models.GetGroupBuyDetailResponse}
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/group-buy/active-detail [post]
func (c *GroupBuyController) GetActiveGroupBuyDetail(ctx *gin.Context) {
	// 获取当前用户ID
	userID := middleware.GetCurrentUser(ctx)
	if userID == 0 {
		utils.Unauthorized(ctx)
		return
	}

	// 根据user_id查询uid
	userRepo := database.NewUserRepository()
	var user models.User
	err := userRepo.FindByID(context.Background(), userID, &user)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	// 调用服务层
	response, err := c.groupBuyService.GetActiveGroupBuyDetail(ctx, user.Uid)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeDatabaseError, err.Error())
		return
	}

	// 返回成功响应
	utils.Success(ctx, response)
}

// JoinGroupBuy 确认参与拼单
// @Summary 确认参与拼单
// @Description 确认参与拼单，创建订单并更新拼单状态
// @Tags 拼单
// @Accept json
// @Produce json
// @Param request body models.JoinGroupBuyRequest true "确认参与拼单请求"
// @Success 200 {object} utils.Response{data=models.JoinGroupBuyResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/group-buy/join [post]
func (c *GroupBuyController) JoinGroupBuy(ctx *gin.Context) {
	// 获取当前用户ID
	userID := middleware.GetCurrentUser(ctx)
	if userID == 0 {
		utils.Unauthorized(ctx)
		return
	}

	// 根据user_id查询uid
	userRepo := database.NewUserRepository()
	var user models.User
	err := userRepo.FindByID(context.Background(), userID, &user)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	// 绑定请求参数
	var req models.JoinGroupBuyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(ctx, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层
	response, err := c.groupBuyService.JoinGroupBuy(ctx, req.GroupBuyNo, user.Uid)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeDatabaseError, err.Error())
		return
	}

	// 返回成功响应
	utils.Success(ctx, response)
}
