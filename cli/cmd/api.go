package cmd

import (
	"fmt"
	"os"

	"github.com/oleoneto/dm/api/server"
	c "github.com/oleoneto/dm/config"
	"github.com/oleoneto/dm/logger"
	"github.com/spf13/cobra"
)

var (
	serverPort         = 3809
	apiNamespacePrefix = "migrations"
	apiVersionPrefix   = "v1"
	apiConfig          = c.APIConfig{}
	apiDebugMode       = false
	apiHost            = "*"

	apiCmd = &cobra.Command{
		Use:   "api",
		Short: "Run a RESTful API",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			validateDatabaseConfig()
			overrideVariablesFromEnvironment()

			apiConfig = c.APIConfig{
				AllowedHost:      apiHost,
				ConnectionString: databaseUrl,
				DebugMode:        apiDebugMode,
				Directory:        directory,
				Namespace:        apiNamespacePrefix,
				Table:            table,
				Version:          apiVersionPrefix,
			}
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("⚡️ Server running on port", serverPort)
		},
		Run: func(cmd *cobra.Command, args []string) {
			app := server.API(apiConfig)
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
	apiCmd.PersistentFlags().StringVar(&apiHost, "hosts", apiHost, "allowed CORS origins")
	apiCmd.PersistentFlags().BoolVar(&apiDebugMode, "debug", apiDebugMode, "shows server debug output")

	apiCmd.MarkFlagRequired("database-url")
	apiCmd.MarkFlagRequired("directory")
	apiCmd.MarkFlagRequired("table")
}
