package cmd

import (
	"os"

	"github.com/cleopatrio/db-migrator-lib/engines/postgresql"
	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/spf13/cobra"
)

var (
	config      string
	directory   string
	engine      string // TODO: Limit choices
	migrator    migrations.Migrator
	databaseUrl string
	table       string

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
	// TODO: Receive from Migration.`Adapter`
	migrator.Engine = postgresql.Postgres{Table: table, Database: databaseUrl}
	migrator.Directory = directory
	migrator.DatabaseUrl = databaseUrl
	migrator.Table = table
}

func init() {
	cobra.OnInitialize(initConfig)

	// CLI configuration
	rootCmd.PersistentFlags().StringVar(&config, "config", "", "config file")

	// Migrator configuration
	rootCmd.PersistentFlags().StringVar(&engine, "engine", "postgres", "database engine")
	rootCmd.PersistentFlags().StringVar(&databaseUrl, "database-url", os.Getenv("DATABASE_URL"), "database url")
	rootCmd.PersistentFlags().StringVar(&directory, "directory", "./migrations", "migrations directory")
	rollbackCmd.PersistentFlags().StringVar(&table, "table", "_migrations", "table wherein migrations are tracked")

	// Sub-commands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rollbackCmd)
	rootCmd.AddCommand(statusCmd)
}
