package dependency

import (
	"github.com/poloniex/polo-local-dev/cmd/util"
	"github.com/poloniex/polo-local-dev/config"
	"github.com/poloniex/polo-local-dev/output"
	"github.com/spf13/cobra"
)

var groupFlag string
var projectFlag string
var allFlag bool

var Command = &cobra.Command{
	Use:   "dep",
	Short: "Project dependency tools",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var explain = &cobra.Command{
	Use:   "explain",
	Short: "Display dependency graph and execute order",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		output.Title("Dependencies")

		dependencyGroup := cmd.Flag("runMode").Value.String()

		projectsToGraph, projectsErr := util.ProjectsFromFlagsWithDeps(groupFlag, projectFlag, allFlag, dependencyGroup)
		if projectsErr != nil {
			output.Warning(projectsErr.Error())
			return
		}

		output.Section("Graph")

		if len(projectsToGraph) <= 1 {
			for projectKey := range projectsToGraph {
				output.Plain(projectKey)
			}
		} else {
			config.GenerateTree(projectsToGraph, dependencyGroup)
		}

		output.Section("Execution Order")
		if len(projectsToGraph) <= 1 {
			for projectKey := range projectsToGraph {
				output.Plain(projectKey)
			}
		} else {
			for _, project := range config.GenerateOrderedSet(projectsToGraph, dependencyGroup) {
				output.Plain(project)
			}
		}
	},
}

var graph = &cobra.Command{
	Use:   "graph",
	Short: "Display dependency graph",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		output.Title("Dependencies")

		dependencyGroup := cmd.Flag("runMode").Value.String()

		projectsToGraph, projectsErr := util.ProjectsFromFlagsWithDeps(groupFlag, projectFlag, allFlag, dependencyGroup)
		if projectsErr != nil {
			output.Warning(projectsErr.Error())
			return
		}

		output.Section("Dependency Graph")
		if len(projectsToGraph) <= 1 {
			for projectKey := range projectsToGraph {
				output.Plain(projectKey)
			}
		} else {
			config.GenerateTree(projectsToGraph, dependencyGroup)
		}
	},
}

var order = &cobra.Command{
	Use:   "order",
	Short: "Display dependency-based execute order",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		output.Title("Dependencies")

		dependencyGroup := cmd.Flag("runMode").Value.String()

		projectsToGraph, projectsErr := util.ProjectsFromFlagsWithDeps(groupFlag, projectFlag, allFlag, dependencyGroup)
		if projectsErr != nil {
			output.Warning(projectsErr.Error())
			return
		}

		output.Section("Execution Order")
		if len(projectsToGraph) <= 1 {
			for projectKey := range projectsToGraph {
				output.Plain(projectKey)
			}
		} else {
			for _, project := range config.GenerateOrderedSet(projectsToGraph, dependencyGroup) {
				output.Plain(project)
			}
		}
	},
}

func init() {
	util.CommonProjectFlags(Command, &groupFlag, &projectFlag, &allFlag)

	explain.PersistentFlags().StringP("runMode", "r", "build", "one of: build, run")
	Command.AddCommand(explain)

	graph.PersistentFlags().StringP("runMode", "r", "build", "one of: build, run")
	Command.AddCommand(graph)

	order.PersistentFlags().StringP("runMode", "r", "build", "one of: build, run")
	Command.AddCommand(order)
}
