package cmd

import (
	"fmt"
	"os"

	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/cleopatrio/db-migrator-lib/stores"
	"github.com/spf13/cobra"
)

var (
	config       string
	runner       = migrations.Runner{}
	storeAdapter migrations.Store
	directory    = "./migrations"
	adapter      = "postgresql"
	databaseUrl  = os.Getenv("DATABASE_URL")
	table        = "_migrations"
	FilePattern  = migrations.FilePattern

	SUPPORTED_ADAPTERS = map[string]migrations.Store{
		"postgresql": stores.Postgres{URL: databaseUrl},
		// "sqlite3":    stores.SQLite3{URL: databaseUrl},
	}

	rootCmd = &cobra.Command{
		Use:   "dm",
		Short: "DM, short for Database Migrator is a migration management tool.",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func validateDatabaseConfig() {
	if databaseUrl == "" {
		fmt.Fprintf(os.Stderr, "No database specified.\nProvide a value for the flag or set DATABASE_URL in your environment.\n")
		os.Exit(1)
	}

	selectedAdapter, ok := SUPPORTED_ADAPTERS[adapter]

	if !ok {
		fmt.Fprintf(os.Stderr, "Unsupported adapter '%v'.\n", adapter)
		os.Exit(1)
	}

	storeAdapter = selectedAdapter
	runner.SetStore(storeAdapter)
	runner.SetSchemaTable(table)
}

func init() {
	// CLI configuration
	// cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&config, "config", config, "config file")

	// Migrator configuration
	rootCmd.PersistentFlags().StringVarP(&adapter, "adapter", "a", adapter, "database adapter")
	rootCmd.PersistentFlags().StringVarP(&directory, "directory", "d", directory, "migrations directory")
	rootCmd.PersistentFlags().StringVarP(&table, "table", "t", table, "table wherein migrations are tracked")

	// Sub-commands
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rollbackCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(validateCmd)

	// Runner configuration
	// These changes can be overridden by validateDatabaseConfig()
	runner.SetStore(SUPPORTED_ADAPTERS[adapter])
	runner.SetSchemaTable(table)
}
