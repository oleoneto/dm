package runner

import (
	"fmt"
	"strings"

	"github.com/oleoneto/dm/pkg/migrator"
)

var trackMigrationStatement = func(schema, table string) string {
	// TODO: Save checksum
	return fmt.Sprintf(`INSERT INTO %v.%v (version, name) VALUES ($1, $2)`, schema, table)
}

var untrackMigrationStatement = func(schema, table string) string {
	return fmt.Sprintf("DELETE FROM %v.%v WHERE version = $1 AND name = $2;", schema, table)
}

// migrateUpStatement - Given a migration, extracts and assembles all migration instructions into a single query
var migrateUpStatement = func(m migrator.Migration) string {
	var query string
	for _, stmt := range m.Changes.Up {
		query += fmt.Sprintf("%v;\n", strings.TrimSuffix(strings.TrimSpace(stmt), ";"))
	}
	return query
}

// migrateDownStatement - Given a migration, extracts and assembles all rollback instructions into a single query
var migrateDownStatement = func(m migrator.Migration) string {
	var query string
	for _, stmt := range m.Changes.Down {
		query += fmt.Sprintf("%v;\n", strings.TrimSuffix(strings.TrimSpace(stmt), ";"))
	}
	return query
}
