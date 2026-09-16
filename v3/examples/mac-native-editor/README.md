# Native Notes

Native Notes is an experimental, WebView-free Wails v3 application. It opens
its editor once on launch, then remains available from a macOS status item. The
close button hides the window so the tray icon can reopen it. The editor uses
only AppKit:

- `NSWindow` and `NSPanel`
- `NSToolbar` and `NSSearchToolbarItem`
- `NSSplitViewController`
- `NSOutlineView` source-list sidebar
- `NSScrollView` and `NSTextView`
- `NSVisualEffectView` behind the Quick Note panel

Run it from the v3 directory:

```sh
GOWORK=off go run -tags wails_native ./examples/mac-native-editor
```

The demo copies its three embedded text files into a temporary directory and
prints that directory on startup. Edits made with the toolbar or Command-S are
written to those files.

The Quick Note toolbar button (also View > Quick Note, Control-Command-N, and
the status item menu) opens a floating capture panel. Its inspector appends the
captured text to the note open in the main window, saves it to
`Quick Note.txt`, or hides the panel. The panel is a second `NativeWindow`
whose `Mac` options exercise the parts of `MacWindow` a document window has no
use for:

```go
Mac: application.MacWindow{
    WindowClass:        application.MacWindowClassPanel,
    PanelPreferences:   application.MacPanelPreferences{FloatingPanel: true},
    WindowLevel:        application.MacWindowLevelFloating,
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
        application.MacWindowCollectionBehaviorFullScreenAuxiliary,
    Backdrop:     application.MacBackdropTranslucent,
    TitleBar:     application.MacTitleBar{Hide: true},
    CornerRadius: 14,
}
```

## Which `MacWindow` options a `NativeWindow` applies

`NativeWindowOptions.Mac` is the same `MacWindow` struct a `WebviewWindow`
uses, and a native window applies each option with the same AppKit semantics.
`NativeWindowOptions` has no `Frameless` field, so `TitleBar.Hide` is the
native frameless switch and follows the `WebviewWindowOptions.Frameless` rules.

| Field | Native window behaviour |
| --- | --- |
| `Backdrop` | `Normal`: opaque window background. `Transparent`: clear window and transparent editor. `Translucent`: `NSVisualEffectView` beneath every pane. `LiquidGlass`: `NSGlassEffectView` on macOS 26+, translucent fallback elsewhere. Sidebar, inspector and content-list panes keep their own material. |
| `LiquidGlass` | `Style`, `Material`, `CornerRadius` and `TintColor` apply; `GroupID` and `GroupSpacing` need `-tags private_mac_apis`, as for a WebView window. |
| `DisableShadow` | `NSWindow.hasShadow`. |
| `TitleBar.Hide` | Frameless. With `CornerType` rounded and `CornerRadius` 0 the AppKit frame is kept with a transparent titlebar, hidden title and hidden window buttons; otherwise the window is borderless. The other `TitleBar` fields are then ignored. |
| `CornerType`, `CornerRadius` | Only with `TitleBar.Hide`: square gives a borderless square window; a radius masks the content to that radius. |
| `TitleBar.AppearsTransparent`, `HideTitle`, `FullSizeContent`, `ToolbarStyle`, `HideToolbarSeparator`, `ShowToolbarWhenFullscreen` | Applied to the `NSWindow` and its attached `NSToolbar`. `UseToolbar` is ignored: a native window only shows a toolbar attached with `SetToolbar`. |
| `ContentLayout` | `EdgeToEdge` keeps the split layout full-size so the editor scrolls beneath the toolbar. `BelowToolbar` removes full-size content, which places the whole layout, sidebar included, beneath the toolbar. `Automatic` follows `TitleBar.FullSizeContent` (edge-to-edge when frameless). |
| `Appearance` | `NSWindow.appearance` by name. |
| `WindowLevel` | `NSWindow.level`. Precedence: `WindowLevel`, then `AlwaysOnTop` or a floating panel, then normal. |
| `CollectionBehavior` | `NSWindow.collectionBehavior`; zero selects `FullScreenPrimary`. |
| `TabbingMode` | `NSWindow.tabbingMode`; the default resolves to disallowed. |
| `DisableEscapeExitsFullscreen` | Escape is swallowed while fullscreen. |
| `WindowClass`, `PanelPreferences` | `Panel` creates an `NSPanel` with the floating, becomes-key-only-if-needed, non-activating and utility preferences. |
| `InvisibleTitleBarHeight` | Not applicable: it sizes the WebView drag region. |
| `EventMapping` | Not applicable: a native window emits no window events yet. |
| `EnableFraudulentWebsiteWarnings`, `WebviewPreferences` | Not applicable: there is no WKWebView. |
| `SplitView`, `Toolbar` | Ignored: use `NativeWindowOptions.SplitView` and `NativeWindowOptions.Toolbar`. |

For a quick process-level sample, find the compiled process and inspect it over
several seconds:

```sh
pid=$(pgrep -n mac-native-editor)
ps -o pid,rss,%cpu,time,command -p "$pid"
top -l 5 -pid "$pid" -stats pid,command,cpu,mem,threads,wq
```

There should be no WebKit WebContent child process associated with Native
Notes. Compare the same measurements with `examples/mac-toolbar` to quantify
the cost of the WKWebView path on the same machine and OS release.

Use the `wails_native` build tag for a genuinely native binary. `NativeOnly`
without the tag suppresses frontend startup at runtime, but the monolithic
application package still links WebKit and includes the frontend and updater
dependency graphs.

Single-instance support is also omitted from the lean native build. Add
`wails_single_instance` when it is required:

```sh
GOWORK=off go run -tags 'wails_native wails_single_instance' \
  ./examples/mac-native-editor
```

For the controlled Swift/AppKit comparison, build both release binaries and
run the complete launch, idle, hidden-window, and document-scaling matrix:

```sh
examples/mac-native-editor/benchmark/run.sh
```

The script prints its results directory. `LAUNCH_ITERATIONS` and `IDLE_SECONDS`
can reduce or extend the default 20 launches and 30-second idle observations.
The checked methodology and August 2026 native/Swift result are documented in
`wep/proposals/macos-native-window/performance.md`. The wider investigation of
normal Wails builds, package reachability, WebKit process cost, GC behaviour,
and optimization priorities is in
`wep/proposals/macos-native-window/performance-investigation.md`.
