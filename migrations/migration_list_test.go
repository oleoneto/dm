package migrations

import "testing"

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
			Up:   "CREATE TABLE users (id SERIAL)",
			Down: "DROP TABLE users;",
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
			Up:   "CREATE TABLE users (id SERIAL)",
			Down: "DROP TABLE users;",
		},
	})

	list.Insert(&Migration{
		Version:  "20221231054542",
		Engine:   "postgresql",
		Name:     "CreateArticles",
		FileName: "20221231054542_create_articles.yaml",
		Changes: Changes{
			Up:   "CREATE TABLE articles (id SERIAL, title VARCHAR NOT NULL)",
			Down: "DROP TABLE articles;",
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
			Up:   "CREATE TABLE users (id SERIAL)",
			Down: "DROP TABLE users;",
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
			Up:   "CREATE TABLE users (id SERIAL)",
			Down: "DROP TABLE users;",
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

	sequence, found := list.Find("CreatePodcasts")

	if !found {
		t.Fatalf(`wanted found = true, but got %v`, found)
	}

	if sequence.size != 2 {
		t.Fatalf(`wanted sequence.size == 2, but got %v`, sequence.size)
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
	list := defaultList()

	list.Reverse()

	if list.GetTail().Name != "CreateUsers" {
		t.Fatalf(`wanted list.GetTail().Name == "CreateUsers", got %v`, list.GetTail().Name)
	}

	if list.GetHead().Name != "CreateArticles" {
		t.Fatalf(`wanted list.GetHead().Name == "CreateArticles", got %v`, list.GetHead().Name)
	}
}

// MARK: -

func defaultList() MigrationList {
	list := MigrationList{}

	list.Insert(&Migration{
		Version:  "20221231054540",
		Engine:   "postgresql",
		Name:     "CreateUsers",
		FileName: "20221231054540_create_users.yaml",
		Changes: Changes{
			Up:   "CREATE TABLE users (id SERIAL)",
			Down: "DROP TABLE users;",
		},
	})

	list.Insert(&Migration{
		Version:  "20221231054541",
		Engine:   "postgresql",
		Name:     "CreatePodcasts",
		FileName: "20221231054541_create_podcasts.yaml",
		Changes: Changes{
			Up:   "CREATE TABLE podcasts (id SERIAL, title VARCHAR NOT NULL)",
			Down: "DROP TABLE podcasts;",
		},
	})

	list.Insert(&Migration{
		Version:  "20221231054542",
		Engine:   "postgresql",
		Name:     "CreateArticles",
		FileName: "20221231054542_create_articles.yaml",
		Changes: Changes{
			Up:   "CREATE TABLE articles (id SERIAL, title VARCHAR NOT NULL)",
			Down: "DROP TABLE articles;",
		},
	})

	return list
}

func defaultMap() map[string]Migration {
	hash := map[string]Migration{
		"20221231054540": Migration{
			Version:  "20221231054540",
			Engine:   "postgresql",
			Name:     "CreateUsers",
			FileName: "20221231054540_create_users.yaml",
			Changes: Changes{
				Up:   "CREATE TABLE users (id SERIAL)",
				Down: "DROP TABLE users;",
			},
		},
		"20221231054541": Migration{
			Version:  "20221231054541",
			Engine:   "postgresql",
			Name:     "CreatePodcasts",
			FileName: "20221231054541_create_podcasts.yaml",
			Changes: Changes{
				Up:   "CREATE TABLE podcasts (id SERIAL, title VARCHAR NOT NULL)",
				Down: "DROP TABLE podcasts;",
			},
		},
		"20221231054542": Migration{
			Version:  "20221231054542",
			Engine:   "postgresql",
			Name:     "CreateArticles",
			FileName: "20221231054542_create_articles.yaml",
			Changes: Changes{
				Up:   "CREATE TABLE articles (id SERIAL, title VARCHAR NOT NULL)",
				Down: "DROP TABLE articles;",
			},
		},
	}

	return hash
}
