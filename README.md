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
