package engines

import (
	"context"
	"database/sql"
)

type SQLite3 struct{ db *sql.DB }

func (e *SQLite3) Exec(context.Context, string, ...any) (sql.Result, error) {
    r, err := e.db.ExecContext(ctx, query, args)
    return r, err
}

func (e *SQLite3) Query(context.Context, string, ...any) (*sql.Rows, error) {
    r, err := e.db.QueryContext(ctx, query, args)
    return r, err
}


func (e *SQLite3) QueryRow(context.Context, string, ...any) (*sql.Row, error) {
    r := e.db.QueryRowContext(ctx, query, args)
    return r, nil
}
