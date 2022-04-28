package postgresql

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/cleopatrio/db-migrator-lib/migrations"
)

func (engine Postgres) Up(changes migrations.Migrations) error {

	if !engine.Validate(changes) {
		return new(migrations.ValidationError)
	}

	engine.acquireDatabaseConnection()

	if engine.IsUpToDate(changes) {
		fmt.Println("Migrations are up-to-date.")
		return nil
	}

	for _, migration := range changes {
		err := engine.performMigration(migration, engine.registerMigration)

		if err != nil {
			fmt.Printf("\nMigration '%v' (%v) failed.\n%v \n", migration.Name, migration.Version, err)

			_ = engine.deregisterMigration(migration)
			return err
		}
	}

	return nil
}

func (engine Postgres) Down(changes migrations.Migrations) error {
	engine.acquireDatabaseConnection()

	if engine.IsEmpty() {
		fmt.Println("No migrations to rollback.")
		return nil
	}

	// Rollback migrations in reverse order to account for entity dependencies.
	sort.Sort(sort.Reverse(changes))

	for _, migration := range changes {
		rows, _ := Pg().Query(context.Background(), migration.Changes.Down)

		err := rows.Scan()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

		engine.deregisterMigration(migration)
	}

	return nil
}

// MARK: - Helpers

func (engine Postgres) performMigration(migration migrations.Migration, callback func(migrations.Migration) error) error {
	rows, _ := Pg().Query(context.Background(), migration.Changes.Up)

	err := rows.Scan()

	if err != nil {
		return err
	}

	err = callback(migration)

	return err
}

func (engine Postgres) registerMigration(migration migrations.Migration) error {
	rows, _ := Pg().Query(
		context.Background(),
		fmt.Sprintf("INSERT INTO %v (version, name) VALUES ($1, $2);", engine.Table),
		migration.Version,
		migration.Name,
	)

	err := rows.Scan()

	if err != nil {
		return err
	}

	fmt.Printf("Added version: %v. Name: %s\n", migration.Version, migration.Name)
	return nil
}

func (engine Postgres) deregisterMigration(migration migrations.Migration) error {
	rows, _ := Pg().Query(
		context.Background(),
		fmt.Sprintf("DELETE FROM %v WHERE version = $1 AND name = $2;", engine.Table),
		migration.Version,
		migration.Name,
	)

	err := rows.Scan()

	if err != nil {
		return err
	}

	fmt.Printf("Removed version: %v, name: %s\n", migration.Version, migration.Name)
	return nil
}
