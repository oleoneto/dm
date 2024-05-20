package migrator

import (
	"context"
	"errors"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/engines"
)

type MigratorProtocol interface {
	Validate(context.Context, ds.Queue[Migration]) error
	MigrateUp(context.Context, ds.Queue[Migration]) error
	MigrateDown(context.Context, ds.Queue[Migration]) error
}

type MigrationsController struct {
	engine engines.SqlEngine
}

func NewMigrationsController(engine engines.SqlEngine) *MigrationsController {
	return &MigrationsController{
		engine: engine,
	}
}

// Up - Applies the given migrations
func (ctr *MigrationsController) MigrateUp(ctx context.Context, data ds.Queue[Migration]) error {
	if err := ctr.Validate(ctx, data); err != nil {
		return err
	}

	_, err := ctr.engine.Exec(ctx, "")
	if err != nil {
		return err
	}

	return nil
}

// Down - Reverts the given migrations
func (ctr *MigrationsController) MigrateDown(ctx context.Context, data ds.Queue[Migration]) error {
	if err := ctr.Validate(ctx, data); err != nil {
		return err
	}

	_, err := ctr.engine.Exec(ctx, "")
	if err != nil {
		return err
	}

	return nil
}

// Validate - Checks if all migrations are valid
func (ctr *MigrationsController) Validate(ctx context.Context, migrations ds.Queue[Migration]) error {
	if migrations.IsEmpty() {
		return errors.New("no migrations found")
	}

	visitedNames := make(map[string]bool)
	visitedVersions := make(map[string]bool)

	node := migrations.Dequeue()

	for node != nil {
		if visitedVersions[node.Version] {
			return errors.New(`duplicate migration version`)
		}

		if visitedNames[node.Name] {
			return errors.New("duplicate migration name")
		}

		if node.Engine != ctr.engine.Name() {
			return errors.New("migration engine mismatch")
		}

		var mismatchedInstructions int
		var mismatchedTables = map[string]string{}

		// TODO: Check CREATE and DROP mismatch
		for _, change := range node.Changes.Up {
			if len(strings.Split(change, " ")) < 3 {
				return errors.New("missing (or invalid) migrate instruction")
			}

			// if migrations.CreateTablePattern.MatchString(change) {
			// 	mismatchedInstructions += 1
			// 	match := migrations.CreateTablePattern.FindStringSubmatch(change)
			// 	table := match[migrations.CreateTablePattern.SubexpIndex("TableName")]
			// 	mismatchedTables[table] = table
			// } else if migrations.DropTablePattern.MatchString(change) {
			// 	mismatchedInstructions -= 1
			// 	match := migrations.DropTablePattern.FindStringSubmatch(change)
			// 	table := match[migrations.DropTablePattern.SubexpIndex("TableName")]
			// 	delete(mismatchedTables, table)
			// }
		}

		if mismatchedInstructions != 0 || len(mismatchedTables) != 0 {
			return errors.New("CREATE and DROP instructions must always be paired")
		}

		version, name, _ := strings.Cut(node.FileName, "_")
		name = strcase.ToCamel(strings.Split(name, ".")[0])

		if node.Version != version {
			return errors.New("version mismatch")
		}

		if node.Name != name {
			return errors.New("name mismatch")
		}

		visitedNames[node.Name] = true
		visitedVersions[node.Version] = true

		node = migrations.Dequeue()
	}

	return nil
}
