package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate the configuration of migration files",
		Run: func(cmd *cobra.Command, args []string) {
			files := Engine.LoadFiles(directory, &FilePattern)

			if len(files) == 0 {
				fmt.Println("No migrations found.")
				return
			}

			migrations := Engine.BuildMigrations(files)
			valid, reason := Engine.Validate(migrations)

			if valid && reason == "" {
				fmt.Println("Migrations are valid.")
				return
			}

			fmt.Println(reason)
		},
	}
)
