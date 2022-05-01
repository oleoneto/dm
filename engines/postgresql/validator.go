package postgresql

import (
	"github.com/cleopatrio/db-migrator-lib/migrations"
)

func (engine Postgres) Validate(changes migrations.MigrationList) (bool, string) {
	return migrations.Validate(changes)
}
