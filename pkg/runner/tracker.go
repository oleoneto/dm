package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/helpers"
	"github.com/oleoneto/dm/pkg/migrator"
	log "github.com/sirupsen/logrus"
)

var _ migrator.MigrationsTrackerProtocol = (*Runner)(nil)

// AppliedMigrations - Returns a list of migrations recorded in the database.
func (r *Runner) AppliedMigrations(ctx context.Context) (*ds.Queue[migrator.Migration], error) {
	if !r.IsTracked(ctx) {
		return ds.NewQueue[migrator.Migration](), nil
	}

	type migrationVersion struct {
		Id        int       `json:"id"`
		Name      string    `json:"name"`
		Version   string    `json:"version"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
	}

	var versions []migrationVersion

	rows, err := r.engine.QueryContext(
		ctx,
		fmt.Sprintf(
			"SELECT id, name, version, created_at FROM %v.%v ORDER BY id DESC",
			r.trackerOptions.Schema, r.trackerOptions.Table,
		),
	)
	if err != nil {
		log.Error(ctx, err.Error(), helpers.GetCurrentFuncName())
		return nil, err
	}

	for rows.Next() {
		var v migrationVersion
		if err := rows.Scan(&v.Id, &v.Name, &v.Version, &v.CreatedAt); err != nil {
			log.Error(ctx, err.Error(), helpers.GetCurrentFuncName())
			return nil, err
		}
		versions = append(versions, v)
	}

	var queue ds.Queue[migrator.Migration]
	for _, item := range versions {
		queue.Enqueue(migrator.Migration{Id: item.Id, Version: item.Version, Name: item.Name})
	}

	return &queue, nil
}

// PendingMigrations - Compares the migrations found in the filesystem and those not recorded in the database.
func (r *Runner) PendingMigrations(ctx context.Context) (*ds.Queue[migrator.Migration], error) {
	if err := r.LoadMigrations(ctx); err != nil {
		return nil, err
	}

	applied, err := r.AppliedMigrations(ctx)
	if err != nil {
		log.Error(ctx, err.Error(), helpers.GetCurrentFuncName())
		return &r.migrations, err
	}

	var pending = ds.Queue[migrator.Migration]{}

	item := r.migrations.Dequeue()
	for item != nil {
		v := applied.Find(func(m migrator.Migration) bool { return m.Version == item.Version })
		if v == nil {
			pending.Enqueue(*item)
		}

		item = r.migrations.Dequeue()
	}

	return &pending, nil
}

// IsEmpty - Return `true` if the tracked table is found and has no rows.
func (r *Runner) IsEmpty(ctx context.Context) bool {
	var count int

	tracked := r.IsTracked(ctx)

	if !tracked {
		return true
	}

	rows, err := r.engine.QueryContext(ctx, fmt.Sprintf(`SELECT COUNT(id) FROM %v`, r.trackerOptions.Table))
	if err != nil {
		log.Error(ctx, err.Error(), helpers.GetCurrentFuncName())
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
	query := fmt.Sprintf(`
	SELECT
		TABLE_SCHEMA,
		TABLE_NAME,
		TABLE_TYPE
	FROM
		information_schema.TABLES
	WHERE
		TABLE_TYPE LIKE 'BASE TABLE'
		AND TABLE_NAME = '%v'
	LIMIT 1`, r.trackerOptions.Table,
	)

	rows, err := r.engine.QueryContext(ctx, query)
	if err != nil {
		log.Error(ctx, err.Error(), helpers.GetCurrentFuncName())
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
func (r *Runner) IsUpToDate(ctx context.Context) bool {
	tracked := r.IsTracked(ctx)

	if !tracked {
		return false
	}

	latest := r.migrations.GetBack()
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

	query := fmt.Sprintf(
		"SELECT id, name, version, created_at FROM %v ORDER BY id DESC LIMIT 1",
		r.trackerOptions.Table,
	)

	var version migrationVersion

	row := r.engine.QueryRowContext(ctx, query)
	if row.Err() != nil {
		return ""
	}

	row.Scan(&version.Id, &version.Name, &version.Version, &version.CreatedAt)

	return version.Version
}

// StartTracking - Creates tracking table.
func (r *Runner) StartTracking(ctx context.Context) error {
	// Nothing to do here...
	if r.IsTracked(ctx) {
		return nil
	}

	query := fmt.Sprintf(`
	CREATE TABLE %v.%v (
		id SERIAL,
		version varchar UNIQUE NOT NULL,
		name varchar UNIQUE NOT NULL,
		created_at timestamp NOT NULL DEFAULT now(),
		PRIMARY KEY(id)
	)`, r.trackerOptions.Schema, r.trackerOptions.Table,
	)

	if _, err := r.engine.ExecContext(ctx, query); err != nil {
		log.Error(ctx, err.Error(), helpers.GetCurrentFuncName())
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

	query := fmt.Sprintf(`DROP TABLE %v.%v`, r.trackerOptions.Schema, r.trackerOptions.Table)

	_, err := r.engine.ExecContext(ctx, query)
	if err != nil {
		log.Error(ctx, err.Error(), helpers.GetCurrentFuncName())
		return err
	}

	return nil
}

func (r *Runner) LoadMigrations(ctx context.Context) error {
	files := r.fileLoaderFunc(r.trackerOptions.MigrationsDirectory, MigrationFileRegexPattern)

	migrations, err := r.migrationLoaderFunc(files, r.trackerOptions.MigrationsDirectory, r.trackerOptions.FileRegexPattern)
	if err != nil {
		log.Error(ctx, err.Error(), helpers.GetCurrentFuncName())
		return err
	}

	r.migrations = *ds.NewFromSlice(migrations)

	return nil
}
