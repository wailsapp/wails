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

A single window opens on launch. It uses `MacWindowTabbingModePreferred`, so it
is willing to accept new tabs. Two buttons drive the demo:

- **Open tabbed window** opens a window with `MacWindowTabbingModePreferred`. On
  macOS 10.12+ it merges into the current window as a new tab.
- **Open non-tabbed window** opens a window with `MacWindowTabbingModeDisallowed`.
  It always opens as a separate window and never tabs, even via Window > Merge
  All Windows.

Open a mix of both to see the difference: tabbed windows stack into one titled
tab bar, while non-tabbed windows stay independent.

## Relevant Code

See the macOS window options in [main.go](main.go).
