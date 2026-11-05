# Welcome to Your New Wails3 Project!
Now that you have your project set up, it's time to explore the custom badge features that Wails3 offers on **Windows**.

## Running on macOS with private APIs

This example configures a translucent macOS backdrop. The webview transparency needed to reveal that backdrop requires the `private_mac_apis` build tag. Without it, the example runs with an opaque webview above the native backdrop.

From this example directory, run:

```bash
wails3 build -tags private_mac_apis
wails3 task run
```

The build command generates bindings and builds the frontend; the run task then launches the resulting app. It requires the Wails CLI, Node.js/npm, and the usual macOS build prerequisites. To use live reload instead, run `EXTRA_TAGS=private_mac_apis wails3 dev`.

Omit `-tags private_mac_apis` from the build command to run with public macOS APIs only. The tag has no effect on Windows, Linux, iOS, or Android. See the [shared private API guide](../README.md#private-macos-apis) for production builds and fallback details.

## Exploring Custom Badge Features

### Creating the Service with Custom Options (Windows Only)

On Windows, you can customize the badge appearance with various options:

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/badge"
import "image/color"

// Create a badge service with custom options
options := badge.Options{
    TextColour:       color.RGBA{255, 255, 255, 255}, // White text
    BackgroundColour: color.RGBA{0, 0, 255, 255},     // Green background
    FontName:         "consolab.ttf",                 // Bold Consolas font
    FontSize:         20,                             // Font size for single character
    SmallFontSize:    14,                             // Font size for multiple characters
}

badgeService := badge.NewWithOptions(options)

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(badgeService),
    },
})
```

## Badge Operations

### Setting a Badge

Set a badge on the application tile/dock icon with the global options applied:

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

### Setting a Custom Badge

Set a badge on the application tile/dock icon with one-off options applied:

#### Go
```go
// Set a default badge
badgeService.SetCustomBadge("")

// Set a numeric badge
badgeService.SetCustomBadge("3")

// Set a text badge
badgeService.SetCustomBadge("New")
```

#### JS
```js
import {SetCustomBadge} from "../bindings/github.com/wailsapp/wails/v3/pkg/services/badge/service";

const options = {
   BackgroundColour: RGBA.createFrom({
         R: 0,
         G: 255,
         B: 255,
         A: 255,
      }),
      FontName: "arialb.ttf", // System font
      FontSize: 16,
      SmallFontSize: 10,
      TextColour: RGBA.createFrom({
         R: 0,
         G: 0,
         B: 0,
         A: 255,
      }),
}

// Set a default badge
SetCustomBadge("", options)

// Set a numeric badge
SetCustomBadge("3", options)

// Set a text badge
SetCustomBadge("New", options)
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