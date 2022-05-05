package migrations

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v2"
)

var FilePattern = *regexp.MustCompile(`(?P<Version>^\d{20})_(?P<Name>[aA-zZ]+).ya?ml$`)

/*
Runner:
	Responsible for running and reverting migrations.
	The migration and rollback algorithm is self-contained within this type.

	For more flexibility, creations, reads, and deletions are responsibilities of
	the underlying store type. You can inject a type that conforms to the `Store` interface,
	and the `Runner` will call the store's appropriate methods when needed.
*/
type Runner struct {
	schemaTable string
	store       Store
}

// MARK: Accessors

func (runner *Runner) SetSchemaTable(table string) {
	runner.schemaTable = table
}

func (runner *Runner) SetStore(store Store) {
	runner.store = store
}

func (runner *Runner) GetSchemaTable() string {
	return runner.schemaTable
}

func (runner *Runner) GetStore() Store {
	return runner.store
}

// MARK: Migration Tracker

func IsTracked(store Store, schemaTable string) bool {
	var schema []TableSchema

	err := store.Read(SchemaTableExists(schemaTable), &schema)

	if err != nil {
		return false
	}

	return schema[0].TableName != ""
}

func IsEmpty(store Store, schemaTable string) bool {
	var count int

	tracked := IsTracked(store, schemaTable)

	if !tracked {
		return true
	}

	err := store.Read(NumberOfAppliedMigrations(schemaTable), &count)

	if err != nil {
		return false
	}

	return count != 0
}

func IsUpToDate(store Store, schemaTable string, migrations MigrationList) bool {
	tracked := IsTracked(store, schemaTable)

	if !tracked {
		StartTracking(store, schemaTable)
	}

	recent := migrations.GetTail()

	version, tracked := Version(store, schemaTable)
	return tracked && (version == recent.Version)
}

func Version(store Store, schemaTable string) (string, bool) {
	var versions []MigratorVersion
	var schema []TableSchema

	err := store.Read(SchemaTableExists(schemaTable), &schema)

	if err != nil || schema[0].TableName == "" {
		return "0", false
	}

	err = store.Read(SelectMigrationsVersion(schemaTable), &versions)

	// TODO: Fix generic type
	// err := adapter.ScanOne(&version)
	// fmt.Println("Here")
	// fmt.Println(err)

	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return "0", false
		} else if strings.Contains(err.Error(), "no rows") {
			return "0", true
		}

		fmt.Printf("%v\n", err)
	}

	return fmt.Sprintf("%v (%v).\nApplied at: %v", versions[0].Version, versions[0].Name, versions[0].CreatedAt), true
}

func StartTracking(store Store, schemaTable string) error {
	err := store.Create(CreateMigrationTable(schemaTable))

	if err != nil {
		return err
	}

	return nil
}

func StopTracking(store Store, schemaTable string) error {
	if !IsTracked(store, schemaTable) {
		return nil
	}

	err := store.Delete(DropMigrationTable(schemaTable))

	if err != nil {
		return err
	}

	return nil
}

// MARK: - Migration Runner

