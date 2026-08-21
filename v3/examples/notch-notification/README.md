# Notch Notification Example

This macOS example demonstrates `NewNotchWindow` with a compact, stateful
system monitor. The frontend continuously updates simulated CPU, memory, and
disk telemetry while Wails owns the native notch shape, placement, focus
behaviour, and animation.

From this directory, run:

```bash
go run .
```

Use **Hide** or press Escape to animate the window away, then press
**Command+Shift+N** from any application to show the same webview again. The
graphs continue from their existing JavaScript state, demonstrating that
`Hide` retains the loaded webview for the next presentation. Use **Quit** to
terminate the accessory application completely.

The telemetry is intentionally simulated so the example remains
self-contained. A real application can replace it with events or a bound
service without changing the notch-window lifecycle.
