package engines

import (
	"context"
	"database/sql"
)

type SQLite3 struct{ db *sql.DB }

var _ SqlEngine = (*SQLite3)(nil)

func (e SQLite3) Name() string { return "sqlite3" }

func (e *SQLite3) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	r, err := e.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (e *SQLite3) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	r, err := e.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (e *SQLite3) QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	r := e.db.QueryRowContext(ctx, query, args...)
	if r.Err() != nil {
		return nil, r.Err()
	}
	return r, nil
}
