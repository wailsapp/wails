# Wails v3 Plugin Guide

Wails v3 introduces the concept of plugins. A plugin is a self-contained module that can extend the functionality of your Wails application. 
This guide will walk you through the structure and functionality of a Wails plugin.

## Plugin Structure

A Wails plugin is a standard Go module, typically consisting of the following files:

- `plugin.go`: This is the main Go file where the plugin's functionality is implemented.
- `plugin.yaml`: This is the plugin's metadata file. It contains information about the plugin such as its name, author, version, and more.
- `assets/`: This directory contains any static assets that the plugin might need.
- `README.md`: This file provides documentation for the plugin.

## Plugin Implementation

In `plugin.go`, a plugin is defined as a struct that implements the `application.Plugin` interface. This interface requires the following methods:

- `Init()`: This method is called when the plugin is initialized.
- `Shutdown()`: This method is called when the application is shutting down.
- `Name()`: This method returns the name of the plugin.
- `CallableByJS()`: This method returns a list of method names that can be called from the frontend.

In addition to these methods, you can define any number of additional methods that implement the plugin's functionality. 
These methods can be called from the frontend using the `wails.Plugin()` function.

## Plugin Metadata

The `plugin.yaml` file contains metadata about the plugin. This includes the plugin's name, description, author, version, website, repository, and license.

## Plugin Assets

Any static assets that the plugin needs can be placed in the `assets/` directory. 
These assets can be accessed by the frontend by requesting them at the plugin base path.
This path is `/wails/plugins/<plugin-name>/`.

### Example

If a plugin named `logopack` has an asset named `logo.png`, the frontend can access it at `/wails/plugins/logopack/logo.png`.

## Plugin Documentation

The `README.md` file provides documentation for the plugin. This should include instructions on how to install and use the plugin, as well as any other information that users of the plugin might find useful.

## Example

Here's an example of how to use a plugin in your Wails application:

```go
    Plugins: map[string]application.Plugin{
        "log": log.NewPlugin(),
    },
```
And here's how you can call a plugin method from the frontend:

```js
    wails.Plugin("log","Debug","hello world")
```

## Support

If you encounter any issues with a plugin, please raise a ticket in the plugin's repository. 

!!! note
    The Wails team does not provide support for third-party plugins.