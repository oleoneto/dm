package postgresql

import (
	"log"
	"strings"

	"github.com/cleopatrio/db-migrator-lib/migrations"
)

const engineName string = "PostgreSQL"

func (engine Postgres) IsUpToDate(changes migrations.Migrations) bool {
	if !engine.IsTracked() {
		engine.StartTracking()
	}

	recent := changes[len(changes)-1]
	version, tracked := engine.Version()
	return tracked && (version == recent.Version)
}

func (engine Postgres) IsTracked() bool {
	_, tracked := engine.Version()
	return tracked
}

func (engine Postgres) Version() (string, bool) {
	version := MigratorVersion{}

	acquireDatabaseConnection()

	err := QueryAndScanOneInterface(
		"SELECT * FROM _migrations ORDER BY version DESC LIMIT 1;",
		&version,
	)

	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			// VERBOSE: fmt.Printf("%s: Database is not yet being tracked.\n", engineName)
			return "", false
		} else if strings.Contains(err.Error(), "no rows") {
			// VERBOSE: fmt.Printf("%s: No migrations yet.\n", engineName)
			return "", true
		}

		log.Fatalf("%s: Error checking status. %v", engineName, err)
	}

	return version.Version, true
}

func (engine Postgres) StartTracking() error {
	if engine.IsTracked() {
		return nil
	}

	return QueryAndScan(
		`CREATE TABLE _migrations (
			id SERIAL,
			version varchar UNIQUE NOT NULL,
			name varchar
		);`,
	)
}

func (engine Postgres) StopTracking() error {
	if !engine.IsTracked() {
		return nil
	}

	return QueryAndScan("DROP TABLE _migrations;")
}

func (engine Postgres) IsEmpty() bool {
	version, tracked := engine.Version()
	return tracked && (version == "")
}
