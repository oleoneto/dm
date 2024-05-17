package engines

import (
	"context"
	"database/sql"
)

type PostgreSQL struct{}

func (e *PostgreSQL) Exec(context.Context, string, ...any) (sql.Result, error) { return nil, nil }

func (e *PostgreSQL) Query(context.Context, string, ...any) (*sql.Rows, error) { return nil, nil }

func (e *PostgreSQL) QueryRow(context.Context, string, ...any) *sql.Row { return nil }
