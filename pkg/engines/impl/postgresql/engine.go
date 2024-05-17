package engines

import (
	"context"
	"database/sql"
)

type PostgreSQL struct{
    db *sql.DB
}

func (e *PostgreSQL) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) { 
    r, err := e.db.ExecContext(ctx, query, args)
    return r, err
}

func (e *PostgreSQL) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) { 
    r, err := e.db.QueryContext(ctx, query, args)
    return r, err
}

func (e *PostgreSQL) QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, error) {
    r := e.db.QueryRowContext(ctx, query, args)
    return r, nil
}
