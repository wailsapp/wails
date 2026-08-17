# Experimental native-content windows for Wails v3

## Status

This is an additive, intentionally temporary v3 API. It exists to validate a
WebView-free Wails application and collect performance data before the common
window interfaces are redesigned for v4.

It does not remove or weaken `WebviewWindow`. The existing WKWebView window can
still host `MacToolbar`, `MacSidebar`, `MacInspector`, and `MacSplitView`.

## Model

`NativeWindow` creates an ordinary `NSWindow`. Its split layout uses:

- `NSSplitViewController` for pane ownership and resizing;
- `NSOutlineView` for a semantic source-list sidebar;
- `NSTextView` in `NSScrollView` for plain-text editing;
- `NSToolbar` for sidebar, search, and document commands.

No WKWebView is allocated for this window.

```go
app := application.New(application.Options{
    NativeOnly: true,
    Mac: application.MacOptions{
        ActivationPolicy: application.ActivationPolicyAccessory,
    },
})

window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
    Title: "Native Notes",
    Width: 980,
    Height: 680,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropNormal,
        TitleBar: application.MacTitleBar{
            FullSizeContent: true,
            ToolbarStyle: application.MacToolbarStyleUnified,
        },
    },
})

sidebar := application.NewMacSidebar()
editor := application.NewMacTextEditor()
split := application.NewMacSplitView()
split.AddSidebar(sidebar)
split.AddTextEditor(editor)

if err := window.SetSplitView(split); err != nil {
    log.Fatal(err)
}
```

`AddPrimaryContent` remains the WKWebView primary-pane operation;
`AddTextEditor` is the native-window primary-pane operation. A layout uses one
or the other.

The same `MacToolbar` type attaches to either host. Toolbar and sidebar item IDs
remain generated internally; applications keep item handles and callbacks.

## NativeOnly

`Options.NativeOnly` skips frontend transport startup, asset-server creation,
and WebView-only event readers. It should be enabled only when the process will
not create a `WebviewWindow`.

Build with `wails_native` to apply the corresponding compile-time boundary:

```sh
GOWORK=off go build -tags 'production wails_native' \
  ./examples/mac-native-editor
```

That build excludes WebKit, frontend transport, HTTP assets, the updater, and
encrypted single-instance support. Add `wails_single_instance` alongside
`wails_native` if the latter is required. A normal build remains unchanged and
can still mix `WebviewWindow` with the native macOS chrome APIs.

## Text editor semantics

- The editor is plain text and editable by default.
- `SetText` is a programmatic update and does not invoke `OnChange`.
- `OnChange` carries no document copy. Call `Text()` only when the complete
  value is needed, such as Save.
- AppKit owns focus, selection, scrolling, accessibility, undo, and the native
  find bar.
- An editor and split layout may each be attached to only one live window.

## Performance

The controlled Swift/AppKit comparison, complete methodology, raw metric
definitions, results, limitations, and optimization findings are in
`performance.md`. The repeatable harness lives under
`examples/mac-native-editor/benchmark`.

## Deliberately deferred

- A v4 common `Window` interface and compatibility migration.
- A separately constructible WKWebView content component.
- `NSDocument`, autosave/versioning, attributed text, syntax highlighting, and
  multiple native content kinds.
- A final v4 package/build-feature design replacing the experimental tags.
