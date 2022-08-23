package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var doctorEnvironmentVariables = map[string]string{
	"APP_ENV":         "local-west",
	"PROFILE":         "local-west",
	"SPOT_MYSQL_PORT": "3307",
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Validate environment",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		for envVarKey, envVarVal := range doctorEnvironmentVariables {
			if os.Getenv(envVarKey) != envVarVal {

			}
		}
	},
}
