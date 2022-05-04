package cmd

import (
	"fmt"
	"regexp"
)

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
	VersionValidationPattern = regexp.MustCompile(`^\d{14}$`)
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
