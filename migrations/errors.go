package migrations

type EngineError struct{}

type ValidationError struct{}

func (error EngineError) Error() string {
	return "engine returned an error"
}

func (error ValidationError) Error() string {
	return "validation error"
}
