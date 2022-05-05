package main

import (
	"fmt"

	cli "github.com/cleopatrio/db-migrator-lib/cli/cmd"
)

func init() {
	fmt.Printf("Database Migrator %v\n", version)
}

func main() {
	cli.Execute()
}
