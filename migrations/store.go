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
