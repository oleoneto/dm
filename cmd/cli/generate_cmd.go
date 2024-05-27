package cli

import (
	"github.com/spf13/cobra"
)

var (
	filecontent       = ""
	readStdin         = false
)

var generateCmd = &cobra.Command{
	Use:   "generate NAME",
	Short: "Generate a database migration file in the migrations directory",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// selectedAdapter, ok := SUPPORTED_ADAPTERS[adapter]

		// if !ok {
		// 	message := logger.ApplicationError{Error: fmt.Sprintf("Unsupported adapter '%v'", adapter)}
		// 	logger.Custom(format, template).WithFormattedOutput(&message, os.Stderr)
		// 	os.Exit(1)
		// }

		// storeAdapter = selectedAdapter
		// runner.SetStore(storeAdapter)
	},
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	generateCmd.PersistentFlags().Var(state.Flags.Extension, "format", "migration file format")
	generateCmd.PersistentFlags().StringVar(&filecontent, "content", filecontent, "file content [can be read from stdin]")
	generateCmd.PersistentFlags().BoolVar(&readStdin, "stdin", readStdin, "read input from stdin")
}
