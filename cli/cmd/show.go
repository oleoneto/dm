package cmd

import (
	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/spf13/cobra"
)

var (
	showCmd = &cobra.Command{
		Use:   "show",
		Short: "Shows the state of applied and pending migrations",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			validateDatabaseConfig()
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}

	allCmd = &cobra.Command{
		Use:   "all",
		Short: "List all migrations for a given application",
		Run: func(cmd *cobra.Command, args []string) {
			files := migrations.LoadFiles(directory, &FilePattern)
			list := migrations.BuildMigrations(files, directory, &FilePattern)
			m := list.ToSlice()

			WithFormattedOutput(&m)
		},
	}

	appliedCmd = &cobra.Command{
		Use:   "applied",
		Short: "List only applied migrations",
		Run: func(cmd *cobra.Command, args []string) {
			loadFromDir := false
			list := runner.AppliedMigrations(directory, &FilePattern, loadFromDir)
			m := list.ToSlice()

			WithFormattedOutput(&m)
		},
	}

	pendingCmd = &cobra.Command{
		Use:     "pending",
		Short:   "List only pending migrations",
		Aliases: []string{"p"},
		PreRun: func(cmd *cobra.Command, args []string) {
			validateDatabaseConfig()
		},
		Run: func(cmd *cobra.Command, args []string) {
			list := runner.PendingMigrations(directory, &FilePattern)
			m := list.ToSlice()

			WithFormattedOutput(&m)
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Shows the most recently applied migration",
		Run: func(cmd *cobra.Command, args []string) {
			version, _ := runner.Version()

			WithFormattedOutput(&version)
		},
	}
)

func init() {
	showCmd.AddCommand(allCmd)
	showCmd.AddCommand(appliedCmd)
	showCmd.AddCommand(pendingCmd)
	showCmd.AddCommand(versionCmd)

	showCmd.PersistentFlags().StringVarP(&databaseUrl, "database-url", "u", databaseUrl, "database url")
	showCmd.MarkFlagRequired("database-url")
	showCmd.MarkFlagRequired("adapter")
	showCmd.MarkFlagRequired("table")
}
