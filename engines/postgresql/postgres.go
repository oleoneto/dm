package postgresql

import (
	"context"
	"log"
	"regexp"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Name        string
	Table       string
	Database    string
	Directory   string
	FilePattern *regexp.Regexp
}

const engineName string = "PostgreSQL"

var pgInstance *pgxpool.Pool

func (engine Postgres) acquireDatabaseConnection() {
	conn, err := pgxpool.Connect(context.Background(), engine.Database)

	pgInstance = conn

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
}

func Pg() *pgxpool.Pool {
	return pgInstance
}
