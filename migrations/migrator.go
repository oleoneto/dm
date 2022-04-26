package migrations

import (
	"io/fs"
	"regexp"
)

type Migrator struct{}

var (
	MIGRATE_UP   = 0
	MIGRATE_DOWN = 1
	FILE_PATTERN = *regexp.MustCompile(`^\d{14}_[aA-zZ]+.ya?ml`)
)

func (M *Migrator) Perform(changes []Migration, mode int, engine Engine) error {
	switch mode {
	case MIGRATE_UP:
		return engine.Up(changes)

	case MIGRATE_DOWN:
		return engine.Down(changes)
	}

	return nil
}

func (M *Migrator) MatchingFiles(dir string, pattern *regexp.Regexp) ([]fs.FileInfo, error) {
	return MatchingFiles(dir, pattern)
}
