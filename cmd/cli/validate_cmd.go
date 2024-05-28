package cli

import (
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:               "validate",
	Short:             "Validate the configuration of migration files",
	PersistentPreRun:  state.BeforeHook,
	PersistentPostRun: state.AfterHook,
	Run:               func(cmd *cobra.Command, args []string) {},
}

type validationOutput struct {
	Message string
	Valid   bool
}

func (v validationOutput) Description() string { return v.Message }
