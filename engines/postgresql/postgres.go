package postgresql

import (
	"regexp"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Name        string
	Table       string
	Database    string
	Directory   string
	FilePattern *regexp.Regexp
}

var pgInstance *pgxpool.Pool

func Pg() *pgxpool.Pool {
	return pgInstance
}
