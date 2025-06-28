package controllers

import (
	"strconv"
	"time"

	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// LoginLogController 登录记录控制器
type LoginLogController struct {
	loginLogService *services.LoginLogService
}

// NewLoginLogController 创建登录记录控制器
func NewLoginLogController() *LoginLogController {
	return &LoginLogController{
		loginLogService: services.NewLoginLogService(),
	}
}

// GetUserLoginHistory 获取用户登录历史
func (c *LoginLogController) GetUserLoginHistory(ctx *gin.Context) {
	uid := ctx.Param("uid")
	if uid == "" {
		utils.ErrorWithMessage(ctx, utils.CodeInvalidParams, "用户UID不能为空")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	logs, total, err := c.loginLogService.GetUserLoginHistory(ctx, uid, page, size)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeOperationFailed, "获取登录历史失败")
		return
	}

	utils.Success(ctx, gin.H{
		"logs":  logs,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetUserLastLogin 获取用户最后登录信息
func (c *LoginLogController) GetUserLastLogin(ctx *gin.Context) {
	uid := ctx.Param("uid")
	if uid == "" {
		utils.ErrorWithMessage(ctx, utils.CodeInvalidParams, "用户UID不能为空")
		return
	}

	log, err := c.loginLogService.GetUserLastLogin(ctx, uid)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeUserNotFound, "未找到登录记录")
		return
	}

	utils.Success(ctx, log)
}

// GetLoginStats 获取登录统计信息
func (c *LoginLogController) GetLoginStats(ctx *gin.Context) {
	uid := ctx.Param("uid")
	if uid == "" {
		utils.ErrorWithMessage(ctx, utils.CodeInvalidParams, "用户UID不能为空")
		return
	}

	stats, err := c.loginLogService.GetLoginStats(ctx, uid)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeOperationFailed, "获取统计信息失败")
		return
	}

	utils.Success(ctx, stats)
}

// GetLoginLogsByTimeRange 按时间范围获取登录记录
func (c *LoginLogController) GetLoginLogsByTimeRange(ctx *gin.Context) {
	uid := ctx.Param("uid")
	if uid == "" {
		utils.ErrorWithMessage(ctx, utils.CodeInvalidParams, "用户UID不能为空")
		return
	}

	startTimeStr := ctx.Query("start_time")
	endTimeStr := ctx.Query("end_time")

	if startTimeStr == "" || endTimeStr == "" {
		utils.ErrorWithMessage(ctx, utils.CodeInvalidParams, "开始时间和结束时间不能为空")
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeInvalidParams, "开始时间格式错误")
		return
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeInvalidParams, "结束时间格式错误")
		return
	}

	logs, err := c.loginLogService.GetLoginLogsByTimeRange(ctx, uid, startTime, endTime)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeOperationFailed, "获取登录记录失败")
		return
	}

	utils.Success(ctx, logs)
}

// GetLoginLogsByIP 按IP地址获取登录记录
func (c *LoginLogController) GetLoginLogsByIP(ctx *gin.Context) {
	uid := ctx.Param("uid")
	ip := ctx.Query("ip")

	if uid == "" {
		utils.ErrorWithMessage(ctx, utils.CodeInvalidParams, "用户UID不能为空")
		return
	}

	if ip == "" {
		utils.ErrorWithMessage(ctx, utils.CodeInvalidParams, "IP地址不能为空")
		return
	}

	logs, err := c.loginLogService.GetLoginLogsByIP(ctx, uid, ip)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeOperationFailed, "获取登录记录失败")
		return
	}

	utils.Success(ctx, logs)
}

// CleanOldLogs 清理旧登录记录
func (c *LoginLogController) CleanOldLogs(ctx *gin.Context) {
	daysStr := ctx.DefaultQuery("days", "90")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 {
		utils.ErrorWithMessage(ctx, utils.CodeInvalidParams, "天数参数错误")
		return
	}

	err = c.loginLogService.CleanOldLogs(ctx, days)
	if err != nil {
		utils.ErrorWithMessage(ctx, utils.CodeOperationFailed, "清理旧记录失败")
		return
	}

	utils.Success(ctx, gin.H{"message": "清理完成"})
}
