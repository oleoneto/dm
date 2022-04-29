package cmd

import (
	"fmt"
	"os"

	"github.com/cleopatrio/db-migrator-lib/engines/postgresql"
	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/spf13/cobra"
)

var (
	config      string
	directory   = "./migrations"
	engine      = "postgresql"
	migrator    migrations.Migrator
	databaseUrl = os.Getenv("DATABASE_URL")
	table       = "_migrations"

	rootCmd = &cobra.Command{
		Use:   "dm",
		Short: "DM, short for Database Migrator is a migration management tool.",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
	if databaseUrl == "" {
		fmt.Fprintf(os.Stderr, "No database specified.\nProvide a value for the flag or set DATABASE_URL in your environment.\n")
		os.Exit(1)
	}

	SUPPORTED_ENGINES := map[string]migrations.Engine{
		"postgresql": postgresql.Postgres{
			Name:      "PostgreSQL",
			Table:     table,
			Database:  databaseUrl,
			Directory: directory,
		},
	}

	selectedEngine, ok := SUPPORTED_ENGINES[engine]

	if !ok {
		fmt.Fprintf(os.Stderr, "Unsupported engine '%v'.\n", engine)
		os.Exit(1)
	}

	// TODO: Receive from Migration.`Adapter`
	migrator.Engine = selectedEngine
	migrator.SupportedEngines = SUPPORTED_ENGINES
	migrator.Directory = directory
	migrator.DatabaseUrl = databaseUrl
	migrator.Table = table
}

func init() {
	cobra.OnInitialize(initConfig)

	// CLI configuration
	rootCmd.PersistentFlags().StringVar(&config, "config", config, "config file")

	// Migrator configuration
	rootCmd.PersistentFlags().StringVar(&engine, "engine", engine, "database engine")
	rootCmd.PersistentFlags().StringVar(&databaseUrl, "database-url", databaseUrl, "database url")
	rootCmd.PersistentFlags().StringVar(&directory, "directory", directory, "migrations directory")
	rootCmd.PersistentFlags().StringVar(&table, "table", table, "table wherein migrations are tracked")

	// Sub-commands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rollbackCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(pendingCmd)
	rootCmd.AddCommand(validateCmd)
}
