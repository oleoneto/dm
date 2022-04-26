package migrations

import (
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

func CurrentFilepath() string {
	_, name, _, _ := runtime.Caller(1)
	path := path.Join(path.Dir(name), ".")

	return path
}
