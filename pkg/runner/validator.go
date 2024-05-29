package runner

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/oleoneto/dm/pkg/helpers"
	"github.com/oleoneto/dm/pkg/migrator"
	log "github.com/sirupsen/logrus"
)

// Validate - Checks if all migrations are valid
func Validate(ctx context.Context, migrations []migrator.Migration) error {
	if len(migrations) == 0 {
		log.WithField("func", helpers.GetCurrentFuncName()).Debugln("No migrations found...")
		return fmt.Errorf("no migrations found")
	}

	visitedNames := make(map[string]bool)
	visitedVersions := make(map[string]bool)

	for i := 0; i < len(migrations); i++ {
		// if item.Engine != r.engine.Name() { return fmt.Errorf("migration engine mismatch %v", item) }

		item := migrations[i]

		if _, ok := visitedVersions[item.Version]; ok {
			return fmt.Errorf(`duplicate (version): %v`, item)
		}

		if _, ok := visitedNames[item.Name]; ok {
			return fmt.Errorf(`duplicate (name): %v`, item)
		}

		var createTablePattern = *regexp.MustCompile(`CREATE TABLE (?P<TableName>\w+)`)
		var dropTablePattern = *regexp.MustCompile(`(DROP TABLE (IF EXISTS )?)(?P<TableName>\w+)`)

		var mismatchedInstructions int
		var mismatchedTables = map[string]string{}

		for _, change := range item.Changes.Up {
			if len(strings.Split(change, " ")) < 3 {
				return fmt.Errorf("missing (or invalid) migrate instruction: %v", item)
			}

			if createTablePattern.MatchString(change) {
				mismatchedInstructions += 1
				match := createTablePattern.FindStringSubmatch(change)
				table := match[createTablePattern.SubexpIndex("TableName")]
				mismatchedTables[table] = table
			}
		}

		for _, change := range item.Changes.Down {
			if len(strings.Split(change, " ")) < 3 {
				return fmt.Errorf(`missing (or invalid) rollback instruction: %v`, item)
			}

			if dropTablePattern.MatchString(change) {
				mismatchedInstructions -= 1
				match := dropTablePattern.FindStringSubmatch(change)
				table := match[dropTablePattern.SubexpIndex("TableName")]
				delete(mismatchedTables, table)
			}
		}

		if mismatchedInstructions != 0 || len(mismatchedTables) != 0 {
			return fmt.Errorf("CREATE and DROP instructions must always be paired: %v", item)
		}

		version, name, _ := strings.Cut(item.FileName, "_")
		name = strcase.ToCamel(strings.Split(name, ".")[0])

		if item.Version != version {
			return fmt.Errorf("version mismatch: %v", item)
		}

		if item.Name != name {
			return fmt.Errorf("name mismatch: %v", item)
		}

		visitedNames[item.Name] = true
		visitedVersions[item.Version] = true
	}

	return nil
}
