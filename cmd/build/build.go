package build

import (
	"bufio"
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

				// Get stdout and stderr pipes
				stderr, _ := shellCmd.StderrPipe()
				stdout, _ := shellCmd.StdoutPipe()
				if err := shellCmd.Start(); err != nil {
					output.Error(err.Error())
					os.Exit(1)
				}

				// Create output writer channels
				outputWriter := make(chan string)
				closeSignal := make(chan bool)
				finished := make(chan bool)

				// Create output writer coroutine
				go output.FifoOutput(6, outputWriter, closeSignal, finished)

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
		}
	},
}

func init() {
	util.CommonProjectFlags(Command, &groupFlag, &projectFlag, &allFlag)
	util.DependencyFlags(Command, &ignoreDepsFlag)
}
