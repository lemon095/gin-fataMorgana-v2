package controllers

import (
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// CurrencyController 货币配置控制器
type CurrencyController struct {
	currencyService *services.CurrencyService
}

// NewCurrencyController 创建货币配置控制器实例
func NewCurrencyController() *CurrencyController {
	return &CurrencyController{
		currencyService: services.NewCurrencyService(),
	}
}

// GetCurrentCurrency 获取当前货币配置
// @Summary 获取当前货币配置
// @Description 获取系统当前使用的货币符号配置
// @Tags 货币配置
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=services.CurrencyConfig} "成功"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/v1/currency/current [post]
func (cc *CurrencyController) GetCurrentCurrency(c *gin.Context) {
	// 获取当前货币配置
	currencyConfig, err := cc.currencyService.GetCurrentCurrency(c.Request.Context())
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeServer, "获取货币配置失败")
		return
	}

	// 返回成功响应
	utils.Success(c, currencyConfig)
} 