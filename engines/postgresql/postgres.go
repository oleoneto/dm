package postgresql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Name     string `default:"PostgreSQL"`
	Table    string
	Database string
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
