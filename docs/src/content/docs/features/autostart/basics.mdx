---
title: Autostart
description: Register your application to launch at user login on macOS, Windows, and Linux
sidebar:
  order: 1
---

import { Tabs, TabItem, Card, CardGrid } from "@astrojs/starlight/components";

## Autostart

`app.Autostart` registers your application to launch automatically when the user logs in. It picks the right native mechanism per platform and resolves symlinked install paths (Homebrew, Scoop) so registrations don't break when the binary is upgraded.

Registration takes effect on the **next login**, not immediately.

## Quick Start

```go
import "github.com/wailsapp/wails/v3/pkg/application"

// Register to launch at login
if err := app.Autostart.Enable(); err != nil {
    app.Logger.Error("autostart enable failed", "error", err)
}

// Stop launching at login
if err := app.Autostart.Disable(); err != nil {
    app.Logger.Error("autostart disable failed", "error", err)
}

// Check status
enabled, err := app.Autostart.IsEnabled()
```

## API

### `Enable`

Registers the application to launch at login with default options.

```go
func (m *AutostartManager) Enable() error
```

Calling `Enable` repeatedly is safe — the registration is overwritten each time, so it's correct to call on every startup if you've persisted the user's preference.

### `EnableWithOptions`

Registers with custom options.

```go
func (m *AutostartManager) EnableWithOptions(opts AutostartOptions) error
```

**`AutostartOptions`:**

| Field | Type | Description |
|---|---|---|
| `Identifier` | `string` | Overrides the auto-derived registration ID. See "Identifier" below. |
| `Arguments` | `[]string` | Extra arguments appended to the executable path when launched at login (e.g. `--hidden`). |

### `Disable`

Removes the autostart registration. Returns `nil` if the application wasn't registered — disable is idempotent.

```go
func (m *AutostartManager) Disable() error
```

### `IsEnabled`

Reports whether a registration exists. Fast — doesn't validate the registered path.

```go
func (m *AutostartManager) IsEnabled() (bool, error)
```

### `Status`

Returns full registration state.

```go
func (m *AutostartManager) Status() (AutostartStatus, error)
```

**`AutostartStatus`:**

| Field | Type | Description |
|---|---|---|
| `Enabled` | `bool` | Whether a registration exists. |
| `Path` | `string` | On-disk location of the registration artefact (plist path, `.desktop` path, or registry sub-key). Empty when `Enabled` is false. |
| `Strategy` | `AutostartStrategy` | Which mechanism registered the app (see [Platform Behaviour](#platform-behaviour)). |

## Platform Behaviour

<Tabs syncKey="platform">
  <TabItem label="macOS" icon="apple">
    Two mechanisms are used depending on how the app is packaged:

    - **macOS 13+, bundled `.app`**: `SMAppService.mainAppService`. Works for sandboxed apps and Mac App Store builds. No TCC automation prompt (the historical AppleScript approach triggered one).
    - **macOS pre-13, or unbundled binary**: a LaunchAgent plist written to `~/Library/LaunchAgents/<identifier>.plist` with `RunAtLoad=true`.

    `Status()` returns `AutostartStrategySMAppService` or `AutostartStrategyLaunchAgent` so callers can tell which path was taken.

    When the app upgrades from unbundled to bundled, both paths are checked on `Status()` and cleaned up by `Disable()` so an orphaned LaunchAgent doesn't keep launching the old build.
  </TabItem>

  <TabItem label="Windows" icon="seti:windows">
    A registry value is added under `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, with the value name set to the autostart `Identifier` and the data set to the quoted executable path plus any `Arguments`.

    Argument quoting follows `CommandLineToArgvW` rules (backslashes doubled before quotes) so paths with spaces or quotes round-trip correctly.

    `Status().Strategy` returns `AutostartStrategyRegistryRun`.
  </TabItem>

  <TabItem label="Linux" icon="linux">
    An XDG autostart entry is written to `$XDG_CONFIG_HOME/autostart/<identifier>.desktop` (defaulting to `~/.config/autostart/`) with:

    ```ini
    Type=Application
    Hidden=false
    X-GNOME-Autostart-enabled=true
    Exec=<executable> <arguments>
    ```

    The `Exec` field is escaped per the [freedesktop.org Desktop Entry spec](https://specifications.freedesktop.org/desktop-entry-spec/) — reserved characters (``"``, `` ` ``, `$`, `\\`) are backslash-escaped and the value is double-quoted when it contains whitespace.

    `Status().Strategy` returns `AutostartStrategyXDGAutostart`.
  </TabItem>

  <TabItem label="iOS / Android / server" icon="information">
    Not supported. All methods return `ErrAutostartNotSupported`. Use `errors.Is(err, application.ErrAutostartNotSupported)` to detect this cleanly:

    ```go
    if err := app.Autostart.Enable(); err != nil {
        if errors.Is(err, application.ErrAutostartNotSupported) {
            // hide the toggle in the UI
            return
        }
        // real failure — surface it
    }
    ```
  </TabItem>
</Tabs>

## Identifier

If `Options.Identifier` is empty, a default is derived from your app's name:

| Platform | Default identifier |
|---|---|
| macOS (bundled) | The app's bundle identifier, e.g. `com.example.MyApp` |
| macOS (unbundled) | `wails.autostart.<slug>`, where `<slug>` is derived from `application.Options.Name` |
| Windows | Slug of `application.Options.Name` (lowercase, non-`A-Za-z0-9._-` stripped, spaces become dashes) |
| Linux | Same slug as Windows |

Identifiers must match `^[A-Za-z0-9._-]+$` and be no longer than 200 characters. Reverse-DNS form is recommended for macOS (it matches how launchd Labels are conventionally written).

When `AutostartOptions.Identifier` is overridden, the same identifier is reused as the registry value name on Windows and the `.desktop` filename on Linux, so a single string identifies the registration cross-platform.

## Stale Detection

`Disable()` and `Status()` locate the registration by **matching the registered executable path against `os.Executable()` (resolved through any symlinks)**, not by looking up the identifier. This means:

- **Changing the identifier between releases is safe.** The old registration is still discoverable by `Status()` and cleaned up by `Disable()` — as long as the executable path is the same.
- **A second copy of the app at a different path won't clobber the first's registration.** Each binary location is tracked independently.
- **Symlinked installs (Homebrew, Scoop) are stable.** `filepath.EvalSymlinks` is applied to `os.Executable()` before matching, so a Homebrew upgrade that swaps the target doesn't strand the entry.

What this does *not* cover: if the user moves or renames the binary to an unrelated path, the old registration becomes orphaned (it points at the now-missing file). Apps that ship from a stable install path don't need to worry about this; apps shipped as portable single-file binaries should call `Disable()` before moving themselves, or always launch through a stable symlink.

## Example

```go
package main

import (
    "errors"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Restore the user's preference on startup
    if userPrefersAutostart() {
        if err := app.Autostart.Enable(); err != nil {
            if !errors.Is(err, application.ErrAutostartNotSupported) {
                app.Logger.Error("autostart", "error", err)
            }
        }
    }

    app.Run()
}
```

A complete runnable example with status / enable / disable buttons is in [`examples/autostart/`](https://github.com/wailsapp/wails/tree/master/v3/examples/autostart).
