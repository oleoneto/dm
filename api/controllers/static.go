package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

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

func (StaticController) Health(ctx *gin.Context) {
	args := []string{"show", "pending"}
	cmd_args := ctx.MustGet("command_flags").([]string)
	cmd_args = append(args, cmd_args...)

	output, err := exec.Command("dm", cmd_args...).Output()

	if err != nil {
		log.Println("Error: failed to execute command", err)
		log.Println("Command args:", cmd_args)
		log.Println("Command output:", string(output))
		return
	}

	err = json.Unmarshal(output, &content.Migrations)

	if err != nil {
		log.Println("Error: failed to encode as JSON", err)
		log.Println("Content:", output)
		return
	}

	content.Count = content.Migrations.Len()

	// Healthy service
	if content.Migrations.Len() == 0 {
		message := APIMessage{Message: "No pending migrations"}
		ctx.IndentedJSON(200, message)
	}

	// Unheathy
	message := APIMessage{Message: fmt.Sprintf("%v pending migrations", content.Migrations.Len())}
	ctx.IndentedJSON(500, message)
}
