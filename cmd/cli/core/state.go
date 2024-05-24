package core

import (
	"fmt"
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
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type OutputFormatEnum struct {
	Allowed []string
	Default string
}

type CommandFlags struct {
	OutputTemplate string
	OutputFormat   *OutputFormatEnum
}

type CommandState struct {
	Writer             *gout.Gout
	Flags              CommandFlags
	ExecutionStartTime time.Time
	ExecutionExitLog   []any
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

func (ofe OutputFormatEnum) String() string { return ofe.Default }

func (ofe *OutputFormatEnum) Type() string { return "string" }

func (ofe *OutputFormatEnum) Set(value string) error {
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

var _ pflag.Value = (*OutputFormatEnum)(nil)

func NewCommandState() CommandState {
	command := CommandState{
		Writer: gout.New(),
		Flags: CommandFlags{
			OutputFormat: &OutputFormatEnum{
				Allowed: []string{"plain", "json", "yaml", "table", "gotemplate", "silent"},
				Default: "json",
			},
		},
	}

	return command
}

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
	default:
		c.Writer.SetFormatter(gPlain.Formatter{})
	}
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
