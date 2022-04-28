package main

import (
	"fmt"

	cli "github.com/cleopatrio/db-migrator-lib/cli/cmd"
)

func init() {
	fmt.Println("Database Migrator v0.1.0-alpha")
}

func main() {
	cli.Execute()
}
