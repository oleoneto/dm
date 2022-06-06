package cmd

import (
	"github.com/cleopatrio/db-migrator-lib/migrations"
	"github.com/spf13/cobra"
)

type ValidationOutput struct {
	Message string
	Valid   bool
}

func (v ValidationOutput) Description() string {
	return v.Message
}

var (
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate the configuration of migration files",
		Run: func(cmd *cobra.Command, args []string) {
			files := migrations.LoadFiles(directory, &FilePattern)

			if len(files) == 0 {
				validationOutput := &ValidationOutput{Message: "No migrations found.", Valid: false}
				WithFormattedOutput(validationOutput)
				return
			}

			list := migrations.BuildMigrations(files, directory, &FilePattern)
			valid, reason := migrations.Validate(list)

			if valid && reason == "" {
				validationOutput := &ValidationOutput{Message: "Migrations are valid.", Valid: valid}
				WithFormattedOutput(validationOutput)
			}

			validationOutput := &ValidationOutput{Message: reason, Valid: valid}
			WithFormattedOutput(validationOutput)
		},
	}
)
