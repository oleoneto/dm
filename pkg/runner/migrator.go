package runner

import (
	"context"
	"fmt"

	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/engines"
	"github.com/oleoneto/dm/pkg/helpers"
	"github.com/oleoneto/dm/pkg/migrator"
	log "github.com/sirupsen/logrus"
)

// --------------------------------------------
// MARK: Migrator
// --------------------------------------------
//

type MigrateUpFunc func(context.Context, ds.Queue[migrator.Migration], engines.SqlEngineProtocol, TrackerOptions) error

type MigrateDownFunc func(context.Context, ds.Queue[migrator.Migration], engines.SqlEngineProtocol, TrackerOptions) error

type ValidateMigrationsFunc func(context.Context, ds.Queue[migrator.Migration]) error

// ApplySome
func (r *Runner) ApplySome(ctx context.Context, migrations ds.Queue[migrator.Migration]) error {
	r.migrations = migrations

	if r.migrations.IsEmpty() {
		log.Debugln("No migrations provided", helpers.GetCurrentFuncName())
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
	if err := r.LoadMigrations(ctx); err != nil {
		return err
	}

	if r.migrations.IsEmpty() {
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
	seq := r.migrations.FindSequence(func(m migrator.Migration) bool { return m.Version == version })
	if seq != nil {
		seq.Dequeue() // Pop last applied
		return r.migrateUpFunc(ctx, *seq, r.engine, r.trackerOptions)
	}

	return r.migrateUpFunc(ctx, r.migrations, r.engine, r.trackerOptions)
}

func (r *Runner) RevertSome(ctx context.Context, migrations ds.Queue[migrator.Migration]) error {
	r.migrations = migrations

	if r.migrations.IsEmpty() {
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

	reversedQueue := r.migrations.Reversed()
	return r.migrateDownFunc(ctx, reversedQueue, r.engine, r.trackerOptions)
}

// RevertMigrations - Migrate down. Reverts the given migrations.
func (r *Runner) Revert(ctx context.Context) error {
	if r.migrations.IsEmpty() {
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

	seq := r.migrations.FindSequence(func(m migrator.Migration) bool { return m.Version == version })
	if seq != nil {
		reversedQueue := seq.Reversed() // Important so rollback occurs in the correct order (DESC)
		return r.migrateDownFunc(ctx, reversedQueue, r.engine, r.trackerOptions)
	}

	reversedQueue := r.migrations.Reversed()
	return r.migrateDownFunc(ctx, reversedQueue, r.engine, r.trackerOptions)
}

// Apply - Migrate up. Applies the given migrations.
func ApplyMigrations(ctx context.Context, migrations ds.Queue[migrator.Migration], engine engines.SqlEngineProtocol, trackerOptions TrackerOptions) error {
	tx, txerr := engine.BeginTx(ctx, nil)
	if txerr != nil {
		log.Error(ctx, txerr.Error(), helpers.GetCurrentFuncName())

		return txerr
	}

	total := migrations.Size()

	var m = migrations.Dequeue()
	for m != nil {
		stmt := migrateUpStatement(*m)
		_, merr := tx.ExecContext(ctx, stmt)
		if merr != nil {
			tx.Rollback()
			log.Error(merr.Error(), helpers.GetCurrentFuncName())
			return merr
		}

		// Track changes
		_, terr := tx.ExecContext(ctx, trackMigrationStatement(trackerOptions.Schema, trackerOptions.Table), m.Version, m.Name)
		if terr != nil {
			tx.Rollback()
			log.Error(terr.Error(), helpers.GetCurrentFuncName())
			return terr
		}

		m = migrations.Dequeue()
	}

	log.Infoln(fmt.Sprintf("%d migrations applied", total))

	return tx.Commit()
}

func RollbackMigrations(ctx context.Context, migrations ds.Queue[migrator.Migration], engine engines.SqlEngineProtocol, trackerOptions TrackerOptions) error {
	tx, txerr := engine.BeginTx(ctx, nil)
	if txerr != nil {
		log.Error(ctx, txerr.Error(), helpers.GetCurrentFuncName())
		return txerr
	}

	total := migrations.Size()

	var m = migrations.Dequeue()
	for m != nil {
		stmt := migrateDownStatement(*m)
		_, merr := tx.ExecContext(ctx, stmt)
		if merr != nil {
			tx.Rollback()
			log.Error(merr.Error(), helpers.GetCurrentFuncName())
			return merr
		}

		// Track changes
		_, terr := tx.ExecContext(ctx, untrackMigrationStatement(trackerOptions.Schema, trackerOptions.Table), m.Version, m.Name)
		if terr != nil {
			tx.Rollback()
			log.Error(terr.Error(), helpers.GetCurrentFuncName())
			return terr
		}

		m = migrations.Dequeue()
	}

	log.Infoln(fmt.Sprintf("%d migrations reverted", total))

	return tx.Commit()
}
