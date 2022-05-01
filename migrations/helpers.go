package migrations

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

// MatchingFiles - Finds all files that statisfy a regex in the specified directory
func MatchingFiles(dir string, pattern *regexp.Regexp) ([]fs.FileInfo, error) {
	matches := []fs.FileInfo{}

	files, err := ioutil.ReadDir(dir)

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

func BuildMigrations(files []fs.FileInfo, dir string, pattern *regexp.Regexp) MigrationList {
	var changes MigrationList

	for _, file := range files {
		var mg Migration

		err := mg.Load(file, dir, pattern)

		if err == nil {
			changes.Insert(&mg)
		}
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

func Validate(changes MigrationList) (bool, string) {
	visitedNames := map[string]bool{}
	visitedVersions := map[string]bool{}

	change := changes.head

	for change != nil {
		if visitedVersions[change.Version] {
			return invalidChange(*change, "duplicate migration version")
		}

		if visitedNames[change.Name] {
			return invalidChange(*change, "duplicate migration name")
		}

		// TODO: Check if migration is using a supported engine
		// if !supportedEngines[change.Engine] {
		// 	return invalidChange(change, "unsupported database engine")
		// }

		if change.Engine == "" {
			return invalidChange(*change, "missing engine")
		}

		if len(strings.Split(change.Changes.Up, " ")) < 5 {
			return invalidChange(*change, "missing (or invalid) migrate instruction")
		}

		if len(strings.Split(change.Changes.Down, " ")) < 3 {
			return invalidChange(*change, "missing (or invalid) rollback instruction")
		}

		version, name, _ := strings.Cut(change.FileName, "_")
		name = strings.Split(name, ".")[0]
		name = strcase.ToCamel(name)

		if change.Version != version {
			return invalidChange(*change, "version mismatch")
		}

		if change.Name != name {
			return invalidChange(*change, "name mismatch")
		}

		visitedNames[change.Name] = true
		visitedVersions[change.Version] = true

		change = change.next
	}

	return true, ""
}

func invalidChange(change Migration, reason string) (bool, string) {
	return false, fmt.Sprintf("Invalid migration: %v.\nReason: %v.\n", change.FileName, reason)
}
