package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"

	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Error   string
	Message string
	Data    interface{}
}

type APIMigrations struct {
	Count      int                   `json:"total"`
	Migrations migrations.Migrations `json:"migrations"`
}

type APIMessage struct {
	Message string `json:"message"`
}

type APIError struct {
	Error string `json:"error"`
}

func (*Controller) CallCommand(args []string, flags []string) (bytes.Buffer, bytes.Buffer, bool) {
	var standardOutput, standardError bytes.Buffer
	var exited bool

	cmd := exec.Command("dm", append(args, flags...)...)

	cmd.Stdout = &standardOutput
	cmd.Stderr = &standardError

	_ = cmd.Run()

	exited = cmd.ProcessState.Exited()

	return standardOutput, standardError, exited
}

func (*Controller) ExitOnError(exitCode int, ctx *gin.Context) {
	ctx.AbortWithStatus(500)
}

func (*Controller) HandleStderr(stderr bytes.Buffer, ctx *gin.Context) {
	response := APIError{}
	err := json.Unmarshal(stderr.Bytes(), &response)

	if err != nil {
		log.Println("Error: unable to parse stderr as JSON")
		response.Error = stderr.String()
	}

	ctx.AbortWithStatusJSON(500, response)
}

func (*Controller) HandleStdout(stdout bytes.Buffer, data interface{}, ctx *gin.Context, status int) {
	err := json.Unmarshal(stdout.Bytes(), &data)

	if err != nil {
		log.Println("Error: unable to parse stdout as JSON")
		ctx.AbortWithStatusJSON(500, err)
	}

	ctx.AbortWithStatusJSON(status, data)
}

func (controller *Controller) HandleCommandErrors(stdout, stderr bytes.Buffer, exited bool, ctx *gin.Context) bool {
	if exited && stderr.Len() == 0 && stdout.Len() == 0 {
		ctx.AbortWithStatus(500)
		return true
	}

	if stderr.Len() != 0 {
		controller.HandleStderr(stderr, ctx)
		return true
	}

	return false
}
