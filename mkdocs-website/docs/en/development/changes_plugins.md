## Plugins

Plugins are a way to extend the functionality of your Wails application.

### Creating a plugin

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
