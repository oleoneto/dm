package postgresql

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/georgysavva/scany/pgxscan"
)

func (engine Postgres) IsTracked() bool {
	_, tracked := engine.Version()
	return tracked
}

func (engine Postgres) IsEmpty() bool {
	version, tracked := engine.Version()
	return tracked && (version == "")
}

func (engine Postgres) IsUpToDate(changes migrations.MigrationList) bool {
	if !engine.IsTracked() {
		engine.StartTracking()
	}

	recent := changes.GetTail()
	version, tracked := engine.Version()
	return tracked && (version == recent.Version)
}

func (engine Postgres) Version() (string, bool) {
	version := migrations.MigratorVersion{}

	engine.acquireDatabaseConnection()

	rows, _ := Pg().Query(
		context.Background(),
		fmt.Sprintf("SELECT * FROM %v ORDER BY version DESC LIMIT 1;", engine.Table),
	)

	err := pgxscan.ScanOne(&version, rows)

	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			// VERBOSE: fmt.Printf("%s: Database is not yet being tracked.\n", engineName)
			return "0", false
		} else if strings.Contains(err.Error(), "no rows") {
			// VERBOSE: fmt.Printf("%s: No migrations yet.\n", engineName)
			return "0", true
		}

		log.Fatalf("%s: Error checking status. %v", engineName, err)
	}

	return version.Version, true
}

func (engine Postgres) StartTracking() error {
	if engine.IsTracked() {
		return nil
	}

	rows, _ := Pg().Query(
		context.Background(),
		fmt.Sprintf(`CREATE TABLE %v (
			id SERIAL,
			version varchar UNIQUE NOT NULL,
			name varchar UNIQUE NOT NULL,

			PRIMARY KEY(id)
		);`, engine.Table),
	)

	return rows.Scan()
}

func (engine Postgres) StopTracking() error {
	if !engine.IsTracked() {
		return nil
	}

	rows, _ := Pg().Query(
		context.Background(),
		fmt.Sprintf("DROP TABLE %v;", engine.Table),
	)

	return rows.Scan()
}

func (engine Postgres) AppliedMigrations() map[string]migrations.Migration {
	var applied migrations.Migrations
	mapping := map[string]migrations.Migration{}

	engine.acquireDatabaseConnection()

	rows, _ := Pg().Query(
		context.Background(),
		fmt.Sprintf("SELECT * FROM %v;", engine.Table),
	)

	err := pgxscan.ScanAll(&applied, rows)

	if err != nil {
		fmt.Printf("%v: An error occurred.\nError: %v\n", engineName, err)
		return mapping
	}

	// FIX: May need to refactor this logic
	for _, m := range applied {
		key := fmt.Sprintf("%v_%v", m.Version, m.Name)
		mapping[key] = m
	}

	return mapping
}
