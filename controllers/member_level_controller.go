package controllers

import (
	"strconv"

	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// MemberLevelController 用户等级配置表控制器
type MemberLevelController struct {
	memberLevelService *services.MemberLevelService
}

// NewMemberLevelController 创建用户等级配置表控制器实例
func NewMemberLevelController() *MemberLevelController {
	return &MemberLevelController{
		memberLevelService: services.NewMemberLevelService(services.NewMemberLevelRepository()),
	}
}

// GetAllLevels 获取所有等级配置
// @Summary 获取所有等级配置
// @Description 获取所有启用的用户等级配置信息
// @Tags 用户等级
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.MemberLevel}
// @Failure 500 {object} utils.Response
// @Router /api/member-levels [get]
func (c *MemberLevelController) GetAllLevels(ctx *gin.Context) {
	levels, err := c.memberLevelService.GetAllLevels(ctx)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeOperationFailed, "获取等级配置失败")
		return
	}

	utils.SuccessWithMessage(ctx, "获取等级配置成功", levels)
}

// GetLevelByLevel 根据等级获取配置
// @Summary 根据等级获取配置
// @Description 根据指定等级获取等级配置信息
// @Tags 用户等级
// @Accept json
// @Produce json
// @Param level path int true "等级"
// @Success 200 {object} utils.Response{data=models.MemberLevel}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/member-levels/{level} [get]
func (c *MemberLevelController) GetLevelByLevel(ctx *gin.Context) {
	levelStr := ctx.Param("level")
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		utils.InvalidParamsWithMessage(ctx, "等级参数无效")
		return
	}

	memberLevel, err := c.memberLevelService.GetLevelByLevel(ctx, level)
	if err != nil {
		utils.NotFound(ctx)
		return
	}

	utils.SuccessWithMessage(ctx, "获取等级配置成功", memberLevel)
}

// CalculateCashback 计算返现金额
// @Summary 计算返现金额
// @Description 根据用户经验值和金额计算返现金额
// @Tags 用户等级
// @Accept json
// @Produce json
// @Param experience query int true "经验值"
// @Param amount query float64 true "金额"
// @Success 200 {object} utils.Response{data=map[string]interface{}}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/member-levels/calculate-cashback [get]
func (c *MemberLevelController) CalculateCashback(ctx *gin.Context) {
	experienceStr := ctx.Query("experience")
	amountStr := ctx.Query("amount")

	experience, err := strconv.Atoi(experienceStr)
	if err != nil {
		utils.InvalidParamsWithMessage(ctx, "经验值参数无效")
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		utils.InvalidParamsWithMessage(ctx, "金额参数无效")
		return
	}

	cashbackAmount, err := c.memberLevelService.CalculateCashback(ctx, experience, amount)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeOperationFailed, "计算返现金额失败")
		return
	}

	result := map[string]interface{}{
		"experience":      experience,
		"amount":          amount,
		"cashback_amount": cashbackAmount,
	}

	utils.SuccessWithMessage(ctx, "计算返现金额成功", result)
}
