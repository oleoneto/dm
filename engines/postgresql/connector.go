package postgresql

import (
	"context"
	"log"
	"regexp"

	"github.com/jackc/pgx/v4/pgxpool"
)

func (engine Postgres) ConnectionPattern() *regexp.Regexp {
	// TODO: Test and validate behavior
	return regexp.MustCompile(`postgres:\/\/(?P<username>\w+)\:(?P<password>\w+)\@(?P<host>[\w+\.?]+)\:(?P<port>\d+)\/(?P<database>\w+)(\?)?\w+$`)
}

func (engine Postgres) Connect() error {
	conn, err := pgxpool.Connect(context.Background(), engine.Database)

	pgInstance = conn

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		return err
	}

	return nil
}

func (engine Postgres) Disconnect() error {
	pgInstance.Close()
	return nil
}
