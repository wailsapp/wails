# mac-window-tabs

This example showcases macOS window tabbing using `MacWindowTabbingMode`.

Window tabbing is a macOS-only feature (NSWindow tabbing, 10.12+), so this
example is macOS only.

## Running

```bash
task dev
```

This uses the `wails3` CLI (via the Taskfile) to generate bindings, build the
frontend, and run the app with live reload. `task run` builds and runs a
non-dev binary instead.

> The `go.mod` includes a `replace` directive pointing at the local Wails
> module, because `MacWindowTabbingMode` is not yet in a published release.
> `go run .` on its own will not work: it skips binding generation and the
> frontend build.

## What to Expect

Four windows open on launch:

- **Tabbing Enabled** and **Opens From Tabbing Enabled** use
  `MacWindowTabbingModePreferred`. On macOS 10.12+ these merge into a single
  tabbed window.
- **Tabbing Disabled** and **Opens From Tabbing Disabled** use
  `MacWindowTabbingModeDisallowed`. These stay as separate windows and never
  tab.

You can also use Window > Merge All Windows to force the tabbable windows
together.

## Relevant Code

See the macOS window options in [main.go](main.go).
