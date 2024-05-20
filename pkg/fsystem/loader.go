package fsystem

import (
	"io/fs"
	"regexp"

	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/migrator"
)

type FileLoaderProtocol interface {
	// LoadFiles - Loads files from the given directory matching the provided regex.
	LoadFiles(dir string, pattern *regexp.Regexp) []fs.FileInfo
}

type MigrationBuilderProtocol interface {
	// Build - Build migrations from files.
	Build([]fs.FileInfo) (*ds.Queue[migrator.Migration], error)

	// LoadMigrations - Loads migrations from the given directory
	LoadMigrations(files []fs.FileInfo, dir string, pattern *regexp.Regexp) ([]migrator.Migration, error)
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

func (fl FileLoader) LoadFiles(dir string, pattern *regexp.Regexp) []fs.FileInfo {
	return []fs.FileInfo{}
}

func (fl FileLoader) Build([]fs.FileInfo) ds.Queue[migrator.Migration] {
	return ds.Queue[migrator.Migration]{}
}

func (fl FileLoader) LoadMigrations(files []fs.FileInfo, dir string, pattern *regexp.Regexp) ([]migrator.Migration, error) {
	return []migrator.Migration{}, nil
}
