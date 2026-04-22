# mac-window-tabs

This example showcases macOS window tabbing using `MacWindowTabbingMode`.

## Running

```bash
go run .
```

## What to Expect

- Two windows are created with `TabbingMode` set to `MacWindowTabbingModePreferred`.
- On macOS 10.12+, the windows will automatically tab together if tabbing is enabled in System Settings.
- You can also use Window > Merge All Windows to force tabs.

## Relevant Code

See the macOS window options in [main.go](main.go).
