# Next Steps

Congratulations on generating a plugin. This guide will help you author your plugin
and provide some tips on how to get started.

## Plugin Structure

The plugin is a standard Go module that adheres to the following interface:

```go
type Plugin interface {
	Name() string
	Init(app *App) error
	Shutdown()
}
```

The `Name()` method returns the name of the plugin. It should follow the Go module naming convention
and have a prefix of `wails-plugin-`, e.g. `github.com/myuser/wails-plugin-example`.

The `Init()` method is called when the plugin is loaded. The `App` parameter is a pointer to the
main application struct. This may be used for showing dialogs, listening for events or even opening 
new windows. The `Init()` method should return an error if it fails to initialise. This method is
called synchronously so the application will not start until it returns.

The `Shutdown()` method is called when the application is shutting down. This is a good place to
perform any cleanup. This method is called synchronously so the application will not exit completely until
it returns.

## Plugin Directory Structure

The plugin directory structure is as follows:

```
plugin-name
├── models.d.ts
├── plugin.js
├── plugin.go
├── README.md
├── go.mod
├── go.sum
└── plugin.toml
```

### `plugin.go`

This file contains the plugin code. It should contain a struct that implements the `Plugin` interface
and a `NewPlugin()` method that returns a pointer to the struct. Methods are exported by capitalising
the first letter of the method name. These methods may be called from the frontend. If methods
accept or return structs, these structs must be exported. 

### `plugin.js`

This file should contain any JavaScript code that may help developers use the plugin.
In the example plugin, this file contains function wrappers for the plugin methods.
It's good to include JSDocs as that will help developers using your plugin.

### `models.d.ts`

This file should contain TypeScript definitions for any structs that are passed
or returned from the plugin. 
`
### `plugin.toml`

This file contains the plugin metadata. It is important to fill this out correctly
as it will be used by the Wails CLI.

### `README.md`

This file should contain a description of the plugin and how to use it. It should
also contain a link to the plugin repository and how to report bugs.

### `go.mod` and `go.sum`

These are standard Go module files. The package name in `go.mod` should match the
name of the plugin, e.g. `github.com/myuser/wails-plugin-example`.

## Promoting your Plugin

Once you have created your plugin, you should promote it on the Wails Discord server
in the `#plugins` channel. You should also open a PR to promote your plugin on the Wails
website. Update the `website/content/plugins.md` file and add your plugin to the list.