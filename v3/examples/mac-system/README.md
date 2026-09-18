# System integration

This example exercises the system integration managers on `application.App`:

- `app.Permissions`: status, request and System Settings deep links for
  camera, microphone, screen recording, accessibility, location,
  notifications, input monitoring and full disk access.
- `app.Power`: a "Keep awake" toggle backed by `PreventSleep`, plus Low Power
  Mode, thermal state and battery information from `State()`.
- `app.Lifecycle`: `HoldTermination` and `SetSuddenTerminationEnabled`.
- `app.Env`: `Accessibility()`, `KeyboardLayout()` and `Locale()`.

The page updates live through `events.Mac.ApplicationDidChangePowerState`,
`ApplicationDidChangeThermalState`, `ApplicationDidChangeAccessibilitySettings`
(also mirrored as `events.Common.AccessibilitySettingsChanged`),
`ApplicationDidChangeKeyboardLayout` and `ApplicationDidChangeLocale`.

The window also sets `WebviewWindowOptions.Permissions` so web content may use
the microphone without the webview's own prompt while the camera still asks
(issue #6067).

## Running

```bash
go run .
```

Then:

- Click **Request** next to a permission that is *not determined*. macOS shows
  its prompt; the row updates when you answer. Prompts for camera, microphone
  and location require the usage-description keys in a bundled app's
  Info.plist (`NSCameraUsageDescription`, `NSMicrophoneUsageDescription`,
  `NSLocationUsageDescription`); unbundled `go run` binaries can only query
  those kinds. Notifications report *unsupported* when unbundled.
- Toggle **Keep awake** and run `pmset -g assertions` in a terminal: a
  `PreventUserIdleSystemSleep` assertion named "mac-system example keep awake"
  appears while the box is ticked.
- Toggle System Settings > Accessibility > Display > Reduce motion, switch the
  keyboard input source from the menu bar, or change the region in
  System Settings > General > Language & Region. Each change is logged at the
  bottom of the page.

## Notes

- `Request` blocks until the user answers, so call it from a goroutine other
  than the main thread (bound service methods already are).
- Sudden termination lets macOS kill the app instantly at logout. It is
  opt-in through the `NSSupportsSuddenTermination` Info.plist key; the toggle
  in this example changes the runtime state only.
- `Env.Locale()` reports the locale AppKit chose for the app, which depends
  on the `CFBundleLocalizations` the bundle declares; `Preferred` is always the
  user's full list.
