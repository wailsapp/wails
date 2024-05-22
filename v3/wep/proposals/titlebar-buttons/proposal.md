# Wails Enhancement Proposal (WEP)

## Customising Window Controls in Wails

**Author**: Lea Anthony
**Created**: 2024-05-20

## Summary

This is a proposal for an API to control the appearance and functionality of window controls in Wails.
This will only be available on Windows and macOS.

## Motivation

We currently do not fully support the ability to customise window controls. 

## Detailed Design

### Controlling Button State

1. A new enum will be added:

```go
      type ButtonState int

      const (
          ButtonEnabled  ButtonState = 0
          ButtonDisabled ButtonState = 1
          ButtonHidden   ButtonState = 2
      )
```

2. These options will be added to the `WebviewWindowOptions` option struct:

```go
   MinimiseButtonState ButtonState
   MaximiseButtonState ButtonState
   CloseButtonState    ButtonState
```

3. These options will be removed from the current Windows/Mac options:

- DisableMinimiseButton
- DisableMaximiseButton
- DisableCloseButton

4. These methods will be added to the `Window` interface:

```go
   SetMinimizeButtonState(state ButtonState)
   SetMaximizeButtonState(state ButtonState)
   SetCloseButtonState(state ButtonState)
```

The settings translate to the following functionality on each platform:

|                       | Windows                | Mac                    |
|-----------------------|------------------------|------------------------|
| Disable Min/Max/Close | Disables Min/Max/Close | Disables Min/Max/Close |
| Hide Min              | Disables Min           | Hides Min button       |
| Hide Max              | Disables Max           | Hides Max button       |
| Hide Close            | Hides all controls     | Hides Close            |

Note: On Windows, it is not possible to hide the Min/Max buttons individually.
However, disabling both will hide both of the controls and only show the 
close button. 

### Controlling Window Style (Windows)

As Windows currently does not have much in the way of controlling the style of the
titlebar, a new option will be added to the `WebviewWindowOptions` option struct:

```go
   ExStyle int
```

If this is set, then the new Window will use the style specified in the `ExStyle` field.

Example:
```go
package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

func main() {
    app := application.New(application.Options{
		Name:        "My Application",
	})
	
    app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
        Windows: application.WindowsWindow{
            ExStyle: w32.WS_EX_TOOLWINDOW | w32.WS_EX_NOREDIRECTIONBITMAP | w32.WS_EX_TOPMOST,
        },
	})
	
	app.Run()
}
```

## Pros/Cons

### Pros

- We bring much needed customisation capabilities to both macOS and Windows

### Cons

- Windows works slightly different to macOS
- No Linux support (doesn't look like it's possible regardless of the solution)

## Alternatives Considered

The alternative is to draw your own titlebar, but this is a lot of work and often doesn't look good.

## Backwards Compatibility

This is not backwards compatible as we remove the old "disable button" options.

## Test Plan

As part of the implementation, the window example will be updated to test the functionality.

## Reference Implementation

There is a reference implementation as part of this proposal.

## Maintenance Plan

This feature will be maintained and supported by the Wails developers.

## Conclusion

This API would be a leap forward in giving developers greater control over their application window appearances.

