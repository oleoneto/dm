package core

import (
	"errors"
	"regexp"
)

func ParseVersionArgs(flag string) (MigrationFilterFlag, error) {
	parsedFlag := MigrationFilterFlag{Value: flag}

	if NameValidationPattern.MatchString(flag) {
		parsedFlag.Type = "name"
		return parsedFlag, nil
	} else if VersionValidationPattern.MatchString(flag) {
		parsedFlag.Type = "version"
		return parsedFlag, nil
	}

	return parsedFlag, errors.New("invalid migration version or name")
}

type MigrationFilterFlag struct {
	Value, Type string
}

var (
	NameValidationPattern    = regexp.MustCompile(`^[a-zA-Z]+(\_?[a-zA-Z])*$`)
	VersionValidationPattern = regexp.MustCompile(`^\d{20}$`)
)
