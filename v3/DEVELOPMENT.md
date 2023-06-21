# Development

**This guide is a work in progress.**

Thanks for wanting to help out with development of Wails! This guide will help you get started.

## Getting Started

- Git clone this repository. Checkout the `v3-alpha` branch.
- Install the CLI: `cd v3/cmd/wails && go install`

- Optional: If you are wanting to use the build system to build frontend code, you will need to install [npm](https://nodejs.org/en/download).

## Building

For simple programs, you can use the standard `go build` command. It's also possible to use `go run`.

Wails also comes with a build system that can be used to build more complex projects. It utilises the awesome [Task](https://taskfile.dev) build system.
For more information, check out the task homepage or run `wails task --help`. 

## Project layout

The project has the following structure:
    
    ```
    v3
    ├── cmd/wails                  // CLI
    ├── examples                   // Examples of Wails apps 
    ├── internal                   // Internal packages
    ├── pkg
    |   ├── application            // The core Wails library
    |   └── events                 // The event definitions
    |   └── mac                    // macOS specific code used by plugins
    |   └── w32                    // Windows specific code
    ├── plugins                    // Supported plugins
    ├── tasks                      // General tasks
    └── Taskfile.yaml              // Development tasks configuration
    ```

## Development

### Updating the runtime

The runtime is located in `v3/internal/runtime`. When the runtime is updated, the following steps need to be taken:

```shell
wails task runtime:build
```