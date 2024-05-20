package runner

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/fsystem"
	"github.com/oleoneto/dm/pkg/migrator"
	"gopkg.in/yaml.v2"
)

var _ fsystem.MigrationBuilderProtocol = (*Runner)(nil)

func (r *Runner) Build(files []fs.DirEntry) (*ds.Queue[migrator.Migration], error) {
	migrations, err := r.LoadMigrations(files, r.trackerOptions.MigrationsDirectory, r.trackerOptions.FileRegexPattern)
	if err != nil {
		return nil, err
	}

	q := ds.NewFromSlice(migrations)

	return q, nil
}

func (r *Runner) LoadMigrations(files []fs.DirEntry, dir string, pattern *regexp.Regexp) ([]migrator.Migration, error) {
	var migrations []migrator.Migration

	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		path, _ = filepath.Abs(path)
		contents, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var m migrator.Migration
		err = yaml.Unmarshal(contents, &m)
		if err != nil {
			return nil, err
		}

		match := pattern.FindStringSubmatch(file.Name())

		m.FileName = file.Name()
		m.Version = match[pattern.SubexpIndex("Version")]

		migrations = append(migrations, m)
	}

	return migrations, nil
}
