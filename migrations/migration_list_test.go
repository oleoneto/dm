package migrations

import (
	"fmt"
	"testing"
)

func TestGetHead(t *testing.T) {
	list := MigrationList{}

	if list.GetHead() != nil {
		t.Fatalf(`wanted list.GetHead() == nil, got %v`, list.GetHead())
	}

	list.head = &Migration{
		Version:  "20221231054540",
		Engine:   "postgresql",
		Name:     "CreateUsers",
		FileName: "20221231054540_create_users.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE users (id SERIAL)"},
			Down: []string{"DROP TABLE users;"},
		},
	}

	if list.GetHead() == nil {
		t.Fatalf(`wanted list.GetHead() != nil, got %v`, list.GetHead())
	}
}

func TestInsert(t *testing.T) {
	list := MigrationList{}

	if list.GetTail() != nil {
		t.Fatalf(`wanted list.GetTail() == nil, got %v`, list.GetTail())
	}

	list.Insert(&Migration{
		Version:  "20221231054540",
		Engine:   "postgresql",
		Name:     "CreateUsers",
		FileName: "20221231054540_create_users.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE users (id SERIAL)"},
			Down: []string{"DROP TABLE users;"},
		},
	})

	list.Insert(&Migration{
		Version:  "20221231054542",
		Engine:   "postgresql",
		Name:     "CreateArticles",
		FileName: "20221231054542_create_articles.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE articles (id SERIAL, title VARCHAR NOT NULL)"},
			Down: []string{"DROP TABLE articles;"},
		},
	})

	if list.size != 2 {
		t.Fatalf(`wanted list.size == 2, got %v`, list.size)
	}

	if list.GetHead().Name != "CreateUsers" {
		t.Fatalf(`wanted list.GetHead().Name == "CreateUsers", got %v`, list.GetHead().Name)
	}

	if list.GetTail().Name != "CreateArticles" {
		t.Fatalf(`wanted list.GetTail().Name == "CreateArticles", got %v`, list.GetTail().Name)
	}
}

func TestRemoveNonExistent(t *testing.T) {
	list := defaultList()

	list.Remove("CreateNonExistent")

	if list.Size() != 3 {
		t.Fatalf(`wanted list.Size() == 1, but got %v`, list.Size())
	}

	if list.GetHead().Name != "CreateUsers" {
		t.Fatalf(`wanted list.GetHead().Name == "CreateUsers", got %v`, list.GetHead().Name)
	}

	if list.GetTail().Name != "CreateArticles" {
		t.Fatalf(`wanted list.GetTail().Name == "CreateArticles", got %v`, list.GetTail().Name)
	}
}

func TestRemoveHead(t *testing.T) {
	list := defaultList()

	list.Remove("CreateUsers")

	if list.Size() != 2 {
		t.Fatalf(`wanted list.Size() == 2, but got %v`, list.Size())
	}

	if list.GetHead().Name != "CreatePodcasts" {
		t.Fatalf(`wanted list.GetTail().Name == "CreatePodcasts", got %v`, list.GetHead().Name)
	}

	if list.GetTail().Name != "CreateArticles" {
		t.Fatalf(`wanted list.GetTail().Name == "CreateArticles", got %v`, list.GetTail().Name)
	}
}

func TestRemoveTail(t *testing.T) {
	list := defaultList()

	list.Remove("CreateArticles")

	if list.Size() != 2 {
		t.Fatalf(`wanted list.Size() == 2, but got %v`, list.Size())
	}

	if list.GetTail().Name == "CreateArticles" {
		t.Fatalf(`wanted list.GetTail().Name != "CreateArticles", got %v`, list.GetTail().Name)
	}

	if list.GetTail().Name != "CreatePodcasts" {
		t.Fatalf(`wanted list.GetTail().Name == "CreatePodcasts", got %v`, list.GetTail().Name)
	}
}

