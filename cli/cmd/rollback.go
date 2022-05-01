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
		Use:   "rollback",
		Short: "Rollback migration(s)",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var version VersionFlag

			if rollbackTo != "" {
				version, err = parsedVersionFlag(rollbackTo)

				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}

			list := Engine.AppliedMigrations()
			list.Reverse()

			if version.Value != "" {
				sequence, found := list.Find(strcase.ToCamel(version.Value))

				if !found {
					return
				}

				list = sequence
			}

			Engine.Down(list)
		},
	}
)

func init() {
	rollbackCmd.PersistentFlags().StringVar(&rollbackTo, "version", "", "rollback this version (and anything applied after it)")
}
