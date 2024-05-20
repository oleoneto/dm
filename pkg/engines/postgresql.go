package engines

import (
	"context"
	"database/sql"
)

type PostgreSQL struct{ db *sql.DB }

var _ SqlEngine = (*PostgreSQL)(nil)

func (e PostgreSQL) Name() string { return "postgresql" }

func (e *PostgreSQL) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	r, err := e.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (e *PostgreSQL) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	r, err := e.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (e *PostgreSQL) QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	r := e.db.QueryRowContext(ctx, query, args...)
	if r.Err() != nil {
		return nil, r.Err()
	}
	return r, nil
}
