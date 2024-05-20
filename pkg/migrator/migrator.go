package migrator

import (
	"context"

	"github.com/oleoneto/dm/pkg/ds"
)

type MigratorProtocol interface {
	Validate(context.Context, ds.Queue[Migration]) error
	Apply(context.Context, ds.Queue[Migration]) error
	Revert(context.Context, ds.Queue[Migration]) error
}
