# Application

The application API assists in creating an application using the Wails
framework.

### New

API: `New(appOptions Options) *App`

`New(appOptions Options)` creates a new application using the given application
options . It applies default values for unspecified options, merges them with
the provided ones, initializes and returns an instance of the application.

In case of an error during initialization, the application is stopped with the
error message provided.

It should be noted that if a global application instance already exists, that
instance will be returned instead of creating a new one.

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

`Get()` returns the global application instance. It's useful when you need to
access the application from different parts of your code.

```go
    // Get the application instance
    app := application.Get()
```

### Capabilities

API: `Capabilities() capabilities.Capabilities`

`Capabilities()` retrieves a map of capabilities that the application currently
has. Capabilities can be about different features the operating system provides,
like webview features.

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

`Quit()` quits the application by destroying windows and potentially other
components.

```go
    // Quit the application
    app.Quit()
```

### IsDarkMode

API: `IsDarkMode() bool`

`IsDarkMode()` checks if the application is running in dark mode. It returns a
boolean indicating whether dark mode is enabled.

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

--8<--
./docs/en/API/application_window.md
./docs/en/API/application_menu.md
./docs/en/API/application_dialogs.md
./docs/en/API/application_events.md
./docs/en/API/application_screens.md
--8<--


## Options

```go title="pkg/application/application_options.go"
--8<--
../v3/pkg/application/application_options.go
--8<--
```
