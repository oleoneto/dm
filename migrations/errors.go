package migrations

type EngineError struct{}

type MigrationError struct{}

type RollbackError struct{}

type ValidationError struct{}

func (error EngineError) Error() string {
	return "engine returned an error"
}

func (error MigrationError) Error() string {
	return "migrations failed"
}

func (error RollbackError) Error() string {
	return "rollback failed"
}

func (error ValidationError) Error() string {
	return "validation error"
}
