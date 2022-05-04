package cmd

import (
	"fmt"

	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/spf13/cobra"
)

var (
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate the configuration of migration files",
		Run: func(cmd *cobra.Command, args []string) {
			files := migrations.LoadFiles(directory, &FilePattern)

			if len(files) == 0 {
				fmt.Println("No migrations found.")
				return
			}

			list := migrations.BuildMigrations(files, directory, &FilePattern)
			valid, reason := migrations.Validate(list)

			if valid && reason == "" {
				fmt.Println("Migrations are valid.")
				return
			}

			fmt.Println(reason)
		},
	}
)
