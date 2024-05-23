package runner_test

import (
	"io/fs"
	"regexp"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/oleoneto/dm/pkg/fsystem"
)

var Postgres = embeddedpostgres.NewDatabase()

type Loader struct{ data []fs.DirEntry }

var _ fsystem.FileLoaderProtocol = (*Loader)(nil)

func (l Loader) LoadFiles(string, *regexp.Regexp) []fs.DirEntry { return l.data }
