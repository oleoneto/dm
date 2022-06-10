package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/oleoneto/dm/config"
)

type ErrorResponse struct {
	Errors []string
}

func ConfigurationMiddleware(configuration config.APIConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errors, isValid := configuration.IsValid()

		if !isValid {
			ctx.AbortWithStatusJSON(config.SERVER_ERROR, ErrorResponse{Errors: errors})
			return
		}

		flags := []string{
			"--output-format", "json",
			"--directory", configuration.Directory,
			"--table", configuration.Table,
		}

		ctx.Set("command_flags", flags)
		ctx.Next()
	}
}