func TestRemoveMiddleNode(t *testing.T) {
	list := defaultList()

	list.Remove("CreatePodcasts")

	if list.Size() != 2 {
		t.Fatalf(`wanted list.Size() == 1, but got %v`, list.Size())
	}

	if list.GetHead().Name != "CreateUsers" {
		t.Fatalf(`wanted list.GetHead().Name == "CreateUsers", got %v`, list.GetHead().Name)
	}

	if list.GetHead().next.Name != "CreateArticles" {
		t.Fatalf(`wanted list.GetTail().Name == "CreateArticles", got %v`, list.GetHead().next.Name)
	}
}

func TestGetTail(t *testing.T) {
	list := MigrationList{}

	if list.GetTail() != nil {
		t.Fatalf(`wanted list.GetTail() == nil, got %v`, list.GetTail())
	}

	list.Insert(&Migration{
		Version:  "20221231054540",
		Engine:   "postgresql",
		Name:     "CreateUsers",
		FileName: "20221231054540_create_users.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE users (id SERIAL)"},
			Down: []string{"DROP TABLE users;"},
		},
	})

	if list.GetTail() == nil {
		t.Fatalf(`wanted list.GetTail() != nil, got %v`, list.GetTail())
	}
}

func TestIsEmpty(t *testing.T) {
	list := MigrationList{}

	if list.size != 0 {
		t.Fatalf(`wanted list.size == 0, got %v`, list.size)
	}

	if list.Size() != 0 {
		t.Fatalf(`wanted list.Size() == 0, got %v`, list.Size())
	}

	if list.IsEmpty() != true {
		t.Fatalf(`wanted list.IsEmpty() == true, got %v`, list.IsEmpty())
	}
}

func TestSize(t *testing.T) {
	list := MigrationList{}

	if list.size != 0 {
		t.Fatalf(`wanted list.size == 0, got %v`, list.size)
	}

	list.Insert(&Migration{
		Version:  "20221231054540",
		Engine:   "postgresql",
		Name:     "CreateUsers",
		FileName: "20221231054540_create_users.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE users (id SERIAL)"},
			Down: []string{"DROP TABLE users;"},
		},
	})

	if list.size != 1 {
		t.Fatalf(`wanted list.size == 1, got %v`, list.size)
	}
}

func TestFindByVersion(t *testing.T) {
	list := defaultList()

	sequence, found := list.Find("20221231054541")

	if !found {
		t.Fatalf(`wanted found = true, but got %v`, found)
	}

	if sequence.size != 2 {
		t.Fatalf(`wanted sequence.size == 2, but got %v`, sequence.size)
	}
}

func TestFindByName(t *testing.T) {
	list := defaultList()

	// Scenario 1: Search for an existing migration should return a non-empty sequence + true
	sequence, found := list.Find("CreatePodcasts")

	if !found {
		t.Fatalf(`wanted found = true, but got %v`, found)
	}

	if sequence.size != 2 {
		t.Fatalf(`wanted sequence.size == 2, but got %v`, sequence.size)
	}

	// Scenario 2: Search for a non-existing migration should return an empty sequence + false
	sequence, found = list.Find("CreateUnknown")

	if found {
		t.Errorf(`wanted found = false, but got %v`, found)
	}

	if sequence.size != 0 {
		t.Errorf(`wanted sequence.size == 0, but got %v`, sequence.size)
	}

	// Scenario 3: Search in an empty list should return an empty sequence + false
	list = MigrationList{}

	sequence, found = list.Find("CreatePodcasts")

	if found {
		t.Errorf(`wanted found = false, but got %v`, found)
	}

	if sequence.size != 0 {
		t.Errorf(`wanted sequence.size == 0, but got %v`, sequence.size)
	}
}

func TestFromMap(t *testing.T) {
	hash := defaultMap()
	list := MigrationList{}

	list.FromMap(hash)

	if len(hash) != list.size {
		t.Fatalf(`wanted len(hash) && list.size == true, but got %v`, list.size == len(hash))
	}
}

