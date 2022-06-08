package migrations

import (
	"io/fs"
	"testing"
	"time"
)

var list = MigrationList{}

func defaultMigrationList() MigrationList {
	res := MigrationList{}

	res.Insert(&Migration{
		Version:  "20221231054530129328",
		Engine:   "postgresql",
		Name:     "CreateUsers",
		FileName: "20221231054530129328_create_users.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE users (id SERIAL, username VARCHAR UNIQUE NOT NULL);"},
			Down: []string{"DROP TABLE users;"},
		},
	})

	res.Insert(&Migration{
		Version:  "20221231054531293821",
		Engine:   "postgresql",
		Name:     "CreateArticles",
		FileName: "20221231054531293821_create_articles.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE articles (id SERIAL, title VARCHAR NOT NULL);"},
			Down: []string{"DROP TABLE articles;"},
		},
	})

	res.Insert(&Migration{
		Version:  "20221231054532123874",
		Engine:   "postgresql",
		Name:     "CreateComments",
		FileName: "20221231054532123874_create_comments.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE comments (id SERIAL, content TEXT NOT NULL);"},
			Down: []string{"DROP TABLE comments;"},
		},
	})

	return res
}

func defaultMigrationFiles() []fs.FileInfo {
	res := []fs.FileInfo{
		MigrationFile{size: 1024, modTime: time.Now(), isDir: false, name: "20220504202422742293_create_users.yaml"},
		MigrationFile{size: 1024, modTime: time.Now(), isDir: false, name: "20220504202443251494_create_articles.yaml"},
		MigrationFile{size: 1024, modTime: time.Now(), isDir: false, name: "20220504202502049236_create_comments.yaml"},
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
		Version:  "20221231054532123874",
		Engine:   "postgresql",
		Name:     "CreateLikes",
		FileName: "20221231054532123874_create_likes.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);"},
			Down: []string{"DROP TABLE likes;"},
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
			Up:   []string{"CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);"},
			Down: []string{"DROP TABLE likes;"},
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
			Up:   []string{"CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);"},
			Down: []string{"DROP TABLE likes;"},
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
			Up:   []string{"CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);"},
			Down: []string{"DROP TABLE likes;"},
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
			Up:   []string{"CREATE TABLE"},
			Down: []string{"DROP TABLE likes;"},
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
			Up:   []string{"CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);"},
			Down: []string{"DROP TABLES;"},
		},
	})

	valid, _ := Validate(list)

	if valid {
		t.Fatalf(`want validate == false (invalid), but got %v`, valid)
	}
}

func TestValidateMissingDropWhenCreateIsPresent(t *testing.T) {
	list = defaultMigrationList()

	// Invalid Down Change (no drop instruction)
	list.Insert(&Migration{
		Version:  "20221231054540",
		Engine:   "postgresql",
		Name:     "CreateLikes",
		FileName: "20221231054540_create_likes.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);"},
			Down: []string{"SELECT * FROM likes;"},
		},
	})

	valid, _ := Validate(list)

	if valid {
		t.Fatalf(`want validate == false (invalid), but got %v`, valid)
	}
}

func TestValidateCreateAndDropDifferentTables(t *testing.T) {
	list = defaultMigrationList()

	// Invalid Down Change (drops a different table)
	list.Insert(&Migration{
		Version:  "20221231054540",
		Engine:   "postgresql",
		Name:     "CreateLikes",
		FileName: "20221231054540_create_likes.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE likes (id SERIAL, content_id INT NOT NULL);"},
			Down: []string{"DROP TABLE accounts;"},
		},
	})

	// Invalid Down Change (drops a different table)
	list.Insert(&Migration{
		Version:  "20221231054541",
		Engine:   "postgresql",
		Name:     "CreateReminders",
		FileName: "20221231054540_create_reminders.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE reminders (id SERIAL, content_id INT NOT NULL, time DATETIME);"},
			Down: []string{"DROP TABLE IF EXISTS accounts;"},
		},
	})

	valid, _ := Validate(list)

	if valid {
		t.Fatalf(`want validate == false (invalid), but got %v`, valid)
	}
}

// MARK: File Matcher
func TestMatchingFilesEmpty(t *testing.T) {
	matchedFiles, _ := MatchingFiles("./empty_dir", &FilePattern)

	if len(matchedFiles) != 0 {
		t.Fatalf(`want len(matches) == 0, but got %v`, len(matchedFiles))
	}
}

func TestMatchingFiles(t *testing.T) {
	matchedFiles, _ := MatchingFiles("../examples", &FilePattern)

	if len(matchedFiles) != 6 {
		t.Fatalf(`want len(matches) == 6, but got %v`, len(matchedFiles))
	}
}

// MARK: File Loader

func TestLoadFilesInEmptyDirectory(t *testing.T) {
	files := LoadFiles("migrations", &FilePattern)

	if len(files) != 0 {
		t.Fatalf(`want size == 0, but got %v`, len(files))
	}
}

func TestLoadFiles(t *testing.T) {
	files := LoadFiles("../examples", &FilePattern)

	if len(files) == 0 {
		t.Fatalf(`want size > 0, but got %v`, len(files))
	}
}

// MARK: Migration Builder

func TestBuildMigrationsEmpty(t *testing.T) {
	list = BuildMigrations([]fs.FileInfo{}, "migrations", &FilePattern)

	if list.Size() != 0 {
		t.Fatalf(`want size == 0, but got %v`, list.Size())
	}
}

func TestBuildMigrationsInEmptyDirectory(t *testing.T) {
	list = BuildMigrations(defaultMigrationFiles(), "./empty_dir", &FilePattern)

	if list.Size() != 0 {
		t.Fatalf(`want size == 0, but got %v`, list.Size())
	}
}

func TestBuildMigrations(t *testing.T) {
	list = BuildMigrations(defaultMigrationFiles(), "../examples", &FilePattern)

	if list.Size() != 3 {
		t.Fatalf(`want size == 3, but got %v`, list.Size())
	}
}
