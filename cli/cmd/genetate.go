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
			version, err := parsedVersionFlag(args[0])

			if err != nil || version.Type != "Name" {
				os.Exit(1)
			}

			files := migrations.LoadFiles(directory, &FilePattern)

			for _, file := range files {
				if strings.Contains(strcase.ToCamel(file.Name()), version.Value) {
					message := logger.ApplicationError{Error: "Error: migration with this name already exists."}
					logger.Custom(format, template).WithFormattedOutput(&message, os.Stderr)
					os.Exit(1)
				}
			}

			migration := runner.Generate(version.Value, directory)

			if migration.FileName == "" {
				message := logger.ApplicationError{Error: "Error: migration file not created."}
				logger.Custom(format, template).WithFormattedOutput(&message, os.Stderr)
				os.Exit(1)
			}
		},
	}
)
