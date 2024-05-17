package runner

import "github.com/oleoneto/dm/pkg/migrator"

type MigratorProtocol interface {
	migrator.TrackerProtocol
	migrator.MigratorProtocol
}

type Runner struct {
	Migrator MigratorProtocol
}

func NewRunner() *Runner {
	return &Runner{}
}