func TestToMap(t *testing.T) {
	list := defaultList()
	hash := list.ToMap()

	if len(hash) != list.size {
		t.Fatalf(`wanted len(hash) && list.size == true, but got %v`, list.size == len(hash))
	}
}

func TestToSlice(t *testing.T) {
	list := defaultList()
	slice := list.ToSlice()

	if len(slice) != list.size {
		t.Fatalf(`wanted len(slice) && list.size == true, but got %v`, list.size == len(slice))
	}
}

func TestReverse(t *testing.T) {
	// Scenario 1: An empty list
	list := MigrationList{}

	list.Reverse()

	if list.head != list.tail && list.head != nil {
		t.Fatalf(`expected an empty list`)
	}

	// Scenario 2: A non-empty list
	list = defaultList()

	list.Reverse()

	if list.GetTail().Name != "CreateUsers" {
		t.Fatalf(`wanted list.GetTail().Name == "CreateUsers", got %v`, list.GetTail().Name)
	}

	if list.GetHead().Name != "CreateArticles" {
		t.Fatalf(`wanted list.GetHead().Name == "CreateArticles", got %v`, list.GetHead().Name)
	}
}

func TestMigrationListDescription(t *testing.T) {
	// Scenario 1: An empty list
	list := MigrationList{}

	if list.Description() != "No migrations in list" {
		t.Errorf(`expected a different description`)
	}

	// Scenario 2: A non-empty list
	list = defaultList()
	description := fmt.Sprintf("%v migrations", list.size)

	if list.Description() != description {
		t.Errorf(`expected a different description`)
	}
}

// MARK: - Helpers

func defaultList() MigrationList {
	list := MigrationList{}

	list.Insert(&Migration{
		Version:  "20221231054540",
		Engine:   "postgresql",
		Name:     "CreateUsers",
		FileName: "20221231054540_create_users.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE users (id SERIAL)"},
			Down: []string{"DROP TABLE users;"},
		},
	})

	list.Insert(&Migration{
		Version:  "20221231054541",
		Engine:   "postgresql",
		Name:     "CreatePodcasts",
		FileName: "20221231054541_create_podcasts.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE podcasts (id SERIAL, title VARCHAR NOT NULL)"},
			Down: []string{"DROP TABLE podcasts;"},
		},
	})

	list.Insert(&Migration{
		Version:  "20221231054542",
		Engine:   "postgresql",
		Name:     "CreateArticles",
		FileName: "20221231054542_create_articles.yaml",
		Changes: Changes{
			Up:   []string{"CREATE TABLE articles (id SERIAL, title VARCHAR NOT NULL)"},
			Down: []string{"DROP TABLE articles;"},
		},
	})

	return list
}

func defaultMap() map[string]Migration {
	hash := map[string]Migration{
		"20221231054540": {
			Version:  "20221231054540",
			Engine:   "postgresql",
			Name:     "CreateUsers",
			FileName: "20221231054540_create_users.yaml",
			Changes: Changes{
				Up:   []string{"CREATE TABLE users (id SERIAL)"},
				Down: []string{"DROP TABLE users;"},
			},
		},
		"20221231054541": {
			Version:  "20221231054541",
			Engine:   "postgresql",
			Name:     "CreatePodcasts",
			FileName: "20221231054541_create_podcasts.yaml",
			Changes: Changes{
				Up:   []string{"CREATE TABLE podcasts (id SERIAL, title VARCHAR NOT NULL)"},
				Down: []string{"DROP TABLE podcasts;"},
			},
		},
		"20221231054542": {
			Version:  "20221231054542",
			Engine:   "postgresql",
			Name:     "CreateArticles",
			FileName: "20221231054542_create_articles.yaml",
			Changes: Changes{
				Up:   []string{"CREATE TABLE articles (id SERIAL, title VARCHAR NOT NULL)"},
				Down: []string{"DROP TABLE articles;"},
			},
		},
	}

	return hash
}
