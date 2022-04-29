package cmd

import (
	"fmt"
	"os"

	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var (
	migrateTo string

	// TODO: Add support for running migrations down to a given version
	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Run migration(s)",
		Run: func(cmd *cobra.Command, args []string) {

			if migrateTo != "" {
				version, err := parsedVersionFlag(migrateTo)

				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

				m := migrator.Build(directory)
				sequence, found := m.Find(strcase.ToCamel(version.Value))

				if found {
					// DEBUG: sequence.Display()
					migrator.Run(sequence, migrations.MigrateUp)
				}

				return
			}

			m := migrator.PendingMigrations(directory)

			if m.Size() > 0 {
				migrator.Run(m, migrations.MigrateUp)
			}
		},
	}
)

func init() {
	migrateCmd.PersistentFlags().StringVar(&migrateTo, "version", "", "run migrations up do this version")
}
