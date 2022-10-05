package fork

import (
	"fmt"
	"github.com/poloniex/polo-local-dev/cmd/util"
	"github.com/poloniex/polo-local-dev/git"
	"github.com/poloniex/polo-local-dev/output"
	"github.com/spf13/cobra"
)

var groupFlag string
var projectFlag string
var allFlag bool

var Command = &cobra.Command{
	Use:   "fork",
	Short: "Fork project",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		projectsToFork, projectsErr := util.ProjectsFromFlags(groupFlag, projectFlag, allFlag)
		if projectsErr != nil {
			output.Warning(projectsErr.Error())
			return
		}

		for _, project := range projectsToFork {
			repo, repoErr := git.GetOrganizationRepo(project.GetRepoName())
			if repoErr != nil {
				output.Error(repoErr.Error())
				continue
			}

			if !repo.GetAllowForking() {
				output.Warning(fmt.Sprintf("%s does not allow forking", repo.GetName()))
				continue
			}

			repoForkErr := git.ForkRepo(repo)
			if repoForkErr != nil {
				output.Warning(repoForkErr.Error())
				continue
			}

			output.Ok(repo.GetName())
		}
	},
}

func init() {
	util.CommonProjectFlags(Command, &groupFlag, &projectFlag, &allFlag)
}
