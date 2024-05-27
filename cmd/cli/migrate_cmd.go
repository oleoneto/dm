package cli

import (
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:     "migrate NAME|VERSION",
	Short:   "Run migration(s)",
	Aliases: []string{"m"},
	Args:    cobra.MaximumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// validateDatabaseConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	// migrateCmd.PersistentFlags().VarP(state.Flags.DatabaseURL, "database-url", "u", "database url")
	// migrateCmd.MarkFlagRequired("database-url")
	migrateCmd.MarkFlagRequired("adapter")
	migrateCmd.MarkFlagRequired("table")
}
