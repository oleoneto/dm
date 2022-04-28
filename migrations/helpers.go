package migrations

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"runtime"
)

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

func BuildMigrations(dir string, pattern *regexp.Regexp) Migrations {
	var changes Migrations

	files, _ := MatchingFiles(dir, pattern)

	for _, file := range files {
		var mg Migration

		mg.Load(file, dir)

		changes = append(changes, mg)
	}

	return changes
}

func ListFiles(dir string, pattern *regexp.Regexp) error {
	files, err := MatchingFiles(dir, pattern)

	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}

	return nil
}

func CurrentFilepath() string {
	_, name, _, _ := runtime.Caller(1)
	path := path.Join(path.Dir(name), ".")

	return path
}