func (runner *Runner) Generate(name string, directory string) Migration {
	var migration Migration

	exclusionPattern := regexp.MustCompile(`(\-?\d{4} \-?\d{2} m=\+\d{1}.\d{9})|(\-)|(\W+)|(\:)|(\.)`)

	/*
		Input : 2022-05-04 18:49:19.478478 -0400 -04 m=+0.001942418
		Output: 20220504184919478478 (20 characters in total)

	*/
	now := time.Now().String()
	timestamp := exclusionPattern.ReplaceAllLiteralString(now, ``)

	// BUG: Regex needs updating. Timestamp sometimes returns < 20 characters.
	// FIX: Use a constant for this value instead.
	for len(timestamp) < 20 {
		now = time.Now().String()
		timestamp = exclusionPattern.ReplaceAllLiteralString(now, ``)
	}

	migration.Schema = 1
	migration.Engine = strings.ToLower(runner.store.Name())
	migration.Changes.Up = ""
	migration.Changes.Down = ""
	migration.Name = name
	migration.Version = timestamp
	migration.FileName = fmt.Sprintf("%v_%v.yaml", timestamp, strcase.ToSnake(name))

	content, err := yaml.Marshal(&migration)

	if err != nil {
		fmt.Printf("Error while marshaling. %v\n", err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%v/%v", directory, migration.FileName), content, 0644)

	if err != nil {
		fmt.Printf("Error while writing to file %v\n", err)
	}

	return migration
}

func (runner *Runner) Up(migrations MigrationList) error {
	runner.beforeAction()

	if migrations.Size() == 0 {
		fmt.Println("No migrations to run.")
		return nil
	}

	valid, reason := Validate(migrations)

	if !valid {
		fmt.Println(reason)
		return new(ValidationError)
	}

	if IsUpToDate(runner.store, runner.schemaTable, migrations) {
		fmt.Println("Migrations are up-to-date.")
		return nil
	}

	migration := migrations.GetHead()

	for migration != nil {
		err := runner.performMigration(*migration)

		if err != nil {
			fmt.Printf("\nMigration '%v' (%v) failed.\n%v \n", migration.Name, migration.Version, err)

			_ = runner.undoMigration(*migration, runner.schemaTable)
			return err
		}

		runner.registerMigration(*migration, runner.schemaTable)

		migration = migration.Next()
	}

	return nil
}

func (runner *Runner) Down(migrations MigrationList) error {
	runner.beforeAction()

	valid, _ := Validate(migrations)

	if !valid {
		return new(ValidationError)
	}

	if IsEmpty(runner.store, runner.schemaTable) {
		fmt.Println("No migrations to rollback.")
		return nil
	}

	migration := migrations.GetHead()

	for migration != nil {
		err := runner.performRollback(*migration)

		if err != nil {
			fmt.Printf("\nRollback '%v' (%v) failed.\n%v \n", migration.Name, migration.Version, err)
			return err
		}

		runner.undoMigration(*migration, runner.schemaTable)

		migration = migration.Next()
	}

	return nil
}

func (runner *Runner) PendingMigrations(directory string, filePattern *regexp.Regexp) MigrationList {
	runner.beforeAction()

	files := LoadFiles(directory, filePattern)
	list := BuildMigrations(files, directory, filePattern)

	migrated := Migrations{}
	res := MigrationList{}

	// NOTE: No migrations in database
	if IsEmpty(runner.store, runner.schemaTable) {
		return list
	}

	err := runner.store.Read(SelectMigrations(runner.schemaTable), &migrated)

	if err != nil {
		fmt.Printf("An error occurred.\nError: %v\n", err)
	}

	migratedHash := migrated.ToHash()

	curr := list.GetHead()

	for curr != nil {
		_, applied := migratedHash[curr.Version]

		if !applied {
			res.Insert(&Migration{
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

func (runner *Runner) AppliedMigrations(directory string, filePattern *regexp.Regexp) MigrationList {
	runner.beforeAction()

	migrated := Migrations{}
	res := MigrationList{}

	// NOTE: No migrations in database
	if IsEmpty(runner.store, runner.schemaTable) {
		return MigrationList{}
	}

	err := runner.store.Read(SelectMigrations(runner.schemaTable), &migrated)

	if err != nil {
		fmt.Printf("An error occurred.\nError: %v\n", err)
	}

	for _, curr := range migrated {
		res.Insert(&Migration{
			Engine:  runner.store.Name(),
			Id:      curr.Id,
			Name:    curr.Name,
			Version: curr.Version,
		})
	}

	return res
}

func (runner *Runner) Version() (string, bool) {
	runner.beforeAction()
	return Version(runner.store, runner.schemaTable)
}

// MARK: Helper for performing migration and rollback

func (runner *Runner) beforeAction() {
	if runner.GetSchemaTable() == "" {
		fmt.Println("No schema table provided.")
		os.Exit(1)
	}

	if runner.store == nil {
		fmt.Println("No store adapter specified.")
		os.Exit(1)
	}
}

func (runner *Runner) performMigration(migration Migration) error {
	err := runner.store.Create(migration.Changes.Up)

	if err != nil {
		return err
	}

	return nil
}

func (runner *Runner) performRollback(migration Migration) error {
	err := runner.store.Delete(migration.Changes.Down)

	if err != nil {
		return err
	}

	return nil
}

func (runner *Runner) registerMigration(migration Migration, table string) error {
	err := runner.store.Create(
		CreateMigrationEntry(table),
		migration.Version,
		migration.Name,
	)

	if err != nil {
		return err
	}

	fmt.Printf("Added version: %v. Name: %s\n", migration.Version, migration.Name)
	return nil
}

func (runner *Runner) undoMigration(migration Migration, table string) error {
	err := runner.store.Delete(
		DropMigrationTable(table),
		migration.Version,
		migration.Name,
	)

	if err != nil {
		return err
	}

	fmt.Printf("Removed version: %v, name: %s\n", migration.Version, migration.Name)
	return nil
}