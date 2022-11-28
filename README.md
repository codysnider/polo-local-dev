# Poloniex Local Development Toolkit

## Setup

To build the latest, run the following script with MakeMeAdmin active:

```bash
./build.sh
```

The last line will be the path of the executable. You should see `/usr/local/bin/gopld`.

Attempt running it with:
```bash
pld
```

## Usage Examples

**Common Flags**

*Note: Not all commands support json and verbose output but will not fail if argument is present.*

```
-h, --help      Help
-j, --json      JSON output
-v, --verbose   Verbose output
```

**Project-based Command Flags**

<a name="project-flags"></a>

```
-a, --all              All projects
-g, --group string     Project group
-h, --help             Help
-i, --ignore-deps      Ignore dependency chain
-p, --project string   Single project
```

### Config

**Reload**

Will refresh with project configurations with default/distributed project config.

```bash
pld config reload
```

### Doctor

Will run a sanity check on local environment state. This checks for system configuration, installed applications and authentication state.

```bash
pld doctor
```

### Build

Uses [project-based flags](#project-flags)

**All**
```bash
pld build --all
```

**Project Group**
```bash
pld build -g frontend
```

**Single Project (with dependencies)**
```bash
pld build -p frontend-login
```

**Single Project (without dependencies)**
```bash
pld build -i -p frontend-login
```

### Start

Uses [project-based flags](#project-flags)

**All**
```bash
pld start --all
```

**Project Group**
```bash
pld start -g frontend
```

**Single Project (with dependencies)**
```bash
pld start -p frontend-login
```

**Single Project (without dependencies)**
```bash
pld start -i -p frontend-login
```

### Clone

Uses [project-based flags](#project-flags)

**All**
```bash
pld clone --all
```

**Project Group**
```bash
pld clone -g support
```

### Fork

Uses [project-based flags](#project-flags)

**All**
```bash
pld fork --all
```

**Project Group**
```bash
pld fork -g support
```

### Dependency Tools

Uses [project-based flags](#project-flags)

**Graph**

Display the dependency tree

```bash
pld dep graph --all
```

**Execution Order**

Display the execution order based on dependencies

```bash
pld dep order --all
```

**Explain**

Combination of graph and execution order

```bash
pld dep explain --all
```

### Project Details

Uses [project-based flags](#project-flags)

Get details about project

```bash
pld project details --all
```

## Config Format

*Note: All project config files must be named in the *.project.json format and each can contain any number of projects. These are converted and copied during PLD install and config reload operations.* 

### Parameters

| Parameter Name     | Parameter Description                                                                                         | Required |
|--------------------|---------------------------------------------------------------------------------------------------------------|----------|
| PROJECT-NAME       | Distinct name of project, must be universally unique in config files                                          | YES      |
| REPO-NAME          | Repository name. Use for both git and filesystem operations                                                   | YES      |
| DOCKER-NAME        | Name for project, must match with service name in docker compose                                              | YES      |
| SERVICE-GROUP      | Optional service group membership for use with `--group` flags                                                | NO       |
| DEFAULT-GIT-BRANCH | Base branch for repo, should be main or master unless a larger epic is in progress. Use for git sync commands | YES      |
| BUILD-BASH-COMMAND | Build-phase bash command, any number can be defined, run in sequence and expect a 0 exit code                 | NO       |
| BUILD-EXEC-PATH    | Filesystem location in which the paired command should execute                                                | NO       |
| RUN-BASH-COMMAND   | Run-phase bash command, any number can be defined, run in sequence and expect a 0 exit code                   | NO       |
| RUN-EXEC-PATH      | Filesystem location in which the paired command should execute                                                | NO       |

### Format Template

```json
{
  "PROJECT-NAME": {
    "repo": "REPO-NAME",
    "name": "DOCKER-NAME",
    "groups": [
      "SERVICE-GROUP",
      "SERVICE-GROUP",
      "SERVICE-GROUP"
    ],
    "default_version": "DEFAULT-GIT-BRANCH",
    "build_cmd": [
      {
        "command": "BUILD-BASH-COMMAND",
        "path": "BUILD-EXEC-PATH"
      },
      {
        "command": "BUILD-BASH-COMMAND",
        "path": "BUILD-EXEC-PATH"
      },
      {
        "command": "BUILD-BASH-COMMAND",
        "path": "BUILD-EXEC-PATH"
      }
    ],
    "run_cmd": [
      {
        "command": "RUN-BASH-COMMAND",
        "path": "RUN-EXEC-PATH"
      },
      {
        "command": "RUN-BASH-COMMAND",
        "path": "RUN-EXEC-PATH"
      },
      {
        "command": "RUN-BASH-COMMAND",
        "path": "RUN-EXEC-PATH"
      }
    ]
  }
}
```
