package util

import (
	"errors"
	"github.com/poloniex/polo-local-dev/config"
)

func ProjectsFromFlags(groupFlag, projectFlag string, allFlag bool) (map[string]config.Project, error) {
	projects := map[string]config.Project{}

	if allFlag {
		for projectKey, projectConfig := range config.ProjectConfigs {
			projects[projectKey] = projectConfig
		}
	} else if groupFlag != "" {
		projects = config.GetProjectsByGroup(groupFlag)
	} else if projectFlag != "" {
		projects[projectFlag] = config.GetProjectByKey(projectFlag)
	}

	if len(projects) == 0 {
		return nil, errors.New("no projects based on parameters")
	}

	return projects, nil
}

func ProjectsFromFlagsWithDeps(groupFlag, projectFlag string, allFlag bool, depGroup string) (map[string]config.Project, error) {

	projects, projectsErr := ProjectsFromFlags(groupFlag, projectFlag, allFlag)
	if projectsErr != nil {
		return projects, projectsErr
	}

	// Generate dependency chain based on user input
	config.GenerateReverseGraph(projects, depGroup)
	fullProjectSet := config.GetProjectsByReverseGraph()

	return fullProjectSet, nil
}
