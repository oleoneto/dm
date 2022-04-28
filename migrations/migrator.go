package migrations

import (
	"fmt"
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
	MIGRATE_UP   = 0
	MIGRATE_DOWN = 1
	FILE_PATTERN = *regexp.MustCompile(`(?P<Version>^\d{14})_(?P<Name>[aA-zZ]+).ya?ml$`)
)

func (instance *Migrator) Build(dir string) Migrations {
	return BuildMigrations(dir, &FILE_PATTERN)
}

func (instance *Migrator) ListFiles(dir string) error {
	return ListFiles(dir, &FILE_PATTERN)
}

func (instance *Migrator) PendingMigrations(dir string) map[string]Migration {
	appliedMigrations := instance.Engine.AppliedMigrations()

	migrationFiles := instance.Build(dir)

	for _, file := range migrationFiles {
		key := fmt.Sprintf("%v_%v", file.Version, file.Name)

		_, applied := appliedMigrations[key]

		if !applied {
			fmt.Printf("Name: %v, Version: %v\n", file.Name, file.Version)
		}
	}

	return appliedMigrations
}

func (instance *Migrator) Run(changes []Migration, mode int) error {
	switch mode {
	case MIGRATE_UP:
		return instance.Engine.Up(changes)

	case MIGRATE_DOWN:
		return instance.Engine.Down(changes)
	}

	return nil
}

func (instance *Migrator) Status() {
	version, _ := instance.Engine.Version()
	fmt.Printf("Current version: %v\n", version)
}

func (instance *Migrator) Validate(changes Migrations) bool {
	return instance.Engine.Validate(changes)
}
