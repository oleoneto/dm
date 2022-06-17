package migrations

import (
	"testing"
)

func TestSchemaTableExists(t *testing.T) {
	table := "schema_migrations"

	query := SchemaTableExists(table)
	formatted := `SELECT 
		TABLE_SCHEMA, 
		TABLE_NAME,
		TABLE_TYPE
		FROM 
			information_schema.TABLES 
		WHERE 
			TABLE_TYPE LIKE 'BASE TABLE' AND
			TABLE_NAME = 'schema_migrations';`

	if query != formatted {
		t.Fatalf(`got incorrect query %v`, query)
	}
}

func TestNumberOfAppliedMigrations(t *testing.T) {
	table := "schema_migrations"

	query := NumberOfAppliedMigrations(table)
	formatted := `SELECT COUNT(id) FROM schema_migrations;`

	if query != formatted {
		t.Fatalf(`got incorrect query %v`, query)
	}
}

func TestCreateMigrationTable(t *testing.T) {
	table := "schema_migrations"

	query := CreateMigrationTable(table)
	formatted := `CREATE TABLE schema_migrations (
		id SERIAL,
		version varchar UNIQUE NOT NULL,
		name varchar UNIQUE NOT NULL,
		created_at timestamp NOT NULL DEFAULT now(),

		PRIMARY KEY(id)
	);`

	if query != formatted {
		t.Fatalf(`got incorrect %v`, query)
	}
}

func TestDropMigrationTable(t *testing.T) {
	table := "schema_migrations"

	query := DropMigrationTable(table)

	if query != `DROP TABLE schema_migrations;` {
		t.Fatalf(`got incorrect %v`, query)
	}
}

func TestSelectMigrations(t *testing.T) {
	table := "schema_migrations"

	query := SelectMigrations(table)

	if query != `SELECT id, name, version FROM schema_migrations;` {
		t.Fatalf(`got incorrect %v`, query)
	}
}

func TestSelectMigrationsVersion(t *testing.T) {
	table := "schema_migrations"

	query := SelectMigrationsVersion(table)

	if query != `SELECT id, name, version, created_at FROM schema_migrations ORDER BY id DESC LIMIT 1;` {
		t.Fatalf(`got incorrect %v`, query)
	}
}

func TestCreateMigrationEntry(t *testing.T) {
	table := "schema_migrations"

	query := CreateMigrationEntry(table)

	if query != `INSERT INTO schema_migrations (version, name) VALUES ($1, $2);` {
		t.Fatalf(`got incorrect %v`, query)
	}
}

func TestDeleteMigrationEntry(t *testing.T) {
	table := "schema_migrations"

	query := DeleteMigrationEntry(table)

	if query != `DELETE FROM schema_migrations WHERE version = $1 AND name = $2;` {
		t.Fatalf(`got incorrect %v`, query)
	}

}

func TestSelectMigrationEntry(t *testing.T) {
	table := "schema_migrations"

	query := SelectMigrationEntry(table)

	if query != `SELECT id, name, version FROM schema_migrations WHERE version = $1 AND name = $2;` {
		t.Fatalf(`got incorrect %v`, query)
	}
}
