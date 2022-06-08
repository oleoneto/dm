package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os/exec"
)

type Controller struct {
	Error   string
	Message string
	Data    interface{}
}

type Migration struct {
	Id       int    `yaml:"-" json:"-"`
	FileName string `yaml:"-"`
	Version  string `yaml:"version"`
	Name     string `yaml:"name"`
}

type APIMigrations struct {
	Count      int         `json:"total"`
	Migrations []Migration `json:"migrations"`
}

type APIMessage struct {
	Message string `json:"message"`
}

type APIError struct {
	Error string `json:"error"`
}

var (
	ERROR_STATUS_CODE            = 500
	SUCCESS_STATUS_CODE          = 200
	STATEFUL_SUCCESS_STATUS_CODE = 202
)

func StatelessExecutionStrategy(args, flags []string) (interface{}, error) {

	stdout, stderr, exited := CallCommand(args, flags)

	outbuf, hasErrors := CheckForCommandErrors(stdout, stderr, exited)

	if hasErrors {
		message := APIError{Error: outbuf.String()}
		return message, errors.New(outbuf.String())
	}

	data, err := ParseStdoutAsJSON(outbuf)

	return data, err
}

func StatefulExecutionStrategy(args, flags []string) (interface{}, error) {
	stdout, stderr, exited := CallCommand(args, flags)

	outbuf, hasErrors := CheckForCommandErrors(stdout, stderr, exited)

	if hasErrors {
		return APIError{Error: outbuf.String()}, errors.New(outbuf.String())
	}

	data, err := ParseStdoutAsJSON(outbuf)

	return data, err
}

func CallCommand(args []string, flags []string) (bytes.Buffer, bytes.Buffer, bool) {
	var standardOutput, standardError bytes.Buffer
	var exited bool

	cmd := exec.Command("dm", append(args, flags...)...)

	cmd.Stdout = &standardOutput
	cmd.Stderr = &standardError

	_ = cmd.Run()

	exited = cmd.ProcessState.Exited()

	return standardOutput, standardError, exited
}

func CheckForCommandErrors(stdout, stderr bytes.Buffer, exited bool) (bytes.Buffer, bool) {
	// A populated stderr should always indicate the presence of errors.
	if stderr.Len() != 0 {
		log.Println("stderr was populated with", stderr.String())
		return stderr, true
	}

	// Command was halted bu no further infomation was provided by the exec.Command
	if exited && stderr.Len() == 0 && stdout.Len() == 0 {
		log.Println("total silence from exec.Command")
		return *bytes.NewBufferString("An error occcurred on the server while running the specified command"), true
	}

	return stdout, false
}

func ParseStdoutAsJSON(stdout bytes.Buffer) (interface{}, error) {
	var data interface{}
	var err error

	message := APIMessage{}
	migrations := APIMigrations{}

	err = ParseStdoutAsMessage(stdout, &message)
	data = message

	if err != nil {
		err = ParseStdoutAsMigrations(stdout, &migrations)
		data = migrations
	}

	return data, err
}

func ParseStdoutAsMigrations(stdout bytes.Buffer, data *APIMigrations) error {
	err := json.Unmarshal(stdout.Bytes(), &data.Migrations)

	if err != nil {
		log.Println("Error: unable to parse stdout as APIMigrations", err)
		return err
	}

	data.Count = len(data.Migrations)

	return nil
}

func ParseStdoutAsMessage(stdout bytes.Buffer, message *APIMessage) error {
	err := json.Unmarshal(stdout.Bytes(), &message)

	if err != nil {
		log.Println("Error: unable to parse stdout as APIMessage", err)
		return err
	}

	return nil
}
