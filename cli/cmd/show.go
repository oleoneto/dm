package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	showCmd = &cobra.Command{
		Use:   "show",
		Short: "Shows the state of applied and pending migrations",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	allCmd = &cobra.Command{
		Use:   "all",
		Short: "List all migrations for a given application",
		Run: func(cmd *cobra.Command, args []string) {
			files := Engine.LoadFiles(directory, &FilePattern)

			for _, file := range files {
				fmt.Println(file.Name())
			}
		},
	}

	appliedCmd = &cobra.Command{
		Use:   "applied",
		Short: "List only applied migrations",
		Run: func(cmd *cobra.Command, args []string) {
			applied := Engine.AppliedMigrations()

			migration := applied.GetHead()

			for migration != nil {
				fmt.Println(migration.FileName)
				migration = migration.Next()
			}
		},
	}

	pendingCmd = &cobra.Command{
		Use:   "pending",
		Short: "List only pending migrations",
		Run: func(cmd *cobra.Command, args []string) {
			pending := Engine.PendingMigrations()

			migration := pending.GetHead()

			for migration != nil {
				fmt.Println(migration.FileName)
				migration = migration.Next()
			}
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Shows the most recently applied migration",
		Run: func(cmd *cobra.Command, args []string) {
			version, _ := Engine.Version()
			fmt.Printf("Current version: %v\n", version)
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	showCmd.AddCommand(allCmd)
	showCmd.AddCommand(appliedCmd)
	showCmd.AddCommand(pendingCmd)
	showCmd.AddCommand(versionCmd)
}
