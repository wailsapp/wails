# Browser Plugin

This plugin provides the ability to open a URL or local file in the default browser.

## Installation

Add the plugin to the `Plugins` option in the Applications options:

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/plugins/browser"
)

func main() {
  browserPlugin := browser.NewPlugin()
  app := application.New(application.Options{
    // ...
    Plugins: map[string]application.Plugin{
        "browser": browserPlugin,
    },
  })
```

## Usage

### Go

You can call the methods exported by the plugin directly:

```go
    browserPlugin.OpenURL("https://www.google.com")
    // or
    browserPlugin.OpenFile("/path/to/file")
```

### Javascript

You can call the methods from the frontend using the Plugin method:

```js
    wails.Plugin("browser","OpenURL","https://www.google.com")
    // or
    wails.Plugin("browser","OpenFile","/path/to/file")
```

## Support

If you find a bug in this plugin, please raise a ticket on the Wails [Issue Tracker](https://github.com/wailsapp/wails/issues). 
