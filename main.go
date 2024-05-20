package main

import (
	"context"
	"time"

	"github.com/oleoneto/dm/pkg/engines"
	"github.com/oleoneto/dm/pkg/migrator"
	"github.com/oleoneto/dm/pkg/runner"
)

// _ "github.com/oleoneto/dm/cli/cmd"

func main() {
	dbEngine := &engines.PostgreSQL{}

	m := migrator.NewMigrationsController(dbEngine)
	r := runner.NewRunner(m, dbEngine, runner.TrackerOptions{})

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()

	r.IsTracked(ctx)
}
