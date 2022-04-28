package cmd

import (
	"github.com/spf13/cobra"
)

var (
	pendingCmd = &cobra.Command{
		Use:   "pending",
		Short: "List all pending migrations",
		Run: func(cmd *cobra.Command, args []string) {
			migrator.PendingMigrations(directory)
		},
	}
)
