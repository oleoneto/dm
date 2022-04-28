package cmd

import (
	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all migrations for a given application",
		Run: func(cmd *cobra.Command, args []string) {
			var migrator migrations.Migrator

			migrator.ListFiles(directory)
		},
	}
)
