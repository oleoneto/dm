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

func parsedVersionFlag(flag string) (VersionFlag, error) {
	parsedFlag := VersionFlag{Value: flag}

	versionPattern := regexp.MustCompile(`^\d{14}$`)
	namePattern := regexp.MustCompile(`^[aA-zZ]+_?[aA-zZ]+$`)

	if namePattern.MatchString(flag) {
		parsedFlag.Type = "Name"
		return parsedFlag, nil
	} else if versionPattern.MatchString(flag) {
		parsedFlag.Type = "Version"
		return parsedFlag, nil
	}

	fmt.Println("Error: invalid migration version or name")

	return parsedFlag, new(InvalidFlagError)
}
