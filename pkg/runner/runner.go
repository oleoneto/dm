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
	migrator       migrator.MigratorProtocol
	trackerOptions TrackerOptions
	migrations     ds.Queue[migrator.Migration]
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
		trackerOptions: trackerOptions,
	}
}
