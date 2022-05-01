package migrations

import (
	"io/fs"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"runtime"
)

// MatchingFiles - Finds all files that statisfy a regex in the specified directory
func MatchingFiles(dir string, pattern *regexp.Regexp) ([]fs.FileInfo, error) {
	matches := []fs.FileInfo{}

	files, err := ioutil.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if pattern.MatchString(file.Name()) {
			matches = append(matches, file)
		}
	}

	return matches, nil
}

func BuildMigrations(files []fs.FileInfo, dir string, pattern *regexp.Regexp) MigrationList {
	var changes MigrationList

	for _, file := range files {
		var mg Migration

		_ = mg.Load(file, dir, pattern)

		changes.Insert(&mg)
	}

	return changes
}

func LoadFiles(dir string, pattern *regexp.Regexp) []fs.FileInfo {
	files, err := MatchingFiles(dir, pattern)

	if err != nil {
		return []fs.FileInfo{}
	}

	return files
}

func CurrentFilepath() string {
	_, name, _, _ := runtime.Caller(1)
	path := path.Join(path.Dir(name), ".")

	return path
}
