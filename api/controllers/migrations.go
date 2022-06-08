package controllers

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

type MigrationsController struct {
	Controller
}

type RequestBody struct {
	Migration string `json:"migration"`
}

var content = APIMigrations{}

// MARK: - Stateless Operations
// ------------------------------------------------------------------

func (controller *MigrationsController) List(ctx *gin.Context) {
	args := []string{"show", "all"}

	stdout, stderr, exited := controller.CallCommand(args, ctx.MustGet("command_flags").([]string))
	data := APIMigrations{}

	if controller.HandleCommandErrors(stdout, stderr, exited, ctx) {
		return
	}

	controller.ParseStdoutAsJSON(stdout, data, ctx)
}

func (controller *MigrationsController) Applied(ctx *gin.Context) {
	args := []string{"show", "applied"}

	stdout, stderr, exited := controller.CallCommand(args, ctx.MustGet("command_flags").([]string))
	data := APIMigrations{}

	if controller.HandleCommandErrors(stdout, stderr, exited, ctx) {
		return
	}

	controller.ParseStdoutAsJSON(stdout, data, ctx)
}

func (controller *MigrationsController) Pending(ctx *gin.Context) {
	args := []string{"show", "pending"}

	stdout, stderr, exited := controller.CallCommand(args, ctx.MustGet("command_flags").([]string))
	data := APIMigrations{}

	if controller.HandleCommandErrors(stdout, stderr, exited, ctx) {
		return
	}

	controller.ParseStdoutAsJSON(stdout, data, ctx)
}

// MARK: - Stateful Operations (will affect the state of the database)
// ------------------------------------------------------------------

func (controller *MigrationsController) Migrate(ctx *gin.Context) {
	args := []string{"migrate"}

	var requestBody RequestBody

	_ = ctx.BindJSON(&requestBody)

	response := []APIMessage{}
	stdout, stderr, exited := controller.CallCommand(args, append(ctx.MustGet("command_flags").([]string), requestBody.Migration))

	if controller.HandleCommandErrors(stdout, stderr, exited, ctx) {
		return
	}

	controller.HandleStdout(stdout, response, ctx, 202)
}

func (controller *MigrationsController) Rollback(ctx *gin.Context) {
	args := []string{"rollback"}

	var requestBody RequestBody

	_ = ctx.BindJSON(&requestBody)

	response := []APIMessage{}
	stdout, stderr, exited := controller.CallCommand(args, append(ctx.MustGet("command_flags").([]string), requestBody.Migration))

	if controller.HandleCommandErrors(stdout, stderr, exited, ctx) {
		return
	}

	controller.HandleStdout(stdout, response, ctx, 202)
}

// MARK: - Helpers
// ------------------------------------------------------------------

func (*MigrationsController) ParseStdoutAsJSON(stdout bytes.Buffer, data APIMigrations, ctx *gin.Context) {
	err := json.Unmarshal(stdout.Bytes(), &data.Migrations)

	if err != nil {
		log.Println("Error: unable to parse stdout as JSON")
		ctx.AbortWithStatusJSON(500, err)
	}

	data.Count = data.Migrations.Len()

	ctx.AbortWithStatusJSON(200, data)
}
