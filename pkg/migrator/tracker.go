package migrator

import (
	"context"

	"github.com/oleoneto/dm/pkg/ds"
)

type MigrationsTrackerProtocol interface {
	// StartTracking - Prepares database for migration tracking
	StartTracking(context.Context) error

	// StopTracking - Stops tracking database migrations
	StopTracking(context.Context) error

	// Version - Return the version of the last applied migration.
	Version(context.Context) string

	// IsUpToDate - Indicator of whether migrations are current or up-to-date
	IsUpToDate(context.Context) bool

	// IsTracked - Indicator of whether the database is being managed by this tool
	IsTracked(context.Context) bool

	// IsEmpty - Indicator of whether the database has any migrations
	IsEmpty(context.Context) bool

	// AppliedMigrations - Returns all applied migrations
	AppliedMigrations(context.Context) *ds.Queue[Migration]

	// PendingMigrations - Returns all non-applied migrations.
	PendingMigrations(context.Context) *ds.Queue[Migration]
}
