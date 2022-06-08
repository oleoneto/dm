package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type StaticController struct {
	Controller
}

func (StaticController) Ping(ctx *gin.Context) {
	message := map[string]string{
		"status":  "OK",
		"message": "System is operational",
	}

	ctx.IndentedJSON(200, message)
}

func (controller *StaticController) Health(ctx *gin.Context) {
	args := []string{"show", "pending"}
	flags := ctx.MustGet("command_flags").([]string)

	response, err := StatelessExecutionStrategy(args, flags)

	if err != nil {
		ctx.IndentedJSON(ERROR_STATUS_CODE, response)
		return
	}

	data, castable := response.(APIMigrations)

	if castable {
		// -- Unhealthy
		if len(data.Migrations) != 0 {
			ctx.IndentedJSON(ERROR_STATUS_CODE, APIError{Error: fmt.Sprintf("%v pending migrations", len(data.Migrations))})
			return
		}

		// -- Healthy
		ctx.IndentedJSON(SUCCESS_STATUS_CODE, APIMessage{Message: "No pending migrations"})
		return
	}

	ctx.IndentedJSON(SUCCESS_STATUS_CODE, response)
}
