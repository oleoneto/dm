package core

import (
	"errors"
	"regexp"
)

func ParseVersionArgs(flag string) (VersionFlag, error) {
	parsedFlag := VersionFlag{Value: flag}

	if NameValidationPattern.MatchString(flag) {
		parsedFlag.Type = "Name"
		return parsedFlag, nil
	} else if VersionValidationPattern.MatchString(flag) {
		parsedFlag.Type = "Version"
		return parsedFlag, nil
	}

	return parsedFlag, errors.New("invalid migration version or name")
}

type VersionFlag struct {
	Value, Type string
}

var (
	NameValidationPattern    = regexp.MustCompile(`^[a-zA-Z]+(\_?[a-zA-Z])*$`)
	VersionValidationPattern = regexp.MustCompile(`^\d{20}$`)
)
