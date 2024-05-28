package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/oleoneto/dm/cmd/cli/core"
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

		state.Runner.LoadMigrations(ctx)
		loaded := state.Runner.Migrations(ctx)

		applied, err := state.Runner.AppliedMigrations(ctx)
		if err != nil {
			log.Fatalln(err)
			return
		}

		if applied == nil || applied.IsEmpty() {
			fmt.Println("No applied migrations to rollback.")
			return
		}

		if filter.Value != "" {
			sequence := applied.FindSequence(func(item migrator.Migration) bool { return item.Version == filter.Value || item.Name == filter.Value })
			if sequence == nil {
				fmt.Printf("Migration with %s %s not found.\n", filter.Type, filter.Value)
				return
			}

			loaded = sequence
		}

		err = state.Runner.RevertSome(ctx, *loaded)
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
