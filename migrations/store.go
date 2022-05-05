package migrations

import "fmt"

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

type ExampleStore struct{}

func (s ExampleStore) Create(query string, options ...interface{}) error {
	return nil
}

func (s ExampleStore) Read(string, interface{}, ...interface{}) error {
	fmt.Println("Read from ExampleStore")
	return nil
}

func (s ExampleStore) Delete(query string, options ...interface{}) error {
	return nil
}

func (s ExampleStore) Connect(string) error {
	return nil
}

func (s ExampleStore) Disconnect(string) error {
	return nil
}

func (s ExampleStore) DatabaseURL() string {
	return "example://database.dev"
}

func (s ExampleStore) Name() string {
	return "Example"
}
