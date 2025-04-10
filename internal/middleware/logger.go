package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// After request
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		log.Printf("[GIN] %v | %3d | %13v | %15s |%-7s %#v\n",
			startTime.Format("2006/01/02 - 15:04:05"),
			statusCode,
			duration,
			clientIP,
			method,
			path,
		)
	}
}
