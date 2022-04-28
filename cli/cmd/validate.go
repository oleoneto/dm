package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration of migration files",
		Run: func(cmd *cobra.Command, args []string) {
			changes := migrator.Build(directory)

			if migrator.Validate(changes) {
				fmt.Println("Migrations are valid.")
			}
		},
	}
)
