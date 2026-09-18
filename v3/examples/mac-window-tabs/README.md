# mac-window-tabs

This example showcases macOS window tabbing using `MacWindowTabbingMode`.

Window tabbing is a macOS-only feature (NSWindow tabbing, 10.12+), so this
example is macOS only.

## Running on macOS with private APIs

This example configures a translucent macOS backdrop. The webview transparency needed to reveal that backdrop requires the `private_mac_apis` build tag. Without it, the example runs with an opaque webview above the native backdrop.

From this example directory, run:

```bash
GOWORK=off wails3 build -tags private_mac_apis
GOWORK=off wails3 task run
```

The build command generates bindings and builds the frontend; the run task then launches the resulting app. It requires the Wails CLI, Node.js/npm, and the usual macOS build prerequisites. `GOWORK=off` selects this example’s module and its local Wails replacement. To use live reload instead, run `GOWORK=off EXTRA_TAGS=private_mac_apis wails3 dev`.

Omit `-tags private_mac_apis` from the build command to run with public macOS APIs only. The tag has no effect on Windows, Linux, iOS, or Android. See the [shared private API guide](../README.md#private-macos-apis) for production builds and fallback details.

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
