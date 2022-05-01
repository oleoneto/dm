package cmd

import (
	"fmt"
	"os"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var (
	migrateTo string

	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Run migration(s)",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var version VersionFlag

			if migrateTo != "" {
				version, err = parsedVersionFlag(migrateTo)

				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}

			list := Engine.PendingMigrations()

			if version.Value != "" {
				sequence, found := list.Find(strcase.ToCamel(version.Value))

				if !found {
					return
				}

				list = sequence
			}

			Engine.Up(list)
		},
	}
)

func init() {
	migrateCmd.PersistentFlags().StringVar(&migrateTo, "version", "", "run migrations up do this version")
}
