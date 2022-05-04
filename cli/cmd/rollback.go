package cmd

import (
	"fmt"
	"os"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var (
	rollbackTo string

	rollbackCmd = &cobra.Command{
		Use:     "rollback",
		Short:   "Rollback migration(s)",
		Aliases: []string{"r"},
		PreRun: func(cmd *cobra.Command, args []string) {
			validateDatabaseConfig()
		},
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var version VersionFlag

			if rollbackTo != "" {
				version, err = parsedVersionFlag(rollbackTo)

				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(INVALID_INPUT_ERROR)
				}
			}

			list := runner.AppliedMigrations(directory, &FilePattern)
			list.Reverse()

			if version.Value != "" {
				sequence, found := list.Find(strcase.ToCamel(version.Value))

				if !found {
					fmt.Fprintln(os.Stderr, "Error: Migration not found.")
					os.Exit(INVALID_INPUT_ERROR)
				}

				list = sequence
			}

			runner.Down(list)
		},
	}
)

func init() {
	rollbackCmd.PersistentFlags().StringVarP(&rollbackTo, "version", "v", "", "rollback this version (and anything applied after it)")
	rollbackCmd.PersistentFlags().StringVarP(&databaseUrl, "database-url", "u", databaseUrl, "database url")
	rollbackCmd.MarkFlagRequired("database-url")
	rollbackCmd.MarkFlagRequired("adapter")
	rollbackCmd.MarkFlagRequired("table")
}
