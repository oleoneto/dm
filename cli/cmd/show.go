package cmd

import (
	"os"

	"github.com/cleopatrio/db-migrator-lib/logger"
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

			logger.Custom(format, template).WithFormattedOutput(&m, os.Stdout)
		},
	}

	appliedCmd = &cobra.Command{
		Use:   "applied",
		Short: "List only applied migrations",
		Run: func(cmd *cobra.Command, args []string) {
			loadFromDir := false
			list := runner.AppliedMigrations(directory, &FilePattern, loadFromDir)
			m := list.ToSlice()

			logger.Custom(format, template).WithFormattedOutput(&m, os.Stdout)
		},
	}

	pendingCmd = &cobra.Command{
		Use:     "pending",
		Short:   "List only pending migrations",
		Aliases: []string{"p"},
		Run: func(cmd *cobra.Command, args []string) {
			list := runner.PendingMigrations(directory, &FilePattern)
			m := list.ToSlice()

			logger.Custom(format, template).WithFormattedOutput(&m, os.Stdout)
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Shows the most recently applied migration",
		Run: func(cmd *cobra.Command, args []string) {
			version, _ := runner.Version()

			logger.Custom(format, template).WithFormattedOutput(&version, os.Stdout)
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
