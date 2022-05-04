package stores

import (
	"context"
	"log"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	instance *pgxpool.Pool
	URL      string
}

func (store *Postgres) Connect() error {
	conn, err := pgxpool.Connect(context.Background(), store.URL)

	store.instance = conn

	if err != nil || conn == nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		return err
	}

	return nil
}

func (store *Postgres) Disconnect() error {
	store.instance.Close()
	return nil
}

func (store Postgres) Name() string {
	return "PostgreSQL"
}

func (store Postgres) DatabaseURL() string {
	return store.URL
}

func (store Postgres) Create(query string, options ...interface{}) error {
	store.Connect()

	_, err := store.instance.Exec(context.Background(), query, options...)

	store.Disconnect()
	return err
}

func (store Postgres) Read(query string, model interface{}, options ...interface{}) error {
	store.Connect()

	rows, err := store.instance.Query(context.Background(), query, options...)

	if err != nil {
		return err
	}

	return pgxscan.ScanAll(model, rows)
}

func (store Postgres) Delete(query string, options ...interface{}) error {
	store.Connect()

	_, err := store.instance.Exec(context.Background(), query, options...)
	return err
}
