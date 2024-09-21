# Wails v3 Build System

## Overview

The Wails v3 build system is a flexible and powerful tool designed to streamline the build process for your Wails applications. 
It leverages [Task](https://taskfile.dev), a task runner that allows you to define and run tasks easily. 
While the v3 build system is the default, Wails encourages a "bring your own tooling" approach, allowing developers to customize their build process as needed.

Learn more about how to use Task in the [official documentation](https://taskfile.dev/usage/).

## Task: The Heart of the Build System

[Task](https://taskfile.dev) is a modern alternative to Make, written in Go. It uses a YAML file to define tasks and their dependencies. In the Wails v3 build system, [Task](https://taskfile.dev) plays a central role in orchestrating the build process.

The main `Taskfile.yml` is located in the project root, while platform-specific tasks are defined in `build/Taskfile.<platform>.yml` files.

```
Project Root
│
├── Taskfile.yml
│
└── build
    ├── Taskfile.windows.yml
    ├── Taskfile.darwin.yml
    ├── Taskfile.linux.yml
    └── Taskfile.common.yml
```

The `Taskfile.common.yml` file contains common tasks that are shared across platforms.

## Taskfile.yml

The `Taskfile.yml` file is the main entry point for the build system. It defines the tasks and their dependencies. Here's the default ``Taskfile.yml`` file:

```yaml
version: '3'

includes:
  common: ./build/Taskfile.common.yml
  windows: ./build/Taskfile.windows.yml
  darwin: ./build/Taskfile.darwin.yml
  linux: ./build/Taskfile.linux.yml

vars:
  APP_NAME: "{{.ProjectName}}"
  BIN_DIR: "bin"
  VITE_PORT: '{{.WAILS_VITE_PORT | default 9245}}'

tasks:
  build:
    summary: Builds the application
    cmds:
      - task: "{{OS}}:build"

  package:
    summary: Packages a production build of the application
    cmds:
      - task: "{{OS}}:package"

  run:
    summary: Runs the application
    cmds:
      - task: "{{OS}}:run"

  dev:
    summary: Runs the application in development mode
    cmds:
      - wails3 dev -config ./build/devmode.config.yaml -port {{.VITE_PORT}}

  dev:reload:
    summary: Reloads the application
    cmds:
      - task: run
```

## Platform-Specific Taskfiles

Each platform has its own Taskfile, located in the `build` directory. These files define the core tasks for that
platform. Each taskfile includes common tasks from the `Taskfile.common.yml` file.

### Windows

Location: `build/Taskfile.windows.yml`

The Windows-specific Taskfile includes tasks for building, packaging, and running the application on Windows. Key features include:

- Building with optional production flags
- Generating Windows ``.syso`` file
- Creating an NSIS installer for packaging

### Linux

Location: `build/Taskfile.linux.yml`

The Linux-specific Taskfile includes tasks for building, packaging, and running the application on Linux. Key features include:

- Building with optional production flags
- Creating an AppImage for packaging
- Generating ``.desktop`` file for Linux applications

### macOS

Location: `build/Taskfile.darwin.yml`

The macOS-specific Taskfile includes tasks for building, packaging, and running the application on macOS. Key features include:

- Building with optional production flags
- Creating an ``.app`` bundle for packaging
- Setting macOS-specific build flags and environment variables

## Wails3 Commands and Task Execution

The `wails3 task` command is an embedded version of taskfile.dev, which executes the tasks defined in your Taskfile.yml.

The `wails3 build` and `wails3 package` commands are aliases for `wails3 task build` and `wails3 task package` respectively. When you run these commands, Wails internally translates them to the appropriate task execution:

- `wails3 build` → `wails3 task build`
- `wails3 package` → `wails3 task package`

## Common Build Process

Across all platforms, the build process typically includes the following steps:

1. Tidying Go modules
2. Building the frontend
3. Generating icons
4. Compiling the Go code with platform-specific flags
5. Packaging the application (platform-specific)

## Customising the Build Process

While the v3 build system provides a solid default configuration, you can easily customise it to fit your project's needs. By modifying the `Taskfile.yml` and platform-specific Taskfiles, you can:

- Add new tasks
- Modify existing tasks
- Change the order of task execution
- Integrate with other tools and scripts

This flexibility allows you to tailor the build process to your specific requirements while still benefiting from the structure provided by the Wails v3 build system.

## Development Mode

The Wails v3 build system includes a powerful development mode that enhances the developer experience by providing live reloading and hot module replacement. 
This mode is activated using the `wails3 dev` command.

### How It Works

When you run `wails3 dev`, the following process occurs:

1. The command checks for an available port, defaulting to 9245 if not specified.
2. It sets up the environment variables for the frontend dev server (Vite).
3. It starts the file watcher using the `refresh` library.

The [refresh](https://github.com/atterpac/refresh) library is responsible for monitoring file changes and triggering rebuilds. 
It uses a configuration file, typically located at `./build/devmode.config.yaml`, to determine which files to watch and what 
actions to take when changes are detected.

### Configuration

The development mode can be configured using the `devmode.config.yaml` file. Here's an example of its structure:

```yaml
config:
  root_path: .
  log_level: warn
  debounce: 1000
  ignore:
    dir:
      - .git
      - node_modules
      - frontend
      - bin
    file:
      - .DS_Store
      - .gitignore
      - .gitkeep
    watched_extension:
      - "*.go"
    git_ignore: true
  executes:
    - cmd: wails3 task common:install:frontend:deps
      type: once
    - cmd: wails3 task common:dev:frontend
      type: background
    - cmd: go mod tidy
      type: blocking
    - cmd: wails3 task build
      type: blocking
    - cmd: wails3 task run
      type: primary
```

This configuration file allows you to:

- Set the root path for file watching
- Configure logging level
- Set a debounce time for file change events
- Ignore specific directories, files, or file extensions
- Define commands to execute on file changes

### Customising Development Mode

You can customise the development mode experience by modifying the `devmode.config.yaml` file. 

Some ways to customise include:

1. Changing the watched directories or files
2. Adjusting the debounce time to control how quickly the system responds to changes
3. Adding or modifying the execute commands to fit your project's needs

You can also specify a custom configuration file and port:

```shell
wails3 dev -config ./path/to/custom/config.yaml -port 8080
```
