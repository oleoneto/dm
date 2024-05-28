package cli

import (
	"github.com/oleoneto/dm/cmd/cli/core"
	"github.com/spf13/cobra"
)

func Execute() error {
	setupGlobalFlags()

	return rootCmd.Execute()
}

var state = core.NewCommandState()

var rootCmd = &cobra.Command{
	Use:               "dm",
	Short:             "DM, short for Database Migrator is a migration management tool.",
	PersistentPreRun:  state.BeforeHook,
	PersistentPostRun: state.AfterHook,
	Run:               func(cmd *cobra.Command, args []string) { cmd.Help() },
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rollbackCmd)
}

func setupGlobalFlags() {
	rootCmd.PersistentFlags().VarP(state.Flags.OutputFormat, "output", "o", "output format")
	rootCmd.PersistentFlags().StringVarP(&state.Flags.OutputTemplate, "output-template", "y", state.Flags.OutputTemplate, "template (used when output format is 'gotemplate')")

	// Migrator configuration
	rootCmd.PersistentFlags().VarP(state.Flags.Engine, "adapter", "a", "database adapter")
	rootCmd.PersistentFlags().StringVarP(&state.Flags.Directory, "directory", "d", state.Flags.Directory, "migrations directory")
	rootCmd.PersistentFlags().StringVarP(&state.Flags.Table, "table", "t", state.Flags.Table, "table wherein migrations are tracked")
}
