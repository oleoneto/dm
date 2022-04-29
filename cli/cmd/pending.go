package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	pendingCmd = &cobra.Command{
		Use:   "pending",
		Short: "List all pending migrations",
		Run: func(cmd *cobra.Command, args []string) {
			pending := migrator.PendingMigrations(directory)

			migration := pending.GetHead()

			for migration != nil {
				fmt.Printf("Name: %v, Version: %v\n", migration.Name, migration.Version)
				migration = migration.Next()
			}
		},
	}
)
