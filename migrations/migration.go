package migrations

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Migration struct {
	Version int    `yaml:"version"`
	Name    string `yaml:"name"`
	Adapter string `yaml:"adapter"`
	Changes struct {
		Up   string `yaml:"up"`
		Down string `yaml:"down"`
	} `yaml:"changes"`
}

func (migration *Migration) Load(file fs.FileInfo, parent string) error {
	path, _ := filepath.Abs(fmt.Sprintf("%v/%v", parent, file.Name()))

	contents, _ := ioutil.ReadFile(path)

	err := yaml.Unmarshal(contents, &migration)

	if err != nil {
		return err
	}

	return nil
}
