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
          ButtonActive   ButtonState = 0
          ButtonInactive ButtonState = 1
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

|                       | Windows                      | Mac                    |
|-----------------------|------------------------------|------------------------|
| Disable Min/Max/Close | Disables Min/Max/Close       | Disables Min/Max/Close |
| Hide Min              | Hides both Min + Max buttons | Hides Min button       |
| Hide Max              | Hides both Min + Max buttons | Hides Max button       |
| Hide Close            | Hides all buttons            | Hides Close            |

## Pros/Cons

### Pros

1. We can support everything that's possible on both macOS and Windows

### Cons

1. Windows works slightly different to macOS

## Alternatives Considered

The alternative is to draw your own titlebar, but this is a lot of work and often doesn't look good.

## Backwards Compatibility

This is not backwards compatible as we remove the old `Disable*Button` options.

## Test Plan

As part of the implementation, an example will be created in `v3/examples` to test the functionality.

## Reference Implementation

There is a reference implementation as part of this proposal.

## Maintenance Plan

This feature will be maintained and supported by the Wails developers.

## Conclusion

This API would be a leap forward in giving developers greater control over their application window appearances.

