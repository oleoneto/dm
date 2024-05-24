package runner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/oleoneto/dm/pkg/migrator"
	"gopkg.in/yaml.v2"
)

// FileLoaderFunc - Loads files from a file system
type FileLoaderFunc func(dir string, pattern *regexp.Regexp) []fs.DirEntry

func LoadFiles(dir string, pattern *regexp.Regexp) []fs.DirEntry {
	var matchingFiles = func(dir string, pattern *regexp.Regexp) ([]fs.DirEntry, error) {
		matches := []fs.DirEntry{}

		files, err := os.ReadDir(dir)

		if err != nil {
			fmt.Println(err)
			return matches, err
		}

		for _, file := range files {
			if pattern.MatchString(file.Name()) {
				matches = append(matches, file)
			}
		}

		return matches, nil
	}

	files, err := matchingFiles(dir, pattern)
	if err != nil {
		return []fs.DirEntry{}
	}

	return files
}

// LoadMigrationsFunc - Loads migrations from the given directory
type LoadMigrationsFunc func(files []fs.DirEntry, dir string, pattern *regexp.Regexp) ([]migrator.Migration, error)

func LoadMigrations(files []fs.DirEntry, dir string, pattern *regexp.Regexp) ([]migrator.Migration, error) {
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
