package main

import (
	"fmt"

	"github.com/cleopatrio/db-migrator-lib/engines"
	"github.com/cleopatrio/db-migrator-lib/migrations"
)

func init() {
	fmt.Println("Database Migrator v0.1.0-alpha")
}

func main() {
	var migrator migrations.Migrator
	var changes []migrations.Migration

	// TODO: Define these values based on command-line args
	// ---
	mode := migrations.MIGRATE_UP
	parent := "./examples"
	var psql engines.Postgres
	// ---

	files, _ := migrator.MatchingFiles(parent, &migrations.FILE_PATTERN)

	for _, file := range files {
		var mg migrations.Migration

		mg.Load(file, parent)

		changes = append(changes, mg)
	}

	migrator.Perform(changes, mode, psql)
}
