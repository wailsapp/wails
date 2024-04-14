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

### Options

```go
type Config struct {
    Filename string
    AutoSave bool
}
```

- `Filename` - The name of the file to store the key/value pairs in. This file will be created in the application's data directory.
- `AutoSave` - If true, the store will be saved to disk after every change. If false, you will need to call `Save()` to persist the changes.

## Usage

### Go

You can call the methods exported by the plugin directly:

```go
    // Set a key
    err := kvstore.Set("url", "https://www.google.com")
    if err != nil {
        // handle error
    }
	// Get a key
    url := kvstore.Get("url").(string)
	
	// Delete a key
    err = kvstore.Delete("url")
    if err != nil {
        // handle error
    }
    
	// If you have not enables AutoSave, you will need to call Save() to persist the changes
    err = kvstore.Save()
    if err != nil {
        // handle error
    }
```

### Javascript

You can call the methods from the frontend using the Plugin method:

```js
    wails.Plugin("kvstore","Set", "url", "https://www.google.com")
    wails.Plugin("kvstore","Get", "url").then((url) => {
    
    })
    wails.Plugin("kvstore","Delete", "url").then((url) => {
    
    })
    
    // or
    wails.Plugin("browser","OpenFile","/path/to/file")
```

## Support

If you find a bug in this plugin, please raise a ticket on the Wails [Issue Tracker](https://github.com/wailsapp/wails/issues). 
