package cmd

import (
	"fmt"
	"os"

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

				files := Engine.LoadFiles(directory, &FilePattern)
				list := Engine.BuildMigrations(files)
				list.Reverse()

				sequence, found := list.Find(strcase.ToCamel(version.Value))

				if found {
					// DEBUG: sequence.Display()
					Engine.Down(sequence)
				}

				return
			}

			files := Engine.LoadFiles(directory, &FilePattern)
			list := Engine.BuildMigrations(files)
			list.Reverse()
			Engine.Down(list)
		},
	}
)

func init() {
	rollbackCmd.PersistentFlags().StringVar(&rollbackTo, "version", "", "rollback to this version")
}
