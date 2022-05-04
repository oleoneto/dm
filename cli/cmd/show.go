package cmd

import (
	"fmt"
	"strings"

	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/iancoleman/strcase"
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

			for _, file := range files {
				version, name, _ := strings.Cut(file.Name(), "_")
				name = strings.Split(name, ".")[0]
				name = strcase.ToCamel(name)

				fmt.Printf("Version: %v (%v)\n", version, name)
			}
		},
	}

	appliedCmd = &cobra.Command{
		Use:   "applied",
		Short: "List only applied migrations",
		Run: func(cmd *cobra.Command, args []string) {
			applied := runner.AppliedMigrations(directory, &FilePattern)

			migration := applied.GetHead()

			for migration != nil {
				fmt.Println(migration.Description())
				migration = migration.Next()
			}
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
			pending := runner.PendingMigrations(directory, &FilePattern)

			migration := pending.GetHead()

			for migration != nil {
				fmt.Println(migration.Description())
				migration = migration.Next()
			}
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Shows the most recently applied migration",
		Run: func(cmd *cobra.Command, args []string) {
			version, _ := runner.Version()
			fmt.Printf("Current version: %v\n", version)
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
