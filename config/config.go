package config

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/poloniex/polo-local-dev/output"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// This bakes the distributed config files into the binary
	//go:embed dist
	distEmbed embed.FS

	// Default path for local configs
	configPath = "~/.pld/"

	// For validation of paths
	pathRegex, _ = regexp.Compile("^(/[^/ ]*)+/?$")

	// ProjectConfigs post-installation and resolution
	ProjectConfigs = map[string]Project{}

	// Config is the PLD application settings
	Config CommonConfig
)

func absolutePath(path string) string {

	// Get user home path
	usr, _ := user.Current()
	dir := usr.HomeDir

	// Expand
	if path == "~" {
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}

	return path
}

func init() {

	// Ensure ~/.pld/ exists
	configAbsolutePath := absolutePath(configPath)
	if _, err := os.Stat(configAbsolutePath); errors.Is(err, os.ErrNotExist) {
		mkdirErr := os.Mkdir(configAbsolutePath, os.ModePerm)
		if mkdirErr != nil {
			output.Error(mkdirErr.Error())
			os.Exit(1)
		}
	}

	// Load common config
	var commonConfigErr error
	Config, commonConfigErr = loadCommonConfig()
	if commonConfigErr != nil {
		output.Error(commonConfigErr.Error())
		os.Exit(1)
	}

	// Load dist and local project configs
	distConfigs, installedConfigs, configLoadErr := loadProjectConfigs()
	if configLoadErr != nil {
		output.Error(configLoadErr.Error())
		os.Exit(1)
	}

	// Search for uninstalled configs and install
	for projectName := range distConfigs {
		if _, installed := installedConfigs[projectName]; !installed {
			installErr := installProjectConfig(projectName, distConfigs[projectName])
			if installErr != nil {
				output.Error(installErr.Error())
				os.Exit(1)
			}
			output.Ok(fmt.Sprintf("Installing config: %s", projectName))
		}
	}

	// TODO: Logic to remove stale/removed config files (e.g. no longer in dist, still in local)

	// Reload local configs
	_, installedConfigs, configLoadErr = loadProjectConfigs()
	if configLoadErr != nil {
		output.Error(configLoadErr.Error())
		os.Exit(1)
	}

	for projectName, project := range installedConfigs {
		ProjectConfigs[projectName] = project
	}
}

func generateCommonConfig() (CommonConfig, error) {
	commonConfig := CommonConfig{}

	workspaceRootPrompt := promptui.Prompt{
		Label: "Workspace root",
		Validate: func(input string) error {
			if !pathRegex.MatchString(absolutePath(input)) {
				return errors.New("invalid path")
			}
			return nil
		},
	}

	workspaceRoot, promptErr := workspaceRootPrompt.Run()
	if promptErr != nil {
		return commonConfig, promptErr
	}
	commonConfig.WorkspaceRoot = absolutePath(workspaceRoot)

	return commonConfig, nil
}

func loadCommonConfig() (CommonConfig, error) {
	var commonConfig CommonConfig
	if _, err := os.Stat(absolutePath(configPath) + "/config.json"); err == nil {
		configFile, fileReadErr := ioutil.ReadFile(absolutePath(configPath) + "/config.json")
		if fileReadErr != nil {
			return commonConfig, fileReadErr
		}

		if parseErr := json.Unmarshal(configFile, &commonConfig); parseErr != nil {
			return commonConfig, parseErr
		}
	} else {
		var commonConfigErr error
		commonConfig, commonConfigErr = generateCommonConfig()
		if commonConfigErr != nil {
			return commonConfig, commonConfigErr
		}

		if configInstallErr := installCommonConfig(commonConfig); configInstallErr != nil {
			return commonConfig, commonConfigErr
		}
	}

	return commonConfig, nil
}

func loadProjectConfigs() (map[string]Project, map[string]Project, error) {

	distConfigs := map[string]Project{}

	// Load embedded configs
	distFiles, distReadErr := distEmbed.ReadDir("dist")
	if distReadErr != nil {
		return nil, nil, distReadErr
	}

	for _, distFilename := range distFiles {

		// Only parse *.project.json files
		if strings.Contains(distFilename.Name(), ".project.json") {

			// Read file contents
			distFile, fileReadErr := distEmbed.ReadFile("dist/" + distFilename.Name())
			if fileReadErr != nil {
				return nil, nil, fileReadErr
			}

			// Unmarshal into ProjectFile
			var distConfig ProjectFile
			if parseErr := json.Unmarshal(distFile, &distConfig); parseErr != nil {
				return nil, nil, parseErr
			}

			// Append to returned set
			for projectName, project := range distConfig {
				distConfigs[projectName] = project
			}
		}
	}

	installedConfigs := map[string]Project{}

	installedFiles, installFolderReadErr := ioutil.ReadDir(absolutePath(configPath))
	if installFolderReadErr != nil {
		return nil, nil, installFolderReadErr
	}

	for _, installedFile := range installedFiles {
		if !installedFile.IsDir() && strings.Contains(installedFile.Name(), ".project.json") {

			jsonFile, fileOpenErr := os.Open(absolutePath(configPath) + "/" + installedFile.Name())
			if fileOpenErr != nil {
				return nil, nil, fileOpenErr
			}
			defer func(jsonFile *os.File) {
				_ = jsonFile.Close()
			}(jsonFile)

			fileBytes, fileReadErr := ioutil.ReadAll(jsonFile)
			if fileReadErr != nil {
				return nil, nil, fileReadErr
			}

			// Unmarshal into ProjectFile
			var installedConfig ProjectFile
			if parseErr := json.Unmarshal(fileBytes, &installedConfig); parseErr != nil {
				return nil, nil, parseErr
			}

			// Append to returned set
			for projectName, project := range installedConfig {
				installedConfigs[projectName] = project
			}
		}
	}

	return distConfigs, installedConfigs, nil
}

func installProjectConfig(name string, project Project) error {

	projectFile := map[string]Project{
		name: project,
	}

	projectJson, jsonErr := json.MarshalIndent(&projectFile, "", "    ")
	if jsonErr != nil {
		return jsonErr
	}

	fileWriteErr := ioutil.WriteFile(fmt.Sprintf("%s/%s.project.json", absolutePath(configPath), name), projectJson, os.ModePerm)
	if fileWriteErr != nil {
		return fileWriteErr
	}

	return nil
}

func installCommonConfig(commonConfig CommonConfig) error {

	configJson, jsonErr := json.MarshalIndent(&commonConfig, "", "    ")
	if jsonErr != nil {
		return jsonErr
	}

	fileWriteErr := ioutil.WriteFile(fmt.Sprintf("%s/config.json", absolutePath(configPath)), configJson, os.ModePerm)
	if fileWriteErr != nil {
		return fileWriteErr
	}

	return nil
}
