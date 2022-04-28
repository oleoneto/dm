package postgresql

import (
	"fmt"
	"os"
	"sort"

	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Name string `default:"PostgreSQL"`
}

type MigratorVersion struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

var pgInstance *pgxpool.Pool

func (engine Postgres) Up(changes migrations.Migrations) error {
	acquireDatabaseConnection()

	if engine.IsUpToDate(changes) {
		fmt.Println("Migrations are up-to-date.")
		return nil
	}

	for _, migration := range changes {
		err := performMigration(migration, registerMigration)

		if err != nil {
			fmt.Printf("\nMigration '%v' (%v) failed.\n%v \n", migration.Name, migration.Version, err)

			_ = deregisterMigration(migration)
			return err
		}
	}

	return nil
}

func (engine Postgres) Down(changes migrations.Migrations) error {
	acquireDatabaseConnection()

	if engine.IsEmpty() {
		fmt.Println("No migrations to rollback.")
		return nil
	}

	// Rollback migrations in reverse order to account for entity dependencies.
	sort.Sort(sort.Reverse(changes))

	for _, migration := range changes {
		err := QueryAndScan(migration.Changes.Down)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

		deregisterMigration(migration)
	}

	return nil
}
