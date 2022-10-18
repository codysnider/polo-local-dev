package start

import (
	"bytes"
	"context"
	"fmt"
	"github.com/poloniex/polo-local-dev/cmd/util"
	"github.com/poloniex/polo-local-dev/config"
	"github.com/poloniex/polo-local-dev/docker"
	"github.com/poloniex/polo-local-dev/output"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"time"
)

var groupFlag string
var projectFlag string
var allFlag bool
var ignoreDepsFlag bool

var Command = &cobra.Command{
	Use:   "start",
	Short: "Start project",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		output.Title("Start")

		var projectsToRun map[string]config.Project
		var projectsErr error
		orderedProjects := []string{}

		if ignoreDepsFlag {
			projectsToRun, projectsErr = util.ProjectsFromFlags(groupFlag, projectFlag, allFlag)
			for projectKey := range projectsToRun {
				orderedProjects = append(orderedProjects, projectKey)
			}
		} else {
			projectsToRun, projectsErr = util.ProjectsFromFlagsWithDeps(groupFlag, projectFlag, allFlag, "run")
			orderedProjects = config.GenerateOrderedSet(projectsToRun, "run")
		}

		if projectsErr != nil {
			output.Warning(projectsErr.Error())
			return
		}

		for _, projectKey := range orderedProjects {
			project := config.GetProjectByKey(projectKey)
			output.Section(project.Name)

			for _, shellCmd := range project.RunPrepare() {

				output.Plain(fmt.Sprintf("Command: %s", shellCmd.String()))
				output.Plain(fmt.Sprintf("Path: %s", shellCmd.Dir))

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
				} else {
					output.Ok("Done without errors")
				}

				if verbose, verboseFlagErr := cmd.Flags().GetBool("verbose"); verbose && verboseFlagErr == nil {
					fmt.Println(shellCmd.Stdout)
					fmt.Println(shellCmd.Stderr)
				}
			}

			projectContainer := project.FindRunningContainer()
			if len(projectContainer.ID) > 0 {
				if !docker.ContainerHasHealthCheck(context.Background(), &projectContainer) {
					output.Warning("Container has no health check")
				} else {
					healthCheckSpinner := output.Spin("Waiting for health check", output.OkString("Health check done"))
					healthCheckSpinner.Start()

					healthyStatus := make(chan bool, 1)
					output.InputCancelFunc(func(healthy chan<- bool) {
						for {
							if docker.ContainerIsHealthy(&projectContainer) {
								healthy <- true
								break
							}
							time.Sleep(time.Second)
						}
					}, time.Second*30, healthyStatus)

					if <-healthyStatus {
						healthCheckSpinner.FinalMSG = output.OkString("Container healthy")
					} else {
						healthCheckSpinner.FinalMSG = output.ErrorString("Container NOT healthy")
					}

					healthCheckSpinner.Stop()
				}
			}
		}
	},
}

func init() {
	util.CommonProjectFlags(Command, &groupFlag, &projectFlag, &allFlag)
	util.DependencyFlags(Command, &ignoreDepsFlag)
}
