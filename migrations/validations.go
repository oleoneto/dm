package migrations

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

func Validate(changes MigrationList) bool {
	visitedNames := map[string]bool{}
	visitedVersions := map[string]bool{}

	change := changes.head

	for change != nil {
		if visitedNames[change.Name] {
			return invalidChange(*change, "duplicate migration name")
		}

		if visitedVersions[change.Version] {
			return invalidChange(*change, "duplicate migration version")
		}

		// TODO: Check if migration is using a supported engine
		// if !supportedEngines[change.Engine] {
		// 	return invalidChange(change, "unsupported database engine")
		// }

		if len(strings.Split(change.Changes.Up, " ")) < 5 {
			return invalidChange(*change, "missing (or invalid) migrate instruction")
		}

		if len(strings.Split(change.Changes.Down, " ")) < 3 {
			return invalidChange(*change, "missing (or invalid) rollback instruction")
		}

		version, name, _ := strings.Cut(change.FileName, "_")
		name = strings.Split(name, ".")[0]
		name = strcase.ToCamel(name)

		if change.Version != version {
			return invalidChange(*change, "version mismatch")
		}

		if change.Name != name {
			return invalidChange(*change, "name mismatch")
		}

		visitedNames[change.Name] = true
		visitedVersions[change.Version] = true

		change = change.next
	}

	return true
}

func invalidChange(change Migration, reason string) bool {
	fmt.Printf("Invalid migration: %v (%v). Reason: %v.\n", change.Name, change.Version, reason)
	return false
}
