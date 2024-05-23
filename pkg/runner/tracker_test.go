package runner_test

import (
	"context"
	"io/fs"
	"testing"

	"github.com/oleoneto/dm/pkg/engines"
	"github.com/oleoneto/dm/pkg/fsystem"
	"github.com/oleoneto/dm/pkg/runner"
	"github.com/stretchr/testify/assert"
)

func TestRunner_IsTracked(t *testing.T) {
	type fields struct {
		engine         engines.SqlEngineProtocol
		loader         fsystem.FileLoaderProtocol
		trackerOptions runner.TrackerOptions
		// migrations          ds.Queue[migrator.Migration]
		migrationLoaderFunc runner.LoadMigrationsFunc
		validatorFunc       runner.ValidateMigrationsFunc
		migrateUpFunc       runner.MigrateUpFunc
		migrateDownFunc     runner.MigrateDownFunc
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "untracked changes - 1",
			fields: fields{
				loader:         Loader{data: []fs.DirEntry{}},
				trackerOptions: runner.TrackerOptions{},
			},
			args: args{
				ctx: context.TODO(),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := runner.NewRunner(
				tt.fields.engine,
				tt.fields.loader,
				tt.fields.trackerOptions,
				tt.fields.migrationLoaderFunc,
				tt.fields.validatorFunc,
				tt.fields.migrateUpFunc,
				tt.fields.migrateDownFunc,
			)

			isTracked := r.IsTracked(tt.args.ctx)
			assert.Equal(t, isTracked, tt.want)

			// if got := r.IsTracked(tt.args.ctx); got != tt.want {
			// 	t.Errorf("Runner.IsTracked() = %v, want %v", got, tt.want)
			// }
		})
	}
}
