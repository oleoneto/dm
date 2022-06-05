package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/iancoleman/strcase"
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
				fmt.Fprintf(os.Stderr, "Unsupported adapter '%v'.\n", adapter)
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
					fmt.Println("Error: a migration with this name already exists.")
					os.Exit(1)
				}
			}

			migration := runner.Generate(version.Value, directory)

			if migration.FileName == "" {
				fmt.Println("Error: migration file not created.")
				os.Exit(1)
			}
		},
	}
)
