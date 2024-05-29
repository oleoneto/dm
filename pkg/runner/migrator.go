package runner

import (
	"context"
	"fmt"

	"github.com/oleoneto/dm/pkg/engines"
	"github.com/oleoneto/dm/pkg/helpers"
	"github.com/oleoneto/dm/pkg/migrator"
	log "github.com/sirupsen/logrus"
)

// --------------------------------------------
// MARK: Migrator
// --------------------------------------------

type MigrateUpFunc func(context.Context, []migrator.Migration, engines.SqlEngineProtocol, TrackerOptions) error

type MigrateDownFunc func(context.Context, []migrator.Migration, engines.SqlEngineProtocol, TrackerOptions) error

type ValidateMigrationsFunc func(context.Context, []migrator.Migration) error

// ApplySome
func (r *Runner) ApplySome(ctx context.Context, migrations []migrator.Migration) error {
	r.migrations = migrations

	if len(r.migrations) == 0 {
		log.WithField("func", helpers.GetCurrentFuncName()).Debugln("No migrations found...")
		return nil
	}

	if err := r.validatorFunc(ctx, r.migrations); err != nil {
		return err
	}

	if err := r.StartTracking(ctx); err != nil {
		return err
	}

	return r.migrateUpFunc(ctx, r.migrations, r.engine, r.trackerOptions)
}

// Apply - Migrate up. Applies all migrations.
func (r *Runner) Apply(ctx context.Context) error {
	if err := r.LoadAllMigrations(ctx); err != nil {
		return err
	}

	if len(r.migrations) == 0 {
		log.WithField("func", helpers.GetCurrentFuncName()).Debugln("No migrations found...")
		return nil
	}

	if err := r.validatorFunc(ctx, r.migrations); err != nil {
		return err
	}

	if err := r.StartTracking(ctx); err != nil {
		return err
	}

	// Determine a starting point
	version := r.Version(ctx)

	seq := helpers.FindLeftSequence(r.migrations, func(m migrator.Migration) bool { return m.Version == version })
	if len(seq) > 0 {
		r.migrations = seq // TODO: Review this...
	}

	return r.migrateUpFunc(ctx, r.migrations, r.engine, r.trackerOptions)
}

func (r *Runner) RevertSome(ctx context.Context, migrations []migrator.Migration) error {
	r.migrations = migrations

	if len(r.migrations) == 0 {
		log.WithField("func", helpers.GetCurrentFuncName()).Debugln("No migrations found...")
		return nil
	}

	// Can't revert if database is not being tracked
	if !r.IsTracked(ctx) {
		return fmt.Errorf("database not currently tracked")
	}

	// Determine if there are any migrations applied.
	version := r.Version(ctx)
	if version == "" {
		return nil
	}

	return r.migrateDownFunc(ctx, r.migrations, r.engine, r.trackerOptions)
}

// RevertMigrations - Migrate down. Reverts the given migrations.
func (r *Runner) Revert(ctx context.Context) error {
	if len(r.migrations) == 0 {
		log.WithField("func", helpers.GetCurrentFuncName()).Debugln("No migrations found...")
		return nil
	}

	// Can't revert if database is not being tracked
	if !r.IsTracked(ctx) {
		return fmt.Errorf("database not currently tracked")
	}

	// Determine if there are any migrations applied.
	version := r.Version(ctx)
	if version == "" {
		return nil
	}

	seq := helpers.FindLeftSequence(r.migrations, func(m migrator.Migration) bool { return m.Version == version })
	if len(seq) > 0 {
		r.migrations = seq
	}

	return r.migrateDownFunc(ctx, r.migrations, r.engine, r.trackerOptions)
}

// Apply - Migrate up. Applies the given migrations.
func ApplyMigrations(ctx context.Context, migrations []migrator.Migration, engine engines.SqlEngineProtocol, trackerOptions TrackerOptions) error {
	tx, txerr := engine.BeginTx(ctx, nil)
	if txerr != nil {
		log.WithField("func", helpers.GetCurrentFuncName()).Error(txerr.Error())

		return txerr
	}

	total := len(migrations)

	for i := 0; i < len(migrations); i++ {
		var m = migrations[i]
		stmt := migrateUpStatement(m)
		_, merr := tx.ExecContext(ctx, stmt)
		if merr != nil {
			tx.Rollback()
			log.WithField("func", helpers.GetCurrentFuncName()).Error(merr.Error())
			return merr
		}

		// Track changes
		_, terr := tx.ExecContext(ctx, trackMigrationStatement(trackerOptions.Schema, trackerOptions.Table), m.Version, m.Name)
		if terr != nil {
			tx.Rollback()
			log.WithField("func", helpers.GetCurrentFuncName()).Error(terr.Error())
			return terr
		}
	}

	log.Infoln(fmt.Sprintf("%d migrations applied", total))

	return tx.Commit()
}

func RollbackMigrations(ctx context.Context, migrations []migrator.Migration, engine engines.SqlEngineProtocol, trackerOptions TrackerOptions) error {
	tx, txerr := engine.BeginTx(ctx, nil)
	if txerr != nil {
		log.WithField("func", helpers.GetCurrentFuncName()).Error(txerr.Error())
		return txerr
	}

	total := len(migrations)

	// Important that rollback occurs in the correct order (DESC)
	for i := len(migrations) - 1; i >= 0; i-- {
		var m = migrations[i]

		stmt := migrateDownStatement(m)
		if stmt == "" {
			return fmt.Errorf("empty rollback instructions")
		}

		_, merr := tx.ExecContext(ctx, stmt)
		if merr != nil {
			tx.Rollback()
			log.WithField("func", helpers.GetCurrentFuncName()).Error(merr.Error())
			return merr
		}

		// Untrack changes
		_, terr := tx.ExecContext(ctx, untrackMigrationStatement(trackerOptions.Schema, trackerOptions.Table), m.Version, m.Name)
		if terr != nil {
			tx.Rollback()
			log.WithField("func", helpers.GetCurrentFuncName()).Error(terr.Error())
			return terr
		}
	}

	log.Infoln(fmt.Sprintf("%d migrations reverted", total))

	return tx.Commit()
}
