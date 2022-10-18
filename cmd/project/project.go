package project

import (
	"github.com/poloniex/polo-local-dev/cmd/util"
	"github.com/poloniex/polo-local-dev/output"
	"github.com/spf13/cobra"
)

var groupFlag string
var projectFlag string
var allFlag bool

var Command = &cobra.Command{
	Use:   "project",
	Short: "Project and group details",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var projectDetails = &cobra.Command{
	Use:   "details",
	Short: "Display project details",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		output.Title("Project Details")

		projects, projectsErr := util.ProjectsFromFlags(groupFlag, projectFlag, allFlag)
		if projectsErr != nil {
			output.Warning(projectsErr.Error())
			return
		}

		for _, proj := range projects {
			output.Section(proj.Name)
			output.Plain(proj.Display())
		}

	},
}

func init() {
	util.CommonProjectFlags(Command, &groupFlag, &projectFlag, &allFlag)

	Command.AddCommand(projectDetails)
}
