package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/oleoneto/dm/pkg/fsystem"
	"github.com/oleoneto/dm/pkg/runner"
	// _ "github.com/oleoneto/dm/cli/cmd"
)

func main() {
	dbEngine, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	loader := fsystem.FileLoader{}

	r := runner.NewRunner(
		dbEngine,
		loader,
		runner.TrackerOptions{
			Schema: "public",
			Table:  "_migrations",
			// MigrationsDirectory: "examples",
			MigrationsDirectory: "/Users/cleopatrio.neto/Developer/Explorations/aluna/migrations",
		},
	)

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	fmt.Println("Tracking?", r.IsTracked(ctx))
	fmt.Println("Is Empty?", r.IsEmpty(ctx))
	fmt.Println("StartTracking", r.StartTracking(ctx))
	// fmt.Println("StopTracking", r.StopTracking(ctx))
	fmt.Println("Tracking?", r.IsTracked(ctx))
	fmt.Println("Applied:", r.AppliedMigrations(ctx))
	fmt.Println("Pending:", r.PendingMigrations(ctx))

	fmt.Println(r.Apply(ctx))
}
