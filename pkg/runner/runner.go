package runner

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/engines"
	"github.com/oleoneto/dm/pkg/fsystem"
	"github.com/oleoneto/dm/pkg/migrator"
	"gopkg.in/yaml.v2"
)

var MigrationFileRegexPattern = func() *regexp.Regexp {
	return regexp.MustCompile(`(?P<Version>^\d{20})_(?P<Name>[aA-zZ]+).yaml|sql$`)
}()

type Runner struct {
	engine         engines.SqlEngineProtocol
	loader         fsystem.FileLoaderProtocol
	migrator       migrator.MigratorProtocol
	trackerOptions TrackerOptions
}

type TrackerOptions struct {
	Schema              string
	Table               string
	MigrationsDirectory string
	FileRegexPattern    *regexp.Regexp
}

var _ migrator.MigrationsTrackerProtocol = (*Runner)(nil)
var _ fsystem.MigrationBuilderProtocol = (*Runner)(nil)

func NewRunner(
	e engines.SqlEngineProtocol,
	l fsystem.FileLoaderProtocol,
	m migrator.MigratorProtocol,
	trackerOptions TrackerOptions,
) *Runner {
	if trackerOptions.MigrationsDirectory == "" {
		trackerOptions.MigrationsDirectory = "migrations"
	}

	if trackerOptions.Table == "" {
		trackerOptions.Table = "_migrations"
	}

	if trackerOptions.Schema == "" {
		trackerOptions.Schema = "public"
	}

	if trackerOptions.FileRegexPattern == nil {
		trackerOptions.FileRegexPattern = MigrationFileRegexPattern
	}

	if e == nil {
		panic("sql engine not set")
	}

	return &Runner{
		engine:         e,
		loader:         l,
		migrator:       m,
		trackerOptions: trackerOptions,
	}
}

// AppliedMigrations - Returns a list of migrations recorded in the database.
func (r *Runner) AppliedMigrations(ctx context.Context) *ds.Queue[migrator.Migration] {
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
			"SELECT id, name, version, created_at FROM %v ORDER BY id DESC",
			r.trackerOptions.Table,
		),
	)
	if err != nil {
		return nil
	}

	if !rows.Next() {
		return nil
	}

	rows.Scan(&versions)

	var queue ds.Queue[migrator.Migration]
	for _, item := range versions {
		queue.Enqueue(migrator.Migration{Id: item.Id, Version: item.Version, Name: item.Name})
	}

	return &queue
}

// PendingMigrations - Compares the migrations found in the filesystem and those not recorded in the database.
func (r *Runner) PendingMigrations(ctx context.Context) *ds.Queue[migrator.Migration] {
	files := r.loader.LoadFiles(r.trackerOptions.MigrationsDirectory, MigrationFileRegexPattern)

	migrations, err := r.Build(files)
	if err != nil {
		return nil
	}
	applied := r.AppliedMigrations(ctx)

	var pending = ds.Queue[migrator.Migration]{}

	item := migrations.Dequeue()
	for item != nil {
		if v := applied.Find(func(m migrator.Migration) bool {
			return m.Version == item.Version
		}); v != nil {
			pending.Enqueue(*v)
		}

		item = migrations.Dequeue()
	}

	return &pending
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
	rows, err := r.engine.QueryContext(
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

	row := r.engine.QueryRowContext(
		ctx,
		fmt.Sprintf(
			"SELECT id, name, version, created_at FROM %v ORDER BY id DESC LIMIT 1",
			r.trackerOptions.Table,
		),
	)
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
	CREATE TABLE %v (
		id SERIAL,
		version varchar UNIQUE NOT NULL,
		name varchar UNIQUE NOT NULL,
		created_at timestamp NOT NULL DEFAULT now(),
		PRIMARY KEY(id)
	)`, r.trackerOptions.Table)

	_, err := r.engine.ExecContext(ctx, query)

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

	query := fmt.Sprintf(`DROP TABLE %v`, r.trackerOptions.Table)

	rows, err := r.engine.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	if n, err := rows.RowsAffected(); n <= 0 {
		return err
	}

	return nil
}

// MARK: - Implement Migration Builder
func (r *Runner) Build(files []fs.FileInfo) (*ds.Queue[migrator.Migration], error) {
	migrations, err := r.LoadMigrations(files, r.trackerOptions.MigrationsDirectory, r.trackerOptions.FileRegexPattern)
	if err != nil {
		return nil, err
	}

	q := ds.NewFromSlice(migrations)

	return q, nil
}

func (r *Runner) LoadMigrations(files []fs.FileInfo, dir string, pattern *regexp.Regexp) ([]migrator.Migration, error) {
	var migrations []migrator.Migration

	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		path, _ = filepath.Abs(path)
		contents, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var m migrator.Migration
		err = yaml.Unmarshal(contents, &m)
		if err != nil {
			return nil, err
		}

		match := pattern.FindStringSubmatch(file.Name())

		m.FileName = file.Name()
		m.Version = match[pattern.SubexpIndex("Version")]

		migrations = append(migrations, m)
	}

	return migrations, nil
}
