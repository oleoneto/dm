package cmd

import (
	"fmt"
	"os"

	"github.com/cleopatrio/db-migrator-lib/logger"
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
	format       = "plain"
	template     = ""

	SUPPORTED_ADAPTERS = map[string]migrations.Store{
		"postgresql": stores.Postgres{URL: databaseUrl},
		// "sqlite3":    stores.SQLite3{URL: databaseUrl},
	}

	rootCmd = &cobra.Command{
		Use:   "dm",
		Short: "DM, short for Database Migrator is a migration management tool.",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	cliVersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Shows the version of the CLI",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Custom(format, template).WithFormattedOutput(&version, os.Stdout)
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func validateDatabaseConfig() {
	if databaseUrl == "" {
		message := logger.ApplicationError{
			Error: "No database specified.\nProvide a value for the flag or set DATABASE_URL in your environment.\n",
		}

		logger.Custom(format, template).WithFormattedOutput(&message, os.Stderr)
		os.Exit(101)
	}

	selectedAdapter, ok := SUPPORTED_ADAPTERS[adapter]

	if !ok {
		message := logger.ApplicationError{Error: fmt.Sprintf("Unsupported adapter '%v'", adapter)}
		logger.Custom(format, template).WithFormattedOutput(&message, os.Stderr)
		os.Exit(102)
	}

	storeAdapter = selectedAdapter
	runner.SetStore(storeAdapter)
	runner.SetSchemaTable(table)
	runner.SetLogger(format, template)
}

func init() {
	// CLI configuration
	// cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&config, "config", config, "config file")

	// Migrator configuration
	rootCmd.PersistentFlags().StringVarP(&adapter, "adapter", "a", adapter, "database adapter")
	rootCmd.PersistentFlags().StringVarP(&directory, "directory", "d", directory, "migrations directory")
	rootCmd.PersistentFlags().StringVarP(&table, "table", "t", table, "table wherein migrations are tracked")
	rootCmd.PersistentFlags().StringVarP(&format, "output-format", "o", format, "output format")
	rootCmd.PersistentFlags().StringVarP(&template, "output-template", "y", template, "template (used when output format is 'gotemplate')")

	// Sub-commands
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rollbackCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(cliVersionCmd)

	// Runner configuration
	// These changes can be overridden by validateDatabaseConfig()
	runner.SetStore(SUPPORTED_ADAPTERS[adapter])
	runner.SetSchemaTable(table)
}
