# Application

The application API assists in creating an application using the Wails framework.

### New

API: `New(appOptions Options) *App`

`New(appOptions Options)` creates a new application using the given application options . It applies default values for unspecified options, merges them with the provided ones, initializes and returns an instance of the application.

In case of an error during initialization, the application is stopped with the error message provided.

It should be noted that if a global application instance already exists, that instance will be returned instead of creating a new one.


```go title="main.go" hl_lines="6-9"
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Name:        "WebviewWindow Demo",
		// Other options
    })
    
	// Rest of application
}
```

### Get

`Get()` returns the global application instance. It's useful when you need to access the application from different parts of your code.

```go
    // Get the application instance
    app := application.Get()
```

### Capabilities

API: `Capabilities() capabilities.Capabilities`

`Capabilities()` retrieves a map of capabilities that the application currently has. Capabilities can be about different features the operating system provides, like webview features.

```go
    // Get the application capabilities
    capabilities := app.Capabilities()
	if capabilities.HasNativeDrag {
		// Do something
    }   
```

### GetPID

API: `GetPID() int`

`GetPID()` returns the Process ID of the application.

```go
    pid := app.GetPID()
```

### Run

API: `Run() error`

`Run()` starts the execution of the application and its components.

```go
    app := application.New(application.Options{ 
	    //options 
	})
    // Run the application
    err := application.Run()
    if err != nil {
        // Handle error
    }   
```

### Quit

API: `Quit()`

`Quit()` quits the application by destroying windows and potentially other components.

```go
    // Quit the application
    app.Quit()
```

### IsDarkMode

API: `IsDarkMode() bool`

`IsDarkMode()` checks if the application is running in dark mode. It returns a boolean indicating whether dark mode is enabled.

```go
    // Check if dark mode is enabled
    if app.IsDarkMode() {
        // Do something
    }
```

### Hide

API: `Hide()`

`Hide()` hides the application window.

```go
    // Hide the application window
    app.Hide()
```

### Show

API: `Show()`

`Show()` shows the application window.

```go
    // Show the application window
    app.Show()
```

### NewWebviewWindow

API: `NewWebviewWindow() *WebviewWindow`

`NewWebviewWindow()` creates a new Webview window with default options, and returns it.

```go
    // Create a new webview window
    window := app.NewWebviewWindow()
```

### NewWebviewWindowWithOptions

API: `NewWebviewWindowWithOptions(windowOptions WebviewWindowOptions) *WebviewWindow`

`NewWebviewWindowWithOptions()` creates a new webview window with custom options. The newly created window is added to a map of windows managed by the application.

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

`OnWindowCreation()` registers a callback function to be called when a window is created.

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

`CurrentWindow()` fetches and returns a pointer to the currently active window in the application. If there is no window, it returns nil.

```go
    // Get the current window
    window := app.CurrentWindow()
```

### NewSystemTray

API: `NewSystemTray() *SystemTray`

`NewSystemTray()` creates and returns a new system tray instance.

```go
    // Create a new system tray
    tray := app.NewSystemTray()
```


### NewMenu

API: `NewMenu() *Menu`

This method, belonging to App struct and not Menu struct, also initialises and returns a new `Menu`.

```go
    // Create a new menu
    menu := app.NewMenu()
```

### RegisterContextMenu

API: `RegisterContextMenu(name string, menu *Menu)`

`RegisterContextMenu()` registers a context menu with a given name. This menu can be used later in the application.

```go

    // Create a new menu
    ctxmenu := app.NewMenu()

    // Register the menu as a context menu
    app.RegisterContextMenu("MyContextMenu", ctxmenu)
```

### SetMenu

API: `SetMenu(menu *Menu)`

`SetMenu()` sets the menu for the application. On Mac, this will be the global menu. For Windows and Linux, this will be the default menu for any new window created.

```go
    // Create a new menu
    menu := app.NewMenu()

    // Set the menu for the application
    app.SetMenu(menu)
```

### ShowAboutDialog

API: `ShowAboutDialog()`

`ShowAboutDialog()` shows an "About" dialog box. It can show the application's name, description and icon.

```go
    // Show the about dialog
    app.ShowAboutDialog()
```

### Info

`InfoDialog()` creates and returns a new instance of `MessageDialog` with an `InfoDialogType`. This dialog is typically used to display informational messages to the user.

### Question

`QuestionDialog()` creates and returns a new instance of `MessageDialog` with a `QuestionDialogType`. This dialog is often used to ask a question to the user and expect a response.


### Warning

`WarningDialog()` creates and returns a new instance of `MessageDialog` with a `WarningDialogType`. As the name suggests, this dialog is primarily used to display warning messages to the user.

### Error

`ErrorDialog()` creates and returns a new instance of `MessageDialog` with an `ErrorDialogType`. This dialog is designed to be used when you need to display an error message to the user.

### OpenFile

`OpenFileDialog()` creates and returns a new `OpenFileDialogStruct`. This dialog prompts the user to select one or more files from their file system.

### SaveFile

`SaveFileDialog()` creates and returns a new `SaveFileDialogStruct`. This dialog prompts the user to choose a location on their file system where a file should be saved.

### OpenDirectory

`OpenDirectoryDialog()` creates and returns a new instance of `MessageDialog` with an `OpenDirectoryDialogType`. This dialog enables the user to choose a directory from their file system.


### On

API: `On(eventType events.ApplicationEventType, callback func(event *Event)) func()`

`On()` registers an event listener for specific application events. The callback function provided will be triggered when the corresponding event occurs. The function returns a function that can be called to remove the listener.

### RegisterHook

API: `RegisterHook(eventType events.ApplicationEventType, callback func(event *Event)) func()`

`RegisterHook()` registers a callback to be run as a hook during specific events. These hooks are run before listeners attached with `On()`. The function returns a function that can be called to remove the hook.

### GetPrimaryScreen

API: `GetPrimaryScreen() (*Screen, error)`

`GetPrimaryScreen()` returns the primary screen of the system.

### GetScreens

API: `GetScreens() ([]*Screen, error)`

`GetScreens()` returns information about all screens attached to the system.

This is a brief summary of the exported methods in the provided `App` struct. Do note that for more detailed functionality or considerations, refer to the actual Go code or further internal documentation.

## Options

```go title="application_options.go"
--8<--
../v3/pkg/application/options_application.go
--8<--
```