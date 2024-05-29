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

var RollbackCmd = &cobra.Command{
	Use:     "rollback NAME|VERSION",
	Short:   "Rollback migration(s)",
	Aliases: []string{"r"},
	Args:    cobra.MaximumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		state.BeforeHook(cmd, args)
		state.ConnectDatabase(cmd, args)
	},
	PersistentPostRun: state.AfterHook,
	Run: func(cmd *cobra.Command, args []string) {
		var filter core.MigrationFilterFlag
		var err error

		if len(args) > 0 && args[0] != "" {
			filter, err = core.ParseVersionArgs(args[0])
			if err != nil {
				log.Fatalln(err)
				return
			}
		}

		ctx := context.TODO()

		applied, err := state.Runner.AppliedMigrations(ctx)
		if err != nil {
			log.Fatalln(err)
			return
		}

		if len(applied) == 0 {
			fmt.Println("No applied migrations to rollback.")
			return
		}

		if filter.Value != "" {
			sequence := helpers.FindLeftSequence(applied, func(item migrator.Migration) bool { return item.Version == filter.Value || item.Name == filter.Value })
			if sequence == nil {
				fmt.Printf("Migration with %s %s not found.\n", filter.Type, filter.Value)
				return
			}

			applied = sequence
		}

		var keys = make(map[string]int)
		for _, m := range applied {
			keys[m.Version] = 0
		}

		state.Runner.LoadMigrations(ctx, keys)

		err = state.Runner.RevertSome(ctx, state.Runner.Migrations(ctx))
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RollbackCmd.PersistentFlags().StringVarP(state.Flags.DatabaseURL, "database-url", "u", *state.Flags.DatabaseURL, "database url")
	RollbackCmd.MarkFlagRequired("database-url")
	RollbackCmd.MarkFlagRequired("adapter")
	RollbackCmd.MarkFlagRequired("table")
}
