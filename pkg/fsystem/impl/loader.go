package impl

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/oleoneto/dm/pkg/migrator"
	"gopkg.in/yaml.v2"
)

// MARK: - Migration loader

func LoadMigration(file fs.FileInfo, parent string, pattern *regexp.Regexp) (*migrator.Migration, error) {
	path := filepath.Join(parent, file.Name())

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

	return &m, nil
}
