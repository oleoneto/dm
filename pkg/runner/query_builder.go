package runner

import (
	"fmt"
	"strings"

	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/migrator"
)

var trackMigrationStatement = func(schema, table string) string {
	// TODO: Save checksum
	return fmt.Sprintf(`INSERT INTO %v.%v (version, name) VALUES ($1, $2)`, schema, table)
}

var untrackMigrationStatement = func(schema, table string) string {
	return fmt.Sprintf("DELETE FROM %v.%v WHERE version = $1 AND name = $2;", schema, table)
}

var migrateUpStatement = func(m migrator.Migration) string {
	var query string
	for _, stmt := range m.Changes.Up {
		query += fmt.Sprintf("%v;\n", strings.TrimSuffix(strings.TrimSpace(stmt), ";"))
	}
	return query
}

var migrateDownStatement = func(m migrator.Migration) string {
	var query string
	for _, stmt := range m.Changes.Down {
		query += fmt.Sprintf("%v;\n", strings.TrimSuffix(strings.TrimSpace(stmt), ";"))
	}
	return query
}

// migrateUpStatements - Given a slice of migrations, extracts and assembles all migration commands into a single query
var migrateUpStatements = func(migrations ds.Queue[migrator.Migration]) string {
	var query string
	for _, m := range *migrations.RawData() {
		query += migrateUpStatement(m)
	}
	return query
}

// migrateDownStatements - Given a slice of migrations, extracts and assembles all migration commands into a single query
var migrateDownStatements = func(migrations ds.Queue[migrator.Migration]) string {
	var query string
	for _, m := range *migrations.RawData() {
		query += migrateDownStatement(m)
	}
	return query
}
