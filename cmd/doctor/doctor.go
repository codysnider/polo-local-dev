package doctor

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-version"
	"github.com/poloniex/polo-local-dev/cmd/util/aws"
	"github.com/poloniex/polo-local-dev/git"
	"github.com/poloniex/polo-local-dev/output"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"regexp"
)

var (
	doctorDockerMinVersion = "18.09.00"
	doctorPythonMinVersion = "3.8"
)

var doctorEnvironmentVariables = map[string]string{
	"APP_ENV":         "local-west",
	"PROFILE":         "local-west",
	"SPOT_MYSQL_PORT": "3307",
}

var Command = &cobra.Command{
	Use:   "doctor",
	Short: "Validate environment",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		output.Title("Doctor")

		output.Section("AWS Auth")

		authOkMsg, authErr := aws.GetCallerIdentity()
		if authErr != nil {
			output.Warning(authErr.Error())
		} else {
			output.Ok(authOkMsg)
		}

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
		} else {
			output.Warning(fmt.Sprintf("Version: %s < %s", installedVersion, minDockerVersion))
		}

		output.Section("Python")

		pythonVersionCommand := exec.Command("python", "--version")
		pythonVersionOutput, pythonVersionErr := pythonVersionCommand.CombinedOutput()
		if pythonVersionErr != nil {
			output.Error("python not installed")
		} else {
			pythonVersionRegex := regexp.MustCompile(`^Python\s([\d|.]+)+`)
			pythonVersionMatches := pythonVersionRegex.FindAllSubmatch(pythonVersionOutput, -1)

			if len(pythonVersionMatches) < 1 {
				output.Error("python not installed")
			} else {
				installedPythonVer, installedPythonVersionParseErr := version.NewVersion(string(pythonVersionMatches[0][1]))
				if installedPythonVersionParseErr != nil {
					output.Error(versionFetchErr.Error())
				}

				minPythonVer, minPythonVersionParseErr := version.NewVersion(doctorPythonMinVersion)
				if minPythonVersionParseErr != nil {
					output.Error(installedVersionParseErr.Error())
				}

				if installedPythonVer.GreaterThanOrEqual(minPythonVer) {
					output.Ok(fmt.Sprintf("Version: %s >= %s", installedPythonVer, minPythonVer))
				} else {
					output.Warning(fmt.Sprintf("Version: %s <= %s", installedPythonVer, minPythonVer))
				}
			}
		}

		output.Section("Github")

		if os.Getenv("GITHUB_TOKEN") == "" {
			output.Error("GITHUB_TOKEN env var is not set")
		} else {
			output.Ok("GITHUB_TOKEN is set")
		}

		org, orgErr := git.GetOrganization()
		if orgErr != nil || org == nil {
			output.Error("Github organization inaccessible")
		} else {
			output.Ok("Github organization accessible")
		}

	},
}
