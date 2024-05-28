package cli

import (
	"context"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/oleoneto/dm/cmd/cli/core"
	"github.com/oleoneto/dm/pkg/migrator"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:               "show",
	Short:             "Shows the state of applied and pending migrations",
	PersistentPreRun:  state.BeforeHook,
	PersistentPostRun: state.AfterHook,
	Run:               func(cmd *cobra.Command, args []string) { cmd.Help() },
}

var (
	allCmd = &cobra.Command{
		Use:    "all",
		Short:  "List all migrations for a given application",
		PreRun: state.ConnectDatabase,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.TODO()

			state.Runner.LoadMigrations(ctx)

			migrations := state.Runner.Migrations(ctx)
			if migrations != nil {
				state.Writer.Print(Migrations(*migrations.RawData()))
			}
		},
	}

	appliedCmd = &cobra.Command{
		Use:    "applied",
		Short:  "List only applied migrations",
		PreRun: state.ConnectDatabase,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.TODO()

			migrations, err := state.Runner.AppliedMigrations(ctx)
			if err != nil {
				state.Writer.Print(err)
			}

			if migrations != nil {
				state.Writer.Print(Migrations(*migrations.RawData()))
			}
		},
	}

	pendingCmd = &cobra.Command{
		Use:     "pending",
		Short:   "List only pending migrations",
		Aliases: []string{"p"},
		PreRun:  state.ConnectDatabase,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.TODO()

			migrations, err := state.Runner.PendingMigrations(ctx)
			if err != nil {
				state.Writer.Print(err)
			}

			if migrations != nil {
				state.Writer.Print(Migrations(*migrations.RawData()))
			}
		},
	}

	migrationVersionCmd = &cobra.Command{
		Use:    "version",
		Short:  "Shows the most recently applied migration",
		PreRun: state.ConnectDatabase,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.TODO()

			if version := state.Runner.Version(ctx); version != "" {
				state.Writer.Print(version)
				return
			}
		},
	}
)

func init() {
	showCmd.AddCommand(allCmd)
	showCmd.AddCommand(appliedCmd)
	showCmd.AddCommand(pendingCmd)
	showCmd.AddCommand(migrationVersionCmd)

	showCmd.PersistentFlags().StringVarP(&state.Flags.DatabaseURL, "database-url", "u", "", "database url")
	showCmd.MarkFlagRequired("database-url")
	showCmd.MarkFlagRequired("adapter")
	showCmd.MarkFlagRequired("table")
}

type Migrations []migrator.Migration

var _ core.TableFormattable = (*Migrations)(nil)

func (data Migrations) TableWriter() table.Writer {
	if data == nil {
		return nil
	}

	t := core.Initialize(core.TableOptions{
		OutputMirror: nil, // Delegate printing to gout tool
		ColumnConfig: &[]table.ColumnConfig{
			{Name: "name", Align: text.AlignLeft},
			{Name: "version", Align: text.AlignLeft},
		},
		Header: table.Row{
			"name",
			"version",
		},
		Footer: &table.Row{"Total", len(data)},
	})

	for _, item := range data {
		t.AppendRow(table.Row{
			item.Name,
			item.Version,
		})
	}

	return t
}
