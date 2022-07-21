package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/oleoneto/dm/logger"
	"github.com/oleoneto/dm/migrations"
	"github.com/spf13/cobra"
)

var (
	SUPPORTED_FORMATS = map[string]int{"yaml": 0, "sql": 1}
	mformat           = "yaml"
	filecontent       = ""
	readStdin         = false

	generateCmd = &cobra.Command{
		Use:   "generate NAME",
		Short: "Generate a database migration file in the migrations directory",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			selectedAdapter, ok := SUPPORTED_ADAPTERS[adapter]

			if !ok {
				message := logger.ApplicationError{Error: fmt.Sprintf("Unsupported adapter '%v'", adapter)}
				logger.Custom(format, template).WithFormattedOutput(&message, os.Stderr)
				os.Exit(1)
			}

			storeAdapter = selectedAdapter
			runner.SetStore(storeAdapter)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// MARK: Guard
			_, ok := SUPPORTED_FORMATS[mformat]

			if !ok {
				message := logger.ApplicationError{Error: "Error: unsupported format."}
				logger.Custom(format, template).WithFormattedOutput(&message, os.Stderr)
				return
			}

			version, err := parsedVersionFlag(args[0])

			if err != nil || version.Type != "Name" {
				os.Exit(1)
			}

			files := migrations.LoadFiles(directory, &FilePattern)

			for _, file := range files {
				if strings.Contains(strcase.ToSnake(file.Name()), strcase.ToSnake(version.Value)) {
					message := logger.ApplicationError{Error: "Error: migration with this name already exists."}
					logger.Custom(format, template).WithFormattedOutput(&message, os.Stderr)
					os.Exit(1)
				}
			}

			if readStdin {
				filecontent, _ = readFromStdin()
			}

			migration := runner.Generate(mformat, filecontent, version.Value, directory)

			if migration.FileName == "" {
				message := logger.ApplicationError{Error: "Error: migration file not created."}
				logger.Custom(format, template).WithFormattedOutput(&message, os.Stderr)
				os.Exit(1)
			}
		},
	}
)

func init() {
	generateCmd.PersistentFlags().StringVar(&mformat, "format", mformat, "migration file format")
	generateCmd.PersistentFlags().StringVar(&filecontent, "content", filecontent, "file content [can be read from stdin]")
	generateCmd.PersistentFlags().BoolVar(&readStdin, "stdin", readStdin, "read input from stdin")
}
