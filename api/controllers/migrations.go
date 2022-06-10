package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/oleoneto/dm/config"
)

type MigrationsController struct {
	Controller
}

type RequestBody struct {
	Migration string `json:"migration"`
}

// MARK: - Stateless Operations
// ------------------------------------------------------------------

func (controller *MigrationsController) List(ctx *gin.Context) {
	args := []string{"show", "all"}
	flags := ctx.MustGet("command_flags").([]string)

	response, err := StatelessExecutionStrategy(args, flags)

	if err != nil {
		ctx.IndentedJSON(config.SERVER_ERROR, response)
		return
	}

	ctx.IndentedJSON(config.SUCCESS, response)
}

func (controller *MigrationsController) Applied(ctx *gin.Context) {
	args := []string{"show", "applied"}
	flags := ctx.MustGet("command_flags").([]string)

	response, err := StatelessExecutionStrategy(args, flags)

	if err != nil {
		ctx.IndentedJSON(config.SERVER_ERROR, response)
		return
	}

	ctx.IndentedJSON(config.SUCCESS, response)
}

func (controller *MigrationsController) Pending(ctx *gin.Context) {
	args := []string{"show", "pending"}
	flags := ctx.MustGet("command_flags").([]string)

	response, err := StatelessExecutionStrategy(args, flags)

	if err != nil {
		ctx.IndentedJSON(config.SERVER_ERROR, response)
		return
	}

	ctx.IndentedJSON(config.SUCCESS, response)
}

// MARK: - Stateful Operations (will affect the state of the database)
// ------------------------------------------------------------------

func (controller *MigrationsController) Migrate(ctx *gin.Context) {
	var requestBody RequestBody
	_ = ctx.Bind(&requestBody)

	args := []string{"migrate"}
	flags := append(ctx.MustGet("command_flags").([]string), requestBody.Migration)

	response, err := StatefulExecutionStrategy(args, flags)

	if err != nil {
		ctx.IndentedJSON(config.SERVER_ERROR, response)
		return
	}

	ctx.IndentedJSON(config.ACCEPTED, response)
}

func (controller *MigrationsController) Rollback(ctx *gin.Context) {
	var requestBody RequestBody
	_ = ctx.Bind(&requestBody)

	args := []string{"rollback"}
	flags := append(ctx.MustGet("command_flags").([]string), requestBody.Migration)

	response, err := StatefulExecutionStrategy(args, flags)

	if err != nil {
		ctx.IndentedJSON(config.SERVER_ERROR, response)
		return
	}

	ctx.IndentedJSON(config.ACCEPTED, response)
}
