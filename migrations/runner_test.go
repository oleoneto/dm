package migrations

import (
	"os"
	"testing"

	"github.com/oleoneto/dm/stores"
)

var testPostgresStore = stores.Postgres{URL: os.Getenv("DATABASE_URL")}

func TestMain(m *testing.M) {
	setUp()

	retCode := m.Run()

	tearDown()

	os.Exit(retCode)
}

func setUp() {
	rebuildDatabaseSchema()
}

func tearDown() {
	rebuildDatabaseSchema()
}

func rebuildDatabaseSchema() {
	testPostgresStore.Delete("DROP SCHEMA public CASCADE;")
	testPostgresStore.Delete("CREATE SCHEMA public;")
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

	// Scenario 1: An empty store without migrations
	empty := IsEmpty(testPostgresStore, table)

	if !empty {
		t.Fatalf(`wanted empty, but got %v`, empty)
	}

	// Scenario 2: A non-empty store with migrations
	table = "test_migrations"
	runner := testRunner()
	runner.schemaTable = table
	runner.store = testPostgresStore

	runner.Up(defaultList())

	empty = IsEmpty(testPostgresStore, table)

	if empty {
		t.Fatalf(`wanted non empty, but got %v`, empty)
	}

	t.Cleanup(rebuildDatabaseSchema)
}

func TestStoreIsTracked(t *testing.T) {
	table := "schema_migrations"

	tracked := IsTracked(testPostgresStore, table)

	if tracked {
		t.Fatalf(`wanted tracked == false, but got %v`, tracked)
	}
}

func TestStoreVersion(t *testing.T) {
	// Scenario 1: Empty store. Not tracked.
	table := "schema_migrations"

	version, tracked := Version(testPostgresStore, table)

	if version.Version != "" || tracked {
		t.Fatalf(`wanted version == "0" && tracked == false, but got (%v, %v)`, version.Version, tracked)
	}

	// Scenario 2: Non-empty store. Tracked.
	runner := testRunner()
	runner.schemaTable = table
	runner.store = testPostgresStore

	list := defaultList()

	runner.Up(list)

	version, tracked = Version(testPostgresStore, table)

	if version.Version != list.tail.Version || !tracked {
		t.Errorf(`wanted version == list.tail.version and tracked = true, but got (%v, %v)`, version.Version, tracked)
	}

	version, tracked = runner.Version()

	if version.Version != list.tail.Version || !tracked {
		t.Errorf(`wanted version == list.tail.version and tracked = true, but got (%v, %v)`, version.Version, tracked)
	}

	t.Cleanup(rebuildDatabaseSchema)
}

func TestStoreIsUpToDate(t *testing.T) {
	table := "schema_migrations"

	upToDate := IsUpToDate(testPostgresStore, table, defaultList())

	if upToDate {
		t.Fatalf(`wanted upToDate == false, but got %v`, upToDate)
	}
}

func TestStartAndStopTrackingStore(t *testing.T) {
	table := "schema_migrations"

	// Start tracking
	tracking := StartTracking(testPostgresStore, table)

	if !tracking {
		t.Errorf(`wanted tracking == true, but got %v`, tracking)
	}

	// Stop tracking
	stopped := StopTracking(testPostgresStore, table)

	if !stopped {
		t.Errorf(`wanted stopped == true, but got %v`, stopped)
	}

	// Stop tracking
	table = "unknown_migrations_table"

	stopped = StopTracking(testPostgresStore, table)

	if !stopped {
		t.Errorf(`wanted stopped == true, but got %v`, stopped)
	}

	t.Cleanup(rebuildDatabaseSchema)
}

// =======================================
// MARK: - Logger

// TODO: Implement tests
func TestRunnerLogError(t *testing.T)  {}
func TestRunnerLogInfo(t *testing.T)   {}
func TestRunnerSetLogger(t *testing.T) {}

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
	runner.SetStore(testPostgresStore)

	if runner.store != testPostgresStore {
		t.Fatalf(`expected %v, but got %v`, testPostgresStore.Name(), runner.store.Name())
	}

	if runner.GetStore() != testPostgresStore {
		t.Fatalf(`expected %v, but got %v`, testPostgresStore.Name(), runner.store.Name())
	}

	if runner.store.Name() != testPostgresStore.Name() {
		t.Fatalf(`expected %v, but got %v`, testPostgresStore.Name(), runner.store.Name())
	}
}

// =======================================

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

	t.Cleanup(rebuildDatabaseSchema)
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

	t.Cleanup(rebuildDatabaseSchema)
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

	t.Cleanup(rebuildDatabaseSchema)
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

	// TODO: Scenario 2: Unable to read migrations from the database
	// TODO: Scenario 3: Applied migrations

	t.Cleanup(rebuildDatabaseSchema)
}

func TestRunnerPendingMigrations(t *testing.T) {
	// TODO: Scenario 1: No migration file found
	// TODO: Scenario 2: Unable to read migrations from the database
	// TODO: Scenario 3: Pending migrations
}

// TODO: Implement  tests
func TestRunnerUp(t *testing.T) {}

func TestRunnerDown(t *testing.T) {}

func TestRunnerGenerate(t *testing.T) {}
