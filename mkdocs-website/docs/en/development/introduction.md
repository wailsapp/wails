# Introduction

!!! note This guide is a work in progress.

Thanks for wanting to help out with development of Wails! This guide will help
you get started.

## Getting Started

- Git clone this repository. Checkout the `v3-alpha` branch.
- Install the CLI: `cd v3/cmd/wails3 && go install`

- Optional: If you are wanting to use the build system to build frontend code,
  you will need to install [npm](https://nodejs.org/en/download).

## Building

For simple programs, you can use the standard `go build` command. It's also
possible to use `go run`.

Wails also comes with a build system that can be used to build more complex
projects. It utilises the awesome [Task](https://taskfile.dev) build system. For
more information, check out the task homepage or run `wails task --help`.

## Project layout

The project has the following structure:

    ```
    v3
    ├── cmd/wails3                  // CLI
    ├── examples                   // Examples of Wails apps
    ├── internal                   // Internal packages
    |   ├── runtime                // The Wails JS runtime
    |   └── templates              // The supported project templates
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

### Alpha Todo List

We are currently tracking known issues and tasks in the
[Alpha Todo List](https://github.com/orgs/wailsapp/projects/6). If you want to
help out, please check this list and follow the instructions in the
[Feedback](../getting-started/feedback.md) page.

### Adding window functionality

The preferred way to add window functionality is to add a new function to the
`pkg/application/webview_window.go` file. This should implement all the
functionality required for all platforms. Any platform specific code should be
called via a `webviewWindowImpl` interface method. This interface is implemented
by each of the target platforms to provide the platform specific functionality.
In some cases, this may do nothing. Once you've added the interface method,
ensure each platform implements it. A good example of this is the `SetMinSize`
method.

- Mac: `webview_window_darwin.go`
- Windows: `webview_window_windows.go`
- Linux: `webview_window_linux.go`

Most, if not all, of the platform specific code should be run on the main
thread. To simplify this, there are a number of `invokeSync` methods defined in
`application.go`.

### Updating the runtime

The runtime is located in `v3/internal/runtime`. When the runtime is updated,
the following steps need to be taken:

```shell
wails3 task runtime:build
```

### Events

Events are defined in `v3/pkg/events`. When adding a new event, the following
steps need to be taken:

- Add the event to the `events.txt` file
- Run `wails3 task events:generate`

There are a number of types of events: platform specific application and window
events + common events. The common events are useful for cross-platform event
handling, but you aren't limited to the "lowest common denominator". You can use
the platform specific events if you need to.

When adding a common event, ensure that the platform specific events are mapped.
An example of this is in `window_webview_darwin.go`:

```go
		// Translate ShouldClose to common WindowClosing event
		w.parent.On(events.Mac.WindowShouldClose, func(_ *WindowEventContext) {
			w.parent.emit(events.Common.WindowClosing)
		})
```

NOTE: We may try to automate this in the future by adding the mapping to the
event definition.

### Plugins

Plugins are a way to extend the functionality of your Wails application.

#### Creating a plugin

Plugins are standard Go structure that adhere to the following interface:

```go
type Plugin interface {
    Name() string
    Init(*application.App) error
    Shutdown()
    CallableByJS() []string
    InjectJS() string
}
```

The `Name()` method returns the name of the plugin. This is used for logging
purposes.

The `Init(*application.App) error` method is called when the plugin is loaded.
The `*application.App` parameter is the application that the plugin is being
loaded into. Any errors will prevent the application from starting.

The `Shutdown()` method is called when the application is shutting down.

The `CallableByJS()` method returns a list of exported functions that can be
called from the frontend. These method names must exactly match the names of the
methods exported by the plugin.

The `InjectJS()` method returns JavaScript that should be injected into all
windows as they are created. This is useful for adding custom JavaScript
functions that complement the plugin.

The built-in plugins can be found in the `v3/plugins` directory. Check them out
for inspiration.

## Tasks

The Wails CLI uses the [Task](https://taskfile.dev) build system. It is imported
as a library and used to run the tasks defined in `Taskfile.yaml`. The main
interfacing with Task happens in `v3/internal/commands/task.go`.

### Upgrading Taskfile

To check if there's an upgrade for Taskfile, run `wails3 task -version` and
check against the Task website.

To upgrade the version of Taskfile used, run:

```shell
wails3 task taskfile:upgrade
```

If there are incompatibilities then they should appear in the
`v3/internal/commands/task.go` file.

Usually the best way to fix incompatibilities is to clone the task repo at
`https://github.com/go-task/task` and look at the git history to determine what
has changed and why.

To check all changes have worked correctly, re-install the CLI and check the
version again:

```shell
wails3 task cli:install
wails3 task -version
```

## Opening a PR

Make sure that all PRs have a ticket associated with them providing context to
the change. If there is no ticket, please create one first. Ensure that all PRs
have updated the CHANGELOG.md file with the changes made. The CHANGELOG.md file
is located in the `mkdocs-website/docs` directory.

## Misc Tasks

### Upgrading Taskfile

The Wails CLI uses the [Task](https://taskfile.dev) build system. It is imported
as a library and used to run the tasks defined in `Taskfile.yaml`. The main
interfacing with Task happens in `v3/internal/commands/task.go`.

To check if there's an upgrade for Taskfile, run `wails3 task -version` and
check against the Task website.

To upgrade the version of Taskfile used, run:

```shell
wails3 task taskfile:upgrade
```

If there are incompatibilities then they should appear in the
`v3/internal/commands/task.go` file.

Usually the best way to fix incompatibilities is to clone the task repo at
`https://github.com/go-task/task` and look at the git history to determine what
has changed and why.

To check all changes have worked correctly, re-install the CLI and check the
version again:

```shell
wails3 task cli:install
wails3 task -version
```
