package start

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
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

	ProjectLoop:
		for _, projectKey := range orderedProjects {
			project := config.GetProjectByKey(projectKey)
			output.Section(project.Name)

			if len(project.RunPrepare()) == 0 {
				output.Plain("No run commands defined")
			}

			for _, shellCmd := range project.RunPrepare() {

				output.Plain(fmt.Sprintf("Command: %s", shellCmd.String()))
				output.Plain(fmt.Sprintf("Path: %s", shellCmd.Dir))

				// Get stdout and stderr pipes
				stderr, _ := shellCmd.StderrPipe()
				stdout, _ := shellCmd.StdoutPipe()
				if err := shellCmd.Start(); err != nil {
					output.Error(err.Error())
					os.Exit(1)
				}

				// Create output writer channels
				outputWriter := make(chan string)
				closeSignal := make(chan bool, 1)
				finished := make(chan bool, 1)

				// Create output writer coroutine
				go output.FifoOutput("Run Command Output", 6, outputWriter, closeSignal, finished)

				// Add STDERR writer coroutine
				go func() {
					scanner := bufio.NewScanner(stderr)
					scanner.Split(bufio.ScanLines)
					for scanner.Scan() {
						m := scanner.Text()
						outputWriter <- m
					}
				}()

				// Add STDOUT writer coroutine
				go func() {
					scanner := bufio.NewScanner(stdout)
					scanner.Split(bufio.ScanLines)
					for scanner.Scan() {
						m := scanner.Text()
						outputWriter <- m
					}
				}()

				// Wait for command to exit
				cmdErr := shellCmd.Wait()

				// Send signal to coroutine to clear output
				closeSignal <- true

				// Block until writer coroutine finished
				<-finished

				if cmdErr != nil {
					if exiterr, ok := cmdErr.(*exec.ExitError); ok {
						output.Error(fmt.Sprintf("Exit Status: %d", exiterr.ExitCode()))
					}
				} else {
					output.Ok("Done")
				}
			}

			projectContainer := project.FindRunningContainer()
			if len(projectContainer.ID) > 0 {
				if !docker.ContainerHasHealthCheck(context.Background(), &projectContainer) {
					output.Warning("Container has no health check")
				} else {

					// Create health check writer channels
					outputWriter := make(chan string)
					closeSignal := make(chan bool, 1)
					finished := make(chan bool, 1)

					// Create output writer coroutine
					go output.FifoOutput("Health Check Output", 6, outputWriter, closeSignal, finished)

					stopHealthcheckUpdates := make(chan bool, 1)
					healthCheckLogs := make(chan *types.HealthcheckResult)
					go docker.ContainerHealthCheckStream(context.Background(), &projectContainer, stopHealthcheckUpdates, healthCheckLogs)

					// Begin status check coroutine
					healthyStatus := make(chan bool, 1)
					go func() {
						output.InputCancelFunc(func(healthy chan<- bool) {
							for {
								if docker.ContainerIsHealthy(&projectContainer) {
									healthy <- true
									return
								}
							}
						}, time.Second*30, healthyStatus)
					}()

					// Wait for healthy status or health check log output
					for {
						select {
						case logEntry := <-healthCheckLogs:

							// Write logs to terminal
							outputWriter <- logEntry.Output

						case healthy := <-healthyStatus:

							// Stop the health check output coroutine
							stopHealthcheckUpdates <- true

							// Send signal to coroutine to clear output
							closeSignal <- true

							// Block until writer cleans up
							<-finished

							// Display health status
							if healthy {
								output.Ok("Healthy")
							} else {
								output.Error("NOT Healthy")
							}

							// Move to next project
							continue ProjectLoop
						}
					}

				}
			}
		}
	},
}

func init() {
	util.CommonProjectFlags(Command, &groupFlag, &projectFlag, &allFlag)
	util.DependencyFlags(Command, &ignoreDepsFlag)
}
