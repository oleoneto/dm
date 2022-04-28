package postgresql

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

func acquireDatabaseConnection() {
	conn, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	pgInstance = conn

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
}

func Pg() *pgxpool.Pool {
	return pgInstance
}

func QueryAndScan(query string, args ...any) error {
	rows, _ := Pg().Query(context.Background(), query, args...)
	return rows.Scan()
}

func QueryAndScanInterface(query string, dest interface{}, args ...any) error {
	rows, _ := Pg().Query(context.Background(), query, args...)
	return pgxscan.ScanAll(dest, rows)
}

func QueryAndScanOneInterface(query string, dest interface{}, args ...any) error {
	rows, _ := Pg().Query(context.Background(), query, args...)
	return pgxscan.ScanOne(dest, rows)
}

// ---

func performMigration(migration migrations.Migration, callback func(migrations.Migration) error) error {
	rows, _ := Pg().Query(context.Background(), migration.Changes.Up)

	err := rows.Scan()

	if err != nil {
		return err
	}

	err = callback(migration)

	return err
}

func registerMigration(migration migrations.Migration) error {
	rows, _ := Pg().Query(
		context.Background(),
		"INSERT INTO _migrations (version, name) VALUES ($1, $2);",
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

func deregisterMigration(migration migrations.Migration) error {
	rows, _ := Pg().Query(
		context.Background(),
		"DELETE FROM _migrations WHERE version = $1 AND name = $2;",
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
