package middleware

import (
	"github.com/gin-gonic/gin"
)

func CorsHeaders(origin string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Origin", origin)
		ctx.Next()
	}
}
