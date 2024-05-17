package migrator

import (
	"github.com/oleoneto/dm/pkg/ds"
)

type TrackerProtocol interface {
	// StartTracking - Prepares database for migration tracking
	StartTracking() error

	// StopTracking - Stops tracking database migrations
	StopTracking() error

	// Version - Return the version of the last applied migration. The returned boolean should indicate if the database is being tracked
	Version() (string, bool)

	// IsUpToDate - Indicator of whether migrations are current or up-to-date
	IsUpToDate(changes ds.Queue[Migration]) bool

	// IsTracked - Indicator of whether the database is being managed by this tool
	IsTracked() bool

	// IsEmpty - Indicator of whether the database has any migrations
	IsEmpty() bool

	// AppliedMigrations - Returns all applied/saved migrations
	AppliedMigrations() ds.Queue[Migration]

	// PendingMigrations - Returns all non-applied/saved migrations.
	PendingMigrations() ds.Queue[Migration]
}
