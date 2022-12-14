package util

import "github.com/spf13/cobra"

func CommonProjectFlags(cmd *cobra.Command, groupFlag, projectFlag *string, allFlag *bool) {
	cmd.PersistentFlags().StringVarP(groupFlag, "group", "g", "", "project group")
	cmd.PersistentFlags().StringVarP(projectFlag, "project", "p", "", "project")
	cmd.PersistentFlags().BoolVarP(allFlag, "all", "a", false, "all projects")
}

func DependencyFlags(cmd *cobra.Command, ignoreDeps *bool) {
	cmd.PersistentFlags().BoolVarP(ignoreDeps, "ignore-deps", "i", false, "ignore dependency chain")
}
