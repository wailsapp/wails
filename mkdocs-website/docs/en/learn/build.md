# Wails v3 Build System

## Overview

The Wails v3 build system is a flexible and powerful tool designed to streamline the build process for your Wails applications. It leverages Taskfile, a task runner that allows you to define and run tasks easily. While the v3 build system is the default, Wails encourages a "bring your own tooling" approach, allowing developers to customize their build process as needed.

## Taskfile: The Heart of the Build System

Taskfile is a modern alternative to Make, written in Go. It uses a YAML file to define tasks and their dependencies. In the Wails v3 build system, Taskfile plays a central role in orchestrating the build process.

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

### Windows (Taskfile.windows.yml)

The Windows-specific Taskfile includes tasks for building, packaging, and running the application on Windows. Key features include:

- Building with optional production flags
- Generating Windows ``.syso`` file
- Creating an NSIS installer for packaging

### Linux (Taskfile.linux.yml)

The Linux-specific Taskfile includes tasks for building, packaging, and running the application on Linux. Key features include:

- Building with optional production flags
- Creating an AppImage for packaging
- Generating ``.desktop`` file for Linux applications

### macOS (Taskfile.darwin.yml)

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

## Customizing the Build Process

While the v3 build system provides a solid default configuration, you can easily customize it to fit your project's needs. By modifying the `Taskfile.yml` and platform-specific Taskfiles, you can:

- Add new tasks
- Modify existing tasks
- Change the order of task execution
- Integrate with other tools and scripts

This flexibility allows you to tailor the build process to your specific requirements while still benefiting from the structure provided by the Wails v3 build system.

## Development Mode

The build system includes a `dev` task for running the application in development mode. This task uses the `wails3 dev` command with a configuration file and a specified Vite port.

## Conclusion

The Wails v3 build system, powered by Taskfile, offers a robust and flexible approach to building your Wails applications. By understanding its structure and flow, you can leverage its capabilities to create efficient and customized build processes for your projects across Windows, Linux, and macOS platforms.