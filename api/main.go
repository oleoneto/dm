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
	namespace           = kingpin.Flag("namespace", "namespace prefix").OverrideDefaultFromEnvar("API_NAMESPACE").String()
	serverPort          = kingpin.Flag("port", "server port").OverrideDefaultFromEnvar("SERVER_PORT").Int()
	versionPrefix       = kingpin.Flag("version-prefix", "API version prefix").OverrideDefaultFromEnvar("API_VERSION").String()
)

func overrideEnvironmentVars() {
	os.Setenv("API_NAMESPACE", *namespace)
	os.Setenv("API_VERSION", *versionPrefix)
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

	if *namespace == "" {
		*namespace = "migrations"
	}

	if *versionPrefix == "" {
		*versionPrefix = "v1"
	}

	app := server.API(*versionPrefix, *namespace)

	app.Run(fmt.Sprintf(":%v", *serverPort))
}
