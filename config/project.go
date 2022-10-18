package config

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/poloniex/polo-local-dev/docker"
	"github.com/poloniex/polo-local-dev/output"
	"github.com/tufin/asciitree"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/tabwriter"
)

var reverseGraph = map[string][]string{}

type ProjectFile map[string]Project

type DependsOn struct {
	Compile []string `json:"compile,omitempty"`
	Run     []string `json:"run,omitempty"`
}

type ReverseDependsOn struct {
	Compile []string `json:"compile,omitempty"`
	Run     []string `json:"run,omitempty"`
}

type ShellCommand struct {
	Command string `json:"command,omitempty"`
	Path    string `json:"path,omitempty"`
}

type Project struct {
	Repo             string           `json:"repo,omitempty"`
	Name             string           `json:"name,omitempty"`
	Groups           []string         `json:"groups,omitempty"`
	DefaultVersion   string           `json:"default_version,omitempty"`
	BuildCmd         []ShellCommand   `json:"build_cmd,omitempty"`
	RunCmd           []ShellCommand   `json:"run_cmd,omitempty"`
	DependsOn        DependsOn        `json:"depends_on,omitempty"`
	ReverseDependsOn ReverseDependsOn `json:"-"`
}

func (p *Project) stringReplacements() map[string]string {
	return map[string]string{
		"#NAME#":           p.Name,
		"#PROJECT_ROOT#":   p.RootPath(),
		"#REPO#":           p.Repo,
		"#WORKSPACE_ROOT#": Config.WorkspaceRoot,
	}
}

func (p *Project) RootPath() string {
	return fmt.Sprintf("%s/%s", Config.WorkspaceRoot, p.Repo)
}

func (p *Project) GetRepoName() string {
	if p.Repo != "" {
		return p.Repo
	}

	return p.Name
}

func (p *Project) BuildPrepare() []*exec.Cmd {
	buildCmds := make([]*exec.Cmd, len(p.BuildCmd))
	for cmdIdx, cmd := range p.BuildCmd {
		for oldString, newString := range p.stringReplacements() {
			cmd.Command = strings.ReplaceAll(cmd.Command, oldString, newString)
		}
		cmdSplit := strings.Split(cmd.Command, " ")
		buildCmds[cmdIdx] = exec.Command(cmdSplit[0], cmdSplit[1:]...)

		for oldString, newString := range p.stringReplacements() {
			cmd.Path = strings.ReplaceAll(cmd.Path, oldString, newString)
		}
		buildCmds[cmdIdx].Dir = cmd.Path
	}

	return buildCmds
}

func (p *Project) RunPrepare() []*exec.Cmd {
	runCmds := make([]*exec.Cmd, len(p.RunCmd))
	for cmdIdx, cmd := range p.RunCmd {
		for oldString, newString := range p.stringReplacements() {
			cmd.Command = strings.ReplaceAll(cmd.Command, oldString, newString)
		}
		cmdSplit := strings.Split(cmd.Command, " ")
		runCmds[cmdIdx] = exec.Command(cmdSplit[0], cmdSplit[1:]...)

		for oldString, newString := range p.stringReplacements() {
			cmd.Path = strings.ReplaceAll(cmd.Path, oldString, newString)
		}
		runCmds[cmdIdx].Dir = cmd.Path
	}

	return runCmds
}

