# Application

![Application API Image](application_api.png)

The Application API provides access to the main application and its functionalities. 

### `New(appOptions Options) *App`

The `New()` function is used to create a new Wails application with the provided options. This is the entry point for creating a Wails app.

```go
app := application.New(wails.Options{
    Name:     "My Wails App",
    Width:    800,
    Height:   600,
    Resizable: true,
})
```

### `Run() error`

The `Run()` method starts the main event loop of the Wails application. This should be called after all windows, system trays, and other components have been set up.

```go
err := app.Run()
if err != nil {
    // Handle error
}
```

### `Quit()`

The `Quit()` method gracefully shuts down the Wails application.

```go
app.Quit()
```

### `SetIcon(icon []byte)`

The `SetIcon()` method sets the icon for the Wails application.

```go
app.SetIcon(iconData)
```

### `SetMenu(menu *Menu)`

The `SetMenu()` method sets the main application menu.

```go
menu := app.NewMenu()
menu.Append(&wails.MenuItem{
    Label: "File",
    Submenu: [...],
})
app.SetMenu(menu)
```

### `ShowAboutDialog()`

The `ShowAboutDialog()` method displays the application's about dialog.

```go
app.ShowAboutDialog()
```

### `CurrentWindow() *WebviewWindow`

The `CurrentWindow()` method returns the current application window.

```go
currentWindow := app.CurrentWindow()
```

### `NewWebviewWindow() *WebviewWindow`

The `NewWebviewWindow()` method creates a new Webview window.

```go
window := app.NewWebviewWindow()
```

### `NewSystemTray() *SystemTray`

The `NewSystemTray()` method creates a new system tray.

```go
systemTray := app.NewSystemTray()
```

### `Clipboard() *Clipboard`

The `Clipboard()` method returns the application's clipboard instance, which can be used to read and write to the system clipboard.

```go
clipboard := app.Clipboard()
clipboard.Write("Hello, world!")
```

### `Environment() EnvironmentInfo`

The `Environment()` method returns information about the current operating system and runtime environment.

```go
env := app.Environment()
fmt.Println(env.OS, env.Arch, env.Debug)
```

### `OnShutdown(f func())`

The `OnShutdown()` method adds a method to be run when the application is shutting down.

```go
app.OnShutdown(func() {
    // Cleanup resources
})
```

### `BrowserOpenURL(url string) error`

The `BrowserOpenURL()` method opens the default system browser and navigates to the specified URL.

```go
err := app.BrowserOpenURL("https://www.example.com")
if err != nil {
    // Handle error
}
```

### `BrowserOpenFile(path string) error`

The `BrowserOpenFile()` method opens the default system application associated with the specified file.

```go
err := app.BrowserOpenFile("/path/to/file.pdf")
if err != nil {
    // Handle error
}
```

### `RegisterContextMenu(name string, menu *Menu)`

The `RegisterContextMenu()` method registers a context menu with the given name. This menu can be referenced by the 
frontend to display custom context menus.

```go
contextMenu := app.NewMenu()
// Add menu items
app.RegisterContextMenu("my-context-menu", contextMenu)
```

### `GetWindowByName(name string) Window`

The `GetWindowByName()` method returns the window with the given name.

```go
window := app.GetWindowByName("My Window")
```

### `OnWindowCreation(callback func(window Window))`

The `OnWindowCreation()` method registers a callback to be called when a new window is created.

```go
app.OnWindowCreation(func(window Window) {
    // Handle window creation
})
```
```