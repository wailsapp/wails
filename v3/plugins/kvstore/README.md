# KVStore Plugin

This plugin provides a simple key/value store for your Wails applications.

## Installation

Add the plugin to the `Plugins` option in the Applications options:

```go
package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/plugins/kvstore"
)

func main() {
	kvstorePlugin := kvstore.NewPlugin(&kvstore.Config{
		Filename: "myapp.db",
	})
	app := application.New(application.Options{
		// ...
		Plugins: map[string]application.Plugin{
			"kvstore": kvstorePlugin,
		},
	})

```

## Usage

### Go

You can call the methods exported by the plugin directly:

```go
    err := kvstore.Set("url", "https://www.google.com")
    if err != nil {
        // handle error
    }
    url := kvstore.Get("url").(string)

	// If you have not enables AutoSave, you will need to call Save() to persist the changes
    err = kvstore.Save()
    if err != nil {
        // handle error
    }
```

### Javascript

You can call the methods from the frontend using the Plugin method:

```js
wails.Plugin("kvstore", "Set", "url", "https://www.google.com");
wails.Plugin("kvstore", "Get", "url").then((url) => {});

// or
wails.Plugin("browser", "OpenFile", "/path/to/file");
```

## Support

If you find a bug in this plugin, please raise a ticket on the Wails
[Issue Tracker](https://github.com/wailsapp/wails/issues).
