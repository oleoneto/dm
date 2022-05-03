package engines

import "fmt"

func CreateMigrationTable(table string) string {
	return fmt.Sprintf(`CREATE TABLE %v (
		id SERIAL,
		version varchar UNIQUE NOT NULL,
		name varchar UNIQUE NOT NULL,
		created_at timestamp NOT NULL DEFAULT now(),

		PRIMARY KEY(id)
	);`, table)
}

func DropMigrationTable(table string) string {
	return fmt.Sprintf("DROP TABLE %v;", table)
}

func SelectMigrations(table string) string {
	return fmt.Sprintf("SELECT id, name, version FROM %v;", table)
}

func SelectMigrationsVersion(table string) string {
	return fmt.Sprintf("SELECT id, name, version, created_at FROM %v ORDER BY id DESC LIMIT 1;", table)
}

func CreateMigrationEntry(table string) string {
	return fmt.Sprintf("INSERT INTO %v (version, name) VALUES ($1, $2);", table)
}

func DeleteMigrationEntry(table string) string {
	return fmt.Sprintf("DELETE FROM %v WHERE version = $1 AND name = $2;", table)
}

func SelectMigrationEntry(table string) string {
	return fmt.Sprintf("SELECT id, name, version FROM %v WHERE version = $1 AND name = $2;", table)
}
