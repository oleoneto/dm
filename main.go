package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/oleoneto/dm/pkg/fsystem"
	"github.com/oleoneto/dm/pkg/logger"
	"github.com/oleoneto/dm/pkg/runner"
	// _ "github.com/oleoneto/dm/cli/cmd"
)

func main() {
	logger.NewLogger()

	dsn := os.Getenv("DATABASE_URL")
	dbEngine, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	defer dbEngine.Close()

	loader := fsystem.FileLoader{}

	r := runner.NewRunner(
		dbEngine,
		loader,
		runner.TrackerOptions{
			Schema: "public",
			Table:  "_migrations",
			// MigrationsDirectory: "examples",
			MigrationsDirectory: "examples/aluna",
		},
		runner.LoadMigrations,
		runner.Validate,
		runner.ApplyMigrations,
		runner.RollbackMigrations,
	)

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()

	// fmt.Println("Tracking?", r.IsTracked(ctx))
	// fmt.Println("Is Empty?", r.IsEmpty(ctx))
	// fmt.Println("StartTracking", r.StartTracking(ctx))
	// fmt.Println("StopTracking", r.StopTracking(ctx))
	fmt.Println("Tracking?", r.IsTracked(ctx))

	applied, aerr := r.AppliedMigrations(ctx)
	if aerr != nil {
		panic(aerr)
	}
	fmt.Println("Applied:", applied)

	pending, perr := r.PendingMigrations(ctx)
	if perr != nil {
		panic(perr)
	}
	fmt.Println("Pending:", pending)

	fmt.Println(r.Apply(ctx))
	// fmt.Println(r.Revert(ctx))
	fmt.Println("Done...!")
}
