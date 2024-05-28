package runner

import (
	"regexp"

	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/engines"
	"github.com/oleoneto/dm/pkg/migrator"
)

var MigrationFileRegexPattern = func() *regexp.Regexp {
	return regexp.MustCompile(`(?P<Version>^\d{20})_(?P<Name>[aA-zZ]+).yaml|sql$`)
}()

type Runner struct {
	engine         engines.SqlEngineProtocol
	trackerOptions TrackerOptions

	// Data
	migrations ds.Queue[migrator.Migration]

	fileLoaderFunc      FileLoaderFunc
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
	trackerOptions *TrackerOptions,
	e engines.SqlEngineProtocol,
	fl FileLoaderFunc,
	lf LoadMigrationsFunc,
	vf ValidateMigrationsFunc,
	muf MigrateUpFunc,
	mdf MigrateDownFunc,
) *Runner {
	var defaultTrackerOptions = TrackerOptions{
		MigrationsDirectory: "migrations",
		Table:               "_migrations",
		Schema:              "public",
		FileRegexPattern:    MigrationFileRegexPattern,
	}

	if trackerOptions == nil {
		trackerOptions = &defaultTrackerOptions
	} else {
		if trackerOptions.MigrationsDirectory == "" {
			trackerOptions.MigrationsDirectory = defaultTrackerOptions.MigrationsDirectory
		}

		if trackerOptions.Table == "" {
			trackerOptions.Table = defaultTrackerOptions.Table
		}

		if trackerOptions.Schema == "" {
			trackerOptions.Schema = defaultTrackerOptions.Schema
		}

		if trackerOptions.FileRegexPattern == nil {
			trackerOptions.FileRegexPattern = defaultTrackerOptions.FileRegexPattern
		}
	}

	if e == nil {
		panic("sql engine not set")
	}

	return &Runner{
		trackerOptions:      *trackerOptions,
		engine:              e,
		fileLoaderFunc:      fl,
		migrationLoaderFunc: lf,
		validatorFunc:       vf,
		migrateUpFunc:       muf,
		migrateDownFunc:     mdf,
	}
}
