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

func (engine Postgres) IsUpToDate(changes migrations.Migrations) bool {
	if !engine.IsTracked() {
		engine.StartTracking()
	}

	recent := changes[len(changes)-1]
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

	rows, _ := Pg().Query(
		context.Background(),
		fmt.Sprintf(`CREATE TABLE %v (
			id SERIAL,
			version varchar UNIQUE NOT NULL,
			name varchar
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
