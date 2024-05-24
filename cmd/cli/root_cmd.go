package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "dm",
	Short: "DM, short for Database Migrator is a migration management tool.",
	Run:   func(cmd *cobra.Command, args []string) { cmd.Help() },
}

func Execute() error { return rootCmd.Execute() }
