package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/oleoneto/dm/cmd/cli/core"
	"github.com/oleoneto/dm/pkg/migrator"
	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
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

		var version core.MigrationFilterFlag
		var err error

		if len(args) > 0 && args[0] != "" {
			version, err = core.ParseVersionArgs(args[0])
			if err != nil {
				log.Fatalln(err)
				return
			}
		}

		migrations, err := state.Runner.PendingMigrations(ctx)
		if err != nil {
			log.Fatalln(err)
			return
		}

		if migrations == nil || migrations.IsEmpty() {
			fmt.Println("No pending migrations to apply.")
			return
		}

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

		err = state.Runner.ApplySome(ctx, *migrations)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	MigrateCmd.PersistentFlags().StringVarP(state.Flags.DatabaseURL, "database-url", "u", *state.Flags.DatabaseURL, "database url")
	MigrateCmd.MarkFlagRequired("database-url")
	MigrateCmd.MarkFlagRequired("adapter")
	MigrateCmd.MarkFlagRequired("table")
}
