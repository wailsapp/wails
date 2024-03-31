### NewWebviewWindow

API: `NewWebviewWindow() *WebviewWindow`

`NewWebviewWindow()` creates a new Webview window with default options, and
returns it.

```go
    // Create a new webview window
    window := app.NewWebviewWindow()
```

### NewWebviewWindowWithOptions

API:
`NewWebviewWindowWithOptions(windowOptions WebviewWindowOptions) *WebviewWindow`

`NewWebviewWindowWithOptions()` creates a new webview window with custom
options. The newly created window is added to a map of windows managed by the
application.

```go
    // Create a new webview window with custom options
    window := app.NewWebviewWindowWithOptions(WebviewWindowOptions{
		Name: "Main",
        Title: "My Window",
        Width: 800,
        Height: 600,
    })
```

### OnWindowCreation

API: `OnWindowCreation(callback func(window *WebviewWindow))`

`OnWindowCreation()` registers a callback function to be called when a window is
created.

```go
    // Register a callback to be called when a window is created
    app.OnWindowCreation(func(window *WebviewWindow) {
        // Do something
    })
```

### GetWindowByName

API: `GetWindowByName(name string) *WebviewWindow`

`GetWindowByName()` fetches and returns a window with a specific name.

```go
    // Get a window by name
    window := app.GetWindowByName("Main")
```

### CurrentWindow

API: `CurrentWindow() *WebviewWindow`

`CurrentWindow()` fetches and returns a pointer to the currently active window
in the application. If there is no window, it returns nil.

```go
    // Get the current window
    window := app.CurrentWindow()
```
