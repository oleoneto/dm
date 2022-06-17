package migrations

import (
	"fmt"
	"io/fs"
	"testing"
	"time"
)

func TestMigrationNodeNext(t *testing.T) {
	list := defaultList()

	if list.head.next != list.head.Next() {
		t.Errorf(`expected head.next to be head.Next()`)
	}
}

// MARK: - Sortable

func TestMigrationsLen(t *testing.T) {
	// Scenario 1: An empty migrations slice
	migrations := Migrations{}

	if migrations.Len() != 0 {
		t.Errorf(`expected length to be 0, but got %v`, migrations.Len())
	}

	// Scenario 2: A non-empty migrations slice
	migrations = Migrations{
		Migration{
			Version:  "20221231054540",
			Engine:   "postgresql",
			Name:     "CreateLikes",
			FileName: "20221231054540_create_likes.yaml",
			Changes: Changes{
				Up:   []string{"CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);"},
				Down: []string{"DROP TABLE likes;"},
			},
		},
	}

	if migrations.Len() != 1 {
		t.Errorf(`expected length to be 1, but got %v`, migrations.Len())
	}
}

func TestMigrationsLess(t *testing.T) {
	migrations := Migrations{
		Migration{
			Version:  "20221231054540",
			Engine:   "postgresql",
			Name:     "CreateLikes",
			FileName: "20221231054540_create_likes.yaml",
			Changes: Changes{
				Up:   []string{"CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);"},
				Down: []string{"DROP TABLE likes;"},
			},
		},
		Migration{
			Version:  "20221231054581",
			Engine:   "postgresql",
			Name:     "CreateComments",
			FileName: "20221231054581_create_comments.yaml",
			Changes: Changes{
				Up:   []string{"CREATE TABLE comments (id SERIAL, content_id INT NOT NULL);"},
				Down: []string{"DROP TABLE comments;"},
			},
		},
	}

	if migrations.Less(1, 0) {
		t.Errorf(`expected %v < %v`, migrations[0], migrations[1])
	}
}

func TestMigrationsSwap(t *testing.T) {
	migrations := Migrations{
		Migration{
			Version:  "20221231054540",
			Engine:   "postgresql",
			Name:     "CreateLikes",
			FileName: "20221231054540_create_likes.yaml",
			Changes: Changes{
				Up:   []string{"CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);"},
				Down: []string{"DROP TABLE likes;"},
			},
		},
		Migration{
			Version:  "20221231054581",
			Engine:   "postgresql",
			Name:     "CreateComments",
			FileName: "20221231054581_create_comments.yaml",
			Changes: Changes{
				Up:   []string{"CREATE TABLE comments (id SERIAL, content_id INT NOT NULL);"},
				Down: []string{"DROP TABLE comments;"},
			},
		},
	}

	migrations.Swap(0, 1)

	if migrations[0].Name != "CreateComments" {
		t.Errorf(`expected [0]%v, [1]%v`, migrations[1].Name, migrations[0].Name)
	}
}

// MARK: - Formattable

func TestMigrationVersionDescription(t *testing.T) {
	version := MigratorVersion{}
	description := fmt.Sprintf("%v (%v).\nApplied at: %v", version.Version, version.Name, version.CreatedAt)

	if version.Description() != description {
		t.Errorf(`expected a different description %v, but got %v`, version.Description(), description)
	}
}

func TestMigrationsVersionDescription(t *testing.T) {
	// Scenario 1: An empty migrations slice
	migrations := Migrations{}

	if migrations.Description() != "No migrations" {
		t.Errorf(`expected a different description, got %v`, migrations.Description())
	}

	// Scenario 2: A non-empty migrations slice
	migrations = Migrations{
		Migration{
			Version:  "20221231054540",
			Engine:   "postgresql",
			Name:     "CreateLikes",
			FileName: "20221231054540_create_likes.yaml",
			Changes: Changes{
				Up:   []string{"CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);"},
				Down: []string{"DROP TABLE likes;"},
			},
		},
		Migration{
			Version:  "20221231054581",
			Engine:   "postgresql",
			Name:     "CreateComments",
			FileName: "20221231054581_create_comments.yaml",
			Changes: Changes{
				Up:   []string{"CREATE TABLE comments (id SERIAL, content_id INT NOT NULL);"},
				Down: []string{"DROP TABLE comments;"},
			},
		},
	}

	description := fmt.Sprintln("Version: 20221231054540 (CreateLikes)")
	description += fmt.Sprintln("Version: 20221231054581 (CreateComments)")

	if migrations.Description() != description {
		t.Errorf(`expected a different description %v, got %v`, migrations.Description(), description)
	}
}

// MARK: - Hashable

func TestMigrationToHash(t *testing.T) {
	// Scenario 1: An empty migrations slice
	migrations := Migrations{}
	hashed := migrations.ToHash()

	if migrations.Len() != len(hashed) && migrations.Len() != 0 {
		t.Errorf(`expected both map and slice to be empty`)
	}

	// Scenario 2: A non-empty migrations slice
	migrations = Migrations{
		Migration{
			Version:  "20221231054540",
			Engine:   "postgresql",
			Name:     "CreateLikes",
			FileName: "20221231054540_create_likes.yaml",
			Changes: Changes{
				Up:   []string{"CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);"},
				Down: []string{"DROP TABLE likes;"},
			},
		},
	}

	hashed = migrations.ToHash()

	if migrations.Len() != len(hashed) && migrations.Len() != 1 {
		t.Errorf(`expected both map and slice to be of size = 1`)
	}

	_, exists := hashed["20221231054540"]

	if !exists {
		t.Errorf(`expected migration to be in hash map`)
	}
}

// MARK: - fs.FileInfo

func TextMigrationFileIsDir(t *testing.T) {
	migrationFile := testMigrationFile()

	if migrationFile.IsDir() != migrationFile.isDir {
		t.Errorf(`expected method and attribute to match`)
	}
}

func TextMigrationFileName(t *testing.T) {
	migrationFile := testMigrationFile()

	if migrationFile.Name() != migrationFile.name {
		t.Errorf(`expected method and attribute to match`)
	}
}

// MARK: - Helpers

func testMigrationFile() MigrationFile {
	return MigrationFile{
		name:    "example.yaml",
		isDir:   false,
		mode:    fs.ModeAppend,
		modTime: time.Now(),
		sys:     nil,
		size:    42,
	}
}
