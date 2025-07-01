package controllers

import (
	"gin-fataMorgana/database"
	"gin-fataMorgana/utils"

	"github.com/gin-gonic/gin"
)

// HealthController 健康检查控制器
type HealthController struct{}

// NewHealthController 创建健康检查控制器实例
func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthCheck 系统健康检查
func (hc *HealthController) HealthCheck(c *gin.Context) {
	// 检查数据库连接
	dbStatus := "healthy"
	if err := database.HealthCheck(); err != nil {
		dbStatus = "unhealthy"
	}

	// 检查Redis连接
	redisStatus := "healthy"
	if err := database.RedisClient.Ping(c.Request.Context()).Err(); err != nil {
		redisStatus = "unhealthy"
	}

	// 返回健康状态
	status := gin.H{
		"status": "healthy",
		"services": gin.H{
			"database": dbStatus,
			"redis":    redisStatus,
			},
	}

	if dbStatus == "unhealthy" || redisStatus == "unhealthy" {
		status["status"] = "unhealthy"
		utils.ErrorWithMessage(c, utils.CodeServer, "系统服务异常")
		return
	}

	utils.Success(c, status)
}

// DatabaseHealth 数据库健康检查
func (hc *HealthController) DatabaseHealth(c *gin.Context) {
	if err := database.HealthCheck(); err != nil {
		utils.ErrorWithMessage(c, utils.CodeServer, "数据库连接失败")
		return
	}
	utils.Success(c, gin.H{"status": "healthy"})
}

// RedisHealth Redis健康检查
func (hc *HealthController) RedisHealth(c *gin.Context) {
	if err := database.RedisClient.Ping(c.Request.Context()).Err(); err != nil {
		utils.ErrorWithMessage(c, utils.CodeServer, "Redis连接失败")
		return
	}
	utils.Success(c, gin.H{"status": "healthy"})
}
