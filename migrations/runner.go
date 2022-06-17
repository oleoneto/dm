package migrations

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/oleoneto/dm/logger"
	"gopkg.in/yaml.v2"
)

var (
	FilePattern        = *regexp.MustCompile(`(?P<Version>^\d{20})_(?P<Name>[aA-zZ]+).yaml$`)
	CreateTablePattern = *regexp.MustCompile(`CREATE TABLE (?P<TableName>\w+)`)
	DropTablePattern   = *regexp.MustCompile(`(DROP TABLE (IF EXISTS )?)(?P<TableName>\w+)`)
)

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
	logger      logger.Logger
}

// MARK: Logger

func (runner *Runner) LogError(err string) {
	message := logger.ApplicationError{Error: err}
	runner.logger.WithFormattedOutput(&message, os.Stderr)
}

func (runner *Runner) LogInfo(info string) {
	message := logger.ApplicationMessage{Message: info}
	runner.logger.WithFormattedOutput(&message, os.Stdout)
}

// MARK: Accessors

func (runner *Runner) SetLogger(format, template string) {
	runner.logger = logger.Custom(format, template)
}

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

// MARK: - Returns `true` if the schemaTable is found in the database. Returns `false` in any other case.
func IsTracked(store Store, schemaTable string) bool {
	var schema []TableSchema

	err := store.Read(SchemaTableExists(schemaTable), &schema)

	if err != nil || len(schema) == 0 {
		return false
	}

	return schema[0].TableName != ""
}

// MARK: - Return `true` if the schemaTable is found and has no rows.
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
		tracking := StartTracking(store, schemaTable)

		if !tracking {
			return false
		}
	}

	recent := migrations.GetTail()

	version, tracked := Version(store, schemaTable)
	return tracked && (version.Version == recent.Version)
}

func Version(store Store, schemaTable string) (MigratorVersion, bool) {
	var versions []MigratorVersion
	var schema []TableSchema

	emptyVersion := MigratorVersion{}

	err := store.Read(SchemaTableExists(schemaTable), &schema)

	if err != nil || len(schema) == 0 || schema[0].TableName == "" {
		return emptyVersion, false
	}

	err = store.Read(SelectMigrationsVersion(schemaTable), &versions)

	if err != nil {
		fmt.Printf("%v\n", err)
		return emptyVersion, false
	}

	if len(versions) == 0 {
		return emptyVersion, true
	}

	return versions[0], true
}

func StartTracking(store Store, schemaTable string) bool {
	err := store.Create(CreateMigrationTable(schemaTable))

	return err != nil
}

func StopTracking(store Store, schemaTable string) bool {
	if !IsTracked(store, schemaTable) {
		return true
	}

	err := store.Delete(DropMigrationTable(schemaTable))

	return err != nil
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

	migration.Schema = 2
	migration.Engine = strings.ToLower(runner.store.Name())
	migration.Changes.Up = []string{""}
	migration.Changes.Down = []string{""}
	migration.Name = name
	migration.Version = timestamp
	migration.FileName = fmt.Sprintf("%v_%v.yaml", timestamp, strcase.ToSnake(name))

	content, err := yaml.Marshal(&migration)

	if err != nil {
		runner.LogError(fmt.Sprintf("Error while marshaling. %v\n", err))
	}

	err = ioutil.WriteFile(fmt.Sprintf("%v/%v", directory, migration.FileName), content, 0644)

	if err != nil {
		runner.LogError(fmt.Sprintf("Error while writing to file %v\n", err))
	}

	return migration
}

func (runner *Runner) Up(migrations MigrationList) error {
	runner.beforeAction()

	if migrations.Size() == 0 {
		runner.LogInfo("No migrations to run.")
		return nil
	}

	valid, reason := Validate(migrations)

	if !valid {
		runner.LogError(reason)
		return new(ValidationError)
	}

	if IsUpToDate(runner.store, runner.schemaTable, migrations) {
		runner.LogInfo("Migrations are up-to-date.")
		return nil
	}

	migration := migrations.GetHead()

	for migration != nil {
		err := runner.performMigration(*migration)

		if err != nil {
			_ = runner.removeMigrationFromSchema(*migration, runner.schemaTable)
			runner.LogError(fmt.Sprintf("\nMigration '%v' (%v) failed.\n%v \n", migration.Name, migration.Version, err))
			return err
		}

		runner.registerMigration(*migration, runner.schemaTable)

		migration = migration.Next()
	}

	runner.logger.ReleaseCachedMessages(os.Stdout)

	return nil
}

func (runner *Runner) Down(migrations MigrationList) error {
	runner.beforeAction()

	valid, reason := Validate(migrations)

	if !valid {
		runner.LogError(reason)
		return new(ValidationError)
	}

	if IsEmpty(runner.store, runner.schemaTable) {
		runner.LogInfo("No migrations to rollback.")
		return nil
	}

	migration := migrations.GetHead()

	for migration != nil {
		// Perform the migration's rollback instruction (down)

		for _, change := range migration.Changes.Down {
			err := runner.store.Delete(change)

			if err != nil {
				runner.LogError(fmt.Sprintf("\nRollback '%v' (%v) failed.\n%v \n", migration.Name, migration.Version, err))
				return err
			}
		}

		err := runner.removeMigrationFromSchema(*migration, runner.schemaTable)

		if err != nil {
			return err
		}

		migration = migration.Next()
	}

	runner.logger.ReleaseCachedMessages(os.Stdout)

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
		runner.LogError(fmt.Sprintf("An error occurred.\nError: %v\n", err))
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

func (runner *Runner) AppliedMigrations(directory string, filePattern *regexp.Regexp, loadFromDir bool) MigrationList {
	runner.beforeAction()

	migrated := Migrations{}
	res := MigrationList{}

	// NOTE: No migrations in database
	if IsEmpty(runner.store, runner.schemaTable) {
		return MigrationList{}
	}

	err := runner.store.Read(SelectMigrations(runner.schemaTable), &migrated)

	if err != nil {
		runner.LogError(fmt.Sprintf("An error occurred.\nError: %v\n", err))
	}

	for _, curr := range migrated {
		m := Migration{
			Engine:   runner.store.Name(),
			Id:       curr.Id,
			Name:     curr.Name,
			Version:  curr.Version,
			FileName: fmt.Sprintf(`%v_%v.yaml`, curr.Version, strcase.ToSnake(curr.Name)),
		}

		// NOTE: Create a representation of the underlying file and
		// use it to load the file stored on the disk
		if loadFromDir {
			m.Load(MigrationFile{name: m.FileName}, directory, filePattern)
		}

		res.Insert(&m)
	}

	return res
}

func (runner *Runner) Version() (MigratorVersion, bool) {
	runner.beforeAction()
	return Version(runner.store, runner.schemaTable)
}

// MARK: Helper for performing migration and rollback

func (runner *Runner) beforeAction() {
	if runner.GetSchemaTable() == "" {
		runner.LogError("No schema table provided.")
		os.Exit(1)
	}

	if runner.store == nil {
		runner.LogError("No store adapter specified.")
		os.Exit(1)
	}
}

func (runner *Runner) performMigration(migration Migration) error {
	for _, change := range migration.Changes.Up {
		err := runner.store.Create(change)

		if err != nil {
			return err
		}
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

	runner.logger.CacheMessage(migration)
	return nil
}

func (runner *Runner) removeMigrationFromSchema(migration Migration, table string) error {
	err := runner.store.Delete(
		DeleteMigrationEntry(table),
		migration.Version,
		migration.Name,
	)

	if err != nil {
		return err
	}

	runner.logger.CacheMessage(migration)
	return nil
}
