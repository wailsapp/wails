# Wails Enhancement Proposal (WEP)

## Title

**Author**: Lea Anthony
**Created**: 2024-05-20

## Summary

Create a better API for controlling the appearance and functionality of window buttons in Wails.
This will only be available on Windows and macOS.

## Motivation

We currently do not fully support the ability to hide or disable the window buttons on the titlebar. 

## Detailed Design

1. A new enum will be added:

```go
      type ButtonState int

      const (
          ButtonEnabled  ButtonState = 0
          ButtonDisabled ButtonState = 1
          ButtonHidden   ButtonState = 2
      )
```

2. The following options will be added to the `WebviewWindowOptions` option struct:

```go
   MinimiseButtonState ButtonState
   MaximiseButtonState ButtonState
   CloseButtonState    ButtonState
```

The following options will be removed from the `WindowsWindow` options:

- DisableMinimiseButton
- DisableMaximiseButton

The following options will be removed from the `MacWindow` options:

- DisableMinimiseButton
- DisableMaximiseButton
- DisableCloseButton

3. The following methods will be added to the `Window` interface:

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

## Pros/Cons

### Pros

- We can support everything that's possible on both macOS and Windows

### Cons

- Windows works slightly different to macOS
- No Linux support (doesn't look like it's possible regardless of the solution)

## Alternatives Considered

The alternative is to draw your own titlebar, but this is a lot of work and often doesn't look good.

## Backwards Compatibility

This is not backwards compatible as we remove the old button disable options.

## Test Plan

As part of the implementation, the window example will be updated to test the functionality.

## Reference Implementation

There is a reference implementation as part of this proposal.

## Maintenance Plan

This feature will be maintained and supported by the Wails developers.

## Conclusion

This API would be a leap forward in giving developers greater control over their application window appearances.

