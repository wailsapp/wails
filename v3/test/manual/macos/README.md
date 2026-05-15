# macOS Manual Tests

Manual test programs for macOS-specific window behavior. These tests can't be
automated in `go test` because they rely on AppKit-level activation state and
visual cues that only a human can verify.

## Running

```bash
cd v3/test/manual/macos/<test-name>
go run .
```

## Tests

### non-activating-panel

Verifies `MacWindow.NonActivatingPanel` and the underlying `NSPanel` migration.

Opens two windows side by side: one with `NonActivatingPanel: true` (the
"Panel"), and one without (the "Normal" window, used as a control).

| # | Action | Expected behavior |
|---|--------|-------------------|
| 1 | Bring another app to the foreground (e.g. Finder) | Both windows STAY VISIBLE. Regression check for `hidesOnDeactivate` — if the panel vanishes on app switch, the `NSPanel` default leaked through. |
| 2 | Click the input field in the **Panel** window | Cursor blinks in the input; the *other* app's menu bar stays in place. The Wails app does not become active in the dock. |
| 3 | Click the input field in the **Normal** window | Cursor blinks in the input; Wails *does* activate (menu bar switches, dock icon highlights). This is the control case. |
| 4 | Open the Window menu in Wails' menu bar (after clicking the Normal window) | Only the Normal window appears in the list. The Panel never becomes the app's main window. |
| 5 | Type into the Panel's input | Characters appear (proves the panel becomes key even though it doesn't activate the app). |
| 6 | Close the Panel via Cmd+W, then quit and re-run | No crash. Regression check for `releasedWhenClosed` — `NSPanel` defaults that to `YES` and the Go side keeps a raw pointer past `[close]`. |

### Notes

- Step 2 is the headline behavior. If it fails, the `windowFocus` /
  `activateIgnoringOtherApps:` guard or the `NSWindowStyleMaskNonactivatingPanel`
  bit isn't being honored.
- Step 1 is the regression most likely to be re-introduced — easy to forget
  that `NSPanel` differs from `NSWindow` on `hidesOnDeactivate`.
- Step 4 verifies the conditional `canBecomeMainWindow` override in
  `webview_window_darwin.m`.
