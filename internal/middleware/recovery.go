package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %s", r)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Server Error",
				})
			}
		}()
		c.Next()
	}
}
