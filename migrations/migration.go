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

type MigrationFile struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
	sys     any
}

type Migration struct {
	Id       int     `yaml:"-"`
	FileName string  `yaml:"-"`
	Version  string  `yaml:"-"`
	Schema   int     `yaml:"schema"`
	Name     string  `yaml:"name"`
	Engine   string  `yaml:"engine"`
	Changes  Changes `yaml:"changes"`
	next     *Migration
	previous *Migration
}

type Changes struct {
	Up   string `yaml:"up"`
	Down string `yaml:"down"`
}

type MigratorVersion struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type TableSchema struct {
	TableSchema string `json:"table_schema" db:"table_schema"`
	TableName   string `json:"table_name" db:"table_name"`
	TableType   string `json:"table_type" db:"table_type"`
}

func (M *Migration) Description() string {
	return fmt.Sprintf("Version: %v (%v)", M.Version, M.Name)
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
	path := filepath.Join(parent, file.Name())

	path, _ = filepath.Abs(path)

	contents, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(contents, &instance)

	if err != nil {
		return err
	}

	match := pattern.FindStringSubmatch(file.Name())

	instance.FileName = file.Name()
	instance.Version = match[pattern.SubexpIndex("Version")]

	return nil
}

// MARK: - Migration file (implements the fs.FileInfo interface)

func (f MigrationFile) Name() string {
	return f.name
}

func (f MigrationFile) Size() int64 {
	return f.size
}

func (f MigrationFile) Mode() fs.FileMode {
	return f.mode
}

func (f MigrationFile) ModTime() time.Time {
	return f.modTime
}

func (f MigrationFile) IsDir() bool {
	return f.isDir
}

func (f MigrationFile) Sys() any {
	return f.sys
}
