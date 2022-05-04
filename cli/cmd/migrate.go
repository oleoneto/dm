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
		Use:     "migrate",
		Short:   "Run migration(s)",
		Aliases: []string{"m"},
		PreRun: func(cmd *cobra.Command, args []string) {
			validateDatabaseConfig()
		},
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var version VersionFlag

			if migrateTo != "" {
				version, err = parsedVersionFlag(migrateTo)

				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(INVALID_INPUT_ERROR)
				}
			}

			list := runner.PendingMigrations(directory, &FilePattern)

			if version.Value != "" {
				sequence, found := list.Find(strcase.ToCamel(version.Value))

				if !found {
					fmt.Fprintln(os.Stderr, "Error: Migration not found.")
					os.Exit(INVALID_INPUT_ERROR)
				}

				list = sequence
			}

			runner.Up(list)
		},
	}
)

func init() {
	migrateCmd.PersistentFlags().StringVarP(&migrateTo, "version", "v", "", "run migrations up do this version")
	migrateCmd.PersistentFlags().StringVarP(&databaseUrl, "database-url", "u", databaseUrl, "database url")
	migrateCmd.MarkFlagRequired("database-url")
	migrateCmd.MarkFlagRequired("adapter")
	migrateCmd.MarkFlagRequired("table")
}
