package runner

import (
	"regexp"

	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/engines"
	"github.com/oleoneto/dm/pkg/fsystem"
	"github.com/oleoneto/dm/pkg/migrator"
)

var MigrationFileRegexPattern = func() *regexp.Regexp {
	return regexp.MustCompile(`(?P<Version>^\d{20})_(?P<Name>[aA-zZ]+).yaml|sql$`)
}()

type Runner struct {
	engine         engines.SqlEngineProtocol
	loader         fsystem.FileLoaderProtocol
	trackerOptions TrackerOptions

	// Data
	migrations ds.Queue[migrator.Migration]

	migrationLoaderFunc LoadMigrationsFunc
	validatorFunc       ValidateMigrationsFunc
	migrateUpFunc       MigrateUpFunc
	migrateDownFunc     MigrateDownFunc
}

type TrackerOptions struct {
	Schema              string
	Table               string
	MigrationsDirectory string
	FileRegexPattern    *regexp.Regexp
}

func NewRunner(
	e engines.SqlEngineProtocol,
	l fsystem.FileLoaderProtocol,
	trackerOptions TrackerOptions,
	lf LoadMigrationsFunc,
	vf ValidateMigrationsFunc,
	muf MigrateUpFunc,
	mdf MigrateDownFunc,
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
		engine:              e,
		loader:              l,
		migrationLoaderFunc: lf,
		trackerOptions:      trackerOptions,
		validatorFunc:       vf,
		migrateUpFunc:       muf,
		migrateDownFunc:     mdf,
	}
}
