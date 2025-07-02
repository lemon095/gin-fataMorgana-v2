package controllers

import (
	"gin-fataMorgana/models"
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// AnnouncementController 公告控制器
type AnnouncementController struct {
	announcementService *services.AnnouncementService
}

// NewAnnouncementController 创建公告控制器实例
func NewAnnouncementController() *AnnouncementController {
	return &AnnouncementController{
		announcementService: services.NewAnnouncementService(),
	}
}

// GetAnnouncementList godoc
// @Summary 获取公告列表
// @Description 获取已发布的公告列表，按创建时间倒序排列，支持缓存（1分钟）
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param request body models.AnnouncementListRequest true "公告列表请求参数"
// @Success 200 {object} utils.Response{data=models.AnnouncementListResponse} "成功返回公告列表"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/announcements/list [post]
func (c *AnnouncementController) GetAnnouncementList(ctx *gin.Context) {
	var request models.AnnouncementListRequest
	
	// 绑定请求参数
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.InvalidParamsWithMessage(ctx, "请求参数错误: "+err.Error())
		return
	}
	
	// 设置默认值
	if request.Page <= 0 {
		request.Page = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = 10
	}
	if request.PageSize > 100 {
		request.PageSize = 100
	}
	
	// 获取公告列表
	result, err := c.announcementService.GetAnnouncementList(ctx, request.Page, request.PageSize)
	if err != nil {
		utils.InternalError(ctx)
		return
	}
	
	// 返回成功响应
	utils.SuccessWithMessage(ctx, "获取公告列表成功", result)
} 