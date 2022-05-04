package migrations

type Store interface {
	// Create - Adds record(s) to the store.
	Create(string, ...interface{}) error

	// Read - Reads record from the store.
	Read(string, *interface{}, ...interface{}) error

	// Delete - Removes record(s) from the store.
	Delete(string, ...interface{}) error
}
