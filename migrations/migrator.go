package migrations

import (
	"fmt"
	"regexp"
)

type Migrator struct {
	Engine    Engine
	Directory string
}

var (
	MIGRATE_UP   = 0
	MIGRATE_DOWN = 1
	FILE_PATTERN = *regexp.MustCompile(`(?P<Version>^\d{14})_(?P<Name>[aA-zZ]+).ya?ml$`)
)

func (instance *Migrator) Build(dir string) []Migration {
	var changes []Migration

	files, _ := MatchingFiles(dir, &FILE_PATTERN)

	for _, file := range files {
		var mg Migration

		mg.Load(file, dir)

		changes = append(changes, mg)
	}

	return changes
}

func (instance *Migrator) ListFiles(dir string) error {
	files, err := MatchingFiles(dir, &FILE_PATTERN)

	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}

	return nil
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
	status, _ := instance.Engine.Version()
	fmt.Println(status)
}
