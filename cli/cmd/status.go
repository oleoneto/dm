package cmd

import (
	"github.com/spf13/cobra"
)

var (
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Shows all currently applied migrations",
		Run: func(cmd *cobra.Command, args []string) {
			migrator.Status()
		},
	}
)
