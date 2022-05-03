package postgresql

import (
	"context"
	"fmt"

	"github.com/cleopatrio/db-migrator-lib/engines"
	"github.com/cleopatrio/db-migrator-lib/migrations"
)

func (engine Postgres) Up(changes migrations.MigrationList) error {
	valid, reason := engine.Validate(changes)

	if !valid {
		fmt.Println(reason)
		return new(migrations.ValidationError)
	}

	engine.Connect()

	if changes.Size() < 1 || engine.IsUpToDate(changes) {
		fmt.Println("Migrations are up-to-date.")
		return nil
	}

	migration := changes.GetHead()

	for migration != nil {
		err := engine.performMigration(*migration, engine.registerMigration)

		if err != nil {
			fmt.Printf("\nMigration '%v' (%v) failed.\n%v \n", migration.Name, migration.Version, err)

			_ = engine.deregisterMigration(*migration)
			return err
		}

		migration = migration.Next()
	}

	return nil
}

func (engine Postgres) Down(changes migrations.MigrationList) error {

	valid, _ := engine.Validate(changes)

	if !valid {
		return new(migrations.ValidationError)
	}

	engine.Connect()

	if engine.IsEmpty() {
		fmt.Println("No migrations to rollback.")
		return nil
	}

	migration := changes.GetHead()

	for migration != nil {
		err := engine.performRollback(*migration, engine.deregisterMigration)

		if err != nil {
			fmt.Printf("\nRollback '%v' (%v) failed.\n%v \n", migration.Name, migration.Version, err)
			return err
		}

		migration = migration.Next()
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

	return callback(migration)
}

func (engine Postgres) performRollback(migration migrations.Migration, callback func(migrations.Migration) error) error {
	rows, err := Pg().Query(context.Background(), migration.Changes.Down)

	if err != nil {
		return err
	}

	err = rows.Scan()

	if err != nil {
		return err
	}

	return callback(migration)
}

func (engine Postgres) registerMigration(migration migrations.Migration) error {
	rows, _ := Pg().Query(
		context.Background(),
		engines.CreateMigrationEntry(engine.Table),
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
		engines.DropMigrationTable(engine.Table),
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
