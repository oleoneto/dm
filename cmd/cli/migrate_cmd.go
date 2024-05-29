package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/oleoneto/dm/cmd/cli/core"
	"github.com/oleoneto/dm/pkg/helpers"
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

		var filter core.MigrationFilterFlag
		var err error

		if len(args) > 0 && args[0] != "" {
			filter, err = core.ParseVersionArgs(args[0])
			if err != nil {
				log.Fatalln(err)
				return
			}
		}

		pending, err := state.Runner.PendingMigrations(ctx)
		if err != nil {
			log.Fatalln(err)
			return
		}

		if len(pending) == 0 {
			fmt.Println("No pending migrations to apply.")
			return
		}

		if filter.Value != "" {
			sequence := helpers.FindLeftSequence(pending, func(item migrator.Migration) bool {
				return item.Version == filter.Value || item.Name == filter.Value
			})

			if sequence == nil {
				fmt.Printf("Migration with %s %s not found.\n", filter.Type, filter.Value)
				return
			}

			pending = sequence
		}

		err = state.Runner.ApplySome(ctx, pending)
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
