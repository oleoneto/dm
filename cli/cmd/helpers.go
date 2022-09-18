package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/oleoneto/dm/logger"
)

type MessageOutput struct {
	Message string
}

type ErrorOutput struct {
	Error string
}

func (t MessageOutput) Description() string {
	return t.Message
}

func (t ErrorOutput) Description() string {
	return t.Error
}

type VersionFlag struct {
	Value string
	Type  string
}

type InvalidFlagError struct{}

func (e InvalidFlagError) Error() string {
	return "invalid flag value"
}

var (
	NameValidationPattern    = regexp.MustCompile(`^[a-zA-Z]+(\_?[a-zA-Z])*$`)
	VersionValidationPattern = regexp.MustCompile(`^\d{20}$`)
)

func parsedVersionFlag(flag string) (VersionFlag, error) {
	parsedFlag := VersionFlag{Value: flag}

	if NameValidationPattern.MatchString(flag) {
		parsedFlag.Type = "Name"
		return parsedFlag, nil
	} else if VersionValidationPattern.MatchString(flag) {
		parsedFlag.Type = "Version"
		return parsedFlag, nil
	}

	message := logger.ApplicationMessage{
		Message: "Error: invalid migration version or name",
	}
	logger.Custom(format, template).WithFormattedOutput(&message, os.Stderr)

	return parsedFlag, new(InvalidFlagError)
}

func readFromStdin() (input string, err error) {
	if flag.NArg() == 0 {
		reader := bufio.NewReader(os.Stdin)

		input, err = reader.ReadString(';')

		if err != nil {
			log.Fatalln("failed to read input")
		}

		// otherwise, we would have a blank line
		input = strings.TrimSpace(input)
	} else {
		input = flag.Arg(0)
	}

	fmt.Println("File content:\n", input)
	return input, err
}