func GetProjectsByGroup(group string) map[string]Project {
	matchingProjects := map[string]Project{}
	for projectName, projectConfig := range ProjectConfigs {
		for _, projectGroup := range projectConfig.Groups {
			if group == projectGroup {
				matchingProjects[projectName] = projectConfig
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

func GetProjectByKey(key string) Project {
	if project, exists := ProjectConfigs[key]; exists {
		return project
	}

	return Project{}
}

func (p *Project) Dependencies(group string) []string {
	if group == "run" {
		return p.DependsOn.Run
	} else if group == "compile" {
		return p.DependsOn.Compile
	}

	return append(p.DependsOn.Run, p.DependsOn.Compile...)
}

func GenerateTree(projects map[string]Project, dependencyGroup string) {
	for projectKey, project := range projects {
		if len(project.Dependencies(dependencyGroup)) == 0 {
			tree := asciitree.Tree{}
			tree.Add(projectKey)
			if _, hasDependencies := reverseGraph[projectKey]; hasDependencies && len(reverseGraph[projectKey]) > 0 {
				for _, child := range reverseGraph[projectKey] {
					appendTree(child, projectKey, &tree)
				}
			}
			tree.Fprint(os.Stdout, true, "       ")
			fmt.Println("")
		}
	}
}

func appendTree(projectKey, parentPath string, tree *asciitree.Tree) *asciitree.Tree {
	if _, hasDependencies := reverseGraph[projectKey]; hasDependencies && len(reverseGraph[projectKey]) > 0 {
		for _, child := range reverseGraph[projectKey] {
			appendTree(child, fmt.Sprintf("%s/%s", parentPath, projectKey), tree)
		}
	} else {
		tree.Add(fmt.Sprintf("%s/%s", parentPath, projectKey))
	}

	return tree
}

func GenerateOrderedSet(projects map[string]Project, dependencyGroup string) []string {

	infiniteLoopFuse := 0

	orderedProjects := []string{}
	for len(projects) > 0 && infiniteLoopFuse < 200 {
		infiniteLoopFuse++

	projectLoop:
		for projectKey, project := range projects {

			if len(project.Dependencies(dependencyGroup)) == 0 {
				orderedProjects = append(orderedProjects, projectKey)
				delete(projects, projectKey)
				break
			} else {
				for _, dependsOnProject := range project.Dependencies(dependencyGroup) {
					dependencySatisfied := false
					for _, orderedProject := range orderedProjects {
						if dependsOnProject == orderedProject {
							dependencySatisfied = true
							break
						}
					}

					if !dependencySatisfied {
						continue projectLoop
					}
				}

				orderedProjects = append(orderedProjects, projectKey)
				delete(projects, projectKey)
				break
			}
		}
	}

	if len(projects) > 0 {
		output.Plain("Unresolved projects:")
		for projectKey := range projects {
			output.Error(projectKey)
		}
		output.Plain("")
	}

	return orderedProjects
}

// GenerateReverseGraph creates the list of "things that depend on me" for each project
func GenerateReverseGraph(projectsToGraph map[string]Project, dependencyGroup string) map[string][]string {

	additionalPassRequired := false

	for projectKey, project := range projectsToGraph {
		if _, exists := reverseGraph[projectKey]; !exists {
			reverseGraph[projectKey] = []string{}
			additionalPassRequired = true
		}
		for _, runDependency := range project.Dependencies(dependencyGroup) {
			found := false
			for _, child := range reverseGraph[runDependency] {
				if child == projectKey {
					found = true
					break
				}
			}
			if !found {
				reverseGraph[runDependency] = append(reverseGraph[runDependency], projectKey)
			}
		}
	}

	if additionalPassRequired {
		GenerateReverseGraph(GetProjectsByReverseGraph(), dependencyGroup)
	}

	return reverseGraph
}

func ProjectMapFromKeySlice(projects []string) map[string]Project {
	fullSet := map[string]Project{}
	for _, projectKey := range projects {
		if _, exists := fullSet[projectKey]; !exists {
			fullSet[projectKey] = GetProjectByKey(projectKey)
		}
	}

	return fullSet
}

// GetProjectsByReverseGraph creates a slice of Project structs instead of the reverse graph strings
func GetProjectsByReverseGraph() map[string]Project {
	fullSet := map[string]Project{}
	for projectKey, children := range reverseGraph {
		if _, exists := fullSet[projectKey]; !exists {
			fullSet[projectKey] = GetProjectByKey(projectKey)
		}
		for _, child := range children {
			if _, exists := fullSet[child]; !exists {
				fullSet[child] = GetProjectByKey(child)
			}
		}
	}

	return fullSet
}

func (p *Project) ContainerNameMatchers() (matchers []*regexp.Regexp) {
	if len(p.RootPath()) == 0 || len(p.Name) == 0 {
		matchers = append(matchers, regexp.MustCompile(fmt.Sprintf("\\/%s(_{1})%s(_{1})\\d+", p.Repo, p.Name)))
	}
	matchers = append(matchers, regexp.MustCompile(fmt.Sprintf("\\/([a-zA-Z\\-]+)?(_{1})%s(_{1})\\d+", p.Name)))
	return
}

func (p *Project) FindRunningContainer() types.Container {
	ctx := context.Background()
	containers, containerErr := docker.Containers(ctx)
	if containerErr != nil {
		output.Error(containerErr.Error())
		return types.Container{}
	}

	matchers := p.ContainerNameMatchers()
	matchingContainers := []types.Container{}
	for _, container := range containers {
		for _, containerName := range container.Names {
			for _, matcher := range matchers {
				if matcher.Match([]byte(containerName)) {
					matchingContainers = append(matchingContainers, container)
				}
			}
		}
	}

	if len(matchingContainers) > 1 {
		output.Warning("Found multiple container matches")
		return types.Container{}
	}

	if len(matchingContainers) == 0 {
		output.Warning("Could not match container")
		return types.Container{}
	}

	output.Ok(fmt.Sprintf("Found container %s", matchingContainers[0].ID[:10]))
	return matchingContainers[0]
}

func (p *Project) Display() string {
	out := strings.Builder{}
	w := tabwriter.NewWriter(&out, 10, 0, 3, ' ', 0)

	if len(p.Name) > 0 {
		_, _ = fmt.Fprintf(w, "System name\t%s\n", p.Name)
	}

	if len(p.DefaultVersion) > 0 {
		_, _ = fmt.Fprintf(w, "Default branch\t%s\n", p.DefaultVersion)
	}

	if len(p.Repo) > 0 {
		_, _ = fmt.Fprintf(w, "Repo\t%s\n", p.Repo)
	}

	for groupIndex, group := range p.Groups {
		if groupIndex == 0 {
			_, _ = fmt.Fprintf(w, "Groups\t%s\n", group)
		} else {
			_, _ = fmt.Fprintf(w, "\t%s\n", group)
		}
	}

	for depIndex, dep := range p.DependsOn.Compile {
		if depIndex == 0 {
			_, _ = fmt.Fprintf(w, "Build Dependencies\t%s\n", dep)
		} else {
			_, _ = fmt.Fprintf(w, "\t%s\n", dep)
		}
	}

	for depIndex, dep := range p.DependsOn.Run {
		if depIndex == 0 {
			_, _ = fmt.Fprintf(w, "Run Dependencies\t%s\n", dep)
		} else {
			_, _ = fmt.Fprintf(w, "\t%s\n", dep)
		}
	}

	for buildCmdIndex, buildCmd := range p.BuildPrepare() {
		if buildCmdIndex == 0 {
			_, _ = fmt.Fprintf(w, "Build Commands\tcd %s && %s\n", buildCmd.Dir, buildCmd.String())
		} else {
			_, _ = fmt.Fprintf(w, "\tcd %s && %s\n", buildCmd.Dir, buildCmd.String())
		}
	}

	for runCmdIndex, runCmd := range p.RunPrepare() {
		if runCmdIndex == 0 {
			_, _ = fmt.Fprintf(w, "Run Commands\tcd %s && %s\n", runCmd.Dir, runCmd.String())
		} else {
			_, _ = fmt.Fprintf(w, "\tcd %s && %s\n", runCmd.Dir, runCmd.String())
		}
	}

	_ = w.Flush()

	return out.String()
}
