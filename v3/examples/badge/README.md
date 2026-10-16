# Welcome to Your New Wails3 Project!
Now that you have your project set up, it's time to explore the basic badge features that Wails3 offers on **macOS** and **Windows**.

## Running on macOS with private APIs

This example configures a translucent macOS backdrop. The webview transparency needed to reveal that backdrop requires the `private_mac_apis` build tag. Without it, the example runs with an opaque webview above the native backdrop.

From this example directory, run:

```bash
wails3 build -tags private_mac_apis
wails3 task run
```

The build command generates bindings and builds the frontend; the run task then launches the resulting app. It requires the Wails CLI, Node.js/npm, and the usual macOS build prerequisites. To use live reload instead, run `EXTRA_TAGS=private_mac_apis wails3 dev`.

Omit `-tags private_mac_apis` from the build command to run with public macOS APIs only. The tag has no effect on Windows, Linux, iOS, or Android. See the [shared private API guide](../README.md#private-macos-apis) for production builds and fallback details.

## Exploring Badge Features

### Creating the Service

First, initialize the badge service:

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/badge"

// Create a new badge service
badgeService := badge.New()

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(badgeService),
    },
})
```

## Badge Operations

### Setting a Badge

Set a badge on the application tile/dock icon:

#### Go
```go
// Set a default badge
badgeService.SetBadge("")

// Set a numeric badge
badgeService.SetBadge("3")

// Set a text badge
badgeService.SetBadge("New")
```

#### JS
```js
import {SetBadge} from "../bindings/github.com/wailsapp/wails/v3/pkg/services/badge/service";

// Set a default badge
SetBadge("")

// Set a numeric badge
SetBadge("3")

// Set a text badge
SetBadge("New")
```

### Removing a Badge

Remove the badge from the application icon:

#### Go
```go
badgeService.RemoveBadge()
```

#### JS
```js
import {RemoveBadge} from "../bindings/github.com/wailsapp/wails/v3/pkg/services/badge/service";

RemoveBadge()
```