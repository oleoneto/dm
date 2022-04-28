package cmd

import (
	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/spf13/cobra"
)

var (
	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Run migration(s)",
		Run: func(cmd *cobra.Command, args []string) {
			m := migrator.Build(directory)
			migrator.Run(m, migrations.MIGRATE_UP)
		},
	}
)
