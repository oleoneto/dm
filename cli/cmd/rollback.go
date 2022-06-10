package cmd

import (
	"os"

	"github.com/iancoleman/strcase"
	"github.com/oleoneto/dm/logger"
	"github.com/spf13/cobra"
)

var (
	rollbackCmd = &cobra.Command{
		Use:     "rollback NAME|VERSION",
		Short:   "Rollback migration(s)",
		Aliases: []string{"r"},
		Args:    cobra.MaximumNArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
				message := logger.ApplicationMessage{Message: "No applied migrations to rollback."}
				logger.Custom(format, template).WithFormattedOutput(&message, os.Stdout)
				return
			}

			list.Reverse()

			if version.Value != "" {
				sequence, found := list.Find(strcase.ToCamel(version.Value))

				if !found {
					message := logger.ApplicationMessage{Message: "Nothing to do."}
					logger.Custom(format, template).WithFormattedOutput(&message, os.Stdout)
					return
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
