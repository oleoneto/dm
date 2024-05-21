package fsystem

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
)

type FileLoaderProtocol interface {
	// LoadFiles - Loads files from the given directory matching the provided regex.
	LoadFiles(dir string, pattern *regexp.Regexp) []fs.DirEntry
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
