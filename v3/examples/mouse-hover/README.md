# Mouse Hover Example

This example demonstrates the `WindowMouseEnter` and `WindowMouseLeave` events, which are fired when the mouse cursor enters or leaves a window.

## Features Demonstrated

1. **Mouse Enter/Leave Events**: The example shows how to listen for `events.Common.WindowMouseEnter` and `events.Common.WindowMouseLeave` events on windows.

2. **FocusOnMouseEnter Option**: The secondary window demonstrates the `FocusOnMouseEnter` option, which automatically focuses the window when the mouse enters it. This is particularly useful for tray popup windows where you want the user to be able to interact with the window immediately without an initial click to focus.

## Use Cases

- **Tray Icon Popup Windows**: When a user hovers over a tray icon and a popup window appears, enabling `FocusOnMouseEnter` allows immediate interaction with the window contents.
- **Tooltip-style Windows**: For custom tooltip or hover card implementations that need to receive keyboard input.
- **Dashboard Widgets**: For multi-window applications where you want hover-to-focus behavior.

## Running the Example

```bash
cd v3/examples/mouse-hover
go run .
```

## Code Highlights

### Listening for Mouse Events

```go
window.OnWindowEvent(events.Common.WindowMouseEnter, func(e *application.WindowEvent) {
    app.Logger.Info("Mouse entered window!")
})

window.OnWindowEvent(events.Common.WindowMouseLeave, func(e *application.WindowEvent) {
    app.Logger.Info("Mouse left window!")
})
```

### Auto-Focus on Mouse Enter

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:             "Auto-Focus Window",
    FocusOnMouseEnter: true, // Window will focus when mouse enters
})
```

## Platform Support

Mouse enter/leave events work on all platforms:
- **Windows**: Uses `WM_MOUSEMOVE` with `TrackMouseEvent` and `WM_MOUSELEAVE`
- **macOS**: Uses `NSTrackingArea` with `NSTrackingActiveAlways` for tracking even when unfocused
- **Linux**: Uses GTK's `enter-notify-event` and `leave-notify-event` signals
