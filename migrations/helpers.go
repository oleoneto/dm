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

// BuildMigrations - Instantiate a list of migrations from the contents of the provided files. Accesses the filesystem.
func BuildMigrations(files []fs.FileInfo, dir string, pattern *regexp.Regexp) MigrationList {
	var migrations MigrationList

	for _, file := range files {
		var mg Migration

		err := mg.Load(file, dir, pattern)

		if err == nil {
			migrations.Insert(&mg)
		}
	}

	return migrations
}

func LoadFiles(dir string, pattern *regexp.Regexp) []fs.FileInfo {
	files, err := MatchingFiles(dir, pattern)

	if err != nil {
		return []fs.FileInfo{}
	}

	return files
}

// Validate - Runs validations on a list of migrations.
func Validate(migrations MigrationList) (bool, string) {
	visitedNames := map[string]bool{}
	visitedVersions := map[string]bool{}

	migration := migrations.head

	for migration != nil {
		if visitedVersions[migration.Version] {
			return invalidMigration(*migration, "duplicate migration version")
		}

		if visitedNames[migration.Name] {
			return invalidMigration(*migration, "duplicate migration name")
		}

		// TODO: Check if migration is using a supported engine
		// if !supportedEngines[migration.Engine] {
		// 	return invalidMigration(migration, "unsupported database engine")
		// }

		if migration.Engine == "" {
			return invalidMigration(*migration, "missing engine")
		}

		if len(strings.Split(migration.Changes.Up, " ")) < 5 {
			return invalidMigration(*migration, "missing (or invalid) migrate instruction")
		}

		if len(strings.Split(migration.Changes.Down, " ")) < 3 {
			return invalidMigration(*migration, "missing (or invalid) rollback instruction")
		}

		version, name, _ := strings.Cut(migration.FileName, "_")
		name = strings.Split(name, ".")[0]
		name = strcase.ToCamel(name)

		if migration.Version != version {
			return invalidMigration(*migration, "version mismatch")
		}

		if migration.Name != name {
			return invalidMigration(*migration, "name mismatch")
		}

		visitedNames[migration.Name] = true
		visitedVersions[migration.Version] = true

		migration = migration.next
	}

	return true, ""
}

func invalidMigration(migration Migration, reason string) (bool, string) {
	return false, fmt.Sprintf("Invalid migration: %v.\nReason: %v.\n", migration.Description(), reason)
}
