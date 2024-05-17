package fsystem

import (
	"io/fs"
	"regexp"

	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/migrator"
)

// MigrationLoaderFunc - Loads migrations from the given directory
type MigrationLoaderFunc func(file fs.FileInfo, sourceDir string, pattern *regexp.Regexp) (*migrator.Migration, error)

// FileLoaderFunc - Loads files from the given directory matching the provided regex
type FileLoaderFunc func(sourceDir string, pattern *regexp.Regexp) []fs.FileInfo

// MigrationBuilder - Given a slice of files, generates
type MigrationBuilder func([]fs.FileInfo) ds.Queue[migrator.Migration]
