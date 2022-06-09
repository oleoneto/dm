package cmd

import (
	"os"
	"regexp"

	"github.com/cleopatrio/db-migrator-lib/logger"
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
