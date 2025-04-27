# Welcome to Your New Wails3 Project!
Now that you have your project set up, it's time to explore the basic badge features that Wails3 offers on **macOS** and **Windows**.

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