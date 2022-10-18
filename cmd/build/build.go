package build

import (
	"bytes"
	"fmt"
	"github.com/poloniex/polo-local-dev/cmd/util"
	"github.com/poloniex/polo-local-dev/config"
	"github.com/poloniex/polo-local-dev/output"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var groupFlag string
var projectFlag string
var allFlag bool
var ignoreDepsFlag bool

var Command = &cobra.Command{
	Use:   "build",
	Short: "Build project",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		output.Title("Build")

		var projectsToBuild map[string]config.Project
		var projectsErr error
		orderedProjects := []string{}

		if ignoreDepsFlag {
			projectsToBuild, projectsErr = util.ProjectsFromFlags(groupFlag, projectFlag, allFlag)
			for projectKey := range projectsToBuild {
				orderedProjects = append(orderedProjects, projectKey)
			}
		} else {
			projectsToBuild, projectsErr = util.ProjectsFromFlagsWithDeps(groupFlag, projectFlag, allFlag, "build")
			orderedProjects = config.GenerateOrderedSet(projectsToBuild, "build")
		}

		if projectsErr != nil {
			output.Warning(projectsErr.Error())
			return
		}

		for _, projectKey := range orderedProjects {
			project := config.GetProjectByKey(projectKey)
			output.Section(project.Name)

			for _, shellCmd := range project.BuildPrepare() {

				output.Plain(fmt.Sprintf("Command: %s", shellCmd.String()))
				output.Plain(fmt.Sprintf("Path: %s", shellCmd.Dir))

				s := output.Spin("Building", "Done")

				var stdoutBuffer, stderrBuffer bytes.Buffer
				shellCmd.Stdout = &stdoutBuffer
				shellCmd.Stderr = &stderrBuffer
				if err := shellCmd.Start(); err != nil {
					output.Error(err.Error())
					os.Exit(1)
				}
				if err := shellCmd.Wait(); err != nil {
					if exiterr, ok := err.(*exec.ExitError); ok {
						output.Error(fmt.Sprintf("Exit Status: %d", exiterr.ExitCode()))
					}
				}

				s.Stop()

				if verbose, verboseFlagErr := cmd.Flags().GetBool("verbose"); verbose && verboseFlagErr == nil {
					fmt.Println(shellCmd.Stdout)
					fmt.Println(shellCmd.Stderr)
				}
			}
		}
	},
}

func init() {
	util.CommonProjectFlags(Command, &groupFlag, &projectFlag, &allFlag)
	util.DependencyFlags(Command, &ignoreDepsFlag)
}
