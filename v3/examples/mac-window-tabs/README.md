# mac-window-tabs

This example showcases macOS window tabbing: the `MacWindowTabbingMode`
window option and the tab group API on `WebviewWindow` and
`MacWindowTabGroup`.

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
> module, because the tab APIs are not yet in a published release.
> `go run .` on its own will not work: it skips binding generation and the
> frontend build.

## What to Expect

A single window opens on launch. It uses `MacWindowTabbingModePreferred`, so it
is willing to accept new tabs. The buttons in the window drive the tabbing
modes:

- **Open tabbed window** opens a window with `MacWindowTabbingModePreferred`. On
  macOS 10.12+ it merges into the current window as a new tab.
- **Open non-tabbed window** opens a window with `MacWindowTabbingModeDisallowed`.
  It always opens as a separate window and never tabs, even via Window > Merge
  All Windows.
- **Add tab via AddTab** opens a window with `MacWindowTabbingModeAutomatic`
  (which does not tab by itself under the default system setting) and attaches
  it with `WebviewWindow.AddTab(window, MacTabOrderAbove)`.
- **Detach this tab** calls `MoveTabToNewWindow` on the key window.
- **Describe tab group** prints `TabGroup().Count()`, `SelectedWindow()` and
  `IsTabBarVisible()` for the key window.

Open a mix of both to see the difference: tabbed windows stack into one titled
tab bar, while non-tabbed windows stay independent.

The **Tabs** menu applies the same API to the key window from Go:

| Menu item               | API                                                   |
| ----------------------- | ----------------------------------------------------- |
| Add Tab (Cmd+T)         | `WebviewWindow.AddTab(other, MacTabOrderAbove)`       |
| Add Native Window Tab   | `WebviewWindow.AddNativeTab(nativeWindow, order)`     |
| Select Next Tab         | `MacWindowTabGroup.SelectNext()`                      |
| Select Previous Tab     | `MacWindowTabGroup.SelectPrevious()`                  |
| Select First Tab        | `MacWindowTabGroup.Select(group.Windows()[0])`        |
| Toggle Tab Bar          | `MacWindowTabGroup.ToggleTabBar()`                    |
| Toggle Tab Overview     | `MacWindowTabGroup.ToggleTabOverview()`               |
| Rename Tab              | `WebviewWindow.SetTabTitle` and `SetTabTooltip`       |
| Move Tab to New Window  | `WebviewWindow.MoveTabToNewWindow()`                  |
| Merge All Windows       | `WebviewWindow.MergeAllWindows()`                     |
| Log Tab Groups          | `app.Window.TabGroups()` and the group getters        |

`WebviewWindow.TabGroup()` returns nil for windows created with tabbing
disallowed (the default) and before the native window exists, so tab
operations belong in menu handlers, service methods or window event callbacks
rather than before `app.Run()`. Every `MacWindowTabGroup` method is nil-safe.

## Relevant Code

See the Tabs menu in [main.go](main.go) and the tab helpers in
[windowservice.go](windowservice.go).
