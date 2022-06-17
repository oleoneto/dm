package migrations

import (
	"log"
	"os"
	"testing"

	"github.com/oleoneto/dm/stores"
)

var (
	emptyStore        = ExampleStore{}
	testPostgresStore = stores.Postgres{URL: os.Getenv("DATABASE_URL")}
)

// TestMain will exec each test, one by one
func TestMain(m *testing.M) {
	setUp()

	retCode := m.Run()

	tearDown()

	os.Exit(retCode)
}

func setUp() {
	log.Println("Creating public schema")
	testPostgresStore.Delete("DROP SCHEMA public CASCADE;")
	testPostgresStore.Delete("CREATE SCHEMA public;")
}

func tearDown() {
	log.Println("Dropping public schema")
}

func testRunner() Runner {
	return Runner{
		store:       testPostgresStore,
		schemaTable: "test_migrations",
	}
}

// ----------------------------------

func TestStoreIsEmpty(t *testing.T) {
	table := "schema_migrations"

	empty := IsEmpty(emptyStore, table)

	if !empty {
		t.Fatalf(`wanted empty, but got %v`, empty)
	}
}

func TestStoreIsTracked(t *testing.T) {
	table := "schema_migrations"

	tracked := IsTracked(emptyStore, table)

	if tracked {
		t.Fatalf(`wanted tracked == false, but got %v`, tracked)
	}
}

func TestStoreVersion(t *testing.T) {
	table := "schema_migrations"

	version, tracked := Version(emptyStore, table)

	if version.Version != "" || tracked {
		t.Fatalf(`wanted version == "0" && tracked == false, but got (%v, %v)`, version.Version, tracked)
	}
}

func TestStoreIsUpToDate(t *testing.T) {
	table := "schema_migrations"

	upToDate := IsUpToDate(emptyStore, table, defaultList())

	if upToDate {
		t.Fatalf(`wanted upToDate == false, but got %v`, upToDate)
	}
}

func TestStartTrackingStore(t *testing.T) {
	table := "schema_migrations"

	state := StartTracking(emptyStore, table)

	if state {
		t.Fatalf(`wanted state == false, but got %v`, state)
	}
}

func TestStopTrackingStore(t *testing.T) {
	table := "schema_migrations"

	state := StopTracking(emptyStore, table)

	// Empty store is not tracked. So, StopTracking should succeed.
	if !state {
		t.Fatalf(`wanted state == false, but got %v`, state)
	}
}

// =======================================
// MARK: - Logger

func TestRunnerLogError(t *testing.T) {
}

func TestRunnerLogInfo(t *testing.T) {
}

func TestRunnerSetLogger(t *testing.T) {
}

// =======================================
// MARK: - Accessors

func TestRunnerSchemaTableAccessors(t *testing.T) {
	runner := Runner{}

	// Set table name
	table := "dm_migrations"
	runner.SetSchemaTable(table)

	if runner.schemaTable == "" || runner.schemaTable != table {
		t.Fatalf(`expected %v, but got %v`, table, runner.schemaTable)
	}

	if runner.GetSchemaTable() == "" || runner.GetSchemaTable() != table {
		t.Fatalf(`expected %v, but got %v`, table, runner.schemaTable)
	}
}

func TestRunnerStoreAccessors(t *testing.T) {
	runner := Runner{}

	// Set store
	runner.SetStore(emptyStore)

	if runner.store != emptyStore {
		t.Fatalf(`expected %v, but got %v`, emptyStore.Name(), runner.store.Name())
	}

	if runner.GetStore() != emptyStore {
		t.Fatalf(`expected %v, but got %v`, emptyStore.Name(), runner.store.Name())
	}

	if runner.store.Name() != emptyStore.Name() {
		t.Fatalf(`expected %v, but got %v`, emptyStore.Name(), runner.store.Name())
	}
}

// MARK: - Postgres

/*
func TestRunnerBeforAction(t *testing.T) {
	runner := Runner{}
	runner.SetSchemaTable("dm_migrations")
	runner.SetStore(emptyStore)

	// FIXME: runner.beforeAction() is not testable
	// Scenario 0: Succeed
	// // Scenario 1: Exit if no table is provided
	// // Scenario 2: Exit if no adapter is provided
}
*/

func TestRunnerPerformMigration(t *testing.T) {
	runner := testRunner()

	// Scenario 1: Given a valid migration, write it to the database.
	err := runner.performMigration(*defaultMigrationList().head)

	if err != nil {
		t.Errorf(`expected no errors, but got %v`, err)
	}

	// Scenario 2: Given an invalid migration, exit and do not write to the database.
	invalidMigrations := []Migration{
		{
			Version:  "20221231054540",
			Engine:   "postgresql",
			Name:     "CreateLikes",
			FileName: "20221231054540_create_likes.yaml",
			Changes: Changes{
				Up:   []string{"CREATE TABLE likes (id, content_id INT NOT NULL)"},
				Down: []string{"DROP TABLE likes;"},
			},
		},
		{
			Version:  "20221231054541",
			Engine:   "postgresql",
			Name:     "CreateLikes",
			FileName: "20221231054541_create_likes.yaml",
			Changes: Changes{
				Up:   []string{"CREATE TABLE likes (id SERIAL, content_id)"},
				Down: []string{"DROP TABLE likes;"},
			},
		},
	}

	for _, migration := range invalidMigrations {
		err = runner.performMigration(migration)

		if err == nil {
			t.Errorf(`expected an error, but got %v`, err)
		}
	}

	// Cleanup
	testPostgresStore.Delete(defaultMigrationList().head.Changes.Down[0])
}

func TestRunnerRegisterMigration(t *testing.T) {
	runner := testRunner()

	// Scenario 1: Given a migration, it should not be added to the schema table if it does not exist
	err := runner.registerMigration(*defaultMigrationList().head, runner.schemaTable)

	if err == nil {
		t.Errorf(`expected a database error, but got %v`, err)
	}

	// Create schema table
	testPostgresStore.Create(CreateMigrationTable(runner.schemaTable))

	// Scenario 2: Given a migration, it should be added to the schema table
	err = runner.registerMigration(*defaultMigrationList().head, runner.schemaTable)

	if err != nil {
		t.Errorf(`expected no errors, but got %v`, err)
	}

	// Scenario 3: Given a duplicate migration, it should be added to the schema table
	err = runner.registerMigration(*defaultMigrationList().head, runner.schemaTable)

	if err == nil {
		t.Errorf(`expected a database error, but got %v`, err)
	}
}

func TestRunnerRemoveMigrationFromSchema(t *testing.T) {
	runner := testRunner()

	// Create schema table
	testPostgresStore.Create(CreateMigrationTable(runner.schemaTable))

	// Scenario 1: Given a migration, remove it from the schema table
	err := runner.removeMigrationFromSchema(*defaultMigrationList().head, runner.schemaTable)

	if err != nil {
		t.Errorf(`expected no errors, but got %v`, err)
	}

	// Scenario 2: Given a migration, error if a non-existing schema table is provided
	err = runner.removeMigrationFromSchema(*defaultMigrationList().head, "wrong_table")

	if err == nil {
		t.Errorf(`expected an error, but got %v`, err)
	}
}

func TestRunnerAppliedMigrations(t *testing.T) {
	runner := testRunner()
	directory := "./examples"
	filePattern := &FilePattern
	loadFromDir := false

	// Scenario 1: No applied migrations
	migrations := runner.AppliedMigrations(directory, filePattern, loadFromDir)

	if migrations.Size() != 0 {
		t.Errorf(`expected no migrations to have been applied, but got %v`, migrations.Size())
	}
}
