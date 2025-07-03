package controllers

import (
	"fmt"
	"gin-fataMorgana/models"
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// AmountConfigController 金额配置控制器
type AmountConfigController struct {
	amountConfigService *services.AmountConfigService
}

// NewAmountConfigController 创建金额配置控制器实例
func NewAmountConfigController() *AmountConfigController {
	return &AmountConfigController{
		amountConfigService: services.NewAmountConfigService(),
	}
}

// GetAmountConfigsByType godoc
// @Summary 根据类型获取金额配置列表
// @Description 根据配置类型获取金额配置列表，支持recharge(充值)和withdraw(提现)类型，只返回激活状态(is_active=true)的配置
// @Tags 金额配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AmountConfigRequest true "金额配置请求参数"
// @Success 200 {object} utils.Response{data=[]models.AmountConfigResponse} "成功返回金额配置列表"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "认证失败"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/amount-config/list [post]
func (c *AmountConfigController) GetAmountConfigsByType(ctx *gin.Context) {
	var request models.AmountConfigRequest

	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.InvalidParamsWithMessage(ctx, "请求参数错误: "+err.Error())
		return
	}

	// 获取金额配置列表
	configs, err := c.amountConfigService.GetAmountConfigsByType(ctx, request.Type)
	if err != nil {
		utils.InternalError(ctx)
		return
	}

	// 返回成功响应
	utils.SuccessWithMessage(ctx, "获取金额配置列表成功", configs)
}

// GetAmountConfigByID godoc
// @Summary 根据ID获取金额配置详情
// @Description 根据配置ID获取金额配置详细信息，只返回激活状态(is_active=true)的配置
// @Tags 金额配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "配置ID"
// @Success 200 {object} utils.Response{data=models.AmountConfigResponse} "成功返回金额配置详情"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "认证失败"
// @Failure 404 {object} utils.Response "配置不存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/amount-config/{id} [get]
func (c *AmountConfigController) GetAmountConfigByID(ctx *gin.Context) {
	// 获取路径参数
	idStr := ctx.Param("id")
	if idStr == "" {
		utils.InvalidParamsWithMessage(ctx, "配置ID不能为空")
		return
	}

	// 解析ID
	var id int64
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		utils.InvalidParamsWithMessage(ctx, "配置ID格式错误")
		return
	}

	// 获取金额配置详情
	config, err := c.amountConfigService.GetAmountConfigByID(ctx, id)
	if err != nil {
		utils.NotFound(ctx)
		return
	}

	// 返回成功响应
	utils.SuccessWithMessage(ctx, "获取金额配置详情成功", config)
}
