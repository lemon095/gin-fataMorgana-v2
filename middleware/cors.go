package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// 允许的域名列表
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:8080", 
			"http://localhost:5173",
			"https://colombiatkadmin.com",
			"http://colombiatkadmin.com",
			"https://www.colombiatkadmin.com",
			"http://www.colombiatkadmin.com",
			"https://colombiatk.com",
			"http://colombiatk.com",
			"https://www.colombiatk.com",
			"http://www.colombiatk.com",
		}
		
		// 检查Origin是否在允许列表中
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}
		
		// 设置CORS头
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-CSRF-Token, X-API-Key, Cache-Control, Pragma, Referer, User-Agent, Accept-Language, token")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "43200") // 12小时

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
} 