package engines

import (
	"context"
	"database/sql"
)

type SQLite3 struct{}

func (e *SQLite3) Exec(context.Context, string, ...any) (sql.Result, error) { return nil, nil }

func (e *SQLite3) Query(context.Context, string, ...any) (*sql.Rows, error) { return nil, nil }

func (e *SQLite3) QueryRow(context.Context, string, ...any) *sql.Row { return nil }
