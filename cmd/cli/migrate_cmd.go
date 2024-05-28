package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/oleoneto/dm/cmd/cli/core"
	"github.com/oleoneto/dm/pkg/migrator"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:     "migrate NAME|VERSION",
	Short:   "Run migration(s)",
	Aliases: []string{"m"},
	Args:    cobra.MaximumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		state.BeforeHook(cmd, args)
		state.ConnectDatabase(cmd, args)
	},
	PersistentPostRun: state.AfterHook,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.TODO()

		var version core.VersionFlag
		var err error

		if len(args) > 0 && args[0] != "" {
			version, err = core.ParseVersionArgs(args[0])
			if err != nil {
				log.Fatal(err)
			}
		}

		migrations, err := state.Runner.PendingMigrations(ctx)

		if version.Value != "" {
			sequence := migrations.FindSequence(func(item migrator.Migration) bool {
				return item.Version == version.Value
			})

			if sequence == nil {
				fmt.Println("Nothing to do.")
				return
			}

			migrations = sequence
		}

		state.Runner.Apply(ctx, migrations)
	},
}

func init() {
	migrateCmd.PersistentFlags().StringVarP(&state.Flags.DatabaseURL, "database-url", "u", "", "database url")
	migrateCmd.MarkFlagRequired("database-url")
	migrateCmd.MarkFlagRequired("adapter")
	migrateCmd.MarkFlagRequired("table")
}
