package runner

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/migrator"
)

// Validate - Checks if all migrations are valid
func Validate(ctx context.Context, migrations ds.Queue[migrator.Migration]) error {
	if migrations.IsEmpty() {
		return fmt.Errorf("no migrations found")
	}

	m := migrations.Dequeue()
	visitedNames := make(map[string]bool)
	visitedVersions := make(map[string]bool)

	for m != nil {
		// if m.Engine != r.engine.Name() { return fmt.Errorf("migration engine mismatch %v", m) }

		if _, ok := visitedVersions[m.Version]; ok {
			return fmt.Errorf(`duplicate (version): %v`, m)
		}

		if _, ok := visitedNames[m.Name]; ok {
			return fmt.Errorf(`duplicate (name): %v`, m)
		}

		var createTablePattern = *regexp.MustCompile(`CREATE TABLE (?P<TableName>\w+)`)
		var dropTablePattern = *regexp.MustCompile(`(DROP TABLE (IF EXISTS )?)(?P<TableName>\w+)`)

		var mismatchedInstructions int
		var mismatchedTables = map[string]string{}

		for _, change := range m.Changes.Up {
			if len(strings.Split(change, " ")) < 3 {
				return fmt.Errorf("missing (or invalid) migrate instruction: %v", m)
			}

			if createTablePattern.MatchString(change) {
				mismatchedInstructions += 1
				match := createTablePattern.FindStringSubmatch(change)
				table := match[createTablePattern.SubexpIndex("TableName")]
				mismatchedTables[table] = table
			}
		}

		for _, change := range m.Changes.Down {
			if len(strings.Split(change, " ")) < 3 {
				return fmt.Errorf(`missing (or invalid) rollback instruction: %v`, m)
			}

			if dropTablePattern.MatchString(change) {
				mismatchedInstructions -= 1
				match := dropTablePattern.FindStringSubmatch(change)
				table := match[dropTablePattern.SubexpIndex("TableName")]
				delete(mismatchedTables, table)
			}
		}

		if mismatchedInstructions != 0 || len(mismatchedTables) != 0 {
			return fmt.Errorf("CREATE and DROP instructions must always be paired: %v", m)
		}

		version, name, _ := strings.Cut(m.FileName, "_")
		name = strcase.ToCamel(strings.Split(name, ".")[0])

		if m.Version != version {
			return fmt.Errorf("version mismatch: %v", m)
		}

		if m.Name != name {
			return fmt.Errorf("name mismatch: %v", m)
		}

		visitedNames[m.Name] = true
		visitedVersions[m.Version] = true

		m = migrations.Dequeue()
	}

	return nil
}
