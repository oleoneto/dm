package cmd

import (
	"fmt"
	"os"

	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var (
	rollbackTo string

	// TODO: Add support for rolling back migrations up to a given version
	rollbackCmd = &cobra.Command{
		Use:   "rollback",
		Short: "Rollback migration(s)",
		Run: func(cmd *cobra.Command, args []string) {

			if rollbackTo != "" {
				version, err := parsedVersionFlag(rollbackTo)

				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

				m := migrator.Build(directory)
				m.Reverse()

				sequence, found := m.Find(strcase.ToCamel(version.Value))

				if found {
					// DEBUG: sequence.Display()
					migrator.Run(sequence, migrations.MigrateDown)
				}

				return
			}

			m := migrator.Build(directory)
			m.Reverse()
			migrator.Run(m, migrations.MigrateDown)
		},
	}
)

func init() {
	rollbackCmd.PersistentFlags().StringVar(&rollbackTo, "version", "", "rollback to this version")
}
