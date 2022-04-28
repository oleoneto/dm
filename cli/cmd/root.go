package cmd

import (
	"github.com/cleopatrio/db-migrator-lib/engines/postgresql"
	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/spf13/cobra"
)

var (
	config    string
	directory string
	engine    string // TODO: Limit choices
	migrator  migrations.Migrator

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
	migrator.Engine = postgresql.Postgres{} // TODO: Receive from Migration.`Adapter`
	migrator.Directory = directory
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&config, "config", "", "config file (default is $HOME/.dm.yaml)")
	rootCmd.PersistentFlags().StringVar(&engine, "engine", "postgres", "database engine (default is postgres)")
	rootCmd.PersistentFlags().StringVar(&directory, "directory", "migrations", "migrations directory (default is $PWD/migrations)")

	// Sub-commands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rollbackCmd)
	rootCmd.AddCommand(statusCmd)
}
