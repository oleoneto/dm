package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/drewstinnett/gout/v2"
	"github.com/drewstinnett/gout/v2/config"
	"github.com/drewstinnett/gout/v2/formats/gotemplate"
	gJSON "github.com/drewstinnett/gout/v2/formats/json"
	gPlain "github.com/drewstinnett/gout/v2/formats/plain"
	gYAML "github.com/drewstinnett/gout/v2/formats/yaml"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/oleoneto/dm/pkg/helpers"
	"github.com/oleoneto/dm/pkg/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type FlagEnum struct {
	Allowed []string
	Default string
}

type CommandFlags struct {
	OutputTemplate string
	OutputFormat   *FlagEnum
	Engine         *FlagEnum
	Extension      *FlagEnum
	Table          string
	Directory      string
	DatabaseURL    *string
}

type CommandState struct {
	Writer             *gout.Gout
	Flags              CommandFlags
	ExecutionStartTime time.Time
	ExecutionExitLog   []any
	Database           *sql.DB
	Runner             *runner.Runner
}

type TableFormattable interface{ TableWriter() table.Writer }

type TableFormatter struct{}

type SilentFormatter struct{}

func (f *TableFormatter) Format(data any) ([]byte, error) {
	tw, ok := data.(TableFormattable)
	if !ok {
		return []byte{}, nil
	}

	return []byte(tw.TableWriter().Render()), nil
}

func (f *SilentFormatter) Format(data any) ([]byte, error) { return []byte{}, nil }

func (ofe FlagEnum) String() string { return ofe.Default }

func (ofe *FlagEnum) Type() string { return "string" }

func (ofe *FlagEnum) Set(value string) error {
	isIncluded := func(opts []string, v string) bool {
		for _, opt := range opts {
			if v == opt {
				return true
			}
		}

		return false
	}

	if !isIncluded(ofe.Allowed, value) {
		return fmt.Errorf("%s is not a supported output format: %s", value, strings.Join(ofe.Allowed, ","))
	}

	ofe.Default = value
	return nil
}

var _ pflag.Value = (*FlagEnum)(nil)

func (c *CommandState) SetFormatter(cmd *cobra.Command, args []string) {
	switch cmd.Flag("output").Value.String() {
	case "table":
		c.Writer.SetFormatter(&TableFormatter{})
	case "json":
		c.Writer.SetFormatter(gJSON.Formatter{})
	case "yaml":
		c.Writer.SetFormatter(gYAML.Formatter{})
	case "gotemplate":
		c.Writer.SetFormatter(gotemplate.Formatter{
			Opts: config.FormatterOpts{"template": c.Flags.OutputTemplate},
		})
	case "silent":
		c.Writer.SetFormatter(&SilentFormatter{})
	case "plain":
		c.Writer.SetFormatter(gPlain.Formatter{})
	default:
		c.Writer.SetFormatter(gJSON.Formatter{})
	}
}

func (c *CommandState) ConnectDatabase(cmd *cobra.Command, args []string) {
	if c.Flags.DatabaseURL == nil || *c.Flags.DatabaseURL == "" {
		log.Fatalln("database-url not set")
		return
	}

	switch cmd.Flag("adapter").Value.String() {
	case "postgresql":
		var err error
		c.Database, err = sql.Open("pgx", *c.Flags.DatabaseURL)
		if err != nil {
			log.Fatal(err)
			return
		}
	case "sqlite3":
		var err error
		c.Database, err = sql.Open("sqlite3_extended", *c.Flags.DatabaseURL)
		if err != nil {
			log.Fatal(err)
			return
		}
	default:
		log.Fatal("database adapter not set")
		return
	}

	c.Runner = runner.NewRunner(
		&runner.TrackerOptions{
			Schema:              "public",
			Table:               c.Flags.Table,
			MigrationsDirectory: c.Flags.Directory,
		},
		c.Database,
		runner.LoadFiles,
		runner.LoadMigrations,
		runner.Validate,
		runner.ApplyMigrations,
		runner.RollbackMigrations,
	)
}

func (c *CommandState) BeforeHook(cmd *cobra.Command, args []string) {
	c.ExecutionStartTime = time.Now()

	c.SetFormatter(cmd, args)
}

func (c *CommandState) AfterHook(cmd *cobra.Command, args []string) {
	fmt.Fprintln(
		os.Stderr,
		append([]any{"Elapsed time:", time.Since(c.ExecutionStartTime)}, c.ExecutionExitLog...)...,
	)
}

func NewCommandState() *CommandState {
	command := CommandState{
		Writer: gout.New(),
		Flags: CommandFlags{
			OutputFormat: &FlagEnum{
				Allowed: []string{"plain", "json", "yaml", "table", "gotemplate", "silent"},
				Default: "json",
			},
			Engine: &FlagEnum{
				Allowed: []string{"postgresql", "sqlite3"},
				Default: "postgresql",
			},
			Extension: &FlagEnum{
				Allowed: []string{"yaml", "sql"},
				Default: "yaml",
			},
			Table:       "_migrations",
			Directory:   "./migrations",
			DatabaseURL: helpers.PointerTo(os.Getenv("DATABASE_URL")),
		},
	}

	return &command
}
