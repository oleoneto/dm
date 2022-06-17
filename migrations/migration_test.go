package migrations

import "testing"

// MARK: - Sortable

func TestMigrationsLen(t *testing.T) {

}

func TestMigrationsLess(t *testing.T) {

}

func TestMigrationsSwap(t *testing.T) {

}

// MARK: - Formattable

func TestMigrationDescription(t *testing.T) {

}

func TestMigrationVersionDescription(t *testing.T) {

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
func TestMigrationFileName(t *testing.T) {}

func TestMigrationFileSize(t *testing.T) {}

func TestMigrationFileMode(t *testing.T) {}

func TestMigrationFileModTime(t *testing.T) {}

func TestMigrationFileIsDir(t *testing.T) {}

func TestMigrationFileSys(t *testing.T) {}
