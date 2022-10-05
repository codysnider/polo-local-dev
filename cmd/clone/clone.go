package clone

import (
	"fmt"
	"github.com/poloniex/polo-local-dev/cmd/util"
	"github.com/poloniex/polo-local-dev/config"
	"github.com/poloniex/polo-local-dev/git"
	"github.com/poloniex/polo-local-dev/output"
	"github.com/spf13/cobra"
	"os"
)

var groupFlag string
var projectFlag string
var allFlag bool
var setRemotesFlag bool

var Command = &cobra.Command{
	Use:   "clone",
	Short: "Clone project",
	Long:  "Clones repo(s) to local environment. If a repo is already cloned, it is skipped. If the -r flag is set, origin will be set to fork and upstream will be set to Poloniex.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		output.Title("Clone")

		projectsToClone, projectsErr := util.ProjectsFromFlags(groupFlag, projectFlag, allFlag)
		if projectsErr != nil {
			output.Warning(projectsErr.Error())
			return
		}

		for _, project := range projectsToClone {

			output.Section(project.Name)

			// Check for folder
			if _, err := os.Stat(fmt.Sprintf("%s/%s", config.Config.WorkspaceRoot, project.Name)); err != nil {
				repo, orgRepoErr := git.GetOrganizationRepo(project.GetRepoName())
				if orgRepoErr != nil {
					output.Error(orgRepoErr.Error())
					continue
				}

				cloneErr := git.CloneRepo(fmt.Sprintf("%s/%s", config.Config.WorkspaceRoot, project.Name), repo)
				if cloneErr != nil {
					output.Error(cloneErr.Error())
					continue
				}
			} else {
				output.Ok("src folder exists")
			}

		}
	},
}

func init() {
	Command.PersistentFlags().BoolVarP(&setRemotesFlag, "remote", "r", true, "set remote")
	util.CommonProjectFlags(Command, &groupFlag, &projectFlag, &allFlag)
}
