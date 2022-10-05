package config

type ProjectFile map[string]Project

type DependsOn struct {
	Compile []interface{} `json:"compile,omitempty"`
	Run     []interface{} `json:"run,omitempty"`
}

type Project struct {
	Repo           string    `json:"repo,omitempty"`
	Name           string    `json:"name,omitempty"`
	Groups         []string  `json:"groups,omitempty"`
	DefaultVersion string    `json:"default_version,omitempty"`
	BuildCmd       []string  `json:"build_cmd,omitempty"`
	DependsOn      DependsOn `json:"depends_on,omitempty"`
	FlywayPath     string    `json:"flyway_path,omitempty"`
}

func (p *Project) GetRepoName() string {
	if p.Repo != "" {
		return p.Repo
	}

	return p.Name
}

func GetProjectsByGroup(group string) []Project {
	matchingProjects := []Project{}
	for _, projectConfig := range ProjectConfigs {
		for _, projectGroup := range projectConfig.Groups {
			if group == projectGroup {
				matchingProjects = append(matchingProjects, projectConfig)
			}
		}
	}
	return matchingProjects
}

func GetProjectByName(name string) Project {
	for _, projectConfig := range ProjectConfigs {
		if projectConfig.Name == name {
			return projectConfig
		}
	}

	return Project{}
}
