# Notch Notification Example

This macOS example demonstrates `NewNotchWindow` with a compact, stateful
system monitor. The frontend continuously updates real CPU, memory, and disk
telemetry from public macOS host and filesystem APIs while Wails owns the
native notch shape, placement, focus behaviour, and animation.

From this directory, run:

```bash
go run .
```

Use the **chevron** button or press Escape to animate the window away, then press
**Command+Shift+N** from any application to show the same webview again. The
graph continues from its existing JavaScript state, demonstrating that
`Hide` retains the loaded webview for the next presentation. Use **Quit** to
terminate the accessory application completely.

The bound service reads cumulative processor ticks and virtual-memory counters
from Mach, plus filesystem capacity through `statfs`. A production application
can use the same service pattern for its own events or telemetry without
changing the notch-window lifecycle.
