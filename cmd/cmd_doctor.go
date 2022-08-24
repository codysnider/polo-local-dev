package cmd

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-version"
	"github.com/poloniex/polo-local-dev/output"
	"github.com/spf13/cobra"
	"os"
)

var doctorDockerMinVersion = "18.09.00"

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

		output.Title("Doctor")

		output.Section("Environment Variables")

		for envVarKey, envVarVal := range doctorEnvironmentVariables {
			if os.Getenv(envVarKey) == "" {
				output.Warning(fmt.Sprintf("doctor env var missing: %s", envVarKey))
				continue
			}
			if os.Getenv(envVarKey) != envVarVal {
				output.Warning(fmt.Sprintf("doctor env var: %s, got: %s, expected: %s", envVarKey, os.Getenv(envVarKey), envVarVal))
				continue
			}

			output.Ok(envVarKey)
		}

		output.Section("Docker")

		cli, dockerClientErr := client.NewClientWithOpts(client.FromEnv)
		if dockerClientErr != nil {
			output.Error(dockerClientErr.Error())
			return
		}

		installedVer, versionFetchErr := cli.ServerVersion(context.Background())
		if versionFetchErr != nil {
			output.Error(versionFetchErr.Error())
			return
		}

		minDockerVersion, minDockerVersionParseErr := version.NewVersion(doctorDockerMinVersion)
		if minDockerVersionParseErr != nil {
			output.Error(minDockerVersionParseErr.Error())
			return
		}

		installedVersion, installedVersionParseErr := version.NewVersion(installedVer.Version)
		if installedVersionParseErr != nil {
			output.Error(installedVersionParseErr.Error())
			return
		}

		if installedVersion.GreaterThanOrEqual(minDockerVersion) {
			output.Ok(fmt.Sprintf("Version: %s >= %s", installedVersion, minDockerVersion))
		}
	},
}
