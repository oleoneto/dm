package main

import (
	"fmt"
	"os"

	"github.com/cleopatrio/db-migrator-lib/api/server"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	databaseURL         = kingpin.Flag("database-url", "database url").OverrideDefaultFromEnvar("DATABASE_URL").Required().String()
	migrationsDirectory = kingpin.Flag("directory", "migrations directory").OverrideDefaultFromEnvar("MIGRATIONS_DIRECTORY").Required().String()
	migrationsTableName = kingpin.Flag("table", "migrations table").OverrideDefaultFromEnvar("MIGRATIONS_TABLE").Required().String()
	serverPort          = kingpin.Flag("port", "Server port").OverrideDefaultFromEnvar("SERVER_PORT").Int()
)

func overrideEnvironmentVars() {
	os.Setenv("DATABASE_URL", *databaseURL)
	os.Setenv("MIGRATIONS_DIRECTORY", *migrationsDirectory)
	os.Setenv("MIGRATIONS_TABLE", *migrationsTableName)
}

func main() {
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	overrideEnvironmentVars()

	if *serverPort == 0 {
		*serverPort = 3809
	}

	app := server.API()

	app.Run(fmt.Sprintf(":%v", *serverPort))
}
