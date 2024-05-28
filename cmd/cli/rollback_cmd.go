package cli

import (
	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:     "rollback NAME|VERSION",
	Short:   "Rollback migration(s)",
	Aliases: []string{"r"},
	Args:    cobra.MaximumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		state.BeforeHook(cmd, args)
		state.ConnectDatabase(cmd, args)
	},
	PersistentPostRun: state.AfterHook,
	Run:               func(cmd *cobra.Command, args []string) {},
}

func init() {
	rollbackCmd.PersistentFlags().StringVarP(&state.Flags.DatabaseURL, "database-url", "u", "", "database url")
	rollbackCmd.MarkFlagRequired("database-url")
	rollbackCmd.MarkFlagRequired("adapter")
	rollbackCmd.MarkFlagRequired("table")
}
