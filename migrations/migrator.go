package migrations

import (
	"fmt"
	"io/fs"
	"regexp"
)

type Migrator struct {
	Engine           Engine
	Directory        string
	DatabaseUrl      string
	Table            string
	SupportedEngines map[string]Engine
}

var (
	MigrateUp            = 0
	MigrateDown          = 1
	MigrationFilePattern = *regexp.MustCompile(`(?P<Version>^\d{14})_(?P<Name>[aA-zZ]+).ya?ml$`)
)

func (instance *Migrator) Build(dir string) MigrationList {
	return BuildMigrations(dir, &MigrationFilePattern)
}

func (instance *Migrator) ListFiles(dir string) []fs.FileInfo {
	return ListFiles(dir, &MigrationFilePattern)
}

func (instance *Migrator) PendingMigrations(dir string) MigrationList {
	appliedMigrations := instance.Engine.AppliedMigrations()
	migrations := instance.Build(dir)
	var sequence MigrationList

	migration := migrations.head

	for migration != nil {
		key := fmt.Sprintf("%v_%v", migration.Version, migration.Name)

		_, applied := appliedMigrations[key]

		if !applied {
			sequence.Insert(migration)
		}

		migration = migration.next
	}

	return sequence
}

func (instance *Migrator) Run(changes MigrationList, mode int) error {
	switch mode {
	case MigrateUp:
		return instance.Engine.Up(changes)

	case MigrateDown:
		return instance.Engine.Down(changes)
	}

	return nil
}

func (instance *Migrator) Status() {
	version, _ := instance.Engine.Version()
	fmt.Printf("Current version: %v\n", version)
}

func (instance *Migrator) Validate(changes MigrationList) bool {
	return instance.Engine.Validate(changes)
}
