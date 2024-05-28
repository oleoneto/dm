package cli

import (
	"github.com/oleoneto/dm/cmd/cli/core"
	"github.com/spf13/cobra"
)

func Execute() error {
	setupGlobalFlags()

	return RootCmd.Execute()
}

var state = core.NewCommandState()

var RootCmd = &cobra.Command{
	Use:               "dm",
	Short:             "DM, short for Database Migrator is a migration management tool.",
	PersistentPreRun:  state.BeforeHook,
	PersistentPostRun: state.AfterHook,
	Run:               func(cmd *cobra.Command, args []string) { cmd.Help() },
}

func init() {
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(ShowCmd)
	RootCmd.AddCommand(MigrateCmd)
	RootCmd.AddCommand(RollbackCmd)
}

func setupGlobalFlags() {
	RootCmd.PersistentFlags().VarP(state.Flags.OutputFormat, "output", "o", "output format")
	RootCmd.PersistentFlags().StringVarP(&state.Flags.OutputTemplate, "output-template", "y", state.Flags.OutputTemplate, "template (used when output format is 'gotemplate')")

	// Migrator configuration
	RootCmd.PersistentFlags().VarP(state.Flags.Engine, "adapter", "a", "database adapter")
	RootCmd.PersistentFlags().StringVarP(&state.Flags.Directory, "directory", "d", state.Flags.Directory, "migrations directory")
	RootCmd.PersistentFlags().StringVarP(&state.Flags.Table, "table", "t", state.Flags.Table, "table wherein migrations are tracked")
}
