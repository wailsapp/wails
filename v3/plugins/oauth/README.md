# oauth Plugin

This plugin provides the ability to initiate an OAuth authentication flow.

## Installation

Add the plugin to the `Plugins` option in the Applications options:

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/plugins/browser"
)

func main() {
    oAuthPlugin := oauth.NewPlugin(oauth.Config{
        Providers: []goth.Provider{
            github.New(
                os.Getenv("clientkey"),
                os.Getenv("secret"),
                "http://localhost:9876/auth/github/callback",
                "email",
                "profile"),
        },
    })

    app := application.New(application.Options{
    // ...
    Plugins: map[string]application.Plugin{
        "oauth": oAuthPlugin,
    },
    })
```

## Usage

### Go

You can start the flow by calling `Start()` on the plugin instance:

```go
	app.Events.On("github-login", func(e *application.WailsEvent) {
        oAuthPlugin.Start()
        oAuthWindow.Show()
    })
```

There is a working example of github auth in the `v3/examples` directory.

## Support

If you find a bug in this plugin, please raise a ticket on the Wails [Issue Tracker](https://github.com/wailsapp/wails/issues). 
