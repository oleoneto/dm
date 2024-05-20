package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/oleoneto/dm/pkg/fsystem"
	"github.com/oleoneto/dm/pkg/migrator"
	"github.com/oleoneto/dm/pkg/runner"
	// _ "github.com/oleoneto/dm/cli/cmd"
)

func main() {
	dbEngine, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	loader := fsystem.FileLoader{}

	migrationController := migrator.NewMigrationsController(dbEngine)
	r := runner.NewRunner(
		dbEngine,
		loader,
		migrationController,
		runner.TrackerOptions{Schema: "public", Table: "_migrations"},
	)

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	fmt.Println("Tracking?", r.IsTracked(ctx))
	fmt.Println("Is Empty?", r.IsEmpty(ctx))
	fmt.Println("StartTracking", r.StartTracking(ctx))
	// fmt.Println("StopTracking", r.StopTracking(ctx))
	fmt.Println("Tracking?", r.IsTracked(ctx))
	fmt.Println("Applied migrations:", r.AppliedMigrations(ctx))
	fmt.Println("Pending migrations", r.PendingMigrations(ctx))
}
