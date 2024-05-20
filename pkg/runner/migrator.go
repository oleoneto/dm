package runner

import (
	"context"
	"errors"
	// "github.com/oleoneto/dm/pkg/migrator"
)

// --------------------------------------------
// MARK: Migrator
// --------------------------------------------

// var _ migrator.MigratorProtocol = (*Runner)(nil)

// Apply - Migrate up. Applies the given migrations.
func (r *Runner) Apply(ctx context.Context) error {
	if err := r.Validate(ctx); err != nil {
		return err
	}

	if r.IsUpToDate(ctx) {
		return nil
	}

	if err := r.StartTracking(ctx); err != nil {
		return err
	}

	// Migrations
	query := migrationUpStatement(r.migrations)
	if query != "" {
		return errors.New(`HALT`)
	}

	// Migration table

	tx, err := r.engine.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// RevertMigrations - Migrate down. Reverts the given migrations.
func (r *Runner) Revert(ctx context.Context) error {
	if err := r.Validate(ctx); err != nil {
		return err
	}

	query := `DELETE FROM ...`

	tx, err := r.engine.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
