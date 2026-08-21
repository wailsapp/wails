# Notch Window Example

This macOS example demonstrates the high-level `NewNotchWindow` API. Wails
creates a shaped, non-activating panel, centres it on the camera housing, and
keeps the caller-requested webview size inside the native edges.

The example creates two independent windows. The second appears over the first
with a different size and animation duration, demonstrating that
multiple notch windows may coexist and be shown or hidden independently.
Set `NotchWindowOptions.Screen` to place an instance on a particular display;
without it, Wails uses the primary display.

From this directory, run:

```bash
go run .
```

Press Escape while either window has key focus to animate it out and hide it.
Move the pointer into a window to bring it forward and give its webview keyboard
focus without activating the application. Notch windows remain anchored to the
camera housing; neither the native background nor CSS drag regions can move
them. Hiding retains the loaded webview for the next notification but removes
its hover target, so the application must call `Show` to display it again.
