package engines

import (
	"fmt"
	"log"

	"github.com/cleopatrio/db-migrator-lib/migrations"
)

type Postgres struct{}

func (engine Postgres) Up(changes []migrations.Migration) error {
	// TODO: Implement database-specific behavior

	if engine.IsUpToDate() {
		log.Println("Nothing to do. Migrations are up-to-date.")
		return nil
	}

	fmt.Println("PostgreSQL: Running migrations up")
	return nil
}

func (engine Postgres) Down(changes []migrations.Migration) error {
	// TODO: Implement database-specific behavior
	fmt.Println("PostgreSQL: Running migrations down")
	return nil
}

func (engine Postgres) IsUpToDate() bool {
	// TODO: Implement behavior
	return false
}
