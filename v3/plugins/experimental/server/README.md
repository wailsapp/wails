# Server Plugin

This plugin provides a simple server for your Wails applications to make them accessible over the local network.
Bidirectional communication occurs over a websocket connection.

## Installation

Add the plugin to the `Plugins` option in the Applications options:

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/plugins/server"
)

func main() {
    app := application.New(application.Options{
        // ...
        Plugins: map[string]application.Plugin{
            "server": server.NewPlugin(&server.Config{
                Host: "0.0.0.0",
                Port: 31115,
            }),
        },
    })

```


## Support

If you find a bug in this plugin, please raise a ticket on the Wails [Issue Tracker](https://github.com/wailsapp/wails/issues).
