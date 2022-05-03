package postgresql

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"regexp"
	"strings"

	"github.com/cleopatrio/db-migrator-lib/engines"
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

	engine.Connect()

	rows, _ := Pg().Query(
		context.Background(),
		engines.SelectMigrationsVersion(engine.Table),
	)

	err := pgxscan.ScanOne(&version, rows)

	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			// VERBOSE: fmt.Printf("%s: Database is not yet being tracked.\n", engine.Name)
			return "0", false
		} else if strings.Contains(err.Error(), "no rows") {
			// VERBOSE: fmt.Printf("%s: No migrations yet.\n", engine.Name)
			return "0", true
		}

		log.Fatalf("%s: Error checking status. %v", engine.Name, err)
	}

	return fmt.Sprintf("%v (%v).\nApplied at: %v", version.Version, version.Name, version.CreatedAt), true
}

func (engine Postgres) StartTracking() error {
	if engine.IsTracked() {
		return nil
	}

	rows, _ := Pg().Query(
		context.Background(),
		engines.CreateMigrationTable(engine.Table),
	)

	return rows.Scan()
}

func (engine Postgres) StopTracking() error {
	if !engine.IsTracked() {
		return nil
	}

	rows, _ := Pg().Query(
		context.Background(),
		engines.DropMigrationTable(engine.Table),
	)

	return rows.Scan()
}

func (engine Postgres) LoadFiles(dir string, pattern *regexp.Regexp) []fs.FileInfo {
	return migrations.LoadFiles(dir, pattern)
}

func (engine Postgres) BuildMigrations(files []fs.FileInfo) migrations.MigrationList {
	return migrations.BuildMigrations(files, engine.Directory, engine.FilePattern)
}

func (engine Postgres) PendingMigrations() migrations.MigrationList {
	files := engine.LoadFiles(engine.Directory, engine.FilePattern)
	list := engine.BuildMigrations(files)

	migrated := migrations.Migrations{}
	res := migrations.MigrationList{}

	// NOTE: No migrations in database
	if engine.IsEmpty() || !engine.IsTracked() {
		return list
	}

	engine.Connect()

	rows, _ := Pg().Query(context.Background(), engines.SelectMigrations(engine.Table))
	err := pgxscan.ScanAll(&migrated, rows)

	if err != nil {
		fmt.Printf("%v: An error occurred.\nError: %v\n", engine.Name, err)
	}

	migratedHash := migrated.ToHash()

	curr := list.GetHead()

	for curr != nil {
		_, applied := migratedHash[curr.Version]

		if !applied {
			res.Insert(&migrations.Migration{
				Changes:  curr.Changes,
				Engine:   curr.Engine,
				FileName: curr.FileName,
				Id:       curr.Id,
				Name:     curr.Name,
				Schema:   curr.Schema,
				Version:  curr.Version,
			})
		}

		curr = curr.Next()
	}

	return res
}

func (engine Postgres) AppliedMigrations() migrations.MigrationList {
	files := engine.LoadFiles(engine.Directory, engine.FilePattern)
	list := engine.BuildMigrations(files)

	migrated := migrations.Migrations{}
	res := migrations.MigrationList{}

	// NOTE: No migrations in database
	if engine.IsEmpty() || !engine.IsTracked() {
		return migrations.MigrationList{}
	}

	engine.Connect()

	rows, _ := Pg().Query(context.Background(), engines.SelectMigrations(engine.Table))
	err := pgxscan.ScanAll(&migrated, rows)

	if err != nil {
		fmt.Printf("%v: An error occurred.\nError: %v\n", engine.Name, err)
	}

	migratedHash := migrated.ToHash()

	curr := list.GetHead()

	for curr != nil {
		_, applied := migratedHash[curr.Version]

		if applied {
			res.Insert(&migrations.Migration{
				Changes:  curr.Changes,
				Engine:   curr.Engine,
				FileName: curr.FileName,
				Id:       curr.Id,
				Name:     curr.Name,
				Schema:   curr.Schema,
				Version:  curr.Version,
			})
		}

		curr = curr.Next()
	}

	return res
}
