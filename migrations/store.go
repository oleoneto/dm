package migrations

type DatabaseConnector interface {
	// Connect - Acquire a connection to the database.
	Connect() error

	// Disconnect - Releases all existing database connections.
	Disconnect() error
}

type Store interface {
	// Create - Adds record(s) to the store.
	Create(string, ...interface{}) error

	// Read - Reads record from the store.
	Read(string, interface{}, ...interface{}) error

	// Delete - Removes record(s) from the store.
	Delete(string, ...interface{}) error

	// Name - A string that identifies this store.
	Name() string

	DatabaseURL() string
}

// MARK: - Example Store

type ExampleStore struct {
	storage map[string]interface{}
}

func (s ExampleStore) Create(query string, options ...interface{}) error {
	s.storage[query] = options
	return nil
}

func (s ExampleStore) Read(query string, model interface{}, options ...interface{}) interface{} {
	return s.storage[query]
}

func (s ExampleStore) Delete(query string, options ...interface{}) error {
	delete(s.storage, query)
	return nil
}

func (s ExampleStore) Connect(string) error {
	return nil
}

func (s ExampleStore) Disconnect(string) error {
	return nil
}
