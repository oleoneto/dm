package runner

import (
	"fmt"

	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/migrator"
)

var trackMigrationStatement = func(table string) string {
	// TODO: Save checksum
	return fmt.Sprintf(`INSERT INTO %v (version, name) VALUES ($1, $2)`, table)
}

var untrackMigrationStatement = func(table string) string {
	return fmt.Sprintf("DELETE FROM %v WHERE version = $1 AND name = $2;", table)
}

// MigrationUpStatement - Given a slice of migrations, extracts and assembles all migration commands into a single query
var migrationUpStatement = func(migrations ds.Queue[migrator.Migration]) string {
	var query string
	for _, m := range *migrations.RawData() {
		for _, stmt := range m.Changes.Up {
			query += stmt
		}
	}
	return query
}

// MigrationDownStatement - Given a slice of migrations, extracts and assembles all migration commands into a single query
var migrationDownStatement = func(migrations ds.Queue[migrator.Migration]) string {
	var query string
	for _, m := range *migrations.RawData() {
		for _, stmt := range m.Changes.Down {
			query += stmt
		}
	}
	return query
}
