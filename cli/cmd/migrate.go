package cmd

import (
	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/spf13/cobra"
)

var (
	migrateToVersion string

	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Run migration(s)",
		Run: func(cmd *cobra.Command, args []string) {
			m := migrator.Build(directory)
			migrator.Run(m, migrations.MIGRATE_UP)
		},
	}
)

func migrateConfig() {}

func init() {
	cobra.OnInitialize(migrateConfig)

	migrateCmd.PersistentFlags().StringVar(&migrateToVersion, "version", "", "run migrations up do this version")
}
