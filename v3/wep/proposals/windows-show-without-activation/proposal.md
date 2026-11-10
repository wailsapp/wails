# Wails Enhancement Proposal (WEP)

## Show a Windows window without activating it

**WEP Number**: (leave blank, assigned on acceptance)  
**Status**: Draft  
**Author**: Alex Marshall (@AlexM1010)  
**Created**: 2026-10-29  
**Discussion**: none  
**Implementor**: Alex Marshall (@AlexM1010)  
**Target**: Wails v3

## Summary

Add `WindowsWindow.ShowWithoutActivation`. When it is set, `Show()` uses
`ShowWindow(SW_SHOWNA)` instead of `SW_SHOW`, so the window appears without
becoming the foreground window. This is the Windows counterpart of the existing
`Mac.PanelPreferences.NonActivating`.

## Motivation

On Windows, `Show()` always calls `ShowWindow(hwnd, SW_SHOW)`, which activates
the window. A status indicator, toast, HUD or on-screen helper therefore takes
keyboard focus from the application the person is typing into. macOS already
avoids this for non-activating panels (`orderFrontRegardless`). Windows has no
equivalent.

The workaround is to create the window hidden and call
`ShowWindow(SW_SHOWNA)` on `NativeWindow()`. That bypasses Wails' own show
path, which does more than reveal the HWND. It records `windowShown` and
`showRequested`, applies content protection, and shows the WebView2 controller
(directly, or through the efficiency-mode timeout). `navigationCompleted` hides
the controller of a window Wails believes was never shown, so a raw reveal can
show an HWND whose WebView stays invisible.

## Detailed Design

- New field on `WindowsWindow` (`webview_window_options.go`):

  ```go
  // ShowWithoutActivation makes Show display the window with SW_SHOWNA, so it
  // does not become the foreground window. Everything else Show does is
  // unchanged. It does not change hit-testing or later keyboard focus.
  // Default: false (SW_SHOW, which activates the window)
  ShowWithoutActivation bool
  ```

- In `windowsWebviewWindow.show()`, choose `SW_SHOWNA` when the option is set.
  The rest of the function is untouched.
- `Focus()` keeps activating the window, as it does for macOS panels. A click
  still activates the window unless the application also sets
  `WS_EX_NOACTIVATE` through the existing `ExStyle`.

The change is about six lines. It has been carried in a vendored beta.28 build.

## Non-Goals

- Making a window click-through or never-activatable. `ExStyle` already
  covers `WS_EX_TRANSPARENT` and `WS_EX_NOACTIVATE`.
- A cross-platform `ShowInactive()` method or a unified option.
- Changing `Focus()`, `Restore()` or the initial show of `Hidden: false`
  windows beyond what goes through `show()`.

## Platform Considerations

- **Windows**: `SW_SHOWNA` shows the window in its current size and position
  without activating it. This is supported on every Windows version Wails
  targets.
- **macOS**: already served by `Mac.PanelPreferences.NonActivating`. No change.
- **Linux**: not covered. GTK's `gtk_window_set_focus_on_map(FALSE)` could
  provide a similar option later, but Wayland compositors may ignore it.

## Pros/Cons

Pros: small, opt-in, keeps Wails' show bookkeeping, and matches existing macOS
behaviour. Cons: one more Windows-only option.

## Alternatives Considered

- **Raw `ShowWindow` on `NativeWindow()`**: skips the bookkeeping above and can
  leave the WebView invisible.
- **A `ShowInactive()` method on `Window`**: needs a decision for every platform
  and adds an API surface. A window that should never steal focus is better
  described when it is created.
- **Inferring it from `ExStyle & WS_EX_NOACTIVATE`**: implicit, and it changes
  behaviour for existing applications that already set that style.

## Backwards Compatibility

None. The option defaults to false, which keeps today's `SW_SHOW`.

## Security and Privacy

None. Not taking focus avoids keystrokes going to a window that appeared
unexpectedly.

## Test Plan

- A unit test for the chosen show command, if the maintainers want one (the
  call site would be split out to make it testable).
- Manual test: with a text editor focused, `Show()` a window that has the option.
  The editor keeps focus and the caret, and the window is visible with its
  content rendered. Without the option, the window takes focus as it does today.

## Reference Implementation

The diff against `master`, which will be submitted after acceptance:

```diff
-	w32.ShowWindow(w.hwnd, w32.SW_SHOW)
+	showCommand := w32.SW_SHOW
+	if w.parent.options.Windows.ShowWithoutActivation {
+		showCommand = w32.SW_SHOWNA
+	}
+	w32.ShowWindow(w.hwnd, showCommand)
```

## Maintenance Plan

One branch in `show()` and one documented field. The implementor will
maintain it.

## Conclusion

A small opt-in flag gives Windows the non-activating show macOS already has,
without applications bypassing Wails' show path.
