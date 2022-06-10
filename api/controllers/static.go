package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/oleoneto/dm/config"
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

func (StaticController) Documentation(ctx *gin.Context) {
	ctx.HTML(config.SUCCESS, "swagger.html", gin.H{"title": "Database Migrator"})
}

func (controller *StaticController) Health(ctx *gin.Context) {
	args := []string{"show", "pending"}
	flags := ctx.MustGet("command_flags").([]string)

	response, err := StatelessExecutionStrategy(args, flags)

	if err != nil {
		ctx.IndentedJSON(config.SERVER_ERROR, response)
		return
	}

	data, castable := response.(APIMigrations)

	if castable {
		// -- Unhealthy
		if len(data.Migrations) != 0 {
			ctx.IndentedJSON(config.SERVER_ERROR, APIError{Error: fmt.Sprintf("%v pending migrations", len(data.Migrations))})
			return
		}

		// -- Healthy
		ctx.IndentedJSON(config.SUCCESS, APIMessage{Message: "No pending migrations"})
		return
	}

	ctx.IndentedJSON(config.SUCCESS, response)
}
