package engines

import (
	"context"
	"database/sql"
)

type SQLite3 struct{ db *sql.DB }

func (e *SQLite3) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
    r, err := e.db.ExecContext(ctx, query, args)
    return r, err
}

func (e *SQLite3) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
    r, err := e.db.QueryContext(ctx, query, args)
    return r, err
}


func (e *SQLite3) QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, error) {
    r := e.db.QueryRowContext(ctx, query, args)
    return r, nil
}
