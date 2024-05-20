package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/engines"
	"github.com/oleoneto/dm/pkg/migrator"
)

type Runner struct {
	migrator       migrator.MigratorProtocol
	engine         engines.SqlEngine
	trackerOptions TrackerOptions
}

type TrackerOptions struct {
	Schema string
	Table  string
}

var _ migrator.MigrationsTrackerProtocol = (*Runner)(nil)

func NewRunner(m migrator.MigratorProtocol, e engines.SqlEngine, trackerOptions TrackerOptions) *Runner {
	if trackerOptions.Table == "" {
		trackerOptions.Table = "_migrations"
	}

	if trackerOptions.Schema == "" {
		trackerOptions.Schema = "public"
	}

	if e == nil {
		panic("sql engine not set")
	}

	return &Runner{
		migrator: m,
		engine:   e,
	}
}

func (r *Runner) AppliedMigrations(ctx context.Context) ds.Queue[migrator.Migration]

func (r *Runner) PendingMigrations(ctx context.Context) ds.Queue[migrator.Migration]

// IsEmpty - Return `true` if the tracked table is found and has no rows.
func (r *Runner) IsEmpty(ctx context.Context) bool {
	var count int

	tracked := r.IsTracked(ctx)

	if !tracked {
		return true
	}

	rows, err := r.engine.Query(ctx, fmt.Sprintf(`SELECT COUNT(id) FROM %v`, r.trackerOptions.Table))
	if err != nil {
		return false
	}

	if !rows.Next() {
		return false
	}

	rows.Scan(&count)

	return count == 0
}

// IsTracker - Returns `true` if the tracked table is found in the database.
func (r *Runner) IsTracked(ctx context.Context) bool {
	rows, err := r.engine.Query(
		ctx,
		fmt.Sprintf(`
		SELECT 
			TABLE_SCHEMA, 
			TABLE_NAME,
			TABLE_TYPE
		FROM 
			information_schema.TABLES 
		WHERE 
			TABLE_TYPE LIKE 'BASE TABLE'
			AND TABLE_NAME = '%v'
		LIMIT 1`, r.trackerOptions.Table),
	)
	if err != nil {
		return false
	}

	if !rows.Next() {
		return false
	}

	type tableSchema struct {
		TableSchema string `json:"table_schema" db:"table_schema"`
		TableName   string `json:"table_name" db:"table_name"`
		TableType   string `json:"table_type" db:"table_type"`
	}

	var schema tableSchema
	rows.Scan(&schema.TableSchema, &schema.TableName, &schema.TableType)

	return schema.TableName != ""
}

// IsUpToDate - Returns `true` if the latest migration file matches the last applied migration.
func (r *Runner) IsUpToDate(ctx context.Context, migrations ds.Queue[migrator.Migration]) bool {
	tracked := r.IsTracked(ctx)

	if !tracked {
		return false
	}

	latest := migrations.GetBack()
	version := r.Version(ctx)

	return version == latest.Version
}

// Version - Returns the version of the last applied migration.
func (r *Runner) Version(ctx context.Context) string {
	type migrationVersion struct {
		Id        int       `json:"id"`
		Name      string    `json:"name"`
		Version   string    `json:"version"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
	}

	var version migrationVersion

	rows, err := r.engine.QueryRow(
		ctx,
		fmt.Sprintf(
			"SELECT id, name, version, created_at FROM %v ORDER BY id DESC LIMIT 1",
			r.trackerOptions.Table,
		),
	)
	if err != nil {
		return ""
	}

	rows.Scan(&version.Id, &version.Name, &version.Version, &version.CreatedAt)

	return version.Version
}

// StartTracking - Creates tracking table.
func (r *Runner) StartTracking(ctx context.Context) error {
	// Nothing to do here...
	if r.IsTracked(ctx) {
		return nil
	}

	_, err := r.engine.Exec(
		ctx,
		fmt.Sprintf(`
		CREATE TABLE %v (
			id SERIAL,
			version varchar UNIQUE NOT NULL,
			name varchar UNIQUE NOT NULL,
			created_at timestamp NOT NULL DEFAULT now(),
			PRIMARY KEY(id)
		)`, r.trackerOptions.Table),
	)

	if err != nil {
		return err
	}

	return nil
}

// StopTracking - Removes tracking table.
func (r *Runner) StopTracking(ctx context.Context) error {
	// Nothing to do here...
	if !r.IsTracked(ctx) {
		return nil
	}

	rows, err := r.engine.Exec(ctx, fmt.Sprintf(`DROP TABLE %v`, r.trackerOptions.Table))
	if err != nil {
		return err
	}

	if n, err := rows.RowsAffected(); n <= 0 {
		return err
	}

	return nil
}
