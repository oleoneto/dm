package cli

import (
	"github.com/oleoneto/dm/cmd/cli/core"
	"github.com/spf13/cobra"
)

var state = core.NewCommandState()

var rootCmd = &cobra.Command{
	Use:   "dm",
	Short: "DM, short for Database Migrator is a migration management tool.",
	Run:   func(cmd *cobra.Command, args []string) { cmd.Help() },
}

func Execute() error { return rootCmd.Execute() }
