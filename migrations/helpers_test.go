package migrations

import (
	"io/fs"
	"regexp"
	"testing"
	"time"
)

var (
	pattern = *regexp.MustCompile(`(?P<Version>^\d{14})_(?P<Name>[aA-zZ]+).ya?ml$`)

	list = MigrationList{}
)

func defaultMigrationList() MigrationList {
	res := MigrationList{}

	res.Insert(&Migration{
		Version:  "20221231054530",
		Engine:   "postgresql",
		Name:     "CreateUsers",
		FileName: "20221231054530_create_users.yaml",
		Changes: Changes{
			Up:   "CREATE TABLE users (id SERIAL, username VARCHAR UNIQUE NOT NULL);",
			Down: "DROP TABLE users;",
		},
	})

	res.Insert(&Migration{
		Version:  "20221231054531",
		Engine:   "postgresql",
		Name:     "CreateArticles",
		FileName: "20221231054531_create_articles.yaml",
		Changes: Changes{
			Up:   "CREATE TABLE articles (id SERIAL, title VARCHAR NOT NULL);",
			Down: "DROP TABLE articles;",
		},
	})

	res.Insert(&Migration{
		Version:  "20221231054532",
		Engine:   "postgresql",
		Name:     "CreateComments",
		FileName: "20221231054532_create_comments.yaml",
		Changes: Changes{
			Up:   "CREATE TABLE comments (id SERIAL, content TEXT NOT NULL);",
			Down: "DROP TABLE comments;",
		},
	})

	return res
}

func defaultMigrationFiles() []fs.FileInfo {
	res := []fs.FileInfo{
		MockFile{size: 1024, modTime: time.Now(), isDir: false, name: "20220420120000_create_users.yml"},
		MockFile{size: 1024, modTime: time.Now(), isDir: false, name: "20220421010000_create_articles.yml"},
		MockFile{size: 1024, modTime: time.Now(), isDir: false, name: "20220423010000_create_comments.yaml"},
	}

	return res
}

// MARK: Validations

func TestValidateEmptyList(t *testing.T) {

	valid, _ := Validate(list)

	if !valid {
		t.Fatalf(`want validate == true (valid), but got %v`, valid)
	}
}

func TestValidate(t *testing.T) {
	valid, _ := Validate(defaultMigrationList())

	if !valid {
		t.Fatalf(`want validate == true (valid), but got %v`, valid)
	}
}

func TestValidateDuplicateVersion(t *testing.T) {
	list = defaultMigrationList()

	// Duplicate Version
	list.Insert(&Migration{
		Version:  "20221231054532",
		Engine:   "postgresql",
		Name:     "CreateLikes",
		FileName: "20221231054532_create_likes.yaml",
		Changes: Changes{
			Up:   "CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);",
			Down: "DROP TABLE likes;",
		},
	})

	valid, _ := Validate(list)

	if valid {
		t.Fatalf(`want validate == false (invalid), but got %v`, valid)
	}
}

func TestValidateDuplicateName(t *testing.T) {
	list = defaultMigrationList()

	// Duplicate Name
	list.Insert(&Migration{
		Version:  "20221231054540",
		Engine:   "postgresql",
		Name:     "CreateComments",
		FileName: "20221231054540_create_comments.yaml",
		Changes: Changes{
			Up:   "CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);",
			Down: "DROP TABLE likes;",
		},
	})

	valid, _ := Validate(list)

	if valid {
		t.Fatalf(`want validate == false (invalid), but got %v`, valid)
	}
}

func TestValidateMismatchedName(t *testing.T) {
	list = defaultMigrationList()

	// Name Mismatch
	list.Insert(&Migration{
		Version:  "20221231054540",
		Engine:   "postgresql",
		Name:     "CreateLikes",
		FileName: "20221231054540_create_comments.yaml",
		Changes: Changes{
			Up:   "CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);",
			Down: "DROP TABLE likes;",
		},
	})

	valid, _ := Validate(list)

	if valid {
		t.Fatalf(`want validate == false (invalid), but got %v`, valid)
	}
}

func TestValidateMissingEngine(t *testing.T) {
	list = defaultMigrationList()

	// Missing Engine
	list.Insert(&Migration{
		Version:  "20221231054540",
		Engine:   "",
		Name:     "CreateLikes",
		FileName: "20221231054540_create_likes.yaml",
		Changes: Changes{
			Up:   "CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);",
			Down: "DROP TABLE likes;",
		},
	})

	valid, _ := Validate(list)

	if valid {
		t.Fatalf(`want validate == false (invalid), but got %v`, valid)
	}
}

func TestValidateInvalidUpChange(t *testing.T) {
	list = defaultMigrationList()

	// Invalid Up Change
	list.Insert(&Migration{
		Version:  "20221231054540",
		Engine:   "postgresql",
		Name:     "CreateLikes",
		FileName: "20221231054540_create_likes.yaml",
		Changes: Changes{
			Up:   "CREATE TABLE",
			Down: "DROP TABLE likes;",
		},
	})

	valid, _ := Validate(list)

	if valid {
		t.Fatalf(`want validate == false (invalid), but got %v`, valid)
	}
}

func TestValidateInvalidDownChange(t *testing.T) {
	list = defaultMigrationList()

	// Invalid Down Change
	list.Insert(&Migration{
		Version:  "20221231054540",
		Engine:   "postgresql",
		Name:     "CreateLikes",
		FileName: "20221231054540_create_likes.yaml",
		Changes: Changes{
			Up:   "CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);",
			Down: "DROP TABLES;",
		},
	})

	valid, _ := Validate(list)

	if valid {
		t.Fatalf(`want validate == false (invalid), but got %v`, valid)
	}
}

// MARK: File Matcher
func TestMatchingFilesEmpty(t *testing.T) {
	matchedFiles, _ := MatchingFiles("./empty_dir", &pattern)

	if len(matchedFiles) != 0 {
		t.Fatalf(`want len(matches) == 0, but got %v`, len(matchedFiles))
	}
}

func TestMatchingFiles(t *testing.T) {
	matchedFiles, _ := MatchingFiles("../examples", &pattern)

	if len(matchedFiles) != 3 {
		t.Fatalf(`want len(matches) == 3, but got %v`, len(matchedFiles))
	}
}

// MARK: Migration Builder

func TestBuildMigrationsEmpty(t *testing.T) {
	list = BuildMigrations([]fs.FileInfo{}, "migrations", &pattern)

	if list.Size() != 0 {
		t.Fatalf(`want size == 0, but got %v`, list.Size())
	}
}

func TestBuildMigrationsInEmptyDirectory(t *testing.T) {
	list = BuildMigrations(defaultMigrationFiles(), "./empty_dir", &pattern)

	if list.Size() != 0 {
		t.Fatalf(`want size == 0, but got %v`, list.Size())
	}
}

func TestBuildMigrations(t *testing.T) {
	list = BuildMigrations(defaultMigrationFiles(), "../examples", &pattern)

	if list.Size() != 3 {
		t.Fatalf(`want size == 3, but got %v`, list.Size())
	}
}

// MARK: - Supporting Definitions for Testing

type MockFile struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
	sys     any
}

func (f MockFile) Name() string {
	return f.name
}

func (f MockFile) Size() int64 {
	return f.size
}

func (f MockFile) Mode() fs.FileMode {
	return f.mode
}

func (f MockFile) ModTime() time.Time {
	return f.modTime
}

func (f MockFile) IsDir() bool {
	return f.isDir
}

func (f MockFile) Sys() any {
	return f.sys
}
