package main

import (
	"github.com/oleoneto/dm/pkg/runner"
)

// _ "github.com/oleoneto/dm/cli/cmd"

func main() {
	r := runner.Runner{}

	r.Migrator.IsEmpty()
}
