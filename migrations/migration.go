package migrations

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"time"

	"gopkg.in/yaml.v2"
)

type Migration struct {
	Id       int
	FileName string
	Version  string
	Schema   int    `yaml:"schema"`
	Name     string `yaml:"name"`
	Engine   string `yaml:"engine"`
	Changes  struct {
		Up   string `yaml:"up"`
		Down string `yaml:"down"`
	} `yaml:"changes"`
	next     *Migration
	previous *Migration
}

type MigratorVersion struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// MARK: - Implements LinkedList behavior
func (M *Migration) Next() *Migration {
	return M.next
}

// MARK: - Implements Sortable Interface

type Migrations []Migration

func (m Migrations) Len() int {
	return len(m)
}

func (m Migrations) Less(left, right int) bool {
	return m[left].Version < m[right].Version
}

func (m Migrations) Swap(left, right int) {
	m[left], m[right] = m[right], m[left]
}

// MARK: - Implements Hashable

func (m Migrations) ToHash() map[string]Migration {
	hash := map[string]Migration{}

	for _, v := range m {
		hash[v.Version] = v
	}

	return hash
}

// MARK: - Migration loader

func (instance *Migration) Load(file fs.FileInfo, parent string, pattern *regexp.Regexp) error {
	path, _ := filepath.Abs(fmt.Sprintf("%v/%v", parent, file.Name()))

	contents, _ := ioutil.ReadFile(path)

	err := yaml.Unmarshal(contents, &instance)

	match := pattern.FindStringSubmatch(file.Name())

	instance.FileName = file.Name()
	instance.Version = match[pattern.SubexpIndex("Version")]

	if err != nil {
		return err
	}

	return nil
}
