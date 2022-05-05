package cmd

import (
	"fmt"
	"os"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var (
	rollbackCmd = &cobra.Command{
		Use:     "rollback NAME|VERSION",
		Short:   "Rollback migration(s)",
		Aliases: []string{"r"},
		Args:    cobra.MaximumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			validateDatabaseConfig()
		},
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var version VersionFlag

			if len(args) > 0 && args[0] != "" {
				version, err = parsedVersionFlag(args[0])

				if err != nil {
					os.Exit(INVALID_INPUT_ERROR)
				}
			}

			loadFromDir := true
			list := runner.AppliedMigrations(directory, &FilePattern, loadFromDir)

			if list.Size() == 0 {
				fmt.Println("No applied migrations to rollback.")
				return
			}

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
	rollbackCmd.PersistentFlags().StringVarP(&databaseUrl, "database-url", "u", databaseUrl, "database url")
	rollbackCmd.MarkFlagRequired("database-url")
	rollbackCmd.MarkFlagRequired("adapter")
	rollbackCmd.MarkFlagRequired("table")
}
