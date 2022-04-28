package postgresql

import (
	"github.com/cleopatrio/db-migrator-lib/migrations"
)

func (engine Postgres) Validate(changes migrations.Migrations) bool {
	return migrations.Validate(changes)
}
