package cmd

import (
	"fmt"
	"regexp"

	"github.com/drewstinnett/go-output-format/formatter"
)

type Formattable interface {
	Description() string
}

type VersionFlag struct {
	Value string
	Type  string
}

type InvalidFlagError struct{}

func (e InvalidFlagError) Error() string {
	return "invalid flag value"
}

func WithFormattedOutput(data Formattable) {
	config := &formatter.Config{
		Format:   format,
		Template: template,
	}

	output, err := formatter.OutputData(data, config)

	if err != nil {
		fmt.Println(err)
		return
	}

	switch config.Format {
	case "plain":
		fmt.Println(data.Description())
	default:
		fmt.Println(string(output))
	}
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

	fmt.Println("Error: invalid migration version or name")

	return parsedFlag, new(InvalidFlagError)
}
