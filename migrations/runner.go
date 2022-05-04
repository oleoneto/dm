package migrations

/*
Runner:
	Responsible for running and reverting migrations.
	The migration and rollback algorithm is self-contained within this type.

	For more flexibility, creations, reads, and deletions are responsibilities of
	the underlying store type. You can inject a type that conforms to the `Store` interface,
	and the `Runner` will call the store's appropriate methods when needed.
*/
type Runner struct{}
