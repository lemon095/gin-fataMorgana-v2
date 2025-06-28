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
	dbErr := database.HealthCheck()

	// 检查Redis连接
	redisErr := database.RedisClient.Ping(c.Request.Context()).Err()

	// 获取数据库统计信息
	dbStats := database.GetDBStats()

	// 获取查询统计信息
	queryStats := database.GetQueryStats()

	healthStatus := "healthy"
	if dbErr != nil || redisErr != nil {
		healthStatus = "unhealthy"
	}

	utils.Success(c, gin.H{
		"status": healthStatus,
		"services": gin.H{
			"database": gin.H{
				"status": func() string {
					if dbErr != nil {
						return "unhealthy"
					}
					return "healthy"
				}(),
				"error": func() string {
					if dbErr != nil {
						return dbErr.Error()
					}
					return ""
				}(),
				"stats":       dbStats,
				"query_stats": queryStats,
			},
			"redis": gin.H{
				"status": func() string {
					if redisErr != nil {
						return "unhealthy"
					}
					return "healthy"
				}(),
				"error": func() string {
					if redisErr != nil {
						return redisErr.Error()
					}
					return ""
				}(),
			},
		},
	})
}

// DatabaseStats 数据库统计信息
func (hc *HealthController) DatabaseStats(c *gin.Context) {
	stats := database.GetDBStats()
	if stats == nil {
		utils.Error(c, utils.CodeDatabaseError)
		return
	}

	utils.Success(c, gin.H{
		"database_stats": stats,
	})
}

// QueryStats 查询统计信息
func (hc *HealthController) QueryStats(c *gin.Context) {
	stats := database.GetQueryStats()
	if stats == nil {
		utils.Error(c, utils.CodeDatabaseError)
		return
	}

	utils.Success(c, gin.H{
		"query_stats": stats,
	})
}

// PerformanceOptimization 性能优化建议
func (hc *HealthController) PerformanceOptimization(c *gin.Context) {
	optimization := database.OptimizeQueries()

	utils.Success(c, gin.H{
		"optimization_recommendations": optimization,
	})
}

// DatabaseHealth 数据库健康检查
func (hc *HealthController) DatabaseHealth(c *gin.Context) {
	err := database.HealthCheck()
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"status":  "healthy",
		"message": "数据库连接正常",
	})
}

// RedisHealth Redis健康检查
func (hc *HealthController) RedisHealth(c *gin.Context) {
	err := database.RedisClient.Ping(c.Request.Context()).Err()
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeRedisError, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"status":  "healthy",
		"message": "Redis连接正常",
	})
}
