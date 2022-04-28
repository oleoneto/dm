package cmd

import (
	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/spf13/cobra"
)

var (
	rollbackCmd = &cobra.Command{
		Use:   "rollback",
		Short: "Rollback migration(s)",
		Run: func(cmd *cobra.Command, args []string) {
			m := migrator.Build(directory)
			migrator.Run(m, migrations.MIGRATE_DOWN)
		},
	}
)
