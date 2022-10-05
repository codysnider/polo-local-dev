package util

import (
	"errors"
	"github.com/poloniex/polo-local-dev/config"
)

func ProjectsFromFlags(groupFlag, projectFlag string, allFlag bool) ([]config.Project, error) {
	projects := []config.Project{}

	if allFlag {
		for _, projectConfig := range config.ProjectConfigs {
			projects = append(projects, projectConfig)
		}
	} else if groupFlag != "" {
		projects = config.GetProjectsByGroup(groupFlag)
	} else if projectFlag != "" {
		projects = append(projects, config.GetProjectByName(projectFlag))
	}

	if len(projects) == 0 {
		return nil, errors.New("no projects based on parameters")
	}

	return projects, nil
}
