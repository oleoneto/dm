package cmd

import (
	"fmt"
	"os"

	"github.com/cleopatrio/db-migrator-lib/api/server"
	c "github.com/cleopatrio/db-migrator-lib/config"
	"github.com/cleopatrio/db-migrator-lib/logger"
	"github.com/spf13/cobra"
)

var (
	serverPort         = 3809
	apiNamespacePrefix = "migrations"
	apiVersionPrefix   = "v1"
	apiConfig          = c.DMConfig{}
	apiDebugMode       = false

	apiCmd = &cobra.Command{
		Use:   "api",
		Short: "Run a RESTful API",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			validateDatabaseConfig()
			overrideVariablesFromEnvironment()

			apiConfig = c.DMConfig{
				ConnectionString: databaseUrl,
				Table:            table,
				Directory:        directory,
			}
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("⚡️ Server running on port", serverPort)
		},
		Run: func(cmd *cobra.Command, args []string) {
			app := server.API(apiVersionPrefix, apiNamespacePrefix, apiDebugMode, apiConfig)
			err := app.Run(fmt.Sprintf(":%v", serverPort))

			if err != nil {
				message := ErrorOutput{Error: err.Error()}
				logger.Custom(format, template).WithFormattedOutput(&message, os.Stderr)
				return
			}
		},
	}
)

func init() {
	apiCmd.PersistentFlags().StringVarP(&databaseUrl, "database-url", "u", databaseUrl, "database url")
	apiCmd.PersistentFlags().IntVarP(&serverPort, "port", "p", serverPort, "server port")
	apiCmd.PersistentFlags().StringVarP(&apiVersionPrefix, "version", "v", apiVersionPrefix, "api version param")
	apiCmd.PersistentFlags().StringVarP(&apiNamespacePrefix, "namespace", "n", apiNamespacePrefix, "api resource namespace param")
	apiCmd.PersistentFlags().BoolVar(&apiDebugMode, "debug", apiDebugMode, "shows server debug output")

	apiCmd.MarkFlagRequired("database-url")
	apiCmd.MarkFlagRequired("directory")
	apiCmd.MarkFlagRequired("table")
}
