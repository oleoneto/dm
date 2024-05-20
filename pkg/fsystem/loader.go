package fsystem

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"

	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/migrator"
)

type FileLoaderProtocol interface {
	// LoadFiles - Loads files from the given directory matching the provided regex.
	LoadFiles(dir string, pattern *regexp.Regexp) []fs.DirEntry
}

type MigrationBuilderProtocol interface {
	// Build - Build migrations from files.
	Build([]fs.DirEntry) (*ds.Queue[migrator.Migration], error)

	// LoadMigrations - Loads migrations from the given directory
	LoadMigrations(files []fs.DirEntry, dir string, pattern *regexp.Regexp) ([]migrator.Migration, error)
}

type FileGeneratorProtocol interface {
	// Generate - Creates a new file with the given content and format at the specified directory.
	GenerateFile(format, content, name, dir string)
}

// ------------------------------------------------------------------------------
// MARK: File Loader
// ------------------------------------------------------------------------------

var _ FileLoaderProtocol = (*FileLoader)(nil)

type FileLoader struct{}

func (fl FileLoader) LoadFiles(dir string, pattern *regexp.Regexp) []fs.DirEntry {
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

func (fl FileLoader) Build([]fs.DirEntry) ds.Queue[migrator.Migration] {
	return ds.Queue[migrator.Migration]{}
}

func (fl FileLoader) LoadMigrations(files []fs.DirEntry, dir string, pattern *regexp.Regexp) ([]migrator.Migration, error) {
	return []migrator.Migration{}, nil
}
